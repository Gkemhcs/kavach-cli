package cmd

import (
	"fmt"
	"runtime"
	"time"

	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// Build information variables set during the build process.
var (
	// Version is the semantic version of the application
	Version = "0.1.0"

	// BuildTime is the time when the binary was built
	BuildTime = "unknown"

	// GitCommit is the git commit hash
	GitCommit = "unknown"

	// GitBranch is the git branch name
	GitBranch = "unknown"

	// GoVersion is the Go version used to build the binary
	GoVersion = runtime.Version()

	// Platform is the target platform (OS/Arch)
	Platform = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
)

// BuildInfo returns a formatted string with build information for the CLI.
func BuildInfo() string {
	return fmt.Sprintf(`🛡️ Kavach CLI Build Information:
  Version:     %s
  Build Time:  %s
  Git Commit:  %s
  Git Branch:  %s
  Go Version:  %s
  Platform:    %s`,
		Version, BuildTime, GitCommit, GitBranch, GoVersion, Platform)
}

// GetVersion returns the current version string.
func GetVersion() string {
	return Version
}

// GetBuildTime returns the build time as a time.Time.
func GetBuildTime() time.Time {
	if BuildTime == "unknown" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, BuildTime)
	if err != nil {
		return time.Time{}
	}
	return t
}

// ProjectInfo returns a formatted string with project information.
func ProjectInfo() string {
	return fmt.Sprintf(`🛡️ Kavach - Enterprise Secret Management Platform

Project Summary:
  Kavach is a comprehensive secret management solution designed for enterprise
  environments. It provides secure storage, synchronization, and distribution
  of secrets across multiple cloud providers and environments.

Key Features:
  • Secure secret storage with encryption at rest
  • Multi-cloud provider synchronization (AWS, GCP, Azure)
  • Role-based access control (RBAC)
  • Organization and group management
  • Environment-specific secret deployment
  • CLI and REST API interfaces
  • Compliance and security best practices

Architecture:
  • Backend: Go-based REST API with Gin framework
  • CLI: Command-line interface for automation
  • Storage: Encrypted secret storage
  • Authentication: JWT-based with RBAC
  • Monitoring: Structured logging with Logrus

Technology Stack:
  • Language: Go %s
  • Framework: Gin (Backend), Cobra (CLI)
  • Database: Configurable (PostgreSQL, etc.)
  • Security: AES-256 encryption, JWT tokens
  • Platform: %s

Project Links:
  • Repository: https://github.com/Gkemhcs/kavach-backend
  • Documentation: https://github.com/Gkemhcs/kavach-backend/docs
  • Issues: https://github.com/Gkemhcs/kavach-backend/issues

Current Version: %s
Build Date: %s`, GoVersion, Platform, Version, BuildTime)
}

// infoCmd is the Cobra command for displaying project information.
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display project information and summary",
	Long: `Display comprehensive information about the Kavach project, its features,
architecture, and technology stack.

This command provides:
  • Project overview and purpose
  • Key features and capabilities
  • Technology stack and architecture
  • Project links and documentation
  • Current version and build information

Examples:
  kavach info
  kavach info --json`,
	Run: func(cmd *cobra.Command, args []string) {
		// Logger is available in context, but info is user-facing so we use fmt for output.
		logger := cmd.Context().Value(CtxLoggerKey).(*utils.Logger)
		defer logger.Close()
		json, _ := cmd.Flags().GetBool("json")
		if json {
			jsonOutput := fmt.Sprintf(`{
  "project": {
    "name": "Kavach",
    "description": "Enterprise Secret Management Platform",
    "version": "%s",
    "buildDate": "%s"
  },
  "summary": {
    "purpose": "Comprehensive secret management solution for enterprise environments",
    "keyFeatures": [
      "Secure secret storage with encryption at rest",
      "Multi-cloud provider synchronization",
      "Role-based access control (RBAC)",
      "Organization and group management",
      "Environment-specific secret deployment",
      "CLI and REST API interfaces"
    ]
  },
  "architecture": {
    "backend": "Go-based REST API with Gin framework",
    "cli": "Command-line interface with Cobra",
    "security": "AES-256 encryption, JWT tokens",
    "storage": "Encrypted secret storage"
  },
  "technology": {
    "language": "Go",
    "goVersion": "%s",
    "platform": "%s",
    "frameworks": ["Gin", "Cobra", "Logrus"]
  },
  "links": {
    "repository": "https://github.com/Gkemhcs/kavach-backend",
    "documentation": "https://github.com/Gkemhcs/kavach-backend/docs",
    "issues": "https://github.com/Gkemhcs/kavach-backend/issues"
  }
}`, Version, BuildTime, GoVersion, Platform)
			fmt.Println(jsonOutput)
			logger.Info("Displayed project info in JSON format", map[string]interface{}{"cmd": "info", "json": true})
			return
		}
		fmt.Println(ProjectInfo())
		logger.Info("Displayed project info in text format", map[string]interface{}{"cmd": "info", "json": false})
	},
}

// init registers the info command with the root command.
func init() {
	rootCmd.AddCommand(infoCmd)
	// Add JSON output flag
	infoCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
}
