# Kavach CLI

A secure and powerful command-line interface for managing secrets, environments, and access controls in the Kavach platform.

## 🚀 Quick Start

### Using Docker (Recommended)

The easiest way to use Kavach CLI is via Docker. We provide **public**, signed, and verified Docker images:

```bash
# Pull and run the latest version (no authentication required)
docker run --rm ghcr.io/gkemhcs/kavach-cli:latest --help

# Run with version info
docker run --rm ghcr.io/gkemhcs/kavach-cli:latest version

# Run a specific command
docker run --rm ghcr.io/gkemhcs/kavach-cli:latest env list --org my-org
```

**Note**: Our Docker images are **publicly accessible** - no authentication required to pull them!

### Verify Docker Image Integrity

All our Docker images are signed with cosign and include SLSA provenance attestations:

```bash
# Install cosign (if not already installed)
go install github.com/sigstore/cosign/cmd/cosign@latest

# Verify image signature
cosign verify ghcr.io/gkemhcs/kavach-cli:latest

# Verify attestation (SLSA provenance)
cosign verify-attestation --type slsaprovenance --keyless ghcr.io/gkemhcs/kavach-cli:latest
```

### Using Our Verification Script

We provide a verification script to easily verify Docker images:

```bash
# Make it executable
chmod +x scripts/verify-docker.sh

# Verify latest image
./scripts/verify-docker.sh

# Verify specific version
./scripts/verify-docker.sh -t v1.0.0

# Verify all available tags
./scripts/verify-docker.sh -a
```

### Docker Image Tags

- `latest` - Latest stable release
- `v1.0.0` - Specific version tags
- `v1.0` - Major.minor version tags
- `main-abc1234` - Branch-based tags for development

### Docker Registry

Our images are hosted on GitHub Container Registry (ghcr.io) and include:
- Multi-platform support (linux/amd64, linux/arm64)
- Optimized layers with GitHub Actions caching
- Security scanning and vulnerability checks
- Automated builds on every release

## 📦 Installation

### Binary Installation

Download pre-built binaries for your platform:

```bash
# Linux AMD64
curl -L -o kavach https://github.com/Gkemhcs/kavach-cli/releases/latest/download/kavach-cli_Linux_x86_64.tar.gz
tar -xzf kavach-cli_Linux_x86_64.tar.gz
chmod +x kavach
sudo mv kavach /usr/local/bin/

# macOS AMD64
curl -L -o kavach https://github.com/Gkemhcs/kavach-cli/releases/latest/download/kavach-cli_Darwin_x86_64.tar.gz
tar -xzf kavach-cli_Darwin_x86_64.tar.gz
chmod +x kavach
sudo mv kavach /usr/local/bin/

# Windows AMD64
# Download from GitHub releases and extract
```

### Verify Binary Integrity

All binaries are signed with cosign:

```bash
# Download signature and certificate
curl -O https://github.com/Gkemhcs/kavach-cli/releases/latest/download/kavach-cli_Linux_x86_64.sig
curl -O https://github.com/Gkemhcs/kavach-cli/releases/latest/download/kavach-cli_Linux_x86_64.pem

# Verify signature
cosign verify-blob --cert kavach-cli_Linux_x86_64.pem --signature kavach-cli_Linux_x86_64.sig kavach
```

### Using Our Verification Script

```bash
# Make it executable
chmod +x scripts/verify-release.sh

# Verify latest release
./scripts/verify-release.sh

# Verify specific version
./scripts/verify-release.sh v1.0.0
```

## 🔧 Development

### Prerequisites

- Go 1.21+
- Docker (for containerized builds)
- cosign (for signing verification)

### Build from Source

```bash
# Clone the repository
git clone https://github.com/Gkemhcs/kavach-cli.git
cd kavach-cli

# Install dependencies
go mod download

# Build binary
go build -o kavach ./main.go

# Run tests
go test ./...

# Build for all platforms
make build-all
```

### Docker Development

```bash
# Build development image
docker build -t kavach-cli:dev .

# Run development image
docker run --rm kavach-cli:dev --help

# Build with specific version
docker build --build-arg VERSION=v1.0.0 -t kavach-cli:v1.0.0 .
```

### Release Process

```bash
# Create a new release
git tag v1.0.0
git push origin v1.0.0

# This triggers GitHub Actions which:
# 1. Builds binaries for all platforms
# 2. Creates Docker images with caching
# 3. Signs all artifacts with cosign
# 4. Creates SLSA attestations
# 5. Pushes to GitHub Container Registry
```

## 🔒 Security Features

### Supply Chain Security

- **Cosign Signatures**: All binaries and Docker images are signed
- **SLSA Attestations**: Provenance information for build process
- **Transparency Log**: Signatures uploaded to Sigstore Rekor
- **Keyless Signing**: Uses OIDC with GitHub Actions

### Docker Security

- **Multi-platform**: Support for linux/amd64 and linux/arm64
- **Non-root User**: Runs as unprivileged user (UID 1001)
- **Minimal Base**: Alpine Linux with only necessary packages
- **Health Checks**: Built-in health monitoring
- **Layer Caching**: Optimized builds with GitHub Actions cache

### Binary Security

- **Static Linking**: No external dependencies
- **Hardened Runtime**: macOS binaries include security features
- **Code Signing**: Proper code signing for macOS
- **Permission Controls**: Secure file permissions

## 📚 Usage Examples

### Basic Commands

```bash
# Get help
kavach --help

# Check version
kavach version

# List environments
kavach env list --org my-org --secret-group backend

# Grant access
kavach env grant production --user john.doe --role editor --org my-org --secret-group backend
```

### Docker Usage

```bash
# Interactive mode
docker run -it --rm ghcr.io/gkemhcs/kavach-cli:latest

# Mount configuration
docker run --rm -v ~/.kavach:/app/config ghcr.io/gkemhcs/kavach-cli:latest

# Environment variables
docker run --rm -e KAVACH_TOKEN=your-token ghcr.io/gkemhcs/kavach-cli:latest
```

## 🐛 Troubleshooting

### Common Issues

**Docker Image Won't Run**
```bash
# Check if image exists
docker images | grep kavach-cli

# Pull latest image
docker pull ghcr.io/gkemhcs/kavach-cli:latest

# Verify image integrity
./scripts/verify-docker.sh
```

**Binary Won't Execute**
```bash
# Check permissions
ls -la kavach

# Make executable
chmod +x kavach

# Verify signature
./scripts/verify-release.sh
```

**Verification Fails**
```bash
# Install cosign
go install github.com/sigstore/cosign/cmd/cosign@latest

# Check cosign version
cosign version

# Verify manually
cosign verify ghcr.io/gkemhcs/kavach-cli:latest
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- [GitHub Repository](https://github.com/Gkemhcs/kavach-cli)
- [GitHub Container Registry](https://ghcr.io/gkemhcs/kavach-cli)
- [Releases](https://github.com/Gkemhcs/kavach-cli/releases)
- [Documentation](https://docs.kavach.dev)
- [Issues](https://github.com/Gkemhcs/kavach-cli/issues)
