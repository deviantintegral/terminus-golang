package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pantheon-systems/terminus-go/pkg/api/models"
)

func TestSitesService_ListBranches(t *testing.T) {
	testSiteID := "12345678-1234-1234-1234-123456789abc"

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/sites/" + testSiteID + "/code-tips"
		if r.URL.Path != expectedPath {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		branches := []map[string]interface{}{
			{"id": "master", "sha": "abc123"},
			{"id": "develop", "sha": "def456"},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(branches)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	sitesService := NewSitesService(client)

	branches, err := sitesService.ListBranches(context.Background(), testSiteID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(branches) != 2 {
		t.Errorf("expected 2 branches, got %d", len(branches))
	}

	if branches[0].ID != "master" {
		t.Errorf("expected branch ID 'master', got '%s'", branches[0].ID)
	}

	if branches[0].SHA != "abc123" {
		t.Errorf("expected SHA 'abc123', got '%s'", branches[0].SHA)
	}
}

func TestSitesService_GetPlans(t *testing.T) {
	testSiteID := "12345678-1234-1234-1234-123456789abc"

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/sites/" + testSiteID + "/plans"
		if r.URL.Path != expectedPath {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		plans := []map[string]interface{}{
			{"id": "plan1", "name": "Basic", "sku": "basic", "billing_cycle": "monthly", "price": 35.0},
			{"id": "plan2", "name": "Performance", "sku": "performance", "billing_cycle": "monthly", "price": 175.0},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(plans)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	sitesService := NewSitesService(client)

	plans, err := sitesService.GetPlans(context.Background(), testSiteID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(plans) != 2 {
		t.Errorf("expected 2 plans, got %d", len(plans))
	}

	if plans[0].ID != "plan1" {
		t.Errorf("expected plan ID 'plan1', got '%s'", plans[0].ID)
	}

	if plans[0].Name != "Basic" {
		t.Errorf("expected name 'Basic', got '%s'", plans[0].Name)
	}

	if plans[0].Price != 35.0 {
		t.Errorf("expected price 35.0, got %f", plans[0].Price)
	}
}

func TestSitesService_ListOrganizations(t *testing.T) {
	testSiteID := "12345678-1234-1234-1234-123456789abc"

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/sites/" + testSiteID + "/memberships/organizations"
		if r.URL.Path != expectedPath {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		memberships := []map[string]interface{}{
			{
				"organization": map[string]interface{}{
					"id":   "org1",
					"name": "Organization 1",
				},
			},
			{
				"organization": map[string]interface{}{
					"id":   "org2",
					"name": "Organization 2",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(memberships)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	sitesService := NewSitesService(client)

	orgs, err := sitesService.ListOrganizations(context.Background(), testSiteID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(orgs) != 2 {
		t.Errorf("expected 2 organizations, got %d", len(orgs))
	}

	if orgs[0].OrgID != "org1" {
		t.Errorf("expected org ID 'org1', got '%s'", orgs[0].OrgID)
	}

	if orgs[0].OrgName != "Organization 1" {
		t.Errorf("expected org name 'Organization 1', got '%s'", orgs[0].OrgName)
	}
}

func TestSitesService_List_ReturnsSites(t *testing.T) {
	testUserID := "user-123"

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/users/" + testUserID + "/memberships/sites"
		if r.URL.Path != expectedPath {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		// Simulate API response that includes upstream field
		sites := []map[string]interface{}{
			{
				"site": map[string]interface{}{
					"id":       "site1",
					"name":     "test-site-1",
					"label":    "Test Site 1",
					"upstream": "wordpress",
				},
			},
			{
				"site": map[string]interface{}{
					"id":    "site2",
					"name":  "test-site-2",
					"label": "Test Site 2",
					"upstream": map[string]interface{}{
						"id":    "upstream-id",
						"label": "WordPress",
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(sites)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	sitesService := NewSitesService(client)

	sites, err := sitesService.List(context.Background(), testUserID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(sites) != 2 {
		t.Errorf("expected 2 sites, got %d", len(sites))
	}

	// Verify sites were unmarshaled correctly
	if sites[0].ID != "site1" {
		t.Errorf("expected site ID 'site1', got '%s'", sites[0].ID)
	}

	// Verify Site includes upstream field in JSON
	jsonData, err := json.Marshal(sites[0])
	if err != nil {
		t.Fatalf("failed to marshal site: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	// Site should include upstream
	if _, exists := result["upstream"]; !exists {
		t.Errorf("Site JSON should include upstream field")
	}
}

func TestSitesService_Get_IncludesUpstream(t *testing.T) {
	testSiteID := "12345678-1234-1234-1234-123456789abc"

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/sites/" + testSiteID
		if r.URL.Path != expectedPath {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		// Simulate API response that includes upstream field
		site := map[string]interface{}{
			"id":       testSiteID,
			"name":     "test-site",
			"label":    "Test Site",
			"upstream": "wordpress",
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(site)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	sitesService := NewSitesService(client)

	site, err := sitesService.Get(context.Background(), testSiteID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify site was unmarshaled correctly
	if site.ID != testSiteID {
		t.Errorf("expected site ID '%s', got '%s'", testSiteID, site.ID)
	}

	// Marshal to JSON to verify upstream is included in Site
	jsonData, err := json.Marshal(site)
	if err != nil {
		t.Fatalf("failed to marshal site: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("failed to unmarshal site JSON: %v", err)
	}

	// Verify upstream field IS in Site output
	if _, exists := result["upstream"]; !exists {
		t.Errorf("upstream field should be present in Site JSON output")
	}

	// Verify other fields are present
	if result["id"] != testSiteID {
		t.Errorf("expected id '%s', got '%v'", testSiteID, result["id"])
	}
}

func TestSiteListItem_ExcludesUpstream(t *testing.T) {
	// This test verifies that SiteListItem excludes upstream
	site := &models.Site{
		ID:       "test-id",
		Name:     "test-name",
		Label:    "Test Label",
		Created:  1234567890,
		Upstream: "wordpress",
	}

	// Convert to list item
	listItem := site.ToListItem()

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

	// Ensure upstream is NOT present in list item
	if _, exists := result["upstream"]; exists {
		t.Errorf("upstream field should not be present in SiteListItem JSON output")
	}
}
