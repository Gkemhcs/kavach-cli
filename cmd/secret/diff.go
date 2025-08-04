package secret

import (
	"fmt"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/secret"
	"github.com/Gkemhcs/kavach-cli/internal/types"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// GetDiffHeaders returns the headers for the secret diff table.
func GetDiffHeaders() []string {
	return []string{
		"Secret Name",
		"Change Type",
		"Details",
	}
}

// ToRenderable converts a list of secret diff changes to a 2D string slice for table rendering.
func ToRenderable(changes []types.SecretDiffChange) [][]string {
	var out [][]string
	for _, change := range changes {
		var changeType string
		var details string

		switch change.Type {
		case "added":
			changeType = "➕ Added"
			details = fmt.Sprintf("New value: %s", truncateValue(change.NewValue, 55))
		case "removed":
			changeType = "❌ Removed"
			details = fmt.Sprintf("Old value: %s", truncateValue(change.OldValue, 55))
		case "modified":
			changeType = "🔄 Modified"
			details = fmt.Sprintf("Old: %s → New: %s", truncateValue(change.OldValue, 25), truncateValue(change.NewValue, 25))
		case "no_change":
			changeType = "✅ No Change"
			details = fmt.Sprintf("Value unchanged: %s", truncateValue(change.OldValue, 55))
		default:
			changeType = "❓ Unknown"
			details = "Unknown change type"
		}

		out = append(out, []string{change.Name, changeType, details})
	}
	return out
}

// NewDiffCommand creates a new command for showing differences between versions
func NewDiffCommand(logger *utils.Logger, cfg *config.Config, secretClient secret.SecretClient) *cobra.Command {
	var (
		fromVersion string
		toVersion   string
		orgName     string
		envName     string
		groupName   string
	)

	cmd := &cobra.Command{
		Use:   "diff",
		Short: "🔍 Show differences between versions",
		Long: `Show differences between two secret versions.

This command displays a detailed comparison of secrets between two versions,
showing what was added, modified, or removed.

Key concepts:
- Compares secrets between two specific versions
- Shows detailed change information for each secret
- Displays change types: added, removed, modified, no change
- Provides summary statistics of changes
- Helps understand what changed between versions

The output includes:
- Secret Name: The name of the secret
- Change Type: Added, Removed, Modified, or No Change
- Details: Specific changes or values (truncated for security)

Prerequisites:
- Both source and target versions must exist
- You must have appropriate permissions for the environment
- Organization, secret group, and environment must be specified

Use cases:
- Auditing changes between deployments
- Understanding what changed in a version
- Reviewing changes before deployment
- Troubleshooting configuration issues
- Compliance and audit requirements

Examples:
  kavach secret diff --from "abc123" --to "def456"
  kavach secret diff --from "v1.0" --to "v1.1"
  kavach secret diff --from "prod-v1" --to "prod-v2"

Note: Secret values are truncated for security. Use 'kavach secret details <version-id>'
to see full secret values for a specific version.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Validate required flags
			if fromVersion == "" {
				fmt.Println("\n🚨 Error: From version is required. Please use --from to specify the source version.")
				logger.Warn("From version not provided for secret diff", map[string]interface{}{"cmd": "secret diff"})
				return
			}

			if toVersion == "" {
				fmt.Println("\n🚨 Error: To version is required. Please use --to to specify the target version.")
				logger.Warn("To version not provided for secret diff", map[string]interface{}{"cmd": "secret diff"})
				return
			}

			// Get diff
			diff, err := secretClient.GetVersionDiff(orgName, groupName, envName, fromVersion, toVersion)
			if err != nil {
				// Handle specific CLI errors
				switch err {
				case cliErrors.ErrSecretVersionNotFound:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("Secret version not found during diff", map[string]interface{}{
						"cmd":  "secret diff",
						"from": fromVersion,
						"to":   toVersion,
					})
					return
				case cliErrors.ErrDecryptionFailed:
					fmt.Printf("\n❌ %v\n", err)
					logger.Error("Decryption failed during diff", err, map[string]interface{}{
						"cmd":  "secret diff",
						"from": fromVersion,
						"to":   toVersion,
					})
					return
				case cliErrors.ErrInternalServer:
					fmt.Printf("\n❌ %v\n", err)
					logger.Error("Internal server error during diff", err, map[string]interface{}{
						"cmd":  "secret diff",
						"from": fromVersion,
						"to":   toVersion,
					})
					return
				case cliErrors.ErrNotLoggedIn:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("User not logged in during diff", map[string]interface{}{
						"cmd":  "secret diff",
						"from": fromVersion,
						"to":   toVersion,
					})
					return
				case cliErrors.ErrAccessDenied:
					fmt.Printf("\n%s\n", err.Error())
					logger.Warn("Access denied during secret diff", map[string]interface{}{
						"cmd":  "secret diff",
						"from": fromVersion,
						"to":   toVersion,
					})
					return
				default:
					fmt.Printf("\n❌ Failed to get version diff: %v\n", err)
					logger.Error("Failed to get version diff", err, map[string]interface{}{
						"cmd":  "secret diff",
						"from": fromVersion,
						"to":   toVersion,
					})
					return
				}
			}

			// Display diff information
			fmt.Printf("\n📊 Version Diff\n")
			fmt.Printf("From Version: %s\n", diff.FromVersion)
			fmt.Printf("To Version:   %s\n", diff.ToVersion)
			fmt.Printf("Total Changes: %d\n\n", len(diff.Changes))

			if len(diff.Changes) == 0 {
				fmt.Println("📝 No differences found between the versions.")
				logger.Info("No differences found between versions", map[string]interface{}{
					"cmd":  "secret diff",
					"from": fromVersion,
					"to":   toVersion,
				})
				return
			}

			// Display changes using the shared table rendering utility
			utils.RenderTable(GetDiffHeaders(), ToRenderable(diff.Changes))

			// Summary
			added := 0
			removed := 0
			modified := 0
			noChange := 0
			for _, change := range diff.Changes {
				switch change.Type {
				case "added":
					added++
				case "removed":
					removed++
				case "modified":
					modified++
				case "no_change":
					noChange++
				}
			}

			fmt.Printf("\n📈 Summary:\n")
			fmt.Printf("   ➕ Added:    %d\n", added)
			fmt.Printf("   ❌ Removed:  %d\n", removed)
			fmt.Printf("   🔄 Modified: %d\n", modified)
			fmt.Printf("   ✅ No Change: %d\n", noChange)

			logger.Info("Version diff retrieved successfully", map[string]interface{}{
				"cmd":       "secret diff",
				"from":      fromVersion,
				"to":        toVersion,
				"changes":   len(diff.Changes),
				"added":     added,
				"removed":   removed,
				"modified":  modified,
				"no_change": noChange,
			})
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&fromVersion, "from", "f", "", "Source version ID")
	cmd.Flags().StringVarP(&toVersion, "to", "t", "", "Target version ID")
	cmd.Flags().StringVarP(&envName, "env", "e", "", "environment name")
	cmd.Flags().StringVarP(&groupName, "group", "g", "", "secret group name")
	cmd.Flags().StringVarP(&orgName, "org", "o", "", "organization name")
	// Mark required flags
	cmd.MarkFlagRequired("from")
	cmd.MarkFlagRequired("to")

	return cmd
}

// truncateValue truncates a value to the specified length and adds ellipsis if needed
func truncateValue(value string, maxLength int) string {
	if len(value) <= maxLength {
		return value
	}
	return value[:maxLength-3] + "..."
}
