package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ArtemStepanov/caddy-admin-ui/internal/caddy"
	"github.com/ArtemStepanov/caddy-admin-ui/internal/config"
	"github.com/ArtemStepanov/caddy-admin-ui/internal/storage"
	"github.com/ArtemStepanov/caddy-admin-ui/internal/version"
)

// validateCaddyURL checks that a URL is a valid http(s) origin (scheme + host + optional port).
func validateCaddyURL(rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("URL must not be empty")
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("URL scheme must be http or https, got %q", u.Scheme)
	}
	if u.Host == "" {
		return fmt.Errorf("URL must have a host")
	}
	if u.Path != "" && u.Path != "/" {
		return fmt.Errorf("URL must not contain a path")
	}
	if u.RawQuery != "" {
		return fmt.Errorf("URL must not contain query parameters")
	}
	if u.Fragment != "" {
		return fmt.Errorf("URL must not contain a fragment")
	}
	if u.User != nil {
		return fmt.Errorf("URL must not contain credentials")
	}
	return nil
}

const (
	driftWarning     = "Manual Caddy changes after the last import or sync are not automatically merged. Re-run import review before syncing after manual edits."
	recoveryGuidance = "Local routes are unchanged. Fix Caddy connectivity or re-run import review before syncing again."
)

type importSummary struct {
	TotalFound        int  `json:"total_found"`
	Editable          int  `json:"editable"`
	ReadOnlyPreserved int  `json:"readonly_preserved"`
	Unsupported       int  `json:"unsupported"`
	LocalOnly         int  `json:"local_only"`
	WillUpdate        int  `json:"will_update"`
	WillReplaceLocal  bool `json:"will_replace_local"`
}

type importGroups struct {
	NewFromCaddy      []importRouteRow `json:"new_from_caddy"`
	WillUpdate        []importRouteRow `json:"will_update"`
	LocalOnly         []importRouteRow `json:"local_only"`
	ReadOnlyPreserved []importRouteRow `json:"readonly_preserved"`
}

type importRouteRow struct {
	RouteID        string          `json:"route_id,omitempty"`
	Domain         string          `json:"domain"`
	Path           string          `json:"path,omitempty"`
	HandlerType    string          `json:"handler_type"`
	Destination    string          `json:"destination,omitempty"`
	SupportStatus  string          `json:"support_status"`
	ReadOnlyReason string          `json:"readonly_reason,omitempty"`
	RawCaddyRoute  json.RawMessage `json:"raw_caddy_route,omitempty"`
	ChangeType     string          `json:"change_type"`
}

// Handler contains all HTTP handlers
type Handler struct {
	store           *storage.SQLiteStorage
	defaultCaddyURL string // fallback URL from env
	lastSyncedAt    time.Time
	lastSyncError   string
}

// NewHandler creates a new handler
func NewHandler(store *storage.SQLiteStorage, defaultCaddyURL string) *Handler {
	return &Handler{
		store:           store,
		defaultCaddyURL: defaultCaddyURL,
	}
}

// getCaddyClient returns a Caddy client using the URL from GlobalConfig (or default)
func (h *Handler) getCaddyClient() *caddy.Client {
	cfg, err := h.store.GetGlobalConfig()
	if err != nil || cfg.CaddyAdminURL == "" {
		return caddy.NewClient(h.defaultCaddyURL)
	}
	return caddy.NewClient(cfg.CaddyAdminURL)
}

// getCaddyURL returns the current Caddy URL from GlobalConfig (or default)
func (h *Handler) getCaddyURL() string {
	cfg, err := h.store.GetGlobalConfig()
	if err != nil || cfg.CaddyAdminURL == "" {
		return h.defaultCaddyURL
	}
	return cfg.CaddyAdminURL
}

func makeEditableRoute(route *storage.Route) {
	route.ReadOnly = false
	route.SupportStatus = storage.SupportStatusEditable
	route.ReadOnlyReason = ""
	route.RawCaddyRoute = nil
}

func readonlyConflict(c *gin.Context) {
	c.JSON(http.StatusConflict, gin.H{"error": "route is read-only and managed outside the UI"})
}

func syncWarning(message string, err error) gin.H {
	return gin.H{
		"status":            "warning",
		"warning":           message + ": " + err.Error(),
		"recovery_guidance": recoveryGuidance,
	}
}

func routeKey(route *storage.Route) string {
	return route.Domain + "\x00" + route.Path + "\x00" + route.HandlerType
}

func routeDestination(route *storage.Route) string {
	switch route.HandlerType {
	case "reverse_proxy":
		var cfg storage.ReverseProxyConfig
		if json.Unmarshal(route.Config, &cfg) == nil && len(cfg.Upstreams) > 0 {
			return cfg.Upstreams[0]
		}
	case "file_server":
		var cfg storage.FileServerConfig
		if json.Unmarshal(route.Config, &cfg) == nil {
			return cfg.Root
		}
	case "redir":
		var cfg storage.RedirectConfig
		if json.Unmarshal(route.Config, &cfg) == nil {
			return cfg.To
		}
	}
	return ""
}

func routeDisplayDomain(route *storage.Route, destination string) string {
	if route.Domain != "" && route.Domain != "*" && route.Domain != "UNKNOWN" {
		return route.Domain
	}
	if route.Path != "" {
		return route.Path
	}
	if destination != "" {
		return destination
	}
	return "unknown route"
}

func routeRow(route *storage.Route, changeType string) importRouteRow {
	destination := routeDestination(route)
	row := importRouteRow{
		RouteID:        route.ID,
		Domain:         routeDisplayDomain(route, destination),
		Path:           route.Path,
		HandlerType:    route.HandlerType,
		Destination:    destination,
		SupportStatus:  route.SupportStatus,
		ReadOnlyReason: route.ReadOnlyReason,
		ChangeType:     changeType,
	}
	if route.IsReadOnly() {
		row.RawCaddyRoute = route.RawCaddyRoute
	}
	return row
}

func preserveMatchingLocalIDs(routes, localRoutes []*storage.Route) {
	locals := map[string]*storage.Route{}
	for _, route := range localRoutes {
		locals[routeKey(route)] = route
	}
	for _, route := range routes {
		key := routeKey(route)
		if local := locals[key]; local != nil {
			route.ID = local.ID
			route.CreatedAt = local.CreatedAt
			delete(locals, key)
		}
	}
}

func buildImportPreview(routes, localRoutes []*storage.Route) (importSummary, importGroups) {
	locals := map[string]*storage.Route{}
	seen := map[string]bool{}
	for _, route := range localRoutes {
		locals[routeKey(route)] = route
	}

	summary := importSummary{TotalFound: len(routes), WillReplaceLocal: true}
	groups := importGroups{
		NewFromCaddy:      []importRouteRow{},
		WillUpdate:        []importRouteRow{},
		LocalOnly:         []importRouteRow{},
		ReadOnlyPreserved: []importRouteRow{},
	}
	for _, route := range routes {
		key := routeKey(route)
		seen[key] = true
		if route.SupportStatus == storage.SupportStatusEditable {
			summary.Editable++
			row := routeRow(route, "new")
			if local := locals[key]; local != nil {
				row.RouteID = local.ID
				row.ChangeType = "update"
				summary.WillUpdate++
				groups.WillUpdate = append(groups.WillUpdate, row)
			} else {
				groups.NewFromCaddy = append(groups.NewFromCaddy, row)
			}
			continue
		}

		summary.ReadOnlyPreserved++
		if route.SupportStatus == storage.SupportStatusUnsupportedReadOnly {
			summary.Unsupported++
		}
		row := routeRow(route, "readonly_preserve")
		if local := locals[key]; local != nil {
			row.RouteID = local.ID
		}
		groups.ReadOnlyPreserved = append(groups.ReadOnlyPreserved, row)
	}

	for _, route := range localRoutes {
		if !seen[routeKey(route)] {
			summary.LocalOnly++
			groups.LocalOnly = append(groups.LocalOnly, routeRow(route, "local_only_remove"))
		}
	}
	return summary, groups
}

func fetchParsedCaddyRoutes(h *Handler) ([]*storage.Route, int, error) {
	raw, err := h.getCaddyClient().GetConfig("")
	if err != nil {
		return nil, http.StatusBadGateway, fmt.Errorf("failed to connect to Caddy")
	}

	var caddyConfig config.CaddyConfig
	if err := json.Unmarshal(raw, &caddyConfig); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to parse Caddy response")
	}

	routes, err := config.ParseCaddyConfig(&caddyConfig)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to parse Caddy config")
	}
	if err := config.ValidateRoutesForBuild(routes); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to validate Caddy routes: %w", err)
	}
	return routes, http.StatusOK, nil
}

// ListRoutes returns all routes
func (h *Handler) ListRoutes(c *gin.Context) {
	routes, err := h.store.ListRoutes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if routes == nil {
		routes = []*storage.Route{}
	}
	c.JSON(http.StatusOK, gin.H{"routes": routes})
}

// CreateRoute creates a new route
func (h *Handler) CreateRoute(c *gin.Context) {
	var route storage.Route
	if err := c.ShouldBindJSON(&route); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if route.Domain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "domain is required"})
		return
	}
	if route.HandlerType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "handler_type is required"})
		return
	}

	route.Enabled = true // New routes are enabled by default
	makeEditableRoute(&route)

	if err := h.store.CreateRoute(&route); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Auto-sync to Caddy
	if err := h.syncToCaddy(); err != nil {
		resp := syncWarning("Route created but sync to Caddy failed", err)
		resp["route"] = route
		c.JSON(http.StatusCreated, resp)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"route": route})
}

// GetRoute returns a single route
func (h *Handler) GetRoute(c *gin.Context) {
	id := c.Param("id")
	route, err := h.store.GetRoute(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"route": route})
}

// UpdateRoute updates an existing route
func (h *Handler) UpdateRoute(c *gin.Context) {
	id := c.Param("id")

	// Check if route exists
	existing, err := h.store.GetRoute(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
		return
	}
	if existing.IsReadOnly() {
		readonlyConflict(c)
		return
	}

	var route storage.Route
	if err := c.ShouldBindJSON(&route); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Preserve ID and timestamps
	route.ID = existing.ID
	route.CreatedAt = existing.CreatedAt
	// Preserve state that isn't editable in the UI
	route.Enabled = existing.Enabled
	makeEditableRoute(&route)

	if err := h.store.UpdateRoute(&route); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Auto-sync to Caddy
	if err := h.syncToCaddy(); err != nil {
		resp := syncWarning("Route updated but sync to Caddy failed", err)
		resp["route"] = route
		c.JSON(http.StatusOK, resp)
		return
	}

	c.JSON(http.StatusOK, gin.H{"route": route})
}

// DeleteRoute deletes a route
func (h *Handler) DeleteRoute(c *gin.Context) {
	id := c.Param("id")

	existing, err := h.store.GetRoute(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
		return
	}
	if existing.IsReadOnly() {
		readonlyConflict(c)
		return
	}

	if err := h.store.DeleteRoute(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Auto-sync to Caddy
	if err := h.syncToCaddy(); err != nil {
		resp := syncWarning("Route deleted but sync to Caddy failed", err)
		resp["message"] = resp["warning"]
		c.JSON(http.StatusOK, resp)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "route deleted"})
}

// ToggleRoute enables/disables a route
func (h *Handler) ToggleRoute(c *gin.Context) {
	id := c.Param("id")

	route, err := h.store.GetRoute(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
		return
	}
	if route.IsReadOnly() {
		readonlyConflict(c)
		return
	}

	route.Enabled = !route.Enabled

	if err := h.store.UpdateRoute(route); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Auto-sync to Caddy
	if err := h.syncToCaddy(); err != nil {
		resp := syncWarning("Route toggled but sync to Caddy failed", err)
		resp["route"] = route
		c.JSON(http.StatusOK, resp)
		return
	}

	c.JSON(http.StatusOK, gin.H{"route": route})
}

// GetConfig returns global configuration
func (h *Handler) GetConfig(c *gin.Context) {
	cfg, err := h.store.GetGlobalConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if cfg.CaddyAdminURL == "" {
		cfg.CaddyAdminURL = h.defaultCaddyURL
	}
	c.JSON(http.StatusOK, gin.H{"config": cfg})
}

// UpdateConfig updates global configuration
func (h *Handler) UpdateConfig(c *gin.Context) {
	var cfg storage.GlobalConfig
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if cfg.CaddyAdminURL != "" {
		if err := validateCaddyURL(cfg.CaddyAdminURL); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if err := h.store.SetGlobalConfig(&cfg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"config": cfg})
}

// GetStatus returns Caddy health status
func (h *Handler) GetStatus(c *gin.Context) {
	caddyClient := h.getCaddyClient()
	caddyURL := h.getCaddyURL()

	start := time.Now()
	err := caddyClient.Health()
	latency := time.Since(start).Milliseconds()

	if err != nil {
		resp := gin.H{
			"status":    "offline",
			"error":     err.Error(),
			"latency":   latency,
			"admin_url": caddyURL,
			"version":   version.Version,
		}
		if !h.lastSyncedAt.IsZero() {
			resp["last_synced_at"] = h.lastSyncedAt.Format(time.RFC3339)
			resp["last_sync_error"] = h.lastSyncError
		}
		c.JSON(http.StatusOK, resp)
		return
	}

	// Get route count
	routes, _ := h.store.ListRoutes()
	routeCount := 0
	if routes != nil {
		routeCount = len(routes)
	}

	resp := gin.H{
		"status":      "online",
		"latency":     latency,
		"admin_url":   caddyURL,
		"route_count": routeCount,
		"version":     version.Version,
	}
	if !h.lastSyncedAt.IsZero() {
		resp["last_synced_at"] = h.lastSyncedAt.Format(time.RFC3339)
		resp["last_sync_error"] = h.lastSyncError
	}
	c.JSON(http.StatusOK, resp)
}

// SyncToCaddy manually triggers sync to Caddy
func (h *Handler) SyncToCaddy(c *gin.Context) {
	if err := h.syncToCaddy(); err != nil {
		resp := syncWarning("Sync to Caddy failed", err)
		resp["error"] = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "synced successfully"})
}

// TestConnection tests connection to a specific Caddy URL
func (h *Handler) TestConnection(c *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validateCaddyURL(req.URL); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	testClient := caddy.NewClient(req.URL)
	start := time.Now()
	err := testClient.Health()
	latency := time.Since(start).Milliseconds()

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"error":   err.Error(),
			"latency": latency,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"latency": latency,
	})
}

// syncToCaddy builds config from routes and loads it into Caddy
func (h *Handler) syncToCaddy() error {
	routes, err := h.store.ListRoutes()
	if err != nil {
		h.lastSyncedAt = time.Now()
		h.lastSyncError = err.Error()
		return err
	}

	globalCfg, err := h.store.GetGlobalConfig()
	if err != nil {
		h.lastSyncedAt = time.Now()
		h.lastSyncError = err.Error()
		return err
	}

	if err := config.ValidateRoutesForBuild(routes); err != nil {
		h.lastSyncedAt = time.Now()
		h.lastSyncError = err.Error()
		return err
	}

	// Build Caddy config from routes
	caddyConfig := config.BuildCaddyConfig(routes, globalCfg)

	// Pretty print for debugging
	data, _ := json.MarshalIndent(caddyConfig, "", "  ")
	_ = data // Could log this if needed

	// Load into Caddy using dynamic client
	err = h.getCaddyClient().LoadConfig(caddyConfig)
	h.lastSyncedAt = time.Now()
	if err != nil {
		h.lastSyncError = err.Error()
	} else {
		h.lastSyncError = ""
	}
	return err
}

// PreviewImport returns what would be imported from Caddy
func (h *Handler) PreviewImport(c *gin.Context) {
	routes, status, err := fetchParsedCaddyRoutes(h)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	localRoutes, err := h.store.ListRoutes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read local routes"})
		return
	}
	summary, groups := buildImportPreview(routes, localRoutes)
	c.JSON(http.StatusOK, gin.H{
		"summary":  summary,
		"groups":   groups,
		"warnings": []string{driftWarning},
	})
}

// ImportFromCaddy pulls config from Caddy and overwrites local routes
func (h *Handler) ImportFromCaddy(c *gin.Context) {
	routes, status, err := fetchParsedCaddyRoutes(h)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error(), "recovery_guidance": recoveryGuidance})
		return
	}
	localRoutes, err := h.store.ListRoutes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read local routes", "recovery_guidance": recoveryGuidance})
		return
	}
	preserveMatchingLocalIDs(routes, localRoutes)

	if err := h.store.ReplaceAllRoutes(routes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save imported routes", "recovery_guidance": recoveryGuidance})
		return
	}

	summary, _ := buildImportPreview(routes, nil)
	c.JSON(http.StatusOK, gin.H{
		"imported":           len(routes),
		"editable":           summary.Editable,
		"readonly_preserved": summary.ReadOnlyPreserved,
		"unsupported":        summary.Unsupported,
		"message":            "Configuration imported successfully",
		"warnings":           []string{driftWarning},
	})
}

// GetRouteDetails returns preserved raw JSON for read-only routes.
func (h *Handler) GetRouteDetails(c *gin.Context) {
	route, err := h.store.GetRoute(c.Param("id"))
	if err != nil || route == nil || !route.IsReadOnly() || len(route.RawCaddyRoute) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "route details not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"route": gin.H{
			"id":              route.ID,
			"domain":          route.Domain,
			"path":            route.Path,
			"handler_type":    route.HandlerType,
			"support_status":  route.SupportStatus,
			"readonly_reason": route.ReadOnlyReason,
		},
		"raw_caddy_route": route.RawCaddyRoute,
	})
}
