package api

import (
	"context"
	"fmt"

	"github.com/pantheon-systems/terminus-go/pkg/api/models"
)

// MultidevService handles multidev environment operations
type MultidevService struct {
	client *Client
}

// NewMultidevService creates a new multidev service
func NewMultidevService(client *Client) *MultidevService {
	return &MultidevService{client: client}
}

// CreateRequest represents a multidev creation request
type CreateMultidevRequest struct {
	FromEnvironment string `json:"from_environment"`
	CloudDevelopmentEnvironmentID string `json:"cloud_development_environment_id"`
}

// Create creates a new multidev environment
func (s *MultidevService) Create(ctx context.Context, siteID, envName, fromEnv string) (*models.Workflow, error) {
	path := fmt.Sprintf("/sites/%s/environments", siteID)

	req := map[string]interface{}{
		"type": "create_cloud_development_environment",
		"params": map[string]interface{}{
			"environment_id":     envName,
			"deploy_from_environment": fromEnv,
		},
	}

	resp, err := s.client.Post(ctx, path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create multidev: %w", err)
	}

	var workflow models.Workflow
	if err := DecodeResponse(resp, &workflow); err != nil {
		return nil, err
	}

	return &workflow, nil
}

// Delete deletes a multidev environment
func (s *MultidevService) Delete(ctx context.Context, siteID, envID string, deleteDB bool) (*models.Workflow, error) {
	path := fmt.Sprintf("/sites/%s/environments/%s", siteID, envID)

	req := map[string]interface{}{
		"delete_branch": true,
	}
	if deleteDB {
		req["delete_db"] = true
	}

	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to delete multidev: %w", err)
	}

	var workflow models.Workflow
	if err := DecodeResponse(resp, &workflow); err != nil {
		return nil, err
	}

	return &workflow, nil
}

// MergeToDev merges a multidev to dev
func (s *MultidevService) MergeToDev(ctx context.Context, siteID, envID string, updateDB bool) (*models.Workflow, error) {
	path := fmt.Sprintf("/sites/%s/environments/%s/workflows", siteID, envID)

	req := map[string]interface{}{
		"type": "merge_dev_into_cloud_development_environment",
		"params": map[string]interface{}{
			"updatedb": updateDB,
		},
	}

	resp, err := s.client.Post(ctx, path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to merge to dev: %w", err)
	}

	var workflow models.Workflow
	if err := DecodeResponse(resp, &workflow); err != nil {
		return nil, err
	}

	return &workflow, nil
}

// MergeFromDev merges dev into a multidev
func (s *MultidevService) MergeFromDev(ctx context.Context, siteID, envID string, updateDB bool) (*models.Workflow, error) {
	path := fmt.Sprintf("/sites/%s/environments/%s/workflows", siteID, envID)

	req := map[string]interface{}{
		"type": "merge_cloud_development_environment_into_dev",
		"params": map[string]interface{}{
			"updatedb": updateDB,
		},
	}

	resp, err := s.client.Post(ctx, path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to merge from dev: %w", err)
	}

	var workflow models.Workflow
	if err := DecodeResponse(resp, &workflow); err != nil {
		return nil, err
	}

	return &workflow, nil
}
