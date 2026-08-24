package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/ArtemStepanov/caddy-admin-ui/internal/storage"
)

type fakeManagedCaddy struct {
	mu          sync.Mutex
	routes      json.RawMessage
	routesETag  string
	rootVersion int
	patches     int
	loads       int
	lastMatch   string
}

func newFakeManagedCaddy(routes string) *fakeManagedCaddy {
	return &fakeManagedCaddy{routes: json.RawMessage(routes), routesETag: `"routes-1"`, rootVersion: 1}
}

func (f *fakeManagedCaddy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/config/":
		w.Header().Set("ETag", fmt.Sprintf(`"root-%d"`, f.rootVersion))
		_, _ = fmt.Fprintf(w, `{"logging":{"logs":{"default":{"level":"INFO"}}},"apps":{"http":{"servers":{"srv0":{"listen":[":8443"],"automatic_https":{"disable":true},"routes":%s},"untouched":{"listen":[":9000"],"routes":[]}}}}}`, f.routes)
	case r.Method == http.MethodGet && r.URL.Path == "/config/apps/http/servers/srv0/routes":
		w.Header().Set("ETag", f.routesETag)
		_, _ = w.Write(f.routes)
	case r.Method == http.MethodPatch && r.URL.Path == "/config/apps/http/servers/srv0/routes":
		f.lastMatch = r.Header.Get("If-Match")
		if f.lastMatch != f.routesETag {
			w.WriteHeader(http.StatusPreconditionFailed)
			return
		}
		var routes []json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&routes); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		f.routes, _ = json.Marshal(routes)
		f.patches++
		f.rootVersion++
		f.routesETag = fmt.Sprintf(`"routes-%d"`, f.patches+1)
		w.WriteHeader(http.StatusOK)
	case r.Method == http.MethodPost && r.URL.Path == "/load":
		f.loads++
		w.WriteHeader(http.StatusOK)
	default:
		http.NotFound(w, r)
	}
}

func requestJSON(t *testing.T, router http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var data []byte
	if body != nil {
		var err error
		data, err = json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
	}
	req := httptest.NewRequest(method, path, bytes.NewReader(data))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func connectManagedCaddy(t *testing.T, router http.Handler, serverURL string) setupPreview {
	t.Helper()
	w := requestJSON(t, router, http.MethodPost, "/api/setup/preview", setupRequest{URL: serverURL, Server: "srv0"})
	if w.Code != http.StatusOK {
		t.Fatalf("preview status %d: %s", w.Code, w.Body.String())
	}
	var preview setupPreview
	if err := json.Unmarshal(w.Body.Bytes(), &preview); err != nil {
		t.Fatal(err)
	}
	w = requestJSON(t, router, http.MethodPost, "/api/setup/confirm", setupRequest{URL: serverURL, Server: "srv0", PreviewToken: preview.PreviewToken})
	if w.Code != http.StatusOK {
		t.Fatalf("confirm status %d: %s", w.Code, w.Body.String())
	}
	return preview
}

func TestManagedSetupAndScopedSyncPreserveExternalConfig(t *testing.T) {
	router, store, cleanup := setupTestRouter(t)
	defer cleanup()
	fake := newFakeManagedCaddy(`[{"match":[{"host":["external.example.com"],"method":["GET"]}],"handle":[{"handler":"custom_handler","secret":"keep"}],"custom":{"keep":true}}]`)
	server := httptest.NewServer(fake)
	defer server.Close()

	preview := connectManagedCaddy(t, router, server.URL)
	if preview.RouteCount != 1 || preview.ReadOnly != 1 || preview.Editable != 0 {
		t.Fatalf("unexpected ownership preview: %+v", preview)
	}
	cfg, err := store.GetGlobalConfig()
	if err != nil || !cfg.SetupComplete || cfg.ManagedServer != "srv0" || cfg.LastETag != `"routes-1"` {
		t.Fatalf("managed config not saved: %+v err=%v", cfg, err)
	}

	w := requestJSON(t, router, http.MethodPost, "/api/routes", map[string]any{
		"domain": "app.example.com", "handler_type": "reverse_proxy", "config": map[string]any{"upstreams": []string{"app:8080"}},
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create status %d: %s", w.Code, w.Body.String())
	}
	fake.mu.Lock()
	patched := append(json.RawMessage(nil), fake.routes...)
	patches, loads, lastMatch := fake.patches, fake.loads, fake.lastMatch
	fake.mu.Unlock()
	if patches != 1 || loads != 0 || lastMatch != `"routes-1"` {
		t.Fatalf("expected one guarded scoped patch, patches=%d loads=%d if-match=%q", patches, loads, lastMatch)
	}
	var remote []map[string]any
	if err := json.Unmarshal(patched, &remote); err != nil {
		t.Fatal(err)
	}
	if len(remote) != 2 || remote[0]["custom"] == nil || remote[1]["@id"] == nil {
		t.Fatalf("external raw route or managed marker lost: %s", patched)
	}
	snapshots, err := store.ListSnapshots()
	if err != nil || len(snapshots) < 2 {
		t.Fatalf("expected setup and pre-write snapshots, got %+v err=%v", snapshots, err)
	}
}

func TestManagedSyncRefusesDriftWithoutPersistingMutation(t *testing.T) {
	router, store, cleanup := setupTestRouter(t)
	defer cleanup()
	fake := newFakeManagedCaddy(`[]`)
	server := httptest.NewServer(fake)
	defer server.Close()
	connectManagedCaddy(t, router, server.URL)

	fake.mu.Lock()
	fake.routesETag = `"external-change"`
	fake.routes = json.RawMessage(`[{"handle":[{"handler":"static_response","body":"outside"}]}]`)
	fake.rootVersion++
	fake.mu.Unlock()
	w := requestJSON(t, router, http.MethodPost, "/api/routes", map[string]any{
		"domain": "blocked.example.com", "handler_type": "reverse_proxy", "config": map[string]any{"upstreams": []string{"app:8080"}},
	})
	if w.Code != http.StatusConflict || !bytes.Contains(w.Body.Bytes(), []byte("changed outside")) {
		t.Fatalf("drift should return conflict, got %d: %s", w.Code, w.Body.String())
	}
	routes, _ := store.ListRoutes()
	if len(routes) != 0 {
		t.Fatalf("drifted mutation reached local storage: %+v", routes)
	}
	fake.mu.Lock()
	patches := fake.patches
	fake.mu.Unlock()
	if patches != 0 {
		t.Fatalf("drifted mutation sent %d patches", patches)
	}
}

func TestSetupConfirmationRejectsTOCTOUChange(t *testing.T) {
	router, store, cleanup := setupTestRouter(t)
	defer cleanup()
	fake := newFakeManagedCaddy(`[]`)
	server := httptest.NewServer(fake)
	defer server.Close()
	w := requestJSON(t, router, http.MethodPost, "/api/setup/preview", setupRequest{URL: server.URL, Server: "srv0"})
	var preview setupPreview
	_ = json.Unmarshal(w.Body.Bytes(), &preview)
	fake.mu.Lock()
	fake.rootVersion++
	fake.mu.Unlock()
	w = requestJSON(t, router, http.MethodPost, "/api/setup/confirm", setupRequest{URL: server.URL, Server: "srv0", PreviewToken: preview.PreviewToken})
	if w.Code != http.StatusConflict {
		t.Fatalf("changed preview should conflict, got %d: %s", w.Code, w.Body.String())
	}
	cfg, _ := store.GetGlobalConfig()
	if cfg.SetupComplete {
		t.Fatal("stale setup preview was persisted")
	}
}

func TestManagedRouteValidationRejectsDuplicateBeforeWrite(t *testing.T) {
	router, store, cleanup := setupTestRouter(t)
	defer cleanup()
	fake := newFakeManagedCaddy(`[]`)
	server := httptest.NewServer(fake)
	defer server.Close()
	connectManagedCaddy(t, router, server.URL)
	body := map[string]any{"domain": "same.example.com", "path": "/api", "handler_type": "reverse_proxy", "config": map[string]any{"upstreams": []string{"app:8080"}}}
	if w := requestJSON(t, router, http.MethodPost, "/api/routes", body); w.Code != http.StatusCreated {
		t.Fatalf("first create status %d: %s", w.Code, w.Body.String())
	}
	w := requestJSON(t, router, http.MethodPost, "/api/routes", body)
	if w.Code == http.StatusCreated || !bytes.Contains(w.Body.Bytes(), []byte("duplicates")) {
		t.Fatalf("duplicate should be rejected: %d %s", w.Code, w.Body.String())
	}
	routes, _ := store.ListRoutes()
	if len(routes) != 1 {
		t.Fatalf("duplicate reached storage: %+v", routes)
	}
}

func TestLocalDraftDoesNotContactCaddyBeforeSetup(t *testing.T) {
	router, store, cleanup := setupTestRouter(t)
	defer cleanup()
	w := requestJSON(t, router, http.MethodPost, "/api/routes", map[string]any{
		"domain": "draft.example.com", "handler_type": "reverse_proxy", "config": map[string]any{"upstreams": []string{"app:8080"}},
	})
	if w.Code != http.StatusCreated || !bytes.Contains(w.Body.Bytes(), []byte("local draft")) {
		t.Fatalf("draft response %d: %s", w.Code, w.Body.String())
	}
	routes, _ := store.ListRoutes()
	if len(routes) != 1 {
		t.Fatalf("local draft not persisted: %+v", routes)
	}
	fake := newFakeManagedCaddy(`[]`)
	server := httptest.NewServer(fake)
	defer server.Close()
	preview := connectManagedCaddy(t, router, server.URL)
	if len(preview.LocalDrafts) != 1 {
		t.Fatalf("setup preview omitted local draft: %+v", preview)
	}
	routes, _ = store.ListRoutes()
	if len(routes) != 1 || routes[0].Domain != "draft.example.com" {
		t.Fatalf("setup confirmation discarded local draft: %+v", routes)
	}
}

func TestEmptyCaddyCanBeExplicitlyBootstrapped(t *testing.T) {
	router, store, cleanup := setupTestRouter(t)
	defer cleanup()
	bootstrapped := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/config/":
			w.Header().Set("ETag", `"root-empty"`)
			_, _ = w.Write([]byte(`{}`))
		case r.Method == http.MethodPut && r.URL.Path == "/config/apps":
			if r.Header.Get("If-Match") != `"root-empty"` {
				w.WriteHeader(http.StatusPreconditionFailed)
				return
			}
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body["http"] == nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			bootstrapped = true
		case bootstrapped && r.Method == http.MethodGet && r.URL.Path == "/config/apps/http/servers/caddy-admin-ui/routes":
			w.Header().Set("ETag", `"routes-empty"`)
			_, _ = w.Write([]byte(`[]`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	w := requestJSON(t, router, http.MethodPost, "/api/setup/preview", setupRequest{URL: server.URL})
	if w.Code != http.StatusOK {
		t.Fatalf("preview status %d: %s", w.Code, w.Body.String())
	}
	var preview setupPreview
	_ = json.Unmarshal(w.Body.Bytes(), &preview)
	if !preview.CaddyEmpty || !preview.CanBootstrap || preview.SelectedServer != "caddy-admin-ui" {
		t.Fatalf("unexpected empty preview: %+v", preview)
	}
	w = requestJSON(t, router, http.MethodPost, "/api/setup/confirm", setupRequest{URL: server.URL, Server: preview.SelectedServer, PreviewToken: preview.PreviewToken})
	if w.Code != http.StatusOK || !bootstrapped {
		t.Fatalf("bootstrap status %d bootstrapped=%v: %s", w.Code, bootstrapped, w.Body.String())
	}
	cfg, _ := store.GetGlobalConfig()
	if !cfg.SetupComplete || cfg.ManagedServer != "caddy-admin-ui" {
		t.Fatalf("bootstrap ownership not saved: %+v", cfg)
	}
}

func TestSnapshotRestoreUsesCurrentETagAndRebuildsLocalState(t *testing.T) {
	router, store, cleanup := setupTestRouter(t)
	defer cleanup()
	fake := newFakeManagedCaddy(`[]`)
	server := httptest.NewServer(fake)
	defer server.Close()
	connectManagedCaddy(t, router, server.URL)
	w := requestJSON(t, router, http.MethodPost, "/api/routes", map[string]any{
		"domain": "restore.example.com", "handler_type": "reverse_proxy", "config": map[string]any{"upstreams": []string{"app:8080"}},
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create status %d: %s", w.Code, w.Body.String())
	}
	var created struct {
		Route struct {
			ID string `json:"id"`
		} `json:"route"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &created)
	w = requestJSON(t, router, http.MethodDelete, "/api/routes/"+created.Route.ID, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("delete status %d: %s", w.Code, w.Body.String())
	}
	snapshots, _ := store.ListSnapshots()
	var beforeDelete *storage.Snapshot
	for _, snapshot := range snapshots {
		if snapshot.Reason == "delete route" {
			beforeDelete = snapshot
			break
		}
	}
	if beforeDelete == nil {
		t.Fatalf("delete snapshot missing: %+v", snapshots)
	}
	w = requestJSON(t, router, http.MethodPost, "/api/snapshots/"+beforeDelete.ID+"/restore", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("restore status %d: %s", w.Code, w.Body.String())
	}
	routes, _ := store.ListRoutes()
	if len(routes) != 1 || routes[0].Domain != "restore.example.com" || routes[0].IsReadOnly() {
		t.Fatalf("snapshot did not restore owned route: %+v", routes)
	}
}

func TestSetupInitializesMissingRoutesWithoutChangingServerSettings(t *testing.T) {
	router, store, cleanup := setupTestRouter(t)
	defer cleanup()
	initialized := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/config/":
			w.Header().Set("ETag", `"root-1"`)
			_, _ = w.Write([]byte(`{"apps":{"http":{"servers":{"srv0":{"listen":[":9443"],"automatic_https":{"disable":true}}}}}}`))
		case r.Method == http.MethodGet && r.URL.Path == "/config/apps/http/servers/srv0":
			w.Header().Set("ETag", `"server-1"`)
			_, _ = w.Write([]byte(`{"listen":[":9443"],"automatic_https":{"disable":true}}`))
		case r.Method == http.MethodPatch && r.URL.Path == "/config/apps/http/servers/srv0":
			if r.Header.Get("If-Match") != `"server-1"` {
				w.WriteHeader(http.StatusPreconditionFailed)
				return
			}
			var body map[string]any
			_ = json.NewDecoder(r.Body).Decode(&body)
			listen := body["listen"].([]any)
			automaticHTTPS := body["automatic_https"].(map[string]any)
			if len(listen) != 1 || listen[0] != ":9443" || automaticHTTPS["disable"] != true || body["routes"] == nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			initialized = true
		case initialized && r.Method == http.MethodGet && r.URL.Path == "/config/apps/http/servers/srv0/routes":
			w.Header().Set("ETag", `"routes-1"`)
			_, _ = w.Write([]byte(`[]`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	w := requestJSON(t, router, http.MethodPost, "/api/setup/preview", setupRequest{URL: server.URL, Server: "srv0"})
	var preview setupPreview
	_ = json.Unmarshal(w.Body.Bytes(), &preview)
	w = requestJSON(t, router, http.MethodPost, "/api/setup/confirm", setupRequest{URL: server.URL, Server: "srv0", PreviewToken: preview.PreviewToken})
	if w.Code != http.StatusOK || !initialized {
		t.Fatalf("missing routes setup status %d initialized=%v: %s", w.Code, initialized, w.Body.String())
	}
	cfg, _ := store.GetGlobalConfig()
	if !cfg.SetupComplete {
		t.Fatalf("setup state missing: %+v", cfg)
	}
}
