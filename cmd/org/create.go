package org

import (
	"fmt"

	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/org"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewCreateOrgCommand returns a Cobra command to create a new organization.
// Handles user feedback, error reporting, and logging.
func NewCreateOrgCommand(logger *utils.Logger, orgClient org.OrgClient) *cobra.Command {
	var description string
	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "🏗️ Create a new organization",
		Long: `Create a new organization in Kavach.

Organizations are the top-level containers that group related resources together.
When you create an organization, you automatically become its owner with full
administrative privileges.

Key features:
- You become the owner of the created organization
- Organization names must be unique across Kavach
- Organizations can contain multiple secret groups and environments
- You can invite other users and assign different roles

Available roles for organization members:
- owner: Full administrative access (create, delete, manage members)
- admin: Administrative access (manage resources and members)
- member: Basic access (view and use resources)
- viewer: Read-only access (view resources only)

Examples:
  kavach org create mycompany --description "My company organization"
  kavach org create project-alpha --description "Alpha project organization"
  kavach org create mycompany                    # Without description

Note: Organization names should be descriptive and follow your naming conventions.
Once created, you can activate the organization to set it as default for
future commands.`,
		Example: `  kavach org create mycompany --description "My company organization"
  kavach org create project-alpha --description "Alpha project organization"
  kavach org create mycompany                    # Without description`,

		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			logger.Info("Creating organization", map[string]interface{}{"cmd": "org create", "org": name})
			err := orgClient.CreateOrganization(name, description)
			if err != nil {
				if err == cliErrors.ErrUnReachableBackend {
					fmt.Println("🚨 Oops! Kavach backend is currently down or not responding.")
					fmt.Println("📡 This may be due to server downtime or high request volume.")
					fmt.Println("📩 If this persists, please drop us a message at 👉 **gudikotieswarmani@gmail.com**")
					logger.Error("Backend unreachable during org create", err, map[string]interface{}{"cmd": "org create", "org": name})
					return nil // or exit gracefully
				}
				if err == cliErrors.ErrNotLoggedIn {
					fmt.Printf("\n🔒 You are not logged in. Please run 'kavach login' to log in.\n")
					logger.Warn("User not logged in during org create", map[string]interface{}{"cmd": "org create", "org": name})
					return nil
				}
				if err == cliErrors.ErrDuplicateOrganisation {
					fmt.Printf("\n❌ Organization '%s' already exists. Please choose a different name.\n", name)
					logger.Warn("Duplicate organization during create", map[string]interface{}{"cmd": "org create", "org": name})
					return nil
				}
				if err == cliErrors.ErrAccessDenied {
					fmt.Printf("\n%s\n", err.Error())
					logger.Warn("Access denied during org create", map[string]interface{}{"cmd": "org create", "org": name})
					return nil
				}
				if err == cliErrors.ErrInvalidToken {
					fmt.Printf("\n🔒 You are not logged in. Please run 'kavach login' to log in.\n")
					logger.Warn("User not logged in during org create", map[string]interface{}{"cmd": "org create", "org": name})
					return nil
				}

				// Check if the error message contains authentication-related text
				if cliErrors.IsAuthenticationError(err) {
					fmt.Printf("\n🔑 Please login again, the session is expired, unable to authenticate you\n")
					logger.Warn("Authentication error during org create", map[string]interface{}{"cmd": "org create", "org": name})
					return nil
				}

				logger.Error("Failed to create organization", err, map[string]interface{}{"cmd": "org create", "org": name})
				return err
			}
			fmt.Printf("\n🎉 Organization '%s' created successfully!\n", name)
			logger.Info("Organization created successfully", map[string]interface{}{"cmd": "org create", "org": name})
			return nil
		},
	}
	cmd.Flags().StringVar(&description, "description", "", "Description of the organization")
	return cmd
}
