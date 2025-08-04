package secret

import (
	"fmt"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	"github.com/Gkemhcs/kavach-cli/internal/secret"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewAddSecretCommand creates a new command for adding secrets to staging
func NewAddSecretCommand(logger *utils.Logger, cfg *config.Config, secretClient secret.SecretClient) *cobra.Command {
	var (
		secretName  string
		secretValue string
	)

	cmd := &cobra.Command{
		Use:   "add",
		Short: "➕ Add a secret to staging",
		Long: `Add a secret to the local staging area for batch operations.

This command adds a secret to the local staging area without immediately
committing it to the environment. This allows you to stage multiple secrets
and commit them all at once using the 'commit' command.

Key concepts:
- Secrets are added to a local staging area first
- Multiple secrets can be staged before committing
- Staging allows for batch operations and review
- Secrets are not available in the environment until committed
- Staging area is local to your current session

Available input methods:
- --name: The name of the secret
- --value: Direct string value for simple secrets

Use cases:
- Staging multiple secrets for a new deployment
- Reviewing secrets before committing to environment
- Batch operations for multiple configuration changes
- Testing secret formats before committing

Examples:
  kavach secret add --name "db-password" --value "mypassword123"
  kavach secret add --name "api-key" --value "sk-1234567890abcdef"
  kavach secret add --name "debug" --value "true"

Note: After adding secrets to staging, use 'kavach secret commit' to
commit them to the environment. Use 'kavach secret list' to see
currently staged secrets.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Validate required flags
			if secretName == "" {
				fmt.Println("\n🚨 Error: Secret name is required. Please use --name to specify the secret name.")
				logger.Warn("Secret name not provided for secret add", map[string]interface{}{"cmd": "secret add"})
				return
			}

			if secretValue == "" {
				fmt.Println("\n🚨 Error: Secret value is required. Please use --value to specify the secret value.")
				logger.Warn("Secret value not provided for secret add", map[string]interface{}{"cmd": "secret add"})
				return
			}

			// Add secret to staging area
			stagingService := secret.NewStagingService(logger)
			err := stagingService.AddSecretToStaging(secretName, secretValue)
			if err != nil {
				fmt.Printf("\n❌ Failed to add secret to staging: %v\n", err)
				logger.Error("Failed to add secret to staging", err, map[string]interface{}{
					"cmd":    "secret add",
					"secret": secretName,
				})
				return
			}

			fmt.Printf("\n✅ Successfully added secret '%s' to staging!\n", secretName)
			fmt.Printf("   Secret Name:  %s\n", secretName)
			fmt.Printf("\n💡 Use 'kavach secret commit --org <org> --group <group> --env <env> --message \"your message\"' to commit staged secrets.\n")

			logger.Info("Secret added to staging successfully", map[string]interface{}{
				"cmd":    "secret add",
				"secret": secretName,
			})
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&secretName, "name", "n", "", "Secret name")
	cmd.Flags().StringVarP(&secretValue, "value", "v", "", "Secret value")

	// Mark required flags
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("value")

	return cmd
}
