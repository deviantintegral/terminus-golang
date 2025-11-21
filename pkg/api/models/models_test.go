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
