# Claude Development Instructions

This file contains instructions for Claude Code when working on this project.

## Pre-Flight Checks

Before making any changes or commits, ALWAYS run:

```bash
golangci-lint run --timeout=5m
```

If there are any linting errors, **STOP** and fix them before proceeding.

**Note on golangci-lint versions:**
- The project uses golangci-lint v2.6.2
- The CI pipeline uses v2.6.2, which is the authoritative version
- A SessionStart hook is configured to automatically install the correct version
- When in doubt, verify changes pass in CI rather than relying solely on local linting

## Development Workflow

### 1. After Every Code Change

Run the linter immediately:
```bash
golangci-lint run --timeout=5m
```

Fix all issues before continuing with further changes.

### 2. Before Every Commit

**MANDATORY CHECKS:**

1. **Run linter:**
   ```bash
   golangci-lint run --timeout=5m
   ```
   - Must return "0 issues"
   - Fix any issues before proceeding

2. **Run tests:**
   ```bash
   go test ./... -v
   ```
   - All tests must pass
   - No test failures allowed

3. **Build verification:**
   ```bash
   go build -o bin/terminus ./cmd/terminus
   ```
   - Build must succeed
   - No compilation errors

### 3. Common Linting Issues and Fixes

#### Unchecked Errors (errcheck)

**Wrong:**
```go
fmt.Println("hello")
resp.Body.Close()
```

**Correct:**
```go
_, _ = fmt.Println("hello")
defer func() { _ = resp.Body.Close() }()
```

#### Unused Parameters (revive)

**Wrong:**
```go
func handler(cmd *cobra.Command, args []string) error {
    // cmd and args not used
}
```

**Correct:**
```go
func handler(_ *cobra.Command, _ []string) error {
    // Explicitly mark as unused
}
```

#### Missing Package Comments (revive)

**Wrong:**
```go
package mypackage
```

**Correct:**
```go
// Package mypackage provides functionality for...
package mypackage
```

#### Response Body Not Closed (bodyclose)

**Wrong:**
```go
resp, err := client.Get(ctx, "/path")
// body not closed
```

**Correct:**
```go
resp, err := client.Get(ctx, "/path")
defer func() { _ = resp.Body.Close() }()
```

Or use DecodeResponse which handles closing:
```go
resp, err := client.Get(ctx, "/path")
if err := DecodeResponse(resp, &target); err != nil {
    return err
}
```

#### Security Issues (gosec)

For intentional file operations, add nolint comments:
```go
data, err := os.ReadFile(path) //nolint:gosec // User-specified config file
```

### 4. Linter Configuration

The project uses `.golangci.yml` with these enabled linters:
- errcheck - unchecked errors
- govet - suspicious constructs
- ineffassign - ineffectual assignments
- staticcheck - static analysis
- unused - unused code
- misspell - spelling
- revive - code quality
- gocritic - comprehensive checks
- gocyclo - complexity
- gosec - security
- unconvert - unnecessary conversions
- prealloc - slice preallocation
- bodyclose - HTTP body closing
- nolintlint - nolint directives

### 5. Quick Fix Commands

If you encounter linting errors, here are some quick fixes:

```bash
# Format code
go fmt ./...

# Tidy modules
go mod tidy

# Run specific linter
golangci-lint run --disable-all -E errcheck
```

## Important Reminders

- ✅ **ALWAYS** run linter after changes
- ✅ **NEVER** commit with linting errors
- ✅ **ALWAYS** run tests before committing
- ✅ Fix all issues immediately, don't accumulate them
- ✅ Use `//nolint` sparingly and only when justified

## Continuous Integration

The CI pipeline will fail if:
- Linting errors exist
- Tests fail
- Build fails

Save time by catching these locally with the checks above.

---

**TL;DR:** Run `golangci-lint run --timeout=5m` after every change and before every commit. Fix all issues immediately.
