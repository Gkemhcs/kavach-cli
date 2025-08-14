package secretgroup

import (
	"fmt"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/secretgroup"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewCreateSecretGroupCommand creates a new command for creating secret groups
func NewCreateSecretGroupCommand(logger *utils.Logger, groupClient secretgroup.SecretGroupClient) *cobra.Command {

	var description string
	var orgName string
	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "🏗️ Create a new secret group",
		Long: `Create a new secret group within the current organization.

Secret groups are logical containers that organize related secrets and environments.
When you create a secret group, you automatically become its owner with full
administrative privileges.

Key features:
- You become the owner of the created secret group
- Secret group names must be unique within the organization
- Secret groups can contain multiple environments (dev, staging, prod)
- You can invite other users and assign different roles
- Secret groups help organize secrets by project, team, or application

Available roles for secret group members:
- owner: Full administrative access (create, delete, manage members)
- admin: Administrative access (manage resources and members)
- member: Basic access (view and use resources)
- viewer: Read-only access (view resources only)

Use cases:
- Organize secrets by application (e.g., "myapp", "backend", "frontend")
- Group secrets by team (e.g., "dev-team", "qa-team", "ops-team")
- Separate secrets by project (e.g., "project-alpha", "project-beta")

Examples:
  kavach group create myapp --description "My application secrets"
  kavach group create backend --description "Backend service secrets"
  kavach group create myapp                    # Without description

Note: Secret group names should be descriptive and follow your naming conventions.
Once created, you can activate the secret group to set it as default for
future commands.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			config, _ := config.LoadCLIConfig()

			if config.Organization == "" && orgName == "" {
				fmt.Println("\n⚠️  You haven't passed an organization and no default organization is set.")
				fmt.Println("💡 Set a default organization using: `kavach org activate <org-name>`")
				return nil
			}
			if orgName == "" {
				msg := fmt.Sprintf("you havent passed organization option so we are creating the secret  group %s under active  organization %s", name, config.Organization)

				cont := utils.ConfirmSecretGroupCreation(msg)
				if !cont {
					fmt.Print("\n exiting \n")
					return nil
				}

			}

			if orgName == "" {
				orgName = config.Organization
			}

			err := groupClient.CreateSecretGroup(name, description, orgName)
			if err != nil {
				// Only return error if it's not a handled/user-friendly case
				if err == cliErrors.ErrNotLoggedIn {
					fmt.Printf("\n🚫You are not logged in , please run kavach login \n")
					return nil
				}
				if err == cliErrors.ErrUnReachableBackend {
					fmt.Println("🚨 Oops! Kavach backend is currently down or not responding.")
					fmt.Println("📡 This may be due to server downtime or high request volume.")
					fmt.Println("📩 If this persists, please drop us a message at 👉 **gudikotieswarmani@gmail.com**")
					return nil // or exit gracefully
				}
				if err == cliErrors.ErrOrganizationNotFound {
					fmt.Printf("\n❌ Organization '%s' does not exist.\n", orgName)
					return nil
				}
				if err == cliErrors.ErrDuplicateSecretGroup {
					fmt.Printf("\n❌ Secret group '%s' already exists in organization '%s'.\n", name, orgName)
					return nil
				}
				if err == cliErrors.ErrAccessDenied {
					fmt.Printf("\n%s\n", err.Error())
					return nil
				}
				// Check if the error message contains authentication-related text
				if cliErrors.IsAuthenticationError(err) {
					fmt.Printf("\n🔑 Please login again, the session is expired, unable to authenticate you\n")
					return nil
				}
				return err
			}
			fmt.Printf("\n 🎉 Secret group '%s' created successfully under organization '%s'.\n", name, orgName)

			return nil
		},
	}
	cmd.Flags().StringVar(&description, "description", "", "Description of the secret group")
	cmd.Flags().StringVarP(&orgName, "organization", "o", "", "Organisation under which we need to create the secret group")
	return cmd

}
