package members

import (
	"fmt"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/groups"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

func NewRemoveGroupMemberCommand(logger *utils.Logger, userGroupMemberClient groups.UserGroupMemberClient) *cobra.Command {

	var groupName string
	var orgName string
	var userName string

	cmd := &cobra.Command{
		Use:   "remove",
		Short: "➖ Remove a user from a user group",
		Long: `Remove a user from a user group to revoke their group-based permissions.

When you remove a user from a user group, they lose access to all resources
where the group has been granted permissions. This is useful for revoking
access when users leave teams or change roles.

Key concepts:
- Removing a user from a group revokes all group-based permissions
- The user loses access to resources at organization, secret group, and environment levels
- Individual user permissions are preserved and remain unaffected
- The user can still access resources through their individual permissions

Prerequisites:
- You must be an admin or owner of the user group
- The user must be a current member of the group
- The user must exist in Kavach

Use cases:
- Removing users who have left the team
- Revoking access when users change roles
- Removing users from role-based groups
- Managing team access through group membership

Examples:
  kavach user-group members remove --user-group developers --user john
  kavach user-group members remove --user-group qa-team --user sarah
  kavach user-group members remove --user-group admins --user mike

Note: After removing a user from a group, they will lose access to all resources
where the group has been granted permissions. Individual user permissions are
not affected. Use 'kavach user-group members list' to verify the removal.`,
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			if userName == "" {
				fmt.Printf("⚠️  Please provide the --user argument (username) to remove a member.\n")
				return nil
			}
			if groupName == "" {
				fmt.Printf("⚠️  Please provide the --user-group argument (user group name) to remove a member.\n")
				return nil
			}
			config, _ := config.LoadCLIConfig()

			if config.Organization == "" && orgName == "" {
				fmt.Println("\n⚠️  You haven't passed an organization and no default organization is set.")
				fmt.Println("💡 Set a default organization using: `kavach org activate <org-name>`")
				return nil
			}
			if orgName == "" {
				msg := fmt.Sprintf("you havent passed organization option so we are adding to  user  group %s under active  organization %s", groupName, config.Organization)

				cont := utils.ConfirmSecretGroupCreation(msg)
				if !cont {
					fmt.Print("\n exiting \n")
					return nil
				}

			}

			if orgName == "" {
				orgName = config.Organization
			}

			err := userGroupMemberClient.RemoveGroupMember(orgName, userName, groupName)
			if err != nil {
				// Only return error if it's not a handled/user-friendly case
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
					fmt.Printf("\n❌ Organization '%s' does not exist.\n", orgName)
					return nil
				}
				if err == cliErrors.ErrUserGroupNotFound {
					fmt.Printf("\n ❌ User group %s doesnt exist in organization %s \n", groupName, orgName)
					return nil
				}
				if err == cliErrors.ErrUserNotFound {
					fmt.Printf("\n❌ User '%s' not found in the organization.\n", userName)
					return nil
				}
				if err == cliErrors.ErrMemberNotFound {
					fmt.Printf("\n❌ User '%s' is not a member of user group '%s'.\n", userName, groupName)
					return nil
				}
				if err == cliErrors.ErrAccessDenied {
					fmt.Printf("\n%s\n", err.Error())
					return nil
				}
				return err
			}
			fmt.Printf("\n🎉 Member '%s' removed successfully from user group '%s'.\n", userName, groupName)

			return nil
		},
	}
	cmd.Flags().StringVarP(&groupName, "user-group", "g", "", "name of UserGroup")
	cmd.Flags().StringVarP(&userName, "user", "u", "", " name of user you want to add tho this group")
	cmd.Flags().StringVarP(&orgName, "organization", "o", "", "Organisation under which we need to create the secret group")
	return cmd

}
