package secret

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/secret"
	"github.com/Gkemhcs/kavach-cli/internal/types"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewPushSecretCommand creates a new command for pushing secrets from .env file
func NewPushSecretCommand(logger *utils.Logger, cfg *config.Config, secretClient secret.SecretClient) *cobra.Command {
	var (
		orgName       string
		groupName     string
		envName       string
		filePath      string
		commitMessage string
	)

	cmd := &cobra.Command{
		Use:   "push",
		Short: "📤 Push secrets from .env file",
		Long: `Push secrets directly from a .env file to create a new version.

This command reads secrets from a .env file and creates a new version with
the specified commit message. This bypasses the staging area and directly
creates secrets in the environment.

Key concepts:
- Reads secrets directly from a .env file
- Bypasses the staging area for immediate deployment
- Creates a new version with your commit message
- Supports standard .env file format (KEY=value)
- Useful for bulk secret deployment

Prerequisites:
- You must have a valid .env file with secrets
- You must have appropriate permissions for the environment
- Organization, secret group, and environment must be specified

The push process:
1. Reads and parses the .env file
2. Validates all secrets in the file
3. Creates a new version in the environment
4. Provides confirmation of successful push

Use cases:
- Deploying secrets from existing .env files
- Bulk secret deployment to environments
- Migrating secrets from other systems
- Automated deployment workflows

Examples:
  kavach secret push --org "my-org" --group "backend" --env "prod" --file ".env" --message "Deploy production secrets"
  kavach secret push --org "my-org" --group "frontend" --env "staging" --file ".env.staging" --message "Update staging config"
  kavach secret push --org "my-org" --group "backend" --env "dev" --file ".env.dev" --message "Deploy development config"

Note: This operation creates a new version with all secrets from the file.
Use 'kavach secret list' to verify the secrets are available in the environment.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Validate required flags
			if orgName == "" {
				fmt.Println("\n🚨 Error: Organization name is required. Please use --org to specify the organization.")
				logger.Warn("Organization name not provided for secret push", map[string]interface{}{"cmd": "secret push"})
				return
			}

			if groupName == "" {
				fmt.Println("\n🚨 Error: Secret group name is required. Please use --group to specify the secret group.")
				logger.Warn("Secret group name not provided for secret push", map[string]interface{}{"cmd": "secret push"})
				return
			}

			if envName == "" {
				fmt.Println("\n🚨 Error: Environment name is required. Please use --env to specify the environment.")
				logger.Warn("Environment name not provided for secret push", map[string]interface{}{"cmd": "secret push"})
				return
			}

			if filePath == "" {
				fmt.Println("\n🚨 Error: File path is required. Please use --file to specify the .env file path.")
				logger.Warn("File path not provided for secret push", map[string]interface{}{"cmd": "secret push"})
				return
			}

			if commitMessage == "" {
				fmt.Println("\n🚨 Error: Commit message is required. Please use --message to specify a commit message.")
				logger.Warn("Commit message not provided for secret push", map[string]interface{}{"cmd": "secret push"})
				return
			}

			// Read and parse .env file
			secrets, err := parseEnvFile(filePath)
			if err != nil {
				fmt.Printf("\n❌ Failed to parse .env file: %v\n", err)
				logger.Error("Failed to parse .env file", err, map[string]interface{}{
					"cmd":  "secret push",
					"file": filePath,
				})
				return
			}

			if len(secrets) == 0 {
				fmt.Printf("\n📝 No secrets found in file '%s'.\n", filePath)
				logger.Info("No secrets found in .env file", map[string]interface{}{
					"cmd":  "secret push",
					"file": filePath,
				})
				return
			}

			// Show what will be pushed
			fmt.Printf("\n📋 Pushing %d secrets from '%s':\n", len(secrets), filePath)
			fmt.Printf("   Organization: %s\n", orgName)
			fmt.Printf("   Secret Group: %s\n", groupName)
			fmt.Printf("   Environment:  %s\n", envName)
			fmt.Printf("   Commit Message: %s\n", commitMessage)
			fmt.Println("\nSecrets to be pushed:")
			for _, secret := range secrets {
				fmt.Printf("   • %s\n", secret.Name)
			}

			// Confirm push
			fmt.Print("\nDo you want to proceed? (y/N): ")
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Push cancelled.")
				logger.Info("Push cancelled by user", map[string]interface{}{"cmd": "secret push"})
				return
			}

			// Create the version
			result, err := secretClient.CreateVersion(orgName, groupName, envName, secrets, commitMessage)
			if err != nil {
				// Handle specific CLI errors
				switch err {
				case cliErrors.ErrEmptySecrets:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("Empty secrets error during push", map[string]interface{}{
						"cmd":   "secret push",
						"org":   orgName,
						"group": groupName,
						"env":   envName,
						"file":  filePath,
					})
					return
				case cliErrors.ErrTooManySecrets:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("Too many secrets error during push", map[string]interface{}{
						"cmd":   "secret push",
						"org":   orgName,
						"group": groupName,
						"env":   envName,
						"file":  filePath,
					})
					return
				case cliErrors.ErrInvalidSecretName:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("Invalid secret name error during push", map[string]interface{}{
						"cmd":   "secret push",
						"org":   orgName,
						"group": groupName,
						"env":   envName,
						"file":  filePath,
					})
					return
				case cliErrors.ErrSecretValueTooLong:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("Secret value too long error during push", map[string]interface{}{
						"cmd":   "secret push",
						"org":   orgName,
						"group": groupName,
						"env":   envName,
						"file":  filePath,
					})
					return
				case cliErrors.ErrEncryptionFailed:
					fmt.Printf("\n❌ %v\n", err)
					logger.Error("Encryption failed during push", err, map[string]interface{}{
						"cmd":   "secret push",
						"org":   orgName,
						"group": groupName,
						"env":   envName,
						"file":  filePath,
					})
					return
				case cliErrors.ErrEnvironmentNotFound:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("Environment not found during push", map[string]interface{}{
						"cmd":   "secret push",
						"org":   orgName,
						"group": groupName,
						"env":   envName,
						"file":  filePath,
					})
					return
				case cliErrors.ErrSecretGroupNotFound:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("Secret group not found during push", map[string]interface{}{
						"cmd":   "secret push",
						"org":   orgName,
						"group": groupName,
						"env":   envName,
						"file":  filePath,
					})
					return
				case cliErrors.ErrOrganizationNotFound:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("Organization not found during push", map[string]interface{}{
						"cmd":   "secret push",
						"org":   orgName,
						"group": groupName,
						"env":   envName,
						"file":  filePath,
					})
					return
				case cliErrors.ErrNotLoggedIn:
					fmt.Printf("\n❌ %v\n", err)
					logger.Warn("User not logged in during push", map[string]interface{}{
						"cmd":   "secret push",
						"org":   orgName,
						"group": groupName,
						"env":   envName,
						"file":  filePath,
					})
					return
				case cliErrors.ErrAccessDenied:
					fmt.Printf("\n%s\n", err.Error())
					logger.Warn("Access denied during secret push", map[string]interface{}{
						"cmd":   "secret push",
						"org":   orgName,
						"group": groupName,
						"env":   envName,
						"file":  filePath,
					})
					return
				default:
					fmt.Printf("\n❌ Failed to push secrets: %v\n", err)
					logger.Error("Failed to push secrets", err, map[string]interface{}{
						"cmd":   "secret push",
						"org":   orgName,
						"group": groupName,
						"env":   envName,
						"file":  filePath,
					})
					return
				}
			}

			fmt.Printf("\n✅ Successfully pushed secrets!\n")
			fmt.Printf("   Version ID: %s\n", result.ID)
			fmt.Printf("   Environment: %s\n", result.EnvironmentID)
			fmt.Printf("   Commit Message: %s\n", result.CommitMessage)
			fmt.Printf("   Secret Count: %d\n", result.SecretCount)
			fmt.Printf("   Created At: %s\n", result.CreatedAt)

			logger.Info("Secrets pushed successfully", map[string]interface{}{
				"cmd":          "secret push",
				"version_id":   result.ID,
				"secret_count": result.SecretCount,
				"org":          orgName,
				"group":        groupName,
				"env":          envName,
				"file":         filePath,
			})
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&orgName, "org", "o", "", "Organization name")
	cmd.Flags().StringVarP(&groupName, "group", "g", "", "Secret group name")
	cmd.Flags().StringVarP(&envName, "env", "e", "", "Environment name")
	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to .env file")
	cmd.Flags().StringVarP(&commitMessage, "message", "m", "", "Commit message")

	// Mark required flags
	cmd.MarkFlagRequired("org")
	cmd.MarkFlagRequired("group")
	cmd.MarkFlagRequired("env")
	cmd.MarkFlagRequired("file")
	cmd.MarkFlagRequired("message")

	return cmd
}

// parseEnvFile reads and parses a .env file
func parseEnvFile(filePath string) ([]types.Secret, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var secrets []types.Secret
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value format
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid format at line %d: %s", lineNumber, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
			(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
			value = value[1 : len(value)-1]
		}

		secrets = append(secrets, types.Secret{
			Name:  key,
			Value: value,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return secrets, nil
}
