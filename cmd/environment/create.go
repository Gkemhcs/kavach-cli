package environment

import (
	"fmt"
	"strings"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	"github.com/Gkemhcs/kavach-cli/internal/environment"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewCreateEnvironmentCommand returns a Cobra command to create a new environment.
// Handles user feedback, error reporting, and logging.
func NewCreateEnvironmentCommand(logger *utils.Logger, envClient environment.EnvironmentClient) *cobra.Command {
	var description string
	var orgName string
	var groupName string
	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "🏗️ Create a new environment",
		Long: `Create a new environment within the current secret group.

Environments are containers where your actual secrets are stored and managed.
Each environment can have different provider configurations and access controls
for different deployment stages.

Key features:
- You become the owner of the created environment
- Environment names must be unique within the secret group
- Environments can have different provider configurations (GCP, Azure, GitHub)
- You can invite other users and assign different roles
- Environments help separate secrets by deployment stage or configuration

Available roles for environment members:
- owner: Full administrative access (create, delete, manage members)
- admin: Administrative access (manage secrets and members)
- member: Basic access (view and use secrets)
- viewer: Read-only access (view secrets only)

Common environment patterns:
- development: For development and testing
- staging: For pre-production testing
- production: For live production systems
- testing: For automated testing environments
- qa: For quality assurance testing

Use cases:
- Separate secrets by deployment stage (dev, staging, prod)
- Different configurations for different regions
- Isolate secrets by team or project
- Different access controls per environment

Examples:
  kavach env create production --description "Production environment"
  kavach env create development --description "Development environment"
  kavach env create production                    # Without description

Note: Environment names should be descriptive and follow your naming conventions.
Once created, you can activate the environment to set it as default for
future commands.`,
		Example: `  kavach env create production --description "Production environment"
  kavach env create development --description "Development environment"
  kavach env create production                    # Without description`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			logger.Info("Creating environment", map[string]interface{}{"cmd": "env create", "env": name, "org": orgName, "group": groupName})
			cfg, _ := config.LoadCLIConfig()
			if cfg.Organization == "" && cfg.SecretGroup == "" && orgName == "" && groupName == "" {
				fmt.Println("\n⚠️  You haven't passed an organization and secret group and no default organization  and secret group are set.")
				fmt.Println("💡 Set a default organization using: `kavach org activate <org-name>` and `kavach group activate <secret-group-name>`")
				logger.Warn("No organization or secret group set for environment creation", map[string]interface{}{"cmd": "env create"})
				return nil
			}
			if orgName == "" {
				orgName = cfg.Organization
			}
			if groupName == "" {
				groupName = cfg.SecretGroup
			}
			msg := fmt.Sprintf("Creating environment '%s' under secret group '%s' in organization '%s'", name, groupName, orgName)
			cont := utils.ConfirmSecretGroupCreation(msg)
			if !cont {
				fmt.Print("\n❌ Exiting environment creation.\n")
				logger.Info("User cancelled environment creation", map[string]interface{}{"cmd": "env create", "env": name, "org": orgName, "group": groupName})
				return nil
			}
			err := envClient.CreateEnvironment(name, groupName, orgName, description)
			if err != nil {
				if err == cliErrors.ErrUnReachableBackend {
					fmt.Println("🚨 Oops! Kavach backend is currently down or not responding.")
					fmt.Println("📡 This may be due to server downtime or high request volume.")
					fmt.Println("📩 If this persists, please drop us a message at 👉 **gudikotieswarmani@gmail.com**")
					logger.Error("Backend unreachable during environment create", err, map[string]interface{}{"cmd": "env create", "env": name, "org": orgName, "group": groupName})
					return nil // or exit gracefully
				}
				if err == cliErrors.ErrOrganizationNotFound {
					fmt.Printf("\n❌ Organization '%s' is not found.\n", orgName)
					logger.Warn("Organization not found during environment create", map[string]interface{}{"cmd": "env create", "org": orgName})
					return nil
				}
				if err == cliErrors.ErrSecretGroupNotFound {
					fmt.Printf("\n❌ Secret group '%s' not found under organization '%s'.\n", groupName, orgName)
					logger.Warn("Secret group not found during environment create", map[string]interface{}{"cmd": "env create", "group": groupName, "org": orgName})
					return nil
				}
				if err == cliErrors.ErrNotLoggedIn {
					fmt.Printf("\n🔒 You are not logged in. Please run 'kavach login' to log in.\n")
					logger.Warn("User not logged in during environment create", map[string]interface{}{"cmd": "env create"})
					return nil
				}
				if err == cliErrors.ErrDuplicateEnvironment {
					fmt.Printf("\n❌ Environment '%s' already exists under secret group '%s' in organization '%s'. Please choose a different name.\n", name, groupName, orgName)
					logger.Warn("Duplicate environment during create", map[string]interface{}{"cmd": "env create", "env": name, "group": groupName, "org": orgName})
					return nil
				}
				if strings.Contains(err.Error(), "not allowed") {
					fmt.Printf("\nAllowed environment names are: dev, prod, staging.\n")
					logger.Warn("Environment name not allowed during create", map[string]interface{}{"cmd": "env create", "env": name})
					return nil
				}
				if err == cliErrors.ErrAccessDenied {
					fmt.Printf("\n%s\n", err.Error())
					logger.Warn("Access denied during environment create", map[string]interface{}{"cmd": "env create", "env": name, "group": groupName, "org": orgName})
					return nil
				}
				logger.Error("Failed to create environment", err, map[string]interface{}{"cmd": "env create", "env": name, "group": groupName, "org": orgName})
				return err
			}
			fmt.Printf("\n🎉 Environment '%s' created successfully!\n", name)
			logger.Info("Environment created successfully", map[string]interface{}{"cmd": "env create", "env": name, "group": groupName, "org": orgName})
			return nil
		},
	}
	cmd.Flags().StringVar(&description, "description", "", "Description of the environment")
	cmd.Flags().StringVarP(&orgName, "organization", "o", "", "Organization under which you want to create the environment")
	cmd.Flags().StringVarP(&groupName, "secret-group", "g", "", "Secret Group under which to create the environment")
	return cmd
}
