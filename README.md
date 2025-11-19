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
