package org

import (
	"fmt"

	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/org"
	"github.com/Gkemhcs/kavach-cli/internal/types"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// validateRole checks if the provided role is valid for organization permissions.
// Valid roles are: admin, viewer, editor
func validateRole(role string) bool {
	return role == "admin" || role == "viewer" || role == "editor"
}

// NewGrantOrgCommand creates a Cobra command for granting organization access permissions.
// This command allows users to grant specific roles (admin, editor, viewer) to users or user groups
// for accessing organizations.
func NewGrantOrgCommand(logger *utils.Logger, orgClient org.OrgClient) *cobra.Command {
	var role string
	var userName string
	var userGroupName string

	cmd := &cobra.Command{
		Use:   "grant <organization-name>",
		Short: "🔑 Grant access permissions to an organization",
		Long: `Grant access permissions to an organization for a user or user group.

This command allows you to assign specific roles to users or user groups for accessing
organizations. The granted permissions will apply to the entire organization and all
its resources (secret groups, environments, etc.).

Key concepts:
- Permissions are granted at the organization level
- Users/groups get access to all resources within the organization
- Role bindings can be updated by granting the same role again
- Only organization owners and admins can grant permissions

Available roles:
- admin: Administrative access (manage resources, members, and grant permissions)
- editor: Basic access (view and use resources, create secret groups/environments)
- viewer: Read-only access (view resources only, cannot modify anything)

Permission hierarchy:
- admin > editor > viewer
- Higher roles inherit all permissions of lower roles

Prerequisites:
- You must be an admin or owner of the organization
- The user or group must exist in Kavach
- You must specify either --user or --group (not both)

Examples:
  kavach org grant "my-company" --user "john.doe" --role editor
  kavach org grant "startup" --group "developers" --role viewer
  kavach org grant "my-company" --user "sarah" --role admin

Note: If a user/group already has a role binding, granting a new role will
update their existing permissions. Use 'kavach org revoke' to remove permissions.`,
		Example: `  kavach org grant "my-company" --user "john.doe" --role editor
  kavach org grant "startup" --group "developers" --role viewer`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			orgName := args[0]

			logger.Info("Granting organization access permissions", map[string]interface{}{
				"cmd": "org grant",
				"org": orgName,
			})

			// Validate that only one of username or group name is provided
			if userGroupName != "" && userName != "" {
				fmt.Printf("\n🚨 Error: Cannot specify both user and group. Please provide either --user or --group, not both.\n")
				return nil
			}

			// Validate that at least one of username or group name is provided
			if userGroupName == "" && userName == "" {
				fmt.Printf("\n🚨 Error: Please specify either a user (--user) or a group (--group) to grant permissions to.\n")
				return nil
			}

			// Validate required parameters
			if role == "" {
				fmt.Printf("\n🚨 Error: Role is required. Please use --role to specify the permission level (admin, editor, or viewer).\n")
				return nil
			}

			if !validateRole(role) {
				fmt.Printf("\n🚨 Error: Invalid role '%s'. Valid roles are: admin, editor, viewer\n", role)
				return nil
			}

			// Prepare the role binding request
			req := types.GrantRoleBindingInput{
				UserName:  userName,
				GroupName: userGroupName,
				Role:      role,
				OrgName:   orgName,
			}

			// Execute the role binding grant
			err := orgClient.GrantRoleBinding(req)
			if err != nil {
				if err == cliErrors.ErrUnReachableBackend {
					fmt.Println("🚨 Error: Unable to connect to Kavach backend")
					fmt.Println("📡 This may be due to server downtime or network connectivity issues.")
					fmt.Println("📩 If this persists, please contact support at 👉 **gudikotieswarmani@gmail.com**")
					logger.Error("Backend unreachable during organization grant", err, map[string]interface{}{
						"cmd": "org grant",
						"org": orgName,
					})
					return nil
				}

				if err == cliErrors.ErrNotLoggedIn {
					fmt.Printf("\n🔒 Error: You are not logged in. Please run 'kavach login' to authenticate.\n")
					logger.Warn("User not logged in during organization grant", map[string]interface{}{
						"cmd": "org grant",
						"org": orgName,
					})
					return nil
				}

				if err == cliErrors.ErrDuplicateRoleBinding {
					fmt.Printf("\n⚠️  Warning: Role binding already exists for this user/group on this organization.\n")
					fmt.Printf("   The existing permissions have been updated to '%s'.\n", role)
					return nil
				}

				if err == cliErrors.ErrOrganizationNotFound {
					fmt.Printf("\n🚨 Error: Organization '%s' not found.\n", orgName)
					fmt.Printf("   Please verify the organization name or create it first.\n")
					return nil
				}

				if err == cliErrors.ErrAccessDenied {
					fmt.Printf("\n%s\n", err.Error())
					logger.Warn("Access denied during org grant", map[string]interface{}{
						"cmd": "org grant",
						"org": orgName,
					})
					return nil
				}

				logger.Error("Failed to grant organization permissions", err, map[string]interface{}{
					"cmd": "org grant",
					"org": orgName,
				})
				return err
			}

			// Success message
			target := userName
			if userName == "" {
				target = userGroupName
			}

			fmt.Printf("\n✅ Success: Granted '%s' permissions to '%s' on organization '%s'\n", role, target, orgName)

			return nil
		},
	}

	// Command flags with improved descriptions
	cmd.Flags().StringVarP(&role, "role", "r", "", "Permission level to grant (admin, editor, viewer)")
	cmd.Flags().StringVarP(&userName, "user", "u", "", "GitHub username to grant permissions to")
	cmd.Flags().StringVarP(&userGroupName, "group", "g", "", "User group name to grant permissions to")

	// Mark required flags
	cmd.MarkFlagRequired("role")

	return cmd
}
