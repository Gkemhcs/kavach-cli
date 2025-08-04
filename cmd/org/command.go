package org

import (
	"github.com/Gkemhcs/kavach-cli/internal/config"
	"github.com/Gkemhcs/kavach-cli/internal/org"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewOrgCommand returns the root Cobra command for organization management.
// Registers all subcommands for organization CRUD and activation.
func NewOrgCommand(logger *utils.Logger, cfg *config.Config, orgClient org.OrgClient) *cobra.Command {
	orgCmd := &cobra.Command{
		Use:   "org",
		Short: "🏢 Manage organizations",
		Long: `Manage organizations in Kavach - create, list, activate, and delete organizations.

Organizations are the top-level containers in Kavach that group related resources
together. Each organization can contain multiple secret groups, environments, and
user groups. Organizations help you organize your secrets and manage access control
at a high level.

Key concepts:
- Organizations are the root containers for all Kavach resources
- Each organization can have multiple members with different roles
- Organizations can be activated to set a default context for commands
- All resources (secret groups, environments) belong to an organization

Available operations:
- create: Create a new organization
- list: List all organizations you have access to
- activate: Set an organization as the default for future commands
- delete: Delete an organization (requires confirmation)
- grant: Grant permissions to users or groups
- revoke: Revoke permissions from users or groups

Examples:
  kavach org create mycompany --description "My company organization"
  kavach org list                    # List all accessible organizations
  kavach org activate mycompany      # Set mycompany as default
  kavach org grant --user john --role admin --org mycompany
  kavach org delete mycompany        # Delete organization (with confirmation)`,
		Example: `  kavach org create mycompany --description "My company organization"
  kavach org list
  kavach org activate mycompany
  kavach org grant --user john --role admin --org mycompany`,
	}
	orgCmd.AddCommand(
		NewCreateOrgCommand(logger, orgClient),
		NewListOrgCommand(logger, orgClient),
		NewDeleteOrgCommand(logger, orgClient),
		NewActivateOrgCommand(logger, orgClient), // Add switch command
		NewGrantOrgCommand(logger, orgClient),
		NewRevokeOrgCommand(logger, orgClient),
	)
	return orgCmd
}
