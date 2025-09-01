#!/bin/bash

# Create mkdocs.yml file
cat > mkdocs.yml << 'EOF'
site_name: ctRestClient
site_description: Documentation for ctRestClient
site_url: https://mvach.github.io/ctRestClient/
theme:
  name: material
  palette:
    primary: indigo
    accent: indigo
  features:
    - navigation.tabs  # This moves the navigation to the top as tabs
    - navigation.tabs.sticky  # Makes the tabs sticky at the top
    - navigation.sections  # Renders sections as groups in the sidebar
    - navigation.expand  # Expands all collapsible sections
    - toc.follow  # Automatically follows the table of contents on the right
    - content.code.copy  # Add copy button to code blocks
use_directory_urls: true
plugins:
  - search
markdown_extensions:
  - admonition
  - pymdownx.details
  - pymdownx.superfences
  - toc:
      permalink: true
nav:
  - 🇩🇪 German: manual-de.md
  - 🇬🇧 English: manual-en.md
EOF