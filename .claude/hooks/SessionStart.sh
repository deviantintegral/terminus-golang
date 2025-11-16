#!/bin/bash

# SessionStart hook for terminus-golang
# This hook installs required development tools when a Claude Code session starts

set -e

echo "=== SessionStart Hook: Installing development tools ==="

# Install pre-commit if not already installed
if ! command -v pre-commit &> /dev/null; then
    echo "Installing pre-commit..."
    if command -v pip3 &> /dev/null; then
        pip3 install pre-commit
    elif command -v pip &> /dev/null; then
        pip install pre-commit
    else
        echo "Warning: pip not found. Skipping pre-commit installation."
        echo "Please install pre-commit manually: https://pre-commit.com/#install"
    fi
else
    echo "pre-commit is already installed ($(pre-commit --version))"
fi

# Install golangci-lint v2.6.2 if not already installed or wrong version
GOLANGCI_VERSION="v2.6.2"
GOLANGCI_INSTALLED_VERSION=""

if command -v golangci-lint &> /dev/null; then
    GOLANGCI_INSTALLED_VERSION=$(golangci-lint --version 2>&1 | grep -oP 'golangci-lint has version \K[^ ]+' || echo "unknown")
fi

if [ "$GOLANGCI_INSTALLED_VERSION" != "$GOLANGCI_VERSION" ]; then
    echo "Installing golangci-lint $GOLANGCI_VERSION..."

    # Download and install golangci-lint
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin $GOLANGCI_VERSION

    # Verify installation
    if command -v golangci-lint &> /dev/null; then
        echo "golangci-lint installed successfully: $(golangci-lint --version)"
    else
        echo "Warning: golangci-lint installation may have failed"
        echo "Please ensure $(go env GOPATH)/bin is in your PATH"
    fi
else
    echo "golangci-lint $GOLANGCI_VERSION is already installed"
fi

# Install pre-commit hooks if not already installed
if [ -f ".pre-commit-config.yaml" ] && command -v pre-commit &> /dev/null; then
    echo "Installing pre-commit hooks..."
    pre-commit install --install-hooks 2>/dev/null || echo "pre-commit hooks already installed"
    pre-commit install --hook-type commit-msg 2>/dev/null || true
fi

echo "=== SessionStart Hook: Complete ==="
