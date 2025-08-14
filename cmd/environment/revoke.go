package environment

import (
	"fmt"

	"github.com/Gkemhcs/kavach-cli/internal/environment"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/types"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewRevokeEnvironmentCommand creates a Cobra command for revoking environment access permissions.
// This command allows users to revoke specific roles (admin, editor, viewer) from users or user groups
// for accessing environments within secret groups.
func NewRevokeEnvironmentCommand(logger *utils.Logger, envClient environment.EnvironmentClient) *cobra.Command {
	var role string
	var userName string
	var userGroupName string
	var orgName string
	var secretGroupName string

	cmd := &cobra.Command{
		Use:   "revoke <environment-name>",
		Short: "🚫 Revoke access permissions from an environment",
		Long: `Revoke access permissions from an environment for a user or user group.

This command allows you to remove specific roles from users or user groups for accessing
environments. When you revoke permissions, the user/group will lose access to the
environment and all its secrets and configurations.

Key concepts:
- Permissions are revoked at the environment level
- Users/groups lose access to all secrets within the environment
- You can revoke specific roles while keeping other roles intact
- Only environment owners and admins can revoke permissions

Available roles to revoke:
- admin: Administrative access (manage secrets, members, and grant permissions)
- editor: Basic access (view and use secrets, create secrets)
- viewer: Read-only access (view secrets only, cannot modify anything)

Important considerations:
- Revoking permissions is immediate and affects all resources in the environment
- Users may lose access to secrets they were working with
- Consider the impact on ongoing work before revoking permissions
- You cannot revoke your own permissions (as an environment owner)

Prerequisites:
- You must be an admin or owner of the environment
- The user or group must have the specified role binding
- You must specify either --user or --group (not both)
- Organization and secret group must be specified

Examples:
  kavach env revoke production --user "john.doe" --role editor --org "my-org" --secret-group "backend"
  kavach env revoke staging --group "developers" --role viewer --org "my-org" --secret-group "frontend"
  kavach env revoke production --user "sarah" --role admin --org "my-org" --secret-group "backend"

Note: If a user/group doesn't have the specified role binding, the command will
show a warning but won't fail. Use 'kavach env grant' to add new permissions.`,
		Example: `  kavach env revoke production --user "john.doe" --role editor --org "my-org" --secret-group "backend"
  kavach env revoke staging --group "developers" --role viewer --org "my-org" --secret-group "frontend"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			environmentName := args[0]

			logger.Info("Revoking environment access permissions", map[string]interface{}{
				"cmd":         "env revoke",
				"environment": environmentName,
				"org":         orgName,
				"secretGroup": secretGroupName,
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
			if orgName == "" {
				fmt.Printf("\n🚨 Error: Organization name is required. Please use --org to specify the organization.\n")
				return nil
			}

			if secretGroupName == "" {
				fmt.Printf("\n🚨 Error: Secret group name is required. Please use --secret-group to specify the secret group.\n")
				return nil
			}

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
				UserName:        userName,
				GroupName:       userGroupName,
				Role:            role,
				OrgName:         orgName,
				SecretGroupName: secretGroupName,
				EnvironmentName: environmentName,
			}

			// Execute the role binding revocation
			err := envClient.RevokeRoleBinding(req)
			if err != nil {
				if err == cliErrors.ErrUnReachableBackend {
					fmt.Println("🚨 Error: Unable to connect to Kavach backend")
					fmt.Println("📡 This may be due to server downtime or network connectivity issues.")
					fmt.Println("📩 If this persists, please contact support at 👉 **gudikotieswarmani@gmail.com**")
					logger.Error("Backend unreachable during environment revoke", err, map[string]interface{}{
						"cmd":         "env revoke",
						"environment": environmentName,
						"org":         orgName,
					})
					return nil
				}

				if err == cliErrors.ErrNotLoggedIn {
					fmt.Printf("\n🔒 Error: You are not logged in. Please run 'kavach login' to authenticate.\n")
					logger.Warn("User not logged in during environment revoke", map[string]interface{}{
						"cmd":         "env revoke",
						"environment": environmentName,
						"org":         orgName,
					})
					return nil
				}

				if err == cliErrors.ErrRoleBindingNotFound {
					fmt.Printf("\n⚠️  Warning: No role binding found for the specified user/group on this environment.\n")
					fmt.Printf("   The user/group may not have had '%s' permissions to begin with.\n", role)
					return nil
				}

				if err == cliErrors.ErrOrganizationNotFound {
					fmt.Printf("\n🚨 Error: Organization '%s' not found.\n", orgName)
					fmt.Printf("   Please verify the organization name.\n")
					return nil
				}

				if err == cliErrors.ErrEnvironmentNotFound {
					fmt.Printf("\n🚨 Error: Environment '%s' not found in secret group '%s'.\n", environmentName, secretGroupName)
					fmt.Printf("   Please verify the environment name and secret group.\n")
					return nil
				}

				if err == cliErrors.ErrSecretGroupNotFound {
					fmt.Printf("\n🚨 Error: Secret group '%s' not found in organization '%s'.\n", secretGroupName, orgName)
					fmt.Printf("   Please verify the secret group name.\n")
					return nil
				}

				if err == cliErrors.ErrAccessDenied {
					fmt.Printf("\n%s\n", err.Error())
					logger.Warn("Access denied during environment revoke", map[string]interface{}{
						"cmd":         "env revoke",
						"environment": environmentName,
						"org":         orgName,
					})
					return nil
				}
				// Check if the error message contains authentication-related text
				if cliErrors.IsAuthenticationError(err) {
					fmt.Printf("\n🔑 Please login again, the session is expired, unable to authenticate you\n")
					logger.Warn("Authentication error during environment revoke", map[string]interface{}{
						"cmd":         "env revoke",
						"environment": environmentName,
						"org":         orgName,
					})
					return nil
				}
				logger.Error("Failed to revoke environment permissions", err, map[string]interface{}{
					"cmd":         "env revoke",
					"environment": environmentName,
					"org":         orgName,
				})
				return err
			}

			// Success message
			target := userName
			if userName == "" {
				target = userGroupName
			}

			fmt.Printf("\n✅ Success: Revoked '%s' permissions from '%s' on environment '%s'\n", role, target, environmentName)
			fmt.Printf("   Organization: %s\n", orgName)
			fmt.Printf("   Secret Group: %s\n", secretGroupName)

			return nil
		},
	}

	// Command flags with improved descriptions
	cmd.Flags().StringVarP(&role, "role", "r", "", "Permission level to revoke (admin, editor, viewer)")
	cmd.Flags().StringVarP(&userName, "user", "u", "", "GitHub username to revoke permissions from")
	cmd.Flags().StringVarP(&userGroupName, "group", "g", "", "User group name to revoke permissions from")
	cmd.Flags().StringVarP(&orgName, "org", "o", "", "Organization name where the environment exists")
	cmd.Flags().StringVarP(&secretGroupName, "secret-group", "s", "", "Secret group name containing the environment")

	// Mark required flags
	cmd.MarkFlagRequired("role")
	cmd.MarkFlagRequired("org")
	cmd.MarkFlagRequired("secret-group")

	return cmd
}
