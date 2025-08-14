package usergroup

import (
	"fmt"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/groups"

	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewDeleteUserGroupCommand creates a new command for deleting user groups
func NewDeleteUserGroupCommand(logger *utils.Logger, userGroupClient groups.UserGroupClient) *cobra.Command {
	var orgName string
	cmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "🗑️ Delete a user group by name",
		Long: `Delete a user group and remove all its members.

This command permanently deletes a user group and removes all users from it.
When a user group is deleted, all its members lose their group-based permissions
and role bindings associated with the group are removed.

⚠️  WARNING: This action is irreversible and will permanently delete the user group
and remove all its members. Make sure you have backups if needed.

Prerequisites:
- You must be the owner or admin of the user group
- The user group must not have any active role bindings
- You'll be prompted for confirmation before deletion

The deletion process:
1. Prompts for confirmation to prevent accidental deletion
2. Verifies you have permission to delete the user group
3. Removes all members from the user group
4. Removes all role bindings associated with the user group
5. Permanently deletes the user group
6. Confirms successful deletion

Examples:
  kavach user-group delete qa-team        # Delete user group (with confirmation)
  kavach user-group list                  # Verify user group is deleted

Note: When a user group is deleted, all its members lose their group-based
permissions. Individual user permissions are not affected.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			config, _ := config.LoadCLIConfig()
			if config.Organization == "" && orgName == "" {
				fmt.Printf("\n⚠️  You haven't passed an organization and no default organization is set.\n")
				fmt.Printf("💡 Set a default organization using: `kavach org activate <org-name>`.\n")
				return nil
			}
			if orgName == "" {
				orgName = config.Organization
			}
			msg := fmt.Sprintf("are you sure to delete user group %s under org  %s if yes click on proceed otherwise cancel", name, orgName)
			if !utils.ConfirmSecretGroupCreation(msg) {
				fmt.Print("\n cancelled the delete operation \n")
				return nil
			}

			err := userGroupClient.DeleteUserGroup(orgName, name)
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
				case cliErrors.ErrUserGroupNotFound:
					fmt.Printf("\n❌ User group '%s' does not exist.\n", name)
					return nil
				case cliErrors.ErrAccessDenied:
					fmt.Printf("\n%s\n", err.Error())
					return nil
				default:
					return err
				}

			}
			fmt.Printf("\n🗑️ UserGroup '%s' deleted successfully.\n", name)
			return nil
		},
	}
	cmd.Flags().StringVarP(&orgName, "organization", "o", "", "Organisation under which we need to delete the secret group")
	return cmd
}
