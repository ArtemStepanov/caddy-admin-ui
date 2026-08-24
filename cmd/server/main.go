package main

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/ArtemStepanov/caddy-admin-ui/internal/api"
	"github.com/ArtemStepanov/caddy-admin-ui/internal/storage"
	"github.com/ArtemStepanov/caddy-admin-ui/internal/version"
)

func main() {
	// Get configuration from environment
	dbPath := getEnv("DB_PATH", "./data/routes.db")
	caddyURL := getEnv("CADDY_ADMIN_URL", "http://localhost:2019")
	listenAddr := getEnv("LISTEN_ADDR", "127.0.0.1:3000")
	adminUser := os.Getenv("ADMIN_USER")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	log.Printf("Starting Caddy Admin UI %s", version.Version)
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
	r.Use(securityHeaders())
	r.Use(csrfProtection())
	// Liveness intentionally reveals no Caddy or application state and remains
	// available to the container health check when UI authentication is enabled.
	r.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	// Fail fast if only one credential is set
	if (adminUser == "") != (adminPassword == "") {
		log.Fatalf("Both ADMIN_USER and ADMIN_PASSWORD must be set together")
	}

	// Optional Basic Auth (protects everything: UI + API)
	if adminUser != "" && adminPassword != "" {
		r.Use(gin.BasicAuth(gin.Accounts{adminUser: adminPassword}))
		log.Printf("  Basic Auth: enabled")
	} else if !isLoopback(listenAddr) {
		if !strings.EqualFold(os.Getenv("ALLOW_INSECURE_NO_AUTH"), "true") {
			log.Fatalf("Refusing unauthenticated non-loopback listener; configure ADMIN_USER and ADMIN_PASSWORD or explicitly set ALLOW_INSECURE_NO_AUTH=true")
		}
		log.Printf("  WARNING: unauthenticated non-loopback access explicitly enabled")
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

func securityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Security-Policy", "default-src 'self'; style-src 'self'; img-src 'self' data:; object-src 'none'; base-uri 'none'; frame-ancestors 'none'")
		c.Header("Referrer-Policy", "no-referrer")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		c.Next()
	}
}

func csrfProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodHead || c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}
		if c.GetHeader("X-Caddy-Admin-UI") != "1" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "missing CSRF protection header"})
			return
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
		if origin := c.GetHeader("Origin"); origin != "" {
			parsed, err := url.Parse(origin)
			if err != nil || !strings.EqualFold(parsed.Host, c.Request.Host) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "cross-origin request rejected"})
				return
			}
		}
		c.Next()
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
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		host = addr
	}
	host = strings.Trim(host, "[]")
	if strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
