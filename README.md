# Kavach CLI

A powerful, enterprise-grade command-line interface for secure secrets management, built with Go and designed for DevOps engineers, developers, and security professionals.

## 🚀 Features

### 🔐 **Secrets Management**
- **Secure Secret Operations**: Create, read, update, and delete secrets with end-to-end encryption
- **Version Control**: Git-like workflow for secret management with commit messages
- **Environment Support**: Manage secrets across multiple environments (dev, staging, prod)
- **Secret Groups**: Organize secrets into logical groups within organizations
- **Bulk Operations**: Manage multiple secrets simultaneously

### 🛡️ **Identity & Access Management**
- **OAuth 2.0 Authentication**: Secure GitHub OAuth integration for seamless login
- **Device Flow Support**: CLI-friendly authentication for headless and CI/CD environments
- **Role-Based Access Control**: Leverage backend RBAC for fine-grained permissions
- **Session Management**: Automatic token refresh and secure credential storage
- **Multi-User Support**: Manage multiple user accounts and organizations

### 🔌 **Multi-Provider Integration**
- **GitHub Secrets**: Sync secrets to GitHub repositories and environments
- **Google Cloud Platform**: Integration with GCP Secret Manager
- **Azure Key Vault**: Sync secrets to Azure Key Vault
- **Provider Credentials**: Secure storage and management of provider API keys
- **Cross-Platform Sync**: Synchronize secrets across multiple cloud providers

### 🏢 **Organization Management**
- **Multi-Organization Support**: Work with multiple organizations from a single CLI
- **User Groups**: Create and manage user groups for bulk permission management
- **Environment Management**: Create and configure environment-specific settings
- **Resource Hierarchy**: Navigate through organizations, secret groups, and environments

### 📊 **Developer Experience**
- **Interactive Prompts**: User-friendly prompts for complex operations
- **Tabular Output**: Clean, formatted output using table writers
- **JSON Support**: Full JSON output for scripting and automation
- **Comprehensive Help**: Detailed help for all commands and subcommands
- **Auto-completion**: Shell completion for bash, zsh, and fish

## 🏗️ Architecture

### **Core Components**

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   CLI Layer     │    │  Client Layer   │    │   Backend API   │
│   (Cobra CMD)   │◄──►│   (HTTP Client) │◄──►│   (Kavach API)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Configuration  │    │   Authentication│    │   Response      │
│   (Viper)       │    │   (OAuth/JWT)   │    │   Processing    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### **Command Structure**
- **Root Commands**: Core CLI functionality and global options
- **Domain Commands**: Organization, secret, provider, and environment management
- **Subcommands**: Specific operations within each domain
- **Global Flags**: Configuration, authentication, and output options

### **Data Flow**
1. **Command Parsing**: Cobra parses command-line arguments and flags
2. **Configuration Loading**: Viper loads configuration from multiple sources
3. **Authentication Check**: Validates user authentication and token validity
4. **API Request**: HTTP client sends requests to backend API
5. **Response Processing**: Processes and formats API responses
6. **Output Rendering**: Displays results in user-friendly format

## 🛠️ Technology Stack

### **Core Technologies**
- **Go 1.23**: High-performance, compiled language
- **Cobra**: Powerful CLI framework for Go applications
- **Viper**: Configuration management with multiple sources
- **Zerolog**: Structured logging with zero allocation

### **HTTP & Networking**
- **Standard HTTP Client**: Go's built-in HTTP client with custom middleware
- **JWT Management**: Secure token handling and refresh
- **OAuth 2.0**: GitHub OAuth integration for authentication

### **User Interface**
- **Table Writer**: Beautiful tabular output formatting
- **Color Support**: Cross-platform color output
- **Interactive Prompts**: User-friendly input handling
- **Progress Indicators**: Visual feedback for long-running operations

### **Development & Testing**
- **Testify**: Testing framework with mocks and assertions
- **GoReleaser**: Automated release management
- **Docker**: Containerization and development environment
- **Make**: Build automation and development workflows

## 📁 Project Structure

```
cli/
├── cmd/                    # Command definitions
│   ├── root.go            # Root command and CLI setup
│   ├── version.go         # Version information
│   ├── status.go          # System status
│   ├── info.go            # System information
│   ├── login/             # Authentication commands
│   ├── logout/            # Logout functionality
│   ├── org/               # Organization management
│   ├── secretgroup/       # Secret group operations
│   ├── environment/       # Environment management
│   ├── secret/            # Secret operations
│   ├── provider/          # Provider integrations
│   └── user-group/        # User group management
├── internal/               # Private application code
│   ├── auth/              # Authentication logic
│   ├── client/            # HTTP client implementations
│   ├── config/            # Configuration management
│   ├── errors/            # Custom error types
│   ├── groups/            # User group operations
│   ├── org/               # Organization client
│   ├── provider/          # Provider client
│   ├── secret/            # Secret client
│   ├── secretgroup/       # Secret group client
│   ├── environment/       # Environment client
│   ├── types/             # Shared type definitions
│   ├── utils/             # Utility functions
│   └── version/           # Version information
├── scripts/                # Build and deployment scripts
├── Dockerfile              # Container image definition
├── Makefile                # Build automation
├── .goreleaser.yml         # Release configuration
└── go.mod                  # Go module definition
```

## 🚀 Quick Start

### **Prerequisites**
- Go 1.23 or later
- Git
- Make (optional, for using Makefile)
- Access to Kavach backend instance

### **Installation Options**

#### **1. From Pre-built Releases**
```bash
# Linux AMD64
curl -L https://github.com/Gkemhcs/kavach-cli/releases/latest/download/kavach-cli_Linux_x86_64.tar.gz | tar -xz
sudo mv kavach /usr/local/bin/

# macOS AMD64
curl -L https://github.com/Gkemhcs/kavach-cli/releases/latest/download/kavach-cli_Darwin_x86_64.tar.gz | tar -xz
sudo mv kavach /usr/local/bin/

# Windows AMD64
# Download from GitHub releases and extract
```

#### **2. From Source**
```bash
# Clone the repository
git clone https://github.com/Gkemhcs/kavach-cli.git
cd kavach-cli

# Build and install
make build
sudo make install
```

#### **3. Using Package Managers**
```bash
# macOS (Homebrew)
brew install Gkemhcs/kavach-cli/kavach-cli

# Windows (Scoop)
scoop bucket add kavach-cli https://github.com/Gkemhcs/scoop-kavach-cli
scoop install kavach-cli
```

### **First-Time Setup**

1. **Verify Installation**
   ```bash
   kavach version
   ```

2. **Login to Backend**
   ```bash
   kavach login
   ```

3. **List Organizations**
   ```bash
   kavach org list
   ```

4. **Check Status**
   ```bash
   kavach status
   ```

## ⚙️ Configuration

### **Configuration Sources**
The CLI loads configuration from multiple sources in order of precedence:

1. **Command Line Flags**: Highest priority
2. **Environment Variables**: `KAVACH_*` prefixed variables
3. **Configuration Files**: `.kavach/config.yaml`
4. **Default Values**: Built-in sensible defaults

### **Environment Variables**

```bash
# Backend Configuration
export KAVACH_BACKEND_ENDPOINT="https://your-backend.com/api/v1/"
export KAVACH_DEVICE_CODE_URL="https://your-backend.com/api/v1/auth/device/code"
export KAVACH_DEVICE_TOKEN_URL="https://your-backend.com/api/v1/auth/device/token"

# Authentication
export KAVACH_TOKEN_FILE_PATH="$HOME/.kavach/credentials.json"
export KAVACH_LOG_DIR_PATH="$HOME/.kavach/"

# Logging Configuration
export KAVACH_LOG_LEVEL="info"
export KAVACH_LOG_MAX_SIZE="10"
export KAVACH_LOG_MAX_BACKUPS="5"
export KAVACH_LOG_MAX_AGE="30"
export KAVACH_LOG_COMPRESS="true"
```

### **Configuration File**
Create `~/.kavach/config.yaml`:

```yaml
backend:
  endpoint: "https://your-backend.com/api/v1/"
  device_code_url: "https://your-backend.com/api/v1/auth/device/code"
  device_token_url: "https://your-backend.com/api/v1/auth/device/token"

auth:
  token_file_path: "~/.kavach/credentials.json"
  auto_refresh: true

logging:
  level: "info"
  max_size: 10
  max_backups: 5
  max_age: 30
  compress: true

output:
  format: "table"
  color: true
  quiet: false
```

## 📚 Command Reference

### **Authentication Commands**

```bash
# Login to backend
kavach login

# Check authentication status
kavach status

# Logout from backend
kavach logout

# Refresh authentication tokens
kavach login --refresh
```

### **Organization Management**

```bash
# List organizations
kavach org list

# Create organization
kavach org create --name "MyOrg" --description "My Organization"

# Get organization details
kavach org get --name "MyOrg"

# Update organization
kavach org update --name "MyOrg" --description "Updated description"

# Delete organization
kavach org delete --name "MyOrg"
```

### **Secret Group Management**

```bash
# List secret groups
kavach secretgroup list --org "MyOrg"

# Create secret group
kavach secretgroup create --org "MyOrg" --name "api-secrets" --description "API Keys"

# Get secret group details
kavach secretgroup get --org "MyOrg" --name "api-secrets"

# Update secret group
kavach secretgroup update --org "MyOrg" --name "api-secrets" --description "Updated description"

# Delete secret group
kavach secretgroup delete --org "MyOrg" --name "api-secrets"
```

### **Environment Management**

```bash
# List environments
kavach environment list --org "MyOrg" --group "api-secrets"

# Create environment
kavach environment create --org "MyOrg" --group "api-secrets" --name "production"

# Get environment details
kavach environment get --org "MyOrg" --group "api-secrets" --name "production"

# Update environment
kavach environment update --org "MyOrg" --group "api-secrets" --name "production" --description "Production environment"

# Delete environment
kavach environment delete --org "MyOrg" --group "api-secrets" --name "production"
```

### **Secret Management**

```bash
# List secrets
kavach secret list --org "MyOrg" --group "api-secrets" --env "production"

# Create secret
kavach secret create --org "MyOrg" --group "api-secrets" --env "production" --name "API_KEY" --value "secret123"

# Get secret value
kavach secret get --org "MyOrg" --group "api-secrets" --env "production" --name "API_KEY"

# Update secret
kavach secret update --org "MyOrg" --group "api-secrets" --env "production" --name "API_KEY" --value "new-secret"

# Delete secret
kavach secret delete --org "MyOrg" --group "api-secrets" --env "production" --name "API_KEY"

# Create secret version
kavach secret version create --org "MyOrg" --group "api-secrets" --env "production" --message "Update API keys"

# List secret versions
kavach secret version list --org "MyOrg" --group "api-secrets" --env "production"

# Rollback to version
kavach secret version rollback --org "MyOrg" --group "api-secrets" --env "production" --version "abc12345"
```

### **Provider Operations**

```bash
# List providers
kavach provider list --org "MyOrg" --group "api-secrets" --env "production"

# Add provider credentials
kavach provider add --org "MyOrg" --group "api-secrets" --env "production" --provider "github" --config-file "github-config.json"

# Sync secrets to provider
kavach provider sync --org "MyOrg" --group "api-secrets" --env "production" --provider "github"

# Get provider status
kavach provider status --org "MyOrg" --group "api-secrets" --env "production" --provider "github"
```

### **User Group Management**

```bash
# List user groups
kavach user-group list --org "MyOrg"

# Create user group
kavach user-group create --org "MyOrg" --name "developers" --description "Development team"

# Add user to group
kavach user-group add-user --org "MyOrg" --group "developers" --username "john.doe"

# Remove user from group
kavach user-group remove-user --org "MyOrg" --group "developers" --username "john.doe"

# Delete user group
kavach user-group delete --org "MyOrg" --name "developers"
```

### **Global Options**

```bash
# Set output format
kavach --output json org list

# Enable debug logging
kavach --log-level debug org list

# Use specific configuration file
kavach --config /path/to/config.yaml org list

# Show help
kavach --help
kavach org --help
kavach org create --help
```

## 🧪 Development Setup

### **Prerequisites**
- Go 1.23 or later
- Git
- Make (optional)
- Access to Kavach backend instance

### **Local Development**

1. **Clone Repository**
   ```bash
   git clone https://github.com/Gkemhcs/kavach-cli.git
   cd kavach-cli
   ```

2. **Install Dependencies**
   ```bash
   go mod download
   go mod tidy
   ```

3. **Build CLI**
   ```bash
   make build
   ```

4. **Run Tests**
   ```bash
   make test
   make test-coverage
   ```

5. **Install Locally**
   ```bash
   make install
   ```

### **Available Make Commands**

```bash
make help          # Show all available commands
make build         # Build for current platform
make build-all     # Build for all supported platforms
make test          # Run tests
make test-race     # Run tests with race detection
make test-coverage # Run tests with coverage report
make lint          # Run linter
make format        # Format code
make install       # Install locally
make clean         # Clean build artifacts
make version       # Show version information
make docker-build  # Build Docker image
make snapshot      # Build snapshot release
```

## 🐳 Docker

### **Build Docker Image**
```bash
make docker-build
```

### **Run Container**
```bash
# Run with volume mount for configuration
docker run -it --rm \
  -v ~/.kavach:/root/.kavach \
  kavach-cli:latest version

# Run with environment variables
docker run -it --rm \
  -e KAVACH_BACKEND_ENDPOINT="https://your-backend.com/api/v1/" \
  kavach-cli:latest login
```

### **Development with Docker**
```bash
# Build development image
docker build -t kavach-cli:dev .

# Run with source code mount
docker run -it --rm \
  -v $(pwd):/app \
  -w /app \
  kavach-cli:dev
```

## 📊 Monitoring & Logging

### **Logging Configuration**
The CLI uses structured logging with configurable levels:

```bash
# Set log level
export KAVACH_LOG_LEVEL="debug"

# Configure log rotation
export KAVACH_LOG_MAX_SIZE="10"      # MB
export KAVACH_LOG_MAX_BACKUPS="5"
export KAVACH_LOG_MAX_AGE="30"       # Days
export KAVACH_LOG_COMPRESS="true"
```

### **Log Output Formats**
- **Console**: Human-readable output with colors
- **JSON**: Structured logging for automation
- **File**: Rotated log files with compression

### **Health Checks**
```bash
# Check CLI status
kavach status

# Check backend connectivity
kavach info

# Verify authentication
kavach login --check
```

## 🔒 Security Features

### **Authentication Security**
- **OAuth 2.0**: Industry-standard authentication protocol
- **JWT Tokens**: Secure token-based authentication
- **Token Refresh**: Automatic token renewal
- **Secure Storage**: Encrypted credential storage

### **Data Protection**
- **HTTPS Only**: All API communications use TLS
- **Token Encryption**: Sensitive tokens are encrypted at rest
- **Secure Headers**: Proper security headers in all requests
- **Input Validation**: Comprehensive input sanitization

### **Access Control**
- **Backend RBAC**: Leverages backend role-based access control
- **Permission Validation**: Client-side permission checking
- **Audit Logging**: All operations logged for compliance
- **Session Management**: Secure session handling

## 🚀 Release Management

### **Version Information**
```bash
# Show version details
kavach version

# Short version
kavach version --short

# JSON format
kavach version --json
```

### **Release Process**
The CLI uses GoReleaser for automated releases:

1. **Create Git Tag**
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

2. **Automated Build**
   - GitHub Actions builds for all platforms
   - Creates GitHub release with artifacts
   - Generates checksums and signatures

3. **Distribution**
   - Binary releases for Linux, macOS, Windows
   - Package manager support (Homebrew, Scoop)
   - Docker images

### **Supported Platforms**
| OS | Architecture | Status | Package Format |
|----|-------------|--------|----------------|
| Linux | AMD64 | ✅ | `.tar.gz` |
| Linux | ARM64 | ✅ | `.tar.gz` |
| macOS | AMD64 | ✅ | `.tar.gz` |
| macOS | ARM64 | ✅ | `.tar.gz` |
| Windows | AMD64 | ✅ | `.zip` |
| Windows | ARM64 | ❌ | N/A |

## 🔄 CI/CD Integration

### **GitHub Actions**
The repository includes comprehensive CI/CD workflows:

- **Release**: Automated releases on tag push
- **Snapshot**: Development builds on every push
- **Testing**: Comprehensive test suite execution
- **Linting**: Code quality and style checks
- **Security**: Dependency vulnerability scanning

### **Release Automation**
```yaml
# Example GitHub Actions workflow
name: Release
on:
  push:
    tags: ['v*']
jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
      - run: make release
```

## 🔍 Troubleshooting

### **Common Issues**

#### **Authentication Problems**
```bash
# Check token validity
kavach status

# Re-authenticate
kavach logout
kavach login

# Verify backend connectivity
kavach info
```

#### **Configuration Issues**
```bash
# Check configuration
kavach --config /path/to/config.yaml version

# Verify environment variables
env | grep KAVACH

# Reset configuration
rm -rf ~/.kavach/
kavach login
```

#### **Network Connectivity**
```bash
# Test backend connectivity
curl -v https://your-backend.com/healthz

# Check DNS resolution
nslookup your-backend.com

# Verify firewall settings
telnet your-backend.com 443
```

### **Debug Mode**
```bash
# Enable debug logging
export KAVACH_LOG_LEVEL="debug"
kavach org list

# Verbose output
kavach --verbose org list

# Check logs
tail -f ~/.kavach/kavach.log
```

### **Getting Help**
```bash
# Command help
kavach --help
kavach org --help
kavach org create --help

# Version information
kavach version --json

# Status check
kavach status
```

## 🤝 Contributing

### **Development Workflow**
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run `make test` and `make lint`
6. Submit a pull request

### **Development Guidelines**
- Follow Go coding standards and conventions
- Add comprehensive tests for new features
- Update documentation for new commands
- Use conventional commit messages
- Ensure all tests pass before submitting

### **Code Quality**
```bash
# Run quality checks
make lint
make format
make test
make test-coverage

# Check for security issues
go list -json -deps ./... | nancy sleuth
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.

## 🆘 Support

### **Documentation**
- [CLI Documentation](https://docs.kavach.gkem.cloud/cli)
- [API Reference](https://docs.kavach.gkem.cloud/api)
- [Examples](https://docs.kavach.gkem.cloud/examples)

### **Community**
- [GitHub Issues](https://github.com/Gkemhcs/kavach-cli/issues)
- [GitHub Discussions](https://github.com/Gkemhcs/kavach-cli/discussions)
- [Discord Server](https://discord.gg/kavach)

### **Getting Help**
- Check the troubleshooting section above
- Search existing issues and discussions
- Create a new issue with detailed information
- Join our community Discord for real-time support

---

**Kavach CLI** - Secure, powerful, and developer-friendly secrets management from the command line.
