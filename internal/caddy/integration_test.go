//go:build integration

package caddy

import (
	"encoding/json"
	"os"
	"testing"
)

func TestScopedPatchWithETagPreservesOtherServers(t *testing.T) {
	adminURL := os.Getenv("CADDY_INTEGRATION_URL")
	if adminURL == "" {
		t.Skip("CADDY_INTEGRATION_URL is not set")
	}
	client := NewClient(adminURL)
	path := "apps/http/servers/srv0/routes"
	before, err := client.GetConfigWithETag(path)
	if err != nil {
		t.Fatal(err)
	}
	if before.ETag == "" {
		t.Fatal("Caddy did not return an ETag")
	}
	routes := []json.RawMessage{json.RawMessage(`{"@id":"integration-route","handle":[{"handler":"static_response","body":"ok"}]}`)}
	if err := client.PatchConfig(path, routes, before.ETag); err != nil {
		t.Fatal(err)
	}
	if err := client.PatchConfig(path, []any{}, before.ETag); err == nil {
		t.Fatal("stale ETag unexpectedly succeeded")
	}
	root, err := client.GetConfig("")
	if err != nil {
		t.Fatal(err)
	}
	var config map[string]any
	if err := json.Unmarshal(root, &config); err != nil {
		t.Fatal(err)
	}
	apps := config["apps"].(map[string]any)
	httpApp := apps["http"].(map[string]any)
	servers := httpApp["servers"].(map[string]any)
	if _, ok := servers["untouched"]; !ok {
		t.Fatal("scoped patch removed an unrelated server")
	}
	logging := config["logging"].(map[string]any)
	if logging["logs"] == nil {
		t.Fatal("scoped patch removed unrelated logging config")
	}
}
