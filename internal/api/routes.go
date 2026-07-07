package api

import (
	"github.com/gin-gonic/gin"

	"github.com/ArtemStepanov/caddy-admin-ui/internal/storage"
)

// SetupRoutes configures all API routes
func SetupRoutes(r *gin.Engine, store *storage.SQLiteStorage, defaultCaddyURL string) {
	h := NewHandler(store, defaultCaddyURL)

	api := r.Group("/api")
	{
		// Routes CRUD
		api.GET("/routes", h.ListRoutes)
		api.POST("/routes", h.CreateRoute)
		api.GET("/routes/:id", h.GetRoute)
		api.GET("/routes/:id/details", h.GetRouteDetails)
		api.PUT("/routes/:id", h.UpdateRoute)
		api.DELETE("/routes/:id", h.DeleteRoute)
		api.POST("/routes/:id/toggle", h.ToggleRoute)

		// Global config
		api.GET("/config", h.GetConfig)
		api.PUT("/config", h.UpdateConfig)

		// Caddy status
		api.GET("/status", h.GetStatus)
		api.POST("/sync", h.SyncToCaddy)
		api.POST("/test-connection", h.TestConnection)
		api.POST("/import-preview", h.PreviewImport)
		api.POST("/import", h.ImportFromCaddy)
	}
}
