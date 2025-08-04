package provider

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/provider"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewListCommand creates a new list command
func NewListCommand(logger *utils.Logger, cfg *config.Config, providerClient provider.ProviderClient) *cobra.Command {
	var (
		// Organization flags
		orgName   string
		groupName string
		envName   string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "📋 List all configured provider credentials",
		Long: `List all configured provider credentials for the current environment.

This command displays a table of all provider credentials configured
for secret synchronization, including their status and configuration details.

Key concepts:
- Shows all configured provider credentials for the environment
- Displays provider type, ID, and configuration summary
- Helps verify provider setup and configuration
- Useful for auditing and troubleshooting
- Shows creation timestamps for credential management

The output includes:
- Provider: Type of provider (github, gcp, azure)
- ID: Unique identifier for the credential
- Configuration: Summary of provider-specific settings
- Created: When the credential was configured

Prerequisites:
- Organization, secret group, and environment must be specified
- You must have appropriate permissions for the environment

Use cases:
- Auditing configured providers
- Verifying provider setup
- Troubleshooting sync issues
- Credential management and rotation
- Compliance and security reviews

Required flags:
  --org    Organization name
  --group  Secret group name
  --env    Environment name

Example:
  kavach provider list --org "myorg" --group "mygroup" --env "prod"

Note: Use 'kavach provider show <provider>' to see detailed configuration
for a specific provider credential.`,
		Run: func(cmd *cobra.Command, args []string) {
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

			// List provider credentials
			fmt.Printf("📋 Listing provider credentials for %s/%s/%s...\n\n", orgName, groupName, envName)

			providers, err := providerClient.ListProviderCredentials(orgName, groupName, envName)
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
				case cliErrors.ErrAccessDenied:
					fmt.Printf("\n%s\n", err.Error())
					return
				default:
					fmt.Fprintf(os.Stderr, "❌ Failed to list provider credentials: %v\n", err)
					os.Exit(1)
				}
			}

			if len(providers) == 0 {
				fmt.Println("📭 No provider credentials configured yet.")
				fmt.Println("   Use 'kavach provider configure <provider>' to add credentials.")
				return
			}

			// Display in table format
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "PROVIDER\tID\tCONFIGURATION\tCREATED\t")
			fmt.Fprintln(w, "--------\t--\t-------------\t-------\t")

			for _, provider := range providers {
				config := getProviderConfigSummary(provider.Provider, provider.Config)
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n",
					provider.Provider,
					provider.ID,
					config,
					provider.CreatedAt,
				)
			}
			w.Flush()

			fmt.Printf("\n✅ Found %d provider credential(s)\n", len(providers))
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

// getProviderConfigSummary returns a summary of the provider configuration
func getProviderConfigSummary(providerType string, config map[string]interface{}) string {
	switch providerType {
	case "github":
		if owner, ok := config["owner"].(string); ok {
			if repo, ok := config["repository"].(string); ok {
				return fmt.Sprintf("%s/%s", owner, repo)
			}
		}
		return "GitHub Repository"
	case "gcp":
		if projectID, ok := config["project_id"].(string); ok {
			return fmt.Sprintf("Project: %s", projectID)
		}
		return "GCP Project"
	case "azure":
		if keyVault, ok := config["key_vault_name"].(string); ok {
			return fmt.Sprintf("Key Vault: %s", keyVault)
		}
		return "Azure Key Vault"
	default:
		return "Unknown Provider"
	}
}
