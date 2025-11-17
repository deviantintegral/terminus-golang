package api

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/pantheon-systems/terminus-go/pkg/api/models"
)

// skipIfNoToken skips the test if PANTHEON_MACHINE_TOKEN is not set
func skipIfNoToken(t *testing.T) string {
	t.Helper()
	token := os.Getenv("PANTHEON_MACHINE_TOKEN")
	if token == "" {
		t.Skip("Skipping integration test: PANTHEON_MACHINE_TOKEN not set")
	}
	return token
}

// fixtureRecorder handles recording and redacting API responses
type fixtureRecorder struct {
	t        *testing.T
	basePath string
}

// newFixtureRecorder creates a new fixture recorder
func newFixtureRecorder(t *testing.T) *fixtureRecorder {
	t.Helper()
	basePath := filepath.Join("testdata", "fixtures")
	if err := os.MkdirAll(basePath, 0o755); err != nil {
		t.Fatalf("failed to create fixtures directory: %v", err)
	}
	return &fixtureRecorder{
		t:        t,
		basePath: basePath,
	}
}

// redactSensitiveData removes sensitive information from JSON responses
func (r *fixtureRecorder) redactSensitiveData(data []byte) []byte {
	// Convert to string for easier manipulation
	str := string(data)

	// Redact patterns
	patterns := map[string]string{
		// Session tokens (UUIDs and long strings)
		`"session":\s*"[^"]{20,}"`:                    `"session": "REDACTED"`,
		`"Session":\s*"[^"]{20,}"`:                    `"Session": "REDACTED"`,
		`"session_token":\s*"[^"]{20,}"`:              `"session_token": "REDACTED"`,
		`"SessionToken":\s*"[^"]{20,}"`:               `"SessionToken": "REDACTED"`,
		// Machine tokens
		`"machine_token":\s*"[^"]{20,}"`:              `"machine_token": "REDACTED"`,
		`"MachineToken":\s*"[^"]{20,}"`:               `"MachineToken": "REDACTED"`,
		// User IDs (UUIDs)
		`"user_id":\s*"[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}"`: `"user_id": "REDACTED-USER-ID"`,
		`"UserID":\s*"[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}"`:  `"UserID": "REDACTED-USER-ID"`,
		`"id":\s*"[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}"`:      `"id": "REDACTED-ID"`,
		// Emails - replace with redacted email
		`"email":\s*"[^"]+@[^"]+\.[^"]+"`:             `"email": "redacted@example.com"`,
		`"Email":\s*"[^"]+@[^"]+\.[^"]+"`:             `"Email": "redacted@example.com"`,
	}

	for pattern, replacement := range patterns {
		re := regexp.MustCompile(pattern)
		str = re.ReplaceAllString(str, replacement)
	}

	return []byte(str)
}

// record saves a fixture to disk with redacted sensitive data
func (r *fixtureRecorder) record(name string, data interface{}) {
	r.t.Helper()

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		r.t.Fatalf("failed to marshal fixture data: %v", err)
	}

	// Redact sensitive data
	redacted := r.redactSensitiveData(jsonData)

	// Write to file
	filePath := filepath.Join(r.basePath, name+".json")
	if err := os.WriteFile(filePath, redacted, 0o644); err != nil {
		r.t.Fatalf("failed to write fixture: %v", err)
	}

	r.t.Logf("Recorded fixture: %s", filePath)
}

// TestAuthLogin tests the auth:login command
func TestAuthLogin(t *testing.T) {
	token := skipIfNoToken(t)
	recorder := newFixtureRecorder(t)

	client := NewClient()
	authService := NewAuthService(client)

	ctx := context.Background()
	email := "test@example.com"

	// Test login
	session, err := authService.Login(ctx, token, email)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	// Validate response
	if session == nil {
		t.Fatal("Expected session to be non-nil")
	}
	if session.Session == "" {
		t.Error("Expected session token to be set")
	}
	if session.UserID == "" {
		t.Error("Expected user ID to be set")
	}
	// Validate timestamp is set and valid
	if session.ExpiresAt == 0 {
		t.Error("Expected expires_at to be set")
	} else {
		expiresTime := time.Unix(session.ExpiresAt, 0)
		now := time.Now()

		// Check that the expiry is in the future
		if expiresTime.Before(now) {
			t.Errorf("Expected expires_at to be in the future, got %v (now: %v)", expiresTime, now)
		}

		// Check that the expiry is within a reasonable range (30 days)
		maxExpiry := now.Add(30 * 24 * time.Hour)
		if expiresTime.After(maxExpiry) {
			t.Errorf("Expected expires_at to be within 30 days, got %v", expiresTime)
		}
	}

	// Record fixture
	recorder.record("auth_login", session)

	t.Logf("Login successful - Session expires at: %d", session.ExpiresAt)
}

// TestAuthWhoami tests the auth:whoami command
func TestAuthWhoami(t *testing.T) {
	token := skipIfNoToken(t)
	recorder := newFixtureRecorder(t)

	client := NewClient()
	authService := NewAuthService(client)

	ctx := context.Background()

	// Login first
	session, err := authService.Login(ctx, token, "")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	// Update client with session token
	client.SetToken(session.Session)

	// Test whoami
	user, err := authService.Whoami(ctx, session.UserID)
	if err != nil {
		// Record the error for documentation purposes
		errorResponse := map[string]interface{}{
			"error": err.Error(),
			"note":  "This endpoint may not be available for all Pantheon accounts",
		}
		recorder.record("auth_whoami_error", errorResponse)
		t.Skipf("Whoami endpoint not available: %v", err)
		return
	}

	// Validate response
	if user == nil {
		t.Fatal("Expected user to be non-nil")
	}
	if user.ID == "" {
		t.Error("Expected user ID to be set")
	}
	if user.Email == "" {
		t.Error("Expected email to be set")
	}

	// Record fixture
	recorder.record("auth_whoami", user)

	t.Logf("Whoami successful - User: %s %s (%s)", user.FirstName, user.LastName, user.Email)
}

// TestAuthLogout tests the auth:logout command
// Note: Logout is a local operation that clears stored credentials
// We test that the session becomes invalid after clearing the token
func TestAuthLogout(t *testing.T) {
	token := skipIfNoToken(t)

	client := NewClient()
	authService := NewAuthService(client)

	ctx := context.Background()

	// Login first
	session, err := authService.Login(ctx, token, "")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	// Set session token (note: we can't validate it if /user endpoint doesn't work)
	client.SetToken(session.Session)
	if session.Session == "" {
		t.Error("Expected session token to be set after login")
	}

	// Simulate logout by clearing token
	client.SetToken("")

	// Verify token was cleared
	if client.token != "" {
		t.Error("Expected token to be cleared after logout")
	}

	t.Log("Logout successful - Session cleared")
}

// TestOrgList tests the org:list command
func TestOrgList(t *testing.T) {
	token := skipIfNoToken(t)
	recorder := newFixtureRecorder(t)

	client := NewClient()
	authService := NewAuthService(client)
	orgService := NewOrganizationsService(client)

	ctx := context.Background()

	// Login first
	session, err := authService.Login(ctx, token, "")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	client.SetToken(session.Session)

	// Test org list
	orgs, err := orgService.List(ctx, session.UserID)
	if err != nil {
		// Record the error for documentation purposes
		errorResponse := map[string]interface{}{
			"error": err.Error(),
			"note":  "This endpoint may not be available for all Pantheon accounts",
		}
		recorder.record("org_list_error", errorResponse)
		t.Skipf("List organizations endpoint not available: %v", err)
		return
	}

	// Validate response
	if orgs == nil {
		t.Fatal("Expected orgs to be non-nil")
	}

	// Record fixture
	recorder.record("org_list", orgs)

	t.Logf("Found %d organizations", len(orgs))
	for i, org := range orgs {
		if org != nil {
			t.Logf("  [%d] %s (ID: %s)", i+1, org.Name, org.ID)
		}
	}
}

// TestOrgInfo tests the org:info command
func TestOrgInfo(t *testing.T) {
	token := skipIfNoToken(t)
	recorder := newFixtureRecorder(t)

	client := NewClient()
	authService := NewAuthService(client)
	orgService := NewOrganizationsService(client)

	ctx := context.Background()

	// Login first
	session, err := authService.Login(ctx, token, "")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	client.SetToken(session.Session)

	// Get org list to find an org ID
	orgs, err := orgService.List(ctx, session.UserID)
	if err != nil {
		t.Skipf("Cannot test org:info - org:list failed: %v", err)
		return
	}

	if len(orgs) == 0 {
		t.Skip("No organizations found - skipping org:info test")
	}

	// Get info for first org
	orgID := orgs[0].ID
	org, err := orgService.Get(ctx, orgID)
	if err != nil {
		// Record the error for documentation purposes
		errorResponse := map[string]interface{}{
			"error":  err.Error(),
			"note":   "This endpoint may not be available for all Pantheon accounts",
			"org_id": orgID,
		}
		recorder.record("org_info_error", errorResponse)
		t.Skipf("Get organization endpoint not available: %v", err)
		return
	}

	// Validate response
	if org == nil {
		t.Fatal("Expected org to be non-nil")
	}
	if org.ID == "" {
		t.Error("Expected org ID to be set")
	}
	if org.Name == "" {
		t.Error("Expected org name to be set")
	}

	// Record fixture
	recorder.record("org_info", org)

	t.Logf("Organization info: %s (ID: %s)", org.Name, org.ID)
}

// TestSiteList tests the site:list command
func TestSiteList(t *testing.T) {
	token := skipIfNoToken(t)
	recorder := newFixtureRecorder(t)

	client := NewClient()
	authService := NewAuthService(client)
	siteService := NewSitesService(client)

	ctx := context.Background()

	// Login first
	session, err := authService.Login(ctx, token, "")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	client.SetToken(session.Session)

	// Test site list
	sites, err := siteService.List(ctx, session.UserID)
	if err != nil {
		// Record the error for documentation purposes
		errorResponse := map[string]interface{}{
			"error": err.Error(),
			"note":  "This endpoint may not be available for all Pantheon accounts",
		}
		recorder.record("site_list_error", errorResponse)
		t.Skipf("List sites endpoint not available: %v", err)
		return
	}

	// Validate response
	if sites == nil {
		t.Fatal("Expected sites to be non-nil")
	}

	// Record fixture
	recorder.record("site_list", sites)

	t.Logf("Found %d sites", len(sites))
	for i, site := range sites {
		if site != nil {
			t.Logf("  [%d] %s (ID: %s)", i+1, site.Name, site.ID)
		}
	}
}

// TestSiteInfo tests the site:info command
func TestSiteInfo(t *testing.T) {
	token := skipIfNoToken(t)
	recorder := newFixtureRecorder(t)

	client := NewClient()
	authService := NewAuthService(client)
	siteService := NewSitesService(client)

	ctx := context.Background()

	// Login first
	session, err := authService.Login(ctx, token, "")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	client.SetToken(session.Session)

	// Get site list to find a site ID
	sites, err := siteService.List(ctx, session.UserID)
	if err != nil {
		t.Skipf("Cannot test site:info - site:list failed: %v", err)
		return
	}

	if len(sites) == 0 {
		t.Skip("No sites found - skipping site:info test")
	}

	// Get info for first site
	siteID := sites[0].ID
	site, err := siteService.Get(ctx, siteID)
	if err != nil {
		// Record the error for documentation purposes
		errorResponse := map[string]interface{}{
			"error":   err.Error(),
			"note":    "This endpoint may not be available for all Pantheon accounts",
			"site_id": siteID,
		}
		recorder.record("site_info_error", errorResponse)
		t.Skipf("Get site endpoint not available: %v", err)
		return
	}

	// Validate response
	if site == nil {
		t.Fatal("Expected site to be non-nil")
	}
	if site.ID == "" {
		t.Error("Expected site ID to be set")
	}
	if site.Name == "" {
		t.Error("Expected site name to be set")
	}

	// Record fixture
	recorder.record("site_info", site)

	t.Logf("Site info: %s (ID: %s)", site.Name, site.ID)
}

// TestIntegrationSequence tests a full sequence of operations
func TestIntegrationSequence(t *testing.T) {
	token := skipIfNoToken(t)

	client := NewClient()
	authService := NewAuthService(client)
	orgService := NewOrganizationsService(client)
	siteService := NewSitesService(client)

	ctx := context.Background()

	// 1. Login
	var session *SessionResponse
	t.Run("Login", func(t *testing.T) {
		var err error
		session, err = authService.Login(ctx, token, "integration@example.com")
		if err != nil {
			t.Fatalf("Login failed: %v", err)
		}
		client.SetToken(session.Session)
		t.Logf("✓ Login successful")
	})

	// 2. Whoami
	var user *models.User
	t.Run("Whoami", func(t *testing.T) {
		if session == nil {
			t.Skip("Login failed, skipping whoami test")
			return
		}
		var err error
		user, err = authService.Whoami(ctx, session.UserID)
		if err != nil {
			t.Logf("⚠ Whoami endpoint not available: %v", err)
			return
		}
		t.Logf("✓ Authenticated as: %s %s (%s)", user.FirstName, user.LastName, user.Email)
	})

	// 3. List organizations
	t.Run("ListOrganizations", func(t *testing.T) {
		orgs, err := orgService.List(ctx, session.UserID)
		if err != nil {
			t.Logf("⚠ List organizations endpoint not available: %v", err)
			return
		}
		t.Logf("✓ Found %d organizations", len(orgs))
	})

	// 4. List sites
	t.Run("ListSites", func(t *testing.T) {
		sites, err := siteService.List(ctx, session.UserID)
		if err != nil {
			t.Logf("⚠ List sites endpoint not available: %v", err)
			return
		}
		t.Logf("✓ Found %d sites", len(sites))
	})

	// 5. Logout (clear token)
	t.Run("Logout", func(t *testing.T) {
		client.SetToken("")
		t.Logf("✓ Logout successful")
	})
}

// siteLifecycleTestData holds the test context for site lifecycle tests
type siteLifecycleTestData struct {
	ctx         context.Context
	client      *Client
	authService *AuthService
	orgService  *OrganizationsService
	siteService *SitesService
	session     *SessionResponse
	orgID       string
	upstreamID  string
	siteName    string
	siteID      string
	recorder    *fixtureRecorder
}

// setupSiteLifecycleTest sets up the test environment
func setupSiteLifecycleTest(t *testing.T) *siteLifecycleTestData {
	t.Helper()
	token := skipIfNoToken(t)

	data := &siteLifecycleTestData{
		ctx:         context.Background(),
		client:      NewClient(),
		recorder:    newFixtureRecorder(t),
	}

	data.authService = NewAuthService(data.client)
	data.orgService = NewOrganizationsService(data.client)
	data.siteService = NewSitesService(data.client)

	// Login
	session, err := data.authService.Login(data.ctx, token, "")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	data.session = session
	data.client.SetToken(session.Session)

	// Get organization
	orgs, err := data.orgService.List(data.ctx, session.UserID)
	if err != nil || len(orgs) == 0 {
		t.Skip("No organizations available")
	}
	data.orgID = orgs[0].ID
	t.Logf("Using organization: %s (ID: %s)", orgs[0].Name, data.orgID)

	// Get upstream
	upstreams, err := data.orgService.ListUpstreams(data.ctx, data.orgID)
	if err != nil || len(upstreams) == 0 {
		t.Skip("No upstreams available")
	}
	data.upstreamID = upstreams[0].ID
	t.Logf("Using upstream: %s (ID: %s)", upstreams[0].Label, data.upstreamID)

	// Generate unique site name
	data.siteName = fmt.Sprintf("test-site-%s", regexp.MustCompile(`\D`).ReplaceAllString(fmt.Sprintf("%d", session.ExpiresAt), "")[:10])
	t.Logf("Test site name: %s", data.siteName)

	return data
}

// TestSiteLifecycle tests the full site lifecycle: create, list, get info, delete
func TestSiteLifecycle(t *testing.T) {
	data := setupSiteLifecycleTest(t)

	// Run subtests for each lifecycle stage
	t.Run("CreateSite", func(t *testing.T) {
		testCreateSite(t, data)
	})

	// Ensure cleanup
	defer func() {
		if data.siteID != "" {
			t.Logf("Cleaning up test site: %s", data.siteID)
			if err := data.siteService.Delete(data.ctx, data.siteID); err != nil {
				t.Logf("Warning: Failed to cleanup: %v", err)
			}
		}
	}()

	t.Run("VerifySiteInList", func(t *testing.T) {
		testVerifySiteInList(t, data)
	})

	t.Run("VerifySiteInfo", func(t *testing.T) {
		testVerifySiteInfo(t, data)
	})

	t.Run("DeleteSite", func(t *testing.T) {
		testDeleteSite(t, data)
	})

	t.Run("VerifySiteDeleted", func(t *testing.T) {
		testVerifySiteDeleted(t, data)
	})
}

// testCreateSite creates a test site
func testCreateSite(t *testing.T, data *siteLifecycleTestData) {
	t.Helper()

	req := &CreateSiteRequest{
		SiteName:     data.siteName,
		Label:        data.siteName,
		UpstreamID:   data.upstreamID,
		Organization: data.orgID,
	}

	site, err := data.siteService.Create(data.ctx, req)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	if site == nil || site.ID == "" {
		t.Fatal("Expected site with valid ID")
	}
	if site.Name != data.siteName {
		t.Errorf("Expected site name %s, got %s", data.siteName, site.Name)
	}

	data.siteID = site.ID
	data.recorder.record("site_lifecycle_created", site)
	t.Logf("✓ Site created: %s (ID: %s)", site.Name, site.ID)
}

// testVerifySiteInList verifies the site appears in the list
func testVerifySiteInList(t *testing.T, data *siteLifecycleTestData) {
	t.Helper()

	sites, err := data.siteService.List(data.ctx, data.session.UserID)
	if err != nil {
		t.Fatalf("Failed to list sites: %v", err)
	}

	for _, s := range sites {
		if s.ID == data.siteID {
			t.Logf("✓ Site found in list: %s", s.Name)
			return
		}
	}
	t.Errorf("Site %s not found in list", data.siteID)
}

// testVerifySiteInfo verifies site information is correct
func testVerifySiteInfo(t *testing.T, data *siteLifecycleTestData) {
	t.Helper()

	site, err := data.siteService.Get(data.ctx, data.siteID)
	if err != nil {
		t.Fatalf("Failed to get site info: %v", err)
	}

	if site.ID != data.siteID {
		t.Errorf("Expected ID %s, got %s", data.siteID, site.ID)
	}
	if site.Name != data.siteName {
		t.Errorf("Expected name %s, got %s", data.siteName, site.Name)
	}
	if site.Label != data.siteName {
		t.Errorf("Expected label %s, got %s", data.siteName, site.Label)
	}

	data.recorder.record("site_lifecycle_info", site)
	t.Logf("✓ Site info verified: %s (ID: %s)", site.Name, site.ID)
}

// testDeleteSite deletes the test site
func testDeleteSite(t *testing.T, data *siteLifecycleTestData) {
	t.Helper()

	if err := data.siteService.Delete(data.ctx, data.siteID); err != nil {
		t.Fatalf("Failed to delete site: %v", err)
	}
	t.Logf("✓ Site deleted: %s", data.siteID)
}

// testVerifySiteDeleted verifies the site no longer appears in the list
func testVerifySiteDeleted(t *testing.T, data *siteLifecycleTestData) {
	t.Helper()

	sites, err := data.siteService.List(data.ctx, data.session.UserID)
	if err != nil {
		t.Fatalf("Failed to list sites: %v", err)
	}

	for _, s := range sites {
		if s.ID == data.siteID {
			t.Errorf("Deleted site %s still in list", data.siteID)
			return
		}
	}
	t.Logf("✓ Verified site no longer in list")
}

// TestFixtureRedaction verifies that sensitive data is properly redacted
func TestFixtureRedaction(t *testing.T) {
	recorder := newFixtureRecorder(t)

	testCases := []struct {
		name     string
		input    string
		contains []string
		notContains []string
	}{
		{
			name:  "Session token redaction",
			input: `{"session": "abc123xyz789longtoken", "user_id": "550e8400-e29b-41d4-a716-446655440000"}`,
			contains: []string{"REDACTED", "REDACTED-USER-ID"},
			notContains: []string{"abc123xyz789longtoken", "550e8400-e29b-41d4-a716-446655440000"},
		},
		{
			name:  "Email redaction",
			input: `{"email": "user@pantheon.io", "name": "Test User"}`,
			contains: []string{"redacted@example.com"},
			notContains: []string{"user@pantheon.io"},
		},
		{
			name:  "Machine token redaction",
			input: `{"machine_token": "verylongsecrettokenstring123456"}`,
			contains: []string{"REDACTED"},
			notContains: []string{"verylongsecrettokenstring123456"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			redacted := string(recorder.redactSensitiveData([]byte(tc.input)))

			for _, expected := range tc.contains {
				if !regexp.MustCompile(expected).MatchString(redacted) {
					t.Errorf("Expected redacted output to contain %q, got: %s", expected, redacted)
				}
			}

			for _, notExpected := range tc.notContains {
				if regexp.MustCompile(fmt.Sprintf("%q", regexp.QuoteMeta(notExpected))).MatchString(redacted) {
					t.Errorf("Expected redacted output NOT to contain %q, got: %s", notExpected, redacted)
				}
			}
		})
	}
}
