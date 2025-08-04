package secret

import (
	"fmt"
	"time"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/secret"
	"github.com/Gkemhcs/kavach-cli/internal/types"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// GetVersionDetailsHeaders returns the headers for the version details table.
func GetVersionDetailsHeaders() []string {
	return []string{
		"Field",
		"Value",
	}
}

// GetSecretsHeaders returns the headers for the secrets table.
func GetSecretsHeaders() []string {
	return []string{
		"Name",
		"Value",
	}
}

// ToRenderableDetails converts version details to a 2D string slice for table rendering.
func ToRenderableDetails(version *types.SecretVersionDetailResponse) [][]string {
	// Format timestamp
	createdAt := "N/A"
	if version.CreatedAt != "" {
		if t, err := time.Parse(time.RFC3339, version.CreatedAt); err == nil {
			createdAt = t.Format("2006-01-02 15:04:05")
		}
	}

	return [][]string{
		{"Version ID", version.ID},
		{"Environment ID", version.EnvironmentID},
		{"Commit Message", version.CommitMessage},
		{"Created At", createdAt},
		{"Secret Count", fmt.Sprintf("%d", len(version.Secrets))},
	}
}

// ToRenderableSecrets converts a list of secrets to a 2D string slice for table rendering.
func ToRenderableSecrets(secrets []types.Secret) [][]string {
	var out [][]string
	for _, secret := range secrets {
		// Truncate value if too long
		value := secret.Value
		if len(value) > 59 {
			value = value[:56] + "..."
		}
		out = append(out, []string{secret.Name, value})
	}
	return out
}

// NewGetVersionDetailsCommand creates a new command for getting version details
func NewGetVersionDetailsCommand(logger *utils.Logger, cfg *config.Config, secretClient secret.SecretClient) *cobra.Command {
	var versionID string
	var envName string
	var orgName string
	var groupName string

	cmd := &cobra.Command{
		Use:   "details",
		Short: "📋 Get version details",
		Long: `Get detailed information about a specific secret version.

This command displays all information about a secret version including
version metadata and all secrets with their values.

Key concepts:
- Shows complete information about a specific version
- Displays version metadata and all secret values
- Useful for auditing and verification purposes
- Shows full secret values (not truncated)
- Provides comprehensive version information

The output includes:
- Version ID: Unique identifier for the version
- Environment ID: The environment this version belongs to
- Commit Message: The message provided when creating the version
- Created At: When the version was created
- Secret Count: Number of secrets in this version
- All Secrets: Complete list with names and values

Prerequisites:
- Version must exist and be accessible
- You must have appropriate permissions for the environment
- Organization, secret group, and environment must be specified

Use cases:
- Auditing specific versions
- Verifying secret values
- Troubleshooting configuration issues
- Compliance and security reviews
- Debugging deployment problems

Examples:
  kavach secret details --version "abc123"
  kavach secret details --version "def456"
  kavach secret details --version "prod-v1.2.3"

Note: This shows the complete secret values. Use 'kavach secret list' to see
version metadata without exposing secret values.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Validate required flags
			if versionID == "" {
				fmt.Println("\n🚨 Error: Version ID is required. Please use --version to specify the version.")
				logger.Warn("Version ID not provided for secret details", map[string]interface{}{"cmd": "secret details"})
				return
			}

			// Get version details
			version, err := secretClient.GetVersionDetails(orgName, groupName, envName, versionID)
			if err != nil {
				// Handle specific CLI errors
				switch err {
				case cliErrors.ErrSecretVersionNotFound:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("Secret version not found during details", map[string]interface{}{
						"cmd":     "secret details",
						"version": versionID,
					})
					return
				case cliErrors.ErrSecretNotFound:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("Secret not found during details", map[string]interface{}{
						"cmd":     "secret details",
						"version": versionID,
					})
					return
				case cliErrors.ErrDecryptionFailed:
					fmt.Printf("\n❌ %v\n", err)
					logger.Error("Decryption failed during details", err, map[string]interface{}{
						"cmd":     "secret details",
						"version": versionID,
					})
					return
				case cliErrors.ErrInternalServer:
					fmt.Printf("\n❌ %v\n", err)
					logger.Error("Internal server error during details", err, map[string]interface{}{
						"cmd":     "secret details",
						"version": versionID,
					})
					return
				case cliErrors.ErrNotLoggedIn:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("User not logged in during details", map[string]interface{}{
						"cmd":     "secret details",
						"version": versionID,
					})
					return
				case cliErrors.ErrAccessDenied:
					fmt.Printf("\n%s\n", err.Error())
					logger.Warn("Access denied during secret details", map[string]interface{}{
						"cmd":     "secret details",
						"version": versionID,
					})
					return
				default:
					fmt.Printf("\n❌ Failed to get version details: %v\n", err)
					logger.Error("Failed to get version details", err, map[string]interface{}{
						"cmd":     "secret details",
						"version": versionID,
					})
					return
				}
			}

			// Display version information using the shared table rendering utility
			fmt.Printf("\n📋 Version Details\n")
			utils.RenderTable(GetVersionDetailsHeaders(), ToRenderableDetails(version))

			// Display secrets
			if len(version.Secrets) > 0 {
				fmt.Printf("\n🔐 Secrets\n")
				utils.RenderTable(GetSecretsHeaders(), ToRenderableSecrets(version.Secrets))
			} else {
				fmt.Printf("\n📝 No secrets found in this version.\n")
			}

			logger.Info("Version details retrieved successfully", map[string]interface{}{
				"cmd":          "secret details",
				"version":      versionID,
				"secret_count": len(version.Secrets),
			})
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&versionID, "version", "v", "", "Version ID to get details for")
	cmd.Flags().StringVarP(&envName, "env", "e", "", "environment name")
	cmd.Flags().StringVarP(&groupName, "group", "g", "", "secret group name")
	cmd.Flags().StringVarP(&orgName, "org", "o", "", "organization name")

	// Mark required flags
	cmd.MarkFlagRequired("version")

	return cmd
}
