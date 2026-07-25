package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/ArtemStepanov/caddy-admin-ui/internal/storage"
)

func init() {
	// Set Gin to test mode to reduce log output
	gin.SetMode(gin.TestMode)
}

// setupTestRouter creates a test router with a temporary database
func setupTestRouter(t *testing.T) (*gin.Engine, *storage.SQLiteStorage, func()) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "api_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	dbPath := filepath.Join(tmpDir, "test.db")
	store, err := storage.NewSQLiteStorage(dbPath)
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to create storage: %v", err)
	}

	router := gin.New()
	// Use a fake Caddy URL that will fail - tests should handle sync errors gracefully
	SetupRoutes(router, store, "http://localhost:29999")

	cleanup := func() {
		store.Close()
		os.RemoveAll(tmpDir)
	}

	return router, store, cleanup
}

func TestListRoutes_Empty(t *testing.T) {
	router, _, cleanup := setupTestRouter(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/routes", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]any
	json.Unmarshal(w.Body.Bytes(), &response)

	routes, ok := response["routes"].([]any)
	if !ok {
		t.Fatal("Expected routes array in response")
	}
	if len(routes) != 0 {
		t.Errorf("Expected 0 routes, got %d", len(routes))
	}
}

func TestCreateRoute_Success(t *testing.T) {
	router, _, cleanup := setupTestRouter(t)
	defer cleanup()

	body := `{
		"domain": "example.com",
		"handler_type": "reverse_proxy",
		"config": {"upstreams": ["localhost:8080"]}
	}`

	req := httptest.NewRequest("POST", "/api/routes", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusCreated, w.Code, w.Body.String())
	}

	var response map[string]any
	json.Unmarshal(w.Body.Bytes(), &response)

	route, ok := response["route"].(map[string]any)
	if !ok {
		t.Fatal("Expected route in response")
	}
	if route["domain"] != "example.com" {
		t.Errorf("Expected domain example.com, got %v", route["domain"])
	}
	if route["id"] == nil {
		t.Error("Expected route to have an ID")
	}

	// May have a warning about Caddy sync failure
	if response["warning"] != nil {
		t.Logf("Expected warning about Caddy sync: %v", response["warning"])
	}
}

func TestCreateRoute_MissingDomain(t *testing.T) {
	router, _, cleanup := setupTestRouter(t)
	defer cleanup()

	body := `{
		"handler_type": "reverse_proxy",
		"config": {"upstreams": ["localhost:8080"]}
	}`

	req := httptest.NewRequest("POST", "/api/routes", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response map[string]any
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["error"] == nil {
		t.Error("Expected error message in response")
	}
}

func TestCreateRoute_MissingHandler(t *testing.T) {
	router, _, cleanup := setupTestRouter(t)
	defer cleanup()

	body := `{
		"domain": "example.com",
		"config": {"upstreams": ["localhost:8080"]}
	}`

	req := httptest.NewRequest("POST", "/api/routes", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetRoute_Success(t *testing.T) {
	router, store, cleanup := setupTestRouter(t)
	defer cleanup()

	// Create a route directly in storage
	route := &storage.Route{
		Domain:      "example.com",
		HandlerType: "reverse_proxy",
		Config:      json.RawMessage(`{"upstreams":["localhost:8080"]}`),
	}
	store.CreateRoute(route)

	req := httptest.NewRequest("GET", "/api/routes/"+route.ID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]any
	json.Unmarshal(w.Body.Bytes(), &response)

	responseRoute, ok := response["route"].(map[string]any)
	if !ok {
		t.Fatal("Expected route in response")
	}
	if responseRoute["id"] != route.ID {
		t.Errorf("Expected ID %s, got %v", route.ID, responseRoute["id"])
	}
}

func TestGetRoute_NotFound(t *testing.T) {
	router, _, cleanup := setupTestRouter(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/routes/non-existing-id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestUpdateRoute_Success(t *testing.T) {
	router, store, cleanup := setupTestRouter(t)
	defer cleanup()

	// Create a route directly in storage
	route := &storage.Route{
		Domain:      "example.com",
		HandlerType: "reverse_proxy",
		Config:      json.RawMessage(`{"upstreams":["localhost:8080"]}`),
	}
	store.CreateRoute(route)

	body := `{
		"domain": "updated.example.com",
		"handler_type": "reverse_proxy",
		"config": {"upstreams": ["localhost:9090"]},
		"enabled": true
	}`

	req := httptest.NewRequest("PUT", "/api/routes/"+route.ID, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	// Verify update
	updated, _ := store.GetRoute(route.ID)
	if updated.Domain != "updated.example.com" {
		t.Errorf("Expected domain updated.example.com, got %s", updated.Domain)
	}
}

func TestUpdateRoute_PreservesEnabled(t *testing.T) {
	router, store, cleanup := setupTestRouter(t)
	defer cleanup()

	// Create an enabled route
	route := &storage.Route{
		Domain:      "example.com",
		HandlerType: "reverse_proxy",
		Config:      json.RawMessage(`{"upstreams":["localhost:8080"]}`),
		Enabled:     true,
	}
	store.CreateRoute(route)

	// Update without "enabled" field
	body := `{
		"domain": "example.com",
		"handler_type": "reverse_proxy",
		"config": {"upstreams": ["localhost:9090"]}
	}`

	req := httptest.NewRequest("PUT", "/api/routes/"+route.ID, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	// Verify update preserved enabled state
	updated, _ := store.GetRoute(route.ID)
	if !updated.Enabled {
		t.Error("Expected route to remain enabled")
	}

	// Verify config was updated
	if !bytes.Contains(updated.Config, []byte("9090")) {
		t.Error("Expected config to be updated")
	}
}

func TestUpdateRoute_EditableRouteHasNoRawCaddyRoute(t *testing.T) {
	router, store, cleanup := setupTestRouter(t)
	defer cleanup()

	// Create an editable route
	route := &storage.Route{
		Domain:      "example.com",
		HandlerType: "reverse_proxy",
		Config:      json.RawMessage(`{"upstreams":["localhost:8080"]}`),
		Enabled:     true,
	}
	store.CreateRoute(route)

	// Update
	body := `{
		"domain": "example.com",
		"handler_type": "reverse_proxy",
		"config": {"upstreams":["localhost:9090"]}
	}`

	req := httptest.NewRequest("PUT", "/api/routes/"+route.ID, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	// Verify editable route remains free of preserved raw config
	updated, _ := store.GetRoute(route.ID)
	if len(updated.RawCaddyRoute) != 0 {
		t.Errorf("Expected RawCaddyRoute to be empty for editable route, got %s", updated.RawCaddyRoute)
	}
}

func TestUpdateRoute_RejectsReadOnlyRoute(t *testing.T) {
	router, store, cleanup := setupTestRouter(t)
	defer cleanup()

	// Create a read-only preserved route
	route := &storage.Route{
		Domain:         "example.com",
		HandlerType:    "reverse_proxy",
		Config:         json.RawMessage(`{}`),
		Enabled:        true,
		SupportStatus:  storage.SupportStatusUnsupportedReadOnly,
		ReadOnlyReason: "unknown handler",
		RawCaddyRoute:  json.RawMessage(`{"original": "data"}`),
	}
	store.CreateRoute(route)

	body := `{
		"domain": "example.com",
		"handler_type": "reverse_proxy",
		"config": {}
	}`

	req := httptest.NewRequest("PUT", "/api/routes/"+route.ID, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusConflict, w.Code, w.Body.String())
	}
}

func TestUpdateRoute_NotFound(t *testing.T) {
	router, _, cleanup := setupTestRouter(t)
	defer cleanup()

	body := `{
		"domain": "example.com",
		"handler_type": "reverse_proxy",
		"config": {}
	}`

	req := httptest.NewRequest("PUT", "/api/routes/non-existing-id", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestDeleteRoute_Success(t *testing.T) {
	router, store, cleanup := setupTestRouter(t)
	defer cleanup()

	// Create a route directly in storage
	route := &storage.Route{
		Domain:      "example.com",
		HandlerType: "reverse_proxy",
		Config:      json.RawMessage(`{}`),
	}
	store.CreateRoute(route)

	req := httptest.NewRequest("DELETE", "/api/routes/"+route.ID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify deletion
	_, err := store.GetRoute(route.ID)
	if err == nil {
		t.Error("Expected route to be deleted")
	}
}

func TestToggleRoute(t *testing.T) {
	router, store, cleanup := setupTestRouter(t)
	defer cleanup()

	// Create an enabled route
	route := &storage.Route{
		Domain:      "example.com",
		HandlerType: "reverse_proxy",
		Config:      json.RawMessage(`{}`),
		Enabled:     true,
	}
	store.CreateRoute(route)

	// Toggle (should disable)
	req := httptest.NewRequest("POST", "/api/routes/"+route.ID+"/toggle", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	updated, _ := store.GetRoute(route.ID)
	if updated.Enabled {
		t.Error("Expected route to be disabled after toggle")
	}

	// Toggle again (should enable)
	req = httptest.NewRequest("POST", "/api/routes/"+route.ID+"/toggle", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	updated, _ = store.GetRoute(route.ID)
	if !updated.Enabled {
		t.Error("Expected route to be enabled after second toggle")
	}
}

func TestToggleRoute_NotFound(t *testing.T) {
	router, _, cleanup := setupTestRouter(t)
	defer cleanup()

	req := httptest.NewRequest("POST", "/api/routes/non-existing-id/toggle", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestGetConfig_UsesDefaultCaddyURL(t *testing.T) {
	router, _, cleanup := setupTestRouter(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/config", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	var response map[string]any
	json.Unmarshal(w.Body.Bytes(), &response)
	config := response["config"].(map[string]any)
	if config["caddy_admin_url"] != "http://localhost:29999" {
		t.Errorf("Expected handler default caddy url, got %v", config["caddy_admin_url"])
	}
}

func TestGetConfig(t *testing.T) {
	router, _, cleanup := setupTestRouter(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/config", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]any
	json.Unmarshal(w.Body.Bytes(), &response)

	config, ok := response["config"].(map[string]any)
	if !ok {
		t.Fatal("Expected config in response")
	}

	// Check defaults
	if config["caddy_admin_url"] != "http://localhost:29999" {
		t.Errorf("Expected default caddy_admin_url, got %v", config["caddy_admin_url"])
	}
}

func TestUpdateConfig(t *testing.T) {
	router, store, cleanup := setupTestRouter(t)
	defer cleanup()

	body := `{
		"caddy_admin_url": "http://custom:2019",
		"enable_encode": false
	}`

	req := httptest.NewRequest("PUT", "/api/config", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify update
	cfg, _ := store.GetGlobalConfig()
	if cfg.CaddyAdminURL != "http://custom:2019" {
		t.Errorf("Expected CaddyAdminURL http://custom:2019, got %s", cfg.CaddyAdminURL)
	}
	if cfg.EnableEncode {
		t.Error("Expected EnableEncode to be false")
	}
}

func TestGetStatus(t *testing.T) {
	router, _, cleanup := setupTestRouter(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/status", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]any
	json.Unmarshal(w.Body.Bytes(), &response)

	// Should be offline since we're using a fake Caddy URL
	if response["status"] != "offline" {
		t.Logf("Status response: %v", response)
	}

	// Should have latency field
	if response["latency"] == nil {
		t.Error("Expected latency in response")
	}

	// Should report the application version
	if response["version"] != "dev" {
		t.Errorf("Expected version dev, got %v", response["version"])
	}
}

func TestSyncToCaddy(t *testing.T) {
	router, _, cleanup := setupTestRouter(t)
	defer cleanup()

	req := httptest.NewRequest("POST", "/api/sync", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should fail because Caddy is not running
	if w.Code != http.StatusInternalServerError {
		t.Logf("Sync response (code %d): %s", w.Code, w.Body.String())
	}
}

func TestTestConnection(t *testing.T) {
	router, _, cleanup := setupTestRouter(t)
	defer cleanup()

	body := `{"url": "http://localhost:29999"}`

	req := httptest.NewRequest("POST", "/api/test-connection", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]any
	json.Unmarshal(w.Body.Bytes(), &response)

	// Should fail to connect
	if response["success"] == true {
		t.Error("Expected connection to fail")
	}
	if response["latency"] == nil {
		t.Error("Expected latency in response")
	}
}

func TestTestConnection_MissingURL(t *testing.T) {
	router, _, cleanup := setupTestRouter(t)
	defer cleanup()

	body := `{}`

	req := httptest.NewRequest("POST", "/api/test-connection", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestNoCORSHeaders(t *testing.T) {
	router, _, cleanup := setupTestRouter(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/routes", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Error("Expected no Access-Control-Allow-Origin header")
	}
}

func TestListRoutes_WithRoutes(t *testing.T) {
	router, store, cleanup := setupTestRouter(t)
	defer cleanup()

	// Create some routes
	store.CreateRoute(&storage.Route{
		Domain:      "a.example.com",
		HandlerType: "reverse_proxy",
		Config:      json.RawMessage(`{}`),
	})
	store.CreateRoute(&storage.Route{
		Domain:      "b.example.com",
		HandlerType: "file_server",
		Config:      json.RawMessage(`{}`),
	})

	req := httptest.NewRequest("GET", "/api/routes", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]any
	json.Unmarshal(w.Body.Bytes(), &response)

	routes, ok := response["routes"].([]any)
	if !ok {
		t.Fatal("Expected routes array in response")
	}
	if len(routes) != 2 {
		t.Errorf("Expected 2 routes, got %d", len(routes))
	}
}

func TestValidateCaddyURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"valid http", "http://caddy:2019", false},
		{"valid https", "https://caddy.example.com:2019", false},
		{"valid localhost", "http://localhost:2019", false},
		{"valid ip", "http://192.168.1.1:2019", false},
		{"wrong scheme ftp", "ftp://caddy:2019", true},
		{"wrong scheme javascript", "javascript:alert(1)", true},
		{"empty string", "", true},
		{"no host", "http://", true},
		{"no scheme", "caddy:2019", true},
		{"unparseable", "://invalid", true},
		{"with path", "http://caddy:2019/admin", true},
		{"with query", "http://caddy:2019?foo=bar", true},
		{"with fragment", "http://caddy:2019#section", true},
		{"with userinfo", "http://user:pass@caddy:2019", true},
		{"trailing slash ok", "http://caddy:2019/", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCaddyURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateCaddyURL(%q) error = %v, wantErr %v", tt.url, err, tt.wantErr)
			}
		})
	}
}

func TestUpdateConfig_RejectsInvalidURL(t *testing.T) {
	router, _, cleanup := setupTestRouter(t)
	defer cleanup()

	body := `{"caddy_admin_url": "ftp://bad:2019", "enable_encode": true}`

	req := httptest.NewRequest("PUT", "/api/config", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}

func TestTestConnection_RejectsInvalidURL(t *testing.T) {
	router, _, cleanup := setupTestRouter(t)
	defer cleanup()

	body := `{"url": "ftp://bad:2019"}`

	req := httptest.NewRequest("POST", "/api/test-connection", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}

func caddyConfigServer(t *testing.T, status int, body string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/config/" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
}

const mixedCaddyConfig = `{"apps":{"http":{"servers":{"srv0":{"routes":[{"match":[{"host":["app.example.com"]}],"handle":[{"handler":"reverse_proxy","upstreams":[{"dial":"localhost:8080"}]}]},{"match":[{"host":["raw.example.com"]}],"handle":[{"handler":"custom_handler","secret":"keep"}]}]}}}}}`

func TestImportPreviewSummaryGroupsAndNoMutation(t *testing.T) {
	router, store, cleanup := setupTestRouter(t)
	defer cleanup()
	server := caddyConfigServer(t, http.StatusOK, mixedCaddyConfig)
	defer server.Close()
	_ = store.SetGlobalConfig(&storage.GlobalConfig{CaddyAdminURL: server.URL, EnableEncode: true})
	local := &storage.Route{Domain: "local.example.com", HandlerType: "reverse_proxy", Config: json.RawMessage(`{}`), Enabled: true}
	_ = store.CreateRoute(local)

	req := httptest.NewRequest("POST", "/api/import-preview", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d: %s", w.Code, w.Body.String())
	}
	var resp struct {
		Summary importSummary `json:"summary"`
		Groups  importGroups  `json:"groups"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Summary.TotalFound != 2 || resp.Summary.Editable != 1 || resp.Summary.ReadOnlyPreserved != 1 || resp.Summary.Unsupported != 1 || resp.Summary.LocalOnly != 1 {
		t.Fatalf("bad summary: %+v", resp.Summary)
	}
	if resp.Groups.WillUpdate == nil {
		t.Fatal("empty groups must be [] not null")
	}
	if len(resp.Groups.NewFromCaddy) != 1 || len(resp.Groups.ReadOnlyPreserved) != 1 || len(resp.Groups.LocalOnly) != 1 {
		t.Fatalf("bad groups: %+v", resp.Groups)
	}
	if !bytes.Contains(resp.Groups.ReadOnlyPreserved[0].RawCaddyRoute, []byte("custom_handler")) {
		t.Fatalf("read-only preview row missing raw JSON: %+v", resp.Groups.ReadOnlyPreserved[0])
	}
	got, _ := store.ListRoutes()
	if len(got) != 1 || got[0].ID != local.ID {
		t.Fatalf("preview mutated storage: %+v", got)
	}
}

func TestImportPreviewFailuresLeaveLocalState(t *testing.T) {
	cases := []struct {
		name   string
		url    string
		body   string
		status int
		want   int
	}{
		{name: "fetch failure", url: "http://localhost:29999", want: http.StatusBadGateway},
		{name: "parse failure", body: `{"secret":"do-not-leak"`, status: http.StatusOK, want: http.StatusInternalServerError},
		{name: "zero routes", body: `{}`, status: http.StatusOK, want: http.StatusOK},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			router, store, cleanup := setupTestRouter(t)
			defer cleanup()
			if tc.body != "" {
				server := caddyConfigServer(t, tc.status, tc.body)
				defer server.Close()
				tc.url = server.URL
			}
			_ = store.SetGlobalConfig(&storage.GlobalConfig{CaddyAdminURL: tc.url, EnableEncode: true})
			_ = store.CreateRoute(&storage.Route{Domain: "local.example.com", HandlerType: "reverse_proxy", Config: json.RawMessage(`{}`), Enabled: true})

			req := httptest.NewRequest("POST", "/api/import-preview", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			if w.Code != tc.want {
				t.Fatalf("status %d want %d: %s", w.Code, tc.want, w.Body.String())
			}
			if bytes.Contains(w.Body.Bytes(), []byte("do-not-leak")) {
				t.Fatal("error leaked response body")
			}
			got, _ := store.ListRoutes()
			if len(got) != 1 {
				t.Fatalf("preview changed storage: %+v", got)
			}
		})
	}
}

func TestImportFromCaddyUpdatesExistingReadOnlyFileServerToEditable(t *testing.T) {
	router, store, cleanup := setupTestRouter(t)
	defer cleanup()
	server := caddyConfigServer(t, http.StatusOK, `{"apps":{"http":{"servers":{"srv0":{"routes":[{"match":[{"host":["files.localhost"]}],"handle":[{"handler":"subroute","routes":[{"handle":[{"handler":"vars","root":"/usr/share/caddy"},{"handler":"file_server","hide":["/etc/caddy/Caddyfile"]}]}]}]}]}}}}}`)
	defer server.Close()
	_ = store.SetGlobalConfig(&storage.GlobalConfig{CaddyAdminURL: server.URL, EnableEncode: true})
	existing := &storage.Route{
		Domain:         "files.localhost",
		HandlerType:    "file_server",
		Config:         json.RawMessage(`{}`),
		Enabled:        true,
		SupportStatus:  storage.SupportStatusPartialReadOnly,
		ReadOnlyReason: "old classification",
		RawCaddyRoute:  json.RawMessage(`{"old":true}`),
	}
	if err := store.CreateRoute(existing); err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest("POST", "/api/import", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status %d: %s", w.Code, w.Body.String())
	}
	got, err := store.GetRoute(existing.ID)
	if err != nil {
		t.Fatalf("existing route id should be preserved: %v", err)
	}
	if got.ReadOnly || got.SupportStatus != storage.SupportStatusEditable || got.ReadOnlyReason != "" || len(got.RawCaddyRoute) != 0 {
		t.Fatalf("files.localhost should be updated to editable: %+v raw=%s", got, got.RawCaddyRoute)
	}
	var cfg storage.FileServerConfig
	if err := json.Unmarshal(got.Config, &cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Root != "/usr/share/caddy" {
		t.Fatalf("expected imported file server root, got %+v", cfg)
	}
}

func TestImportFromCaddyTransactionalSuccessAndFailure(t *testing.T) {
	router, store, cleanup := setupTestRouter(t)
	defer cleanup()
	server := caddyConfigServer(t, http.StatusOK, mixedCaddyConfig)
	defer server.Close()
	_ = store.SetGlobalConfig(&storage.GlobalConfig{CaddyAdminURL: server.URL, EnableEncode: true})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest("POST", "/api/import", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["imported"].(float64) != 2 || resp["editable"].(float64) != 1 || resp["readonly_preserved"].(float64) != 1 {
		t.Fatalf("bad import response: %v", resp)
	}
	routes, _ := store.ListRoutes()
	if len(routes) != 2 {
		t.Fatalf("expected 2 routes, got %+v", routes)
	}

	bad := caddyConfigServer(t, http.StatusOK, `{"bad"`)
	defer bad.Close()
	_ = store.SetGlobalConfig(&storage.GlobalConfig{CaddyAdminURL: bad.URL, EnableEncode: true})
	w = httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest("POST", "/api/import", nil))
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status %d: %s", w.Code, w.Body.String())
	}
	after, _ := store.ListRoutes()
	if len(after) != 2 {
		t.Fatalf("failed import changed routes: %+v", after)
	}
}

func TestImportFromCaddyDoesNotReuseLocalIDForDuplicateKeys(t *testing.T) {
	router, store, cleanup := setupTestRouter(t)
	defer cleanup()
	server := caddyConfigServer(t, http.StatusOK, `{"apps":{"http":{"servers":{"srv0":{"routes":[{"match":[{"host":["dup.example.com"]}],"handle":[{"handler":"custom_one"}]},{"match":[{"host":["dup.example.com"]}],"handle":[{"handler":"custom_two"}]}]}}}}}`)
	defer server.Close()
	_ = store.SetGlobalConfig(&storage.GlobalConfig{CaddyAdminURL: server.URL, EnableEncode: true})
	existing := &storage.Route{Domain: "dup.example.com", HandlerType: "unknown", Config: json.RawMessage(`{}`), Enabled: true, SupportStatus: storage.SupportStatusUnsupportedReadOnly, ReadOnlyReason: "unknown handler", RawCaddyRoute: json.RawMessage(`{"handle":[{"handler":"old"}]}`)}
	if err := store.CreateRoute(existing); err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest("POST", "/api/import", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("duplicate imported keys should not collide, got %d: %s", w.Code, w.Body.String())
	}
	routes, _ := store.ListRoutes()
	ids := map[string]bool{}
	for _, route := range routes {
		ids[route.ID] = true
	}
	if len(routes) != 2 || !ids[existing.ID] || len(ids) != 2 {
		t.Fatalf("expected two unique routes preserving one local id, got %+v", routes)
	}
}

func TestImportPreviewPreservesDynamicUpstreams(t *testing.T) {
	router, store, cleanup := setupTestRouter(t)
	defer cleanup()
	server := caddyConfigServer(t, http.StatusOK, `{"apps":{"http":{"servers":{"srv0":{"routes":[{"match":[{"host":["dynamic.example.com"]}],"handle":[{"handler":"reverse_proxy","dynamic_upstreams":{"source":"srv"}}]}]}}}}}`)
	defer server.Close()
	_ = store.SetGlobalConfig(&storage.GlobalConfig{CaddyAdminURL: server.URL, EnableEncode: true})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest("POST", "/api/import-preview", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("dynamic upstreams should preview as read-only, got %d: %s", w.Code, w.Body.String())
	}
	var resp struct {
		Summary importSummary `json:"summary"`
		Groups  importGroups  `json:"groups"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Summary.ReadOnlyPreserved != 1 || len(resp.Groups.ReadOnlyPreserved) != 1 || resp.Groups.ReadOnlyPreserved[0].HandlerType != "reverse_proxy" {
		t.Fatalf("dynamic upstream route not preserved read-only: %+v", resp)
	}
}

func TestSyncRejectsInvalidReadOnlyRawRouteWithoutDeletingRoute(t *testing.T) {
	router, store, cleanup := setupTestRouter(t)
	defer cleanup()
	route := &storage.Route{
		Domain:         "broken.example.com",
		HandlerType:    "unknown",
		Config:         json.RawMessage(`{}`),
		Enabled:        true,
		SupportStatus:  storage.SupportStatusUnsupportedReadOnly,
		ReadOnlyReason: "unknown handler",
		RawCaddyRoute:  json.RawMessage(`{"not":"a route"}`),
	}
	if err := store.CreateRoute(route); err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest("POST", "/api/sync", nil))
	if w.Code != http.StatusInternalServerError || !bytes.Contains(w.Body.Bytes(), []byte("recovery_guidance")) {
		t.Fatalf("sync should fail with recovery guidance, got %d: %s", w.Code, w.Body.String())
	}
	if _, err := store.GetRoute(route.ID); err != nil {
		t.Fatal("sync failure deleted route")
	}
}

func TestSyncRejectsInvalidEditableRouteConfigWithoutDeletingRoute(t *testing.T) {
	router, store, cleanup := setupTestRouter(t)
	defer cleanup()
	route := &storage.Route{Domain: "broken.example.com", HandlerType: "reverse_proxy", Config: json.RawMessage(`{}`), Enabled: true, SupportStatus: storage.SupportStatusEditable}
	if err := store.CreateRoute(route); err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest("POST", "/api/sync", nil))
	if w.Code != http.StatusInternalServerError || !bytes.Contains(w.Body.Bytes(), []byte("invalid reverse_proxy config")) {
		t.Fatalf("sync should reject invalid config, got %d: %s", w.Code, w.Body.String())
	}
	if _, err := store.GetRoute(route.ID); err != nil {
		t.Fatal("sync failure deleted route")
	}
}

func TestImportRejectsInvalidEditableRouteConfigWithoutChangingLocalRoutes(t *testing.T) {
	router, store, cleanup := setupTestRouter(t)
	defer cleanup()
	server := caddyConfigServer(t, http.StatusOK, `{"apps":{"http":{"servers":{"srv0":{"routes":[{"match":[{"host":["broken.example.com"]}],"handle":[{"handler":"reverse_proxy"}]}]}}}}}`)
	defer server.Close()
	_ = store.SetGlobalConfig(&storage.GlobalConfig{CaddyAdminURL: server.URL, EnableEncode: true})
	local := &storage.Route{Domain: "keep.example.com", HandlerType: "file_server", Config: json.RawMessage(`{}`), Enabled: true}
	if err := store.CreateRoute(local); err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest("POST", "/api/import", nil))
	if w.Code != http.StatusInternalServerError || !bytes.Contains(w.Body.Bytes(), []byte("invalid reverse_proxy config")) {
		t.Fatalf("import should reject invalid config, got %d: %s", w.Code, w.Body.String())
	}
	routes, _ := store.ListRoutes()
	if len(routes) != 1 || routes[0].ID != local.ID {
		t.Fatalf("failed import changed local routes: %+v", routes)
	}
}

func TestImportPreviewUsesFallbackSummaryWhenRouteHasNoHost(t *testing.T) {
	router, store, cleanup := setupTestRouter(t)
	defer cleanup()
	server := caddyConfigServer(t, http.StatusOK, `{"apps":{"http":{"servers":{"srv0":{"routes":[{"match":[{"path":["/api/*"]}],"handle":[{"handler":"reverse_proxy","upstreams":[{"dial":"localhost:8080"}]}]},{"handle":[{"handler":"file_server","root":"/srv/www"}]}]}}}}}`)
	defer server.Close()
	_ = store.SetGlobalConfig(&storage.GlobalConfig{CaddyAdminURL: server.URL, EnableEncode: true})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest("POST", "/api/import-preview", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status %d: %s", w.Code, w.Body.String())
	}
	var resp struct {
		Groups importGroups `json:"groups"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Groups.NewFromCaddy) != 2 {
		t.Fatalf("expected 2 rows, got %+v", resp.Groups.NewFromCaddy)
	}
	if resp.Groups.NewFromCaddy[0].Domain != "/api/*" || resp.Groups.NewFromCaddy[1].Domain != "/srv/www" {
		t.Fatalf("fallback summaries missing: %+v", resp.Groups.NewFromCaddy)
	}
}

func TestReadOnlyRouteMutationRejectionDetailsAndSyncWarning(t *testing.T) {
	router, store, cleanup := setupTestRouter(t)
	defer cleanup()
	route := &storage.Route{
		Domain:         "raw.example.com",
		HandlerType:    "unknown",
		Config:         json.RawMessage(`{}`),
		Enabled:        true,
		SupportStatus:  storage.SupportStatusUnsupportedReadOnly,
		ReadOnlyReason: "unknown handler",
		RawCaddyRoute:  json.RawMessage(`{"match":[{"host":["raw.example.com"]}],"handle":[{"handler":"custom"}]}`),
	}
	_ = store.CreateRoute(route)

	for _, req := range []*http.Request{
		httptest.NewRequest("PUT", "/api/routes/"+route.ID, bytes.NewBufferString(`{"domain":"x","handler_type":"unknown","config":{}}`)),
		httptest.NewRequest("DELETE", "/api/routes/"+route.ID, nil),
		httptest.NewRequest("POST", "/api/routes/"+route.ID+"/toggle", nil),
	} {
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusConflict {
			t.Fatalf("%s got %d: %s", req.Method, w.Code, w.Body.String())
		}
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest("GET", "/api/routes/"+route.ID+"/details", nil))
	if w.Code != http.StatusOK || !bytes.Contains(w.Body.Bytes(), []byte("custom")) {
		t.Fatalf("details status %d: %s", w.Code, w.Body.String())
	}

	w = httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest("POST", "/api/sync", nil))
	if w.Code != http.StatusInternalServerError || !bytes.Contains(w.Body.Bytes(), []byte("recovery_guidance")) {
		t.Fatalf("sync warning status %d: %s", w.Code, w.Body.String())
	}
	if _, err := store.GetRoute(route.ID); err != nil {
		t.Fatal("sync failure deleted route")
	}
}

func TestCreateUpdateRouteIgnoreClientReadOnlyFields(t *testing.T) {
	router, store, cleanup := setupTestRouter(t)
	defer cleanup()

	createBody := `{"domain":"client.example.com","handler_type":"reverse_proxy","config":{"upstreams":["localhost:8080"]},"support_status":"partial_readonly","readonly_reason":"client says readonly","readonly":true}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/routes", bytes.NewBufferString(createBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("create status %d: %s", w.Code, w.Body.String())
	}
	var created struct {
		Route storage.Route `json:"route"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	got, err := store.GetRoute(created.Route.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.ReadOnly || got.SupportStatus != storage.SupportStatusEditable || got.ReadOnlyReason != "" {
		t.Fatalf("create persisted client read-only fields: %+v", got)
	}

	updateBody := `{"domain":"client.example.com","handler_type":"reverse_proxy","config":{"upstreams":["localhost:8081"]},"support_status":"unsupported_readonly","readonly_reason":"client says readonly","readonly":true}`
	w = httptest.NewRecorder()
	req = httptest.NewRequest("PUT", "/api/routes/"+got.ID, bytes.NewBufferString(updateBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("update status %d: %s", w.Code, w.Body.String())
	}
	got, err = store.GetRoute(got.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.ReadOnly || got.SupportStatus != storage.SupportStatusEditable || got.ReadOnlyReason != "" {
		t.Fatalf("update persisted client read-only fields: %+v", got)
	}
}

func TestCreateRoute_WithHeaders(t *testing.T) {
	router, _, cleanup := setupTestRouter(t)
	defer cleanup()

	body := `{
		"domain": "example.com",
		"handler_type": "reverse_proxy",
		"config": {"upstreams": ["localhost:8080"]},
		"headers": {
			"set": {"X-Frame-Options": "DENY"},
			"delete": ["Server"]
		}
	}`

	req := httptest.NewRequest("POST", "/api/routes", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusCreated, w.Code, w.Body.String())
	}
}

// TestImportRouteDecisionOutcomes asserts the single post-parse safety gate returns
// one consistent accept (editable), preserve (read-only with raw content), or reject
// (invalid) decision for every imported-route case, and that rejection leaves local
// routes unchanged. Covers SC-001/SC-002 for the import endpoint.
func TestImportRouteDecisionOutcomes(t *testing.T) {
	cases := []struct {
		name        string
		config      string
		wantStatus  int
		editable    int // expected accept count (0 if not an accept)
		preserved   int // expected readonly_preserved count (0 if not a preserve)
		unsupported int // expected unsupported subset count
	}{
		{
			name:       "editable route accepted",
			config:      `{"apps":{"http":{"servers":{"srv0":{"routes":[{"match":[{"host":["ed.example.com"]}],"handle":[{"handler":"reverse_proxy","upstreams":[{"dial":"localhost:8080"}]}]}]}}}}}`,
			wantStatus:  http.StatusOK,
			editable:   1,
			preserved:  0,
		},
		{
			name:       "partially supported route preserved read-only",
			config:      `{"apps":{"http":{"servers":{"srv0":{"routes":[{"match":[{"host":["partial.example.com"]}],"handle":[{"handler":"file_server"},{"handler":"vars","root":"/srv"}]}]}}}}}`,
			wantStatus:  http.StatusOK,
			preserved:  1,
		},
		{
			name:       "unsupported route preserved read-only",
			config:      `{"apps":{"http":{"servers":{"srv0":{"routes":[{"match":[{"host":["unsup.example.com"]}],"handle":[{"handler":"custom_handler","secret":"keep"}]}]}}}}}`,
			wantStatus:  http.StatusOK,
			preserved:  1,
			unsupported: 1,
		},
		{
			// Legacy-preserved routes (DB backfill target) share the same partial_readonly
			// preserve outcome; a dynamic-upstreams proxy is the parser's preserve case
			// that retains raw content and a read-only reason.
			name:       "legacy-preserved route retained read-only with raw content",
			config:      `{"apps":{"http":{"servers":{"srv0":{"routes":[{"match":[{"host":["legacy.example.com"]}],"handle":[{"handler":"reverse_proxy","dynamic_upstreams":{"source":"srv"}}]}]}}}}}`,
			wantStatus:  http.StatusOK,
			preserved:  1,
		},
		{
			name:       "invalid editable route rejected before local replacement",
			config:      `{"apps":{"http":{"servers":{"srv0":{"routes":[{"match":[{"host":["broken.example.com"]}],"handle":[{"handler":"reverse_proxy"}]}]}}}}}`,
			wantStatus:  http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			router, store, cleanup := setupTestRouter(t)
			defer cleanup()
			server := caddyConfigServer(t, http.StatusOK, tc.config)
			defer server.Close()
			_ = store.SetGlobalConfig(&storage.GlobalConfig{CaddyAdminURL: server.URL, EnableEncode: true})
			local := &storage.Route{Domain: "keep.example.com", HandlerType: "file_server", Config: json.RawMessage(`{}`), Enabled: true}
			if err := store.CreateRoute(local); err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, httptest.NewRequest("POST", "/api/import", nil))
			if w.Code != tc.wantStatus {
				t.Fatalf("status %d want %d: %s", w.Code, tc.wantStatus, w.Body.String())
			}

			if tc.wantStatus == http.StatusOK {
				var resp map[string]any
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatal(err)
				}
				if int(resp["editable"].(float64)) != tc.editable {
					t.Fatalf("editable=%v want %d", resp["editable"], tc.editable)
				}
				if int(resp["readonly_preserved"].(float64)) != tc.preserved {
					t.Fatalf("readonly_preserved=%v want %d", resp["readonly_preserved"], tc.preserved)
				}
				if tc.unsupported > 0 && int(resp["unsupported"].(float64)) != tc.unsupported {
					t.Fatalf("unsupported=%v want %d", resp["unsupported"], tc.unsupported)
				}
			}

			if tc.preserved > 0 {
				got, _ := store.ListRoutes()
				var preserved *storage.Route
				for _, r := range got {
					if r.IsReadOnly() && len(r.RawCaddyRoute) > 0 && r.ReadOnlyReason != "" {
					preserved = r
					break
				}
				}
				if preserved == nil {
					t.Fatalf("no preserved read-only route with raw content and reason: %+v", got)
				}
			}

			if tc.wantStatus == http.StatusInternalServerError {
				got, _ := store.ListRoutes()
				if len(got) != 1 || got[0].ID != local.ID {
					t.Fatalf("rejected import changed local routes: %+v", got)
				}
			}
		})
	}
}
