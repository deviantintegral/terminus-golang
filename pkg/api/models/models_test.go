// Package models defines data structures for Pantheon API resources.
package models

import (
	"encoding/json"
	"testing"
)

func TestSite_MarshalJSON_ExcludesUpstream(t *testing.T) {
	site := &Site{
		ID:       "test-site-id",
		Name:     "test-site",
		Label:    "Test Site",
		Created:  1234567890,
		Upstream: "should-not-appear-in-json",
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

	// Verify that the upstream field is NOT present in the JSON
	if _, exists := result["upstream"]; exists {
		t.Errorf("upstream field should not be present in JSON output, but it was found")
	}

	// Verify other fields are present
	if result["id"] != "test-site-id" {
		t.Errorf("expected id 'test-site-id', got '%v'", result["id"])
	}
	if result["name"] != "test-site" {
		t.Errorf("expected name 'test-site', got '%v'", result["name"])
	}
}

func TestSite_UnmarshalJSON_IgnoresUpstream(t *testing.T) {
	// Test with upstream as a string - should be ignored during unmarshal
	jsonData := []byte(`{
		"id": "test-site-id",
		"name": "test-site",
		"label": "Test Site",
		"upstream": "wordpress"
	}`)

	var site Site
	if err := json.Unmarshal(jsonData, &site); err != nil {
		t.Fatalf("failed to unmarshal site: %v", err)
	}

	// Verify the upstream field is NOT populated (due to json:"-" tag)
	if site.Upstream != nil {
		t.Errorf("upstream field should be ignored during unmarshal, but got '%v'", site.Upstream)
	}

	// Verify other fields are populated correctly
	if site.ID != "test-site-id" {
		t.Errorf("expected id 'test-site-id', got '%s'", site.ID)
	}
	if site.Name != "test-site" {
		t.Errorf("expected name 'test-site', got '%s'", site.Name)
	}
}

func TestSite_UnmarshalJSON_IgnoresUpstreamObject(t *testing.T) {
	// Test with upstream as an object - should be ignored during unmarshal
	jsonData := []byte(`{
		"id": "test-site-id",
		"name": "test-site",
		"label": "Test Site",
		"upstream": {
			"id": "upstream-id",
			"label": "WordPress"
		}
	}`)

	var site Site
	if err := json.Unmarshal(jsonData, &site); err != nil {
		t.Fatalf("failed to unmarshal site: %v", err)
	}

	// Verify the upstream field is NOT populated (due to json:"-" tag)
	if site.Upstream != nil {
		t.Errorf("upstream field should be ignored during unmarshal, but got '%v'", site.Upstream)
	}

	// Verify other fields are populated correctly
	if site.ID != "test-site-id" {
		t.Errorf("expected id 'test-site-id', got '%s'", site.ID)
	}
}

func TestSite_RoundTrip_DropsUpstream(t *testing.T) {
	// This test verifies that when we unmarshal a site with upstream,
	// then marshal it again, the upstream field is dropped
	originalJSON := []byte(`{
		"id": "test-site-id",
		"name": "test-site",
		"label": "Test Site",
		"upstream": "wordpress"
	}`)

	var site Site
	if err := json.Unmarshal(originalJSON, &site); err != nil {
		t.Fatalf("failed to unmarshal site: %v", err)
	}

	// Marshal it back
	newJSON, err := json.Marshal(site)
	if err != nil {
		t.Fatalf("failed to marshal site: %v", err)
	}

	// Unmarshal to a map to verify structure
	var result map[string]interface{}
	if err := json.Unmarshal(newJSON, &result); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	// Verify upstream is not in the output
	if _, exists := result["upstream"]; exists {
		t.Errorf("upstream field should not be present in re-marshaled JSON")
	}

	// Verify other fields are preserved
	if result["id"] != "test-site-id" {
		t.Errorf("expected id 'test-site-id', got '%v'", result["id"])
	}
}
