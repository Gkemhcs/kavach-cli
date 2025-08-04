package provider

import (
	"github.com/Gkemhcs/kavach-cli/internal/config"
	"github.com/Gkemhcs/kavach-cli/internal/provider"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewConfigureCommand creates a new configure command with all its subcommands
func NewConfigureCommand(logger *utils.Logger, cfg *config.Config, providerClient provider.ProviderClient) *cobra.Command {
	configureCmd := &cobra.Command{
		Use:   "configure",
		Short: "🔧 Configure provider credentials for secret synchronization",
		Long: `Configure provider credentials for secret synchronization.

This command allows you to set up credentials for different cloud providers
to enable secret synchronization from Kavach to external platforms.

Key concepts:
- Configures credentials for external cloud providers
- Enables automated secret synchronization
- Supports multiple provider types simultaneously
- Validates credentials before saving
- Integrates with cloud-native secret management

Available providers:
  • github - Configure GitHub Personal Access Token
    - Sync secrets to GitHub repository secrets
    - Integrate with GitHub Actions workflows
    - Support for GitHub environments

  • gcp - Configure Google Cloud Service Account
    - Sync secrets to Google Cloud Secret Manager
    - Support for custom replication settings
    - Regional secret management

  • azure - Configure Azure Service Principal
    - Sync secrets to Azure Key Vault
    - Support for custom secret prefixes
    - Enterprise-grade secret management

Prerequisites:
- Valid credentials for the target cloud provider
- Appropriate permissions in the cloud platform
- Organization, secret group, and environment must exist

Use cases:
- Automated secret deployment to cloud platforms
- Integration with CI/CD pipelines
- Multi-cloud secret management
- Disaster recovery and backup strategies

Examples:
  kavach provider configure github --token "ghp_xxx" --owner "myorg" --repo "myrepo"
  kavach provider configure gcp --key-file "service-account.json" --project-id "my-project"
  kavach provider configure azure --tenant-id "xxx" --client-id "xxx" --client-secret "xxx"

Note: After configuration, use 'kavach secret sync --provider <provider>' to sync
secrets to the configured provider.`,
	}

	// Add subcommands
	configureCmd.AddCommand(
		NewConfigureGitHubCommand(logger, cfg, providerClient),
		NewConfigureGCPCommand(logger, cfg, providerClient),
		NewConfigureAzureCommand(logger, cfg, providerClient),
	)

	return configureCmd
}
