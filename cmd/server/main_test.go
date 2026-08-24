package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestIsLoopback(t *testing.T) {
	tests := []struct {
		addr string
		want bool
	}{
		{"127.0.0.1:3000", true},
		{"127.0.0.2:3000", true},
		{"localhost:3000", true},
		{"[::1]:3000", true},
		{":3000", false},
		{"0.0.0.0:3000", false},
		{"192.168.1.1:3000", false},
		{"10.0.0.1:8080", false},
	}

	for _, tt := range tests {
		t.Run(tt.addr, func(t *testing.T) {
			got := isLoopback(tt.addr)
			if got != tt.want {
				t.Errorf("isLoopback(%q) = %v, want %v", tt.addr, got, tt.want)
			}
		})
	}
}

func TestCSRFProtection(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(csrfProtection())
	router.POST("/api/change", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	for _, tc := range []struct {
		name   string
		header string
		origin string
		want   int
	}{
		{name: "missing header", want: http.StatusForbidden},
		{name: "cross origin", header: "1", origin: "https://evil.example", want: http.StatusForbidden},
		{name: "same origin", header: "1", origin: "http://admin.example", want: http.StatusNoContent},
		{name: "non browser client", header: "1", want: http.StatusNoContent},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "http://admin.example/api/change", nil)
			req.Header.Set("X-Caddy-Admin-UI", tc.header)
			req.Header.Set("Origin", tc.origin)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			if w.Code != tc.want {
				t.Fatalf("status %d want %d: %s", w.Code, tc.want, w.Body.String())
			}
		})
	}
}
