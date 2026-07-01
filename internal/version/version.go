// Package version holds the application version string.
package version

// Version is injected at build time via
// -ldflags "-X github.com/ArtemStepanov/caddy-admin-ui/internal/version.Version=...".
var Version = "dev"
