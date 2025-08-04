package provider

import (
	"github.com/Gkemhcs/kavach-cli/internal/config"
	"github.com/Gkemhcs/kavach-cli/internal/provider"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewProviderCommand creates a new provider command with all its subcommands
func NewProviderCommand(logger *utils.Logger, cfg *config.Config, providerClient provider.ProviderClient) *cobra.Command {
	providerCmd := &cobra.Command{
		Use:   "provider",
		Short: "🔄 Manage cloud provider credentials for secret synchronization",
		Long: `Manage cloud provider credentials for secret synchronization.

This command allows you to configure, list, view, and delete provider credentials
for syncing secrets to various cloud platforms like GitHub, GCP, and Azure.

Key concepts:
- Manages credentials for external cloud providers
- Enables automated secret synchronization
- Supports multiple provider types simultaneously
- Provides comprehensive credential management
- Integrates with cloud-native secret management

Supported providers:
  • GitHub - Sync secrets to GitHub repositories
    - GitHub repository secrets
    - GitHub Actions integration
    - Environment-specific secrets

  • GCP - Sync secrets to Google Cloud Secret Manager
    - Regional secret management
    - Custom replication settings
    - Version control and rotation

  • Azure - Sync secrets to Azure Key Vault
    - Enterprise-grade secret management
    - Custom secret prefixes
    - Azure Key Vault integration

Available operations:
- configure: Set up new provider credentials
- list: View all configured providers
- show: Display detailed provider information
- update: Modify existing provider credentials
- delete: Remove provider credentials

Prerequisites:
- Valid credentials for the target cloud provider
- Appropriate permissions in the cloud platform
- Organization, secret group, and environment must exist

Use cases:
- Multi-cloud secret management
- Automated secret deployment
- CI/CD pipeline integration
- Disaster recovery and backup
- Compliance and security management

Examples:
  kavach provider configure github --token "ghp_xxx" --owner "myorg" --repo "myrepo"
  kavach provider configure gcp --key-file "service-account.json" --project-id "my-project"
  kavach provider configure azure --tenant-id "xxx" --client-id "xxx" --client-secret "xxx"
  kavach provider list
  kavach provider show github
  kavach provider delete gcp

Note: After configuring providers, use 'kavach secret sync --provider <provider>'
to sync secrets to the configured platforms.`,
	}

	// Add subcommands
	providerCmd.AddCommand(
		NewConfigureCommand(logger, cfg, providerClient),
		NewListCommand(logger, cfg, providerClient),
		NewShowCommand(logger, cfg, providerClient),
		NewUpdateCommand(logger, cfg, providerClient),
		NewDeleteCommand(logger, cfg, providerClient),
	)

	return providerCmd
}
