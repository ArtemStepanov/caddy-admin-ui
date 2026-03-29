package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/ArtemStepanov/caddy-admin-ui/internal/api"
	"github.com/ArtemStepanov/caddy-admin-ui/internal/storage"
)

func main() {
	// Get configuration from environment
	dbPath := getEnv("DB_PATH", "./data/routes.db")
	caddyURL := getEnv("CADDY_ADMIN_URL", "http://localhost:2019")
	listenAddr := getEnv("LISTEN_ADDR", "127.0.0.1:3000")
	adminUser := os.Getenv("ADMIN_USER")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	log.Printf("Starting Caddy Orchestrator Lite")
	log.Printf("  Database: %s", dbPath)
	log.Printf("  Default Caddy Admin URL: %s", caddyURL)
	log.Printf("  Listen Address: %s", listenAddr)

	// Initialize storage
	store, err := storage.NewSQLiteStorage(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer store.Close()

	// Initialize Gin router
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// Fail fast if only one credential is set
	if (adminUser == "") != (adminPassword == "") {
		log.Fatalf("Both ADMIN_USER and ADMIN_PASSWORD must be set together")
	}

	// Optional Basic Auth (protects everything: UI + API)
	if adminUser != "" && adminPassword != "" {
		r.Use(gin.BasicAuth(gin.Accounts{adminUser: adminPassword}))
		log.Printf("  Basic Auth: enabled")
	} else if !isLoopback(listenAddr) {
		log.Printf("  WARNING: No authentication configured on a non-loopback address")
	}

	// Setup API routes
	api.SetupRoutes(r, store, caddyURL)

	// Serve static files (frontend)
	webDir := getEnv("WEB_DIR", "./web/dist")
	if _, err := os.Stat(webDir); err == nil {
		// Assets with long-term caching (1 year)
		assets := r.Group("/assets")
		assets.Use(func(c *gin.Context) {
			c.Header("Cache-Control", "public, max-age=31536000, immutable")
		})
		assets.Static("/", filepath.Join(webDir, "assets"))

		// SPA entry point with no-cache (always validate)
		serveIndex := func(c *gin.Context) {
			c.Header("Cache-Control", "no-cache")
			c.File(filepath.Join(webDir, "index.html"))
		}

		r.GET("/", serveIndex)
		r.StaticFile("/favicon.svg", filepath.Join(webDir, "favicon.svg"))
		r.NoRoute(serveIndex)
	} else {
		log.Printf("Warning: Web directory not found at %s", webDir)
	}

	log.Printf("Server starting on %s", listenAddr)
	if err := r.Run(listenAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// isLoopback checks if the listen address is bound to a loopback interface.
// An empty host (e.g. ":3000") means all interfaces, which is NOT loopback.
func isLoopback(addr string) bool {
	host := addr
	if i := strings.LastIndex(addr, ":"); i != -1 {
		host = addr[:i]
	}
	return host == "127.0.0.1" || host == "::1" || host == "localhost"
}
