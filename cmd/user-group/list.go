package usergroup

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
		"User Group Id",
		"User Group  Name",
		"Description",
		"Organization Name",
	}
}

func NewListUserGroupCommand(logger *utils.Logger, userGroupclient groups.UserGroupClient) *cobra.Command {
	var orgName string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "📋 List your user groups",
		Long: `List all user groups in the current organization.

This command displays a table of all user groups within the active organization,
showing basic information about each user group.

The output includes:
- User Group ID: Unique identifier for the user group
- User Group Name: Human-readable name of the user group
- Description: Optional description of the user group's purpose
- Organization Name: The organization this user group belongs to

Available roles:
- owner: Full administrative access (create, delete, manage members)
- admin: Administrative access (manage members and permissions)
- member: Basic access (view group and members)
- viewer: Read-only access (view group only)

Examples:
  kavach user-group list                    # List all user groups in current org
  kavach user-group list --help            # Show detailed help

Note: User groups help organize users for easier permission management.
Use 'kavach user-group members list <group-name>' to see members of a specific group.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			config, _ := config.LoadCLIConfig()
			if orgName == "" {
				if config.Organization == "" {
					fmt.Printf("\n⚠️  You haven't passed an organization and no default organization is set.\n")
					fmt.Printf("💡 Set a default organization using: `kavach org activate <org-name>`.\n")
					return nil
				}
				orgName = config.Organization

			}

			data, err := userGroupclient.ListUserGroups(orgName)

			if err != nil {
				if err == cliErrors.ErrUnReachableBackend {
					fmt.Println("🚨 Oops! Kavach backend is currently down or not responding.")
					fmt.Println("📡 This may be due to server downtime or high request volume.")
					fmt.Println("📩 If this persists, please drop us a message at 👉 **gudikotieswarmani@gmail.com**")
					return nil // or exit gracefully
				}
				if err == cliErrors.ErrNotLoggedIn {
					fmt.Printf("\n🔒 You are not logged in. Please run 'kavach login' to log in.\n")
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
	cmd.Flags().StringVarP(&orgName, "organization", "o", "", "Organisation under which we need to list the secret groups")
	return cmd
}

func ToRenderable(data []groups.ListGroupsByOrgRow, orgName string) [][]string {
	var out [][]string

	for _, userGroup := range data {

		out = append(out, []string{userGroup.ID, userGroup.Name, userGroup.Description, orgName})

	}
	return out
}
