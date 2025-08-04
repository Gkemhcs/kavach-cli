package provider

import (
	"fmt"
	"os"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/provider"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewUpdateCommand creates a new update command
func NewUpdateCommand(logger *utils.Logger, cfg *config.Config, providerClient provider.ProviderClient) *cobra.Command {
	var (
		// Organization flags
		orgName   string
		groupName string
		envName   string
		// GitHub flags
		githubToken            string
		githubOwner            string
		githubRepository       string
		githubEnvironment      string
		githubSecretVisibility string
		// GCP flags
		gcpKeyFile     string
		gcpProjectID   string
		gcpLocation    string
		gcpPrefix      string
		gcpReplication string
		// Azure flags
		azureTenantID       string
		azureClientID       string
		azureClientSecret   string
		azureSubscriptionID string
		azureResourceGroup  string
		azureKeyVaultName   string
		azurePrefix         string
	)

	cmd := &cobra.Command{
		Use:   "update [provider]",
		Short: "🔄 Update a provider credential",
		Long: `Update a provider credential for secret synchronization.

This command allows you to update the credentials and configuration for an existing
provider credential. You can update credentials, configuration settings, or both.

Key concepts:
- Updates existing provider credential configuration
- Supports partial updates (only specified fields)
- Validates new credentials before saving
- Maintains audit trail of changes
- Preserves existing configuration for unspecified fields

Prerequisites:
- Provider credential must already exist
- You must have appropriate permissions
- Organization, secret group, and environment must be specified
- New credentials must be valid for the provider

The update process:
1. Validates the existing provider credential exists
2. Shows current configuration details
3. Updates only specified fields
4. Validates new configuration
5. Saves updated credential

Use cases:
- Rotating provider credentials
- Updating configuration settings
- Changing provider endpoints or settings
- Maintaining credential security
- Adapting to provider changes

Examples:
  # Update GitHub provider credentials
  kavach provider update github --env=prod --org=my-org --group=my-group \
    --github-token "new-token" \
    --github-owner "new-owner" \
    --github-repository "new-repo" \
    --github-environment "production" \
    --github-secret-visibility "selected"

  # Update GCP provider credentials
  kavach provider update gcp --env=prod --org=my-org --group=my-group \
    --gcp-key-file /path/to/service-account.json \
    --gcp-project-id my-project \
    --gcp-location us-central1 \
    --gcp-prefix "kavach/" \
    --gcp-replication automatic

  # Update Azure provider credentials
  kavach provider update azure --env=prod --org=my-org --group=my-group \
    --azure-tenant-id "tenant-id" \
    --azure-client-id "client-id" \
    --azure-client-secret "client-secret" \
    --azure-subscription-id "subscription-id" \
    --azure-resource-group "resource-group" \
    --azure-key-vault-name "key-vault-name" \
    --azure-prefix "kavach/"

Note: Only specified fields will be updated. Use 'kavach provider show <provider>'
to view current configuration before updating.`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			providerName := args[0]

			// Validate provider name
			if err := provider.ValidateProviderName(providerName); err != nil {
				fmt.Fprintf(os.Stderr, "❌ %v\n", err)
				os.Exit(1)
			}

			// Validate required flags
			if orgName == "" {
				fmt.Fprintln(os.Stderr, "❌ Error: --org is required")
				os.Exit(1)
			}
			if groupName == "" {
				fmt.Fprintln(os.Stderr, "❌ Error: --group is required")
				os.Exit(1)
			}
			if envName == "" {
				fmt.Fprintln(os.Stderr, "❌ Error: --env is required")
				os.Exit(1)
			}

			// Check if provider credential exists
			existingProvider, err := providerClient.GetProviderCredential(orgName, groupName, envName, providerName)
			if err != nil {
				switch err {
				case cliErrors.ErrNotLoggedIn:
					fmt.Printf("\n🔒 You are not logged in. Please run 'kavach login' to log in.\n")
					return
				case cliErrors.ErrUnReachableBackend:
					fmt.Println("🚨 Oops! Kavach backend is currently down or not responding.")
					fmt.Println("📡 This may be due to server downtime or high request volume.")
					fmt.Println("📩 If this persists, please drop us a message at 👉 **gudikotieswarmani@gmail.com**")
					return
				case cliErrors.ErrOrganizationNotFound:
					fmt.Printf("\n❌ Organization '%s' does not exist.\n", orgName)
					return
				case cliErrors.ErrSecretGroupNotFound:
					fmt.Printf("\n❌ Secret group '%s' does not exist.\n", groupName)
					return
				case cliErrors.ErrEnvironmentNotFound:
					fmt.Printf("\n❌ Environment '%s' does not exist.\n", envName)
					return
				case cliErrors.ErrProviderCredentialNotFound:
					fmt.Printf("\n❌ Provider credential for '%s' not found.\n", providerName)
					fmt.Printf("💡 Use 'kavach provider configure %s' to create a new credential\n", providerName)
					return
				case cliErrors.ErrAccessDenied:
					fmt.Printf("\n%s\n", err.Error())
					return
				default:
					fmt.Fprintf(os.Stderr, "❌ Failed to get existing provider credential: %v\n", err)
					fmt.Fprintf(os.Stderr, "💡 Use 'kavach provider configure %s' to create a new credential\n", providerName)
					os.Exit(1)
				}
			}

			fmt.Printf("🔄 Updating %s provider credential for %s/%s/%s\n", providerName, orgName, groupName, envName)
			fmt.Printf("Current configuration:\n")
			fmt.Printf("   ID: %s\n", existingProvider.ID)
			fmt.Printf("   Type: %s\n", existingProvider.Provider)
			fmt.Printf("   Created: %s\n", existingProvider.CreatedAt)
			fmt.Printf("   Updated: %s\n", existingProvider.UpdatedAt)

			// Show current configuration summary
			currentConfig := getProviderConfigSummary(existingProvider.Provider, existingProvider.Config)
			fmt.Printf("   Configuration: %s\n\n", currentConfig)

			// Get new credentials and config based on provider type
			var credentials, newConfig map[string]interface{}
			var updateErr error

			switch providerName {
			case "github":
				credentials, newConfig, updateErr = getGitHubUpdateInput(githubToken, githubOwner, githubRepository, githubEnvironment, githubSecretVisibility)
			case "gcp":
				credentials, newConfig, updateErr = getGCPUpdateInput(gcpKeyFile, gcpProjectID, gcpLocation, gcpPrefix, gcpReplication)
			case "azure":
				credentials, newConfig, updateErr = getAzureUpdateInput(azureTenantID, azureClientID, azureClientSecret, azureSubscriptionID, azureResourceGroup, azureKeyVaultName, azurePrefix)
			default:
				fmt.Fprintf(os.Stderr, "❌ Unsupported provider: %s\n", providerName)
				os.Exit(1)
			}

			if updateErr != nil {
				fmt.Fprintf(os.Stderr, "❌ Failed to get update input: %v\n", updateErr)
				os.Exit(1)
			}

			// Update provider credential
			fmt.Printf("🔄 Updating %s provider credential...\n", providerName)

			updatedProvider, err := providerClient.UpdateProviderCredential(orgName, groupName, envName, providerName, credentials, newConfig)
			if err != nil {
				switch err {
				case cliErrors.ErrNotLoggedIn:
					fmt.Printf("\n🔒 You are not logged in. Please run 'kavach login' to log in.\n")
					return
				case cliErrors.ErrUnReachableBackend:
					fmt.Println("🚨 Oops! Kavach backend is currently down or not responding.")
					fmt.Println("📡 This may be due to server downtime or high request volume.")
					fmt.Println("📩 If this persists, please drop us a message at 👉 **gudikotieswarmani@gmail.com**")
					return
				case cliErrors.ErrOrganizationNotFound:
					fmt.Printf("\n❌ Organization '%s' does not exist.\n", orgName)
					return
				case cliErrors.ErrSecretGroupNotFound:
					fmt.Printf("\n❌ Secret group '%s' does not exist.\n", groupName)
					return
				case cliErrors.ErrEnvironmentNotFound:
					fmt.Printf("\n❌ Environment '%s' does not exist.\n", envName)
					return
				case cliErrors.ErrProviderCredentialNotFound:
					fmt.Printf("\n❌ Provider credential for '%s' not found.\n", providerName)
					return
				case cliErrors.ErrProviderCredentialUpdateFailed:
					fmt.Printf("\n❌ Failed to update provider credential.\n")
					return
				case cliErrors.ErrProviderCredentialValidationFailed:
					fmt.Printf("\n❌ Provider credential validation failed. Please check your credentials.\n")
					return
				case cliErrors.ErrAccessDenied:
					fmt.Printf("\n%s\n", err.Error())
					return
				default:
					fmt.Fprintf(os.Stderr, "❌ Failed to update provider credential: %v\n", err)
					os.Exit(1)
				}
			}

			fmt.Printf("✅ Successfully updated %s provider credential!\n", providerName)
			fmt.Printf("   ID: %s\n", updatedProvider.ID)
			fmt.Printf("   Updated: %s\n", updatedProvider.UpdatedAt)

			// Show new configuration summary
			finalConfig := getProviderConfigSummary(updatedProvider.Provider, updatedProvider.Config)
			fmt.Printf("   New Configuration: %s\n", finalConfig)
		},
	}

	// Add organization flags
	cmd.Flags().StringVar(&orgName, "org", "", "Organization name (required)")
	cmd.Flags().StringVar(&groupName, "group", "", "Secret group name (required)")
	cmd.Flags().StringVar(&envName, "env", "", "Environment name (required)")

	// Add GitHub flags
	cmd.Flags().StringVar(&githubToken, "github-token", "", "GitHub Personal Access Token")
	cmd.Flags().StringVar(&githubOwner, "github-owner", "", "GitHub organization or username")
	cmd.Flags().StringVar(&githubRepository, "github-repository", "", "GitHub repository name")
	cmd.Flags().StringVar(&githubEnvironment, "github-environment", "", "GitHub environment name")
	cmd.Flags().StringVar(&githubSecretVisibility, "github-secret-visibility", "", "Secret visibility: all, selected, private")

	// Add GCP flags
	cmd.Flags().StringVar(&gcpKeyFile, "gcp-key-file", "", "GCP Service Account JSON file path")
	cmd.Flags().StringVar(&gcpProjectID, "gcp-project-id", "", "GCP Project ID")
	cmd.Flags().StringVar(&gcpLocation, "gcp-location", "", "Secret Manager location")
	cmd.Flags().StringVar(&gcpPrefix, "gcp-prefix", "", "GCP Secret name prefix")
	cmd.Flags().StringVar(&gcpReplication, "gcp-replication", "automatic", "GCP Secret replication type")

	// Azure flags
	cmd.Flags().StringVar(&azureTenantID, "azure-tenant-id", "", "Azure tenant ID")
	cmd.Flags().StringVar(&azureClientID, "azure-client-id", "", "Azure client ID")
	cmd.Flags().StringVar(&azureClientSecret, "azure-client-secret", "", "Azure client secret")
	cmd.Flags().StringVar(&azureSubscriptionID, "azure-subscription-id", "", "Azure subscription ID")
	cmd.Flags().StringVar(&azureResourceGroup, "azure-resource-group", "", "Azure resource group")
	cmd.Flags().StringVar(&azureKeyVaultName, "azure-key-vault-name", "", "Azure Key Vault name")
	cmd.Flags().StringVar(&azurePrefix, "azure-prefix", "", "Azure Secret name prefix")

	// Mark required flags
	cmd.MarkFlagRequired("org")
	cmd.MarkFlagRequired("group")
	cmd.MarkFlagRequired("env")

	return cmd
}

// getGitHubUpdateInput processes GitHub update input from flags
func getGitHubUpdateInput(token, owner, repository, environment, secretVisibility string) (map[string]interface{}, map[string]interface{}, error) {
	credentials := map[string]interface{}{}
	if token != "" {
		credentials["token"] = token
	}

	config := map[string]interface{}{}
	if owner != "" {
		config["owner"] = owner
	}
	if repository != "" {
		config["repository"] = repository
	}
	if environment != "" {
		config["environment"] = environment
	}
	if secretVisibility != "" {
		config["secret_visibility"] = secretVisibility
	}

	return credentials, config, nil
}

// getGCPUpdateInput processes GCP update input from flags
func getGCPUpdateInput(keyFile, projectID, location, prefix, replication string) (map[string]interface{}, map[string]interface{}, error) {
	credentials := map[string]interface{}{}
	if keyFile != "" {
		// Read and parse the service account key file
		gcpCreds, err := provider.ReadGCPServiceAccountFile(keyFile)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse GCP service account key: %v", err)
		}
		credentials = gcpCreds
	}

	config := map[string]interface{}{}
	if projectID != "" {
		config["project_id"] = projectID
	}
	if location != "" {
		config["secret_manager_location"] = location
	}
	if prefix != "" {
		config["prefix"] = prefix
	}
	if replication != "" {
		config["replication"] = replication
	}

	return credentials, config, nil
}

// getAzureUpdateInput processes Azure update input from flags
func getAzureUpdateInput(tenantID, clientID, clientSecret, subscriptionID, resourceGroup, keyVaultName, prefix string) (map[string]interface{}, map[string]interface{}, error) {
	credentials := map[string]interface{}{}
	if tenantID != "" {
		credentials["tenant_id"] = tenantID
	}
	if clientID != "" {
		credentials["client_id"] = clientID
	}
	if clientSecret != "" {
		credentials["client_secret"] = clientSecret
	}

	config := map[string]interface{}{}
	if subscriptionID != "" {
		config["subscription_id"] = subscriptionID
	}
	if resourceGroup != "" {
		config["resource_group"] = resourceGroup
	}
	if keyVaultName != "" {
		config["key_vault_name"] = keyVaultName
	}
	if prefix != "" {
		config["prefix"] = prefix
	}

	return credentials, config, nil
}
