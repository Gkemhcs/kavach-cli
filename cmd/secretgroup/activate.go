package secretgroup

import (
	"fmt"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/secretgroup"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewActivateSecretGroupCommand creates a new command for activating secret groups
func NewActivateSecretGroupCommand(logger *utils.Logger, groupClient secretgroup.SecretGroupClient) *cobra.Command {
	var orgName string
	cmd := &cobra.Command{
		Use:   "activate <group_name>",
		Short: "🔄 Switch to a different secret group",
		Long: `Activate a secret group to set it as the default context for future commands.

When you activate a secret group, it becomes the default context for all
subsequent CLI commands. This means you won't need to specify the secret group
explicitly in other commands like environment creation, secret management, etc.

Key benefits:
- Reduces the need to specify --group flag in every command
- Provides a consistent working context within the organization
- Makes command usage more convenient
- Helps avoid accidentally working in the wrong secret group

The activation process:
1. Verifies you have access to the specified secret group
2. Sets the secret group as the default context
3. Saves the configuration locally
4. Confirms the activation

Examples:
  kavach group activate myapp      # Set myapp as default secret group
  kavach group activate backend    # Switch to backend secret group
  kavach group list                # See which group is active (🟢)

Note: You can still override the active secret group by explicitly
specifying --group flag in individual commands. The active secret group
is used only when no secret group is explicitly provided.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupName := args[0]
			cfg, _ := config.LoadCLIConfig()

			if cfg.Organization == "" && orgName == "" {
				fmt.Println("\n⚠️  You haven't passed an organization and no default organization is set.")
				fmt.Println("💡 Set a default organization using: `kavach org activate <org-name>`")
				return nil
			}
			if orgName == "" {
				orgName = cfg.Organization
			}
			_, err := groupClient.GetSecretGroupByName(orgName, groupName)
			if err != nil {

				if err == cliErrors.ErrNotLoggedIn {
					fmt.Printf("\n🚫You are not logged in , please run kavach login \n")
					return nil
				}
				if err == cliErrors.ErrUnReachableBackend {
					fmt.Println("🚨 Oops! Kavach backend is currently down or not responding.")
					fmt.Println("📡 This may be due to server downtime or high request volume.")
					fmt.Println("📩 If this persists, please drop us a message at 👉 **gudikotieswarmani@gmail.com**")
					return nil // or exit gracefully
				}
				if err == cliErrors.ErrOrganizationNotFound {
					fmt.Printf("\n❌ Organization '%s' not found. Cannot activate secret group '%s'.\n", orgName, groupName)
					return nil
				}
				logger.Error("Failed to get organization", err)
				if err == cliErrors.ErrSecretGroupNotFound {
					fmt.Printf("\n❌ Secret group '%s' not found in organization '%s'. Cannot activate.\n", groupName, orgName)
					return nil
				}
				if err == cliErrors.ErrAccessDenied {
					fmt.Printf("\n%s\n", err.Error())
					return nil
				}
				// Check if the error message contains authentication-related text
				if cliErrors.IsAuthenticationError(err) {
					fmt.Printf("\n🔑 Please login again, the session is expired, unable to authenticate you\n")
					return nil
				}
				return err
			}

			cfg.Organization = orgName

			cfg.SecretGroup = groupName
			if err := config.SaveCLIConfig(cfg); err != nil {
				logger.Error("Failed to save CLI config", err)
				return err
			}
			fmt.Printf("\n✅ Secret group '%s' is now active under organization '%s'.\n", groupName, orgName)
			return nil
		},
	}
	cmd.Flags().StringVarP(&orgName, "organization", "o", "", "Organisation under which we need to activate the secret group")

	return cmd
}
