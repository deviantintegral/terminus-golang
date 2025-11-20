# Goroutine Implementation Opportunities

## Executive Summary

This document identifies opportunities to introduce concurrent programming patterns using goroutines in the Terminus Go CLI. The analysis reveals several high-impact areas where parallelization can significantly improve performance, particularly around API pagination, workflow monitoring, and batch operations.

**Key Findings:**
- Current implementation is entirely sequential with minimal concurrency
- Primary bottleneck: Network I/O from sequential API calls
- Estimated performance improvements: 40-70% reduction in execution time for multi-item operations
- All code already uses `context.Context`, making goroutine integration straightforward

---

## Project Context

**Terminus Go** is a CLI tool for managing Pantheon sites. It performs numerous I/O-bound operations:
- Fetching paginated lists of sites, organizations, and resources
- Monitoring asynchronous workflows
- Managing backups, environments, and domains
- Making multiple API calls for batch operations

The codebase is well-structured with:
- Clean service layer abstraction (`pkg/api/`)
- Proper context usage throughout
- HTTP client with built-in retry logic
- No existing goroutine patterns (except basic ticker usage in workflow polling)

---

## High Priority Opportunities

### 1. Parallel Pagination (HIGHEST IMPACT)

**Current State:** `pkg/api/client.go:329-379` - `GetPaged()`

```go
// Current: Sequential page fetching
for {
    resp, err := c.Get(ctx, path)  // Network call
    // ... process page
    page++  // Next page sequentially
}
```

**Problem:**
- Large datasets require multiple sequential API calls
- 5 pages × 200ms latency = 1 second total
- User must wait for each page to complete before next begins

**Opportunity:**
Implement parallel page fetching with worker pool pattern.

**Implementation Approach:**
```go
func (c *Client) GetPagedParallel(ctx context.Context, basePath string) ([]json.RawMessage, error) {
    // 1. Fetch first page to determine total pages (if API provides count)
    // 2. Launch bounded worker pool (e.g., 3-5 workers to respect rate limits)
    // 3. Distribute pages across workers
    // 4. Collect results and maintain order
    // 5. Handle errors with error group pattern
}
```

**Technical Considerations:**
- Use `errgroup` package for coordinated error handling
- Implement bounded worker pool (3-5 concurrent requests) to avoid rate limiting
- Preserve result ordering (important for consistent table output)
- Add configuration option to control concurrency level
- Maintain backward compatibility with sequential `GetPaged()`

**Affected Services:**
- `Sites.List()` - Site listings
- `Organizations.List()` - Organization listings
- `SitesService.GetTeam()` - Team member listings
- `OrganizationsService.ListMembers()` - Member listings
- `OrganizationsService.ListUpstreams()` - Upstream listings

**Estimated Impact:**
- 50-70% reduction in time for large paginated results
- Most noticeable with 100+ items requiring 2+ pages

**Complexity:** Medium (2-3 days)

---

### 2. Concurrent Workflow Monitoring (HIGH IMPACT)

**Current State:** `pkg/api/workflows.go:87-121` - `Wait()`

```go
// Current: Single workflow polling
func (s *WorkflowsService) Wait(ctx context.Context, siteID, workflowID string, opts *WaitOptions) (*models.Workflow, error) {
    ticker := time.NewTicker(opts.PollInterval)
    for {
        workflow, err := s.Get(timeoutCtx, siteID, workflowID)
        // ... check if finished
    }
}
```

**Problem:**
- Can only monitor one workflow at a time
- Batch operations (deploy to multiple envs) must wait sequentially
- No way to track parallel operations in progress

**Opportunity:**
Add `WaitMultiple()` method to monitor multiple workflows concurrently.

**Implementation Approach:**
```go
type MultiWaitOptions struct {
    Workflows     []WorkflowIdentifier  // siteID + workflowID pairs
    PollInterval  time.Duration
    Timeout       time.Duration
    OnProgress    func(workflowID string, workflow *models.Workflow)
}

func (s *WorkflowsService) WaitMultiple(ctx context.Context, opts *MultiWaitOptions) (map[string]*models.Workflow, error) {
    // 1. Launch goroutine for each workflow
    // 2. Shared ticker to coordinate polling
    // 3. Collect results with channel
    // 4. Return when all complete or first error
    // 5. Support timeout for entire batch
}
```

**Technical Considerations:**
- Use single ticker shared across all goroutines (more efficient than N tickers)
- Implement with `sync.WaitGroup` or `errgroup` for coordination
- Support partial success mode (return completed workflows even if some fail)
- Add progress reporting for each workflow individually
- Respect global timeout across all workflows

**Use Cases:**
- Deploy to multiple environments simultaneously
- Monitor multiple backup creations
- Track batch domain operations
- Watch parallel multidev operations

**Estimated Impact:**
- Linear time reduction: N workflows in ~1x time instead of Nx time
- For 3 deployments taking 2 minutes each: 2 minutes vs 6 minutes

**Complexity:** Medium (2-3 days)

---

### 3. Batch Backup Operations (MEDIUM-HIGH IMPACT)

**Current State:** `pkg/api/backups.go:89-111` - `CreateElement()`

**Problem:**
- Creating backups for database, files, and code requires 3 sequential API calls
- Downloading multiple backup elements is sequential
- No way to create backups across multiple environments in parallel

**Opportunity:**
Implement parallel backup creation and download operations.

**Implementation Approach:**

```go
// Create multiple backup elements in parallel
func (s *BackupsService) CreateMultipleElements(ctx context.Context, siteID, envID string, elements []string) (map[string]*models.Workflow, error) {
    // Launch goroutine per element (code, database, files)
    // Collect workflow IDs
    // Return map of element -> workflow
}

// Download multiple elements concurrently
func (s *BackupsService) DownloadMultiple(ctx context.Context, siteID, envID, backupID string, elements map[string]string) error {
    // elements: map[element]outputPath
    // Download each element in parallel
    // Use bounded concurrency (2-3 downloads at once)
}
```

**Technical Considerations:**
- Backup creation API calls are independent and can run concurrently
- Downloads should be bounded (2-3 concurrent) to avoid overwhelming network
- Consider implementing progress reporting for large downloads
- Add checksum verification for downloaded files

**Use Cases:**
- `backup:create` with `--elements=all` flag
- Bulk backup downloads
- Creating backups across dev/test/live environments simultaneously

**Estimated Impact:**
- 3x speedup for full backup creation (3 elements)
- 2-3x speedup for multi-element downloads

**Complexity:** Low-Medium (1-2 days)

---

## Medium Priority Opportunities

### 4. Concurrent Organization Detail Fetching

**Current State:** `internal/commands/org.go:59-81` - `runOrgList()`

**Problem:**
When listing organizations, the initial call returns summaries. Getting detailed info requires additional calls:

```go
orgs, err := orgsService.List(ctx, userID)  // Summary list
// If user wants details, must fetch each individually:
for _, org := range orgs {
    details, err := orgsService.Get(ctx, org.ID)  // Sequential
}
```

**Opportunity:**
Add `GetMultiple()` method or fetch details concurrently in commands.

**Implementation Approach:**

```go
func (s *OrganizationsService) GetMultiple(ctx context.Context, orgIDs []string) (map[string]*models.Organization, error) {
    results := make(map[string]*models.Organization)
    var mu sync.Mutex

    g, ctx := errgroup.WithContext(ctx)
    g.SetLimit(5)  // Bounded concurrency

    for _, orgID := range orgIDs {
        orgID := orgID  // Capture for closure
        g.Go(func() error {
            org, err := s.Get(ctx, orgID)
            if err != nil {
                return err
            }
            mu.Lock()
            results[orgID] = org
            mu.Unlock()
            return nil
        })
    }

    if err := g.Wait(); err != nil {
        return nil, err
    }
    return results, nil
}
```

**Use Cases:**
- `org:list` with `--detailed` flag
- Getting multiple organization details for reporting

**Estimated Impact:**
- For users with 10 organizations: 500ms vs 2-3 seconds
- Scales linearly with organization count

**Complexity:** Low (1 day)

---

### 5. Parallel Site Operations

**Opportunity:**
Many site-related operations could be parallelized when operating on multiple sites.

**Potential New Commands:**
```bash
# Batch site info
terminus sites:info site1,site2,site3

# Parallel deployments
terminus env:deploy site1.dev,site2.dev,site3.dev

# Batch cache clearing
terminus env:clear-cache site1.live,site2.live
```

**Implementation Approach:**
- Add batch variants of existing commands
- Use worker pool with bounded concurrency
- Provide progress reporting for each site
- Support fail-fast or continue-on-error modes

**Estimated Impact:**
- N sites processed in ~1x time instead of Nx time
- Particularly valuable for agencies managing many sites

**Complexity:** Medium (2-3 days per command type)

---

### 6. Parallel Domain Operations

**Current State:** `internal/commands/domain.go`

**Opportunities:**
- Add multiple domains to a site concurrently
- Remove multiple domains in parallel
- Check DNS/SSL status for multiple domains simultaneously

**Implementation:**
```go
func (s *DomainsService) AddMultiple(ctx context.Context, siteID, envID string, domains []string) (map[string]error, error) {
    // Add domains in parallel
    // Return map showing success/failure for each
}
```

**Use Cases:**
- Bulk domain management
- Migration scenarios with many domains

**Estimated Impact:**
- 10 domains: 30 seconds vs 5 minutes
- Moderate impact (domain operations less frequent than other operations)

**Complexity:** Low (1 day)

---

## Low Priority / Future Enhancements

### 7. Output Formatting Parallelization

**Current State:** `pkg/output/output.go`

**Opportunity:**
For very large datasets (100+ items), parallel formatting could provide minor improvements.

**Rationale for Low Priority:**
- Most CLI operations return <50 items
- Formatting is CPU-bound, not I/O-bound
- Potential improvements are marginal (<10%)
- Added complexity may not justify minimal gains

**Recommendation:** Defer until profiling shows formatting as bottleneck.

---

### 8. Advanced Polling Strategies

**Opportunities:**
- Exponential backoff for long-running workflows (poll frequently initially, then slow down)
- Request coalescing (if multiple commands watch same workflow, share 1 HTTP request)
- Intelligent polling based on workflow type

**Example:**
```go
// Adaptive polling
func (s *WorkflowsService) WaitAdaptive(ctx context.Context, siteID, workflowID string) (*models.Workflow, error) {
    intervals := []time.Duration{1*time.Second, 2*time.Second, 5*time.Second, 10*time.Second}
    // Start with frequent polling, back off over time
}
```

**Estimated Impact:**
- Reduced API calls for long-running operations
- Faster feedback for quick operations
- Lower load on API servers

**Complexity:** Medium (2 days)

---

## Implementation Guidelines

### General Patterns

#### 1. Error Group Pattern (Recommended)
```go
import "golang.org/x/sync/errgroup"

func ParallelOperation(ctx context.Context, items []string) error {
    g, ctx := errgroup.WithContext(ctx)
    g.SetLimit(5)  // Max 5 concurrent goroutines

    for _, item := range items {
        item := item  // Capture for closure
        g.Go(func() error {
            return processItem(ctx, item)
        })
    }

    return g.Wait()  // Wait for all and return first error
}
```

**Advantages:**
- Automatic error propagation
- Context cancellation on first error
- Simple, clean API
- Built-in concurrency limiting

#### 2. Worker Pool Pattern
```go
func WorkerPoolOperation(ctx context.Context, items []string) error {
    const numWorkers = 5
    jobs := make(chan string, len(items))
    errors := make(chan error, len(items))

    // Start workers
    var wg sync.WaitGroup
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for item := range jobs {
                if err := processItem(ctx, item); err != nil {
                    errors <- err
                    return
                }
            }
        }()
    }

    // Send jobs
    for _, item := range items {
        jobs <- item
    }
    close(jobs)

    // Wait and collect errors
    wg.Wait()
    close(errors)

    for err := range errors {
        if err != nil {
            return err
        }
    }
    return nil
}
```

**Use when:**
- Need fine-grained control over worker behavior
- Want to limit concurrent resource usage (e.g., open files)
- Need complex error handling strategies

#### 3. Channel-Based Coordination
```go
func ChannelCoordination(ctx context.Context, items []string) ([]Result, error) {
    results := make(chan Result, len(items))

    for _, item := range items {
        go func(item string) {
            result, err := processItem(ctx, item)
            results <- Result{Value: result, Error: err}
        }(item)
    }

    // Collect results
    var collected []Result
    for i := 0; i < len(items); i++ {
        collected = append(collected, <-results)
    }
    return collected, nil
}
```

**Use when:**
- Need to collect results from all goroutines regardless of errors
- Want partial success handling
- Need progress reporting

---

### Critical Considerations

#### 1. Rate Limiting

**Problem:** Pantheon API likely has rate limits. Unbounded concurrency could trigger 429 responses.

**Solution:**
```go
// Global rate limiter
import "golang.org/x/time/rate"

type Client struct {
    // ... existing fields
    rateLimiter *rate.Limiter
}

// In NewClient
c.rateLimiter = rate.NewLimiter(rate.Limit(10), 20)  // 10 req/s, burst of 20

// Before each request
if err := c.rateLimiter.Wait(ctx); err != nil {
    return nil, err
}
```

**Recommendations:**
- Start conservative: 3-5 concurrent requests
- Make concurrency configurable via environment variable
- Monitor for 429 responses and back off automatically
- Add metrics/logging to track concurrent request counts

#### 2. Memory Management

**Problem:** Fetching many pages concurrently could consume significant memory.

**Solution:**
- Use bounded worker pools (don't launch N goroutines for N pages)
- Stream results when possible rather than buffering all in memory
- Consider pagination limits for very large datasets

#### 3. Error Handling

**Strategies:**

1. **Fail-Fast (Default):**
   - First error cancels all other operations
   - Use `errgroup` with context cancellation
   - Best for critical operations

2. **Fail-Soft (Optional):**
   - Collect all errors, return partial results
   - Useful for batch operations where partial success is acceptable
   - Provide detailed error report showing which items failed

```go
type BatchResult struct {
    Successes map[string]Result
    Failures  map[string]error
}
```

3. **Retry with Backoff:**
   - Leverage existing retry logic in `client.doWithRetry()`
   - Consider per-item retry counts in batch operations

#### 4. Context Management

**Best Practices:**
- Always pass and respect `context.Context`
- Use `context.WithTimeout()` for operations with deadlines
- Cancel parent context to stop all child goroutines
- Don't create goroutines that might outlive parent context

```go
// Good: Context propagation
func ParallelFetch(ctx context.Context, items []string) error {
    g, ctx := errgroup.WithContext(ctx)  // Derives new context
    // Children inherit cancellation
}
```

#### 5. Testing Concurrent Code

**Strategies:**
- Use `go test -race` to detect data races
- Test with various concurrency limits (1, 5, 100)
- Simulate API errors and timeouts
- Test context cancellation mid-operation
- Verify result ordering where required

**Example Test:**
```go
func TestGetPagedParallel(t *testing.T) {
    t.Parallel()  // Run in parallel with other tests

    // Test with race detector
    // go test -race ./pkg/api

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    results, err := client.GetPagedParallel(ctx, "/test")
    require.NoError(t, err)
    assert.Len(t, results, expectedCount)
}
```

#### 6. Observability

**Recommendations:**
- Add debug logging for goroutine launch/completion
- Track concurrent operation counts
- Log timing improvements (sequential vs parallel)
- Consider adding OpenTelemetry spans for distributed tracing

```go
if c.logger != nil {
    c.logger.Debug("Starting parallel fetch: %d pages with %d workers", totalPages, numWorkers)
}
```

---

## Migration Strategy

### Phase 1: Foundation (Week 1)
1. Implement `GetPagedParallel()` with feature flag
2. Add comprehensive tests with race detector
3. Test with real API under various loads
4. Add rate limiting infrastructure
5. Create benchmarks comparing sequential vs parallel

### Phase 2: Workflow Enhancements (Week 2)
1. Implement `WaitMultiple()` for workflows
2. Add progress reporting for concurrent operations
3. Test with multiple simultaneous deployments
4. Update documentation with new capabilities

### Phase 3: Backup Operations (Week 3)
1. Implement parallel backup creation
2. Add concurrent download support
3. Test with large backup files
4. Add progress bars for downloads

### Phase 4: Batch Commands (Week 4)
1. Add batch variants of common commands
2. Implement error reporting for partial failures
3. Create user documentation
4. Add example scripts for common batch operations

### Phase 5: Optimization & Polish (Week 5)
1. Profile and optimize based on real-world usage
2. Tune concurrency limits based on API feedback
3. Add advanced features (adaptive polling, request coalescing)
4. Performance testing and benchmarking

---

## Configuration

### Environment Variables
```bash
# Control concurrency levels
export TERMINUS_MAX_CONCURRENT_REQUESTS=5

# Enable parallel pagination
export TERMINUS_PARALLEL_PAGINATION=true

# Workflow polling concurrency
export TERMINUS_MAX_CONCURRENT_WORKFLOWS=10
```

### Config File
```yaml
# ~/.terminus/config.yml
performance:
  max_concurrent_requests: 5
  parallel_pagination: true
  max_concurrent_workflows: 10

  # Rate limiting
  requests_per_second: 10
  burst_size: 20
```

---

## Success Metrics

### Performance Targets
- **Pagination:** 50-70% reduction for 100+ item lists
- **Workflows:** Linear scaling (3 workflows = 1x time, not 3x)
- **Batch Operations:** 40-60% improvement for 5+ items

### Quality Metrics
- Zero data races (verified with `-race` detector)
- No increase in API error rates
- No regression in existing functionality
- Test coverage maintained >80%

### User Experience
- Clear progress indication for parallel operations
- Helpful error messages when operations fail
- No breaking changes to existing commands
- Opt-in for new parallel features initially

---

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|-----------|
| API rate limiting triggered | High | Conservative concurrency limits (3-5), configurable, automatic backoff |
| Data races in shared state | High | Thorough testing with `-race`, use channels/mutexes properly, code review |
| Increased memory usage | Medium | Bounded worker pools, streaming where possible, monitoring |
| Increased complexity | Medium | Clear documentation, comprehensive tests, gradual rollout |
| Ordering issues in output | Low | Maintain result ordering, document any ordering changes |
| Context leaks | Medium | Careful context management, linting with `govet`, timeout tests |

---

## Open Questions

1. **API Rate Limits:** What are Pantheon's actual rate limits? Need to test and document.
2. **User Preferences:** Should parallel operations be opt-in or opt-out initially?
3. **Backward Compatibility:** Keep both sequential and parallel methods, or replace entirely?
4. **Error Reporting:** How to best display partial failures in batch operations?
5. **Progress Reporting:** ASCII progress bars, or simple status messages?

---

## Conclusion

The Terminus Go codebase presents excellent opportunities for goroutine-based parallelization. The code is already well-structured with proper context usage and clean service abstractions, making concurrent enhancements relatively straightforward.

**Recommended Starting Point:**
Begin with parallel pagination (`GetPagedParallel()`) as it:
- Has the broadest impact across many commands
- Is relatively self-contained
- Provides measurable performance improvements
- Serves as a pattern for other enhancements

**Expected Overall Impact:**
- 40-70% improvement in multi-item operations
- Better user experience for batch workflows
- Positions codebase for future scalability
- Maintains code quality and maintainability

**Next Steps:**
1. Review and approve this plan
2. Set up performance benchmarking infrastructure
3. Implement Phase 1 (pagination) with feature flag
4. Gather feedback and metrics
5. Iterate and expand to other areas

---

*Document Version: 1.0*
*Created: 2025-11-20*
*Author: Claude Code Analysis*
