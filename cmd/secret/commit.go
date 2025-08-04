package secret

import (
	"fmt"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/secret"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewCommitSecretCommand creates a new command for committing staged secrets
func NewCommitSecretCommand(logger *utils.Logger, cfg *config.Config, secretClient secret.SecretClient) *cobra.Command {
	var (
		orgName       string
		groupName     string
		envName       string
		commitMessage string
	)

	cmd := &cobra.Command{
		Use:   "commit",
		Short: "💾 Commit staged secrets to the environment",
		Long: `Commit all staged secrets to the current environment.

This command takes all secrets that have been added to the staging area
and commits them to the current environment. After committing, the secrets
will be available in the environment and can be accessed by applications.

Key concepts:
- Commits all secrets from the local staging area to the environment
- Creates a new version of each secret in the environment
- Staging area is cleared after successful commit
- Secrets become immediately available to applications
- Commit creates an audit trail with your message

Prerequisites:
- You must have secrets staged using 'kavach secret add'
- You must have appropriate permissions for the environment
- Organization, secret group, and environment must be specified

The commit process:
1. Validates all staged secrets
2. Creates new versions in the environment
3. Clears the local staging area
4. Provides confirmation of successful commit

Use cases:
- Deploying new configuration to production
- Updating multiple secrets at once
- Creating audit trails for secret changes
- Batch operations for environment updates

Examples:
  kavach secret commit --org "my-org" --group "backend" --env "prod" --message "Add API keys for production"
  kavach secret commit --org "my-org" --group "frontend" --env "staging" --message "Update database configuration"
  kavach secret commit --org "my-org" --group "backend" --env "staging" --message "Deploy new config"

Note: After committing, the staging area is cleared. Use 'kavach secret list'
to verify the secrets are available in the environment.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Validate required flags
			if orgName == "" {
				fmt.Println("\n🚨 Error: Organization name is required. Please use --org to specify the organization.")
				logger.Warn("Organization name not provided for secret commit", map[string]interface{}{"cmd": "secret commit"})
				return
			}

			if groupName == "" {
				fmt.Println("\n🚨 Error: Secret group name is required. Please use --group to specify the secret group.")
				logger.Warn("Secret group name not provided for secret commit", map[string]interface{}{"cmd": "secret commit"})
				return
			}

			if envName == "" {
				fmt.Println("\n🚨 Error: Environment name is required. Please use --env to specify the environment.")
				logger.Warn("Environment name not provided for secret commit", map[string]interface{}{"cmd": "secret commit"})
				return
			}

			if commitMessage == "" {
				fmt.Println("\n🚨 Error: Commit message is required. Please use --message to specify a commit message.")
				logger.Warn("Commit message not provided for secret commit", map[string]interface{}{"cmd": "secret commit"})
				return
			}

			// Get staged secrets
			stagingService := secret.NewStagingService(logger)
			staged, err := stagingService.GetStagedSecrets()
			if err != nil {
				fmt.Printf("\n❌ Failed to get staged secrets: %v\n", err)
				logger.Error("Failed to get staged secrets", err, map[string]interface{}{"cmd": "secret commit"})
				return
			}

			// Check if there are any staged secrets
			if len(staged.Secrets) == 0 {
				fmt.Println("\n📝 No secrets in staging area. Use 'kavach secret add' to add secrets first.")
				logger.Info("No secrets in staging area", map[string]interface{}{"cmd": "secret commit"})
				return
			}

			// Show what will be committed
			fmt.Printf("\n📋 Committing %d secrets:\n", len(staged.Secrets))
			fmt.Printf("   Organization: %s\n", orgName)
			fmt.Printf("   Secret Group: %s\n", groupName)
			fmt.Printf("   Environment:  %s\n", envName)
			fmt.Printf("   Commit Message: %s\n", commitMessage)
			fmt.Println("\nSecrets to be committed:")
			for _, secret := range staged.Secrets {
				fmt.Printf("   • %s\n", secret.Name)
			}

			// Confirm commit
			fmt.Print("\nDo you want to proceed? (y/N): ")
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Commit cancelled.")
				logger.Info("Commit cancelled by user", map[string]interface{}{"cmd": "secret commit"})
				return
			}

			// Create the version using the client (which will handle ID resolution)
			result, err := secretClient.CreateVersion(
				orgName,
				groupName,
				envName,
				staged.Secrets,
				commitMessage,
			)
			if err != nil {
				// Handle specific CLI errors
				switch err {
				case cliErrors.ErrEmptySecrets:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("Empty secrets error during commit", map[string]interface{}{
						"cmd":   "secret commit",
						"org":   orgName,
						"group": groupName,
						"env":   envName,
					})
					return
				case cliErrors.ErrTooManySecrets:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("Too many secrets error during commit", map[string]interface{}{
						"cmd":   "secret commit",
						"org":   orgName,
						"group": groupName,
						"env":   envName,
					})
					return
				case cliErrors.ErrInvalidSecretName:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("Invalid secret name error during commit", map[string]interface{}{
						"cmd":   "secret commit",
						"org":   orgName,
						"group": groupName,
						"env":   envName,
					})
					return
				case cliErrors.ErrSecretValueTooLong:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("Secret value too long error during commit", map[string]interface{}{
						"cmd":   "secret commit",
						"org":   orgName,
						"group": groupName,
						"env":   envName,
					})
					return
				case cliErrors.ErrEncryptionFailed:
					fmt.Printf("\n❌ %v\n", err)
					logger.Error("Encryption failed during commit", err, map[string]interface{}{
						"cmd":   "secret commit",
						"org":   orgName,
						"group": groupName,
						"env":   envName,
					})
					return
				case cliErrors.ErrEnvironmentNotFound:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("Environment not found during commit", map[string]interface{}{
						"cmd":   "secret commit",
						"org":   orgName,
						"group": groupName,
						"env":   envName,
					})
					return
				case cliErrors.ErrSecretGroupNotFound:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("Secret group not found during commit", map[string]interface{}{
						"cmd":   "secret commit",
						"org":   orgName,
						"group": groupName,
						"env":   envName,
					})
					return
				case cliErrors.ErrOrganizationNotFound:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("Organization not found during commit", map[string]interface{}{
						"cmd":   "secret commit",
						"org":   orgName,
						"group": groupName,
						"env":   envName,
					})
					return
				case cliErrors.ErrNotLoggedIn:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("User not logged in during commit", map[string]interface{}{
						"cmd":   "secret commit",
						"org":   orgName,
						"group": groupName,
						"env":   envName,
					})
					return
				case cliErrors.ErrAccessDenied:
					fmt.Printf("\n%s\n", err.Error())
					logger.Warn("Access denied during secret commit", map[string]interface{}{
						"cmd":   "secret commit",
						"org":   orgName,
						"group": groupName,
						"env":   envName,
					})
					return
				default:
					fmt.Printf("\n❌ Failed to commit secrets: %v\n", err)
					logger.Error("Failed to commit secrets", err, map[string]interface{}{
						"cmd":   "secret commit",
						"org":   orgName,
						"group": groupName,
						"env":   envName,
					})
					return
				}
			}

			// Clear staging area
			if err := stagingService.ClearStaging(); err != nil {
				fmt.Printf("\n⚠️  Warning: Failed to clear staging area: %v\n", err)
				logger.Warn("Failed to clear staging area", map[string]interface{}{"cmd": "secret commit", "error": err.Error()})
			}

			fmt.Printf("\n✅ Successfully committed secrets!\n")
			fmt.Printf("   Version ID: %s\n", result.ID)
			fmt.Printf("   Environment: %s\n", result.EnvironmentID)
			fmt.Printf("   Commit Message: %s\n", result.CommitMessage)
			fmt.Printf("   Secret Count: %d\n", result.SecretCount)
			fmt.Printf("   Created At: %s\n", result.CreatedAt)

			logger.Info("Secrets committed successfully", map[string]interface{}{
				"cmd":          "secret commit",
				"version_id":   result.ID,
				"secret_count": result.SecretCount,
				"org":          orgName,
				"group":        groupName,
				"env":          envName,
			})
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&orgName, "org", "o", "", "Organization name")
	cmd.Flags().StringVarP(&groupName, "group", "g", "", "Secret group name")
	cmd.Flags().StringVarP(&envName, "env", "e", "", "Environment name")
	cmd.Flags().StringVarP(&commitMessage, "message", "m", "", "Commit message")

	// Mark required flags
	cmd.MarkFlagRequired("org")
	cmd.MarkFlagRequired("group")
	cmd.MarkFlagRequired("env")
	cmd.MarkFlagRequired("message")

	return cmd
}
