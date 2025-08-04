package org

import (
	"fmt"

	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/org"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewDeleteOrgCommand returns a Cobra command to delete an organization by name.
// Handles user confirmation, error reporting, and logging.
func NewDeleteOrgCommand(logger *utils.Logger, orgClient org.OrgClient) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <name>",
		Short: "🗑️ Delete an organization by name",
		Long: `Delete an organization and all its associated resources.

This command permanently deletes an organization and all its contents including:
- All secret groups within the organization
- All environments within those secret groups
- All secrets stored in those environments
- All user groups and member associations
- All role bindings and permissions

⚠️  WARNING: This action is irreversible and will permanently delete all data
associated with the organization. Make sure you have backups if needed.

Prerequisites:
- You must be the owner of the organization
- The organization must not have any child resources (secret groups, environments)
- You'll be prompted for confirmation before deletion

The deletion process:
1. Prompts for confirmation to prevent accidental deletion
2. Verifies you have permission to delete the organization
3. Checks for existing child resources
4. Permanently deletes the organization and all its data
5. Confirms successful deletion

Examples:
  kavach org delete mycompany        # Delete organization (with confirmation)
  kavach org list                    # Verify organization is deleted

Note: If the organization contains secret groups or environments, you'll need
to delete those first before you can delete the organization.`,
		Example: `  kavach org delete mycompany
  # Delete organization (will prompt for confirmation)`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			logger.Info("Attempting to delete organization", map[string]interface{}{"cmd": "org delete", "org": name})
			msg := fmt.Sprintf("Are you sure you want to delete organization '%s'? If yes, click on proceed, otherwise cancel.", name)
			if !utils.ConfirmSecretGroupCreation(msg) {
				fmt.Print("\n❌ Cancelled the delete operation.\n")
				logger.Info("User cancelled organization delete", map[string]interface{}{"cmd": "org delete", "org": name})
				return nil
			}
			err := orgClient.DeleteOrganization(name)
			if err != nil {
				if err == cliErrors.ErrUnReachableBackend {
					fmt.Println("🚨 Oops! Kavach backend is currently down or not responding.")
					fmt.Println("📡 This may be due to server downtime or high request volume.")
					fmt.Println("📩 If this persists, please drop us a message at 👉 **gudikotieswarmani@gmail.com**")
					logger.Error("Backend unreachable during org delete", err, map[string]interface{}{"cmd": "org delete", "org": name})
					return nil // or exit gracefully
				}
				if err == cliErrors.ErrNotLoggedIn {
					fmt.Printf("\n🔒 You are not logged in. Please run 'kavach login' to log in.\n")
					logger.Warn("User not logged in during org delete", map[string]interface{}{"cmd": "org delete", "org": name})
					return nil
				}
				if err == cliErrors.ErrOrganizationNotFound {
					fmt.Printf("\n❌ Organization '%s' does not exist.\n", name)
					logger.Warn("Organization not found during delete", map[string]interface{}{"cmd": "org delete", "org": name})
					return nil
				}
				if err == cliErrors.ErrForeignKeyViolation {
					fmt.Println("🚨 cannot delete secret group  as it contain child resources like secret groups,environments and secrets")
					fmt.Printf("\n 🚨 first delete all child resources to delete organization %s \n", name)

					return nil
				}
				if err == cliErrors.ErrAccessDenied {
					fmt.Printf("\n%s\n", err.Error())
					logger.Warn("Access denied during org delete", map[string]interface{}{"cmd": "org delete", "org": name})
					return nil
				}
				logger.Error("Failed to delete organization", err, map[string]interface{}{"cmd": "org delete", "org": name})
				return err
			}
			fmt.Printf("\n🗑️ Organization '%s' deleted successfully.\n", name)
			logger.Info("Organization deleted successfully", map[string]interface{}{"cmd": "org delete", "org": name})
			return nil
		},
	}
}
