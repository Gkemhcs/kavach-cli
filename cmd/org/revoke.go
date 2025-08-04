package org

import (
	"fmt"

	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/org"
	"github.com/Gkemhcs/kavach-cli/internal/types"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewRevokeOrgCommand creates a Cobra command for revoking organization access permissions.
// This command allows users to revoke specific roles (admin, editor, viewer) from users or user groups
// for accessing organizations.
func NewRevokeOrgCommand(logger *utils.Logger, orgClient org.OrgClient) *cobra.Command {
	var role string
	var userName string
	var userGroupName string

	cmd := &cobra.Command{
		Use:   "revoke <organization-name>",
		Short: "🚫 Revoke access permissions from an organization",
		Long: `Revoke access permissions from an organization for a user or user group.

This command allows you to remove specific roles from users or user groups for accessing
organizations. When you revoke permissions, the user/group will lose access to the
organization and all its resources (secret groups, environments, etc.).

Key concepts:
- Permissions are revoked at the organization level
- Users/groups lose access to all resources within the organization
- You can revoke specific roles while keeping other roles intact
- Only organization owners and admins can revoke permissions

Available roles to revoke:
- admin: Administrative access (manage resources, members, and grant permissions)
- editor: Basic access (view and use resources, create secret groups/environments)
- viewer: Read-only access (view resources only, cannot modify anything)

Important considerations:
- Revoking permissions is immediate and affects all resources in the organization
- Users may lose access to secrets and configurations they were working with
- Consider the impact on ongoing work before revoking permissions
- You cannot revoke your own permissions (as an organization owner)

Prerequisites:
- You must be an admin or owner of the organization
- The user or group must have the specified role binding
- You must specify either --user or --group (not both)

Examples:
  kavach org revoke "my-company" --user "john.doe" --role editor
  kavach org revoke "startup" --group "developers" --role viewer
  kavach org revoke "my-company" --user "sarah" --role admin

Note: If a user/group doesn't have the specified role binding, the command will
show a warning but won't fail. Use 'kavach org grant' to add new permissions.`,
		Example: `  kavach org revoke "my-company" --user "john.doe" --role editor
  kavach org revoke "startup" --group "developers" --role viewer`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			orgName := args[0]

			logger.Info("Revoking organization access permissions", map[string]interface{}{
				"cmd": "org revoke",
				"org": orgName,
			})

			// Validate that only one of username or group name is provided
			if userGroupName != "" && userName != "" {
				fmt.Printf("\n🚨 Error: Cannot specify both user and group. Please provide either --user or --group, not both.\n")
				return nil
			}

			// Validate that at least one of username or group name is provided
			if userGroupName == "" && userName == "" {
				fmt.Printf("\n🚨 Error: Please specify either a user (--user) or a group (--group) to revoke permissions from.\n")
				return nil
			}

			// Validate required parameters
			if role == "" {
				fmt.Printf("\n🚨 Error: Role is required. Please use --role to specify the permission level to revoke (admin, editor, or viewer).\n")
				return nil
			}

			if !validateRole(role) {
				fmt.Printf("\n🚨 Error: Invalid role '%s'. Valid roles are: admin, editor, viewer\n", role)
				return nil
			}

			// Prepare the role binding revocation request
			req := types.RevokeRoleBindingInput{
				UserName:  userName,
				GroupName: userGroupName,
				Role:      role,
				OrgName:   orgName,
			}

			// Execute the role binding revocation
			err := orgClient.RevokeRoleBinding(req)
			if err != nil {
				if err == cliErrors.ErrUnReachableBackend {
					fmt.Println("🚨 Error: Unable to connect to Kavach backend")
					fmt.Println("📡 This may be due to server downtime or network connectivity issues.")
					fmt.Println("📩 If this persists, please contact support at 👉 **gudikotieswarmani@gmail.com**")
					logger.Error("Backend unreachable during organization revoke", err, map[string]interface{}{
						"cmd": "org revoke",
						"org": orgName,
					})
					return nil
				}

				if err == cliErrors.ErrNotLoggedIn {
					fmt.Printf("\n🔒 Error: You are not logged in. Please run 'kavach login' to authenticate.\n")
					logger.Warn("User not logged in during organization revoke", map[string]interface{}{
						"cmd": "org revoke",
						"org": orgName,
					})
					return nil
				}

				if err == cliErrors.ErrRoleBindingNotFound {
					fmt.Printf("\n⚠️  Warning: No role binding found for the specified user/group on this organization.\n")
					fmt.Printf("   The user/group may not have had '%s' permissions to begin with.\n", role)
					return nil
				}

				if err == cliErrors.ErrOrganizationNotFound {
					fmt.Printf("\n🚨 Error: Organization '%s' not found.\n", orgName)
					fmt.Printf("   Please verify the organization name.\n")
					return nil
				}

				if err == cliErrors.ErrAccessDenied {
					fmt.Printf("\n%s\n", err.Error())
					logger.Warn("Access denied during org revoke", map[string]interface{}{
						"cmd": "org revoke",
						"org": orgName,
					})
					return nil
				}

				logger.Error("Failed to revoke organization permissions", err, map[string]interface{}{
					"cmd": "org revoke",
					"org": orgName,
				})
				return err
			}

			// Success message
			target := userName
			if userName == "" {
				target = userGroupName
			}

			fmt.Printf("\n✅ Success: Revoked '%s' permissions from '%s' on organization '%s'\n", role, target, orgName)

			return nil
		},
	}

	// Command flags with improved descriptions
	cmd.Flags().StringVarP(&role, "role", "r", "", "Permission level to revoke (admin, editor, viewer)")
	cmd.Flags().StringVarP(&userName, "user", "u", "", "GitHub username to revoke permissions from")
	cmd.Flags().StringVarP(&userGroupName, "group", "g", "", "User group name to revoke permissions from")

	// Mark required flags
	cmd.MarkFlagRequired("role")

	return cmd
}
