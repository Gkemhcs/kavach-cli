package environment

import (
	"fmt"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	"github.com/Gkemhcs/kavach-cli/internal/environment"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// GetListEnvironmentHeaders returns the headers for the environment list table.
func GetListEnvironmentHeaders() []string {
	return []string{
		"Environment ID",
		"Environment Name",
		"SecretGroup Name",
		"Role",
		"Active",
	}
}

// NewListEnvironmentCommand returns a Cobra command to list the user's environments.
// Handles user feedback, error reporting, and logging.
func NewListEnvironmentCommand(logger *utils.Logger, envClient environment.EnvironmentClient) *cobra.Command {
	var orgName string
	var groupName string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "📋 List your environments",
		Long: `List all environments in the current secret group.

This command displays a table of all environments within the active secret group,
showing your role in each environment and which one is currently active.

The output includes:
- Environment ID: Unique identifier for the environment
- Environment Name: Human-readable name of the environment
- Secret Group Name: The secret group this environment belongs to
- Role: Your role in the environment (owner, admin, member, viewer)
- Active: Indicates which environment is currently set as default (🟢)

Available roles:
- owner: Full administrative access (create, delete, manage members)
- admin: Administrative access (manage secrets and members)
- member: Basic access (view and use secrets)
- viewer: Read-only access (view secrets only)

Examples:
  kavach env list                    # List all environments in current secret group
  kavach env list --help            # Show detailed help

Note: The active environment (marked with 🟢) is used as the default
context for other commands when no environment is explicitly specified.
Use 'kavach env activate <name>' to change the active environment.`,
		Example: `  kavach env list
  # List all environments in the current secret group`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Info("Listing environments", map[string]interface{}{"cmd": "env list"})
			cfg, _ := config.LoadCLIConfig()
			if cfg.Organization == "" && cfg.SecretGroup == "" && orgName == "" && groupName == "" {
				fmt.Println("\n⚠️  You haven't passed an organization and secret group and no default organization  and secret group are set.")
				fmt.Println("💡 Set a default organization using: `kavach org activate <org-name>` and `kavach group activate <secret-group-name>`")
				logger.Warn("No organization or secret group set for environment list", map[string]interface{}{"cmd": "env list"})
				return nil
			}
			if orgName == "" {
				orgName = cfg.Organization
			}
			if groupName == "" {
				groupName = cfg.SecretGroup
			}
			data, err := envClient.ListEnvironment(orgName, groupName)
			if err != nil {
				if err == cliErrors.ErrUnReachableBackend {
					fmt.Println("🚨 Oops! Kavach backend is currently down or not responding.")
					fmt.Println("📡 This may be due to server downtime or high request volume.")
					fmt.Println("📩 If this persists, please drop us a message at 👉 **gudikotieswarmani@gmail.com**")
					logger.Error("Backend unreachable during environment list", err, map[string]interface{}{"cmd": "env list"})
					return nil // or exit gracefully
				}
				if err == cliErrors.ErrNotLoggedIn {
					fmt.Printf("\n🔒 You are not logged in. Please run 'kavach login' to log in.\n")
					logger.Warn("User not logged in during environment list", map[string]interface{}{"cmd": "env list"})
					return nil
				}
				if err == cliErrors.ErrAccessDenied {
					fmt.Printf("\n%s\n", err.Error())
					logger.Warn("Access denied during environment list", map[string]interface{}{"cmd": "env list"})
					return nil
				}
				logger.Error("Failed to list environments", err, map[string]interface{}{"cmd": "env list"})
				return err
			}
			utils.RenderTable(GetListEnvironmentHeaders(), ToRenderable(data))
			logger.Info("Displayed environment list table", map[string]interface{}{"cmd": "env list", "count": len(data)})
			return nil
		},
	}
	cmd.Flags().StringVarP(&orgName, "organization", "o", "", "Organization under which to list the environments")
	cmd.Flags().StringVarP(&groupName, "group", "g", "", "Secret group under which to list the environments")
	return cmd
}

// ToRenderable converts a list of environments to a 2D string slice for table rendering.
func ToRenderable(data []environment.ListEnvironmentsWithMemberRow) [][]string {
	var out [][]string
	config, _ := config.LoadCLIConfig()
	for _, env := range data {
		active := ""
		if config.Environment == env.Name {
			active = "🟢"
		}
		out = append(out, []string{env.EnvironmentID, env.Name, env.SecretGroupName, env.Role, active})
	}
	return out
}
