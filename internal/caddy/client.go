package caddy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const maxResponseBytes = 16 << 20

// APIError preserves Caddy's HTTP status so callers can distinguish drift (412)
// from connectivity and validation failures.
type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Caddy Admin API returned status %d", e.StatusCode)
}

// ConfigResponse is a Caddy config fragment and its concurrency token.
type ConfigResponse struct {
	Body json.RawMessage
	ETag string
}

// Client is a small Caddy Admin API client.
type Client struct {
	adminURL   string
	httpClient *http.Client
}

// NewClient creates a new Caddy client.
func NewClient(adminURL string) *Client {
	return &Client{
		adminURL: strings.TrimRight(adminURL, "/"),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

// GetConfig retrieves the current Caddy configuration at the given path.
func (c *Client) GetConfig(path string) (json.RawMessage, error) {
	response, err := c.GetConfigWithETag(path)
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}

// GetConfigWithETag retrieves a config fragment and its Caddy ETag.
func (c *Client) GetConfigWithETag(path string) (*ConfigResponse, error) {
	resp, err := c.httpClient.Get(c.configURL(path))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Caddy: %w", err)
	}
	defer resp.Body.Close()
	body, err := readBody(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{StatusCode: resp.StatusCode, Body: string(body)}
	}
	return &ConfigResponse{Body: body, ETag: resp.Header.Get("ETag")}, nil
}

// PatchConfig replaces exactly one config value and guards the write with If-Match.
func (c *Client) PatchConfig(path string, config any, etag string) error {
	if etag == "" {
		return fmt.Errorf("ETag is required for Caddy config writes")
	}
	return c.writeConfig(http.MethodPatch, path, config, etag)
}

// PutConfig creates exactly one config value and guards the write with If-Match.
func (c *Client) PutConfig(path string, config any, etag string) error {
	if etag == "" {
		return fmt.Errorf("ETag is required for Caddy config writes")
	}
	return c.writeConfig(http.MethodPut, path, config, etag)
}

func (c *Client) writeConfig(method, path string, config any, etag string) error {
	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	req, err := http.NewRequest(method, c.configURL(path), bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("If-Match", etag)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to Caddy: %w", err)
	}
	defer resp.Body.Close()
	body, _ := readBody(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{StatusCode: resp.StatusCode, Body: string(body)}
	}
	return nil
}

// Health checks if Caddy is responsive.
func (c *Client) Health() error {
	_, err := c.GetConfigWithETag("")
	return err
}

// GetAdminURL returns the configured admin URL.
func (c *Client) GetAdminURL() string {
	return c.adminURL
}

func (c *Client) configURL(path string) string {
	url := c.adminURL + "/config/"
	if path != "" {
		url += strings.TrimLeft(path, "/")
	}
	return url
}

func readBody(reader io.Reader) ([]byte, error) {
	body, err := io.ReadAll(io.LimitReader(reader, maxResponseBytes+1))
	if err != nil {
		return nil, err
	}
	if len(body) > maxResponseBytes {
		return nil, fmt.Errorf("Caddy response exceeds %d bytes", maxResponseBytes)
	}
	return body, nil
}
