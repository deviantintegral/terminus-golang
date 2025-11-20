package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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
