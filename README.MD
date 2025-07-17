# ctRestClient

A command-line tool for exporting ChurchTools groups to CSV files, designed for creating mail merge documents with tools like Microsoft Word or Adobe InDesign.

## Features

- Export ChurchTools groups to CSV format
- Secure token management using KeePass databases
- Multi-instance support for multiple ChurchTools installations
- Configurable data fields and value transformations
- Cross-platform support (Windows, macOS, Linux)

## Quick Start

For detailed user instructions, see:
- [docs/manual-en.md](docs/manual-en.md) (English user manual)
- [docs/manual-de.md](docs/manual-de.md) (German user manual)

### Basic Usage

```bash
# Download the appropriate binary for your platform
./ctRestClient-linux-amd64 -c config.yml -k tokens.kdbx
```

### Configuration Example

```yaml
instances:
  - hostname: your-church.church.tools
    token_name: your-token-name
    groups:
    - name: "Youth Group"
      fields: [id, firstName, lastName, email]
```

## Development

### Prerequisites

- Go 1.21 or higher
- Git

### Getting Started

```bash
# Clone the repository
git clone https://github.com/mvach/ctRestClient.git
cd ctRestClient

# Install dependencies
go mod download

# Build for all platforms
./scripts/build_binaries

# Run tests
go test ./...
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for specific packages
go test ./app

# Show test coverage
go test -cover ./...

# Run integration tests (requires test environment)
go test ./integration
```

### Generating Test Fakes

The project uses Counterfeiter for generating test fakes:

```bash
# Generate all fakes
./scripts/generate_fakes

# Or generate manually for specific modules
go generate ./app
```

### Building

```bash
# Build for current platform
go build -o bin/ctRestClient ./main.go

# Build for all platforms
./scripts/build_binaries

# The script creates binaries in the bin/ directory:
# - ctRestClient-windows-amd64.exe
# - ctRestClient-darwin-amd64
# - ctRestClient-darwin-arm64
# - ctRestClient-linux-amd64
```

## Documentation

- [docs/manual-en.md](docs/manual-en.md) - Complete English user manual
- [docs/manual-de.md](docs/manual-de.md) - Complete German user manual

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
