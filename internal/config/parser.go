package config

import (
	"encoding/json"
	"strings"

	"github.com/ArtemStepanov/caddy-admin-ui/internal/storage"
	"github.com/google/uuid"
)

// ParseCaddyConfig converts a Caddy configuration to a list of storage routes
func ParseCaddyConfig(cfg *CaddyConfig) ([]*storage.Route, error) {
	var routes []*storage.Route

	if cfg.Apps == nil || cfg.Apps.HTTP == nil || cfg.Apps.HTTP.Servers == nil {
		return routes, nil
	}

	for _, server := range cfg.Apps.HTTP.Servers {
		for _, caddyRoute := range server.Routes {
			parsedRoute, err := parseRoute(caddyRoute)
			if err != nil {
				// Log error but continue? Or skip?
				// For now, we'll skip invalid routes but maybe we should still import them as "raw"
				// If we fail to parse, let's treat it as an unknown handler type
				parsedRoute = createRawRoute(caddyRoute)
			}
			routes = append(routes, parsedRoute)
		}
	}

	return routes, nil
}

func parseRoute(r Route) (*storage.Route, error) {
	// marshal raw route first
	rawJSON, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	storageRoute := &storage.Route{
		ID:            uuid.New().String(),
		Enabled:       true,
		RawCaddyRoute: rawJSON,
	}

	// 1. Extract Matchers
	if len(r.Match) > 0 {
		// We only look at the first match block for simplicity,
		// as our builder only creates one.
		match := r.Match[0]

		if len(match.Host) > 0 {
			storageRoute.Domain = strings.Join(match.Host, ", ")
		} else {
			// No host matcher -> Global route
			storageRoute.Domain = "*"
		}

		if len(match.Path) > 0 {
			storageRoute.Path = match.Path[0] // We only support single path in UI model for now
		}
	} else {
		// No matchers -> Global route matching everything
		storageRoute.Domain = "*"
	}

	// 2. Extract Handlers
	// We look for known handlers: reverse_proxy, file_server, static_response (redir)
	// We also look for headers and encode

	// Collect all handlers, unwrapping subroutes so we can find the actual handlers
	// Caddy wraps handlers in "subroute" when configured via Caddyfile
	allHandlers := flattenHandlers(r.Handle)

	var mainHandlerFound bool

	for _, h := range allHandlers {
		handlerType, ok := h["handler"].(string)
		if !ok {
			continue
		}

		switch handlerType {
		case "reverse_proxy":
			if mainHandlerFound {
				continue
			} // Only one main handler supported

			cfg, err := parseReverseProxy(h)
			if err == nil {
				storageRoute.HandlerType = "reverse_proxy"
				storageRoute.Config, _ = json.Marshal(cfg)
				mainHandlerFound = true
			}

		case "file_server":
			if mainHandlerFound {
				continue
			}

			cfg, err := parseFileServer(h)
			if err == nil {
				storageRoute.HandlerType = "file_server"
				storageRoute.Config, _ = json.Marshal(cfg)
				mainHandlerFound = true
			}

		case "static_response":
			// Check if it's a redirect (has Location header)
			if headers, ok := h["headers"].(map[string]any); ok {
				if _, hasLoc := headers["Location"]; hasLoc {
					if mainHandlerFound {
						continue
					}

					cfg, err := parseRedirect(h)
					if err == nil {
						storageRoute.HandlerType = "redir"
						storageRoute.Config, _ = json.Marshal(cfg)
						mainHandlerFound = true
					}
				}
			}

		case "headers":
			// Parse headers
			cfg, err := parseHeaders(h)
			if err == nil {
				storageRoute.Headers = cfg
			}

		case "rewrite":
			// Parse rewrite handler for strip_path_prefix
			if prefix, ok := h["strip_path_prefix"].(string); ok {
				storageRoute.StripPathPrefix = prefix
			}

		case "encode":
			// We just ignore encode handler as it's global setting in our model usually,
			// or implied. But if we want to support per-route encode, we'd need to add it to model.
			// For now, we ignore it, but it's preserved in RawCaddyRoute.
		}
	}

	if mainHandlerFound {
		// If we successfully extracted all handlers via flatten, check whether
		// the route can be fully rebuilt from our model (all handlers are managed).
		// If so, clear RawCaddyRoute so the route goes through normal buildRoute
		// (editable) instead of buildRouteMerged (read-only preservation).
		if allHandlersManaged(allHandlers) {
			storageRoute.RawCaddyRoute = nil
		}
	} else {
		// Flattening didn't work (complex subroute). Do a deep search through
		// the entire handler tree to at least identify the handler type for display.
		// The original JSON is preserved in RawCaddyRoute, so it will be synced back as is.
		if detected := detectDeepHandlerType(r.Handle); detected != "" {
			// Normalize: Caddy uses "static_response" but our model uses "redir"
			if detected == "static_response" {
				detected = "redir"
			}
			storageRoute.HandlerType = detected
		} else {
			storageRoute.HandlerType = "unknown"
		}
		storageRoute.Config = json.RawMessage("{}")
	}

	return storageRoute, nil
}

// flattenHandlers unwraps simple subroute handlers to extract the actual handlers within.
// Caddy wraps handlers in "subroute" when configured via Caddyfile.
// Only unwraps subroutes whose inner routes have no matchers (simple wrappers).
// Subroutes with nested matchers (e.g. from handle_path) are left as-is to avoid
// losing routing semantics.
func flattenHandlers(handlers []Handler) []Handler {
	var result []Handler
	for _, h := range handlers {
		handlerType, _ := h["handler"].(string)
		if handlerType == "subroute" {
			if nested := extractSimpleSubrouteHandlers(h); nested != nil {
				result = append(result, nested...)
				continue
			}
		}
		result = append(result, h)
	}
	return result
}

// extractSimpleSubrouteHandlers extracts handlers from a subroute only if it
// contains exactly one inner route with no match conditions and no terminal flag.
// This covers the common Caddyfile pattern where handlers are wrapped in a trivial
// subroute. Returns nil for complex subroutes (multiple routes, matchers, terminal).
//
// Note: this is used for type DETECTION only. The original structure is preserved
// in RawCaddyRoute. Unrecognized handlers (crowdsec, appsec, etc.) are extracted
// here but simply ignored by the caller's switch statement.
func extractSimpleSubrouteHandlers(h Handler) []Handler {
	routes, ok := h["routes"].([]any)
	if !ok || len(routes) != 1 {
		return nil
	}

	routeMap, ok := routes[0].(map[string]any)
	if !ok {
		return nil
	}

	// If the inner route has matchers, this is a complex subroute — don't flatten
	if match, hasMatch := routeMap["match"]; hasMatch {
		if matchSlice, ok := match.([]any); ok && len(matchSlice) > 0 {
			return nil
		}
	}

	// If the inner route has terminal flag, don't flatten (route-level flow control)
	if terminal, ok := routeMap["terminal"].(bool); ok && terminal {
		return nil
	}

	nestedHandlers, ok := routeMap["handle"].([]any)
	if !ok {
		return nil
	}

	var result []Handler
	for _, nh := range nestedHandlers {
		if nhMap, ok := nh.(map[string]any); ok {
			// Recursively flatten in case of nested subroutes
			result = append(result, flattenHandlers([]Handler{nhMap})...)
		}
	}
	return result
}

// allHandlersManaged returns true if every handler in the list is a type
// we can fully represent and rebuild in our model.
func allHandlersManaged(handlers []Handler) bool {
	managed := map[string]bool{
		"reverse_proxy":   true,
		"file_server":     true,
		"static_response": true,
		"headers":         true,
		"encode":          true,
		"rewrite":         true,
	}
	for _, h := range handlers {
		hType, _ := h["handler"].(string)
		if !managed[hType] {
			return false
		}
	}
	return true
}

// detectDeepHandlerType recursively searches through all handlers (including
// nested subroutes) to find a main handler type. Used as a fallback when
// flattenHandlers can't extract handlers from complex subroute structures.
// Prioritizes reverse_proxy and file_server over static_response, since
// static_response is often used as auxiliary middleware (e.g. websocket blocking).
func detectDeepHandlerType(handlers []Handler) string {
	var found []string
	collectDeepHandlerTypes(handlers, &found)

	// Prefer reverse_proxy > file_server > static_response
	for _, t := range found {
		if t == "reverse_proxy" {
			return t
		}
	}
	for _, t := range found {
		if t == "file_server" {
			return t
		}
	}
	for _, t := range found {
		if t == "static_response" {
			return t
		}
	}
	return ""
}

func collectDeepHandlerTypes(handlers []Handler, found *[]string) {
	mainTypes := map[string]bool{
		"reverse_proxy":   true,
		"file_server":     true,
		"static_response": true,
	}

	for _, h := range handlers {
		hType, _ := h["handler"].(string)
		if mainTypes[hType] {
			*found = append(*found, hType)
		}
		if hType == "subroute" {
			if routes, ok := h["routes"].([]any); ok {
				for _, r := range routes {
					if routeMap, ok := r.(map[string]any); ok {
						if nested, ok := routeMap["handle"].([]any); ok {
							var nestedHandlers []Handler
							for _, nh := range nested {
								if nhMap, ok := nh.(map[string]any); ok {
									nestedHandlers = append(nestedHandlers, nhMap)
								}
							}
							collectDeepHandlerTypes(nestedHandlers, found)
						}
					}
				}
			}
		}
	}
}

func createRawRoute(r Route) *storage.Route {
	rawJSON, _ := json.Marshal(r)
	return &storage.Route{
		ID:            uuid.New().String(),
		Domain:        "UNKNOWN",
		HandlerType:   "unknown",
		Enabled:       true,
		RawCaddyRoute: rawJSON,
		Config:        json.RawMessage("{}"),
	}
}

func parseReverseProxy(h Handler) (*storage.ReverseProxyConfig, error) {
	cfg := &storage.ReverseProxyConfig{}

	// Upstreams
	if upstreams, ok := h["upstreams"].([]any); ok {
		for _, u := range upstreams {
			if uMap, ok := u.(map[string]any); ok {
				if dial, ok := uMap["dial"].(string); ok {
					cfg.Upstreams = append(cfg.Upstreams, dial)
				}
			}
		}
	}

	// Headers
	if headers, ok := h["headers"].(map[string]any); ok {
		if req, ok := headers["request"].(map[string]any); ok {
			if set, ok := req["set"].(map[string]any); ok {
				cfg.Headers = make(map[string]string)
				for k, v := range set {
					// v is usually []any or []string
					if vSlice, ok := v.([]any); ok && len(vSlice) > 0 {
						if vStr, ok := vSlice[0].(string); ok {
							cfg.Headers[k] = vStr
						}
					}
				}
			}
		}
	}

	// Load balancing
	if lb, ok := h["load_balancing"].(map[string]any); ok {
		if sel, ok := lb["selection_policy"].(map[string]any); ok {
			if policy, ok := sel["policy"].(string); ok {
				cfg.LoadBalancing = policy
			}
		}
	}

	return cfg, nil
}

func parseFileServer(h Handler) (*storage.FileServerConfig, error) {
	cfg := &storage.FileServerConfig{}

	if root, ok := h["root"].(string); ok {
		cfg.Root = root
	}

	if _, ok := h["browse"]; ok {
		cfg.Browse = true
	}

	if index, ok := h["index_names"].([]any); ok {
		for _, i := range index {
			if s, ok := i.(string); ok {
				cfg.Index = append(cfg.Index, s)
			}
		}
	}

	// Precompressed
	if _, ok := h["precompressed"]; ok {
		cfg.Precompressed = true
	}

	return cfg, nil
}

func parseRedirect(h Handler) (*storage.RedirectConfig, error) {
	cfg := &storage.RedirectConfig{}

	if codeStr, ok := h["status_code"].(float64); ok { // json unmarshals numbers to float64
		cfg.Code = int(codeStr)
	}

	if headers, ok := h["headers"].(map[string]any); ok {
		if loc, ok := headers["Location"].([]any); ok && len(loc) > 0 {
			if s, ok := loc[0].(string); ok {
				cfg.To = s
			}
		}
	}

	return cfg, nil
}

func parseHeaders(h Handler) (*storage.HeaderConfig, error) {
	cfg := &storage.HeaderConfig{}

	if resp, ok := h["response"].(map[string]any); ok {
		// Set
		if set, ok := resp["set"].(map[string]any); ok {
			cfg.Set = make(map[string]string)
			for k, v := range set {
				if vSlice, ok := v.([]any); ok && len(vSlice) > 0 {
					if vStr, ok := vSlice[0].(string); ok {
						cfg.Set[k] = vStr
					}
				}
			}
		}

		// Add
		if add, ok := resp["add"].(map[string]any); ok {
			cfg.Add = make(map[string]string)
			for k, v := range add {
				if vSlice, ok := v.([]any); ok && len(vSlice) > 0 {
					if vStr, ok := vSlice[0].(string); ok {
						cfg.Add[k] = vStr
					}
				}
			}
		}

		// Delete
		if del, ok := resp["delete"].([]any); ok {
			for _, d := range del {
				if s, ok := d.(string); ok {
					cfg.Delete = append(cfg.Delete, s)
				}
			}
		}
	}

	return cfg, nil
}
