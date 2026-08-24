package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ArtemStepanov/caddy-admin-ui/internal/caddy"
	"github.com/ArtemStepanov/caddy-admin-ui/internal/config"
	"github.com/ArtemStepanov/caddy-admin-ui/internal/storage"
)

var (
	errSetupRequired = errors.New("finish Caddy setup before syncing routes")
	errConfigDrift   = errors.New("Caddy routes changed outside Caddy Admin UI; review and import them before retrying")
	errInvalidRoutes = errors.New("route validation failed")
	serverNameRE     = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)
)

type setupRequest struct {
	URL          string `json:"url"`
	Server       string `json:"server"`
	PreviewToken string `json:"preview_token"`
}

type setupPreview struct {
	URL             string           `json:"url"`
	Servers         []string         `json:"servers"`
	SelectedServer  string           `json:"selected_server"`
	Routes          []*storage.Route `json:"routes"`
	LocalDrafts     []*storage.Route `json:"local_drafts"`
	RouteCount      int              `json:"route_count"`
	Editable        int              `json:"editable"`
	ReadOnly        int              `json:"readonly"`
	CaddyEmpty      bool             `json:"caddy_empty"`
	CanBootstrap    bool             `json:"can_bootstrap"`
	PreviewToken    string           `json:"preview_token"`
	OwnershipNotice string           `json:"ownership_notice"`
}

type caddyInventory struct {
	Root          json.RawMessage
	RootETag      string
	Empty         bool
	Servers       []string
	RouteArray    map[string]json.RawMessage
	RoutesPresent map[string]bool
}

func managedRoutesPath(server string) string {
	return "apps/http/servers/" + server + "/routes"
}

func validateServerName(server string) error {
	if !serverNameRE.MatchString(server) {
		return fmt.Errorf("server name may contain only letters, digits, dot, underscore, and hyphen")
	}
	return nil
}

func fetchInventory(client *caddy.Client) (*caddyInventory, error) {
	response, err := client.GetConfigWithETag("")
	if err != nil {
		return nil, err
	}
	inventory := &caddyInventory{
		Root:          append(json.RawMessage(nil), response.Body...),
		RootETag:      response.ETag,
		RouteArray:    make(map[string]json.RawMessage),
		RoutesPresent: make(map[string]bool),
	}
	var root map[string]json.RawMessage
	if err := json.Unmarshal(response.Body, &root); err != nil {
		return nil, fmt.Errorf("parse Caddy config: %w", err)
	}
	inventory.Empty = len(root) == 0
	var apps map[string]json.RawMessage
	if raw := root["apps"]; len(raw) > 0 {
		if err := json.Unmarshal(raw, &apps); err != nil {
			return nil, fmt.Errorf("parse Caddy apps: %w", err)
		}
	}
	var httpApp map[string]json.RawMessage
	if raw := apps["http"]; len(raw) > 0 {
		if err := json.Unmarshal(raw, &httpApp); err != nil {
			return nil, fmt.Errorf("parse Caddy HTTP app: %w", err)
		}
	}
	var servers map[string]json.RawMessage
	if raw := httpApp["servers"]; len(raw) > 0 {
		if err := json.Unmarshal(raw, &servers); err != nil {
			return nil, fmt.Errorf("parse Caddy HTTP servers: %w", err)
		}
	}
	for name, rawServer := range servers {
		inventory.Servers = append(inventory.Servers, name)
		var server map[string]json.RawMessage
		if err := json.Unmarshal(rawServer, &server); err != nil {
			return nil, fmt.Errorf("parse Caddy server %q: %w", name, err)
		}
		routes, routesPresent := server["routes"]
		if len(routes) == 0 || string(routes) == "null" {
			routes = json.RawMessage("[]")
			routesPresent = false
		}
		inventory.RouteArray[name] = routes
		inventory.RoutesPresent[name] = routesPresent
	}
	sort.Strings(inventory.Servers)
	return inventory, nil
}

func setupToken(url, server string, inventory *caddyInventory) string {
	hash := sha256.New()
	hash.Write([]byte(strings.TrimRight(url, "/")))
	hash.Write([]byte{0})
	hash.Write([]byte(server))
	hash.Write([]byte{0})
	hash.Write([]byte(inventory.RootETag))
	hash.Write([]byte{0})
	hash.Write(inventory.Root)
	return hex.EncodeToString(hash.Sum(nil))
}

func (h *Handler) buildSetupPreview(req setupRequest) (*setupPreview, *caddyInventory, error) {
	url := strings.TrimRight(req.URL, "/")
	if url == "" {
		url = h.getCaddyURL()
	}
	if err := validateCaddyURL(url); err != nil {
		return nil, nil, err
	}
	inventory, err := fetchInventory(caddy.NewClient(url))
	if err != nil {
		return nil, nil, err
	}
	server := req.Server
	if server == "" {
		cfg, _ := h.store.GetGlobalConfig()
		if cfg != nil {
			server = cfg.ManagedServer
		}
	}
	if server == "" && len(inventory.Servers) > 0 {
		server = inventory.Servers[0]
	}
	if server == "" && inventory.Empty {
		server = "caddy-admin-ui"
	}
	if server == "" {
		return nil, nil, fmt.Errorf("Caddy has no HTTP servers and is not empty; create a server before setup")
	}
	if err := validateServerName(server); err != nil {
		return nil, nil, err
	}
	rawRoutes, exists := inventory.RouteArray[server]
	if !exists && !inventory.Empty {
		return nil, nil, fmt.Errorf("HTTP server %q does not exist; select an existing server", server)
	}
	if len(rawRoutes) == 0 {
		rawRoutes = json.RawMessage("[]")
	}
	routes, err := config.ParseCaddyRoutes(rawRoutes)
	if err != nil {
		return nil, nil, err
	}
	localRoutes, err := h.store.ListRoutes()
	if err != nil {
		return nil, nil, fmt.Errorf("read local route drafts: %w", err)
	}
	localDrafts := make([]*storage.Route, 0)
	remoteIDs := make(map[string]bool, len(routes))
	for _, route := range routes {
		remoteIDs[route.ID] = true
	}
	for _, route := range localRoutes {
		if !route.IsReadOnly() && !remoteIDs[route.ID] {
			localDrafts = append(localDrafts, route)
		}
	}
	preview := &setupPreview{
		URL:             url,
		Servers:         inventory.Servers,
		SelectedServer:  server,
		Routes:          routes,
		LocalDrafts:     localDrafts,
		RouteCount:      len(routes),
		CaddyEmpty:      inventory.Empty,
		CanBootstrap:    inventory.Empty,
		PreviewToken:    setupToken(url, server, inventory),
		OwnershipNotice: "Caddy Admin UI will manage only this server's routes array. Other Caddy apps and server settings remain untouched.",
	}
	for _, route := range routes {
		if route.IsReadOnly() {
			preview.ReadOnly++
		} else {
			preview.Editable++
		}
	}
	return preview, inventory, nil
}

// PreviewSetup inventories Caddy and returns a content-derived confirmation token.
func (h *Handler) PreviewSetup(c *gin.Context) {
	var req setupRequest
	if c.Request.ContentLength != 0 {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	preview, _, err := h.buildSetupPreview(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, preview)
}

// ConfirmSetup re-fetches Caddy, verifies the preview token, then adopts one route array.
func (h *Handler) ConfirmSetup(c *gin.Context) {
	var req setupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.PreviewToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "preview_token is required"})
		return
	}

	h.syncMu.Lock()
	defer h.syncMu.Unlock()
	preview, inventory, err := h.buildSetupPreview(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	if preview.PreviewToken != req.PreviewToken {
		c.JSON(http.StatusConflict, gin.H{"error": "Caddy changed after preview; review setup again"})
		return
	}
	client := caddy.NewClient(preview.URL)
	if inventory.Empty {
		apps := map[string]any{"http": map[string]any{"servers": map[string]any{
			preview.SelectedServer: map[string]any{"listen": []string{":80", ":443"}, "routes": []any{}},
		}}}
		if inventory.RootETag == "" {
			c.JSON(http.StatusBadGateway, gin.H{"error": "Caddy root response did not include an ETag"})
			return
		}
		if err := client.PutConfig("apps", apps, inventory.RootETag); err != nil {
			var apiErr *caddy.APIError
			if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusPreconditionFailed {
				c.JSON(http.StatusConflict, gin.H{"error": "Caddy changed during setup; review setup again"})
				return
			}
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to bootstrap empty Caddy: " + err.Error()})
			return
		}
	} else if !inventory.RoutesPresent[preview.SelectedServer] {
		serverPath := "apps/http/servers/" + preview.SelectedServer
		serverConfig, err := client.GetConfigWithETag(serverPath)
		if err != nil || serverConfig.ETag == "" {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to initialize the selected server's routes array"})
			return
		}
		var selectedServer map[string]any
		if err := json.Unmarshal(serverConfig.Body, &selectedServer); err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to parse the selected Caddy server"})
			return
		}
		selectedServer["routes"] = []any{}
		if err := client.PatchConfig(serverPath, selectedServer, serverConfig.ETag); err != nil {
			var apiErr *caddy.APIError
			if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusPreconditionFailed {
				c.JSON(http.StatusConflict, gin.H{"error": "Caddy changed during setup; review setup again"})
				return
			}
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to initialize the selected server's routes array"})
			return
		}
	}
	remote, err := client.GetConfigWithETag(managedRoutesPath(preview.SelectedServer))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to read selected routes: " + err.Error()})
		return
	}
	if remote.ETag == "" {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Caddy route response did not include an ETag"})
		return
	}
	routes, err := config.ParseCaddyRoutes(remote.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	importedCount := len(routes)
	retainedDrafts := 0
	localRoutes, err := h.store.ListRoutes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read local route drafts"})
		return
	}
	remoteIDs := make(map[string]bool, len(routes))
	for _, route := range routes {
		remoteIDs[route.ID] = true
	}
	for _, route := range localRoutes {
		if !route.IsReadOnly() && !remoteIDs[route.ID] {
			route.Position = len(routes)
			routes = append(routes, route)
			retainedDrafts++
		}
	}
	if err := h.store.CreateSnapshot(&storage.Snapshot{Server: preview.SelectedServer, ETag: remote.ETag, Routes: remote.Body, Reason: "initial setup"}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save initial snapshot"})
		return
	}
	if err := h.store.ReplaceAllRoutes(routes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save imported routes"})
		return
	}
	cfg, err := h.store.GetGlobalConfig()
	if err != nil {
		_ = h.store.ReplaceAllRoutes(localRoutes)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	cfg.CaddyAdminURL = preview.URL
	cfg.ManagedServer = preview.SelectedServer
	cfg.SetupComplete = true
	cfg.LastETag = remote.ETag
	if err := h.store.SetGlobalConfig(cfg); err != nil {
		_ = h.store.ReplaceAllRoutes(localRoutes)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"config": cfg, "imported": importedCount, "local_drafts": retainedDrafts, "message": "Caddy ownership confirmed"})
}

func cloneRoutes(routes []*storage.Route) ([]*storage.Route, error) {
	data, err := json.Marshal(routes)
	if err != nil {
		return nil, err
	}
	var cloned []*storage.Route
	if err := json.Unmarshal(data, &cloned); err != nil {
		return nil, err
	}
	for i, route := range cloned {
		route.RawCaddyRoute = append(json.RawMessage(nil), routes[i].RawCaddyRoute...)
	}
	return cloned, nil
}

func (h *Handler) mutateRoutes(reason string, mutate func([]*storage.Route) ([]*storage.Route, error)) error {
	h.syncMu.Lock()
	defer h.syncMu.Unlock()
	current, err := h.store.ListRoutes()
	if err != nil {
		return err
	}
	candidate, err := cloneRoutes(current)
	if err != nil {
		return err
	}
	candidate, err = mutate(candidate)
	if err != nil {
		return err
	}
	for position, route := range candidate {
		route.Position = position
	}
	if err := config.ValidateRoutesForBuild(candidate); err != nil {
		return fmt.Errorf("%w: %v", errInvalidRoutes, err)
	}
	cfg, err := h.store.GetGlobalConfig()
	if err != nil {
		return err
	}
	if !cfg.SetupComplete {
		if err := h.store.ReplaceAllRoutes(candidate); err != nil {
			return err
		}
		return errSetupRequired
	}
	if err := h.syncRoutesLocked(candidate, reason, cfg); err != nil {
		return err
	}
	if err := h.store.ReplaceAllRoutes(candidate); err != nil {
		latestCfg, cfgErr := h.store.GetGlobalConfig()
		if cfgErr == nil {
			if rollbackErr := h.syncRoutesLocked(current, "automatic rollback after local persistence failure", latestCfg); rollbackErr != nil {
				return fmt.Errorf("save local routes: %w; Caddy rollback also failed: %v", err, rollbackErr)
			}
		}
		return fmt.Errorf("save local routes after Caddy update: %w", err)
	}
	return nil
}

func (h *Handler) syncRoutesLocked(routes []*storage.Route, reason string, cfg *storage.GlobalConfig) error {
	if cfg == nil || !cfg.SetupComplete {
		return errSetupRequired
	}
	if err := validateServerName(cfg.ManagedServer); err != nil {
		return err
	}
	client := caddy.NewClient(h.getCaddyURL())
	path := managedRoutesPath(cfg.ManagedServer)
	remote, err := client.GetConfigWithETag(path)
	if err != nil {
		h.setSyncStatus(err)
		return err
	}
	if remote.ETag == "" {
		err := fmt.Errorf("Caddy route response did not include an ETag")
		h.setSyncStatus(err)
		return err
	}
	if cfg.LastETag != "" && remote.ETag != cfg.LastETag {
		h.setSyncStatus(errConfigDrift)
		return errConfigDrift
	}
	if err := h.store.CreateSnapshot(&storage.Snapshot{Server: cfg.ManagedServer, ETag: remote.ETag, Routes: remote.Body, Reason: reason}); err != nil {
		return fmt.Errorf("save pre-write snapshot: %w", err)
	}
	built, err := config.BuildCaddyRoutes(routes, cfg)
	if err != nil {
		return err
	}
	if err := client.PatchConfig(path, built, remote.ETag); err != nil {
		var apiErr *caddy.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusPreconditionFailed {
			err = errConfigDrift
		}
		h.setSyncStatus(err)
		return err
	}
	after, err := client.GetConfigWithETag(path)
	if err != nil {
		h.setSyncStatus(err)
		return fmt.Errorf("routes applied but new ETag could not be read: %w", err)
	}
	if after.ETag == "" {
		err := fmt.Errorf("routes applied but Caddy did not return a new ETag")
		h.setSyncStatus(err)
		return err
	}
	cfg.LastETag = after.ETag
	if err := h.store.SetGlobalConfig(cfg); err != nil {
		var previous []json.RawMessage
		if json.Unmarshal(remote.Body, &previous) == nil {
			_ = client.PatchConfig(path, previous, after.ETag)
		}
		return fmt.Errorf("save new Caddy ETag: %w", err)
	}
	h.setSyncStatus(nil)
	return nil
}

func (h *Handler) setSyncStatus(err error) {
	h.statusMu.Lock()
	defer h.statusMu.Unlock()
	h.lastSyncedAt = time.Now()
	if err == nil {
		h.lastSyncError = ""
	} else {
		h.lastSyncError = err.Error()
	}
}

func (h *Handler) getSyncStatus() (time.Time, string) {
	h.statusMu.RLock()
	defer h.statusMu.RUnlock()
	return h.lastSyncedAt, h.lastSyncError
}

// ListSnapshots returns recent pre-write restore points.
func (h *Handler) ListSnapshots(c *gin.Context) {
	snapshots, err := h.store.ListSnapshots()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if snapshots == nil {
		snapshots = []*storage.Snapshot{}
	}
	c.JSON(http.StatusOK, gin.H{"snapshots": snapshots})
}

// ExportSnapshot downloads a snapshot as a standalone JSON route array.
func (h *Handler) ExportSnapshot(c *gin.Context) {
	snapshot, err := h.store.GetSnapshot(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "snapshot not found"})
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="caddy-routes-%s.json"`, snapshot.ID))
	c.Data(http.StatusOK, "application/json; charset=utf-8", snapshot.Routes)
}

// RestoreSnapshot applies a prior route array with the same drift protection as normal writes.
func (h *Handler) RestoreSnapshot(c *gin.Context) {
	h.syncMu.Lock()
	defer h.syncMu.Unlock()
	snapshot, err := h.store.GetSnapshot(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "snapshot not found"})
		return
	}
	cfg, err := h.store.GetGlobalConfig()
	if err != nil || !cfg.SetupComplete {
		c.JSON(http.StatusConflict, gin.H{"error": errSetupRequired.Error()})
		return
	}
	if snapshot.Server != cfg.ManagedServer {
		c.JSON(http.StatusConflict, gin.H{"error": "snapshot belongs to a different managed server"})
		return
	}
	client := caddy.NewClient(h.getCaddyURL())
	path := managedRoutesPath(cfg.ManagedServer)
	remote, err := client.GetConfigWithETag(path)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	if remote.ETag == "" {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Caddy route response did not include an ETag"})
		return
	}
	if cfg.LastETag != "" && remote.ETag != cfg.LastETag {
		c.JSON(http.StatusConflict, gin.H{"error": errConfigDrift.Error()})
		return
	}
	var restored []json.RawMessage
	if err := json.Unmarshal(snapshot.Routes, &restored); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "snapshot route data is invalid"})
		return
	}
	if err := h.store.CreateSnapshot(&storage.Snapshot{Server: cfg.ManagedServer, ETag: remote.ETag, Routes: remote.Body, Reason: "before snapshot restore"}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := client.PatchConfig(path, restored, remote.ETag); err != nil {
		var apiErr *caddy.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusPreconditionFailed {
			c.JSON(http.StatusConflict, gin.H{"error": errConfigDrift.Error()})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	after, err := client.GetConfigWithETag(path)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	if after.ETag == "" {
		c.JSON(http.StatusBadGateway, gin.H{"error": "snapshot applied but Caddy did not return a new ETag"})
		return
	}
	routes, err := config.ParseCaddyRoutes(after.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := h.store.ReplaceAllRoutes(routes); err != nil {
		var previous []json.RawMessage
		if json.Unmarshal(remote.Body, &previous) == nil {
			_ = client.PatchConfig(path, previous, after.ETag)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	cfg.LastETag = after.ETag
	if err := h.store.SetGlobalConfig(cfg); err != nil {
		var previous []json.RawMessage
		if json.Unmarshal(remote.Body, &previous) == nil {
			_ = client.PatchConfig(path, previous, after.ETag)
			if previousRoutes, parseErr := config.ParseCaddyRoutes(remote.Body); parseErr == nil {
				_ = h.store.ReplaceAllRoutes(previousRoutes)
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.setSyncStatus(nil)
	c.JSON(http.StatusOK, gin.H{"message": "snapshot restored", "restored": len(routes)})
}
