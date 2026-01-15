package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestEnvironmentsService_GetMetrics_WithEnvironment(t *testing.T) {
	testSiteID := "12345678-1234-1234-1234-123456789abc"
	testEnvID := "dev"

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/sites/" + testSiteID + "/environments/" + testEnvID + "/traffic"
		if r.URL.Path != expectedPath {
			t.Errorf("unexpected path: %s, expected: %s", r.URL.Path, expectedPath)
		}

		// Check duration parameter
		duration := r.URL.Query().Get("duration")
		if duration != "28d" {
			t.Errorf("expected duration '28d', got '%s'", duration)
		}

		// Return mock API response matching Pantheon API structure
		response := map[string]interface{}{
			"timeseries": []map[string]interface{}{
				{
					"timestamp":    1766016000,
					"visits":       100,
					"pages_served": 500,
					"cache_hits":   400,
					"cache_misses": 100,
				},
				{
					"timestamp":    1766102400,
					"visits":       0,
					"pages_served": 0,
					"cache_hits":   0,
					"cache_misses": 0,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	envsService := NewEnvironmentsService(client)

	metrics, err := envsService.GetMetrics(context.Background(), testSiteID, testEnvID, "28d")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(metrics) != 2 {
		t.Errorf("expected 2 metrics, got %d", len(metrics))
	}

	// Check first metric (with traffic)
	if metrics[0].Visits != 100 {
		t.Errorf("expected visits 100, got %d", metrics[0].Visits)
	}
	if metrics[0].PagesServed != 500 {
		t.Errorf("expected pages_served 500, got %d", metrics[0].PagesServed)
	}
	if metrics[0].CacheHits != 400 {
		t.Errorf("expected cache_hits 400, got %d", metrics[0].CacheHits)
	}
	if metrics[0].CacheHitRatio != "80.00%" {
		t.Errorf("expected cache_hit_ratio '80.00%%', got '%s'", metrics[0].CacheHitRatio)
	}

	// Check second metric (zero traffic - should show "--")
	if metrics[1].CacheHitRatio != "--" {
		t.Errorf("expected cache_hit_ratio '--' for zero traffic, got '%s'", metrics[1].CacheHitRatio)
	}
}

func TestEnvironmentsService_GetMetrics_SiteLevel(t *testing.T) {
	testSiteID := "12345678-1234-1234-1234-123456789abc"

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/sites/" + testSiteID + "/traffic"
		if r.URL.Path != expectedPath {
			t.Errorf("unexpected path: %s, expected: %s", r.URL.Path, expectedPath)
		}

		// Return mock API response
		response := map[string]interface{}{
			"timeseries": []map[string]interface{}{
				{
					"timestamp":    1766016000,
					"visits":       200,
					"pages_served": 1000,
					"cache_hits":   900,
					"cache_misses": 100,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	envsService := NewEnvironmentsService(client)

	// Empty envID means site-level metrics
	metrics, err := envsService.GetMetrics(context.Background(), testSiteID, "", "12w")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(metrics) != 1 {
		t.Errorf("expected 1 metric, got %d", len(metrics))
	}

	if metrics[0].Visits != 200 {
		t.Errorf("expected visits 200, got %d", metrics[0].Visits)
	}
	if metrics[0].CacheHitRatio != "90.00%" {
		t.Errorf("expected cache_hit_ratio '90.00%%', got '%s'", metrics[0].CacheHitRatio)
	}
}

func TestEnvironmentsService_GetMetrics_WithSiteName(t *testing.T) {
	testSiteName := "my-site"
	testSiteID := "12345678-1234-1234-1234-123456789abc"
	testEnvID := "dev"

	requestCount := 0

	// Create a test server that handles both site lookup and metrics
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		if r.URL.Path == "/site-names/"+testSiteName {
			// Site name lookup
			response := map[string]interface{}{
				"id": testSiteID,
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(response)
			return
		}

		expectedPath := "/sites/" + testSiteID + "/environments/" + testEnvID + "/traffic"
		if r.URL.Path != expectedPath {
			t.Errorf("unexpected path: %s, expected: %s", r.URL.Path, expectedPath)
		}

		response := map[string]interface{}{
			"timeseries": []map[string]interface{}{
				{
					"timestamp":    1766016000,
					"visits":       50,
					"pages_served": 100,
					"cache_hits":   75,
					"cache_misses": 25,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	envsService := NewEnvironmentsService(client)

	// Use site name instead of UUID
	metrics, err := envsService.GetMetrics(context.Background(), testSiteName, testEnvID, "28d")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(metrics) != 1 {
		t.Errorf("expected 1 metric, got %d", len(metrics))
	}

	// Verify site name was resolved (2 requests: lookup + metrics)
	if requestCount != 2 {
		t.Errorf("expected 2 requests (site lookup + metrics), got %d", requestCount)
	}
}

func TestEnvironmentsService_GetMetrics_DatetimeFormat(t *testing.T) {
	testSiteID := "12345678-1234-1234-1234-123456789abc"

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return response with specific timestamp
		response := map[string]interface{}{
			"timeseries": []map[string]interface{}{
				{
					"timestamp":    1609459200, // 2021-01-01 00:00:00 UTC
					"visits":       10,
					"pages_served": 20,
					"cache_hits":   15,
					"cache_misses": 5,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	envsService := NewEnvironmentsService(client)

	metrics, err := envsService.GetMetrics(context.Background(), testSiteID, "dev", "28d")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check datetime format is ISO 8601 with time
	expectedDatetime := "2021-01-01T00:00:00"
	if metrics[0].Datetime != expectedDatetime {
		t.Errorf("expected datetime '%s', got '%s'", expectedDatetime, metrics[0].Datetime)
	}
}

func TestEnvironmentsService_GetMetrics_CacheHitRatioCalculation(t *testing.T) {
	testSiteID := "12345678-1234-1234-1234-123456789abc"

	tests := []struct {
		name          string
		pagesServed   int64
		cacheHits     int64
		expectedRatio string
	}{
		{
			name:          "zero pages served",
			pagesServed:   0,
			cacheHits:     0,
			expectedRatio: "--",
		},
		{
			name:          "100% cache hit",
			pagesServed:   100,
			cacheHits:     100,
			expectedRatio: "100.00%",
		},
		{
			name:          "80% cache hit",
			pagesServed:   500,
			cacheHits:     400,
			expectedRatio: "80.00%",
		},
		{
			name:          "0% cache hit",
			pagesServed:   100,
			cacheHits:     0,
			expectedRatio: "0.00%",
		},
		{
			name:          "fractional percentage",
			pagesServed:   16489,
			cacheHits:     8388,
			expectedRatio: "50.87%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				response := map[string]interface{}{
					"timeseries": []map[string]interface{}{
						{
							"timestamp":    1766016000,
							"visits":       10,
							"pages_served": tt.pagesServed,
							"cache_hits":   tt.cacheHits,
							"cache_misses": tt.pagesServed - tt.cacheHits,
						},
					},
				}

				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(response)
			}))
			defer server.Close()

			client := NewClient(
				WithBaseURL(server.URL),
				WithToken("test-token"),
				WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
			)

			envsService := NewEnvironmentsService(client)
			metrics, err := envsService.GetMetrics(context.Background(), testSiteID, "dev", "28d")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if metrics[0].CacheHitRatio != tt.expectedRatio {
				t.Errorf("expected cache_hit_ratio '%s', got '%s'", tt.expectedRatio, metrics[0].CacheHitRatio)
			}
		})
	}
}

func TestFormatTimestampISO8601(t *testing.T) {
	tests := []struct {
		name      string
		timestamp int64
		expected  string
	}{
		{
			name:      "epoch",
			timestamp: 0,
			expected:  "1970-01-01T00:00:00",
		},
		{
			name:      "2021-01-01",
			timestamp: 1609459200,
			expected:  "2021-01-01T00:00:00",
		},
		{
			name:      "2025-12-18",
			timestamp: 1766016000,
			expected:  "2025-12-18T00:00:00",
		},
		{
			name:      "with time component",
			timestamp: 1609502400, // 2021-01-01 12:00:00 UTC
			expected:  "2021-01-01T12:00:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatTimestampISO8601(tt.timestamp)
			if result != tt.expected {
				t.Errorf("formatTimestampISO8601(%d) = '%s', want '%s'", tt.timestamp, result, tt.expected)
			}
		})
	}
}

func TestEnvironmentsService_GetMetrics_EmptyTimeseries(t *testing.T) {
	testSiteID := "12345678-1234-1234-1234-123456789abc"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"timeseries": []map[string]interface{}{},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	envsService := NewEnvironmentsService(client)

	metrics, err := envsService.GetMetrics(context.Background(), testSiteID, "dev", "28d")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(metrics) != 0 {
		t.Errorf("expected 0 metrics for empty timeseries, got %d", len(metrics))
	}
}
