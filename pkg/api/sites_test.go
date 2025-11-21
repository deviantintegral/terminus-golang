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

func TestIsUUID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid UUID", "12345678-1234-1234-1234-123456789abc", true},
		{"valid UUID uppercase", "12345678-1234-1234-1234-123456789ABC", true},
		{"site name", "my-site-name", false},
		{"short string", "short", false},
		{"empty string", "", false},
		{"UUID without dashes", "12345678123412341234123456789abc", false},
		{"wrong dash positions", "1234567-81234-1234-1234-123456789abc", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsUUID(tt.input)
			if result != tt.expected {
				t.Errorf("IsUUID(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestResolveSiteNameToID(t *testing.T) {
	siteName := "my-site"
	siteID := "12345678-1234-1234-1234-123456789abc"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/site-names/" + siteName
		if r.URL.Path != expectedPath {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"id": siteID})
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	resolvedID, err := ResolveSiteNameToID(context.Background(), client, siteName)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resolvedID != siteID {
		t.Errorf("expected ID %s, got %s", siteID, resolvedID)
	}
}

func TestResolveSiteNameToID_EmptyID(t *testing.T) {
	siteName := "my-site"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"id": ""})
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	_, err := ResolveSiteNameToID(context.Background(), client, siteName)
	if err == nil {
		t.Fatal("expected error for empty ID")
	}
}

func TestEnsureSiteUUID(t *testing.T) {
	t.Run("already UUID", func(t *testing.T) {
		uuid := "12345678-1234-1234-1234-123456789abc"
		client := NewClient()

		result, err := EnsureSiteUUID(context.Background(), client, uuid)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result != uuid {
			t.Errorf("expected %s, got %s", uuid, result)
		}
	})

	t.Run("site name resolution", func(t *testing.T) {
		siteName := "my-site"
		siteID := "12345678-1234-1234-1234-123456789abc"

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"id": siteID})
		}))
		defer server.Close()

		client := NewClient(
			WithBaseURL(server.URL),
			WithToken("test-token"),
			WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
		)

		result, err := EnsureSiteUUID(context.Background(), client, siteName)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result != siteID {
			t.Errorf("expected %s, got %s", siteID, result)
		}
	})
}

func TestNewSitesService(t *testing.T) {
	client := NewClient()
	service := NewSitesService(client)

	if service == nil {
		t.Fatal("expected service to be created")
	}

	if service.client != client {
		t.Error("expected service to have correct client")
	}
}

func TestSitesService_Get_WithSiteName(t *testing.T) {
	siteName := "my-site"
	siteID := "12345678-1234-1234-1234-123456789abc"

	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		switch r.URL.Path {
		case "/site-names/" + siteName:
			// Name resolution
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"id": siteID})
		case "/sites/" + siteID:
			// Get site
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"id":    siteID,
				"name":  siteName,
				"label": "My Site",
			})
		default:
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	service := NewSitesService(client)
	site, err := service.Get(context.Background(), siteName)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if site.ID != siteID {
		t.Errorf("expected site ID %s, got %s", siteID, site.ID)
	}

	if callCount != 2 {
		t.Errorf("expected 2 API calls (resolve + get), got %d", callCount)
	}
}

func TestSitesService_Update(t *testing.T) {
	siteID := "12345678-1234-1234-1234-123456789abc"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT request, got %s", r.Method)
		}

		expectedPath := "/sites/" + siteID
		if r.URL.Path != expectedPath {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req UpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.Label != "Updated Label" {
			t.Errorf("expected label 'Updated Label', got '%s'", req.Label)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"id":    siteID,
			"label": req.Label,
		})
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	service := NewSitesService(client)
	site, err := service.Update(context.Background(), siteID, &UpdateRequest{
		Label: "Updated Label",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if site.Label != "Updated Label" {
		t.Errorf("expected label 'Updated Label', got '%s'", site.Label)
	}
}

func TestSitesService_ListByOrganization(t *testing.T) {
	orgID := "org-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/organizations/" + orgID + "/memberships/sites"
		if r.URL.Path != expectedPath {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		memberships := []map[string]interface{}{
			{
				"site": map[string]interface{}{
					"id":    "site1",
					"name":  "site-1",
					"label": "Site 1",
				},
			},
			{
				"site": map[string]interface{}{
					"id":    "site2",
					"name":  "site-2",
					"label": "Site 2",
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

	service := NewSitesService(client)
	sites, err := service.ListByOrganization(context.Background(), orgID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(sites) != 2 {
		t.Errorf("expected 2 sites, got %d", len(sites))
	}

	if sites[0].ID != "site1" {
		t.Errorf("expected site ID 'site1', got '%s'", sites[0].ID)
	}
}

func TestSitesService_GetTeam(t *testing.T) {
	siteID := "12345678-1234-1234-1234-123456789abc"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/sites/" + siteID + "/memberships/users"
		if r.URL.Path != expectedPath {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		team := []map[string]interface{}{
			{
				"id":        "user1",
				"email":     "user1@example.com",
				"firstname": "John",
				"lastname":  "Doe",
				"role":      "owner",
			},
			{
				"id":        "user2",
				"email":     "user2@example.com",
				"firstname": "Jane",
				"lastname":  "Smith",
				"role":      "team_member",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(team)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	service := NewSitesService(client)
	team, err := service.GetTeam(context.Background(), siteID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(team) != 2 {
		t.Errorf("expected 2 team members, got %d", len(team))
	}

	if team[0].Email != "user1@example.com" {
		t.Errorf("expected email 'user1@example.com', got '%s'", team[0].Email)
	}
}

func TestSitesService_AddTeamMember(t *testing.T) {
	siteID := "12345678-1234-1234-1234-123456789abc"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		expectedPath := "/sites/" + siteID + "/memberships/users"
		if r.URL.Path != expectedPath {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req AddTeamMemberRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"id":        "user123",
			"email":     req.Email,
			"firstname": "New",
			"lastname":  "User",
			"role":      req.Role,
		})
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	service := NewSitesService(client)
	member, err := service.AddTeamMember(context.Background(), siteID, &AddTeamMemberRequest{
		Email: "newuser@example.com",
		Role:  "team_member",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if member.Email != "newuser@example.com" {
		t.Errorf("expected email 'newuser@example.com', got '%s'", member.Email)
	}

	if member.Role != "team_member" {
		t.Errorf("expected role 'team_member', got '%s'", member.Role)
	}
}

func TestSitesService_RemoveTeamMember(t *testing.T) {
	siteID := "12345678-1234-1234-1234-123456789abc"
	userID := "user123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE request, got %s", r.Method)
		}

		expectedPath := "/sites/" + siteID + "/memberships/users/" + userID
		if r.URL.Path != expectedPath {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	service := NewSitesService(client)
	err := service.RemoveTeamMember(context.Background(), siteID, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSitesService_GetTags(t *testing.T) {
	siteID := "12345678-1234-1234-1234-123456789abc"
	orgID := "org-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/organizations/" + orgID + "/tags/sites/" + siteID
		if r.URL.Path != expectedPath {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		tags := []map[string]interface{}{
			{"id": "tag1", "name": "production", "site_id": siteID, "org_id": orgID},
			{"id": "tag2", "name": "wordpress", "site_id": siteID, "org_id": orgID},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(tags)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	service := NewSitesService(client)
	tags, err := service.GetTags(context.Background(), siteID, orgID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(tags))
	}

	if tags[0].Name != "production" {
		t.Errorf("expected tag name 'production', got '%s'", tags[0].Name)
	}
}

func TestSitesService_AddTag(t *testing.T) {
	siteID := "12345678-1234-1234-1234-123456789abc"
	orgID := "org-123"
	tagName := "production"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		expectedPath := "/organizations/" + orgID + "/tags/" + tagName + "/sites"
		if r.URL.Path != expectedPath {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req struct {
			SiteID string `json:"site_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.SiteID != siteID {
			t.Errorf("expected site_id '%s', got '%s'", siteID, req.SiteID)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	service := NewSitesService(client)
	err := service.AddTag(context.Background(), siteID, orgID, tagName)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSitesService_RemoveTag(t *testing.T) {
	siteID := "12345678-1234-1234-1234-123456789abc"
	orgID := "org-123"
	tagName := "production"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE request, got %s", r.Method)
		}

		expectedPath := "/organizations/" + orgID + "/tags/" + tagName + "/sites/" + siteID
		if r.URL.Path != expectedPath {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	service := NewSitesService(client)
	err := service.RemoveTag(context.Background(), siteID, orgID, tagName)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSitesService_GetPlan(t *testing.T) {
	siteID := "12345678-1234-1234-1234-123456789abc"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/sites/" + siteID + "/plan"
		if r.URL.Path != expectedPath {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		plan := map[string]interface{}{
			"id":                    "plan1",
			"name":                  "Basic",
			"sku":                   "basic",
			"billing_cycle":         "monthly",
			"price":                 35.0,
			"multidev_environments": 0,
			"automated_backups":     false,
			"cache_server":          false,
			"custom_upstreams":      false,
			"new_relic":             false,
			"secure_runtime_access": false,
			"storage_gb":            5,
			"support_plan":          "basic",
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(plan)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	service := NewSitesService(client)
	plan, err := service.GetPlan(context.Background(), siteID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if plan.ID != "plan1" {
		t.Errorf("expected plan ID 'plan1', got '%s'", plan.ID)
	}

	if plan.Name != "Basic" {
		t.Errorf("expected name 'Basic', got '%s'", plan.Name)
	}

	if plan.Price != 35.0 {
		t.Errorf("expected price 35.0, got %f", plan.Price)
	}
}
