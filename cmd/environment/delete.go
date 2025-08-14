package environment

import (
	"fmt"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	"github.com/Gkemhcs/kavach-cli/internal/environment"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"

	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewDeleteEnvironmentCommand creates a new command for deleting environments
func NewDeleteEnvironmentCommand(logger *utils.Logger, envClient environment.EnvironmentClient) *cobra.Command {
	var orgName string
	var groupName string
	cmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "🗑️ Delete an environment by name",
		Long: `Delete an environment and all its associated resources.

This command permanently deletes an environment and all its contents including:
- All secrets stored in the environment
- All provider configurations
- All user groups and member associations
- All role bindings and permissions

⚠️  WARNING: This action is irreversible and will permanently delete all data
associated with the environment. Make sure you have backups if needed.

Prerequisites:
- You must be the owner or admin of the environment
- The environment must not have any child resources (secrets)
- You'll be prompted for confirmation before deletion

The deletion process:
1. Prompts for confirmation to prevent accidental deletion
2. Verifies you have permission to delete the environment
3. Checks for existing child resources (secrets)
4. Permanently deletes the environment and all its data
5. Confirms successful deletion

Examples:
  kavach env delete testing        # Delete environment (with confirmation)
  kavach env list                  # Verify environment is deleted

Note: If the environment contains secrets, you'll need to delete
those first before you can delete the environment.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			cfg, _ := config.LoadCLIConfig()
			if cfg.Organization == "" && cfg.SecretGroup == "" && orgName == "" && groupName == "" {

				fmt.Println("\n⚠️  You haven't passed an organization and secret group and no default organization  and secret group are set.")
				fmt.Println("💡 Set a default organization using: `kavach org activate <org-name>` and `kavach group activate <secret-group-name>`")
				return nil

			}
			if orgName == "" {
				orgName = cfg.Organization
			}
			if groupName == "" {
				groupName = cfg.SecretGroup
			}

			msg := fmt.Sprintf("are you sure to delete  the environment %s under secret group %s under org  %s if yes click on proceed otherwise cancel", name, groupName, orgName)
			if !utils.ConfirmSecretGroupCreation(msg) {
				fmt.Print("\n cancelled the delete operation \n")
				return nil
			}

			err := envClient.DeleteEnvironment(orgName, groupName, name)
			if err != nil {
				switch err {
				case cliErrors.ErrUnReachableBackend:
					fmt.Println("🚨 Oops! Kavach backend is currently down or not responding.")
					fmt.Println("📡 This may be due to server downtime or high request volume.")
					fmt.Println("📩 If this persists, please drop us a message at 👉 **gudikotieswarmani@gmail.com**")
					return nil // or exit gracefully
				case cliErrors.ErrNotLoggedIn:
					fmt.Printf("\n🔒 You are not logged in. Please run 'kavach login' to log in.\n")
					return nil
				case cliErrors.ErrOrganizationNotFound:
					fmt.Printf("\n❌ Organization '%s' does not exist.\n", orgName)
					return nil
				case cliErrors.ErrSecretGroupNotFound:
					fmt.Printf("\n❌ SecretGroup '%s' does not exist.\n", groupName)
					return nil
				case cliErrors.ErrEnvironmentNotFound:
					fmt.Printf("\n❌ Environment '%s' does not exist.\n", name)
					return nil
				case cliErrors.ErrAccessDenied:
					fmt.Printf("\n%s\n", err.Error())
					logger.Warn("Access denied during environment delete", map[string]interface{}{"cmd": "env delete", "env": name, "group": groupName, "org": orgName})
					return nil
				default:
					// Check if the error message contains authentication-related text
					if cliErrors.IsAuthenticationError(err) {
						fmt.Printf("\n🔑 Please login again, the session is expired, unable to authenticate you\n")
						logger.Warn("Authentication error during environment delete", map[string]interface{}{"cmd": "env delete", "env": name, "org": orgName, "group": groupName})
						return nil
					}
					return err
				}

			}
			fmt.Printf("\n🗑️ Environment %s deleted successfully.\n", name)
			return nil
		},
	}
	cmd.Flags().StringVarP(&orgName, "organization", "o", "", "Organisation under which we need to delete the environment")
	cmd.Flags().StringVarP(&groupName, "group", "g", "", "secret group under which you want to delete the environment")
	return cmd
}
