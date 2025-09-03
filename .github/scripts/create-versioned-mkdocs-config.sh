#!/bin/bash

VERSION="${1:-latest}"

# Remove 'v' prefix from version for artifact names
VERSION_NO_V="${VERSION#v}"

echo "Processing documentation for version: $VERSION"

# Create mkdocs.yml file with version-specific content
echo "Creating mkdocs.yml for version: $VERSION..."
cat > mkdocs.yml << EOF
site_name: ctRestClient
site_description: Documentation for ctRestClient
site_url: https://mvach.github.io/ctRestClient/
theme:
  name: material
  palette:
    primary: indigo
    accent: indigo
  features:
    - navigation.tabs
    - navigation.tabs.sticky
    - navigation.sections
    - navigation.expand
    - toc.follow
    - content.code.copy
    - navigation.top
use_directory_urls: true
plugins:
  - search
markdown_extensions:
  - admonition
  - pymdownx.details
  - pymdownx.superfences
  - pymdownx.emoji:
      emoji_index: !!python/name:material.extensions.emoji.twemoji
      emoji_generator: !!python/name:material.extensions.emoji.to_svg
  - toc:
      permalink: true
nav:
  - 🇩🇪 German: manual-de.md
  - 🇬🇧 English: manual-en.md
extra:
  version:
    provider: mike
    default: latest
EOF

# Process template files (*.md.tmp) to create versioned documentation
if [ -d "docs" ]; then
    # Process German manual template
    if [ -f "docs/manual-de.md.tmp" ]; then
        echo "Processing docs/manual-de.md.tmp -> docs/manual-de.md..."
        sed -e "s/\${version}/$VERSION/g" -e "s/\${version_no_v}/$VERSION_NO_V/g" "docs/manual-de.md.tmp" > "docs/manual-de.md"
    fi
    
    # Process English manual template
    if [ -f "docs/manual-en.md.tmp" ]; then
        echo "Processing docs/manual-en.md.tmp -> docs/manual-en.md..."
        sed -e "s/\${version}/$VERSION/g" -e "s/\${version_no_v}/$VERSION_NO_V/g" "docs/manual-en.md.tmp" > "docs/manual-en.md"
    fi
fi

echo "Template variables replaced successfully for version: $VERSION"
