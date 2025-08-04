package secret

import (
	"fmt"
	"strings"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/secret"
	"github.com/Gkemhcs/kavach-cli/internal/types"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewCreateSecretCommand creates a new command for creating secret versions
func NewCreateSecretCommand(logger *utils.Logger, cfg *config.Config, secretClient secret.SecretClient) *cobra.Command {
	var (
		commitMessage string
		envName       string
		orgName       string
		groupName     string
		secretsInput  []string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new secret version",
		Long: `Create a new version of secrets for an environment.

This command allows you to create a new version containing multiple secrets with a commit message.
Secrets should be provided in the format "name=value" and can be specified multiple times.

Examples:
  kavach secret create --env "prod" --message "Add API keys" --secret "API_KEY=abc123" --secret "DB_URL=postgres://..."
  kavach secret create --env "staging" --message "Update config" --secret "DEBUG=true" --secret "PORT=8080"`,
		Run: func(cmd *cobra.Command, args []string) {
			// Validate required flags
			if envName == "" {
				fmt.Println("\n🚨 Error: Environment ID is required. Please use --env to specify the environment.")
				logger.Warn("Environment ID not provided for secret create", map[string]interface{}{"cmd": "secret create"})
				return
			}

			if commitMessage == "" {
				fmt.Println("\n🚨 Error: Commit message is required. Please use --message to specify a commit message.")
				logger.Warn("Commit message not provided for secret create", map[string]interface{}{"cmd": "secret create"})
				return
			}

			if len(secretsInput) == 0 {
				fmt.Println("\n🚨 Error: At least one secret is required. Please use --secret to specify secrets.")
				logger.Warn("No secrets provided for secret create", map[string]interface{}{"cmd": "secret create"})
				return
			}

			// Parse secrets from input
			secrets := make([]types.Secret, 0, len(secretsInput))
			for _, secretInput := range secretsInput {
				parts := strings.SplitN(secretInput, "=", 2)
				if len(parts) != 2 {
					fmt.Printf("\n🚨 Error: Invalid secret format '%s'. Use format 'name=value'.\n", secretInput)
					logger.Warn("Invalid secret format", map[string]interface{}{"cmd": "secret create", "secret": secretInput})
					return
				}
				secrets = append(secrets, types.Secret{
					Name:  strings.TrimSpace(parts[0]),
					Value: strings.TrimSpace(parts[1]),
				})
			}

			// Create the secret version
			result, err := secretClient.CreateVersion(orgName, groupName, envName, secrets, commitMessage)
			if err != nil {
				// Handle specific CLI errors
				switch err {
				case cliErrors.ErrEmptySecrets:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("Empty secrets error during create", map[string]interface{}{
						"cmd": "secret create",
						"env": envName,
					})
					return
				case cliErrors.ErrTooManySecrets:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("Too many secrets error during create", map[string]interface{}{
						"cmd": "secret create",
						"env": envName,
					})
					return
				case cliErrors.ErrInvalidSecretName:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("Invalid secret name error during create", map[string]interface{}{
						"cmd": "secret create",
						"env": envName,
					})
					return
				case cliErrors.ErrSecretValueTooLong:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("Secret value too long error during create", map[string]interface{}{
						"cmd": "secret create",
						"env": envName,
					})
					return
				case cliErrors.ErrEncryptionFailed:
					fmt.Printf("\n❌ %v\n", err)
					logger.Error("Encryption failed during create", err, map[string]interface{}{
						"cmd": "secret create",
						"env": envName,
					})
					return
				case cliErrors.ErrEnvironmentNotFound:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("Environment not found during create", map[string]interface{}{
						"cmd": "secret create",
						"env": envName,
					})
					return
				case cliErrors.ErrSecretGroupNotFound:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("Secret group not found during create", map[string]interface{}{
						"cmd": "secret create",
						"env": envName,
					})
					return
				case cliErrors.ErrOrganizationNotFound:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("Organization not found during create", map[string]interface{}{
						"cmd": "secret create",
						"env": envName,
					})
					return
				case cliErrors.ErrNotLoggedIn:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("User not logged in during create", map[string]interface{}{
						"cmd": "secret create",
						"env": envName,
					})
					return
				case cliErrors.ErrAccessDenied:
					fmt.Printf("\n%s\n", err.Error())
					logger.Warn("Access denied during secret create", map[string]interface{}{
						"cmd": "secret create",
						"env": envName,
					})
					return
				default:
					fmt.Printf("\n❌ Failed to create secret version: %v\n", err)
					logger.Error("Failed to create secret version", err, map[string]interface{}{
						"cmd": "secret create",
						"env": envName,
					})
					return
				}
			}

			fmt.Printf("\n✅ Successfully created secret version!\n")
			fmt.Printf("   Version ID: %s\n", result.ID)
			fmt.Printf("   Environment: %s\n", result.EnvironmentID)
			fmt.Printf("   Commit Message: %s\n", result.CommitMessage)
			fmt.Printf("   Secret Count: %d\n", result.SecretCount)
			fmt.Printf("   Created At: %s\n", result.CreatedAt)

			logger.Info("Secret version created successfully", map[string]interface{}{
				"cmd":          "secret create",
				"env":          envName,
				"version_id":   result.ID,
				"secret_count": result.SecretCount,
			})
		},
	}

	// Add flags

	cmd.Flags().StringVarP(&commitMessage, "message", "m", "", "Commit message for the new version")
	cmd.Flags().StringArrayVarP(&secretsInput, "secret", "s", []string{}, "Secret in format 'name=value' (can be specified multiple times)")
	cmd.Flags().StringVarP(&orgName, "org", "o", "", "Organization name")
	cmd.Flags().StringVarP(&groupName, "group", "g", "", "Secret group name")
	cmd.Flags().StringVarP(&envName, "env", "e", "", "Environment name")

	// Mark required flags
	cmd.MarkFlagRequired("env")
	cmd.MarkFlagRequired("message")
	cmd.MarkFlagRequired("secret")

	return cmd
}
