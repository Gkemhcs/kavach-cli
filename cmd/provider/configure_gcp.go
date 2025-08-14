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

// NewConfigureGCPCommand creates a new configure gcp command
func NewConfigureGCPCommand(logger *utils.Logger, cfg *config.Config, providerClient provider.ProviderClient) *cobra.Command {
	var (
		// GCP configuration flags
		gcpKeyFile               string
		gcpProjectID             string
		gcpSecretManagerLocation string
		gcpPrefix                string
		gcpReplication           string
		gcpDisableOlderVersions  bool
		// Organization flags
		orgName   string
		groupName string
		envName   string
	)

	cmd := &cobra.Command{
		Use:   "gcp",
		Short: "☁️ Configure Google Cloud Platform provider credentials",
		Long: `Configure Google Cloud Platform provider credentials for secret synchronization.

This command sets up GCP Service Account credentials to sync secrets
to Google Cloud Secret Manager.

Key concepts:
- Configures GCP Service Account for secret synchronization
- Enables syncing secrets to Google Cloud Secret Manager
- Supports custom secret naming with prefixes
- Configurable replication and versioning settings
- Integrates with GCP Secret Manager for secure storage

Prerequisites:
- GCP Service Account with Secret Manager permissions
- GCP Project must exist and be accessible
- Service Account must have appropriate Secret Manager roles
- Organization, secret group, and environment must exist

Required GCP permissions:
- Secret Manager Admin (for managing secrets)
- Secret Manager Secret Accessor (for reading secrets)
- Secret Manager Secret Version Manager (for versioning)

Required flags:
  --key-file     Path to GCP Service Account JSON file
  --project-id   GCP Project ID
  --org          Organization name
  --group        Secret group name
  --env          Environment name

Optional flags:
  --secret-manager-location Secret Manager location (default: "us-central1")
  --prefix                  Secret name prefix (e.g., "kavach/")
  --replication             Replication type: automatic, user-managed (default: "automatic")
  --disable-older-versions  Disable older secret versions when creating new ones

Replication options:
- automatic: Google-managed replication across regions
- user-managed: Custom replication configuration

Use cases:
- Syncing secrets to Google Cloud Secret Manager
- Integrating with GCP-based applications
- Maintaining secrets in Google Cloud environment
- Automated secret deployment to GCP

Example:
  kavach provider configure gcp \
    --key-file "service-account.json" \
    --project-id "my-project-123" \
    --secret-manager-location "us-west1" \
    --prefix "kavach/" \
    --replication "automatic" \
    --disable-older-versions \
    --org "myorg" \
    --group "mygroup" \
    --env "prod"

Note: After configuration, use 'kavach secret sync --provider gcp' to sync
secrets to Google Cloud Secret Manager.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Validate required flags
			if gcpKeyFile == "" {
				fmt.Fprintln(os.Stderr, "❌ Error: --key-file is required")
				os.Exit(1)
			}
			if gcpProjectID == "" {
				fmt.Fprintln(os.Stderr, "❌ Error: --project-id is required")
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
			if gcpSecretManagerLocation == "" {
				gcpSecretManagerLocation = ""
			}
			if gcpReplication == "" {
				gcpReplication = "automatic"
			}

			// Validate replication type
			validReplications := []string{"automatic", "user-managed"}
			valid := false
			for _, r := range validReplications {
				if r == gcpReplication {
					valid = true
					break
				}
			}
			if !valid {
				fmt.Fprintf(os.Stderr, "❌ Error: --replication must be one of: %v\n", validReplications)
				os.Exit(1)
			}

			// Read and validate service account file
			fmt.Printf("📁 Reading GCP service account file: %s\n", gcpKeyFile)
			credentials, err := provider.ReadGCPServiceAccountFile(gcpKeyFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "❌ Failed to read service account file: %v\n", err)
				os.Exit(1)
			}

			// Prepare config
			config := map[string]interface{}{
				"project_id":              gcpProjectID,
				"secret_manager_location": gcpSecretManagerLocation,
			}

			if gcpPrefix != "" {
				config["prefix"] = gcpPrefix
			}
			if gcpReplication != "" {
				config["replication"] = gcpReplication
			}
			if gcpDisableOlderVersions {
				config["disable_older_versions"] = gcpDisableOlderVersions
			}

			// Create provider credential
			fmt.Printf("🔧 Configuring GCP provider for %s/%s/%s...\n", orgName, groupName, envName)

			result, err := providerClient.CreateProviderCredential(orgName, groupName, envName, "gcp", credentials, config)
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
					fmt.Printf("\n⚠️ Provider credential for GCP already exists.\n")
					return
				case cliErrors.ErrProviderCredentialValidationFailed:
					fmt.Printf("\n❌ GCP provider credential validation failed. Please check your credentials.\n")
					return
				case cliErrors.ErrProviderCredentialCreateFailed:
					fmt.Printf("\n❌ Failed to create GCP provider credential.\n")
					return
				case cliErrors.ErrAccessDenied:
					fmt.Printf("\n%s\n", err.Error())
					return
				default:
					fmt.Fprintf(os.Stderr, "❌ Failed to configure GCP provider: %v\n", err)
					os.Exit(1)
				}
			}

			fmt.Printf("✅ Successfully configured GCP provider!\n")
			fmt.Printf("   Provider ID: %s\n", result.ID)
			fmt.Printf("   Project ID: %s\n", gcpProjectID)
			fmt.Printf("   Service Account: %s\n", credentials["client_email"])
			fmt.Printf("   Secret Manager Location: %s\n", gcpSecretManagerLocation)
			if gcpPrefix != "" {
				fmt.Printf("   Secret Prefix: %s\n", gcpPrefix)
			}
			fmt.Printf("   Replication: %s\n", gcpReplication)
			if gcpDisableOlderVersions {
				fmt.Printf("   Disable Older Versions: enabled\n")
			}
			fmt.Printf("   Created: %s\n", result.CreatedAt)
		},
	}

	// Add flags
	cmd.Flags().StringVar(&gcpKeyFile, "key-file", "", "Path to GCP Service Account JSON file (required)")
	cmd.Flags().StringVar(&gcpProjectID, "project-id", "", "GCP Project ID (required)")
	cmd.Flags().StringVar(&gcpSecretManagerLocation, "secret-manager-location", "us-central1", "Secret Manager location")
	cmd.Flags().StringVar(&gcpPrefix, "prefix", "", "Secret name prefix")
	cmd.Flags().StringVar(&gcpReplication, "replication", "automatic", "Replication type: automatic, user-managed")
	cmd.Flags().BoolVar(&gcpDisableOlderVersions, "disable-older-versions", false, "Disable older secret versions when creating new ones")
	cmd.Flags().StringVar(&orgName, "org", "", "Organization name (required)")
	cmd.Flags().StringVar(&groupName, "group", "", "Secret group name (required)")
	cmd.Flags().StringVar(&envName, "env", "", "Environment name (required)")

	// Mark required flags
	cmd.MarkFlagRequired("key-file")
	cmd.MarkFlagRequired("project-id")
	cmd.MarkFlagRequired("org")
	cmd.MarkFlagRequired("group")
	cmd.MarkFlagRequired("env")

	return cmd
}
