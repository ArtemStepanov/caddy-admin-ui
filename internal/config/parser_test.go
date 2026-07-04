package config

import (
	"encoding/json"
	"testing"

	"github.com/ArtemStepanov/caddy-admin-ui/internal/storage"
)

func TestParseCaddyConfig_Empty(t *testing.T) {
	cfg := &CaddyConfig{} // Empty config, not nil
	routes, err := ParseCaddyConfig(cfg)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(routes) != 0 {
		t.Errorf("Expected 0 routes, got %d", len(routes))
	}
}

func TestParseCaddyConfig_NoApps(t *testing.T) {
	cfg := &CaddyConfig{}
	routes, err := ParseCaddyConfig(cfg)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(routes) != 0 {
		t.Errorf("Expected 0 routes, got %d", len(routes))
	}
}

func TestParseCaddyConfig_ReverseProxy(t *testing.T) {
	cfg := &CaddyConfig{
		Apps: &Apps{
			HTTP: &HTTPApp{
				Servers: map[string]*Server{
					"srv0": {
						Listen: []string{":443"},
						Routes: []Route{
							{
								Match: []Match{
									{Host: []string{"example.com"}},
								},
								Handle: []Handler{
									{
										"handler": "reverse_proxy",
										"upstreams": []any{
											map[string]any{"dial": "localhost:8080"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	routes, err := ParseCaddyConfig(cfg)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(routes) != 1 {
		t.Fatalf("Expected 1 route, got %d", len(routes))
	}

	route := routes[0]
	if route.Domain != "example.com" {
		t.Errorf("Expected domain example.com, got %s", route.Domain)
	}
	if route.HandlerType != "reverse_proxy" {
		t.Errorf("Expected handler_type reverse_proxy, got %s", route.HandlerType)
	}
	if !route.Enabled {
		t.Error("Expected route to be enabled")
	}

	var cfg2 storage.ReverseProxyConfig
	json.Unmarshal(route.Config, &cfg2)
	if len(cfg2.Upstreams) != 1 || cfg2.Upstreams[0] != "localhost:8080" {
		t.Errorf("Expected upstreams [localhost:8080], got %v", cfg2.Upstreams)
	}
}

func TestParseCaddyConfig_FileServer(t *testing.T) {
	cfg := &CaddyConfig{
		Apps: &Apps{
			HTTP: &HTTPApp{
				Servers: map[string]*Server{
					"srv0": {
						Listen: []string{":443"},
						Routes: []Route{
							{
								Match: []Match{
									{Host: []string{"static.example.com"}},
								},
								Handle: []Handler{
									{
										"handler": "file_server",
										"root":    "/var/www",
										"browse":  map[string]any{},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	routes, err := ParseCaddyConfig(cfg)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	route := routes[0]
	if route.HandlerType != "file_server" {
		t.Errorf("Expected handler_type file_server, got %s", route.HandlerType)
	}

	var fsCfg storage.FileServerConfig
	json.Unmarshal(route.Config, &fsCfg)
	if fsCfg.Root != "/var/www" {
		t.Errorf("Expected root /var/www, got %s", fsCfg.Root)
	}
	if !fsCfg.Browse {
		t.Error("Expected browse to be true")
	}
}

func TestParseCaddyConfig_Redirect(t *testing.T) {
	cfg := &CaddyConfig{
		Apps: &Apps{
			HTTP: &HTTPApp{
				Servers: map[string]*Server{
					"srv0": {
						Listen: []string{":443"},
						Routes: []Route{
							{
								Match: []Match{
									{Host: []string{"old.example.com"}},
								},
								Handle: []Handler{
									{
										"handler":     "static_response",
										"status_code": float64(301), // JSON unmarshals to float64
										"headers": map[string]any{
											"Location": []any{"https://new.example.com"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	routes, err := ParseCaddyConfig(cfg)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	route := routes[0]
	if route.HandlerType != "redir" {
		t.Errorf("Expected handler_type redir, got %s", route.HandlerType)
	}

	var redirCfg storage.RedirectConfig
	json.Unmarshal(route.Config, &redirCfg)
	if redirCfg.To != "https://new.example.com" {
		t.Errorf("Expected to https://new.example.com, got %s", redirCfg.To)
	}
	if redirCfg.Code != 301 {
		t.Errorf("Expected code 301, got %d", redirCfg.Code)
	}
}

func TestParseCaddyConfig_Headers(t *testing.T) {
	cfg := &CaddyConfig{
		Apps: &Apps{
			HTTP: &HTTPApp{
				Servers: map[string]*Server{
					"srv0": {
						Listen: []string{":443"},
						Routes: []Route{
							{
								Match: []Match{
									{Host: []string{"example.com"}},
								},
								Handle: []Handler{
									{
										"handler": "headers",
										"response": map[string]any{
											"set": map[string]any{
												"X-Frame-Options": []any{"DENY"},
											},
											"add": map[string]any{
												"X-Custom": []any{"value"},
											},
											"delete": []any{"Server"},
										},
									},
									{
										"handler": "reverse_proxy",
										"upstreams": []any{
											map[string]any{"dial": "localhost:8080"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	routes, err := ParseCaddyConfig(cfg)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	route := routes[0]
	if route.Headers == nil {
		t.Fatal("Expected headers to be parsed")
	}
	if route.Headers.Set["X-Frame-Options"] != "DENY" {
		t.Errorf("Expected Set X-Frame-Options=DENY, got %v", route.Headers.Set)
	}
	if route.Headers.Add["X-Custom"] != "value" {
		t.Errorf("Expected Add X-Custom=value, got %v", route.Headers.Add)
	}
	if len(route.Headers.Delete) != 1 || route.Headers.Delete[0] != "Server" {
		t.Errorf("Expected Delete [Server], got %v", route.Headers.Delete)
	}
}

func TestParseCaddyConfig_CaddyfileFileServerWithRootVarsEditable(t *testing.T) {
	var cfg CaddyConfig
	if err := json.Unmarshal([]byte(`{
		"apps":{"http":{"servers":{"srv0":{"routes":[{
			"match":[{"host":["files.localhost"]}],
			"handle":[{"handler":"subroute","routes":[{"handle":[
				{"handler":"vars","root":"/usr/share/caddy"},
				{"handler":"file_server","hide":["/etc/caddy/Caddyfile"]}
			]}]}]
		}]}}}}
	}`), &cfg); err != nil {
		t.Fatal(err)
	}

	routes, err := ParseCaddyConfig(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(routes) != 1 {
		t.Fatalf("expected 1 route, got %d", len(routes))
	}
	route := routes[0]
	if route.ReadOnly || route.SupportStatus != storage.SupportStatusEditable || len(route.RawCaddyRoute) != 0 {
		t.Fatalf("file server should be editable: %+v raw=%s", route, route.RawCaddyRoute)
	}
	var fs storage.FileServerConfig
	if err := json.Unmarshal(route.Config, &fs); err != nil {
		t.Fatal(err)
	}
	if fs.Root != "/usr/share/caddy" || len(fs.Hide) != 1 || fs.Hide[0] != "/etc/caddy/Caddyfile" {
		t.Fatalf("bad file server config: %+v", fs)
	}
}

func TestParseCaddyConfig_ClassifiesSupportStatus(t *testing.T) {
	cfg := &CaddyConfig{
		Apps: &Apps{HTTP: &HTTPApp{Servers: map[string]*Server{
			"srv0": {Routes: []Route{
				{Match: []Match{{Host: []string{"editable.example.com"}}}, Handle: []Handler{{"handler": "reverse_proxy", "upstreams": []any{map[string]any{"dial": "localhost:8080"}}}}},
				{Match: []Match{{Host: []string{"partial.example.com"}}}, Handle: []Handler{{"handler": "crowdsec"}, {"handler": "reverse_proxy", "upstreams": []any{map[string]any{"dial": "localhost:8081"}}}}},
				{Match: []Match{{Host: []string{"unknown.example.com"}}}, Handle: []Handler{{"handler": "custom_handler"}}},
			}},
		}}},
	}

	routes, err := ParseCaddyConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if routes[0].SupportStatus != storage.SupportStatusEditable || routes[0].ReadOnly || routes[0].ReadOnlyReason != "" {
		t.Fatalf("editable classification failed: %+v", routes[0])
	}
	if routes[1].SupportStatus != storage.SupportStatusPartialReadOnly || !routes[1].ReadOnly || routes[1].ReadOnlyReason == "" {
		t.Fatalf("partial classification failed: %+v", routes[1])
	}
	if routes[2].SupportStatus != storage.SupportStatusUnsupportedReadOnly || !routes[2].ReadOnly || routes[2].ReadOnlyReason == "" {
		t.Fatalf("unsupported classification failed: %+v", routes[2])
	}
}

func TestParseCaddyConfig_UnknownHandler(t *testing.T) {
	cfg := &CaddyConfig{
		Apps: &Apps{
			HTTP: &HTTPApp{
				Servers: map[string]*Server{
					"srv0": {
						Listen: []string{":443"},
						Routes: []Route{
							{
								Match: []Match{
									{Host: []string{"example.com"}},
								},
								Handle: []Handler{
									{
										"handler": "acme_server",
										"ca":      "custom",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	routes, err := ParseCaddyConfig(cfg)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	route := routes[0]
	if route.HandlerType != "unknown" {
		t.Errorf("Expected handler_type unknown, got %s", route.HandlerType)
	}
	// Should preserve raw route
	if len(route.RawCaddyRoute) == 0 {
		t.Error("Expected RawCaddyRoute to be preserved for unknown handler")
	}
}

func TestParseCaddyConfig_GlobalRoute(t *testing.T) {
	// Route without host matcher (matches all)
	cfg := &CaddyConfig{
		Apps: &Apps{
			HTTP: &HTTPApp{
				Servers: map[string]*Server{
					"srv0": {
						Listen: []string{":443"},
						Routes: []Route{
							{
								Match: []Match{
									{Path: []string{"/health"}},
								},
								Handle: []Handler{
									{
										"handler": "reverse_proxy",
										"upstreams": []any{
											map[string]any{"dial": "localhost:8080"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	routes, err := ParseCaddyConfig(cfg)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	route := routes[0]
	if route.Domain != "*" {
		t.Errorf("Expected domain *, got %s", route.Domain)
	}
	if route.Path != "/health" {
		t.Errorf("Expected path /health, got %s", route.Path)
	}
}

func TestParseCaddyConfig_MultiHost(t *testing.T) {
	cfg := &CaddyConfig{
		Apps: &Apps{
			HTTP: &HTTPApp{
				Servers: map[string]*Server{
					"srv0": {
						Listen: []string{":443"},
						Routes: []Route{
							{
								Match: []Match{
									{Host: []string{"example.com", "www.example.com"}},
								},
								Handle: []Handler{
									{
										"handler": "reverse_proxy",
										"upstreams": []any{
											map[string]any{"dial": "localhost:8080"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	routes, err := ParseCaddyConfig(cfg)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	route := routes[0]
	if route.Domain != "example.com, www.example.com" {
		t.Errorf("Expected domain 'example.com, www.example.com', got %s", route.Domain)
	}
}

func TestRoundTrip_ReverseProxy(t *testing.T) {
	// Create a route, build Caddy config, parse it back
	original := &storage.Route{
		Domain:      "example.com",
		Path:        "/api",
		HandlerType: "reverse_proxy",
		Config:      json.RawMessage(`{"upstreams":["localhost:8080"]}`),
		Enabled:     true,
	}

	// Build Caddy config
	caddyConfig := BuildCaddyConfig([]*storage.Route{original}, nil)

	// Parse it back
	routes, err := ParseCaddyConfig(caddyConfig)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	if len(routes) != 1 {
		t.Fatalf("Expected 1 route, got %d", len(routes))
	}

	parsed := routes[0]
	if parsed.Domain != original.Domain {
		t.Errorf("Domain mismatch: expected %s, got %s", original.Domain, parsed.Domain)
	}
	if parsed.Path != original.Path {
		t.Errorf("Path mismatch: expected %s, got %s", original.Path, parsed.Path)
	}
	if parsed.HandlerType != original.HandlerType {
		t.Errorf("HandlerType mismatch: expected %s, got %s", original.HandlerType, parsed.HandlerType)
	}
}

func TestRoundTrip_FileServer(t *testing.T) {
	original := &storage.Route{
		Domain:      "static.example.com",
		HandlerType: "file_server",
		Config:      json.RawMessage(`{"root":"/var/www","browse":true}`),
		Enabled:     true,
	}

	caddyConfig := BuildCaddyConfig([]*storage.Route{original}, nil)
	routes, err := ParseCaddyConfig(caddyConfig)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	parsed := routes[0]
	if parsed.HandlerType != "file_server" {
		t.Errorf("Expected file_server, got %s", parsed.HandlerType)
	}

	var cfg storage.FileServerConfig
	json.Unmarshal(parsed.Config, &cfg)
	if cfg.Root != "/var/www" {
		t.Errorf("Expected root /var/www, got %s", cfg.Root)
	}
}

func TestRoundTrip_Redirect(t *testing.T) {
	original := &storage.Route{
		Domain:      "old.example.com",
		HandlerType: "redir",
		Config:      json.RawMessage(`{"to":"https://new.example.com","code":301}`),
		Enabled:     true,
	}

	// Build Caddy config
	caddyConfig := BuildCaddyConfig([]*storage.Route{original}, nil)

	// Simulate real round-trip through JSON (as it would happen with Caddy API)
	jsonBytes, err := json.Marshal(caddyConfig)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var parsedCaddyConfig CaddyConfig
	if err := json.Unmarshal(jsonBytes, &parsedCaddyConfig); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Parse the JSON-roundtripped config
	routes, err := ParseCaddyConfig(&parsedCaddyConfig)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	parsed := routes[0]
	if parsed.HandlerType != "redir" {
		t.Errorf("Expected redir, got %s", parsed.HandlerType)
	}

	var cfg storage.RedirectConfig
	json.Unmarshal(parsed.Config, &cfg)
	if cfg.To != "https://new.example.com" {
		t.Errorf("Expected to https://new.example.com, got %s", cfg.To)
	}
}

func TestParseCaddyConfig_Rewrite(t *testing.T) {
	cfg := &CaddyConfig{
		Apps: &Apps{
			HTTP: &HTTPApp{
				Servers: map[string]*Server{
					"srv0": {
						Listen: []string{":443"},
						Routes: []Route{
							{
								Match: []Match{
									{Host: []string{"example.com"}, Path: []string{"/api/*"}},
								},
								Handle: []Handler{
									{
										"handler":           "rewrite",
										"strip_path_prefix": "/api",
									},
									{
										"handler": "reverse_proxy",
										"upstreams": []any{
											map[string]any{"dial": "localhost:8080"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	routes, err := ParseCaddyConfig(cfg)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	route := routes[0]
	if route.StripPathPrefix != "/api" {
		t.Errorf("Expected StripPathPrefix /api, got %s", route.StripPathPrefix)
	}
	if route.HandlerType != "reverse_proxy" {
		t.Errorf("Expected handler_type reverse_proxy, got %s", route.HandlerType)
	}
}

func TestRoundTrip_StripPathPrefix(t *testing.T) {
	original := &storage.Route{
		Domain:          "example.com",
		Path:            "/api/*",
		HandlerType:     "reverse_proxy",
		Config:          json.RawMessage(`{"upstreams":["localhost:8080"]}`),
		StripPathPrefix: "/api",
		Enabled:         true,
	}

	// Build Caddy config
	caddyConfig := BuildCaddyConfig([]*storage.Route{original}, nil)

	// Simulate real round-trip through JSON
	jsonBytes, err := json.Marshal(caddyConfig)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var parsedCaddyConfig CaddyConfig
	if err := json.Unmarshal(jsonBytes, &parsedCaddyConfig); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Parse the JSON-roundtripped config
	routes, err := ParseCaddyConfig(&parsedCaddyConfig)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	parsed := routes[0]
	if parsed.StripPathPrefix != "/api" {
		t.Errorf("Expected StripPathPrefix /api, got %s", parsed.StripPathPrefix)
	}
	if parsed.HandlerType != "reverse_proxy" {
		t.Errorf("Expected handler_type reverse_proxy, got %s", parsed.HandlerType)
	}
	if parsed.Path != "/api/*" {
		t.Errorf("Expected path /api/*, got %s", parsed.Path)
	}
}

func TestParseCaddyConfig_SubrouteReverseProxy(t *testing.T) {
	// This is how Caddy generates config from Caddyfile — handlers are wrapped in subroute
	cfg := &CaddyConfig{
		Apps: &Apps{
			HTTP: &HTTPApp{
				Servers: map[string]*Server{
					"srv0": {
						Listen: []string{":443"},
						Routes: []Route{
							{
								Match: []Match{
									{Host: []string{"example.com"}},
								},
								Handle: []Handler{
									{
										"handler": "subroute",
										"routes": []any{
											map[string]any{
												"handle": []any{
													map[string]any{
														"handler": "reverse_proxy",
														"upstreams": []any{
															map[string]any{"dial": "localhost:8080"},
														},
													},
												},
											},
										},
									},
								},
								Terminal: true,
							},
						},
					},
				},
			},
		},
	}

	routes, err := ParseCaddyConfig(cfg)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(routes) != 1 {
		t.Fatalf("Expected 1 route, got %d", len(routes))
	}

	route := routes[0]
	if route.Domain != "example.com" {
		t.Errorf("Expected domain example.com, got %s", route.Domain)
	}
	if route.HandlerType != "reverse_proxy" {
		t.Errorf("Expected handler_type reverse_proxy, got %s", route.HandlerType)
	}

	var rpCfg storage.ReverseProxyConfig
	json.Unmarshal(route.Config, &rpCfg)
	if len(rpCfg.Upstreams) != 1 || rpCfg.Upstreams[0] != "localhost:8080" {
		t.Errorf("Expected upstreams [localhost:8080], got %v", rpCfg.Upstreams)
	}
}

func TestParseCaddyConfig_SubrouteWithHeaders(t *testing.T) {
	// Subroute containing both headers and reverse_proxy handlers
	cfg := &CaddyConfig{
		Apps: &Apps{
			HTTP: &HTTPApp{
				Servers: map[string]*Server{
					"srv0": {
						Listen: []string{":443"},
						Routes: []Route{
							{
								Match: []Match{
									{Host: []string{"example.com"}},
								},
								Handle: []Handler{
									{
										"handler": "subroute",
										"routes": []any{
											map[string]any{
												"handle": []any{
													map[string]any{
														"handler": "headers",
														"response": map[string]any{
															"set": map[string]any{
																"X-Frame-Options": []any{"DENY"},
															},
														},
													},
													map[string]any{
														"handler": "reverse_proxy",
														"upstreams": []any{
															map[string]any{"dial": "localhost:9000"},
														},
													},
												},
											},
										},
									},
								},
								Terminal: true,
							},
						},
					},
				},
			},
		},
	}

	routes, err := ParseCaddyConfig(cfg)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	route := routes[0]
	if route.HandlerType != "reverse_proxy" {
		t.Errorf("Expected handler_type reverse_proxy, got %s", route.HandlerType)
	}
	if route.Headers == nil {
		t.Fatal("Expected headers to be parsed")
	}
	if route.Headers.Set["X-Frame-Options"] != "DENY" {
		t.Errorf("Expected Set X-Frame-Options=DENY, got %v", route.Headers.Set)
	}
}

func TestParseCaddyConfig_SubrouteFromJSON(t *testing.T) {
	// Simulate actual Caddy API JSON response (as returned by Caddyfile-configured instances)
	raw := []byte(`{
		"apps": {
			"http": {
				"servers": {
					"srv0": {
						"listen": [":443"],
						"routes": [
							{
								"match": [{"host": ["app.example.com"]}],
								"handle": [{
									"handler": "subroute",
									"routes": [{
										"handle": [{
											"handler": "reverse_proxy",
											"upstreams": [{"dial": "localhost:3000"}]
										}]
									}]
								}],
								"terminal": true
							}
						]
					}
				}
			}
		}
	}`)

	var cfg CaddyConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	routes, err := ParseCaddyConfig(&cfg)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(routes) != 1 {
		t.Fatalf("Expected 1 route, got %d", len(routes))
	}

	route := routes[0]
	if route.HandlerType != "reverse_proxy" {
		t.Errorf("Expected handler_type reverse_proxy, got %s", route.HandlerType)
	}
	if route.Domain != "app.example.com" {
		t.Errorf("Expected domain app.example.com, got %s", route.Domain)
	}
}

func TestParseCaddyConfig_RealCaddyfileStructure(t *testing.T) {
	// Real Caddy config from Caddyfile: nested subroutes with crowdsec + appsec + reverse_proxy
	raw := []byte(`{
		"apps": {
			"http": {
				"servers": {
					"srv0": {
						"listen": [":443"],
						"routes": [
							{
								"match": [{"host": ["app.example.com"]}],
								"handle": [{
									"handler": "subroute",
									"routes": [{
										"handle": [{
											"handler": "subroute",
											"routes": [{
												"handle": [
													{"handler": "crowdsec"},
													{"handler": "appsec"},
													{
														"handler": "reverse_proxy",
														"headers": {
															"request": {
																"set": {
																	"X-Forwarded-For": ["{http.vars.client_ip}"],
																	"X-Real-Ip": ["{http.vars.client_ip}"]
																}
															}
														},
														"upstreams": [{"dial": "docker-hp:2111"}]
													}
												]
											}]
										}]
									}]
								}],
								"terminal": true
							}
						]
					}
				}
			}
		}
	}`)

	var cfg CaddyConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	routes, err := ParseCaddyConfig(&cfg)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	route := routes[0]
	if route.Domain != "app.example.com" {
		t.Errorf("Expected domain app.example.com, got %s", route.Domain)
	}
	if route.HandlerType != "reverse_proxy" {
		t.Errorf("Expected handler_type reverse_proxy, got %s", route.HandlerType)
	}

	var rpCfg storage.ReverseProxyConfig
	json.Unmarshal(route.Config, &rpCfg)
	if len(rpCfg.Upstreams) != 1 || rpCfg.Upstreams[0] != "docker-hp:2111" {
		t.Errorf("Expected upstreams [docker-hp:2111], got %v", rpCfg.Upstreams)
	}
	// Route has unmanaged middleware (crowdsec, appsec), should be read-only
	if !route.ReadOnly {
		t.Error("Expected ReadOnly to be true for route with unmanaged middleware")
	}
}

func TestParseCaddyConfig_NestedSubroutes(t *testing.T) {
	// Subroute inside subroute (both without matchers)
	cfg := &CaddyConfig{
		Apps: &Apps{
			HTTP: &HTTPApp{
				Servers: map[string]*Server{
					"srv0": {
						Listen: []string{":443"},
						Routes: []Route{
							{
								Match: []Match{
									{Host: []string{"example.com"}},
								},
								Handle: []Handler{
									{
										"handler": "subroute",
										"routes": []any{
											map[string]any{
												"handle": []any{
													map[string]any{
														"handler": "subroute",
														"routes": []any{
															map[string]any{
																"handle": []any{
																	map[string]any{
																		"handler": "file_server",
																		"root":    "/srv",
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
								Terminal: true,
							},
						},
					},
				},
			},
		},
	}

	routes, err := ParseCaddyConfig(cfg)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	route := routes[0]
	if route.HandlerType != "file_server" {
		t.Errorf("Expected handler_type file_server, got %s", route.HandlerType)
	}
}

func TestParseCaddyConfig_SubrouteWithMatchers(t *testing.T) {
	// Subroute with nested matchers (e.g. from handle_path) should stay "unknown"
	// to avoid losing routing semantics
	cfg := &CaddyConfig{
		Apps: &Apps{
			HTTP: &HTTPApp{
				Servers: map[string]*Server{
					"srv0": {
						Listen: []string{":443"},
						Routes: []Route{
							{
								Match: []Match{
									{Host: []string{"example.com"}},
								},
								Handle: []Handler{
									{
										"handler": "subroute",
										"routes": []any{
											map[string]any{
												"match": []any{
													map[string]any{
														"path": []any{"/api/*"},
													},
												},
												"handle": []any{
													map[string]any{
														"handler": "reverse_proxy",
														"upstreams": []any{
															map[string]any{"dial": "localhost:8080"},
														},
													},
												},
											},
										},
									},
								},
								Terminal: true,
							},
						},
					},
				},
			},
		},
	}

	routes, err := ParseCaddyConfig(cfg)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	route := routes[0]
	// Deep detection should find reverse_proxy even inside complex subroute
	if route.HandlerType != "reverse_proxy" {
		t.Errorf("Expected handler_type reverse_proxy (deep detected), got %s", route.HandlerType)
	}
	// Raw route should be preserved
	if len(route.RawCaddyRoute) == 0 {
		t.Error("Expected RawCaddyRoute to be preserved")
	}
	// Complex subroute should be read-only
	if !route.ReadOnly {
		t.Error("Expected ReadOnly to be true for complex subroute with matchers")
	}
}

func TestParseCaddyConfig_SubrouteWithUnmanagedHandler(t *testing.T) {
	// Subroute containing unmanaged handlers (crowdsec, appsec) alongside reverse_proxy.
	// Should detect reverse_proxy type — unmanaged handlers are preserved in RawCaddyRoute.
	cfg := &CaddyConfig{
		Apps: &Apps{
			HTTP: &HTTPApp{
				Servers: map[string]*Server{
					"srv0": {
						Listen: []string{":443"},
						Routes: []Route{
							{
								Match: []Match{
									{Host: []string{"example.com"}},
								},
								Handle: []Handler{
									{
										"handler": "subroute",
										"routes": []any{
											map[string]any{
												"handle": []any{
													map[string]any{
														"handler": "crowdsec",
													},
													map[string]any{
														"handler": "reverse_proxy",
														"upstreams": []any{
															map[string]any{"dial": "localhost:8080"},
														},
													},
												},
											},
										},
									},
								},
								Terminal: true,
							},
						},
					},
				},
			},
		},
	}

	routes, err := ParseCaddyConfig(cfg)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	route := routes[0]
	if route.HandlerType != "reverse_proxy" {
		t.Errorf("Expected handler_type reverse_proxy, got %s", route.HandlerType)
	}
	// RawCaddyRoute should preserve the full original including crowdsec
	if len(route.RawCaddyRoute) == 0 {
		t.Error("Expected RawCaddyRoute to be preserved")
	}
	// Route has unmanaged middleware (crowdsec), should be read-only
	if !route.ReadOnly {
		t.Error("Expected ReadOnly to be true for route with unmanaged middleware")
	}
}

func TestParseCaddyConfig_MultiRouteSubroute(t *testing.T) {
	// Subroute with multiple inner routes should stay "unknown" even without matchers
	// because route boundaries and terminal flags change Caddy's behavior
	cfg := &CaddyConfig{
		Apps: &Apps{
			HTTP: &HTTPApp{
				Servers: map[string]*Server{
					"srv0": {
						Listen: []string{":443"},
						Routes: []Route{
							{
								Match: []Match{
									{Host: []string{"example.com"}},
								},
								Handle: []Handler{
									{
										"handler": "subroute",
										"routes": []any{
											map[string]any{
												"handle": []any{
													map[string]any{
														"handler": "reverse_proxy",
														"upstreams": []any{
															map[string]any{"dial": "localhost:8080"},
														},
													},
												},
											},
											map[string]any{
												"handle": []any{
													map[string]any{
														"handler": "headers",
														"response": map[string]any{
															"set": map[string]any{
																"X-Debug": []any{"true"},
															},
														},
													},
												},
											},
										},
									},
								},
								Terminal: true,
							},
						},
					},
				},
			},
		},
	}

	routes, err := ParseCaddyConfig(cfg)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	route := routes[0]
	// Deep detection should find reverse_proxy even in multi-route subroute
	if route.HandlerType != "reverse_proxy" {
		t.Errorf("Expected handler_type reverse_proxy (deep detected), got %s", route.HandlerType)
	}
}

func TestRoundTrip_UnknownSubroutePreserved(t *testing.T) {
	// Complex subroute (with matchers) should be preserved as-is on round-trip
	raw := []byte(`{
		"apps": {
			"http": {
				"servers": {
					"srv0": {
						"listen": [":443"],
						"routes": [
							{
								"match": [{"host": ["example.com"]}],
								"handle": [{
									"handler": "subroute",
									"routes": [
										{
											"match": [{"path": ["/api/*"]}],
											"handle": [{
												"handler": "reverse_proxy",
												"upstreams": [{"dial": "localhost:8080"}]
											}]
										}
									]
								}],
								"terminal": true
							}
						]
					}
				}
			}
		}
	}`)

	var cfg CaddyConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Import — should detect reverse_proxy via deep search despite nested matchers
	routes, err := ParseCaddyConfig(&cfg)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	imported := routes[0]
	if imported.HandlerType != "reverse_proxy" {
		t.Fatalf("Expected reverse_proxy, got %s", imported.HandlerType)
	}

	// Export — should preserve the subroute handler (not strip it)
	exportedConfig := BuildCaddyConfig(routes, nil)

	exportedJSON, err := json.Marshal(exportedConfig)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var exportedCfg CaddyConfig
	if err := json.Unmarshal(exportedJSON, &exportedCfg); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	server := exportedCfg.Apps.HTTP.Servers["srv0"]
	if len(server.Routes) != 1 {
		t.Fatalf("Expected 1 exported route, got %d", len(server.Routes))
	}

	exportedRoute := server.Routes[0]
	// The subroute should be preserved (not stripped)
	if len(exportedRoute.Handle) == 0 {
		t.Fatal("Expected handlers to be preserved for unknown subroute route")
	}

	handlerType, _ := exportedRoute.Handle[0]["handler"].(string)
	if handlerType != "subroute" {
		t.Errorf("Expected preserved subroute handler, got %s", handlerType)
	}
}

func TestRoundTrip_TrivialSubrouteEditable(t *testing.T) {
	// A trivial subroute with only managed handlers should be fully editable:
	// RawCaddyRoute is cleared, so edits to upstreams/headers take effect on export.
	raw := []byte(`{
		"apps": {
			"http": {
				"servers": {
					"srv0": {
						"listen": [":443"],
						"routes": [
							{
								"match": [{"host": ["app.example.com"]}],
								"handle": [{
									"handler": "subroute",
									"routes": [{
										"handle": [{
											"handler": "reverse_proxy",
											"upstreams": [{"dial": "localhost:3000"}]
										}]
									}]
								}],
								"terminal": true
							}
						]
					}
				}
			}
		}
	}`)

	var cfg CaddyConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	routes, err := ParseCaddyConfig(&cfg)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	imported := routes[0]
	if imported.HandlerType != "reverse_proxy" {
		t.Fatalf("Expected reverse_proxy, got %s", imported.HandlerType)
	}
	// RawCaddyRoute should be cleared for fully-managed subroute
	if len(imported.RawCaddyRoute) != 0 {
		t.Error("Expected RawCaddyRoute to be cleared for trivial subroute (fully editable)")
	}
	// Trivial subroute with all managed handlers should NOT be read-only
	if imported.ReadOnly {
		t.Error("Expected ReadOnly to be false for trivial subroute (fully editable)")
	}

	// Simulate user editing upstreams
	imported.Config = json.RawMessage(`{"upstreams":["localhost:9999"]}`)

	// Export — should use normal buildRoute (not buildRouteMerged) since RawCaddyRoute is nil
	exportedConfig := BuildCaddyConfig(routes, nil)
	exportedJSON, _ := json.Marshal(exportedConfig)
	var exportedCfg CaddyConfig
	json.Unmarshal(exportedJSON, &exportedCfg)

	server := exportedCfg.Apps.HTTP.Servers["srv0"]
	exportedRoute := server.Routes[0]

	// Should be a flat reverse_proxy handler (not subroute), with the edited upstream
	if len(exportedRoute.Handle) != 1 {
		t.Fatalf("Expected 1 handler, got %d", len(exportedRoute.Handle))
	}
	handlerType, _ := exportedRoute.Handle[0]["handler"].(string)
	if handlerType != "reverse_proxy" {
		t.Errorf("Expected flat reverse_proxy, got %s", handlerType)
	}
	upstreams, _ := exportedRoute.Handle[0]["upstreams"].([]any)
	if len(upstreams) != 1 {
		t.Fatalf("Expected 1 upstream, got %d", len(upstreams))
	}
	upstream, _ := upstreams[0].(map[string]any)
	if dial, _ := upstream["dial"].(string); dial != "localhost:9999" {
		t.Errorf("Expected edited upstream localhost:9999, got %s", dial)
	}
}

func TestRoundTrip_SubrouteImport(t *testing.T) {
	// Simulate importing a subroute-wrapped route with unmanaged middleware (crowdsec),
	// then exporting it. The export should preserve the subroute structure, not duplicate handlers.
	raw := []byte(`{
		"apps": {
			"http": {
				"servers": {
					"srv0": {
						"listen": [":443"],
						"routes": [
							{
								"match": [{"host": ["app.example.com"]}],
								"handle": [{
									"handler": "subroute",
									"routes": [{
										"handle": [
											{"handler": "crowdsec"},
											{
												"handler": "reverse_proxy",
												"upstreams": [{"dial": "localhost:3000"}]
											}
										]
									}]
								}],
								"terminal": true
							}
						]
					}
				}
			}
		}
	}`)

	var cfg CaddyConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Import
	routes, err := ParseCaddyConfig(&cfg)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	if len(routes) != 1 {
		t.Fatalf("Expected 1 route, got %d", len(routes))
	}

	imported := routes[0]
	if imported.HandlerType != "reverse_proxy" {
		t.Fatalf("Expected reverse_proxy, got %s", imported.HandlerType)
	}

	// Export (rebuild)
	exportedConfig := BuildCaddyConfig(routes, nil)

	// Serialize and re-parse to verify clean round-trip
	exportedJSON, err := json.Marshal(exportedConfig)
	if err != nil {
		t.Fatalf("Failed to marshal exported config: %v", err)
	}

	var exportedCfg CaddyConfig
	if err := json.Unmarshal(exportedJSON, &exportedCfg); err != nil {
		t.Fatalf("Failed to unmarshal exported config: %v", err)
	}

	server := exportedCfg.Apps.HTTP.Servers["srv0"]
	if len(server.Routes) != 1 {
		t.Fatalf("Expected 1 exported route, got %d", len(server.Routes))
	}

	exportedRoute := server.Routes[0]
	// Should preserve the subroute handler, NOT add a duplicate reverse_proxy
	if len(exportedRoute.Handle) != 1 {
		t.Errorf("Expected 1 top-level handler (subroute), got %d (possible handler duplication)", len(exportedRoute.Handle))
	}

	handlerType, _ := exportedRoute.Handle[0]["handler"].(string)
	if handlerType != "subroute" {
		t.Errorf("Expected subroute handler preserved, got %s", handlerType)
	}
}

func TestParseCaddyConfig_SubrouteNonRedirectStaticResponse(t *testing.T) {
	// A static_response without a Location header (e.g. error page, websocket blocker)
	// should NOT be labeled "redir" — it should be "unknown" since we can't represent it.
	cfg := &CaddyConfig{
		Apps: &Apps{
			HTTP: &HTTPApp{
				Servers: map[string]*Server{
					"srv0": {
						Listen: []string{":443"},
						Routes: []Route{
							{
								Match: []Match{
									{Host: []string{"example.com"}},
								},
								Handle: []Handler{
									{
										"handler": "subroute",
										"routes": []any{
											map[string]any{
												"match": []any{
													map[string]any{
														"path": []any{"/blocked"},
													},
												},
												"handle": []any{
													map[string]any{
														"handler":     "static_response",
														"status_code": float64(403),
														"body":        "Forbidden",
													},
												},
											},
										},
									},
								},
								Terminal: true,
							},
						},
					},
				},
			},
		},
	}

	routes, err := ParseCaddyConfig(cfg)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	route := routes[0]
	// Non-redirect static_response should be "unknown", not "redir"
	if route.HandlerType == "redir" {
		t.Errorf("Non-redirect static_response should not be labeled 'redir', got %s", route.HandlerType)
	}
	if route.HandlerType != "unknown" {
		t.Errorf("Expected handler_type 'unknown' for non-redirect static_response, got %s", route.HandlerType)
	}
}

func TestParseCaddyConfig_SubrouteRedirectStaticResponse(t *testing.T) {
	// A static_response WITH a Location header should be detected as "redir" via deep search.
	cfg := &CaddyConfig{
		Apps: &Apps{
			HTTP: &HTTPApp{
				Servers: map[string]*Server{
					"srv0": {
						Listen: []string{":443"},
						Routes: []Route{
							{
								Match: []Match{
									{Host: []string{"old.example.com"}},
								},
								Handle: []Handler{
									{
										"handler": "subroute",
										"routes": []any{
											map[string]any{
												"match": []any{
													map[string]any{
														"path": []any{"/*"},
													},
												},
												"handle": []any{
													map[string]any{
														"handler":     "static_response",
														"status_code": float64(301),
														"headers": map[string]any{
															"Location": []any{"https://new.example.com{http.request.uri}"},
														},
													},
												},
											},
										},
									},
								},
								Terminal: true,
							},
						},
					},
				},
			},
		},
	}

	routes, err := ParseCaddyConfig(cfg)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	route := routes[0]
	if route.HandlerType != "redir" {
		t.Errorf("Expected handler_type 'redir' for redirect static_response, got %s", route.HandlerType)
	}
}

func TestParseCaddyConfig_SubrouteEmptyHandlers(t *testing.T) {
	// Subroute with an inner route that has no handlers — should fall through to unknown.
	cfg := &CaddyConfig{
		Apps: &Apps{
			HTTP: &HTTPApp{
				Servers: map[string]*Server{
					"srv0": {
						Listen: []string{":443"},
						Routes: []Route{
							{
								Match: []Match{
									{Host: []string{"example.com"}},
								},
								Handle: []Handler{
									{
										"handler": "subroute",
										"routes": []any{
											map[string]any{
												"handle": []any{},
											},
										},
									},
								},
								Terminal: true,
							},
						},
					},
				},
			},
		},
	}

	routes, err := ParseCaddyConfig(cfg)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	route := routes[0]
	if route.HandlerType != "unknown" {
		t.Errorf("Expected handler_type 'unknown' for empty subroute, got %s", route.HandlerType)
	}
}
