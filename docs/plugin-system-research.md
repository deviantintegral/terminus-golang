# Plugin System Research for Terminus Go

## Executive Summary

This document evaluates plugin system approaches for extending Terminus Go with third-party commands. Based on research into Go ecosystem best practices and the specific requirements (ease of development, native command integration), we recommend a **hybrid approach** combining kubectl-style PATH discovery with an optional directory-based plugin system.

## Requirements

1. Third-party extensibility without modifying core Terminus code
2. Easy for Go developers to implement
3. Plugins appear as native `terminus` commands (not separate programs)
4. Cross-platform support (Linux, macOS, Windows)

## Plugin System Approaches Evaluated

### 1. Go Native Plugin Package

**How it works:** Compile plugins as `.so` shared libraries, load at runtime using Go's `plugin` package.

**Pros:**
- Fast function calls (no RPC overhead)
- Can export Cobra commands directly
- Tight integration with host application

**Cons:**
- **Cross-platform issues:** Only works on Linux, macOS, FreeBSD; Windows support is experimental
- **Version matching required:** All shared dependencies must be exact same versions
- **Build complexity:** Plugins must be compiled with same Go version and flags
- **No unloading:** Cannot unload plugins once loaded

**Verdict:** ❌ **Not recommended** - Too fragile for third-party development; the strict version matching requirements make it impractical for independent plugin development.

### 2. HashiCorp go-plugin (RPC-based)

**How it works:** Plugins are separate binaries that communicate with the host via gRPC/net-rpc.

**Pros:**
- Battle-tested in Terraform, Vault, Nomad, Packer
- Plugin crashes don't crash host
- Cross-platform support
- Plugins can be written in any language
- Protocol versioning for compatibility

**Cons:**
- RPC overhead (~30-50μs per call)
- More complex to implement
- Requires defining gRPC/protobuf interfaces
- Plugins are processes, not direct Cobra commands

**Verdict:** ⚠️ **Overkill for this use case** - Better suited for plugins that need to maintain state or perform complex operations. The RPC overhead and complexity aren't justified for simple command extensions.

### 3. External Binary Discovery (kubectl-style) ⭐ RECOMMENDED

**How it works:** Discover executable files in PATH (or plugin directory) with a naming convention like `terminus-*`.

**Pros:**
- **Simplest for developers:** Write any executable, get it in PATH, done
- **True command integration:** `terminus foo` automatically invokes `terminus-foo`
- **Cross-platform:** Works everywhere
- **Language agnostic:** Plugins can be Go, Python, shell scripts, etc.
- **No shared dependencies:** Each plugin is independent
- **Proven pattern:** Used by kubectl, git, docker

**Cons:**
- Subprocess overhead per invocation
- Plugins can't directly access Terminus internals (API client, config)
- Need to standardize how plugins receive authentication/config

**Verdict:** ✅ **Recommended** - Best balance of simplicity, reliability, and developer experience.

### 4. Yaegi (Go Interpreter)

**How it works:** Interpret Go source files at runtime.

**Pros:**
- Plugins as simple `.go` files
- Sandboxing capabilities
- No compilation needed

**Cons:**
- Performance overhead
- Not all Go features supported
- Limited to Go language
- Still experimental for production use

**Verdict:** ❌ **Not recommended** - Performance concerns and limited compatibility make this unsuitable for a production CLI tool.

---

## Recommended Architecture

### Primary: External Binary Discovery

Adopt the kubectl-style plugin pattern with enhancements for Terminus-specific needs.

#### Naming Convention

```
terminus-<plugin-name>[-<subcommand>]
```

Examples:
- `terminus-backup-scheduler` → `terminus backup-scheduler`
- `terminus-site-audit` → `terminus site-audit`
- `terminus-wpms-network` → `terminus wpms-network`

#### Discovery Locations

1. **PATH** - Standard system PATH (highest priority)
2. **Plugin Directory** - `~/.terminus/plugins-3.x/` (already configured in config.go:59)
3. **Project-local** - `./.terminus/plugins/` (optional, for project-specific plugins)

#### How It Works

1. User runs `terminus foo bar`
2. Terminus checks if `foo` is a built-in command
3. If not found, searches for `terminus-foo` executable:
   - First in `~/.terminus/plugins-3.x/`
   - Then in PATH
4. If found, executes `terminus-foo bar` with environment variables for context
5. If not found, shows "unknown command" error

#### Plugin Contract

Plugins receive context via environment variables:

```bash
TERMINUS_PLUGIN=1                      # Indicates running as plugin
TERMINUS_SESSION_TOKEN=<token>         # Authentication token
TERMINUS_HOST=terminus.pantheon.io     # API host
TERMINUS_FORMAT=table                  # Output format preference
TERMINUS_CONFIG_DIR=~/.terminus        # Config directory
```

Plugins can use the `pkg/api` package directly by importing it as a Go library.

#### Plugin Metadata (Optional)

Plugins can implement a `--terminus-plugin-info` flag to return JSON metadata:

```json
{
  "name": "backup-scheduler",
  "version": "1.0.0",
  "description": "Schedule automated backups for sites",
  "author": "Third Party Inc",
  "commands": [
    {
      "name": "schedule",
      "description": "Create a backup schedule"
    }
  ]
}
```

This enables `terminus plugin list` to show all available plugins with descriptions.

---

## Implementation Plan

### Phase 1: Core Plugin Discovery

1. **Plugin loader in root command** (`internal/commands/root.go`)
   - Add `PersistentPreRunE` logic to check for plugins
   - Implement binary discovery in plugin directory and PATH
   - Execute plugins as subprocesses with environment context

2. **Plugin list command**
   ```
   terminus plugin list
   ```
   - Scan plugin locations
   - Display found plugins with optional metadata

3. **Plugin path command**
   ```
   terminus plugin path
   ```
   - Show plugin search paths for developer reference

### Phase 2: Developer Experience

1. **Plugin template/scaffold**
   ```
   terminus plugin init my-plugin
   ```
   - Generate Go project structure
   - Include Cobra boilerplate
   - Pre-configured to import `pkg/api`

2. **Documentation**
   - Plugin development guide
   - API client usage examples
   - Best practices for output formatting

3. **Testing utilities**
   - Mock API responses
   - Test harness for plugin commands

### Phase 3: Distribution (Future)

1. **Plugin registry** (optional)
   - Central repository of community plugins
   - Installation command: `terminus plugin install <name>`

2. **Version management**
   - Track installed plugin versions
   - Update notifications

---

## Example Plugin Implementation

### Simple Go Plugin

```go
// cmd/terminus-site-audit/main.go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/pantheon-systems/terminus-go/pkg/api"
    "github.com/spf13/cobra"
)

func main() {
    rootCmd := &cobra.Command{
        Use:   "terminus-site-audit",
        Short: "Audit site configurations",
        RunE:  runAudit,
    }

    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}

func runAudit(cmd *cobra.Command, args []string) error {
    // Get auth from environment (set by Terminus)
    token := os.Getenv("TERMINUS_SESSION_TOKEN")
    if token == "" {
        return fmt.Errorf("not authenticated - run 'terminus auth login' first")
    }

    // Create API client
    client := api.NewClient(
        api.WithToken(token),
    )

    // Use Terminus API
    sitesService := api.NewSitesService(client)
    sites, err := sitesService.List(context.Background())
    if err != nil {
        return err
    }

    // Perform audit...
    for _, site := range sites {
        fmt.Printf("Auditing site: %s\n", site.Name)
        // ... audit logic
    }

    return nil
}
```

### Build and Install

```bash
# Build the plugin
go build -o terminus-site-audit ./cmd/terminus-site-audit

# Install to plugin directory
cp terminus-site-audit ~/.terminus/plugins-3.x/

# Or add to PATH
mv terminus-site-audit /usr/local/bin/

# Use it
terminus site-audit
```

---

## Comparison with PHP Terminus

| Feature | PHP Terminus | Proposed Go Approach |
|---------|-------------|---------------------|
| Plugin format | PHP classes with annotations | External binaries |
| Discovery | Composer autoloading | PATH/directory scanning |
| API access | Direct object access | Via `pkg/api` import or env vars |
| Installation | Composer require | Copy binary or go install |
| Language | PHP only | Any language |

The Go approach is simpler and more flexible, though it trades some tight integration for reliability and ease of development.

---

## Security Considerations

1. **Plugin isolation:** Plugins run as separate processes, limiting damage from bugs
2. **Permission model:** Plugins inherit user's filesystem permissions
3. **Token handling:** Session tokens passed via environment (not command line for security)
4. **Code signing (future):** Consider plugin verification for enterprise environments

---

## Alternatives Considered but Rejected

### Lua/JavaScript Scripting

Could embed a scripting engine for lightweight plugins. Rejected because:
- Adds large dependency
- Different language than host
- Limited access to Go ecosystem

### WASM Plugins

WebAssembly could provide sandboxed plugins. Rejected because:
- Immature Go→WASM tooling
- Performance overhead
- Complex for simple use cases

---

## Conclusion

The **external binary discovery pattern** (kubectl-style) provides the best balance of:

- **Developer simplicity:** Write a Go program, put it in the right place, done
- **Reliability:** No shared library versioning issues
- **Flexibility:** Plugins can be in any language
- **User experience:** Commands appear native to Terminus

The existing `~/.terminus/plugins-3.x/` directory in the config provides a natural home for plugins, and the Cobra-based architecture makes it straightforward to implement command delegation.

Estimated implementation effort: **2-3 days** for Phase 1 (core discovery), **1 week** for full Phase 1+2.

---

## Next Steps

1. Review and approve this architecture
2. Create detailed technical specification
3. Implement Phase 1 plugin discovery
4. Write developer documentation
5. Create example plugins

---

## References

- [kubectl Plugin System](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/)
- [HashiCorp go-plugin](https://github.com/hashicorp/go-plugin)
- [Cobra Issue #691: Plugin Support](https://github.com/spf13/cobra/issues/691)
- [Go Plugin Package](https://pkg.go.dev/plugin)
