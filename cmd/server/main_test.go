package main

import "testing"

func TestIsLoopback(t *testing.T) {
	tests := []struct {
		addr string
		want bool
	}{
		{"127.0.0.1:3000", true},
		{"localhost:3000", true},
		{"::1:3000", true},
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
