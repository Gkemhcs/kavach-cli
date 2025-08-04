package environment

import (
	"github.com/Gkemhcs/kavach-cli/internal/config"
	envClient "github.com/Gkemhcs/kavach-cli/internal/environment"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewEnvironmentCommand returns the root Cobra command for environment management.
// Registers all subcommands for environment CRUD and activation.
func NewEnvironmentCommand(logger *utils.Logger, cfg *config.Config, envClient envClient.EnvironmentClient) *cobra.Command {
	environmentCmd := &cobra.Command{
		Use:   "env",
		Short: "🌍 Manage environments",
		Long: `Manage environments in Kavach - create, list, activate, and delete environments.

Environments are the containers where your actual secrets are stored and managed.
Each environment belongs to a secret group and can have different configurations
for different deployment stages (development, staging, production).

Key concepts:
- Environments belong to secret groups and contain actual secrets
- Each environment can have different provider configurations (GCP, Azure, GitHub)
- Environments help separate secrets by deployment stage or configuration
- Access control can be managed at the environment level
- Environments can be activated to set a default context for commands

Resource hierarchy:
Organization → Secret Groups → Environments → Secrets

Available operations:
- create: Create a new environment within a secret group
- list: List all environments in the current secret group
- activate: Set an environment as the default for future commands
- delete: Delete an environment and all its secrets
- grant: Grant permissions to users or groups on an environment
- revoke: Revoke permissions from users or groups on an environment

Common environment patterns:
- development: For development and testing
- staging: For pre-production testing
- production: For live production systems
- testing: For automated testing environments

Examples:
  kavach env create production --description "Production environment"
  kavach env list                    # List all environments
  kavach env activate production     # Set production as default
  kavach env grant --user john --role admin --env production
  kavach env delete testing          # Delete environment (with confirmation)`,
		Example: `  kavach env create production --description "Production environment"
  kavach env list
  kavach env activate production
  kavach env grant --user john --role admin --env production`,
	}
	environmentCmd.AddCommand(
		NewCreateEnvironmentCommand(logger, envClient),
		NewListEnvironmentCommand(logger, envClient),
		NewActivateEnvironmentCommand(logger, envClient),
		NewDeleteEnvironmentCommand(logger, envClient),
		NewGrantEnvironmentCommand(logger, envClient),
		NewRevokeEnvironmentCommand(logger, envClient),
	)
	return environmentCmd
}
