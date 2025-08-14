package secretgroup

import (
	"fmt"

	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/secretgroup"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewListBindingsCommand returns a Cobra command to list secret group role bindings.
func NewListBindingsCommand(logger *utils.Logger, secretGroupClient secretgroup.SecretGroupClient) *cobra.Command {
	var orgName string

	cmd := &cobra.Command{
		Use:   "list-bindings",
		Short: "List all role bindings for a secret group",
		Long: `List all role bindings for a secret group with resolved user and group names.

This command shows all role bindings (both direct and inherited) for the specified secret group,
including user and group names instead of just IDs.

Examples:
  kavach group list-bindings group-1 --org org-1
  kavach group list-bindings my-group -o my-org`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupName := args[0]

			logger.Info("Listing secret group role bindings", map[string]interface{}{
				"cmd":   "group list-bindings",
				"org":   orgName,
				"group": groupName,
			})

			// List role bindings
			bindings, err := secretGroupClient.ListRoleBindings(orgName, groupName)
			if err != nil {
				logger.Error("Failed to list secret group role bindings", err, map[string]interface{}{
					"cmd":   "group list-bindings",
					"org":   orgName,
					"group": groupName,
				})

				// Handle specific error types
				switch err {
				case cliErrors.ErrNoRoleBindingsFound:
					fmt.Printf("📋 No role bindings found for secret group '%s' in organization '%s'\n", groupName, orgName)
					return nil
				case cliErrors.ErrOrganizationNotFound:
					fmt.Printf("❌ Organization '%s' not found\n", orgName)
					return nil
				case cliErrors.ErrSecretGroupNotFound:
					fmt.Printf("❌ Secret group '%s' not found in organization '%s'\n", groupName, orgName)
					return nil
				case cliErrors.ErrPermissionDeniedForRoleBindings:
					fmt.Printf("🚫 You don't have permission to view role bindings for secret group '%s' in organization '%s'\n", groupName, orgName)
					return nil
				case cliErrors.ErrInvalidResourceID:
					fmt.Printf("❌ Invalid resource ID provided\n")
					return nil
				case cliErrors.ErrRoleBindingsListFailed:
					fmt.Printf("❌ Failed to list role bindings for secret group '%s' in organization '%s'. Please try again\n", groupName, orgName)
					return nil
				default:
					// For any other errors, return the error as is
					return err
				}
			}

			// Display results (we know bindings exist since we didn't get ErrNoRoleBindingsFound)

			fmt.Printf("Role bindings for secret group '%s' in organization '%s':\n", groupName, orgName)
			fmt.Printf("Total bindings: %d\n\n", len(bindings))

			// Group bindings by type for better display
			directBindings := make([]secretgroup.RoleBinding, 0)
			inheritedFromOrg := make([]secretgroup.RoleBinding, 0)

			for _, binding := range bindings {
				if binding.BindingType == "direct" {
					directBindings = append(directBindings, binding)
				} else if binding.SourceType == "organization" {
					inheritedFromOrg = append(inheritedFromOrg, binding)
				}
			}

			// Display direct bindings
			if len(directBindings) > 0 {
				fmt.Println("Direct Bindings")
				fmt.Println("---------------")
				utils.RenderTable(getRoleBindingsHeaders(), toRenderableBindings(directBindings))
				fmt.Println()
			}

			// Display inherited bindings from organization
			if len(inheritedFromOrg) > 0 {
				fmt.Printf("Inherited from Organization: %s\n", orgName)
				fmt.Println("---------------------------------------")
				utils.RenderTable(getRoleBindingsHeaders(), toRenderableBindings(inheritedFromOrg))
			}

			logger.Info("Successfully displayed secret group role bindings", map[string]interface{}{
				"cmd":   "group list-bindings",
				"org":   orgName,
				"group": groupName,
				"count": len(bindings),
			})

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&orgName, "org", "o", "", "Organization name (required)")
	cmd.MarkFlagRequired("org")

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
func toRenderableBindings(bindings []secretgroup.RoleBinding) [][]string {
	var out [][]string
	for _, binding := range bindings {
		switch binding.EntityType {
		case "user":
			var userName string
			if binding.EntityName != nil && *binding.EntityName != "" {
				userName = *binding.EntityName
			} else {
				userName = "Unknown User"
			}
			out = append(out, []string{"👤 User", userName, binding.Role})
		case "group":
			var groupName string
			if binding.GroupName != nil && *binding.GroupName != "" {
				groupName = *binding.GroupName
			} else {
				groupName = "Unknown Group"
			}
			out = append(out, []string{"👥 Group", groupName, binding.Role})
		}
	}
	return out
}
