#!/bin/bash

# Local release script for ctRestClient
# This script runs GoReleaser in snapshot mode to test the release locally
# Usage: ./scripts/local_release.sh

set -e

echo "Running local release test..."

# Check if GoReleaser is installed
GORELEASER_PATH=""
if command -v goreleaser &> /dev/null; then
    GORELEASER_PATH="goreleaser"
elif [ -f "/home/matthias/.asdf/installs/golang/1.24.4/bin/goreleaser" ]; then
    GORELEASER_PATH="/home/matthias/.asdf/installs/golang/1.24.4/bin/goreleaser"
else
    echo "GoReleaser not found. Installing..."
    go install github.com/goreleaser/goreleaser@latest
    # Try to find it after installation
    if [ -f "/home/matthias/.asdf/installs/golang/1.24.4/bin/goreleaser" ]; then
        GORELEASER_PATH="/home/matthias/.asdf/installs/golang/1.24.4/bin/goreleaser"
    elif command -v goreleaser &> /dev/null; then
        GORELEASER_PATH="goreleaser"
    else
        echo "Failed to install or locate GoReleaser"
        exit 1
    fi
fi

# Run GoReleaser in snapshot mode
echo "Building release artifacts locally..."
"$GORELEASER_PATH" release --snapshot --clean

echo "Local release test completed!"
echo "Check the 'dist/' directory for generated binaries and archives."
echo ""
echo "Generated files:"
ls -la dist/ | grep -E "\.(tar\.gz|exe)$|checksums\.txt" | head -10
