package usergroup

import (
	"fmt"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/groups"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewCreateUserGroupCommand creates a new command for creating user groups
func NewCreateUserGroupCommand(logger *utils.Logger, userGroupClient groups.UserGroupClient) *cobra.Command {

	var description string
	var orgName string
	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "🏗️ Create a new user group",
		Long: `Create a new user group within the current organization.

User groups are collections of users that can be assigned permissions together.
Instead of managing permissions for each user individually, you can create user
groups and assign permissions to the entire group at once.

Key features:
- You become the owner of the created user group
- User group names must be unique within the organization
- User groups can contain multiple users
- You can invite other users and assign different roles
- User groups help organize users by team, department, or role

Available roles for user group members:
- owner: Full administrative access (create, delete, manage members)
- admin: Administrative access (manage members and permissions)
- member: Basic access (view group and members)
- viewer: Read-only access (view group only)

Use cases:
- Organize users by team (e.g., "developers", "qa-team", "ops-team")
- Group users by department (e.g., "engineering", "marketing", "finance")
- Create role-based groups (e.g., "admins", "viewers", "editors")
- Simplify permission management for large teams

Examples:
  kavach user-group create developers --description "Development team"
  kavach user-group create qa-team --description "QA team"
  kavach user-group create developers                    # Without description

Note: User group names should be descriptive and follow your naming conventions.
Once created, you can add members using 'kavach user-group members add'.`,
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
				msg := fmt.Sprintf("you havent passed organization option so we are creating the user  group %s under active  organization %s", name, config.Organization)

				cont := utils.ConfirmSecretGroupCreation(msg)
				if !cont {
					fmt.Print("\n exiting \n")
					return nil
				}

			}

			if orgName == "" {
				orgName = config.Organization
			}

			err := userGroupClient.CreateUserGroup(orgName, name, description)
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
				if err == cliErrors.ErrDuplicateUserGroup {
					fmt.Printf("\n❌ User group '%s' already exists in organization '%s'.\n", name, orgName)
					return nil
				}
				if err == cliErrors.ErrAccessDenied {
					fmt.Printf("\n%s\n", err.Error())
					return nil
				}

				return err
			}
			fmt.Printf("\n 🎉 User group '%s' created successfully under organization '%s'.\n", name, orgName)

			return nil
		},
	}
	cmd.Flags().StringVar(&description, "description", "", "Description of the secret group")
	cmd.Flags().StringVarP(&orgName, "organization", "o", "", "Organisation under which we need to create the secret group")
	return cmd

}
