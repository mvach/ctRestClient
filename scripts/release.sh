#!/bin/bash

# Release script for ctRestClient
# Usage: ./scripts/release.sh <version>
# Example: ./scripts/release.sh v1.0.0

set -e

if [ $# -ne 2 ]; then
    echo "Usage: $0 <version> <origin>"
    echo "Example: $0 v1.0.0 origin"
    exit 1
fi

ORIGIN=$1
VERSION=$2


# Validate version format (should start with v)
if [[ ! $VERSION =~ ^v ]]; then
    echo "Version should start with 'v' (e.g., v1.0.0)"
    exit 1
fi

echo "Creating release $VERSION..."

# Check if working directory is clean
if [ -n "$(git status --porcelain)" ]; then
    echo "Working directory is not clean. Please commit or stash changes first."
    exit 1
fi

# Create and push tag
echo "Creating tag $VERSION..."
git tag "$VERSION"

echo "Pushing tag to $ORIGIN..."
git push "$ORIGIN" "$VERSION"

echo "Release $VERSION created successfully!"
echo "GitHub Actions will now build and publish the release automatically."
echo "Check the Actions tab in your repository for progress."