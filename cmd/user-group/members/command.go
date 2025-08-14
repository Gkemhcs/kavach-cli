package members

import (
	"github.com/Gkemhcs/kavach-cli/internal/config"
	"github.com/Gkemhcs/kavach-cli/internal/groups"

	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewUserGroupMemberCommand creates a new command for managing user group members
func NewUserGroupMemberCommand(logger *utils.Logger, cfg *config.Config, userGroupMemberClient groups.UserGroupMemberClient) *cobra.Command {
	groupMemberCmd := &cobra.Command{
		Use:   "members",
		Short: "👥 Manage user group members",
		Long: `Manage members within user groups in Kavach - add, list, and remove members.

User group members are the individual users who belong to a user group. Managing
members allows you to control who has access to resources through group-based
permissions.

Key concepts:
- Members are individual users who belong to a user group
- Adding members to a group automatically grants them group-based permissions
- Removing members from a group revokes their group-based permissions
- Individual user permissions are not affected by group membership changes
- You can view all members of a group to understand current access

Available operations:
- add: Add a user to a user group
- list: List all members of a user group
- remove: Remove a user from a user group

Use cases:
- Onboarding new team members to existing groups
- Removing users who have left the team
- Auditing group membership for security compliance
- Managing access control through group membership

Examples:
  kavach user-group members add --group developers --user john
  kavach user-group members list --group developers
  kavach user-group members remove --group developers --user john`,
	}

	groupMemberCmd.AddCommand(
		NewAddGroupMemberCommand(logger, userGroupMemberClient),
		NewRemoveGroupMemberCommand(logger, userGroupMemberClient),
		NewListGroupMemberCommand(logger, userGroupMemberClient),
	)
	return groupMemberCmd
}
