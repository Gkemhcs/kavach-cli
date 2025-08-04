package org

import (
	"fmt"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/org"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// GetListOrganizationHeaders returns the headers for the organization list table.
func GetListOrganizationHeaders() []string {
	return []string{
		"Org Id",
		"Org Name",
		"Role",
		"Active",
	}
}

// NewListOrgCommand returns a Cobra command to list the user's organizations.
// Handles user feedback, error reporting, and logging.
func NewListOrgCommand(logger *utils.Logger, orgClient org.OrgClient) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "📋 List your organizations",
		Long: `List all organizations you have access to in Kavach.

This command displays a table of all organizations where you are a member,
showing your role in each organization and which one is currently active.

The output includes:
- Organization ID: Unique identifier for the organization
- Organization Name: Human-readable name of the organization
- Role: Your role in the organization (owner, admin, member, viewer)
- Active: Indicates which organization is currently set as default (🟢)

Available roles:
- owner: Full administrative access (create, delete, manage members)
- admin: Administrative access (manage resources and members)
- member: Basic access (view and use resources)
- viewer: Read-only access (view resources only)

Examples:
  kavach org list                    # List all accessible organizations
  kavach org list --help            # Show detailed help

Note: The active organization (marked with 🟢) is used as the default
context for other commands when no organization is explicitly specified.
Use 'kavach org activate <name>' to change the active organization.`,
		Example: `  kavach org list
  # List all organizations you have access to`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Info("Listing organizations", map[string]interface{}{"cmd": "org list"})
			data, err := orgClient.ListMyOrganizations()
			if err != nil {
				if err == cliErrors.ErrUnReachableBackend {
					fmt.Println("🚨 Oops! Kavach backend is currently down or not responding.")
					fmt.Println("📡 This may be due to server downtime or high request volume.")
					fmt.Println("📩 If this persists, please drop us a message at 👉 **gudikotieswarmani@gmail.com**")
					logger.Error("Backend unreachable during org list", err, map[string]interface{}{"cmd": "org list"})
					return nil // or exit gracefully
				}
				if err == cliErrors.ErrNotLoggedIn {
					fmt.Printf("\n🔒 You are not logged in. Please run 'kavach login' to log in.\n")
					logger.Warn("User not logged in during org list", map[string]interface{}{"cmd": "org list"})
					return nil
				}
				if err == cliErrors.ErrAccessDenied {
					fmt.Printf("\n%s\n", err.Error())
					logger.Warn("Access denied during org list", map[string]interface{}{"cmd": "org list"})
					return nil
				}
				logger.Error("Failed to list organizations", err, map[string]interface{}{"cmd": "org list"})
				return err
			}
			utils.RenderTable(GetListOrganizationHeaders(), ToRenderable(data))
			logger.Info("Displayed organization list table", map[string]interface{}{"cmd": "org list", "count": len(data)})
			return nil
		},
	}
}

// ToRenderable converts a list of organizations to a 2D string slice for table rendering.
func ToRenderable(data []org.ListMembersOfOrganizationRow) [][]string {
	var out [][]string
	config, _ := config.LoadCLIConfig()
	for _, org := range data {
		active := ""
		if config.Organization == org.Name {
			active = "🟢"
		}
		out = append(out, []string{org.OrgID.String(), org.Name, org.Role, active})
	}
	return out
}
