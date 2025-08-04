package usergroup

import (
	"github.com/Gkemhcs/kavach-cli/cmd/user-group/members"
	"github.com/Gkemhcs/kavach-cli/internal/config"
	"github.com/Gkemhcs/kavach-cli/internal/groups"

	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

func NewUserGroupCommand(logger *utils.Logger, cfg *config.Config, userGroupClient groups.UserGroupClient, userGroupMemberClient groups.UserGroupMemberClient) *cobra.Command {
	userGroupCmd := &cobra.Command{
		Use:   "user-group",
		Short: "👥 Manage user groups",
		Long: `Manage user groups in Kavach - create, list, and delete user groups.

User groups are collections of users that can be assigned permissions together.
Instead of managing permissions for each user individually, you can create user
groups and assign permissions to the entire group at once.

Key concepts:
- User groups belong to organizations and can contain multiple users
- User groups can be assigned roles at organization, secret group, or environment levels
- User groups simplify permission management for teams and departments
- You can add or remove users from groups without changing their permissions
- User groups help organize users by team, department, or role

Resource hierarchy:
Organization → User Groups → Users

Available operations:
- create: Create a new user group within an organization
- list: List all user groups in the current organization
- delete: Delete a user group and remove all its members
- members: Manage members within user groups (add, list, remove)

Use cases:
- Organize users by team (e.g., "developers", "qa-team", "ops-team")
- Group users by department (e.g., "engineering", "marketing", "finance")
- Create role-based groups (e.g., "admins", "viewers", "editors")
- Simplify permission management for large teams

Examples:
  kavach user-group create developers --description "Development team"
  kavach user-group list                    # List all user groups
  kavach user-group delete qa-team          # Delete user group (with confirmation)
  kavach user-group members add --group developers --user john`,
	}

	userGroupCmd.AddCommand(
		NewCreateUserGroupCommand(logger, userGroupClient),
		NewDeleteUserGroupCommand(logger, userGroupClient),
		NewListUserGroupCommand(logger, userGroupClient),
		members.NewUserGroupMemberCommand(logger, cfg, userGroupMemberClient),
	)
	return userGroupCmd
}
