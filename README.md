# Kavach CLI

Enterprise-grade secret management and synchronization CLI tool.

## 🚀 Features

- **Multi-Platform Support**: Linux, macOS, Windows (AMD64, ARM64)
- **Cloud Provider Integration**: Azure Key Vault, GCP Secret Manager, GitHub Actions
- **Role-Based Access Control**: Fine-grained permissions and user management
- **Version Control**: Git-like workflow for secret management
- **Secure**: End-to-end encryption and secure authentication

## 📦 Installation

### From Releases

Download the latest release for your platform from [GitHub Releases](https://github.com/Gkemhcs/kavach-cli/releases).

```bash
# Linux
curl -L https://github.com/Gkemhcs/kavach-cli/releases/latest/download/kavach-cli_Linux_x86_64.tar.gz | tar -xz
sudo mv kavach /usr/local/bin/

# macOS
brew install Gkemhcs/kavach-cli/kavach-cli

# Windows
scoop install kavach-cli
```

### From Source

```bash
# Clone the repository
git clone https://github.com/Gkemhcs/kavach-cli.git
cd kavach-cli

# Build and install
make install
```

## 🔧 Development Setup

### Prerequisites

- Go 1.21 or higher
- Git
- Make (optional, for using Makefile)

### Quick Start

```bash
# Clone the repository
git clone https://github.com/Gkemhcs/kavach-cli.git
cd kavach-cli

# Setup development environment
make setup

# Build for current platform
make build

# Run tests
make test

# Install locally
make install
```

### Available Make Commands

```bash
make help          # Show all available commands
make build         # Build for current platform
make build-all     # Build for all platforms
make test          # Run tests
make lint          # Run linter
make format        # Format code
make install       # Install locally
make clean         # Clean build artifacts
make version       # Show version information
```

## 🚀 Building Releases

### Local Development Build

```bash
# Build with version information
make build

# Test the binary
./kavach version
./kavach --help
```

### Multi-Platform Build

```bash
# Build for all supported platforms
make build-all

# Binaries will be created in dist/ directory
ls dist/
```

### Tag-Based Versioning

The CLI uses git tags for versioning. Simply create and push a tag to trigger a release:

```bash
# Create and push a tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

This will automatically:
- Build binaries for all platforms
- Create a GitHub release
- Upload artifacts
- Generate checksums

#### Snapshot Release (Development)

```bash
# Build snapshot release (no tag needed)
make snapshot

# Or directly with GoReleaser
goreleaser build --snapshot --clean
```

## 🐳 Docker

### Build Docker Image

```bash
# Build image
make docker-build

# Run container
docker run --rm kavach-cli:latest version
```

### Using Docker for Development

```bash
# Build development image
docker build -t kavach-cli:dev .

# Run with volume mount for development
docker run -it --rm -v $(pwd):/app kavach-cli:dev
```

## 📋 Version Information

The CLI includes comprehensive version information:

```bash
# Show version
kavach version

# Short version
kavach version --short

# JSON format
kavach version --json
```

### Version Variables

The following information is embedded in the binary:

- **Version**: Semantic version from git tag
- **Build Time**: When the binary was built
- **Git Commit**: Git commit hash
- **Git Branch**: Git branch name
- **Go Version**: Go runtime version
- **Platform**: Target OS/Architecture

## 🔄 CI/CD Integration

### GitHub Actions

The repository includes GitHub Actions workflows:

- **Release**: Automatically builds and releases when tags are pushed
- **Snapshot**: Builds snapshot releases on every push to main

### Release Process

1. **Create a tag**:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

2. **GitHub Actions automatically**:
   - Builds binaries for all platforms
   - Creates GitHub release
   - Uploads artifacts
   - Generates checksums

3. **Verify the release**:
   ```bash
   # Download and test
   wget https://github.com/Gkemhcs/kavach-cli/releases/download/v1.0.0/kavach-cli_Linux_x86_64.tar.gz
   tar -xzf kavach-cli_Linux_x86_64.tar.gz
   ./kavach version
   ```

## 🏗️ Supported Platforms

| OS | Architecture | Status |
|----|-------------|--------|
| Linux | AMD64 | ✅ |
| Linux | ARM64 | ✅ |
| macOS | AMD64 | ✅ |
| macOS | ARM64 | ✅ |
| Windows | AMD64 | ✅ |
| Windows | ARM64 | ❌ |

## 📦 Package Managers

### Homebrew (macOS)

```bash
brew install Gkemhcs/kavach-cli/kavach-cli
```

### Scoop (Windows)

```bash
scoop bucket add kavach-cli https://github.com/Gkemhcs/scoop-kavach-cli
scoop install kavach-cli
```

## 🔍 Troubleshooting

### Common Issues

1. **Permission Denied**:
   ```bash
   chmod +x kavach
   ```

2. **Go Version Too Old**:
   ```bash
   # Update Go to 1.21+
   go version
   ```

3. **Build Fails**:
   ```bash
   # Clean and rebuild
   make clean
   make build
   ```

### Debug Information

```bash
# Show debug information
kavach version --json

# Check environment
echo $GOPATH
echo $GOROOT
go env
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run `make test` and `make lint`
6. Submit a pull request

### Development Guidelines

- Follow Go coding standards
- Add tests for new features
- Update documentation
- Use conventional commit messages

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- [Documentation](https://github.com/Gkemhcs/kavach-docs)
- [Issues](https://github.com/Gkemhcs/kavach-cli/issues)
- [Discussions](https://github.com/Gkemhcs/kavach-cli/discussions)
- [Releases](https://github.com/Gkemhcs/kavach-cli/releases)
