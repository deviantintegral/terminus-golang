package api

import (
	"context"
	"fmt"
	"time"

	"github.com/pantheon-systems/terminus-go/pkg/api/models"
)

// WorkflowsService handles workflow-related operations
type WorkflowsService struct {
	client *Client
}

// NewWorkflowsService creates a new workflows service
func NewWorkflowsService(client *Client) *WorkflowsService {
	return &WorkflowsService{client: client}
}

// List returns all workflows for a site
func (s *WorkflowsService) List(ctx context.Context, siteID string) ([]*models.Workflow, error) {
	path := fmt.Sprintf("/sites/%s/workflows", siteID)
	resp, err := s.client.Get(ctx, path) //nolint:bodyclose // DecodeResponse closes body
	if err != nil {
		return nil, fmt.Errorf("failed to list workflows: %w", err)
	}

	var workflows []*models.Workflow
	if err := DecodeResponse(resp, &workflows); err != nil {
		return nil, err
	}

	return workflows, nil
}

// ListForEnvironment returns workflows for a specific environment
func (s *WorkflowsService) ListForEnvironment(ctx context.Context, siteID, envID string) ([]*models.Workflow, error) {
	path := fmt.Sprintf("/sites/%s/environments/%s/workflows", siteID, envID)
	resp, err := s.client.Get(ctx, path) //nolint:bodyclose // DecodeResponse closes body
	if err != nil {
		return nil, fmt.Errorf("failed to list environment workflows: %w", err)
	}

	var workflows []*models.Workflow
	if err := DecodeResponse(resp, &workflows); err != nil {
		return nil, err
	}

	return workflows, nil
}

// Get returns a specific workflow
func (s *WorkflowsService) Get(ctx context.Context, siteID, workflowID string) (*models.Workflow, error) {
	path := fmt.Sprintf("/sites/%s/workflows/%s", siteID, workflowID)
	resp, err := s.client.Get(ctx, path) //nolint:bodyclose // DecodeResponse closes body
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow: %w", err)
	}

	var workflow models.Workflow
	if err := DecodeResponse(resp, &workflow); err != nil {
		return nil, err
	}

	return &workflow, nil
}

// WaitOptions configures workflow wait behavior
type WaitOptions struct {
	// PollInterval is how often to check workflow status
	PollInterval time.Duration
	// Timeout is the maximum time to wait
	Timeout time.Duration
	// OnProgress is called on each poll with the current workflow state
	OnProgress func(*models.Workflow)
}

// DefaultWaitOptions returns default wait options
func DefaultWaitOptions() *WaitOptions {
	return &WaitOptions{
		PollInterval: 3 * time.Second,
		Timeout:      30 * time.Minute,
	}
}

// Wait waits for a workflow to complete
func (s *WorkflowsService) Wait(ctx context.Context, siteID, workflowID string, opts *WaitOptions) (*models.Workflow, error) {
	if opts == nil {
		opts = DefaultWaitOptions()
	}

	// Create a context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, opts.Timeout)
	defer cancel()

	ticker := time.NewTicker(opts.PollInterval)
	defer ticker.Stop()

	for {
		workflow, err := s.Get(timeoutCtx, siteID, workflowID)
		if err != nil {
			return nil, fmt.Errorf("failed to check workflow status: %w", err)
		}

		if opts.OnProgress != nil {
			opts.OnProgress(workflow)
		}

		if workflow.IsFinished() {
			return workflow, nil
		}

		select {
		case <-timeoutCtx.Done():
			return nil, fmt.Errorf("workflow did not complete within timeout")
		case <-ticker.C:
			// Continue polling
		}
	}
}

// WatchOptions configures workflow watch behavior
type WatchOptions struct {
	PollInterval time.Duration
	OnUpdate     func(*models.Workflow)
}

// Watch watches a workflow and calls OnUpdate on each status change
func (s *WorkflowsService) Watch(ctx context.Context, siteID, workflowID string, opts *WatchOptions) error {
	if opts == nil || opts.OnUpdate == nil {
		return fmt.Errorf("OnUpdate callback is required")
	}

	pollInterval := 3 * time.Second
	if opts.PollInterval > 0 {
		pollInterval = opts.PollInterval
	}

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	var lastStatus string

	for {
		workflow, err := s.Get(ctx, siteID, workflowID)
		if err != nil {
			return fmt.Errorf("failed to check workflow status: %w", err)
		}

		// Call update callback if status changed
		currentStatus := fmt.Sprintf("%s:%s:%d", workflow.Result, workflow.CurrentOperation, workflow.Step)
		if currentStatus != lastStatus {
			opts.OnUpdate(workflow)
			lastStatus = currentStatus
		}

		if workflow.IsFinished() {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			// Continue watching
		}
	}
}
