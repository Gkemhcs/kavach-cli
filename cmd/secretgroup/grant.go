package secretgroup

import (
	"fmt"

	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/secretgroup"
	"github.com/Gkemhcs/kavach-cli/internal/types"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// validateRole checks if the provided role is valid for secret group permissions.
// Valid roles are: admin, viewer, editor
func validateRole(role string) bool {
	return role == "admin" || role == "viewer" || role == "editor"
}

// NewGrantSecretGroupCommand creates a Cobra command for granting secret group access permissions.
// This command allows users to grant specific roles (admin, editor, viewer) to users or user groups
// for accessing secret groups within organizations.
func NewGrantSecretGroupCommand(logger *utils.Logger, secretGroupClient secretgroup.SecretGroupClient) *cobra.Command {
	var role string
	var userName string
	var userGroupName string
	var orgName string

	cmd := &cobra.Command{
		Use:   "grant <secret-group-name>",
		Short: "🔑 Grant access permissions to a secret group",
		Long: `Grant access permissions to a secret group for a user or user group.

This command allows you to assign specific roles to users or user groups for accessing
secret groups. The granted permissions will apply to the entire secret group and all
its environments and secrets.

Key concepts:
- Permissions are granted at the secret group level
- Users/groups get access to all environments and secrets within the secret group
- Role bindings can be updated by granting the same role again
- Only secret group owners and admins can grant permissions

Available roles:
- admin: Administrative access (manage environments, secrets, and grant permissions)
- editor: Basic access (view and use secrets, create environments)
- viewer: Read-only access (view secrets only, cannot modify anything)

Permission hierarchy:
- admin > editor > viewer
- Higher roles inherit all permissions of lower roles

Prerequisites:
- You must be an admin or owner of the secret group
- The user or group must exist in Kavach
- You must specify either --user or --group (not both)
- Organization must be specified with --org

Examples:
  kavach group grant "backend-secrets" --user "john.doe" --role editor --org "my-company"
  kavach group grant "frontend-secrets" --group "developers" --role viewer --org "startup"
  kavach group grant "myapp" --user "sarah" --role admin --org "my-company"

Note: If a user/group already has a role binding, granting a new role will
update their existing permissions. Use 'kavach group revoke' to remove permissions.`,
		Example: `  kavach group grant "backend-secrets" --user "john.doe" --role editor --org "my-company"
  kavach group grant "frontend-secrets" --group "developers" --role viewer --org "startup"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			secretGroupName := args[0]

			logger.Info("Granting secret group access permissions", map[string]interface{}{
				"cmd":         "group grant",
				"secretGroup": secretGroupName,
				"org":         orgName,
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

			if orgName == "" {
				fmt.Printf("\n🚨 Error: Organization name is required. Please use --org to specify the organization.\n")
				return nil
			}

			if !validateRole(role) {
				fmt.Printf("\n🚨 Error: Invalid role '%s'. Valid roles are: admin, editor, viewer\n", role)
				return nil
			}

			// Prepare the role binding request
			req := types.GrantRoleBindingInput{
				UserName:        userName,
				GroupName:       userGroupName,
				Role:            role,
				OrgName:         orgName,
				SecretGroupName: secretGroupName,
			}

			// Execute the role binding grant
			err := secretGroupClient.GrantRoleBinding(req)
			if err != nil {
				if err == cliErrors.ErrUnReachableBackend {
					fmt.Println("🚨 Error: Unable to connect to Kavach backend")
					fmt.Println("📡 This may be due to server downtime or network connectivity issues.")
					fmt.Println("📩 If this persists, please contact support at 👉 **gudikotieswarmani@gmail.com**")
					logger.Error("Backend unreachable during secret group grant", err, map[string]interface{}{
						"cmd":         "group grant",
						"secretGroup": secretGroupName,
						"org":         orgName,
					})
					return nil
				}

				if err == cliErrors.ErrNotLoggedIn {
					fmt.Printf("\n🔒 Error: You are not logged in. Please run 'kavach login' to authenticate.\n")
					logger.Warn("User not logged in during secret group grant", map[string]interface{}{
						"cmd":         "group grant",
						"secretGroup": secretGroupName,
						"org":         orgName,
					})
					return nil
				}

				if err == cliErrors.ErrDuplicateRoleBinding {
					fmt.Printf("\n⚠️  Warning: Role binding already exists for this user/group on this secret group.\n")
					fmt.Printf("   The existing permissions have been updated to '%s'.\n", role)
					return nil
				}

				if err == cliErrors.ErrOrganizationNotFound {
					fmt.Printf("\n🚨 Error: Organization '%s' not found.\n", orgName)
					fmt.Printf("   Please verify the organization name or create it first.\n")
					return nil
				}

				if err == cliErrors.ErrSecretGroupNotFound {
					fmt.Printf("\n🚨 Error: Secret group '%s' not found in organization '%s'.\n", secretGroupName, orgName)
					fmt.Printf("   Please verify the secret group name or create it first.\n")
					return nil
				}

				if err == cliErrors.ErrAccessDenied {
					fmt.Printf("\n%s\n", err.Error())
					logger.Warn("Access denied during secret group grant", map[string]interface{}{
						"cmd":         "group grant",
						"secretGroup": secretGroupName,
						"org":         orgName,
					})
					return nil
				}

				logger.Error("Failed to grant secret group permissions", err, map[string]interface{}{
					"cmd":         "group grant",
					"secretGroup": secretGroupName,
					"org":         orgName,
				})
				return err
			}

			// Success message
			target := userName
			if userName == "" {
				target = userGroupName
			}

			fmt.Printf("\n✅ Success: Granted '%s' permissions to '%s' on secret group '%s'\n", role, target, secretGroupName)
			fmt.Printf("   Organization: %s\n", orgName)

			return nil
		},
	}

	// Command flags with improved descriptions
	cmd.Flags().StringVarP(&role, "role", "r", "", "Permission level to grant (admin, editor, viewer)")
	cmd.Flags().StringVarP(&userName, "user", "u", "", "GitHub username to grant permissions to")
	cmd.Flags().StringVarP(&userGroupName, "group", "g", "", "User group name to grant permissions to")
	cmd.Flags().StringVarP(&orgName, "org", "o", "", "Organization name where the secret group exists")

	// Mark required flags
	cmd.MarkFlagRequired("role")
	cmd.MarkFlagRequired("org")

	return cmd
}
