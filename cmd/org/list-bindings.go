package org

import (
	"fmt"

	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/org"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewListBindingsCommand returns a Cobra command to list organization role bindings.
func NewListBindingsCommand(logger *utils.Logger, orgClient org.OrgClient) *cobra.Command {
	return &cobra.Command{
		Use:   "list-bindings",
		Short: "List all role bindings for an organization",
		Long: `List all role bindings for an organization with resolved user and group names.

This command shows all role bindings (both direct and inherited) for the specified organization,
including user and group names instead of just IDs.

Examples:
  kavach org list-bindings org-1
  kavach org list-bindings my-organization`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			orgName := args[0]

			logger.Info("Listing organization role bindings", map[string]interface{}{
				"cmd": "org list-bindings",
				"org": orgName,
			})

			// List role bindings
			bindings, err := orgClient.ListRoleBindings(orgName)
			if err != nil {
				logger.Error("Failed to list organization role bindings", err, map[string]interface{}{
					"cmd": "org list-bindings",
					"org": orgName,
				})

				// Handle specific error types
				switch err {
				case cliErrors.ErrNoRoleBindingsFound:
					fmt.Printf("📋 No role bindings found for organization '%s'\n", orgName)
					return nil
				case cliErrors.ErrOrganizationNotFound:
					fmt.Printf("❌ Organization '%s' not found\n", orgName)
					return nil
				case cliErrors.ErrPermissionDeniedForRoleBindings:
					fmt.Printf("🚫 You don't have permission to view role bindings for organization '%s'\n", orgName)
					return nil
				case cliErrors.ErrInvalidResourceID:
					fmt.Printf("❌ Invalid organization ID provided\n")
					return nil
				case cliErrors.ErrRoleBindingsListFailed:
					fmt.Printf("❌ Failed to list role bindings for organization '%s'. Please try again\n", orgName)
					return nil
				default:
					// For any other errors, return the error as is
					return err
				}
			}

			// Display results (we know bindings exist since we didn't get ErrNoRoleBindingsFound)

			fmt.Printf("Role bindings for organization '%s':\n", orgName)
			fmt.Printf("Total bindings: %d\n\n", len(bindings))

			// Group bindings by type for better display
			directBindings := make([]org.RoleBinding, 0)
			inheritedFromOrg := make([]org.RoleBinding, 0)

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

			logger.Info("Successfully displayed organization role bindings", map[string]interface{}{
				"cmd":   "org list-bindings",
				"org":   orgName,
				"count": len(bindings),
			})

			return nil
		},
	}
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
func toRenderableBindings(bindings []org.RoleBinding) [][]string {
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
