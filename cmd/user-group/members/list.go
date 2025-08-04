package members

import (
	"fmt"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/groups"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

func GetUserGroupHeaders() []string {
	return []string{
		"UserEmail",
		"UserName",
		"UserGroupName",
	}
}

func NewListGroupMemberCommand(logger *utils.Logger, userGroupMemberClient groups.UserGroupMemberClient) *cobra.Command {

	var groupName string
	var orgName string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "📋 List members of a user group",
		Long: `List all members of a user group to see who has group-based access.

This command displays a table of all users who are members of the specified
user group. This helps you understand who has access to resources through
group-based permissions.

Key concepts:
- Shows all users who belong to the specified user group
- Members gain access to resources where the group has permissions
- Useful for auditing group membership for security compliance
- Helps verify access control and team organization

The output includes:
- User Email: The user's email address
- User Name: The user's display name
- User Group Name: The name of the user group

Use cases:
- Auditing group membership for security compliance
- Verifying team access to resources
- Understanding who has group-based permissions
- Planning access control changes

Examples:
  kavach user-group members list --user-group developers
  kavach user-group members list --user-group qa-team
  kavach user-group members list --user-group admins

Note: This shows only the members of the specified group. To see all user
groups, use 'kavach user-group list'.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			if groupName == "" {
				fmt.Printf("please enter the usergroupname in which you want to list member")
				return nil
			}
			config, _ := config.LoadCLIConfig()

			if config.Organization == "" && orgName == "" {
				fmt.Println("\n⚠️  You haven't passed an organization and no default organization is set.")
				fmt.Println("💡 Set a default organization using: `kavach org activate <org-name>`")
				return nil
			}
			if orgName == "" {
				fmt.Printf("\n⚠️  You haven't passed an organization and no default organization is set.")
				fmt.Printf("💡 Set a default organization using: `kavach org activate <org-name>`")
				return nil
			}

			if orgName == "" {
				orgName = config.Organization
			}

			data, err := userGroupMemberClient.ListUserGroupMembers(orgName, groupName)
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
				if err == cliErrors.ErrAccessDenied {
					fmt.Printf("\n%s\n", err.Error())
					return nil
				}

				return err
			}
			utils.RenderTable(GetUserGroupHeaders(), ToRenderable(data, orgName))
			return nil
		},
	}
	cmd.Flags().StringVarP(&groupName, "user-group", "g", "", "name of UserGroup")

	cmd.Flags().StringVarP(&orgName, "organization", "o", "", "Organisation under which we need to create the secret group")
	return cmd

}

func ToRenderable(data []groups.ListGroupMembersRow, groupName string) [][]string {
	var out [][]string

	for _, user := range data {

		out = append(out, []string{user.Email, user.Name, groupName})

	}
	return out
}
