# Makefile for terminus-golang local development
#
# For CI builds, use GoReleaser directly via the CI workflow.
# This Makefile is for local development convenience.

# Version is determined by git describe:
# - On a tag: returns the tag (e.g., "v0.6.0")
# - Off a tag: returns "<tag>-<commits>-g<hash>" or just "<hash>" if no tags
# - With dirty tree: appends "-dirty"
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -X github.com/deviantintegral/terminus-golang/pkg/version.version=$(VERSION)

.PHONY: build build-all clean test lint

# Build for current platform
build:
	go build -ldflags "$(LDFLAGS)" -o bin/terminus ./cmd/terminus

# Build for all platforms using GoReleaser (snapshot mode)
build-all:
	goreleaser build --snapshot --clean

# Run tests
test:
	go test -v -race ./...

# Run linter
lint:
	golangci-lint run --timeout=5m

# Clean build artifacts
clean:
	rm -rf bin/ dist/

# Default target
.DEFAULT_GOAL := build
