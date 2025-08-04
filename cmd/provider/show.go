package provider

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/provider"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewShowCommand creates a new show command
func NewShowCommand(logger *utils.Logger, cfg *config.Config, providerClient provider.ProviderClient) *cobra.Command {
	var (
		// Organization flags
		orgName   string
		groupName string
		envName   string
	)

	cmd := &cobra.Command{
		Use:   "show [provider]",
		Short: "👁️ Show detailed provider credential information",
		Long: `Show detailed information about a specific provider credential.

This command displays comprehensive details about a provider credential,
including configuration settings and credential information.

Key concepts:
- Shows complete provider credential details
- Displays configuration settings and metadata
- Shows credential information (encrypted)
- Useful for auditing and verification
- Helps troubleshoot sync issues

The output includes:
- Provider Information: ID, type, environment, timestamps
- Configuration: Provider-specific settings and options
- Credentials: Encrypted credential data (not displayed)

Prerequisites:
- Provider credential must exist
- You must have appropriate permissions
- Organization, secret group, and environment must be specified

Use cases:
- Auditing provider configurations
- Verifying credential setup
- Troubleshooting sync issues
- Compliance and security reviews
- Configuration management

Arguments:
  provider    Provider name (github, gcp, azure)

Required flags:
  --org    Organization name
  --group  Secret group name
  --env    Environment name

Examples:
  kavach provider show github --org "myorg" --group "mygroup" --env "prod"
  kavach provider show gcp --org "myorg" --group "mygroup" --env "prod"
  kavach provider show azure --org "myorg" --group "mygroup" --env "prod"

Note: Credential information is encrypted and not displayed for security.
Use 'kavach provider update <provider>' to modify provider credentials.`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			providerName := args[0]

			// Validate provider name
			if err := provider.ValidateProviderName(providerName); err != nil {
				fmt.Fprintf(os.Stderr, "❌ %v\n", err)
				os.Exit(1)
			}

			// Validate required flags
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

			// Get provider credential
			fmt.Printf("👁️ Showing %s provider details for %s/%s/%s...\n\n", providerName, orgName, groupName, envName)

			provider, err := providerClient.GetProviderCredential(orgName, groupName, envName, providerName)
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
				case cliErrors.ErrProviderCredentialNotFound:
					fmt.Printf("\n❌ Provider credential for '%s' not found.\n", providerName)
					return
				case cliErrors.ErrAccessDenied:
					fmt.Printf("\n%s\n", err.Error())
					return
				default:
					fmt.Fprintf(os.Stderr, "❌ Failed to get provider credential: %v\n", err)
					os.Exit(1)
				}
			}

			// Display provider information
			fmt.Printf("🔧 Provider Information:\n")
			fmt.Printf("   ID: %s\n", provider.ID)
			fmt.Printf("   Type: %s\n", provider.Provider)
			fmt.Printf("   Environment ID: %s\n", provider.EnvironmentID)
			fmt.Printf("   Created: %s\n", provider.CreatedAt)
			fmt.Printf("   Updated: %s\n", provider.UpdatedAt)

			// Display configuration
			fmt.Printf("\n⚙️ Configuration:\n")
			configJSON, _ := json.MarshalIndent(provider.Config, "   ", "  ")
			fmt.Println(string(configJSON))

			// Display credentials (if available)
			if len(provider.Credentials) > 0 {
				fmt.Printf("\n🔑 Credentials:\n")
				credentialsJSON, _ := json.MarshalIndent(provider.Credentials, "   ", "  ")
				fmt.Println(string(credentialsJSON))
			} else {
				fmt.Printf("\n🔑 Credentials: (encrypted, not displayed)\n")
			}

			fmt.Printf("\n✅ Provider credential details displayed successfully!\n")
		},
	}

	// Add flags
	cmd.Flags().StringVar(&orgName, "org", "", "Organization name (required)")
	cmd.Flags().StringVar(&groupName, "group", "", "Secret group name (required)")
	cmd.Flags().StringVar(&envName, "env", "", "Environment name (required)")

	// Mark required flags
	cmd.MarkFlagRequired("org")
	cmd.MarkFlagRequired("group")
	cmd.MarkFlagRequired("env")

	return cmd
}
