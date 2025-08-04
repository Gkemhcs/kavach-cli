package secret

import (
	"fmt"
	"strings"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/secret"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewSyncSecretCommand creates a new sync command for secrets
func NewSyncSecretCommand(logger *utils.Logger, cfg *config.Config, secretClient secret.SecretClient) *cobra.Command {
	var (
		envName   string
		groupName string
		orgName   string
		provider  string
		versionID string
	)

	syncCmd := &cobra.Command{
		Use:   "sync",
		Short: "🔄 Sync secrets to external providers",
		Long: `Sync secrets from a specific environment to external providers like GitHub, GCP, or Azure.

This command allows you to synchronize your secrets with external systems for deployment
or configuration management purposes. You can specify a particular version to sync or
use the latest version by default.

Key concepts:
- Synchronizes secrets from Kavach to external cloud providers
- Uses pre-configured provider credentials
- Supports multiple provider types simultaneously
- Can sync specific versions or latest version
- Provides detailed sync results and error reporting

Supported providers:
- GitHub: Sync to GitHub repository secrets
- GCP: Sync to Google Cloud Secret Manager
- Azure: Sync to Azure Key Vault

Prerequisites:
- Provider must be configured using 'kavach provider configure'
- You must have appropriate permissions for the environment
- External provider credentials must be valid

The sync process:
1. Retrieves secrets from the specified version
2. Connects to configured external providers
3. Creates or updates secrets in external systems
4. Provides detailed sync results and statistics

Use cases:
- Deploying secrets to cloud environments
- Maintaining synchronization across platforms
- Automated secret deployment workflows
- Disaster recovery and backup strategies

Examples:
  kavach secret sync --env production --group backend --org myorg --provider github
  kavach secret sync --env staging --group frontend --org myorg --provider gcp --version-id v1.2.3
  kavach secret sync --env dev --group backend --org myorg --provider azure

Note: Use 'kavach provider list' to see configured providers and 'kavach secret list'
to see available versions for syncing.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate required flags
			if envName == "" {
				return fmt.Errorf("❌ environment name is required (--env)")
			}
			if groupName == "" {
				return fmt.Errorf("❌ secret group name is required (--group)")
			}
			if orgName == "" {
				return fmt.Errorf("❌ organization name is required (--org)")
			}
			if provider == "" {
				return fmt.Errorf("❌ provider is required (--provider)")
			}

			// Validate provider
			validProviders := []string{"github", "gcp", "azure"}
			isValidProvider := false
			for _, p := range validProviders {
				if strings.ToLower(provider) == p {
					provider = p
					isValidProvider = true
					break
				}
			}
			if !isValidProvider {
				return fmt.Errorf("❌ invalid provider '%s'. Supported providers: %s", provider, strings.Join(validProviders, ", "))

			}

			logger.Info("Starting secret sync", map[string]interface{}{
				"cmd":        "secret sync",
				"env":        envName,
				"group":      groupName,
				"org":        orgName,
				"provider":   provider,
				"version_id": versionID,
			})

			// Call the sync method
			result, err := secretClient.SyncSecrets(orgName, groupName, envName, provider, versionID)
			if err != nil {
				logger.Error("Secret sync failed", err, map[string]interface{}{
					"cmd":      "secret sync",
					"env":      envName,
					"group":    groupName,
					"org":      orgName,
					"provider": provider,
				})

				// Handle specific error types
				switch err {
				case cliErrors.ErrNotLoggedIn:
					return err
				case cliErrors.ErrConnectionFailed:
					return err
				case cliErrors.ErrOrganizationNotFound:
					return err
				case cliErrors.ErrSecretGroupNotFound:
					return err
				case cliErrors.ErrEnvironmentNotFound:
					return err
				case cliErrors.ErrInvalidProviderType:
					return err
				case cliErrors.ErrProviderCredentialNotFound:
					return err
				case cliErrors.ErrProviderCredentialValidationFailed:
					return err
				case cliErrors.ErrProviderSyncFailed:
					return err
				case cliErrors.ErrNoSecretsToSync:
					return err
				case cliErrors.ErrSecretVersionNotFound:
					return err
				case cliErrors.ErrDecryptionFailed:
					return err
				case cliErrors.ErrInternalServer:
					return err
				case cliErrors.ErrGitHubEnvironmentNotFound:
					return err
				case cliErrors.ErrGitHubEncryptionFailed:
					return err
				case cliErrors.ErrGCPInvalidLocation:
					return err
				case cliErrors.ErrAccessDenied:
					return err
				default:
					// For any other errors, return a generic message
					return fmt.Errorf("❌ Secret sync failed: %v", err)
				}
			}

			// Display results
			fmt.Printf("🔄 Secret sync completed successfully!\n\n")
			fmt.Printf("📊 Sync Summary:\n")
			fmt.Printf("   Provider: %s\n", result.Provider)
			fmt.Printf("   Status: %s\n", result.Status)
			fmt.Printf("   Total Secrets: %d\n", result.TotalCount)
			fmt.Printf("   Successfully Synced: %d\n", result.SyncedCount)
			fmt.Printf("   Failed: %d\n", result.FailedCount)
			fmt.Printf("   Synced At: %s\n\n", result.SyncedAt)

			if result.Message != "" {
				fmt.Printf("📝 Message: %s\n\n", result.Message)
			}

			// Show detailed results if there are any
			if len(result.Results) > 0 {
				fmt.Printf("📋 Detailed Results:\n")
				for _, res := range result.Results {
					if res.Success {
						fmt.Printf("   ✅ %s: Synced successfully\n", res.Name)
					} else {
						fmt.Printf("   ❌ %s: Failed - %s\n", res.Name, res.Error)
					}
				}
				fmt.Println()
			}

			// Show errors if any
			if len(result.Errors) > 0 {
				fmt.Printf("⚠️  Errors:\n")
				for _, err := range result.Errors {
					fmt.Printf("   • %s\n", err)
				}
				fmt.Println()
			}

			logger.Info("Secret sync completed", map[string]interface{}{
				"cmd":          "secret sync",
				"env":          envName,
				"group":        groupName,
				"org":          orgName,
				"provider":     provider,
				"status":       result.Status,
				"synced_count": result.SyncedCount,
				"failed_count": result.FailedCount,
			})

			return nil
		},
	}

	// Add flags
	syncCmd.Flags().StringVarP(&envName, "env", "e", "", "Environment name (required)")
	syncCmd.Flags().StringVarP(&groupName, "group", "g", "", "Secret group name (required)")
	syncCmd.Flags().StringVarP(&orgName, "org", "o", "", "Organization name (required)")
	syncCmd.Flags().StringVarP(&provider, "provider", "p", "", "Provider type: github, gcp, or azure (required)")
	syncCmd.Flags().StringVarP(&versionID, "version-id", "v", "", "Specific version ID to sync (optional, defaults to latest version)")

	// Mark required flags
	syncCmd.MarkFlagRequired("env")
	syncCmd.MarkFlagRequired("group")
	syncCmd.MarkFlagRequired("org")
	syncCmd.MarkFlagRequired("provider")

	return syncCmd
}
