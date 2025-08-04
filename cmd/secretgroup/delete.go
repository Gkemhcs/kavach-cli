package secretgroup

import (
	"fmt"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"

	"github.com/Gkemhcs/kavach-cli/internal/secretgroup"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

func NewDeleteSecretGroupCommand(logger *utils.Logger, groupClient secretgroup.SecretGroupClient) *cobra.Command {
	var orgName string
	cmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "🗑️ Delete a secret group by name",
		Long: `Delete a secret group and all its associated resources.

This command permanently deletes a secret group and all its contents including:
- All environments within the secret group
- All secrets stored in those environments
- All user groups and member associations
- All role bindings and permissions

⚠️  WARNING: This action is irreversible and will permanently delete all data
associated with the secret group. Make sure you have backups if needed.

Prerequisites:
- You must be the owner or admin of the secret group
- The secret group must not have any child resources (environments)
- You'll be prompted for confirmation before deletion

The deletion process:
1. Prompts for confirmation to prevent accidental deletion
2. Verifies you have permission to delete the secret group
3. Checks for existing child resources (environments)
4. Permanently deletes the secret group and all its data
5. Confirms successful deletion

Examples:
  kavach group delete myapp        # Delete secret group (with confirmation)
  kavach group list                # Verify secret group is deleted

Note: If the secret group contains environments, you'll need to delete
those first before you can delete the secret group.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			config, _ := config.LoadCLIConfig()
			if config.Organization == "" && orgName == "" {
				fmt.Println("\n⚠️  You haven't passed an organization and no default organization is set.")
				fmt.Println("💡 Set a default organization using: `kavach org activate <org-name>`")
				return nil
			}
			if orgName == "" {
				orgName = config.Organization
			}
			msg := fmt.Sprintf("are you sure to delete secret group %s under org  %s if yes click on proceed otherwise cancel", name, orgName)
			if !utils.ConfirmSecretGroupCreation(msg) {
				fmt.Print("\n cancelled the delete operation \n")
				return nil
			}

			err := groupClient.DeleteSecretGroupByName(orgName, name)
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
					fmt.Printf("\n❌ SecretGroup '%s' does not exist.\n", name)
					return nil
				case cliErrors.ErrForeignKeyViolation:
					fmt.Println("🚨 cannot delete secret group  as it contain child resources like environments and secrets")
					fmt.Printf("\n 🚨 first delete all child resources to delete secret group %s \n", name)
					return nil
				case cliErrors.ErrAccessDenied:
					fmt.Printf("\n%s\n", err.Error())
					return nil
				default:
					return err
				}

			}
			fmt.Printf("\n🗑️ SecretGroup '%s' deleted successfully.\n", name)
			return nil
		},
	}
	cmd.Flags().StringVarP(&orgName, "organization", "o", "", "Organisation under which we need to delete the secret group")
	return cmd
}
