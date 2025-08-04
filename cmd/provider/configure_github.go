package provider

import (
	"fmt"
	"os"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/provider"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewConfigureGitHubCommand creates a new configure github command
func NewConfigureGitHubCommand(logger *utils.Logger, cfg *config.Config, providerClient provider.ProviderClient) *cobra.Command {
	var (
		// GitHub configuration flags
		githubToken            string
		githubOwner            string
		githubRepository       string
		githubEnvironment      string
		githubSecretVisibility string
		// Organization flags
		orgName   string
		groupName string
		envName   string
	)

	cmd := &cobra.Command{
		Use:   "github",
		Short: "🐙 Configure GitHub provider credentials",
		Long: `Configure GitHub provider credentials for secret synchronization.

This command sets up GitHub Personal Access Token credentials to sync secrets
to GitHub repositories using GitHub Secrets API.

Key concepts:
- Configures GitHub Personal Access Token for secret synchronization
- Enables syncing secrets to GitHub repository secrets
- Supports GitHub environments for different deployment stages
- Configurable secret visibility settings
- Integrates with GitHub Actions and workflows

Prerequisites:
- GitHub Personal Access Token with appropriate permissions
- GitHub repository must exist and be accessible
- Token must have repository secrets permissions
- Organization, secret group, and environment must exist

Required GitHub permissions:
- repo (Full control of private repositories)
- workflow (Update GitHub Action workflows)
- write:packages (Upload packages to GitHub Package Registry)

Required flags:
  --token        GitHub Personal Access Token with repo scope
  --owner        GitHub organization or username
  --repository   GitHub repository name
  --org          Organization name
  --group        Secret group name
  --env          Environment name

Optional flags:
  --environment      GitHub environment name (default: "default")
  --secret-visibility Secret visibility: all, selected, private (default: "private")

Secret visibility options:
- private: Only accessible to the repository
- selected: Accessible to selected repositories
- all: Accessible to all repositories in the organization

Use cases:
- Syncing secrets to GitHub repository secrets
- Integrating with GitHub Actions workflows
- Maintaining secrets for CI/CD pipelines
- Automated secret deployment to GitHub

Example:
  kavach provider configure github \
    --token "ghp_xxxxxxxxxxxxxxxxxxxx" \
    --owner "myorg" \
    --repository "myrepo" \
    --environment "production" \
    --secret-visibility "private" \
    --org "myorg" \
    --group "mygroup" \
    --env "prod"

Note: After configuration, use 'kavach secret sync --provider github' to sync
secrets to GitHub repository secrets.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Validate required flags
			if githubToken == "" {
				fmt.Fprintln(os.Stderr, "❌ Error: --token is required")
				os.Exit(1)
			}
			if githubOwner == "" {
				fmt.Fprintln(os.Stderr, "❌ Error: --owner is required")
				os.Exit(1)
			}
			if githubRepository == "" {
				fmt.Fprintln(os.Stderr, "❌ Error: --repository is required")
				os.Exit(1)
			}
			if orgName == "" {
				fmt.Fprintln(os.Stderr, "❌ Error: --org is required")
				os.Exit(1)
			}
			if groupName == "" {
				fmt.Fprintln(os.Stderr, "❌ Error: --group is required")
				os.Exit(1)
			}
			if envName == "" {
				fmt.Fprintln(os.Stderr, "❌ Error: --env is required")
				os.Exit(1)
			}

			// Set defaults
			if githubEnvironment == "" {
				githubEnvironment = "default"
			}
			if githubSecretVisibility == "" {
				githubSecretVisibility = "private"
			}

			// Validate secret visibility
			validVisibilities := []string{"all", "selected", "private"}
			valid := false
			for _, v := range validVisibilities {
				if v == githubSecretVisibility {
					valid = true
					break
				}
			}
			if !valid {
				fmt.Fprintf(os.Stderr, "❌ Error: --secret-visibility must be one of: %v\n", validVisibilities)
				os.Exit(1)
			}

			// Prepare credentials and config
			credentials := map[string]interface{}{
				"token": githubToken,
			}

			config := map[string]interface{}{
				"owner":             githubOwner,
				"repository":        githubRepository,
				"environment":       githubEnvironment,
				"secret_visibility": githubSecretVisibility,
			}

			// Create provider credential
			fmt.Printf("🔧 Configuring GitHub provider for %s/%s/%s...\n", orgName, groupName, envName)

			result, err := providerClient.CreateProviderCredential(orgName, groupName, envName, "github", credentials, config)
			if err != nil {
				switch err {
				case cliErrors.ErrNotLoggedIn:
					fmt.Printf("\n🔒 You are not logged in. Please run 'kavach login' to log in.\n")
					return
				case cliErrors.ErrUnReachableBackend:
					fmt.Println("🚨 Oops! Kavach backend is currently down or not responding.")
					fmt.Println("📡 This may be due to server downtime or high request volume.")
					fmt.Println("📩 If this persists, please drop us a message at 👉 **gudikotieswarmani@gmail.com**")
					return
				case cliErrors.ErrOrganizationNotFound:
					fmt.Printf("\n❌ Organization '%s' does not exist.\n", orgName)
					return
				case cliErrors.ErrSecretGroupNotFound:
					fmt.Printf("\n❌ Secret group '%s' does not exist.\n", groupName)
					return
				case cliErrors.ErrEnvironmentNotFound:
					fmt.Printf("\n❌ Environment '%s' does not exist.\n", envName)
					return
				case cliErrors.ErrProviderCredentialExists:
					fmt.Printf("\n⚠️ Provider credential for GitHub already exists.\n")
					return
				case cliErrors.ErrProviderCredentialValidationFailed:
					fmt.Printf("\n❌ GitHub provider credential validation failed. Please check your credentials.\n")
					return
				case cliErrors.ErrProviderCredentialCreateFailed:
					fmt.Printf("\n❌ Failed to create GitHub provider credential.\n")
					return
				case cliErrors.ErrGitHubEnvironmentNotFound:
					fmt.Printf("\n❌ GitHub environment '%s' not found in the repository.\n", githubEnvironment)
					return
				case cliErrors.ErrAccessDenied:
					fmt.Printf("\n%s\n", err.Error())
					return
				default:
					fmt.Fprintf(os.Stderr, "❌ Failed to configure GitHub provider: %v\n", err)
					os.Exit(1)
				}
			}

			fmt.Printf("✅ Successfully configured GitHub provider!\n")
			fmt.Printf("   Provider ID: %s\n", result.ID)
			fmt.Printf("   Repository: %s/%s\n", githubOwner, githubRepository)
			fmt.Printf("   Environment: %s\n", githubEnvironment)
			fmt.Printf("   Secret Visibility: %s\n", githubSecretVisibility)
			fmt.Printf("   Created: %s\n", result.CreatedAt)
		},
	}

	// Add flags
	cmd.Flags().StringVar(&githubToken, "token", "", "GitHub Personal Access Token (required)")
	cmd.Flags().StringVar(&githubOwner, "owner", "", "GitHub organization or username (required)")
	cmd.Flags().StringVar(&githubRepository, "repository", "", "GitHub repository name (required)")
	cmd.Flags().StringVar(&githubEnvironment, "environment", "default", "GitHub environment name")
	cmd.Flags().StringVar(&githubSecretVisibility, "secret-visibility", "private", "Secret visibility: all, selected, private")
	cmd.Flags().StringVar(&orgName, "org", "", "Organization name (required)")
	cmd.Flags().StringVar(&groupName, "group", "", "Secret group name (required)")
	cmd.Flags().StringVar(&envName, "env", "", "Environment name (required)")

	// Mark required flags
	cmd.MarkFlagRequired("token")
	cmd.MarkFlagRequired("owner")
	cmd.MarkFlagRequired("repository")
	cmd.MarkFlagRequired("org")
	cmd.MarkFlagRequired("group")
	cmd.MarkFlagRequired("env")

	return cmd
}
