package environment

import (
	"fmt"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	"github.com/Gkemhcs/kavach-cli/internal/environment"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewActivateEnvironmentCommand returns a Cobra command to activate (switch to) a different environment.
// Handles user feedback, error reporting, and logging.
func NewActivateEnvironmentCommand(logger *utils.Logger, envClient environment.EnvironmentClient) *cobra.Command {
	var orgName string
	var groupName string
	cmd := &cobra.Command{
		Use:   "activate <env_name>",
		Short: "🔄 Switch to a different environment",
		Long: `Activate an environment to set it as the default context for future commands.

When you activate an environment, it becomes the default context for all
subsequent CLI commands. This means you won't need to specify the environment
explicitly in other commands like secret management, provider configuration, etc.

Key benefits:
- Reduces the need to specify --env flag in every command
- Provides a consistent working context within the secret group
- Makes command usage more convenient
- Helps avoid accidentally working in the wrong environment

The activation process:
1. Verifies you have access to the specified environment
2. Sets the environment as the default context
3. Saves the configuration locally
4. Confirms the activation

Examples:
  kavach env activate production      # Set production as default environment
  kavach env activate development     # Switch to development environment
  kavach env list                     # See which environment is active (🟢)

Note: You can still override the active environment by explicitly
specifying --env flag in individual commands. The active environment
is used only when no environment is explicitly provided.`,
		Example: `  kavach env activate production
  # Set production as the default environment for future commands`,
		Args: cobra.ExactArgs(1), // Requires exactly one argument: environment name
		RunE: func(cmd *cobra.Command, args []string) error {
			envName := args[0] // The environment to activate
			logger.Info("Activating environment", map[string]interface{}{"cmd": "env activate", "env": envName, "org": orgName, "group": groupName})
			cfg, _ := config.LoadCLIConfig()
			if cfg.Organization == "" && cfg.SecretGroup == "" && orgName == "" && groupName == "" {
				fmt.Println("\n⚠️  You haven't passed an organization and secret group and no default organization  and secret group are set.")
				fmt.Println("💡 Set a default organization using: `kavach org activate <org-name>` and `kavach group activate <secret-group-name>`")
				logger.Warn("No organization or secret group set for environment activation", map[string]interface{}{"cmd": "env activate"})
				return nil
			}
			if orgName == "" {
				orgName = cfg.Organization
			}
			if groupName == "" {
				groupName = cfg.SecretGroup
			}
			_, err := envClient.GetEnvironmentbyName(orgName, groupName, envName)
			if err != nil {
				if err == cliErrors.ErrNotLoggedIn {
					fmt.Printf("\n🚫 You are not logged in. Please run 'kavach login'.\n")
					logger.Warn("User not logged in during environment activate", map[string]interface{}{"cmd": "env activate"})
					return nil
				}
				if err == cliErrors.ErrUnReachableBackend {
					fmt.Println("🚨 Oops! Kavach backend is currently down or not responding.")
					fmt.Println("📡 This may be due to server downtime or high request volume.")
					fmt.Println("📩 If this persists, please drop us a message at 👉 **gudikotieswarmani@gmail.com**")
					logger.Error("Backend unreachable during environment activate", err, map[string]interface{}{"cmd": "env activate"})
					return nil // or exit gracefully
				}
				if err == cliErrors.ErrOrganizationNotFound {
					fmt.Printf("\n❌ Organization '%s' not found. Cannot activate environment '%s'.\n", orgName, envName)
					logger.Warn("Organization not found during environment activate", map[string]interface{}{"cmd": "env activate", "org": orgName})
					return nil
				}
				logger.Error("Failed to get organization during environment activate", err, map[string]interface{}{"cmd": "env activate", "org": orgName, "group": groupName})
				if err == cliErrors.ErrSecretGroupNotFound {
					fmt.Printf("\n❌ Secret group '%s' not found in organization '%s'. Cannot activate.\n", groupName, orgName)
					logger.Warn("Secret group not found during environment activate", map[string]interface{}{"cmd": "env activate", "group": groupName, "org": orgName})
					return nil
				}
				if err == cliErrors.ErrEnvironmentNotFound {
					fmt.Printf("\n❌ Environment '%s' can't be activated as it doesn't exist under secret group '%s' in organization '%s'.\n", envName, groupName, orgName)
					logger.Warn("Environment not found during activate", map[string]interface{}{"cmd": "env activate", "env": envName, "group": groupName, "org": orgName})
					return nil
				}
				if err == cliErrors.ErrAccessDenied {
					fmt.Printf("\n%s\n", err.Error())
					logger.Warn("Access denied during environment activate", map[string]interface{}{"cmd": "env activate", "env": envName, "group": groupName, "org": orgName})
					return nil
				}
				return err
			}
			cfg.Organization = orgName
			cfg.SecretGroup = groupName
			cfg.Environment = envName
			if err := config.SaveCLIConfig(cfg); err != nil {
				logger.Error("Failed to save CLI config during environment activate", err, map[string]interface{}{"cmd": "env activate", "env": envName, "group": groupName, "org": orgName})
				return err
			}
			fmt.Printf("\n✅ Environment '%s' is now active under secret group '%s' in organization '%s'.\n", envName, groupName, orgName)
			logger.Info("Environment activated", map[string]interface{}{"cmd": "env activate", "env": envName, "group": groupName, "org": orgName})
			return nil
		},
	}
	cmd.Flags().StringVarP(&orgName, "organization", "o", "", "Organization under which to activate the environment")
	cmd.Flags().StringVarP(&groupName, "group", "g", "", "Secret group under which to activate the environment")
	return cmd
}
