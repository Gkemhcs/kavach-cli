package environment

import (
	"fmt"

	"github.com/Gkemhcs/kavach-cli/internal/environment"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/types"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// validateRole checks if the provided role is valid for environment permissions.
// Valid roles are: admin, viewer, editor
func validateRole(role string) bool {
	return role == "admin" || role == "viewer" || role == "editor"
}

// NewGrantEnvironmentCommand creates a Cobra command for granting environment access permissions.
// This command allows users to grant specific roles (admin, editor, viewer) to users or user groups
// for accessing environments within secret groups.
func NewGrantEnvironmentCommand(logger *utils.Logger, envClient environment.EnvironmentClient) *cobra.Command {
	var role string
	var userName string
	var userGroupName string
	var orgName string
	var secretGroupName string

	cmd := &cobra.Command{
		Use:   "grant <environment-name>",
		Short: "🔑 Grant access permissions to an environment",
		Long: `Grant access permissions to an environment for a user or user group.

This command allows you to assign specific roles to users or user groups for accessing
environments. The granted permissions will apply to the entire environment and all
its secrets and configurations.

Key concepts:
- Permissions are granted at the environment level
- Users/groups get access to all secrets within the environment
- Role bindings can be updated by granting the same role again
- Only environment owners and admins can grant permissions

Available roles:
- admin: Administrative access (manage secrets, members, and grant permissions)
- editor: Basic access (view and use secrets, create secrets)
- viewer: Read-only access (view secrets only, cannot modify anything)

Permission hierarchy:
- admin > editor > viewer
- Higher roles inherit all permissions of lower roles

Prerequisites:
- You must be an admin or owner of the environment
- The user or group must exist in Kavach
- You must specify either --user or --group (not both)
- Organization and secret group must be specified

Examples:
  kavach env grant production --user "john.doe" --role editor --org "my-org" --secret-group "backend"
  kavach env grant staging --group "developers" --role viewer --org "my-org" --secret-group "frontend"
  kavach env grant production --user "sarah" --role admin --org "my-org" --secret-group "backend"

Note: If a user/group already has a role binding, granting a new role will
update their existing permissions. Use 'kavach env revoke' to remove permissions.`,
		Example: `  kavach env grant production --user "john.doe" --role editor --org "my-org" --secret-group "backend"
  kavach env grant staging --group "developers" --role viewer --org "my-org" --secret-group "frontend"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			environmentName := args[0]

			logger.Info("Granting environment access permissions", map[string]interface{}{
				"cmd":         "env grant",
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
				fmt.Printf("\n🚨 Error: Please specify either a user (--user) or a group (--group) to grant permissions to.\n")
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
				fmt.Printf("\n🚨 Error: Role is required. Please use --role to specify the permission level (admin, editor, or viewer).\n")
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
				EnvironmentName: environmentName,
			}

			// Execute the role binding grant
			err := envClient.GrantRoleBinding(req)
			if err != nil {
				if err == cliErrors.ErrUnReachableBackend {
					fmt.Println("🚨 Error: Unable to connect to Kavach backend")
					fmt.Println("📡 This may be due to server downtime or network connectivity issues.")
					fmt.Println("📩 If this persists, please contact support at 👉 **gudikotieswarmani@gmail.com**")
					logger.Error("Backend unreachable during environment grant", err, map[string]interface{}{
						"cmd":         "env grant",
						"environment": environmentName,
						"org":         orgName,
					})
					return nil
				}

				if err == cliErrors.ErrNotLoggedIn {
					fmt.Printf("\n🔒 Error: You are not logged in. Please run 'kavach login' to authenticate.\n")
					logger.Warn("User not logged in during environment grant", map[string]interface{}{
						"cmd":         "env grant",
						"environment": environmentName,
						"org":         orgName,
					})
					return nil
				}

				if err == cliErrors.ErrDuplicateRoleBinding {
					fmt.Printf("\n⚠️  Warning: Role binding already exists for this user/group on this environment.\n")
					fmt.Printf("   The existing permissions have been updated to '%s'.\n", role)
					return nil
				}

				if err == cliErrors.ErrEnvironmentNotFound {
					fmt.Printf("\n🚨 Error: Environment '%s' not found in secret group '%s'.\n", environmentName, secretGroupName)
					fmt.Printf("   Please verify the environment name and secret group.\n")
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
					logger.Warn("Access denied during environment grant", map[string]interface{}{
						"cmd":         "env grant",
						"environment": environmentName,
						"org":         orgName,
					})
					return nil
				}
				// Check if the error message contains authentication-related text
				if cliErrors.IsAuthenticationError(err) {
					fmt.Printf("\n🔑 Please login again, the session is expired, unable to authenticate you\n")
					logger.Warn("Authentication error during environment grant", map[string]interface{}{
						"cmd":         "env grant",
						"environment": environmentName,
						"org":         orgName,
					})
					return nil
				}
				logger.Error("Failed to grant environment permissions", err, map[string]interface{}{
					"cmd":         "env grant",
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

			fmt.Printf("\n✅ Success: Granted '%s' permissions to '%s' on environment '%s'\n", role, target, environmentName)
			fmt.Printf("   Organization: %s\n", orgName)
			fmt.Printf("   Secret Group: %s\n", secretGroupName)

			return nil
		},
	}

	// Command flags with improved descriptions
	cmd.Flags().StringVarP(&role, "role", "r", "", "Permission level to grant (admin, editor, viewer)")
	cmd.Flags().StringVarP(&userName, "user", "u", "", "GitHub username to grant permissions to")
	cmd.Flags().StringVarP(&userGroupName, "group", "g", "", "User group name to grant permissions to")
	cmd.Flags().StringVarP(&orgName, "org", "o", "", "Organization name where the environment exists")
	cmd.Flags().StringVarP(&secretGroupName, "secret-group", "s", "", "Secret group name containing the environment")

	// Mark required flags
	cmd.MarkFlagRequired("role")
	cmd.MarkFlagRequired("org")
	cmd.MarkFlagRequired("secret-group")

	return cmd
}
