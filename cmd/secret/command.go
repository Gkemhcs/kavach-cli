package secret

import (
	"github.com/Gkemhcs/kavach-cli/internal/config"
	"github.com/Gkemhcs/kavach-cli/internal/secret"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewSecretCommand creates a new secret command with all its subcommands
func NewSecretCommand(logger *utils.Logger, cfg *config.Config, secretClient secret.SecretClient) *cobra.Command {
	secretCmd := &cobra.Command{
		Use:   "secret",
		Short: "🔐 Manage secrets",
		Long: `Manage secrets in Kavach - create, list, and manage secrets.

Secrets are the actual sensitive data stored in your environments. They can be
configuration values, API keys, passwords, certificates, or any other sensitive
information your applications need.

Key concepts:
- Secrets are stored within environments and belong to secret groups
- Each secret can have multiple versions for tracking changes
- Secrets can be synced to external providers (GCP, Azure, GitHub)
- Secrets support different types (string, file, certificate)
- Access control is managed at the environment level

Resource hierarchy:
Organization → Secret Groups → Environments → Secrets

Available operations:
- add: Add a secret to the local staging area (for batch operations)
- commit: Commit staged secrets to the environment
- push: Push secrets to external providers
- list: List all secrets in the current environment
- details: Show detailed information about a specific secret
- export: Export secrets to .env file
- sync: Sync secrets from external providers
- rollback: Rollback to a previous version of a secret
- diff: Show differences between secret versions

Use cases:
- Store application configuration values
- Manage API keys and tokens
- Store database credentials
- Manage SSL certificates
- Store sensitive configuration data

Examples:
  kavach secret add db-password --value "mypassword123"
  kavach secret list                    # List all secrets
  kavach secret details api-key         # Show secret details
  kavach secret export --version abc123 # Export secrets to .env file
  kavach secret add config.json --file ./config.json
  kavach secret commit                  # Commit staged secrets
  kavach secret push                    # Push to external providers`,
	}

	// Add subcommands
	secretCmd.AddCommand(
		NewAddSecretCommand(logger, cfg, secretClient),
		NewCommitSecretCommand(logger, cfg, secretClient),
		NewPushSecretCommand(logger, cfg, secretClient),
		NewListVersionsCommand(logger, cfg, secretClient),
		NewGetVersionDetailsCommand(logger, cfg, secretClient),
		NewExportCommand(logger, cfg, secretClient),
		NewRollbackCommand(logger, cfg, secretClient),
		NewDiffCommand(logger, cfg, secretClient),
		NewSyncSecretCommand(logger, cfg, secretClient),
	)

	return secretCmd
}
