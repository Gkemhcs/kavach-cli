package provider

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/provider"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewDeleteCommand creates a new delete command
func NewDeleteCommand(logger *utils.Logger, cfg *config.Config, providerClient provider.ProviderClient) *cobra.Command {
	var (
		// Organization flags
		orgName   string
		groupName string
		envName   string
	)

	cmd := &cobra.Command{
		Use:   "delete [provider]",
		Short: "🗑️ Delete a provider credential",
		Long: `Delete a provider credential for secret synchronization.

This command removes a provider credential configuration. This action
cannot be undone and will stop secret synchronization to that provider.

Key concepts:
- Permanently removes provider credential configuration
- Stops secret synchronization to the specified provider
- Requires confirmation before deletion
- Cannot be undone - credential must be reconfigured
- Affects all future sync operations to that provider

Prerequisites:
- Provider credential must exist
- You must have appropriate permissions
- Organization, secret group, and environment must be specified

The deletion process:
1. Validates the provider credential exists
2. Shows current configuration details
3. Requests user confirmation
4. Permanently removes the credential
5. Stops secret synchronization

Use cases:
- Removing unused provider configurations
- Switching to different provider credentials
- Cleaning up old or invalid configurations
- Security maintenance and credential rotation

Arguments:
  provider    Provider name (github, gcp, azure)

Required flags:
  --org    Organization name
  --group  Secret group name
  --env    Environment name

Examples:
  kavach provider delete github --org "myorg" --group "mygroup" --env "prod"
  kavach provider delete gcp --org "myorg" --group "mygroup" --env "prod"
  kavach provider delete azure --org "myorg" --group "mygroup" --env "prod"

Warning: This action cannot be undone. After deletion, you must reconfigure
the provider using 'kavach provider configure <provider>' to resume syncing.`,
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

			// Show what will be deleted
			fmt.Printf("🗑️ About to delete %s provider credential for %s/%s/%s\n", providerName, orgName, groupName, envName)
			fmt.Printf("⚠️  This action cannot be undone!\n\n")

			// Get provider details for confirmation
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

			// Show provider details
			fmt.Printf("Provider Details:\n")
			fmt.Printf("   ID: %s\n", provider.ID)
			fmt.Printf("   Type: %s\n", provider.Provider)
			fmt.Printf("   Created: %s\n", provider.CreatedAt)

			// Show configuration summary
			config := getProviderConfigSummary(provider.Provider, provider.Config)
			fmt.Printf("   Configuration: %s\n\n", config)

			// Ask for confirmation
			fmt.Printf("Are you sure you want to delete this provider credential? (y/N): ")
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintf(os.Stderr, "❌ Failed to read confirmation: %v\n", err)
				os.Exit(1)
			}

			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				fmt.Println("❌ Deletion cancelled.")
				return
			}

			// Delete provider credential
			fmt.Printf("🗑️ Deleting %s provider credential...\n", providerName)

			err = providerClient.DeleteProviderCredential(orgName, groupName, envName, providerName)
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
					fmt.Fprintf(os.Stderr, "❌ Failed to delete provider credential: %v\n", err)
					os.Exit(1)
				}
			}

			fmt.Printf("✅ Successfully deleted %s provider credential!\n", providerName)
			fmt.Printf("   Secret synchronization to %s has been disabled.\n", providerName)
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
