package caddy

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetConfigWithETagAndGuardedPatch(t *testing.T) {
	var method, ifMatch string
	var body []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/config/apps/http/servers/srv0/routes" {
			http.NotFound(w, r)
			return
		}
		if r.Method == http.MethodGet {
			w.Header().Set("ETag", `"abc"`)
			_, _ = w.Write([]byte(`[]`))
			return
		}
		method, ifMatch = r.Method, r.Header.Get("If-Match")
		var err error
		body, err = io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	response, err := client.GetConfigWithETag("apps/http/servers/srv0/routes")
	if err != nil || response.ETag != `"abc"` || string(response.Body) != "[]" {
		t.Fatalf("unexpected GET response: %+v err=%v", response, err)
	}
	if err := client.PatchConfig("apps/http/servers/srv0/routes", []any{}, response.ETag); err != nil {
		t.Fatal(err)
	}
	if method != http.MethodPatch || ifMatch != `"abc"` || !json.Valid(body) {
		t.Fatalf("bad patch method=%s if-match=%q body=%s", method, ifMatch, body)
	}
}

func TestClientDoesNotFollowRedirects(t *testing.T) {
	targetCalled := false
	target := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { targetCalled = true }))
	defer target.Close()
	source := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL, http.StatusFound)
	}))
	defer source.Close()
	_, err := NewClient(source.URL).GetConfig("")
	if err == nil || targetCalled {
		t.Fatalf("redirect should be rejected, err=%v targetCalled=%v", err, targetCalled)
	}
}

func TestAPIErrorPreservesStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusPreconditionFailed)
		_, _ = w.Write([]byte("stale"))
	}))
	defer server.Close()
	err := NewClient(server.URL).PatchConfig("routes", []any{}, `"old"`)
	apiErr, ok := err.(*APIError)
	if !ok || apiErr.StatusCode != http.StatusPreconditionFailed {
		t.Fatalf("expected APIError 412, got %#v", err)
	}
}

func TestConfigWritesRequireETag(t *testing.T) {
	client := NewClient("http://127.0.0.1:1")
	if err := client.PatchConfig("routes", []any{}, ""); err == nil {
		t.Fatal("PATCH without ETag should fail before sending")
	}
	if err := client.PutConfig("apps", map[string]any{}, ""); err == nil {
		t.Fatal("PUT without ETag should fail before sending")
	}
}
