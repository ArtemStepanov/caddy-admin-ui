package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ArtemStepanov/caddy-admin-ui/internal/caddy"
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
	recoveryGuidance = "Review the live managed route array and latest snapshot before retrying; run setup again if ownership or the ETag changed."
)

// Handler contains all HTTP handlers
type Handler struct {
	store           *storage.SQLiteStorage
	defaultCaddyURL string // fallback URL from env
	lastSyncedAt    time.Time
	lastSyncError   string
	syncMu          sync.Mutex
	statusMu        sync.RWMutex
}

// NewHandler creates a new handler
func NewHandler(store *storage.SQLiteStorage, defaultCaddyURL string) *Handler {
	return &Handler{
		store:           store,
		defaultCaddyURL: defaultCaddyURL,
	}
}

// Readiness reports whether both the local database and the configured Caddy
// Admin API are reachable. It deliberately returns no topology or error details.
func (h *Handler) Readiness(c *gin.Context) {
	if _, err := h.store.GetGlobalConfig(); err != nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}
	if err := h.getCaddyClient().Health(); err != nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}
	c.Status(http.StatusNoContent)
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
	route.ID = uuid.NewString()
	route.CreatedAt = time.Time{}
	route.UpdatedAt = time.Time{}
	makeEditableRoute(&route)

	err := h.mutateRoutes("create route", func(routes []*storage.Route) ([]*storage.Route, error) {
		route.Position = len(routes)
		return append(routes, &route), nil
	})
	if errors.Is(err, errSetupRequired) {
		resp := syncWarning("Route saved as a local draft; Caddy setup is not complete", err)
		resp["route"] = route
		c.JSON(http.StatusCreated, resp)
		return
	}
	if err != nil {
		c.JSON(mutationErrorStatus(err), gin.H{"error": err.Error(), "recovery_guidance": recoveryGuidance})
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

	err = h.mutateRoutes("update route", func(routes []*storage.Route) ([]*storage.Route, error) {
		for i, candidate := range routes {
			if candidate.ID == id {
				route.Position = candidate.Position
				routes[i] = &route
				return routes, nil
			}
		}
		return nil, fmt.Errorf("route not found")
	})
	if errors.Is(err, errSetupRequired) {
		resp := syncWarning("Route saved as a local draft; Caddy setup is not complete", err)
		resp["route"] = route
		c.JSON(http.StatusOK, resp)
		return
	}
	if err != nil {
		c.JSON(mutationErrorStatus(err), gin.H{"error": err.Error(), "recovery_guidance": recoveryGuidance})
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

	err = h.mutateRoutes("delete route", func(routes []*storage.Route) ([]*storage.Route, error) {
		for i, candidate := range routes {
			if candidate.ID == id {
				return append(routes[:i], routes[i+1:]...), nil
			}
		}
		return nil, fmt.Errorf("route not found")
	})
	if errors.Is(err, errSetupRequired) {
		resp := syncWarning("Local draft deleted; Caddy setup is not complete", err)
		resp["message"] = resp["warning"]
		c.JSON(http.StatusOK, resp)
		return
	}
	if err != nil {
		c.JSON(mutationErrorStatus(err), gin.H{"error": err.Error(), "recovery_guidance": recoveryGuidance})
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

	err = h.mutateRoutes("toggle route", func(routes []*storage.Route) ([]*storage.Route, error) {
		for _, candidate := range routes {
			if candidate.ID == id {
				candidate.Enabled = !candidate.Enabled
				route = candidate
				return routes, nil
			}
		}
		return nil, fmt.Errorf("route not found")
	})
	if errors.Is(err, errSetupRequired) {
		resp := syncWarning("Route updated as a local draft; Caddy setup is not complete", err)
		resp["route"] = route
		c.JSON(http.StatusOK, resp)
		return
	}
	if err != nil {
		c.JSON(mutationErrorStatus(err), gin.H{"error": err.Error(), "recovery_guidance": recoveryGuidance})
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

	existing, err := h.store.GetGlobalConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	existingURL := existing.CaddyAdminURL
	if existingURL == "" {
		existingURL = h.defaultCaddyURL
	}
	// Ownership and concurrency state can only be changed by the setup flow.
	cfg.ManagedServer = ""
	cfg.SetupComplete = false
	cfg.LastETag = ""
	if strings.TrimRight(existingURL, "/") == strings.TrimRight(cfg.CaddyAdminURL, "/") {
		cfg.ManagedServer = existing.ManagedServer
		cfg.SetupComplete = existing.SetupComplete
		cfg.LastETag = existing.LastETag
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
	globalCfg, _ := h.store.GetGlobalConfig()
	setupComplete := globalCfg != nil && globalCfg.SetupComplete
	managedServer := ""
	if globalCfg != nil {
		managedServer = globalCfg.ManagedServer
	}

	start := time.Now()
	err := caddyClient.Health()
	latency := time.Since(start).Milliseconds()
	lastSyncedAt, lastSyncError := h.getSyncStatus()

	if err != nil {
		resp := gin.H{
			"status":         "offline",
			"error":          err.Error(),
			"latency":        latency,
			"admin_url":      caddyURL,
			"version":        version.Version,
			"setup_complete": setupComplete,
			"managed_server": managedServer,
		}
		if !lastSyncedAt.IsZero() {
			resp["last_synced_at"] = lastSyncedAt.Format(time.RFC3339)
			resp["last_sync_error"] = lastSyncError
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
		"status":         "online",
		"latency":        latency,
		"admin_url":      caddyURL,
		"route_count":    routeCount,
		"version":        version.Version,
		"setup_complete": setupComplete,
		"managed_server": managedServer,
	}
	if !lastSyncedAt.IsZero() {
		resp["last_synced_at"] = lastSyncedAt.Format(time.RFC3339)
		resp["last_sync_error"] = lastSyncError
	}
	c.JSON(http.StatusOK, resp)
}

// SyncToCaddy manually triggers sync to Caddy
func (h *Handler) SyncToCaddy(c *gin.Context) {
	if err := h.syncToCaddy(); err != nil {
		resp := syncWarning("Sync to Caddy failed", err)
		resp["error"] = err.Error()
		c.JSON(mutationErrorStatus(err), resp)
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

// syncToCaddy updates only the selected HTTP server's route array.
func (h *Handler) syncToCaddy() error {
	h.syncMu.Lock()
	defer h.syncMu.Unlock()
	routes, err := h.store.ListRoutes()
	if err != nil {
		h.setSyncStatus(err)
		return err
	}

	globalCfg, err := h.store.GetGlobalConfig()
	if err != nil {
		h.setSyncStatus(err)
		return err
	}
	return h.syncRoutesLocked(routes, "manual sync", globalCfg)
}

func mutationErrorStatus(err error) int {
	if errors.Is(err, errInvalidRoutes) {
		return http.StatusBadRequest
	}
	if errors.Is(err, errSetupRequired) || errors.Is(err, errConfigDrift) {
		return http.StatusConflict
	}
	return http.StatusBadGateway
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
