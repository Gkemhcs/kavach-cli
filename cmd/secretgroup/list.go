package secretgroup

import (
	"fmt"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"

	"github.com/Gkemhcs/kavach-cli/internal/secretgroup"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

func GetListSecretGroupHeaders() []string {
	return []string{
		"Secret Group Id",
		"Secret Group  Name",
		"Organization Name",
		"Role",
		"Active",
	}
}

func NewListSecretGroupCommand(logger *utils.Logger, groupclient secretgroup.SecretGroupClient) *cobra.Command {
	var orgName string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "📋 List your secret groups",
		Long: `List all secret groups in the current organization.

This command displays a table of all secret groups within the active organization,
showing your role in each secret group and which one is currently active.

The output includes:
- Secret Group ID: Unique identifier for the secret group
- Secret Group Name: Human-readable name of the secret group
- Organization Name: The organization this secret group belongs to
- Role: Your role in the secret group (owner, admin, member, viewer)
- Active: Indicates which secret group is currently set as default (🟢)

Available roles:
- owner: Full administrative access (create, delete, manage members)
- admin: Administrative access (manage resources and members)
- member: Basic access (view and use resources)
- viewer: Read-only access (view resources only)

Examples:
  kavach group list                    # List all secret groups in current org
  kavach group list --help            # Show detailed help

Note: The active secret group (marked with 🟢) is used as the default
context for other commands when no secret group is explicitly specified.
Use 'kavach group activate <name>' to change the active secret group.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			config, _ := config.LoadCLIConfig()
			if orgName == "" {
				if config.Organization == "" {
					fmt.Println("\n⚠️  You haven't passed an organization and no default organization is set.")
					fmt.Println("💡 Set a default organization using: `kavach org activate <org-name>`")
					return nil
				}
				orgName = config.Organization

			}

			data, err := groupclient.ListSecretGroups(orgName)

			if err != nil {
				if err == cliErrors.ErrUnReachableBackend {
					fmt.Println("🚨 Oops! Kavach backend is currently down or not responding.")
					fmt.Println("📡 This may be due to server downtime or high request volume.")
					fmt.Println("📩 If this persists, please drop us a message at 👉 **gudikotieswarmani@gmail.com**")
					return nil // or exit gracefully
				}
				if err == cliErrors.ErrNotLoggedIn {
					fmt.Printf("\n🔒 You are not logged in. Please run 'kavach login' to log in.\n")
					return nil
				}
				if err == cliErrors.ErrAccessDenied {
					fmt.Printf("\n%s\n", err.Error())
					return nil
				}
				return err
			}

			utils.RenderTable(GetListSecretGroupHeaders(), ToRenderable(data))
			return nil
		},
	}
	cmd.Flags().StringVarP(&orgName, "organization", "o", "", "Organisation under which we need to list the secret groups")
	return cmd
}

func ToRenderable(data []secretgroup.ListSecretGroupsWithMemberRow) [][]string {
	var out [][]string
	config, _ := config.LoadCLIConfig()

	for _, group := range data {
		if config.SecretGroup == group.Name {
			out = append(out, []string{group.SecretGroupID.String(), group.Name, group.OrganizationName, group.Role, "🟢"})
		} else {
			out = append(out, []string{group.SecretGroupID.String(), group.Name, group.OrganizationName, group.Role, ""})
		}

	}
	return out
}
