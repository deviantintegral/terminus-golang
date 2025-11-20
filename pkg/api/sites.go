package api

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pantheon-systems/terminus-go/pkg/api/models"
)

// SitesService handles site-related operations
type SitesService struct {
	client *Client
}

// NewSitesService creates a new sites service
func NewSitesService(client *Client) *SitesService {
	return &SitesService{client: client}
}

// List returns all sites accessible to the authenticated user
func (s *SitesService) List(ctx context.Context, userID string) ([]*models.Site, error) {
	// Get user sites using memberships endpoint
	path := fmt.Sprintf("/users/%s/memberships/sites", userID)

	rawResults, err := s.client.GetPaged(ctx, path)
	if err != nil {
		return nil, err
	}

	sites := make([]*models.Site, 0, len(rawResults))
	for _, raw := range rawResults {
		var membership struct {
			Site *models.Site `json:"site"`
		}
		if err := json.Unmarshal(raw, &membership); err != nil {
			// Try direct unmarshal in case the API returns sites directly
			var site models.Site
			if err := json.Unmarshal(raw, &site); err != nil {
				return nil, fmt.Errorf("failed to decode site: %w", err)
			}
			sites = append(sites, &site)
		} else if membership.Site != nil {
			sites = append(sites, membership.Site)
		}
	}

	return sites, nil
}

// Get returns a specific site by ID or name
func (s *SitesService) Get(ctx context.Context, siteIdentifier string) (*models.Site, error) {
	// First, try to resolve the identifier to a UUID if it's a name
	siteID := siteIdentifier

	// Check if the identifier looks like a UUID (contains dashes in UUID pattern)
	// If not, try to resolve it as a name
	if !isUUID(siteIdentifier) {
		resolvedID, err := s.resolveNameToID(ctx, siteIdentifier)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve site name: %w", err)
		}
		siteID = resolvedID
	}

	path := fmt.Sprintf("/sites/%s", siteID)
	resp, err := s.client.Get(ctx, path) //nolint:bodyclose // DecodeResponse closes body
	if err != nil {
		return nil, fmt.Errorf("failed to get site: %w", err)
	}

	var site models.Site
	if err := DecodeResponse(resp, &site); err != nil {
		return nil, err
	}

	return &site, nil
}

// resolveNameToID converts a site name to a site UUID
func (s *SitesService) resolveNameToID(ctx context.Context, siteName string) (string, error) {
	path := fmt.Sprintf("/site-names/%s", siteName)
	resp, err := s.client.Get(ctx, path) //nolint:bodyclose // DecodeResponse closes body
	if err != nil {
		return "", fmt.Errorf("site not found: %w", err)
	}

	var result struct {
		ID string `json:"id"`
	}
	if err := DecodeResponse(resp, &result); err != nil {
		return "", err
	}

	if result.ID == "" {
		return "", fmt.Errorf("site name resolution returned empty ID")
	}

	return result.ID, nil
}

// isUUID checks if a string looks like a UUID
func isUUID(s string) bool {
	// Simple check: UUIDs have dashes and are 36 characters long
	// Format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	return len(s) == 36 && s[8] == '-' && s[13] == '-' && s[18] == '-' && s[23] == '-'
}

// CreateSiteRequest represents a site creation request
type CreateSiteRequest struct {
	SiteName     string `json:"site_name"`
	Label        string `json:"label,omitempty"`
	UpstreamID   string `json:"upstream_id"`
	Organization string `json:"organization_id,omitempty"`
	Region       string `json:"preferred_zone,omitempty"`
}

// Create creates a new site using workflows
func (s *SitesService) Create(ctx context.Context, userID string, req *CreateSiteRequest) (*models.Site, error) {
	workflowsService := NewWorkflowsService(s.client)

	// Step 1: Create the site via user workflow
	createParams := map[string]interface{}{
		"site_name": req.SiteName,
		"label":     req.Label,
	}

	// Add optional organization and region if provided
	if req.Organization != "" {
		createParams["organization_id"] = req.Organization
	}
	if req.Region != "" {
		createParams["preferred_zone"] = req.Region
	}

	createWorkflow, err := workflowsService.CreateForUser(ctx, userID, "create_site", createParams)
	if err != nil {
		return nil, fmt.Errorf("failed to start site creation workflow: %w", err)
	}

	// Wait for the site creation workflow to complete
	completedWorkflow, err := workflowsService.WaitForUser(ctx, userID, createWorkflow.ID, nil)
	if err != nil {
		return nil, fmt.Errorf("site creation workflow failed: %w", err)
	}

	// Extract the site name from the workflow params to look up the site
	siteName, ok := completedWorkflow.Params["site_name"].(string)
	if !ok || siteName == "" {
		return nil, fmt.Errorf("failed to get site_name from workflow params")
	}

	// Get the created site details from the user's site list
	// We do this instead of using Get() because the site-names endpoint
	// may not be immediately available after site creation
	// Retry a few times with delays to handle propagation
	var site *models.Site
	retryDelays := []time.Duration{
		0,               // immediate
		1 * time.Second, // 1s
		2 * time.Second, // 2s
		4 * time.Second, // 4s
		8 * time.Second, // 8s
	}

	for _, delay := range retryDelays {
		if delay > 0 {
			time.Sleep(delay)
		}

		sites, err := s.List(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to list sites after creation: %w", err)
		}

		// Find the site with the matching name
		for _, s := range sites {
			if s.Name == siteName {
				site = s
				break
			}
		}

		if site != nil {
			break
		}
	}

	if site == nil {
		return nil, fmt.Errorf("created site %s not found in site list after retries", siteName)
	}

	return site, nil

}

// Delete deletes a site using the delete_site workflow
func (s *SitesService) Delete(ctx context.Context, siteID string) error {
	workflowsService := NewWorkflowsService(s.client)

	// Create a delete_site workflow
	workflow, err := workflowsService.CreateForSite(ctx, siteID, "delete_site", map[string]interface{}{})
	if err != nil {
		return fmt.Errorf("failed to start site deletion workflow: %w", err)
	}

	// Wait for the workflow to complete using the user endpoint
	// Note: We can't use the site endpoint to monitor the workflow because
	// the site is deleted during workflow execution. Instead, we use the user
	// endpoint which remains available throughout the deletion process.
	completedWorkflow, err := workflowsService.WaitForUser(ctx, workflow.UserID, workflow.ID, nil)
	if err != nil {
		return fmt.Errorf("site deletion workflow failed: %w", err)
	}

	if !completedWorkflow.IsSuccessful() {
		return fmt.Errorf("site deletion workflow failed: %s", completedWorkflow.GetMessage())
	}

	return nil
}

// UpdateRequest represents a site update request
type UpdateRequest struct {
	Label        string `json:"label,omitempty"`
	ServiceLevel string `json:"service_level,omitempty"`
}

// Update updates a site
func (s *SitesService) Update(ctx context.Context, siteID string, req *UpdateRequest) (*models.Site, error) {
	path := fmt.Sprintf("/sites/%s", siteID)
	resp, err := s.client.Put(ctx, path, req) //nolint:bodyclose // DecodeResponse closes body
	if err != nil {
		return nil, fmt.Errorf("failed to update site: %w", err)
	}

	var site models.Site
	if err := DecodeResponse(resp, &site); err != nil {
		return nil, err
	}

	return &site, nil
}

// ListByOrganization returns sites for a specific organization
func (s *SitesService) ListByOrganization(ctx context.Context, orgID string) ([]*models.Site, error) {
	path := fmt.Sprintf("/organizations/%s/memberships/sites", orgID)

	rawResults, err := s.client.GetPaged(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list organization sites: %w", err)
	}

	sites := make([]*models.Site, 0, len(rawResults))
	for _, raw := range rawResults {
		var membership struct {
			Site *models.Site `json:"site"`
		}
		if err := json.Unmarshal(raw, &membership); err != nil {
			return nil, fmt.Errorf("failed to decode site membership: %w", err)
		}
		if membership.Site != nil {
			sites = append(sites, membership.Site)
		}
	}

	return sites, nil
}

// GetTeam returns team members for a site
func (s *SitesService) GetTeam(ctx context.Context, siteID string) ([]*models.TeamMember, error) {
	path := fmt.Sprintf("/sites/%s/team", siteID)
	resp, err := s.client.Get(ctx, path) //nolint:bodyclose // DecodeResponse closes body
	if err != nil {
		return nil, fmt.Errorf("failed to get site team: %w", err)
	}

	var team []*models.TeamMember
	if err := DecodeResponse(resp, &team); err != nil {
		return nil, err
	}

	return team, nil
}

// AddTeamMemberRequest represents a request to add a team member
type AddTeamMemberRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

// AddTeamMember adds a team member to a site
func (s *SitesService) AddTeamMember(ctx context.Context, siteID string, req *AddTeamMemberRequest) (*models.TeamMember, error) {
	path := fmt.Sprintf("/sites/%s/team", siteID)
	resp, err := s.client.Post(ctx, path, req) //nolint:bodyclose // DecodeResponse closes body
	if err != nil {
		return nil, fmt.Errorf("failed to add team member: %w", err)
	}

	var member models.TeamMember
	if err := DecodeResponse(resp, &member); err != nil {
		return nil, err
	}

	return &member, nil
}

// RemoveTeamMember removes a team member from a site
func (s *SitesService) RemoveTeamMember(ctx context.Context, siteID, userID string) error {
	path := fmt.Sprintf("/sites/%s/team/%s", siteID, userID)
	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to remove team member: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("remove team member failed with status %d", resp.StatusCode)
	}

	return nil
}

// GetTags returns tags for a site
func (s *SitesService) GetTags(ctx context.Context, siteID, orgID string) ([]*models.Tag, error) {
	path := fmt.Sprintf("/organizations/%s/tags/sites/%s", orgID, siteID)
	resp, err := s.client.Get(ctx, path) //nolint:bodyclose // DecodeResponse closes body
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	var tags []*models.Tag
	if err := DecodeResponse(resp, &tags); err != nil {
		return nil, err
	}

	return tags, nil
}

// AddTag adds a tag to a site
func (s *SitesService) AddTag(ctx context.Context, siteID, orgID, tagName string) error {
	path := fmt.Sprintf("/organizations/%s/tags/%s/sites", orgID, tagName)

	req := struct {
		SiteID string `json:"site_id"`
	}{
		SiteID: siteID,
	}

	resp, err := s.client.Post(ctx, path, req)
	if err != nil {
		return fmt.Errorf("failed to add tag: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("add tag failed with status %d", resp.StatusCode)
	}

	return nil
}

// RemoveTag removes a tag from a site
func (s *SitesService) RemoveTag(ctx context.Context, siteID, orgID, tagName string) error {
	path := fmt.Sprintf("/organizations/%s/tags/%s/sites/%s", orgID, tagName, siteID)
	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to remove tag: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("remove tag failed with status %d", resp.StatusCode)
	}

	return nil
}

// GetPlan returns the plan for a site
func (s *SitesService) GetPlan(ctx context.Context, siteID string) (*models.Plan, error) {
	path := fmt.Sprintf("/sites/%s/plan", siteID)
	resp, err := s.client.Get(ctx, path) //nolint:bodyclose // DecodeResponse closes body
	if err != nil {
		return nil, fmt.Errorf("failed to get site plan: %w", err)
	}

	var plan models.Plan
	if err := DecodeResponse(resp, &plan); err != nil {
		return nil, err
	}

	return &plan, nil
}

// ListBranches returns git branches for a site
func (s *SitesService) ListBranches(ctx context.Context, siteID string) ([]*models.Branch, error) {
	path := fmt.Sprintf("/sites/%s/code-tips", siteID)

	rawResults, err := s.client.GetPaged(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list branches: %w", err)
	}

	branches := make([]*models.Branch, 0, len(rawResults))
	for _, raw := range rawResults {
		var branch models.Branch
		if err := json.Unmarshal(raw, &branch); err != nil {
			return nil, fmt.Errorf("failed to decode branch: %w", err)
		}
		branches = append(branches, &branch)
	}

	return branches, nil
}

// GetPlans returns available plans for a site
func (s *SitesService) GetPlans(ctx context.Context, siteID string) ([]*models.Plan, error) {
	path := fmt.Sprintf("/sites/%s/plans", siteID)

	rawResults, err := s.client.GetPaged(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get plans: %w", err)
	}

	plans := make([]*models.Plan, 0, len(rawResults))
	for _, raw := range rawResults {
		var plan models.Plan
		if err := json.Unmarshal(raw, &plan); err != nil {
			return nil, fmt.Errorf("failed to decode plan: %w", err)
		}
		plans = append(plans, &plan)
	}

	return plans, nil
}

// ListOrganizations returns organizations that a site belongs to
func (s *SitesService) ListOrganizations(ctx context.Context, siteID string) ([]*models.SiteOrganizationMembership, error) {
	path := fmt.Sprintf("/sites/%s/memberships/organizations", siteID)

	rawResults, err := s.client.GetPaged(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list site organizations: %w", err)
	}

	memberships := make([]*models.SiteOrganizationMembership, 0, len(rawResults))
	for _, raw := range rawResults {
		var membership struct {
			Organization *models.Organization `json:"organization"`
		}
		if err := json.Unmarshal(raw, &membership); err != nil {
			return nil, fmt.Errorf("failed to decode organization membership: %w", err)
		}
		if membership.Organization != nil {
			memberships = append(memberships, &models.SiteOrganizationMembership{
				OrgID:   membership.Organization.ID,
				OrgName: membership.Organization.Name,
			})
		}
	}

	return memberships, nil
}
