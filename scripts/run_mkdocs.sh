#!/bin/bash

# Script to run MkDocs locally with the config file in the docs folder

# Check if we're in the virtual environment
if [[ "$VIRTUAL_ENV" == "" ]]; then
  echo "Creating/activating virtual environment..."
  if [ ! -d ".venv" ]; then
    python3 -m venv .venv
  fi
  source .venv/bin/activate
fi

# Install MkDocs if not already installed
if ! command -v mkdocs &> /dev/null; then
  echo "Installing MkDocs and dependencies..."
  pip install mkdocs mkdocs-material
fi

# Create mkdocs.yml in the root directory
echo "Starting MkDocs server..."
bash .github/scripts/create-versioned-mkdocs-config.sh "latest"

# Run MkDocs
echo "Access documentation at http://127.0.0.1:8000/"
echo "Press Ctrl+C to stop the server"
mkdocs serve

# Clean up after MkDocs is stopped
echo "Cleaning up..."
rm mkdocs.yml
