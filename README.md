# Terminus Go

[![CI](https://github.com/deviantintegral/terminus-golang/actions/workflows/ci.yml/badge.svg)](https://github.com/deviantintegral/terminus-golang/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/deviantintegral/terminus-golang)](https://github.com/deviantintegral/terminus-golang/releases/latest)

A command-line interface for managing Pantheon sites, written in Go. This is a complete rewrite of the [Pantheon Terminus](https://github.com/pantheon-systems/terminus) CLI tool.

## Features

- **Command-line Interface**: Full-featured CLI for managing Pantheon sites, environments, workflows, and more
- **Go API Package**: Standalone Go package for programmatic access to the Pantheon API
- **Multiple Output Formats**: Table, JSON, YAML, CSV, and list formats
- **Comprehensive Commands**: Support for all major Pantheon operations
- **Session Management**: Secure token storage and session handling
- **Workflow Management**: Monitor and wait for asynchronous operations

## Installation

### From Source

```bash
go install github.com/pantheon-systems/terminus-go/cmd/terminus@latest
```

### Build Manually

```bash
git clone https://github.com/pantheon-systems/terminus-go.git
cd terminus-go
go build -o terminus ./cmd/terminus
```

## Quick Start

### Authentication

First, log in with your machine token:

```bash
terminus auth:login --machine-token=YOUR_TOKEN --email=your@email.com
```

### List Your Sites

```bash
terminus site list
```

### Get Site Information

```bash
terminus site info my-site
```

### List Environments

```bash
terminus env list my-site
```

### Deploy to Test

```bash
terminus env deploy my-site.test
```

### Create a Backup

```bash
terminus backup create my-site.live
```

## Available Commands

### Authentication
- `auth:login` - Log in to Pantheon
- `auth:logout` - Log out of Pantheon
- `auth:whoami` - Show current user information

### Site Management
- `site list` - List all sites
- `site info <site>` - Show site information
- `site create <name>` - Create a new site
- `site delete <site>` - Delete a site
- `site team list <site>` - List team members

### Environment Management
- `env list <site>` - List environments
- `env info <site>.<env>` - Show environment information
- `env clear-cache <site>.<env>` - Clear environment cache
- `env deploy <site>.<env>` - Deploy code to an environment
- `env clone-content <site>.<env>` - Clone database/files between environments
- `env commit <site>.<env>` - Commit changes in SFTP mode
- `env wipe <site>.<env>` - Wipe environment content
- `env connection set <site>.<env> <mode>` - Set connection mode (git/sftp)

### Workflow Management
- `workflow list <site>` - List workflows
- `workflow info <site> <workflow-id>` - Show workflow information
- `workflow wait <site> <workflow-id>` - Wait for a workflow to complete
- `workflow watch <site> <workflow-id>` - Watch a workflow with live updates

### Backup Management
- `backup list <site>.<env>` - List backups
- `backup create <site>.<env>` - Create a backup
- `backup get <site>.<env>` - Download a backup
- `backup restore <site>.<env>` - Restore from a backup
- `backup automatic info <site>.<env>` - Show backup schedule
- `backup automatic enable <site>.<env>` - Enable automatic backups
- `backup automatic disable <site>.<env>` - Disable automatic backups

### Organization Management
- `org list` - List organizations
- `org info <org>` - Show organization information
- `org people list <org>` - List organization members
- `org site list <org>` - List organization sites
- `org upstreams list <org>` - List organization upstreams

### Domain Management
- `domain list <site>.<env>` - List domains
- `domain add <site>.<env> <domain>` - Add a domain
- `domain remove <site>.<env> <domain>` - Remove a domain
- `domain dns <site>.<env> <domain>` - Show DNS recommendations

### Multidev Management
- `multidev create <site>.<multidev>` - Create a multidev environment
- `multidev delete <site>.<multidev>` - Delete a multidev environment
- `multidev merge-to-dev <site>.<multidev>` - Merge multidev to dev
- `multidev merge-from-dev <site>.<multidev>` - Merge dev into multidev

## Global Flags

- `--format` - Output format (table, json, yaml, csv, list)
- `--fields` - Comma-separated list of fields to display
- `--yes, -y` - Answer yes to all prompts
- `--quiet, -q` - Suppress output
- `--verbose, -v` - Verbose output

## Output Formats

### Table (Default)
```bash
terminus site list
```

### JSON
```bash
terminus site list --format=json
```

### YAML
```bash
terminus site info my-site --format=yaml
```

### CSV
```bash
terminus site list --format=csv
```

### Field Filtering
```bash
terminus site list --fields=name,id,framework
```

## Configuration

Terminus Go supports multiple configuration sources (in priority order):

1. Environment variables (`TERMINUS_*`)
2. `.env` file in current directory
3. User config file (`~/.terminus/config.yml`)
4. Default values

### Configuration File

Create `~/.terminus/config.yml`:

```yaml
TERMINUS_HOST: terminus.pantheon.io
TERMINUS_PORT: 443
TERMINUS_PROTOCOL: https
TERMINUS_TIMEOUT: 86400
```

### Environment Variables

```bash
export TERMINUS_HOST=terminus.pantheon.io
export TERMINUS_CACHE_DIR=~/.terminus/cache
```

## Using as a Go Package

Terminus Go can be used as a library in your Go applications:

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/pantheon-systems/terminus-go/pkg/api"
    "github.com/pantheon-systems/terminus-go/pkg/session"
)

func main() {
    // Create API client
    client := api.NewClient(
        api.WithToken("your-session-token"),
    )

    // Create services
    sitesService := api.NewSitesService(client)

    // List sites
    ctx := context.Background()
    sites, err := sitesService.List(ctx)
    if err != nil {
        log.Fatal(err)
    }

    for _, site := range sites {
        fmt.Printf("Site: %s (%s)\n", site.Name, site.ID)
    }
}
```

### API Services

The following services are available:

- `AuthService` - Authentication operations
- `SitesService` - Site management
- `EnvironmentsService` - Environment operations
- `WorkflowsService` - Workflow monitoring
- `BackupsService` - Backup management
- `OrganizationsService` - Organization management
- `DomainsService` - Domain management
- `MultidevService` - Multidev operations

## Development

### Prerequisites

- Go 1.24 or higher
- golangci-lint
- pre-commit and golanci-linter (optional)

### Setup

```bash
# Clone repository
git clone https://github.com/pantheon-systems/terminus-go.git
cd terminus-go

# Install dependencies
go mod download

# Build
go build -o bin/terminus ./cmd/terminus

# Run tests
go test ./...

# Run linter
golangci-lint run
```

### Install Pre-commit Hooks

```bash
pre-commit install
pre-commit install --hook-type commit-msg
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test ./... -coverprofile=coverage.out

# View coverage
go tool cover -html=coverage.out
```

## Project Structure

```
terminus-go/
├── cmd/
│   └── terminus/          # CLI entry point
├── pkg/
│   ├── api/              # Pantheon API client (public)
│   │   ├── client.go     # HTTP client
│   │   ├── auth.go       # Authentication
│   │   ├── sites.go      # Sites API
│   │   ├── environments.go
│   │   ├── workflows.go
│   │   ├── backups.go
│   │   ├── organizations.go
│   │   ├── domains.go
│   │   ├── multidev.go
│   │   └── models/       # API data models
│   ├── config/           # Configuration management
│   ├── session/          # Session/token storage
│   └── output/           # Output formatting
├── internal/
│   └── commands/         # CLI commands (private)
└── test/
    └── fixtures/         # Test fixtures
```

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`go test ./...`)
5. Run linter (`golangci-lint run`)
6. Commit your changes using conventional commits
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Conventional Commits

This project uses [Conventional Commits](https://www.conventionalcommits.org/) for automated changelog generation:

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `test:` - Test changes
- `refactor:` - Code refactoring
- `chore:` - Maintenance tasks

## License

[MIT License](LICENSE)

## Acknowledgments

- Original [Pantheon Terminus](https://github.com/pantheon-systems/terminus) project
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management

## Support

For issues and questions:

- GitHub Issues: https://github.com/pantheon-systems/terminus-go/issues
- Pantheon Documentation: https://pantheon.io/docs

## Differences from PHP Terminus

While this Go version aims for feature parity with the PHP Terminus, there are some differences:

1. **Performance**: Go version is faster due to compiled nature
2. **Binary Distribution**: Single binary with no runtime dependencies
3. **API Package**: Go version provides a standalone API package for programmatic use
4. **Configuration**: Simplified configuration with fewer layers
5. **Plugin System**: Not yet implemented (planned for future release)

## Roadmap

- [x] Core API client
- [x] Authentication commands
- [x] Site management
- [x] Environment operations
- [x] Workflow management
- [x] Backup operations
- [x] Organization management
- [x] Domain management
- [x] Multidev operations
- [ ] HTTPS/SSL management
- [ ] Connection info commands
- [ ] SSH key management
- [ ] Import/export operations
- [ ] New Relic integration
- [ ] Redis/Solr management
- [ ] Plugin system
- [ ] Shell completion
- [ ] Man page generation

## Implementation Status

This section tracks the implementation status of all PHP Terminus commands in this Go version. Commands are compared against [PHP Terminus 3.x](https://github.com/pantheon-systems/terminus/tree/3.x).

**Legend:**
- ✅ = Implemented / Tested
- ❌ = Not implemented / Not tested

### aliases

| Command | Description | Implemented | Human Tested |
|---------|-------------|:-----------:|:------------:|
| `aliases` | Print all site aliases | ❌ | ❌ |

### art

| Command | Description | Implemented | Human Tested |
|---------|-------------|:-----------:|:------------:|
| `art` | Display Pantheon ASCII art | ✅ | ❌ |
| `art:list` | List available ASCII art | ✅ | ❌ |

### auth

| Command | Description | Implemented | Human Tested |
|---------|-------------|:-----------:|:------------:|
| `auth:login` | Log in to Pantheon using a machine token | ✅ | ❌ |
| `auth:logout` | Log out of Pantheon and delete saved session | ✅ | ❌ |
| `auth:whoami` | Display current user information | ✅ | ❌ |

### backup

| Command | Description | Implemented | Human Tested |
|---------|-------------|:-----------:|:------------:|
| `backup:automatic:disable` | Disable automatic backups for an environment | ✅ | ❌ |
| `backup:automatic:enable` | Enable automatic backups for an environment | ✅ | ❌ |
| `backup:automatic:info` | Show automatic backup schedule | ✅ | ❌ |
| `backup:create` | Create a backup of an environment | ✅ | ❌ |
| `backup:get` | Download a specific backup | ✅ | ❌ |
| `backup:info` | Show information about a specific backup | ✅ | ❌ |
| `backup:list` | List backups for an environment | ✅ | ❌ |
| `backup:restore` | Restore an environment from a backup | ✅ | ❌ |

### branch

| Command | Description | Implemented | Human Tested |
|---------|-------------|:-----------:|:------------:|
| `branch:list` | List git branches for a site | ❌ | ❌ |

### connection

| Command | Description | Implemented | Human Tested |
|---------|-------------|:-----------:|:------------:|
| `connection:info` | Show connection info for an environment | ✅ | ❌ |
| `connection:set` | Set connection mode (git/sftp) | ✅ | ❌ |

### dashboard

| Command | Description | Implemented | Human Tested |
|---------|-------------|:-----------:|:------------:|
| `dashboard:view` | Open site dashboard in a browser | ❌ | ❌ |

### domain

| Command | Description | Implemented | Human Tested |
|---------|-------------|:-----------:|:------------:|
| `domain:add` | Add a domain to an environment | ✅ | ❌ |
| `domain:dns` | Show DNS recommendations for a domain | ✅ | ❌ |
| `domain:list` | List domains for an environment | ✅ | ❌ |
| `domain:lookup` | Find the site associated with a domain | ❌ | ❌ |
| `domain:primary:add` | Add a primary domain to an environment | ❌ | ❌ |
| `domain:primary:remove` | Remove primary domain designation | ❌ | ❌ |
| `domain:remove` | Remove a domain from an environment | ✅ | ❌ |

### env

| Command | Description | Implemented | Human Tested |
|---------|-------------|:-----------:|:------------:|
| `env:clear-cache` | Clear caches for an environment | ✅ | ❌ |
| `env:clone-content` | Clone database and/or files between environments | ✅ | ❌ |
| `env:code-log` | Show code log for an environment | ❌ | ❌ |
| `env:code-rebuild` | Rebuild code for an environment | ❌ | ❌ |
| `env:commit` | Commit changes in SFTP mode | ✅ | ❌ |
| `env:deploy` | Deploy code to an environment | ✅ | ❌ |
| `env:diffstat` | Show diff statistics for an environment | ❌ | ❌ |
| `env:info` | Show environment information | ✅ | ❌ |
| `env:list` | List environments for a site | ✅ | ❌ |
| `env:metrics` | Show environment metrics | ❌ | ❌ |
| `env:rotate-random-seed` | Rotate the Drupal hash salt | ❌ | ❌ |
| `env:view` | Open environment in a browser | ❌ | ❌ |
| `env:wake` | Wake a sleeping environment | ❌ | ❌ |
| `env:wipe` | Wipe database and files from an environment | ✅ | ❌ |

### https

| Command | Description | Implemented | Human Tested |
|---------|-------------|:-----------:|:------------:|
| `https:info` | Show HTTPS/SSL information | ❌ | ❌ |
| `https:remove` | Remove HTTPS certificate | ❌ | ❌ |
| `https:set` | Enable HTTPS with a certificate | ❌ | ❌ |

### import

| Command | Description | Implemented | Human Tested |
|---------|-------------|:-----------:|:------------:|
| `import:complete` | Complete site import | ❌ | ❌ |
| `import:database` | Import database to an environment | ❌ | ❌ |
| `import:files` | Import files to an environment | ❌ | ❌ |
| `import:site` | Import a site archive | ❌ | ❌ |

### local

| Command | Description | Implemented | Human Tested |
|---------|-------------|:-----------:|:------------:|
| `local:clone` | Clone a Pantheon site locally | ❌ | ❌ |
| `local:commitAndPush` | Commit and push local changes | ❌ | ❌ |
| `local:dockerize` | Create Docker setup for local development | ❌ | ❌ |
| `local:getLiveDB` | Download database from live environment | ❌ | ❌ |
| `local:getLiveFiles` | Download files from live environment | ❌ | ❌ |

### lock

| Command | Description | Implemented | Human Tested |
|---------|-------------|:-----------:|:------------:|
| `lock:disable` | Disable HTTP basic auth for an environment | ✅ | ❌ |
| `lock:enable` | Enable HTTP basic auth for an environment | ✅ | ❌ |
| `lock:info` | Show lock status for an environment | ✅ | ❌ |

### machine-token

| Command | Description | Implemented | Human Tested |
|---------|-------------|:-----------:|:------------:|
| `machine-token:delete` | Delete a machine token | ❌ | ❌ |
| `machine-token:delete-all` | Delete all machine tokens | ❌ | ❌ |
| `machine-token:list` | List machine tokens | ❌ | ❌ |

### multidev

| Command | Description | Implemented | Human Tested |
|---------|-------------|:-----------:|:------------:|
| `multidev:create` | Create a multidev environment | ✅ | ❌ |
| `multidev:delete` | Delete a multidev environment | ✅ | ❌ |
| `multidev:list` | List multidev environments | ❌ | ❌ |
| `multidev:merge-from-dev` | Merge code from dev into multidev | ✅ | ❌ |
| `multidev:merge-to-dev` | Merge code from multidev to dev | ✅ | ❌ |

### new-relic

| Command | Description | Implemented | Human Tested |
|---------|-------------|:-----------:|:------------:|
| `new-relic:disable` | Disable New Relic for a site | ❌ | ❌ |
| `new-relic:enable` | Enable New Relic for a site | ❌ | ❌ |
| `new-relic:info` | Show New Relic information | ❌ | ❌ |

### org

| Command | Description | Implemented | Human Tested |
|---------|-------------|:-----------:|:------------:|
| `org:info` | Show organization information | ✅ | ❌ |
| `org:list` | List organizations | ✅ | ❌ |
| `org:people:list` | List organization members | ✅ | ❌ |
| `org:site:list` | List sites belonging to an organization | ✅ | ❌ |
| `org:upstream:list` | List upstreams for an organization | ✅ | ❌ |

### owner

| Command | Description | Implemented | Human Tested |
|---------|-------------|:-----------:|:------------:|
| `owner:set` | Change site owner | ❌ | ❌ |

### payment-method

| Command | Description | Implemented | Human Tested |
|---------|-------------|:-----------:|:------------:|
| `payment-method:add` | Add a payment method | ❌ | ❌ |
| `payment-method:list` | List payment methods | ❌ | ❌ |
| `payment-method:remove` | Remove a payment method | ❌ | ❌ |

### plan

| Command | Description | Implemented | Human Tested |
|---------|-------------|:-----------:|:------------:|
| `plan:info` | Show site plan information | ✅ | ❌ |
| `plan:list` | List available plans | ❌ | ❌ |
| `plan:set` | Change the site plan | ❌ | ❌ |

### redis

| Command | Description | Implemented | Human Tested |
|---------|-------------|:-----------:|:------------:|
| `redis:disable` | Disable Redis for a site | ✅ | ❌ |
| `redis:enable` | Enable Redis for a site | ✅ | ❌ |

### remote

| Command | Description | Implemented | Human Tested |
|---------|-------------|:-----------:|:------------:|
| `drush` | Run a Drush command on a site | ❌ | ❌ |
| `wp` | Run a WP-CLI command on a site | ❌ | ❌ |

### self

| Command | Description | Implemented | Human Tested |
|---------|-------------|:-----------:|:------------:|
| `self:clear-cache` | Clear Terminus cache | ❌ | ❌ |
| `self:config:dump` | Dump Terminus configuration | ❌ | ❌ |
| `self:console` | Open interactive console | ❌ | ❌ |
| `self:info` | Show Terminus information | ✅ | ❌ |
| `self:plugin:create` | Create a new plugin | ❌ | ❌ |
| `self:plugin:install` | Install a plugin | ❌ | ❌ |
| `self:plugin:list` | List installed plugins | ❌ | ❌ |
| `self:plugin:reload` | Reload plugins | ❌ | ❌ |
| `self:plugin:search` | Search for plugins | ❌ | ❌ |
| `self:plugin:uninstall` | Uninstall a plugin | ❌ | ❌ |
| `self:plugin:update` | Update a plugin | ❌ | ❌ |

### service-level

| Command | Description | Implemented | Human Tested |
|---------|-------------|:-----------:|:------------:|
| `service-level:set` | Set the service level of a site | ❌ | ❌ |

### site

| Command | Description | Implemented | Human Tested |
|---------|-------------|:-----------:|:------------:|
| `site:create` | Create a new site | ✅ | ❌ |
| `site:delete` | Delete a site | ✅ | ❌ |
| `site:info` | Show site information | ✅ | ❌ |
| `site:label` | Set site label | ❌ | ❌ |
| `site:list` | List sites | ✅ | ❌ |
| `site:lookup` | Look up a site by UUID | ❌ | ❌ |
| `site:org:add` | Add site to an organization | ❌ | ❌ |
| `site:org:list` | List organizations a site belongs to | ❌ | ❌ |
| `site:org:remove` | Remove site from an organization | ❌ | ❌ |
| `site:team:add` | Add a user to the site team | ❌ | ❌ |
| `site:team:list` | List site team members | ✅ | ❌ |
| `site:team:remove` | Remove a user from the site team | ❌ | ❌ |
| `site:team:role` | Change a team member's role | ❌ | ❌ |
| `site:upstream:clear-cache` | Clear upstream cache | ❌ | ❌ |
| `site:upstream:set` | Set the upstream for a site | ❌ | ❌ |

### solr

| Command | Description | Implemented | Human Tested |
|---------|-------------|:-----------:|:------------:|
| `solr:disable` | Disable Solr for a site | ❌ | ❌ |
| `solr:enable` | Enable Solr for a site | ❌ | ❌ |

### ssh-key

| Command | Description | Implemented | Human Tested |
|---------|-------------|:-----------:|:------------:|
| `ssh-key:add` | Add an SSH key to your account | ❌ | ❌ |
| `ssh-key:list` | List SSH keys on your account | ❌ | ❌ |
| `ssh-key:remove` | Remove an SSH key from your account | ❌ | ❌ |

### tag

| Command | Description | Implemented | Human Tested |
|---------|-------------|:-----------:|:------------:|
| `tag:add` | Add a tag to a site | ❌ | ❌ |
| `tag:list` | List tags for a site | ❌ | ❌ |
| `tag:remove` | Remove a tag from a site | ❌ | ❌ |

### upstream

| Command | Description | Implemented | Human Tested |
|---------|-------------|:-----------:|:------------:|
| `upstream:info` | Show upstream information | ✅ | ❌ |
| `upstream:list` | List upstreams | ✅ | ❌ |
| `upstream:updates:apply` | Apply upstream updates to a site | ❌ | ❌ |
| `upstream:updates:list` | List available upstream updates | ❌ | ❌ |
| `upstream:updates:status` | Check for upstream updates | ❌ | ❌ |

### workflow

| Command | Description | Implemented | Human Tested |
|---------|-------------|:-----------:|:------------:|
| `workflow:info` | Show workflow information | ✅ | ❌ |
| `workflow:list` | List workflows for a site | ✅ | ❌ |
| `workflow:wait` | Wait for a workflow to complete | ✅ | ❌ |
| `workflow:watch` | Watch a workflow with live progress | ✅ | ❌ |

### Implementation Summary

| Status | Count |
|--------|-------|
| **Total Commands** | 113 |
| **Implemented** | 48 |
| **Not Implemented** | 65 |
| **Implementation Progress** | 42% |
