#!/bin/bash

# SessionStart hook for terminus-golang
# This hook installs required development tools when a Claude Code session starts

set -e

echo "=== SessionStart Hook: Installing development tools ==="

# Install pre-commit if not already installed
if ! command -v pre-commit &> /dev/null; then
    echo "Installing pre-commit..."
    if command -v brew &> /dev/null; then
        brew install pre-commit
    elif command -v pip3 &> /dev/null; then
        pip3 install pre-commit
    elif command -v pip &> /dev/null; then
        pip install pre-commit
    else
        echo "Warning: brew or pip not found. Skipping pre-commit installation."
        echo "Please install pre-commit manually: https://pre-commit.com/#install"
    fi
else
    echo "pre-commit is already installed ($(pre-commit --version))"
fi

# Install golangci-lint v2.6.2 if not already installed or wrong version
GOLANGCI_VERSION="v2.6.2"
GOLANGCI_INSTALLED_VERSION=""

if ! command -v golangci-lint &> /dev/null; then
  echo "Installing golangci-lint $GOLANGCI_VERSION..."

  # Download and install golangci-lint
  if command -v brew &> /dev/null; then
      brew install golangci-lint
  else
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin $GOLANGCI_VERSION
  fi

  # Verify installation
  if command -v golangci-lint &> /dev/null; then
      echo "golangci-lint installed successfully: $(golangci-lint --version)"
  else
      echo "Warning: golangci-lint installation may have failed"
      echo "Please ensure $(go env GOPATH)/bin is in your PATH"
  fi
fi

# Ensure GOPATH/bin is in PATH so installed Go tools are accessible
export PATH="$(go env GOPATH)/bin:$PATH"

# Install goimports
echo "Installing goimports..."
go install golang.org/x/tools/cmd/goimports@latest

# Verify goimports installation
if command -v goimports &> /dev/null; then
    echo "goimports installed successfully: $(which goimports)"
else
    echo "Warning: goimports installation may have failed"
    echo "Please ensure $(go env GOPATH)/bin is in your PATH"
fi

# Install pre-commit hooks if not already installed
if [ -f ".pre-commit-config.yaml" ] && command -v pre-commit &> /dev/null; then
    echo "Installing pre-commit hooks..."
    pre-commit install --install-hooks
    pre-commit install --hook-type commit-msg
fi

echo "=== SessionStart Hook: Complete ==="
