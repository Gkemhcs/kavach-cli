package secret

import (
	"fmt"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/secret"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewRollbackCommand creates a new command for rolling back to a previous version
func NewRollbackCommand(logger *utils.Logger, cfg *config.Config, secretClient secret.SecretClient) *cobra.Command {
	var (
		versionID     string
		commitMessage string
		envName       string
		groupName     string
		orgName       string
	)

	cmd := &cobra.Command{
		Use:   "rollback",
		Short: "⏪ Rollback to a previous version",
		Long: `Rollback to a previous secret version.

This command creates a new version by copying all secrets from a previous version.
This is useful for reverting changes or restoring a known good state.

Key concepts:
- Creates a new version with all secrets from a target version
- Useful for reverting problematic changes
- Maintains audit trail with your commit message
- Preserves the original target version
- Creates a complete snapshot of the target version

Prerequisites:
- You must have appropriate permissions for the environment
- Target version must exist and be accessible
- Organization, secret group, and environment must be specified

The rollback process:
1. Validates the target version exists
2. Copies all secrets from the target version
3. Creates a new version with your commit message
4. Provides confirmation of successful rollback

Use cases:
- Reverting broken deployments
- Restoring from known good states
- Emergency rollbacks in production
- Testing rollback procedures
- Disaster recovery scenarios

Examples:
  kavach secret rollback --env "prod" --version "abc123" --message "Rollback to stable version"
  kavach secret rollback --env "staging" --version "def456" --message "Revert broken changes"
  kavach secret rollback --env "dev" --version "v1.2.3" --message "Restore working configuration"

Note: This creates a new version, not a replacement. The original target version
remains unchanged. Use 'kavach secret list' to see all versions.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Validate required flags
			if envName == "" {
				fmt.Println("\n🚨 Error: Environment ID is required. Please use --env to specify the environment.")
				logger.Warn("Environment ID not provided for secret rollback", map[string]interface{}{"cmd": "secret rollback"})
				return
			}

			if versionID == "" {
				fmt.Println("\n🚨 Error: Version ID is required. Please use --version to specify the target version.")
				logger.Warn("Version ID not provided for secret rollback", map[string]interface{}{"cmd": "secret rollback"})
				return
			}

			if commitMessage == "" {
				fmt.Println("\n🚨 Error: Commit message is required. Please use --message to specify a commit message.")
				logger.Warn("Commit message not provided for secret rollback", map[string]interface{}{"cmd": "secret rollback"})
				return
			}

			// Confirm rollback
			fmt.Printf("\n⚠️  You are about to rollback environment '%s' to version '%s'.\n", envName, versionID)
			fmt.Printf("   This will create a new version with all secrets from the target version.\n")
			fmt.Printf("   Commit message: %s\n", commitMessage)
			fmt.Print("\nDo you want to continue? (y/N): ")

			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Rollback cancelled.")
				logger.Info("Rollback cancelled by user", map[string]interface{}{"cmd": "secret rollback", "env": envName})
				return
			}

			// Perform rollback
			result, err := secretClient.RollbackToVersion(orgName, groupName, envName, versionID, commitMessage)
			if err != nil {
				// Handle specific CLI errors
				switch err {
				case cliErrors.ErrEnvironmentNotFound:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("Environment not found during rollback", map[string]interface{}{
						"cmd":            "secret rollback",
						"env":            envName,
						"target_version": versionID,
					})
					return
				case cliErrors.ErrTargetSecretVersionNotFound:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("Target secret version not found during rollback", map[string]interface{}{
						"cmd":            "secret rollback",
						"env":            envName,
						"target_version": versionID,
					})
					return
				case cliErrors.ErrEnvironmentsMisMatch:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("Environment mismatch during rollback", map[string]interface{}{
						"cmd":            "secret rollback",
						"env":            envName,
						"target_version": versionID,
					})
					return
				case cliErrors.ErrRollbackFailed:
					fmt.Printf("\n❌ %v\n", err)
					logger.Error("Rollback failed", err, map[string]interface{}{
						"cmd":            "secret rollback",
						"env":            envName,
						"target_version": versionID,
					})
					return
				case cliErrors.ErrCopySecretsFailed:
					fmt.Printf("\n❌ %v\n", err)
					logger.Error("Failed to copy secrets during rollback", err, map[string]interface{}{
						"cmd":            "secret rollback",
						"env":            envName,
						"target_version": versionID,
					})
					return
				case cliErrors.ErrEncryptionFailed:
					fmt.Printf("\n❌ %v\n", err)
					logger.Error("Encryption failed during rollback", err, map[string]interface{}{
						"cmd":            "secret rollback",
						"env":            envName,
						"target_version": versionID,
					})
					return
				case cliErrors.ErrInternalServer:
					fmt.Printf("\n❌ %v\n", err)
					logger.Error("Internal server error during rollback", err, map[string]interface{}{
						"cmd":            "secret rollback",
						"env":            envName,
						"target_version": versionID,
					})
					return
				case cliErrors.ErrNotLoggedIn:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("User not logged in during rollback", map[string]interface{}{
						"cmd":            "secret rollback",
						"env":            envName,
						"target_version": versionID,
					})
					return
				case cliErrors.ErrAccessDenied:
					fmt.Printf("\n%s\n", err.Error())
					logger.Warn("Access denied during secret rollback", map[string]interface{}{
						"cmd":            "secret rollback",
						"env":            envName,
						"target_version": versionID,
					})
					return
				default:
					fmt.Printf("\n❌ Failed to rollback to version: %v\n", err)
					logger.Error("Failed to rollback to version", err, map[string]interface{}{
						"cmd":            "secret rollback",
						"env":            envName,
						"target_version": versionID,
					})
					return
				}
			}

			fmt.Printf("\n✅ Successfully rolled back to version!\n")
			fmt.Printf("   New Version ID: %s\n", result.ID)
			fmt.Printf("   Environment: %s\n", result.EnvironmentID)
			fmt.Printf("   Target Version: %s\n", versionID)
			fmt.Printf("   Commit Message: %s\n", result.CommitMessage)
			fmt.Printf("   Secret Count: %d\n", result.SecretCount)
			fmt.Printf("   Created At: %s\n", result.CreatedAt)

			logger.Info("Successfully rolled back to version", map[string]interface{}{
				"cmd":            "secret rollback",
				"env":            envName,
				"target_version": versionID,
				"new_version":    result.ID,
				"secret_count":   result.SecretCount,
			})
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&versionID, "version", "v", "", "Target version ID to rollback to")
	cmd.Flags().StringVarP(&commitMessage, "message", "m", "", "Commit message for the rollback")
	cmd.Flags().StringVarP(&envName, "env", "e", "", "environment name")
	cmd.Flags().StringVarP(&groupName, "group", "g", "", "secret group name")
	cmd.Flags().StringVarP(&orgName, "org", "o", "", "organization name")
	// Mark required flags
	cmd.MarkFlagRequired("env")
	cmd.MarkFlagRequired("version")
	cmd.MarkFlagRequired("message")

	return cmd
}
