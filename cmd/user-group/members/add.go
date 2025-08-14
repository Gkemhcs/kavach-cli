package members

import (
	"fmt"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/groups"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewAddGroupMemberCommand creates a new command for adding members to user groups
func NewAddGroupMemberCommand(logger *utils.Logger, userGroupMemberClient groups.UserGroupMemberClient) *cobra.Command {

	var groupName string
	var orgName string
	var userName string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "➕ Add a user to a user group",
		Long: `Add a user to a user group to grant them group-based permissions.

When you add a user to a user group, they automatically gain access to all
resources where the group has been granted permissions. This is an efficient
way to manage access control for multiple users at once.

Key concepts:
- Adding a user to a group grants them all group-based permissions
- The user gains access to resources at organization, secret group, and environment levels
- Individual user permissions are preserved and work alongside group permissions
- You can add users to multiple groups for different permission sets

Prerequisites:
- You must be an admin or owner of the user group
- The user must exist in Kavach
- The user must not already be a member of the group

Use cases:
- Onboarding new team members to existing groups
- Granting access to new team members
- Adding users to role-based groups (admins, viewers, etc.)
- Managing team access through group membership

Examples:
  kavach user-group members add --user-group developers --user john
  kavach user-group members add --user-group qa-team --user sarah
  kavach user-group members add --user-group admins --user mike

Note: After adding a user to a group, they will have access to all resources
where the group has been granted permissions. Use 'kavach user-group members list'
to verify the addition.`,
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			if userName == "" {
				fmt.Printf("⚠️  Please provide the --user argument (username) to add a member.\n")
				return nil
			}
			if groupName == "" {
				fmt.Printf("⚠️  Please provide the --user-group argument (user group name) to add a member.\n")
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

			err := userGroupMemberClient.AddGroupMember(orgName, userName, groupName)
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
				if err == cliErrors.ErrDuplicateMember {
					fmt.Printf("\n❌ User '%s' is already a member of user group '%s'.\n", userName, groupName)
					return nil
				}
				if err == cliErrors.ErrUserNotFound {
					fmt.Printf("\n❌ User '%s' not found in the organization.\n", userName)
					return nil
				}
				if err == cliErrors.ErrAccessDenied {
					fmt.Printf("\n%s\n", err.Error())
					return nil
				}
				return err
			}
			fmt.Printf("\n🎉 Member '%s' added successfully to user group '%s'.\n", userName, groupName)

			return nil
		},
	}
	cmd.Flags().StringVarP(&groupName, "user-group", "g", "", "name of UserGroup")
	cmd.Flags().StringVarP(&userName, "user", "u", "", " name of user you want to add tho this group")
	cmd.Flags().StringVarP(&orgName, "organization", "o", "", "Organisation under which we need to create the secret group")
	return cmd

}
