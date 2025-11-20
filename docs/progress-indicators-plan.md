# Progress Indicators Plan

## Executive Summary

This document outlines a comprehensive plan for implementing rich progress indicators in terminus-golang for commands that involve multiple Pantheon API calls or require waiting for workflows to complete. The goal is to provide users with clear, informative feedback during long-running operations that can take 2-30 minutes to complete.

## Current State

### Existing Implementation

**Location:** `internal/commands/workflow.go:143-186`

The current `waitForWorkflow()` function uses a basic indeterminate spinner from the `github.com/schollz/progressbar/v3` library (v3.18.0):

```go
bar = progressbar.NewOptions(-1,
    progressbar.OptionSetDescription(description),
    progressbar.OptionSetWriter(os.Stderr),
    progressbar.OptionSpinnerType(14),
    progressbar.OptionFullWidth(),
)
```

**Limitations:**
- Shows only a static description (e.g., "Creating backup")
- No operation details or current step information
- No elapsed time or progress indication
- No visibility into what operation is currently running
- Users cannot distinguish between a 30-second operation and a 10-minute operation

### Workflow Data Available

**Location:** `pkg/api/models/models.go:48-67`

The `Workflow` model provides rich data that is currently unused:

```go
type Workflow struct {
    ID               string
    Type             string
    Description      string
    StartedAt        float64
    FinishedAt       float64
    TotalTime        float64
    CurrentOperation string    // Currently running operation
    Step             int       // Current step number
    Result           string    // succeeded/failed/aborted/running
    Operations       []Operation  // Historical operation list
    Active           bool
    HasActiveOps     bool
}

type Operation struct {
    ID          string
    Type        string
    Description string
    Result      string
    Duration    float64
}
```

### Commands Requiring Progress Indicators

**15+ commands** across multiple categories trigger workflows:

1. **Backup Operations** (4 commands)
   - `backup:create` - Creates backup (5-10 minutes)
   - `backup:restore` - Restores from backup (10-15 minutes)

2. **Environment Operations** (5 commands)
   - `env:deploy` - Deploys code (2-5 minutes)
   - `env:clone-content` - Clones database/files (5-15 minutes)
   - `env:clear-cache` - Clears environment cache (1-3 minutes)
   - `env:commit` - Commits code changes (1-3 minutes)
   - `connection:set` - Modifies connection mode (2-5 minutes)

3. **Multidev Operations** (4 commands)
   - `multidev:create` - Creates new multidev (3-8 minutes)
   - `multidev:delete` - Deletes multidev (2-5 minutes)
   - `multidev:merge-to-dev` - Merges to dev (2-5 minutes)
   - `multidev:merge-from-dev` - Merges from dev (2-5 minutes)

4. **Feature Operations** (2 commands)
   - `redis:enable` - Enables Redis (2-5 minutes)
   - `redis:disable` - Disables Redis (2-5 minutes)

All these commands currently use: `waitForWorkflow(siteID, workflowID, description)`

## Requirements

### Functional Requirements

1. **Real-time Operation Display**
   - Show current operation being executed
   - Display step number (e.g., "Step 3 of 5" if total steps known)
   - Update dynamically as workflow progresses

2. **Time Information**
   - Show elapsed time since workflow started
   - Update time display continuously (e.g., "Running for 2m 34s")

3. **Status Indication**
   - Visual spinner/animation to show activity
   - Clear final status (success/failure)
   - Color coding for different states (if terminal supports it)

4. **Error Information**
   - Display failure details when workflow fails
   - Show which operation failed
   - Include error messages from the workflow

5. **User Control**
   - Respect `--quiet` flag (suppress progress indicators)
   - Support `--format=json` output (no progress bars in structured output)
   - Allow graceful interruption (Ctrl+C)

### Non-Functional Requirements

1. **Performance**
   - Minimal overhead (<100ms per update)
   - No impact on API polling frequency (3 seconds)
   - Efficient terminal updates

2. **Compatibility**
   - Work on all terminal types
   - Graceful degradation for non-TTY output (pipes, redirects)
   - Support both Linux and macOS

3. **Maintainability**
   - Reusable progress indicator components
   - Consistent interface across all commands
   - Easy to test

## Proposed Solution

### Architecture

#### 1. Progress Indicator Abstraction Layer

**Location:** `pkg/output/progress.go` (new file)

Create reusable progress indicator types:

```go
// ProgressIndicator manages progress display for long-running operations
type ProgressIndicator interface {
    Start()
    Update(info ProgressInfo)
    Finish(result ProgressResult)
    Error(err error)
}

// ProgressInfo contains current operation state
type ProgressInfo struct {
    CurrentOperation string
    Step            int
    TotalSteps      int  // 0 if unknown
    ElapsedTime     time.Duration
    Message         string
}

// ProgressResult contains final operation result
type ProgressResult struct {
    Success     bool
    Duration    time.Duration
    Message     string
}
```

#### 2. Workflow Progress Indicator Implementation

**Location:** `pkg/output/workflow_progress.go` (new file)

Concrete implementation using progressbar/v3:

```go
type WorkflowProgressIndicator struct {
    bar         *progressbar.ProgressBar
    startTime   time.Time
    description string
}
```

**Features:**
- Dynamic description updates
- Elapsed time tracking
- Step counter display
- Color-coded output (green for success, red for failure)

#### 3. Enhanced Callback in WaitOptions

**Location:** `pkg/api/workflows.go:69-77`

Enhance the `OnProgress` callback to provide richer information:

```go
type WaitOptions struct {
    PollInterval time.Duration
    Timeout      time.Duration
    // Enhanced callback with start time and elapsed duration
    OnProgress   func(workflow *models.Workflow, elapsed time.Duration)
}
```

#### 4. Updated waitForWorkflow Helper

**Location:** `internal/commands/workflow.go:143-186`

Refactor to use the new progress indicator:

```go
func waitForWorkflow(siteID, workflowID, description string) error {
    workflowsService := api.NewWorkflowsService(cliContext.APIClient)

    // Create progress indicator
    var progress output.ProgressIndicator
    if !quietFlag && formatFlag == "" {
        progress = output.NewWorkflowProgressIndicator(description)
        progress.Start()
    }

    startTime := time.Now()

    opts := &api.WaitOptions{
        PollInterval: 3 * time.Second,
        Timeout:      30 * time.Minute,
        OnProgress: func(w *models.Workflow, elapsed time.Duration) {
            if progress != nil {
                progress.Update(output.ProgressInfo{
                    CurrentOperation: w.CurrentOperation,
                    Step:            w.Step,
                    ElapsedTime:     elapsed,
                })
            }
        },
    }

    workflow, err := workflowsService.Wait(getContext(), siteID, workflowID, opts)

    if progress != nil {
        if err != nil {
            progress.Error(err)
        } else if workflow.IsSuccessful() {
            progress.Finish(output.ProgressResult{
                Success: true,
                Duration: time.Since(startTime),
            })
        } else {
            progress.Finish(output.ProgressResult{
                Success: false,
                Duration: time.Since(startTime),
                Message: workflow.GetMessage(),
            })
        }
    }

    return handleWorkflowResult(workflow, description)
}
```

### Display Format Examples

#### Basic Progress (Current Operation)
```
Creating backup... [⠋] Step 2 | backup_database | Elapsed: 1m 23s
```

#### With Step Count (if available)
```
Creating backup... [⠙] Step 2/5 | clone_files | Elapsed: 3m 45s
```

#### Success
```
✓ Creating backup completed successfully in 5m 32s
```

#### Failure
```
✗ Creating backup failed after 2m 15s: Database connection timeout
```

### Progressive Enhancement Strategy

Implement in phases to manage complexity and risk:

#### Phase 1: Enhanced Current Operation Display
**Effort:** Small | **Impact:** High

- Display `CurrentOperation` field in progress bar description
- Show elapsed time
- Minimal changes to existing code

**Example output:**
```
Creating backup [⠋] backup_database (1m 23s)
```

**Files to modify:**
- `internal/commands/workflow.go` - Update `waitForWorkflow()` callback

**Changes:**
```go
OnProgress: func(w *models.Workflow) {
    if bar != nil {
        elapsed := time.Since(startTime)
        desc := fmt.Sprintf("%s [%s] %s",
            description,
            w.CurrentOperation,
            formatDuration(elapsed))
        bar.Describe(desc)
        _ = bar.Add(1)
    }
}
```

#### Phase 2: Step Counter and Time Display
**Effort:** Medium | **Impact:** Medium

- Add step counter to display
- Calculate total steps if possible (from Operations history)
- Improve time formatting

**Example output:**
```
Creating backup [⠙] Step 2/5 | clone_files | Elapsed: 3m 45s
```

**Files to modify:**
- `pkg/api/workflows.go` - Enhance `WaitOptions` to track start time
- `internal/commands/workflow.go` - Update display logic

#### Phase 3: Progress Abstraction Layer
**Effort:** Large | **Impact:** Medium (maintainability)

- Create `pkg/output/progress.go` with interfaces
- Implement `WorkflowProgressIndicator`
- Refactor all commands to use abstraction

**Benefits:**
- Easier testing (mock progress indicators)
- Consistent progress display across all commands
- Future enhancement flexibility (e.g., different progress styles)

**New files:**
- `pkg/output/progress.go` - Interfaces and types
- `pkg/output/workflow_progress.go` - Implementation
- `pkg/output/workflow_progress_test.go` - Tests

#### Phase 4: Advanced Features
**Effort:** Large | **Impact:** Low (nice-to-have)

Optional enhancements:
- ETA calculation based on operation history
- Percentage completion (if determinable)
- Detailed operation history display
- Interactive mode (show/hide details)

### Testing Strategy

#### Unit Tests

**Location:** `pkg/output/workflow_progress_test.go`

```go
func TestProgressIndicator_Update(t *testing.T) {
    // Test update logic
}

func TestProgressIndicator_Finish(t *testing.T) {
    // Test completion states
}
```

#### Integration Tests

**Location:** `internal/commands/workflow_test.go`

```go
func TestWaitForWorkflow_WithProgress(t *testing.T) {
    // Test with mock workflow service
}
```

#### Manual Testing

Test commands that require real API access:
- `backup:create` - Long-running, multiple operations
- `env:deploy` - Medium duration, clear steps
- `multidev:create` - Variable duration

Use `PANTHEON_MACHINE_TOKEN_TESTING` environment variable for live testing.

### Error Handling

1. **Non-TTY Environments**
   - Detect if output is a terminal using `term.IsTerminal()`
   - Fall back to simple text output when piped or redirected
   - Example: `Creating backup... Done (5m 32s)`

2. **Workflow Failures**
   - Display clear error message from workflow
   - Show which operation failed
   - Include workflow ID for debugging
   - Example: `Failed at step 3 (clone_database): Connection timeout [workflow: abc123]`

3. **Timeout Handling**
   - Show clear timeout message
   - Display last known operation
   - Suggest workflow:wait command to resume monitoring

4. **Interruption (Ctrl+C)**
   - Clean up progress bar before exit
   - Show workflow ID so user can monitor separately
   - Example: `Interrupted. Workflow abc123 is still running. Use 'workflow:wait site123 abc123' to monitor.`

### Backward Compatibility

**CLI Flags:**
- `--quiet` - Suppresses all progress indicators (existing behavior)
- `--format=json` - No progress bars, only structured output
- `--format=yaml` - No progress bars, only structured output

**Environment Variables:**
- `TERMINUS_QUIET=1` - Alternative to --quiet flag
- `CI=true` - Auto-detect CI environments and disable progress bars

**Output Streams:**
- Progress bars write to STDERR
- Actual command output writes to STDOUT
- Allows piping output while seeing progress

## Implementation Checklist

### Phase 1 (Recommended Starting Point)

- [ ] Add elapsed time tracking to `waitForWorkflow()`
- [ ] Display `CurrentOperation` in progress bar description
- [ ] Format elapsed time as human-readable duration
- [ ] Test with `backup:create` command
- [ ] Test in TTY and non-TTY environments
- [ ] Update documentation with examples

**Estimated effort:** 2-4 hours
**Files changed:** 1 (`internal/commands/workflow.go`)
**Lines of code:** ~30-50

### Phase 2

- [ ] Add step counter to display
- [ ] Calculate total steps from workflow Operations
- [ ] Enhance time formatting (e.g., "3m 45s" vs "3 minutes 45 seconds")
- [ ] Update `WaitOptions` callback signature
- [ ] Test with multiple workflow types
- [ ] Handle edge cases (workflows with no operations)

**Estimated effort:** 4-6 hours
**Files changed:** 2 (`pkg/api/workflows.go`, `internal/commands/workflow.go`)
**Lines of code:** ~80-120

### Phase 3

- [ ] Design progress indicator interfaces
- [ ] Create `pkg/output/progress.go` with abstractions
- [ ] Implement `WorkflowProgressIndicator`
- [ ] Write unit tests for progress indicators
- [ ] Refactor `waitForWorkflow()` to use abstraction
- [ ] Update all workflow-triggering commands
- [ ] Add integration tests
- [ ] Update documentation

**Estimated effort:** 8-12 hours
**Files changed:** ~10-15
**Lines of code:** ~300-500

### Phase 4 (Optional)

- [ ] Implement ETA calculation
- [ ] Add percentage completion
- [ ] Create operation history view
- [ ] Add interactive mode support
- [ ] Performance optimization
- [ ] Advanced testing scenarios

**Estimated effort:** 12-16 hours
**Files changed:** ~5-10
**Lines of code:** ~200-400

## Risks and Mitigations

### Risk 1: Terminal Compatibility Issues
**Probability:** Medium | **Impact:** Medium

Some terminals may not support ANSI escape codes or spinner animations.

**Mitigation:**
- Use progressbar/v3 which handles terminal detection
- Fall back to simple text updates on incompatible terminals
- Test on multiple terminal types (bash, zsh, fish, Windows PowerShell)

### Risk 2: Performance Impact
**Probability:** Low | **Impact:** Low

Frequent terminal updates might impact performance on slow connections.

**Mitigation:**
- Update progress only on workflow poll (every 3 seconds)
- Use efficient string formatting
- Avoid unnecessary allocations in hot path

### Risk 3: Workflow API Changes
**Probability:** Low | **Impact:** Medium

Pantheon API might change workflow structure.

**Mitigation:**
- Defensive coding (check for nil/empty values)
- Graceful degradation (fall back to basic spinner if fields missing)
- Version API client appropriately

### Risk 4: User Experience Regression
**Probability:** Low | **Impact:** High

Users might prefer simpler output or find updates distracting.

**Mitigation:**
- Keep `--quiet` flag for minimal output
- Respect CI environment detection
- Gather user feedback before Phase 3+ features
- Make progress indicators optional via config

## Success Metrics

### User Experience
- Reduced user uncertainty during long operations
- Fewer "is it stuck?" support questions
- Positive feedback on progress visibility

### Technical Metrics
- <100ms overhead per progress update
- Zero test failures introduced
- 100% backward compatibility with existing scripts
- All 15+ workflow commands display progress

### Code Quality
- >80% test coverage for progress code
- Zero linting errors
- Clean abstraction boundaries
- Reusable across all workflow operations

## Future Enhancements

### Parallel Operation Support
Some operations might run multiple workflows in parallel. Consider:
- Multi-line progress display
- Aggregate progress tracking
- Operation prioritization

### Alternative Progress Styles
Allow users to choose progress style via config:
- Spinner (current)
- Bar graph (for deterministic operations)
- Minimal (just status updates)
- Verbose (full operation log)

### JSON Progress Stream
For scripting and automation:
```json
{"status": "running", "operation": "backup_database", "step": 2, "elapsed": 83}
{"status": "running", "operation": "clone_files", "step": 3, "elapsed": 125}
{"status": "success", "duration": 332}
```

### Progress Persistence
Save progress to file for resumable operations:
```bash
terminus backup:create site.env --save-progress=/tmp/backup-progress.json
# Later, if interrupted:
terminus backup:resume --progress=/tmp/backup-progress.json
```

## References

### Libraries
- **progressbar/v3:** https://github.com/schollz/progressbar
  - Current version: v3.18.0
  - Features: Spinners, bars, custom themes, terminal detection

### Similar Tools
- **Terraform:** Shows resource operations with progress
- **Kubernetes kubectl:** Displays rollout progress
- **Docker:** Shows layer download progress
- **GitHub CLI (gh):** Progress for clones and uploads

### Pantheon API
- Workflow structure documentation
- Operation types and states
- Polling best practices

## Appendix

### Workflow Operation Types (Examples)

Common operation names seen in workflows:
- `backup_database`
- `backup_files`
- `clone_database`
- `clone_files`
- `clear_cache`
- `deploy_code`
- `enable_addon`
- `disable_addon`
- `create_environment`
- `delete_environment`
- `merge_code`

### Time Duration Formatting

Proposed format for elapsed time:
- Under 1 minute: "45s"
- 1-60 minutes: "3m 45s"
- Over 60 minutes: "1h 23m"

### Color Scheme (if supported)

- **Running:** Blue spinner
- **Success:** Green checkmark
- **Failure:** Red X
- **Warning:** Yellow exclamation
- **Info:** Cyan text

### Configuration File Support

Future: Allow users to configure progress in `~/.terminus/config.yml`:

```yaml
progress:
  enabled: true
  style: spinner  # spinner, bar, minimal, verbose
  show_elapsed_time: true
  show_step_counter: true
  color: auto  # auto, always, never
```

## Conclusion

Implementing rich progress indicators will significantly improve the user experience for long-running Pantheon operations. The phased approach allows for incremental delivery of value while managing implementation risk.

**Recommended starting point:** Phase 1 (Enhanced Current Operation Display)
- Quick to implement (2-4 hours)
- High user value
- Low risk
- Foundation for future enhancements

Once Phase 1 is validated, proceed to Phase 2 (Step Counter) and then Phase 3 (Abstraction Layer) for long-term maintainability.
