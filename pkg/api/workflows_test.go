package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/deviantintegral/terminus-golang/pkg/api/models"
)

func TestNewWorkflowsService(t *testing.T) {
	client := NewClient()
	service := NewWorkflowsService(client)

	if service == nil {
		t.Fatal("expected service to be created")
	}

	if service.client != client {
		t.Error("expected service to have correct client")
	}
}

func TestWorkflowsService_List(t *testing.T) {
	siteID := "12345678-1234-1234-1234-123456789abc"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/sites/" + siteID + "/workflows"
		if r.URL.Path != expectedPath {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		workflows := []map[string]interface{}{
			{
				"id":          "wf1",
				"type":        "deploy",
				"description": "Deploy to environment",
				"site_id":     siteID,
				"result":      "succeeded",
				"finished_at": 1234567890.0,
			},
			{
				"id":          "wf2",
				"type":        "clear_cache",
				"description": "Clear cache",
				"site_id":     siteID,
				"result":      "succeeded",
				"finished_at": 1234567891.0,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(workflows)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	service := NewWorkflowsService(client)
	workflows, err := service.List(context.Background(), siteID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(workflows) != 2 {
		t.Errorf("expected 2 workflows, got %d", len(workflows))
	}

	if workflows[0].ID != "wf1" {
		t.Errorf("expected workflow ID 'wf1', got '%s'", workflows[0].ID)
	}

	if workflows[0].Type != "deploy" {
		t.Errorf("expected type 'deploy', got '%s'", workflows[0].Type)
	}
}

func TestWorkflowsService_ListForEnvironment(t *testing.T) {
	siteID := "12345678-1234-1234-1234-123456789abc"
	envID := "dev"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/sites/" + siteID + "/environments/" + envID + "/workflows"
		if r.URL.Path != expectedPath {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		workflows := []map[string]interface{}{
			{
				"id":          "wf1",
				"type":        "clear_cache",
				"description": "Clear cache for environment",
				"site_id":     siteID,
				"environment": envID,
				"result":      "succeeded",
				"finished_at": 1234567890.0,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(workflows)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	service := NewWorkflowsService(client)
	workflows, err := service.ListForEnvironment(context.Background(), siteID, envID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(workflows) != 1 {
		t.Errorf("expected 1 workflow, got %d", len(workflows))
	}

	if workflows[0].EnvironmentID != envID {
		t.Errorf("expected environment ID '%s', got '%s'", envID, workflows[0].EnvironmentID)
	}
}

func TestWorkflowsService_Get(t *testing.T) {
	siteID := "12345678-1234-1234-1234-123456789abc"
	workflowID := "wf-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/sites/" + siteID + "/workflows/" + workflowID
		if r.URL.Path != expectedPath {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		workflow := map[string]interface{}{
			"id":          workflowID,
			"type":        "deploy",
			"description": "Deploy to environment",
			"site_id":     siteID,
			"result":      "succeeded",
			"finished_at": 1234567890.0,
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(workflow)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	service := NewWorkflowsService(client)
	workflow, err := service.Get(context.Background(), siteID, workflowID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if workflow.ID != workflowID {
		t.Errorf("expected workflow ID '%s', got '%s'", workflowID, workflow.ID)
	}

	if workflow.Type != "deploy" {
		t.Errorf("expected type 'deploy', got '%s'", workflow.Type)
	}
}

func TestWorkflowsService_Wait(t *testing.T) {
	siteID := "12345678-1234-1234-1234-123456789abc"
	workflowID := "wf-123"

	// Track the number of polls
	pollCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/sites/" + siteID + "/workflows/" + workflowID
		if r.URL.Path != expectedPath {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		pollCount++

		var workflow map[string]interface{}
		if pollCount < 3 {
			// First two polls: workflow is still running
			workflow = map[string]interface{}{
				"id":          workflowID,
				"type":        "deploy",
				"description": "Deploy to environment",
				"site_id":     siteID,
				"result":      "",
				"finished_at": 0.0,
			}
		} else {
			// Third poll: workflow is finished
			workflow = map[string]interface{}{
				"id":          workflowID,
				"type":        "deploy",
				"description": "Deploy to environment",
				"site_id":     siteID,
				"result":      "succeeded",
				"finished_at": 1234567890.0,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(workflow)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	service := NewWorkflowsService(client)

	progressCallCount := 0
	opts := &WaitOptions{
		PollInterval: 100 * time.Millisecond,
		Timeout:      5 * time.Second,
		OnProgress: func(w *models.Workflow) {
			progressCallCount++
		},
	}

	workflow, err := service.Wait(context.Background(), siteID, workflowID, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if workflow.Result != "succeeded" {
		t.Errorf("expected result 'succeeded', got '%s'", workflow.Result)
	}

	if pollCount < 3 {
		t.Errorf("expected at least 3 polls, got %d", pollCount)
	}

	if progressCallCount < 3 {
		t.Errorf("expected at least 3 progress calls, got %d", progressCallCount)
	}
}

func TestWorkflowsService_Wait_Timeout(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping timeout test in short mode")
	}

	siteID := "12345678-1234-1234-1234-123456789abc"
	workflowID := "wf-123"

	pollCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pollCount++

		// Always return a workflow that's not finished
		workflow := map[string]interface{}{
			"id":          workflowID,
			"type":        "deploy",
			"description": "Deploy to environment",
			"site_id":     siteID,
			"result":      "",
			"finished_at": 0.0,
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(workflow)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	service := NewWorkflowsService(client)

	opts := &WaitOptions{
		PollInterval: 50 * time.Millisecond,
		Timeout:      200 * time.Millisecond, // Short timeout
	}

	_, err := service.Wait(context.Background(), siteID, workflowID, opts)
	if err == nil {
		t.Fatal("expected timeout error")
	}

	// Verify that we polled multiple times before timing out
	if pollCount < 2 {
		t.Errorf("expected at least 2 polls, got %d", pollCount)
	}
}

func TestWorkflowsService_Wait_DefaultOptions(t *testing.T) {
	siteID := "12345678-1234-1234-1234-123456789abc"
	workflowID := "wf-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		workflow := map[string]interface{}{
			"id":          workflowID,
			"type":        "deploy",
			"description": "Deploy to environment",
			"site_id":     siteID,
			"result":      "succeeded",
			"finished_at": 1234567890.0,
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(workflow)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	service := NewWorkflowsService(client)

	// Pass nil options to use defaults
	workflow, err := service.Wait(context.Background(), siteID, workflowID, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if workflow.Result != "succeeded" {
		t.Errorf("expected result 'succeeded', got '%s'", workflow.Result)
	}
}

func TestDefaultWaitOptions(t *testing.T) {
	opts := DefaultWaitOptions()

	if opts.PollInterval != 3*time.Second {
		t.Errorf("expected poll interval 3s, got %v", opts.PollInterval)
	}

	if opts.Timeout != 30*time.Minute {
		t.Errorf("expected timeout 30m, got %v", opts.Timeout)
	}
}

func TestWorkflowsService_Watch(t *testing.T) {
	siteID := "12345678-1234-1234-1234-123456789abc"
	workflowID := "wf-123"

	// Track the number of polls
	pollCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/sites/" + siteID + "/workflows/" + workflowID
		if r.URL.Path != expectedPath {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		pollCount++

		var workflow map[string]interface{}
		if pollCount < 3 {
			// First two polls: workflow is running
			workflow = map[string]interface{}{
				"id":                workflowID,
				"type":              "deploy",
				"description":       "Deploy to environment",
				"site_id":           siteID,
				"result":            "",
				"current_operation": "sync_code",
				"step":              pollCount,
				"finished_at":       0.0,
			}
		} else {
			// Third poll: workflow is finished
			workflow = map[string]interface{}{
				"id":                workflowID,
				"type":              "deploy",
				"description":       "Deploy to environment",
				"site_id":           siteID,
				"result":            "succeeded",
				"current_operation": "done",
				"step":              3,
				"finished_at":       1234567890.0,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(workflow)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	service := NewWorkflowsService(client)

	updateCallCount := 0
	opts := &WatchOptions{
		PollInterval: 100 * time.Millisecond,
		OnUpdate: func(w *models.Workflow) {
			updateCallCount++
		},
	}

	err := service.Watch(context.Background(), siteID, workflowID, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if updateCallCount < 3 {
		t.Errorf("expected at least 3 update calls, got %d", updateCallCount)
	}
}

func TestWorkflowsService_Watch_NoCallback(t *testing.T) {
	siteID := "12345678-1234-1234-1234-123456789abc"
	workflowID := "wf-123"

	client := NewClient()
	service := NewWorkflowsService(client)

	// Pass nil options (no callback)
	err := service.Watch(context.Background(), siteID, workflowID, nil)
	if err == nil {
		t.Fatal("expected error for missing OnUpdate callback")
	}

	// Pass options without callback
	opts := &WatchOptions{
		PollInterval: 100 * time.Millisecond,
	}
	err = service.Watch(context.Background(), siteID, workflowID, opts)
	if err == nil {
		t.Fatal("expected error for missing OnUpdate callback")
	}
}

func TestWorkflowsService_CreateForUser(t *testing.T) {
	userID := "user-123"
	workflowType := "create_site"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		expectedPath := "/users/" + userID + "/workflows"
		if r.URL.Path != expectedPath {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if reqBody["type"] != workflowType {
			t.Errorf("expected type '%s', got '%v'", workflowType, reqBody["type"])
		}

		params, ok := reqBody["params"].(map[string]interface{})
		if !ok {
			t.Fatal("expected params to be a map")
		}

		if params["site_name"] != "my-site" {
			t.Errorf("expected site_name 'my-site', got '%v'", params["site_name"])
		}

		workflow := map[string]interface{}{
			"id":          "wf-123",
			"type":        workflowType,
			"description": "Create site",
			"user_id":     userID,
			"result":      "",
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(workflow)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	service := NewWorkflowsService(client)
	workflow, err := service.CreateForUser(context.Background(), userID, workflowType, map[string]interface{}{
		"site_name": "my-site",
		"label":     "My Site",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if workflow.Type != workflowType {
		t.Errorf("expected type '%s', got '%s'", workflowType, workflow.Type)
	}

	if workflow.UserID != userID {
		t.Errorf("expected user ID '%s', got '%s'", userID, workflow.UserID)
	}
}

func TestWorkflowsService_GetForUser(t *testing.T) {
	userID := "user-123"
	workflowID := "wf-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/users/" + userID + "/workflows/" + workflowID
		if r.URL.Path != expectedPath {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		workflow := map[string]interface{}{
			"id":          workflowID,
			"type":        "create_site",
			"description": "Create site",
			"user_id":     userID,
			"result":      "succeeded",
			"finished_at": 1234567890.0,
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(workflow)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	service := NewWorkflowsService(client)
	workflow, err := service.GetForUser(context.Background(), userID, workflowID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if workflow.ID != workflowID {
		t.Errorf("expected workflow ID '%s', got '%s'", workflowID, workflow.ID)
	}

	if workflow.UserID != userID {
		t.Errorf("expected user ID '%s', got '%s'", userID, workflow.UserID)
	}
}

func TestWorkflowsService_WaitForUser(t *testing.T) {
	userID := "user-123"
	workflowID := "wf-123"

	// Track the number of polls
	pollCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/users/" + userID + "/workflows/" + workflowID
		if r.URL.Path != expectedPath {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		pollCount++

		var workflow map[string]interface{}
		if pollCount < 3 {
			// First two polls: workflow is still running
			workflow = map[string]interface{}{
				"id":          workflowID,
				"type":        "create_site",
				"description": "Create site",
				"user_id":     userID,
				"result":      "",
				"finished_at": 0.0,
			}
		} else {
			// Third poll: workflow is finished
			workflow = map[string]interface{}{
				"id":          workflowID,
				"type":        "create_site",
				"description": "Create site",
				"user_id":     userID,
				"result":      "succeeded",
				"finished_at": 1234567890.0,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(workflow)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	service := NewWorkflowsService(client)

	progressCallCount := 0
	opts := &WaitOptions{
		PollInterval: 100 * time.Millisecond,
		Timeout:      5 * time.Second,
		OnProgress: func(w *models.Workflow) {
			progressCallCount++
		},
	}

	workflow, err := service.WaitForUser(context.Background(), userID, workflowID, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if workflow.Result != "succeeded" {
		t.Errorf("expected result 'succeeded', got '%s'", workflow.Result)
	}

	if pollCount < 3 {
		t.Errorf("expected at least 3 polls, got %d", pollCount)
	}

	if progressCallCount < 3 {
		t.Errorf("expected at least 3 progress calls, got %d", progressCallCount)
	}
}

func TestWorkflowsService_WaitForUser_DefaultOptions(t *testing.T) {
	userID := "user-123"
	workflowID := "wf-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		workflow := map[string]interface{}{
			"id":          workflowID,
			"type":        "create_site",
			"description": "Create site",
			"user_id":     userID,
			"result":      "succeeded",
			"finished_at": 1234567890.0,
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(workflow)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	service := NewWorkflowsService(client)

	// Pass nil options to use defaults
	workflow, err := service.WaitForUser(context.Background(), userID, workflowID, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if workflow.Result != "succeeded" {
		t.Errorf("expected result 'succeeded', got '%s'", workflow.Result)
	}
}

func TestWorkflowsService_CreateForSite(t *testing.T) {
	siteID := "12345678-1234-1234-1234-123456789abc"
	workflowType := "clear_cache"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		expectedPath := "/sites/" + siteID + "/workflows"
		if r.URL.Path != expectedPath {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if reqBody["type"] != workflowType {
			t.Errorf("expected type '%s', got '%v'", workflowType, reqBody["type"])
		}

		workflow := map[string]interface{}{
			"id":          "wf-123",
			"type":        workflowType,
			"description": "Clear cache",
			"site_id":     siteID,
			"result":      "",
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(workflow)
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithToken("test-token"),
		WithHTTPClient(&http.Client{Timeout: 5 * time.Second}),
	)

	service := NewWorkflowsService(client)
	workflow, err := service.CreateForSite(context.Background(), siteID, workflowType, map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if workflow.Type != workflowType {
		t.Errorf("expected type '%s', got '%s'", workflowType, workflow.Type)
	}

	if workflow.SiteID != siteID {
		t.Errorf("expected site ID '%s', got '%s'", siteID, workflow.SiteID)
	}
}
