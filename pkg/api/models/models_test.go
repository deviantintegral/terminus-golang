// Package models defines data structures for Pantheon API resources.
package models

import (
	"encoding/json"
	"testing"
)

func TestSite_MarshalJSON_IncludesUpstream(t *testing.T) {
	site := &Site{
		ID:       "test-site-id",
		Name:     "test-site",
		Label:    "Test Site",
		Created:  1234567890,
		Upstream: "wordpress",
	}

	jsonData, err := json.Marshal(site)
	if err != nil {
		t.Fatalf("failed to marshal site: %v", err)
	}

	// Unmarshal to a map to check the raw JSON structure
	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	// Verify that the upstream field IS present in full Site JSON
	if _, exists := result["upstream"]; !exists {
		t.Errorf("upstream field should be present in Site JSON output")
	}

	// Verify other fields are present
	if result["id"] != "test-site-id" {
		t.Errorf("expected id 'test-site-id', got '%v'", result["id"])
	}
	if result["name"] != "test-site" {
		t.Errorf("expected name 'test-site', got '%v'", result["name"])
	}
}

func TestSiteListItem_MarshalJSON_ExcludesUpstream(t *testing.T) {
	site := &Site{
		ID:       "test-site-id",
		Name:     "test-site",
		Label:    "Test Site",
		Created:  1234567890,
		Upstream: "wordpress",
	}

	// Convert to list item
	listItem := site.ToListItem()

	jsonData, err := json.Marshal(listItem)
	if err != nil {
		t.Fatalf("failed to marshal site list item: %v", err)
	}

	// Unmarshal to a map to check the raw JSON structure
	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	// Verify that the upstream field is NOT present in list output
	if _, exists := result["upstream"]; exists {
		t.Errorf("upstream field should not be present in SiteListItem JSON output")
	}

	// Verify other fields are present
	if result["id"] != "test-site-id" {
		t.Errorf("expected id 'test-site-id', got '%v'", result["id"])
	}
	if result["name"] != "test-site" {
		t.Errorf("expected name 'test-site', got '%v'", result["name"])
	}
}

func TestSite_ToListItem_CopiesAllFields(t *testing.T) {
	site := &Site{
		ID:                 "test-id",
		Name:               "test-name",
		Label:              "Test Label",
		Created:            1234567890,
		Framework:          "drupal",
		Organization:       "org-123",
		Service:            "free",
		PlanName:           "Sandbox",
		Upstream:           "wordpress",
		PHP:                "8.1",
		Holder:             "user",
		HolderID:           "holder-123",
		Owner:              "owner-123",
		Frozen:             false,
		IsFrozen:           false,
		PreferredZone:      "us-east",
		PreferredZoneLabel: "US East",
		Info:               map[string]interface{}{"foo": "bar"},
	}

	listItem := site.ToListItem()

	// Verify all fields are copied except Upstream
	if listItem.ID != site.ID {
		t.Errorf("expected ID '%s', got '%s'", site.ID, listItem.ID)
	}
	if listItem.Name != site.Name {
		t.Errorf("expected Name '%s', got '%s'", site.Name, listItem.Name)
	}
	if listItem.Label != site.Label {
		t.Errorf("expected Label '%s', got '%s'", site.Label, listItem.Label)
	}
	if listItem.Created != site.Created {
		t.Errorf("expected Created '%d', got '%d'", site.Created, listItem.Created)
	}
	if listItem.Framework != site.Framework {
		t.Errorf("expected Framework '%s', got '%s'", site.Framework, listItem.Framework)
	}
	if listItem.Organization != site.Organization {
		t.Errorf("expected Organization '%s', got '%s'", site.Organization, listItem.Organization)
	}
	if listItem.Service != site.Service {
		t.Errorf("expected Service '%s', got '%s'", site.Service, listItem.Service)
	}
	if listItem.PlanName != site.PlanName {
		t.Errorf("expected PlanName '%s', got '%s'", site.PlanName, listItem.PlanName)
	}
	if listItem.PHP != site.PHP {
		t.Errorf("expected PHP '%s', got '%s'", site.PHP, listItem.PHP)
	}
	if listItem.PreferredZoneLabel != site.PreferredZoneLabel {
		t.Errorf("expected PreferredZoneLabel '%s', got '%s'", site.PreferredZoneLabel, listItem.PreferredZoneLabel)
	}
}

func TestSiteListItem_JSONOutput_Structure(t *testing.T) {
	// This test verifies the complete JSON structure of a SiteListItem
	// to ensure upstream is excluded and other fields are included
	listItem := &SiteListItem{
		ID:      "test-id",
		Name:    "test-name",
		Label:   "Test Label",
		Created: 1234567890,
	}

	jsonData, err := json.Marshal(listItem)
	if err != nil {
		t.Fatalf("failed to marshal site list item: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	expectedFields := []string{"id", "name", "label", "created"}
	for _, field := range expectedFields {
		if _, exists := result[field]; !exists {
			t.Errorf("expected field '%s' to be present in JSON", field)
		}
	}

	// Ensure upstream is NOT present
	if _, exists := result["upstream"]; exists {
		t.Errorf("upstream field should not be present in SiteListItem JSON output")
	}
}

func TestMetrics_Serialize(t *testing.T) {
	tests := []struct {
		name           string
		metrics        *Metrics
		expectedPeriod string
		expectedRatio  string
	}{
		{
			name: "full datetime format",
			metrics: &Metrics{
				Datetime:      "2025-12-18T00:00:00",
				Visits:        100,
				PagesServed:   500,
				CacheHits:     400,
				CacheMisses:   100,
				CacheHitRatio: "80%",
			},
			expectedPeriod: "2025-12-18",
			expectedRatio:  "80%",
		},
		{
			name: "zero values with dash ratio",
			metrics: &Metrics{
				Datetime:      "2025-01-01T00:00:00",
				Visits:        0,
				PagesServed:   0,
				CacheHits:     0,
				CacheMisses:   0,
				CacheHitRatio: "--",
			},
			expectedPeriod: "2025-01-01",
			expectedRatio:  "--",
		},
		{
			name: "date only format",
			metrics: &Metrics{
				Datetime:      "2025-06-15",
				Visits:        50,
				PagesServed:   200,
				CacheHits:     180,
				CacheMisses:   20,
				CacheHitRatio: "90%",
			},
			expectedPeriod: "2025-06-15",
			expectedRatio:  "90%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := tt.metrics.Serialize()

			// Check we have the right number of fields
			if len(fields) != 6 {
				t.Errorf("expected 6 fields, got %d", len(fields))
			}

			// Check field names and values
			expectedNames := []string{"Period", "Visits", "Pages Served", "Cache Hits", "Cache Misses", "Cache Hit Ratio"}
			for i, expectedName := range expectedNames {
				if fields[i].Name != expectedName {
					t.Errorf("expected field %d name to be '%s', got '%s'", i, expectedName, fields[i].Name)
				}
			}

			// Check Period is date-only
			if fields[0].Value != tt.expectedPeriod {
				t.Errorf("expected Period '%s', got '%s'", tt.expectedPeriod, fields[0].Value)
			}

			// Check Cache Hit Ratio
			if fields[5].Value != tt.expectedRatio {
				t.Errorf("expected Cache Hit Ratio '%s', got '%s'", tt.expectedRatio, fields[5].Value)
			}
		})
	}
}

func TestMetrics_DefaultFields(t *testing.T) {
	metrics := &Metrics{}
	fields := metrics.DefaultFields()

	expectedFields := []string{"Period", "Visits", "Pages Served", "Cache Hits", "Cache Misses", "Cache Hit Ratio"}

	if len(fields) != len(expectedFields) {
		t.Errorf("expected %d default fields, got %d", len(expectedFields), len(fields))
	}

	for i, expected := range expectedFields {
		if fields[i] != expected {
			t.Errorf("expected default field %d to be '%s', got '%s'", i, expected, fields[i])
		}
	}
}

func TestMetrics_MarshalJSON(t *testing.T) {
	metrics := &Metrics{
		Datetime:      "2025-12-18T00:00:00",
		Visits:        100,
		PagesServed:   500,
		CacheHits:     400,
		CacheMisses:   100,
		CacheHitRatio: "80%",
	}

	jsonData, err := json.Marshal(metrics)
	if err != nil {
		t.Fatalf("failed to marshal metrics: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	// Verify JSON field names match PHP terminus (snake_case)
	expectedFields := []string{"datetime", "visits", "pages_served", "cache_hits", "cache_misses", "cache_hit_ratio"}
	for _, field := range expectedFields {
		if _, exists := result[field]; !exists {
			t.Errorf("expected field '%s' to be present in JSON", field)
		}
	}

	// Verify values
	if result["datetime"] != "2025-12-18T00:00:00" {
		t.Errorf("expected datetime '2025-12-18T00:00:00', got '%v'", result["datetime"])
	}
	if result["cache_hit_ratio"] != "80%" {
		t.Errorf("expected cache_hit_ratio '80%%', got '%v'", result["cache_hit_ratio"])
	}
}

func TestMetrics_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"datetime": "2025-12-18T00:00:00",
		"visits": 100,
		"pages_served": 500,
		"cache_hits": 400,
		"cache_misses": 100,
		"cache_hit_ratio": "--"
	}`

	var metrics Metrics
	if err := json.Unmarshal([]byte(jsonData), &metrics); err != nil {
		t.Fatalf("failed to unmarshal metrics: %v", err)
	}

	if metrics.Datetime != "2025-12-18T00:00:00" {
		t.Errorf("expected datetime '2025-12-18T00:00:00', got '%s'", metrics.Datetime)
	}
	if metrics.Visits != 100 {
		t.Errorf("expected visits 100, got %d", metrics.Visits)
	}
	if metrics.PagesServed != 500 {
		t.Errorf("expected pages_served 500, got %d", metrics.PagesServed)
	}
	if metrics.CacheHits != 400 {
		t.Errorf("expected cache_hits 400, got %d", metrics.CacheHits)
	}
	if metrics.CacheMisses != 100 {
		t.Errorf("expected cache_misses 100, got %d", metrics.CacheMisses)
	}
	if metrics.CacheHitRatio != "--" {
		t.Errorf("expected cache_hit_ratio '--', got '%s'", metrics.CacheHitRatio)
	}
}
