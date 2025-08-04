package org

import (
	"fmt"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/org"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewActivateOrgCommand returns a Cobra command to activate (switch to) a different organization.
// Handles user feedback, error reporting, and logging.
func NewActivateOrgCommand(logger *utils.Logger, orgClient org.OrgClient) *cobra.Command {
	return &cobra.Command{
		Use:   "activate <org_name>",
		Short: "🔄 Switch to a different organization",
		Long: `Activate an organization to set it as the default context for future commands.

When you activate an organization, it becomes the default context for all
subsequent CLI commands. This means you won't need to specify the organization
explicitly in other commands like secret group creation, environment management, etc.

Key benefits:
- Reduces the need to specify --org flag in every command
- Provides a consistent working context
- Makes command usage more convenient
- Helps avoid accidentally working in the wrong organization

The activation process:
1. Verifies you have access to the specified organization
2. Sets the organization as the default context
3. Saves the configuration locally
4. Confirms the activation

Examples:
  kavach org activate mycompany      # Set mycompany as default
  kavach org activate project-alpha  # Switch to project-alpha
  kavach org list                    # See which org is active (🟢)

Note: You can still override the active organization by explicitly
specifying --org flag in individual commands. The active organization
is used only when no organization is explicitly provided.`,
		Example: `  kavach org activate mycompany
  # Set mycompany as the default organization for future commands`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			orgName := args[0]
			logger.Info("Activating organization", map[string]interface{}{"cmd": "org activate", "org": orgName})
			orgObj, err := orgClient.GetOrganizationByName(orgName)
			if err != nil {
				if err == cliErrors.ErrNotLoggedIn {
					fmt.Printf("\n🔒 You are not logged in. Please run 'kavach login' to log in.\n")
					logger.Warn("User not logged in during org activate", map[string]interface{}{"cmd": "org activate", "org": orgName})
					return nil
				}
				if err == cliErrors.ErrUnReachableBackend {
					fmt.Println("🚨 Oops! Kavach backend is currently down or not responding.")
					fmt.Println("📡 This may be due to server downtime or high request volume.")
					fmt.Println("📩 If this persists, please drop us a message at 👉 **gudikotieswarmani@gmail.com**")
					logger.Error("Backend unreachable during org activate", err, map[string]interface{}{"cmd": "org activate", "org": orgName})
					return nil // or exit gracefully
				}
				if err == cliErrors.ErrOrganizationNotFound {
					fmt.Printf("\n❌ Organization '%s' not found. Cannot activate.\n", orgName)
					logger.Warn("Organization not found during activate", map[string]interface{}{"cmd": "org activate", "org": orgName})
					return nil
				}
				if err == cliErrors.ErrAccessDenied {
					fmt.Printf("\n%s\n", err.Error())
					logger.Warn("Access denied during org activate", map[string]interface{}{"cmd": "org activate", "org": orgName})
					return nil
				}
				logger.Error("Failed to get organization during activate", err, map[string]interface{}{"cmd": "org activate", "org": orgName})
				return err
			}
			cfg, err := config.LoadCLIConfig()
			if err != nil {
				logger.Error("Failed to load CLI config during org activate", err, map[string]interface{}{"cmd": "org activate", "org": orgName})
				return err
			}
			cfg.Organization = orgObj.Name
			if err := config.SaveCLIConfig(cfg); err != nil {
				logger.Error("Failed to save CLI config during org activate", err, map[string]interface{}{"cmd": "org activate", "org": orgName})
				return err
			}
			fmt.Printf("\n✅ Organization '%s' is now active.\n", orgObj.Name)
			logger.Info("Organization activated", map[string]interface{}{"cmd": "org activate", "org": orgObj.Name})
			return nil
		},
	}
}
