package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"runtime"
	"time"

	"github.com/google/uuid"
)

const (
	// DefaultBaseURL is the base URL for the Pantheon API
	DefaultBaseURL = "https://terminus.pantheon.io:443/api"

	// DefaultTimeout is the default timeout for HTTP requests
	DefaultTimeout = 86400 * time.Second

	// MaxRetries is the maximum number of retry attempts
	MaxRetries = 5

	// InitialBackoff is the initial backoff duration for retries
	InitialBackoff = 1 * time.Second
)

// Client is the HTTP client for the Pantheon API
type Client struct {
	baseURL    string
	httpClient *http.Client
	userAgent  string
	token      string
	logger     Logger
}

// Logger is an interface for logging
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// NewClient creates a new API client
func NewClient(options ...ClientOption) *Client {
	c := &Client{
		baseURL: DefaultBaseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		userAgent: fmt.Sprintf("Terminus-Go/0.0.0 (go_version=%s; os=%s; arch=%s)",
			runtime.Version(), runtime.GOOS, runtime.GOARCH),
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}

// WithBaseURL sets a custom base URL
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithToken sets the authentication token
func WithToken(token string) ClientOption {
	return func(c *Client) {
		c.token = token
	}
}

// WithLogger sets a custom logger
func WithLogger(logger Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// SetToken updates the authentication token
func (c *Client) SetToken(token string) {
	c.token = token
}

// Request makes an HTTP request to the API with retry logic
func (c *Client) Request(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	fullURL := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Add trace ID
	traceID := uuid.New().String()
	req.Header.Set("X-Pantheon-Trace-Id", traceID)

	if c.logger != nil {
		c.logger.Debug("API Request: %s %s (trace: %s)", method, fullURL, traceID)
	}

	// Execute request with retry logic
	return c.doWithRetry(req)
}

// doWithRetry executes an HTTP request with exponential backoff retry logic
func (c *Client) doWithRetry(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for attempt := 0; attempt <= MaxRetries; attempt++ {
		// Clone the request body for retries
		var bodyReader io.Reader
		if req.Body != nil {
			bodyBytes, readErr := io.ReadAll(req.Body)
			if readErr != nil {
				return nil, fmt.Errorf("failed to read request body: %w", readErr)
			}
			_ = req.Body.Close()
			bodyReader = bytes.NewReader(bodyBytes)
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		resp, err = c.httpClient.Do(req)

		if err == nil && !shouldRetry(resp.StatusCode) {
			// Success or non-retriable error
			return resp, nil
		}

		if err != nil && c.logger != nil {
			c.logger.Warn("Request failed (attempt %d/%d): %v", attempt+1, MaxRetries+1, err)
		} else if c.logger != nil {
			c.logger.Warn("Request returned %d (attempt %d/%d)", resp.StatusCode, attempt+1, MaxRetries+1)
		}

		// Don't sleep on the last attempt
		if attempt < MaxRetries {
			backoff := time.Duration(math.Pow(2, float64(attempt))) * InitialBackoff
			if c.logger != nil {
				c.logger.Debug("Retrying after %v", backoff)
			}
			time.Sleep(backoff)

			// Restore body for retry
			if bodyReader != nil {
				req.Body = io.NopCloser(bodyReader)
			}
		}

		// Close failed response body
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}

	if err != nil {
		return nil, fmt.Errorf("request failed after %d attempts: %w", MaxRetries+1, err)
	}

	return resp, fmt.Errorf("request failed with status %d after %d attempts", resp.StatusCode, MaxRetries+1)
}

// shouldRetry determines if a status code should trigger a retry
func shouldRetry(statusCode int) bool {
	// Retry on 5xx server errors and 429 rate limiting
	return statusCode >= 500 || statusCode == 429
}

// Get makes a GET request
func (c *Client) Get(ctx context.Context, path string) (*http.Response, error) {
	return c.Request(ctx, http.MethodGet, path, nil)
}

// Post makes a POST request
func (c *Client) Post(ctx context.Context, path string, body interface{}) (*http.Response, error) {
	return c.Request(ctx, http.MethodPost, path, body)
}

// Put makes a PUT request
func (c *Client) Put(ctx context.Context, path string, body interface{}) (*http.Response, error) {
	return c.Request(ctx, http.MethodPut, path, body)
}

// Delete makes a DELETE request
func (c *Client) Delete(ctx context.Context, path string) (*http.Response, error) {
	return c.Request(ctx, http.MethodDelete, path, nil)
}

// Patch makes a PATCH request
func (c *Client) Patch(ctx context.Context, path string, body interface{}) (*http.Response, error) {
	return c.Request(ctx, http.MethodPatch, path, body)
}

// GetPaged makes paginated GET requests and returns all results
func (c *Client) GetPaged(ctx context.Context, basePath string) ([]json.RawMessage, error) {
	var allResults []json.RawMessage
	page := 1
	limit := 100

	for {
		// Build URL with pagination parameters
		u, err := url.Parse(c.baseURL + basePath)
		if err != nil {
			return nil, fmt.Errorf("failed to parse URL: %w", err)
		}

		q := u.Query()
		q.Set("limit", fmt.Sprintf("%d", limit))
		q.Set("page", fmt.Sprintf("%d", page))
		u.RawQuery = q.Encode()

		path := u.Path + "?" + u.RawQuery
		// Remove baseURL from path since Request will add it back
		if len(path) > len(c.baseURL) {
			path = path[len(c.baseURL):]
		}

		resp, err := c.Get(ctx, path)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
		}

		var results []json.RawMessage
		if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
			_ = resp.Body.Close()
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		_ = resp.Body.Close()

		if len(results) == 0 {
			break
		}

		allResults = append(allResults, results...)

		// If we got fewer results than the limit, we're done
		if len(results) < limit {
			break
		}

		page++
	}

	return allResults, nil
}

// DecodeResponse decodes a JSON response into a target struct
func DecodeResponse(resp *http.Response, target interface{}) error {
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return &Error{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

// Error represents an API error response
type Error struct {
	StatusCode int
	Message    string
}

func (e *Error) Error() string {
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

// IsNotFound returns true if the error is a 404 Not Found
func IsNotFound(err error) bool {
	if apiErr, ok := err.(*Error); ok {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}

// IsConflict returns true if the error is a 409 Conflict
func IsConflict(err error) bool {
	if apiErr, ok := err.(*Error); ok {
		return apiErr.StatusCode == http.StatusConflict
	}
	return false
}
