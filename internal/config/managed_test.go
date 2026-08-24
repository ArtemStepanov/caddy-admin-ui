package config

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/ArtemStepanov/caddy-admin-ui/internal/storage"
)

func TestParseCaddyRoutesPreservesExternalJSONAndOrder(t *testing.T) {
	raw := json.RawMessage(`[{"match":[{"host":["one.example.com"],"method":["GET"]}],"handle":[{"handler":"custom","secret":"keep"}],"extension":{"nested":true}},{"@id":"caddy-admin-ui-route-9d87c74b-f335-4db4-b18f-d71550934295","match":[{"host":["two.example.com"]}],"handle":[{"handler":"reverse_proxy","upstreams":[{"dial":"app:8080"}]}],"terminal":true}]`)
	routes, err := ParseCaddyRoutes(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(routes) != 2 || routes[0].Position != 0 || routes[1].Position != 1 {
		t.Fatalf("route order lost: %+v", routes)
	}
	if !routes[0].IsReadOnly() || !bytes.Contains(routes[0].RawCaddyRoute, []byte(`"extension"`)) {
		t.Fatalf("external route was not preserved: %+v raw=%s", routes[0], routes[0].RawCaddyRoute)
	}
	if routes[1].IsReadOnly() || routes[1].ID != "9d87c74b-f335-4db4-b18f-d71550934295" {
		t.Fatalf("owned route not recognized: %+v", routes[1])
	}
}

func TestBuildCaddyRoutesPreservesReadOnlyAndMarksOwned(t *testing.T) {
	external := json.RawMessage(`{"handle":[{"handler":"custom","secret":"keep"}],"custom":true}`)
	routes := []*storage.Route{
		{ID: "external", Domain: "*", HandlerType: "unknown", Config: json.RawMessage(`{}`), Enabled: true, ReadOnly: true, SupportStatus: storage.SupportStatusUnsupportedReadOnly, ReadOnlyReason: "external", RawCaddyRoute: external},
		{ID: "owned", Domain: "app.example.com", HandlerType: "reverse_proxy", Config: json.RawMessage(`{"upstreams":["app:8080"]}`), Enabled: true},
	}
	built, err := BuildCaddyRoutes(routes, &storage.GlobalConfig{EnableEncode: true})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(built[0], external) {
		t.Fatalf("external route changed: %s", built[0])
	}
	var owned map[string]any
	if err := json.Unmarshal(built[1], &owned); err != nil {
		t.Fatal(err)
	}
	if owned["@id"] != ManagedRouteIDPrefix+"owned" {
		t.Fatalf("managed marker missing: %s", built[1])
	}
}

func TestValidateRoutesRejectsDuplicateOwnedMatch(t *testing.T) {
	routes := []*storage.Route{
		{ID: "one", Domain: "Example.com", Path: "api", HandlerType: "file_server", Config: json.RawMessage(`{}`), Enabled: true},
		{ID: "two", Domain: "example.com", Path: "/api", HandlerType: "file_server", Config: json.RawMessage(`{}`), Enabled: true},
	}
	if err := ValidateRoutesForBuild(routes); err == nil {
		t.Fatal("duplicate host/path should be rejected")
	}
}
