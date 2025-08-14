package environment

import (
	"fmt"

	"github.com/Gkemhcs/kavach-cli/internal/environment"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewListBindingsCommand returns a Cobra command to list environment role bindings.
func NewListBindingsCommand(logger *utils.Logger, envClient environment.EnvironmentClient) *cobra.Command {
	var orgName, groupName string

	cmd := &cobra.Command{
		Use:   "list-bindings",
		Short: "List all role bindings for an environment",
		Long: `List all role bindings for an environment with resolved user and group names.

This command shows all role bindings (both direct and inherited) for the specified environment,
including user and group names instead of just IDs.

Examples:
  kavach env list-bindings prod --org org-1 --group group-1
  kavach env list-bindings my-env -o my-org -g my-group`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			envName := args[0]

			logger.Info("Listing environment role bindings", map[string]interface{}{
				"cmd":   "env list-bindings",
				"org":   orgName,
				"group": groupName,
				"env":   envName,
			})

			// List role bindings
			bindings, err := envClient.ListRoleBindings(orgName, groupName, envName)
			if err != nil {
				logger.Error("Failed to list environment role bindings", err, map[string]interface{}{
					"cmd":   "env list-bindings",
					"org":   orgName,
					"group": groupName,
					"env":   envName,
				})

				// Handle specific error types
				switch err {
				case cliErrors.ErrNoRoleBindingsFound:
					fmt.Printf("📋 No role bindings found for environment '%s' in organization '%s' and secret group '%s'\n", envName, orgName, groupName)
					return nil
				case cliErrors.ErrOrganizationNotFound:
					fmt.Printf("❌ Organization '%s' not found\n", orgName)
					return nil
				case cliErrors.ErrSecretGroupNotFound:
					fmt.Printf("❌ Secret group '%s' not found in organization '%s'\n", groupName, orgName)
					return nil
				case cliErrors.ErrEnvironmentNotFound:
					fmt.Printf("❌ Environment '%s' not found in secret group '%s' of organization '%s'\n", envName, groupName, orgName)
					return nil
				case cliErrors.ErrPermissionDeniedForRoleBindings:
					fmt.Printf("🚫 You don't have permission to view role bindings for environment '%s' in organization '%s' and secret group '%s'\n", envName, orgName, groupName)
					return nil
				case cliErrors.ErrInvalidResourceID:
					fmt.Printf("❌ Invalid resource ID provided\n")
					return nil
				case cliErrors.ErrRoleBindingsListFailed:
					fmt.Printf("❌ Failed to list role bindings for environment '%s' in organization '%s' and secret group '%s'. Please try again\n", envName, orgName, groupName)
					return nil
				default:
					// Check if the error message contains authentication-related text
					if cliErrors.IsAuthenticationError(err) {
						fmt.Printf("\n🔑 Please login again, the session is expired, unable to authenticate you\n")
						logger.Warn("Authentication error during environment list-bindings", map[string]interface{}{
							"cmd":   "env list-bindings",
							"org":   orgName,
							"group": groupName,
							"env":   envName,
						})
						return nil
					}
					// For any other errors, return the error as is
					return err
				}
			}

			// Display results (we know bindings exist since we didn't get ErrNoRoleBindingsFound)

			fmt.Printf("Role bindings for environment '%s' in organization '%s' and secret group '%s':\n", envName, orgName, groupName)
			fmt.Printf("Total bindings: %d\n\n", len(bindings))

			// Group bindings by type for better display
			directBindings := make([]environment.RoleBinding, 0)
			inheritedFromOrg := make([]environment.RoleBinding, 0)
			inheritedFromSecretGroup := make([]environment.RoleBinding, 0)

			for _, binding := range bindings {
				if binding.BindingType == "direct" {
					directBindings = append(directBindings, binding)
				} else if binding.SourceType == "organization" {
					inheritedFromOrg = append(inheritedFromOrg, binding)
				} else if binding.SourceType == "secret_group" {
					inheritedFromSecretGroup = append(inheritedFromSecretGroup, binding)
				}
			}

			// Display direct bindings
			if len(directBindings) > 0 {
				fmt.Println("Direct Bindings")
				fmt.Println("---------------")
				utils.RenderTable(getRoleBindingsHeaders(), toRenderableBindings(directBindings))
				fmt.Println()
			}

			// Display inherited bindings from secret group
			if len(inheritedFromSecretGroup) > 0 {
				fmt.Printf("Inherited from Secret Group: %s\n", groupName)
				fmt.Println("--------------------------------------------")
				utils.RenderTable(getRoleBindingsHeaders(), toRenderableBindings(inheritedFromSecretGroup))
				fmt.Println()
			}

			// Display inherited bindings from organization
			if len(inheritedFromOrg) > 0 {
				fmt.Printf("Inherited from Organization: %s\n", orgName)
				fmt.Println("---------------------------------------")
				utils.RenderTable(getRoleBindingsHeaders(), toRenderableBindings(inheritedFromOrg))
			}

			logger.Info("Successfully displayed environment role bindings", map[string]interface{}{
				"cmd":   "env list-bindings",
				"org":   orgName,
				"group": groupName,
				"env":   envName,
				"count": len(bindings),
			})

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&orgName, "org", "o", "", "Organization name (required)")
	cmd.Flags().StringVarP(&groupName, "group", "g", "", "Secret group name (required)")
	cmd.MarkFlagRequired("org")
	cmd.MarkFlagRequired("group")

	return cmd
}

// getRoleBindingsHeaders returns the headers for the role bindings table
func getRoleBindingsHeaders() []string {
	return []string{
		"Type",
		"Name",
		"Role",
	}
}

// toRenderableBindings converts role bindings to a 2D string slice for table rendering
func toRenderableBindings(bindings []environment.RoleBinding) [][]string {
	var out [][]string
	for _, binding := range bindings {
		switch binding.EntityType {
		case "user":
			var userName string
			if binding.EntityName != "" {
				userName = binding.EntityName
			} else {
				userName = "Unknown User"
			}
			out = append(out, []string{"👤 User", userName, binding.Role})
		case "group":
			var groupName string
			if binding.GroupName != "" {
				groupName = binding.GroupName
			} else {
				groupName = "Unknown Group"
			}
			out = append(out, []string{"👥 Group", groupName, binding.Role})
		}
	}
	return out
}
