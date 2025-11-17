package api

import (
	"context"
	"encoding/json"
	"fmt"

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
func (s *SitesService) Get(ctx context.Context, siteID string) (*models.Site, error) {
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

// CreateSiteRequest represents a site creation request
type CreateSiteRequest struct {
	SiteName     string `json:"site_name"`
	Label        string `json:"label,omitempty"`
	UpstreamID   string `json:"upstream_id"`
	Organization string `json:"organization_id,omitempty"`
	Region       string `json:"preferred_zone,omitempty"`
}

// Create creates a new site
func (s *SitesService) Create(ctx context.Context, req *CreateSiteRequest) (*models.Site, error) {
	resp, err := s.client.Post(ctx, "/sites", req) //nolint:bodyclose // DecodeResponse closes body
	if err != nil {
		return nil, fmt.Errorf("failed to create site: %w", err)
	}

	var site models.Site
	if err := DecodeResponse(resp, &site); err != nil {
		return nil, err
	}

	return &site, nil
}

// Delete deletes a site
func (s *SitesService) Delete(ctx context.Context, siteID string) error {
	path := fmt.Sprintf("/sites/%s", siteID)
	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to delete site: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("delete failed with status %d", resp.StatusCode)
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
