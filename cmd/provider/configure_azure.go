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

// NewConfigureAzureCommand creates a new configure azure command
func NewConfigureAzureCommand(logger *utils.Logger, cfg *config.Config, providerClient provider.ProviderClient) *cobra.Command {
	var (
		// Azure configuration flags
		azureTenantID             string
		azureClientID             string
		azureClientSecret         string
		azureSubscriptionID       string
		azureDisableOlderVersions bool
		azureResourceGroup        string
		azureKeyVaultName         string
		azurePrefix               string
		// Organization flags
		orgName   string
		groupName string
		envName   string
	)

	cmd := &cobra.Command{
		Use:   "azure",
		Short: "🔷 Configure Azure provider credentials",
		Long: `Configure Azure provider credentials for secret synchronization.

This command sets up Azure Service Principal credentials to sync secrets
to Azure Key Vault.

Key concepts:
- Configures Azure Service Principal for secret synchronization
- Enables syncing secrets to Azure Key Vault
- Supports custom secret naming with prefixes
- Validates credentials before saving
- Integrates with Azure Key Vault for secure storage

Prerequisites:
- Azure Service Principal with Key Vault permissions
- Azure Key Vault must exist and be accessible
- Service Principal must have appropriate Key Vault roles
- Organization, secret group, and environment must exist

Required Azure permissions:
- Key Vault Secrets Officer (for creating/updating secrets)
- Key Vault Secrets User (for reading secrets)
- Key Vault Administrator (for managing vault settings)

Required flags:
  --tenant-id         Azure Tenant ID
  --client-id         Azure Client ID (Service Principal ID)
  --client-secret     Azure Client Secret
  --subscription-id   Azure Subscription ID
  --resource-group    Azure Resource Group name
  --key-vault-name    Azure Key Vault name
  --org               Organization name
  --group             Secret group name
  --env               Environment name

Optional flags:
  --prefix            Secret name prefix (e.g., "kavach/")

Use cases:
- Syncing secrets to Azure Key Vault
- Integrating with Azure-based applications
- Maintaining secrets in Azure cloud environment
- Automated secret deployment to Azure

Example:
  kavach provider configure azure \
    --tenant-id "00000000-0000-0000-0000-000000000000" \
    --client-id "11111111-1111-1111-1111-111111111111" \
    --client-secret "my-secret-key" \
    --subscription-id "22222222-2222-2222-2222-222222222222" \
    --resource-group "my-resource-group" \
    --key-vault-name "my-key-vault" \
    --prefix "kavach/" \
    --org "myorg" \
    --group "mygroup" \
    --env "prod"

Note: After configuration, use 'kavach secret sync --provider azure' to sync
secrets to Azure Key Vault.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Validate required flags
			if azureTenantID == "" {
				fmt.Fprintln(os.Stderr, "❌ Error: --tenant-id is required")
				os.Exit(1)
			}
			if azureClientID == "" {
				fmt.Fprintln(os.Stderr, "❌ Error: --client-id is required")
				os.Exit(1)
			}
			if azureClientSecret == "" {
				fmt.Fprintln(os.Stderr, "❌ Error: --client-secret is required")
				os.Exit(1)
			}
			if azureSubscriptionID == "" {
				fmt.Fprintln(os.Stderr, "❌ Error: --subscription-id is required")
				os.Exit(1)
			}
			if azureResourceGroup == "" {
				fmt.Fprintln(os.Stderr, "❌ Error: --resource-group is required")
				os.Exit(1)
			}
			if azureKeyVaultName == "" {
				fmt.Fprintln(os.Stderr, "❌ Error: --key-vault-name is required")
				os.Exit(1)
			}
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

			// Prepare credentials and config
			credentials := map[string]interface{}{
				"tenant_id":     azureTenantID,
				"client_id":     azureClientID,
				"client_secret": azureClientSecret,
			}

			config := map[string]interface{}{
				"subscription_id": azureSubscriptionID,
				"resource_group":  azureResourceGroup,
				"key_vault_name":  azureKeyVaultName,
			}
			if azureDisableOlderVersions {
				config["disable_older_versions"] = azureDisableOlderVersions
			}

			if azurePrefix != "" {
				config["prefix"] = azurePrefix
			}

			if azureDisableOlderVersions {
				config["delete_older_versions"] = azureDisableOlderVersions
			}

			// Create provider credential
			fmt.Printf("🔧 Configuring Azure provider for %s/%s/%s...\n", orgName, groupName, envName)

			result, err := providerClient.CreateProviderCredential(orgName, groupName, envName, "azure", credentials, config)
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
				case cliErrors.ErrProviderCredentialExists:
					fmt.Printf("\n⚠️ Provider credential for Azure already exists.\n")
					return
				case cliErrors.ErrProviderCredentialValidationFailed:
					fmt.Printf("\n❌ Azure provider credential validation failed. Please check your credentials.\n")
					return
				case cliErrors.ErrProviderCredentialCreateFailed:
					fmt.Printf("\n❌ Failed to create Azure provider credential.\n")
					return
				case cliErrors.ErrAccessDenied:
					fmt.Printf("\n%s\n", err.Error())
					return
				default:
					fmt.Fprintf(os.Stderr, "❌ Failed to configure Azure provider: %v\n", err)
					os.Exit(1)
				}
			}

			fmt.Printf("✅ Successfully configured Azure provider!\n")
			fmt.Printf("   Provider ID: %s\n", result.ID)
			fmt.Printf("   Tenant ID: %s\n", azureTenantID)
			fmt.Printf("   Client ID: %s\n", azureClientID)
			fmt.Printf("   Subscription ID: %s\n", azureSubscriptionID)
			fmt.Printf("   Resource Group: %s\n", azureResourceGroup)
			fmt.Printf("   Key Vault: %s\n", azureKeyVaultName)
			if azurePrefix != "" {
				fmt.Printf("   Secret Prefix: %s\n", azurePrefix)
			}
			fmt.Printf("   Created: %s\n", result.CreatedAt)
		},
	}

	// Add flags
	cmd.Flags().StringVar(&azureTenantID, "tenant-id", "", "Azure Tenant ID (required)")
	cmd.Flags().StringVar(&azureClientID, "client-id", "", "Azure Client ID (Service Principal ID) (required)")
	cmd.Flags().StringVar(&azureClientSecret, "client-secret", "", "Azure Client Secret (required)")
	cmd.Flags().StringVar(&azureSubscriptionID, "subscription-id", "", "Azure Subscription ID (required)")
	cmd.Flags().StringVar(&azureResourceGroup, "resource-group", "", "Azure Resource Group name (required)")
	cmd.Flags().StringVar(&azureKeyVaultName, "key-vault-name", "", "Azure Key Vault name (required)")
	cmd.Flags().BoolVar(&azureDisableOlderVersions, "disable-older-versions", false, "Disable older secret versions when creating new ones")
	cmd.Flags().StringVar(&azurePrefix, "prefix", "", "Secret name prefix")
	cmd.Flags().StringVar(&orgName, "org", "", "Organization name (required)")
	cmd.Flags().StringVar(&groupName, "group", "", "Secret group name (required)")
	cmd.Flags().StringVar(&envName, "env", "", "Environment name (required)")

	// Mark required flags
	cmd.MarkFlagRequired("tenant-id")
	cmd.MarkFlagRequired("client-id")
	cmd.MarkFlagRequired("client-secret")
	cmd.MarkFlagRequired("subscription-id")
	cmd.MarkFlagRequired("resource-group")
	cmd.MarkFlagRequired("key-vault-name")
	cmd.MarkFlagRequired("org")
	cmd.MarkFlagRequired("group")
	cmd.MarkFlagRequired("env")

	return cmd
}
