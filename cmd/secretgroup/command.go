package secretgroup

import (
	"github.com/Gkemhcs/kavach-cli/internal/config"

	"github.com/Gkemhcs/kavach-cli/internal/secretgroup"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewSecretGroupCommand creates a new command for managing secret groups
func NewSecretGroupCommand(logger *utils.Logger, cfg *config.Config, groupClient secretgroup.SecretGroupClient) *cobra.Command {
	groupCmd := &cobra.Command{
		Use:   "group",
		Short: "🔐 Manage secret groups",
		Long: `Manage secret groups in Kavach - create, list, activate, and delete secret groups.

Secret groups are logical containers that organize related secrets within an organization.
Each secret group can contain multiple environments (like development, staging, production)
and provides a way to manage access control and organization of your secrets.

Key concepts:
- Secret groups belong to organizations and contain environments
- Each secret group can have multiple environments with different configurations
- Secret groups help organize secrets by project, team, or application
- Access control can be managed at the secret group level
- Secret groups can be activated to set a default context for commands

Resource hierarchy:
Organization → Secret Groups → Environments → Secrets

Available operations:
- create: Create a new secret group within an organization
- list: List all secret groups in the current organization
- activate: Set a secret group as the default for future commands
- delete: Delete a secret group and all its environments/secrets
- grant: Grant permissions to users or groups on a secret group
- revoke: Revoke permissions from users or groups on a secret group

Examples:
  kavach group create myapp --description "My application secrets"
  kavach group list                    # List all secret groups
  kavach group activate myapp          # Set myapp as default
  kavach group grant --user john --role admin --group myapp
  kavach group delete myapp            # Delete secret group (with confirmation)`,
	}

	groupCmd.AddCommand(
		NewCreateSecretGroupCommand(logger, groupClient),
		NewListSecretGroupCommand(logger, groupClient),
		NewActivateSecretGroupCommand(logger, groupClient),
		NewDeleteSecretGroupCommand(logger, groupClient),
		NewGrantSecretGroupCommand(logger, groupClient),
		NewRevokeSecretGroupCommand(logger, groupClient),
		NewListBindingsCommand(logger, groupClient),
	)
	return groupCmd
}
