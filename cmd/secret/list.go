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

// GetListVersionsHeaders returns the headers for the secret versions list table.
func GetListVersionsHeaders() []string {
	return []string{
		"Version ID",
		"Commit Message",
		"Created At",
		"Secret Count",
	}
}

// ToRenderableVersions converts a list of secret versions to a 2D string slice for table rendering.
func ToRenderableVersions(versions []types.SecretVersionResponse) [][]string {
	var out [][]string
	for _, version := range versions {
		// Truncate commit message if too long
		commitMsg := version.CommitMessage
		if len(commitMsg) > 18 {
			commitMsg = commitMsg[:15] + "..."
		}

		// Format timestamp
		createdAt := "N/A"
		if version.CreatedAt != "" {
			if t, err := time.Parse(time.RFC3339, version.CreatedAt); err == nil {
				createdAt = t.Format("2006-01-02 15:04:05")
			}
		}

		out = append(out, []string{
			version.ID,
			commitMsg,
			createdAt,
			fmt.Sprintf("%d", version.SecretCount),
		})
	}
	return out
}

// NewListVersionsCommand creates a new command for listing secret versions
func NewListVersionsCommand(logger *utils.Logger, cfg *config.Config, secretClient secret.SecretClient) *cobra.Command {
	var groupName string
	var orgName string
	var envName string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "📋 List secret versions",
		Long: `List all secret versions for an environment.

This command displays a table of all secret versions with their details including
version ID, commit message, creation date, and secret count.

Key concepts:
- Shows all versions of secrets in the environment
- Displays version metadata without exposing actual values
- Useful for version history and audit trails
- Shows commit messages for each version
- Helps track changes over time

The output includes:
- Version ID: Unique identifier for the version
- Commit Message: The message provided when creating the version
- Created At: When the version was created
- Secret Count: Number of secrets in this version
- Status: Current status of the version

Use cases:
- Auditing secret version history
- Tracking changes over time
- Finding specific versions by commit message
- Understanding when secrets were updated
- Planning rollbacks to previous versions

Examples:
  kavach secret list --env "prod"
  kavach secret list --env "staging"
  kavach secret list --env "dev"

Note: This shows version metadata, not the actual secret values.
Use 'kavach secret details <version-id>' to see detailed information
about a specific version.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Validate required flags
			if envName == "" {
				fmt.Println("\n🚨 Error: Environment ID is required. Please use --env to specify the environment.")
				logger.Warn("Environment ID not provided for secret list", map[string]interface{}{"cmd": "secret list"})
				return
			}

			// Get versions
			versions, err := secretClient.ListVersions(orgName, groupName, envName)
			if err != nil {
				// Handle specific CLI errors
				switch err {
				case cliErrors.ErrEnvironmentNotFound:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("Environment not found during list", map[string]interface{}{
						"cmd": "secret list",
						"env": envName,
					})
					return
				case cliErrors.ErrSecretGroupNotFound:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("Secret group not found during list", map[string]interface{}{
						"cmd": "secret list",
						"env": envName,
					})
					return
				case cliErrors.ErrOrganizationNotFound:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("Organization not found during list", map[string]interface{}{
						"cmd": "secret list",
						"env": envName,
					})
					return
				case cliErrors.ErrInternalServer:
					fmt.Printf("\n❌ %v\n", err)
					logger.Error("Internal server error during list", err, map[string]interface{}{
						"cmd": "secret list",
						"env": envName,
					})
					return
				case cliErrors.ErrNotLoggedIn:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("User not logged in during list", map[string]interface{}{
						"cmd": "secret list",
						"env": envName,
					})
					return
				case cliErrors.ErrAccessDenied:
					fmt.Printf("\n%s\n", err.Error())
					logger.Warn("Access denied during secret list", map[string]interface{}{
						"cmd": "secret list",
						"env": envName,
					})
					return
				default:
					fmt.Printf("\n❌ Failed to list secret versions: %v\n", err)
					logger.Error("Failed to list secret versions", err, map[string]interface{}{
						"cmd": "secret list",
						"env": envName,
					})
					return
				}
			}

			if len(versions) == 0 {
				fmt.Printf("\n📝 No secret versions found for environment '%s'.\n", envName)
				logger.Info("No secret versions found", map[string]interface{}{"cmd": "secret list", "env": envName})
				return
			}

			// Display versions using the shared table rendering utility
			fmt.Printf("\n📋 Secret Versions for Environment: %s\n", envName)
			utils.RenderTable(GetListVersionsHeaders(), ToRenderableVersions(versions))
			fmt.Printf("Total versions: %d\n", len(versions))

			logger.Info("Secret versions listed successfully", map[string]interface{}{
				"cmd":   "secret list",
				"env":   envName,
				"count": len(versions),
			})
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&envName, "env", "e", "", "environment name")
	cmd.Flags().StringVarP(&groupName, "group", "g", "", "secret group name")
	cmd.Flags().StringVarP(&orgName, "org", "o", "", "organization name")
	// Mark required flags
	cmd.MarkFlagRequired("env")

	return cmd
}
