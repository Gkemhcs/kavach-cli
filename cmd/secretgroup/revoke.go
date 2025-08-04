package secretgroup

import (
	"fmt"

	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/secretgroup"
	"github.com/Gkemhcs/kavach-cli/internal/types"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewRevokeSecretGroupCommand creates a Cobra command for revoking secret group access permissions.
// This command allows users to revoke specific roles (admin, editor, viewer) from users or user groups
// for accessing secret groups within organizations.
func NewRevokeSecretGroupCommand(logger *utils.Logger, secretGroupClient secretgroup.SecretGroupClient) *cobra.Command {
	var role string
	var userName string
	var userGroupName string
	var orgName string

	cmd := &cobra.Command{
		Use:   "revoke <secret-group-name>",
		Short: "🚫 Revoke access permissions from a secret group",
		Long: `Revoke access permissions from a secret group for a user or user group.

This command allows you to remove specific roles from users or user groups for accessing
secret groups. When you revoke permissions, the user/group will lose access to the
secret group and all its environments and secrets.

Key concepts:
- Permissions are revoked at the secret group level
- Users/groups lose access to all environments and secrets within the secret group
- You can revoke specific roles while keeping other roles intact
- Only secret group owners and admins can revoke permissions

Available roles to revoke:
- admin: Administrative access (manage environments, secrets, and grant permissions)
- editor: Basic access (view and use secrets, create environments)
- viewer: Read-only access (view secrets only, cannot modify anything)

Important considerations:
- Revoking permissions is immediate and affects all resources in the secret group
- Users may lose access to secrets and configurations they were working with
- Consider the impact on ongoing work before revoking permissions
- You cannot revoke your own permissions (as a secret group owner)

Prerequisites:
- You must be an admin or owner of the secret group
- The user or group must have the specified role binding
- You must specify either --user or --group (not both)
- Organization must be specified with --org

Examples:
  kavach group revoke "backend-secrets" --user "john.doe" --role editor --org "my-company"
  kavach group revoke "frontend-secrets" --group "developers" --role viewer --org "startup"
  kavach group revoke "myapp" --user "sarah" --role admin --org "my-company"

Note: If a user/group doesn't have the specified role binding, the command will
show a warning but won't fail. Use 'kavach group grant' to add new permissions.`,
		Example: `kavach group revoke "backend-secrets" --user "john.doe" --role editor --org "my-company"
kavach group revoke "frontend-secrets" --group "developers" --role viewer --org "startup"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			secretGroupName := args[0]

			logger.Info("Revoking secret group access permissions", map[string]interface{}{
				"cmd":         "group revoke",
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
				fmt.Printf("\n🚨 Error: Please specify either a user (--user) or a group (--group) to revoke permissions from.\n")
				return nil
			}

			// Validate required parameters
			if role == "" {
				fmt.Printf("\n🚨 Error: Role is required. Please use --role to specify the permission level to revoke (admin, editor, or viewer).\n")
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

			// Prepare the role binding revocation request
			req := types.RevokeRoleBindingInput{
				UserName:        userName,
				GroupName:       userGroupName,
				Role:            role,
				OrgName:         orgName,
				SecretGroupName: secretGroupName,
			}

			// Execute the role binding revocation
			err := secretGroupClient.RevokeRoleBinding(req)
			if err != nil {
				if err == cliErrors.ErrUnReachableBackend {
					fmt.Println("🚨 Error: Unable to connect to Kavach backend")
					fmt.Println("📡 This may be due to server downtime or network connectivity issues.")
					fmt.Println("📩 If this persists, please contact support at 👉 **gudikotieswarmani@gmail.com**")
					logger.Error("Backend unreachable during secret group revoke", err, map[string]interface{}{
						"cmd":         "group revoke",
						"secretGroup": secretGroupName,
						"org":         orgName,
					})
					return nil
				}

				if err == cliErrors.ErrNotLoggedIn {
					fmt.Printf("\n🔒 Error: You are not logged in. Please run 'kavach login' to authenticate.\n")
					logger.Warn("User not logged in during secret group revoke", map[string]interface{}{
						"cmd":         "group revoke",
						"secretGroup": secretGroupName,
						"org":         orgName,
					})
					return nil
				}

				if err == cliErrors.ErrRoleBindingNotFound {
					fmt.Printf("\n⚠️  Warning: No role binding found for the specified user/group on this secret group.\n")
					fmt.Printf("   The user/group may not have had '%s' permissions to begin with.\n", role)
					return nil
				}

				if err == cliErrors.ErrOrganizationNotFound {
					fmt.Printf("\n🚨 Error: Organization '%s' not found.\n", orgName)
					fmt.Printf("   Please verify the organization name.\n")
					return nil
				}

				if err == cliErrors.ErrSecretGroupNotFound {
					fmt.Printf("\n🚨 Error: Secret group '%s' not found in organization '%s'.\n", secretGroupName, orgName)
					fmt.Printf("   Please verify the secret group name.\n")
					return nil
				}

				if err == cliErrors.ErrAccessDenied {
					fmt.Printf("\n%s\n", err.Error())
					logger.Warn("Access denied during secret group revoke", map[string]interface{}{
						"cmd":         "group revoke",
						"secretGroup": secretGroupName,
						"org":         orgName,
					})
					return nil
				}

				logger.Error("Failed to revoke secret group permissions", err, map[string]interface{}{
					"cmd":         "group revoke",
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

			fmt.Printf("\n✅ Success: Revoked '%s' permissions from '%s' on secret group '%s'\n", role, target, secretGroupName)
			fmt.Printf("   Organization: %s\n", orgName)

			return nil
		},
	}

	// Command flags with improved descriptions
	cmd.Flags().StringVarP(&role, "role", "r", "", "Permission level to revoke (admin, editor, viewer)")
	cmd.Flags().StringVarP(&userName, "user", "u", "", "GitHub username to revoke permissions from")
	cmd.Flags().StringVarP(&userGroupName, "group", "g", "", "User group name to revoke permissions from")
	cmd.Flags().StringVarP(&orgName, "org", "o", "", "Organization name where the secret group exists")

	// Mark required flags
	cmd.MarkFlagRequired("role")
	cmd.MarkFlagRequired("org")

	return cmd
}
