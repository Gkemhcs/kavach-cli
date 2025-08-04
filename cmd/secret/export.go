package secret

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/secret"
	"github.com/Gkemhcs/kavach-cli/internal/types"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewExportCommand creates a new command for exporting secrets to .env file
func NewExportCommand(logger *utils.Logger, cfg *config.Config, secretClient secret.SecretClient) *cobra.Command {
	var (
		versionID string
		envName   string
		orgName   string
		groupName string
		outputDir string
		fileName  string
	)

	cmd := &cobra.Command{
		Use:   "export",
		Short: "📤 Export secrets to .env file",
		Long: `Export secrets from a specific version to a .env file.

This command exports all secrets from a specific version to a .env file
instead of displaying them in the terminal. This is useful for creating
environment files for deployment or backup purposes.

Key concepts:
- Exports secrets to a .env file format (KEY=value)
- Supports custom output directory and filename
- Overwrites existing .env files
- Useful for deployment and backup scenarios
- Maintains secret values in file format

Prerequisites:
- Version must exist and be accessible
- You must have appropriate permissions for the environment
- Organization, secret group, and environment must be specified

The export process:
1. Retrieves secrets from the specified version
2. Formats them as KEY=value pairs
3. Writes to the specified .env file
4. Overwrites existing file if it exists

Use cases:
- Creating deployment environment files
- Backing up secret configurations
- Migrating secrets to other systems
- Local development setup
- Disaster recovery preparation

Examples:
  kavach secret export --version "abc123" --output-dir "./config"
  kavach secret export --version "prod-v1.2.3" --file "production.env"
  kavach secret export --version "dev-latest" --output-dir "./envs" --file "development.env"

Note: The .env file will contain actual secret values. Ensure proper
file permissions and security measures are in place.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runExportCommand(logger, secretClient, versionID, envName, orgName, groupName, outputDir, fileName)
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&versionID, "version", "v", "", "Version ID to export secrets from")
	cmd.Flags().StringVarP(&envName, "env", "e", "", "Environment name")
	cmd.Flags().StringVarP(&groupName, "group", "g", "", "Secret group name")
	cmd.Flags().StringVarP(&orgName, "org", "o", "", "Organization name")
	cmd.Flags().StringVarP(&outputDir, "output-dir", "d", ".", "Output directory for .env file")
	cmd.Flags().StringVarP(&fileName, "file", "f", ".env", "Output filename (default: .env)")

	// Mark required flags
	cmd.MarkFlagRequired("version")

	return cmd
}

// runExportCommand handles the export logic with proper error handling
func runExportCommand(logger *utils.Logger, secretClient secret.SecretClient, versionID, envName, orgName, groupName, outputDir, fileName string) error {
	// Validate required parameters
	if err := validateExportParams(versionID, outputDir, fileName); err != nil {
		logger.Warn("Invalid export parameters", map[string]interface{}{
			"cmd": "secret export",
			"err": err.Error(),
		})
		return err
	}

	// Set default values
	outputDir = getDefaultOutputDir(outputDir)
	fileName = getDefaultFileName(fileName)

	// Prepare file path
	filePath, err := prepareOutputPath(outputDir, fileName)
	if err != nil {
		logger.Error("Failed to prepare output path", err, map[string]interface{}{
			"cmd":       "secret export",
			"outputDir": outputDir,
			"fileName":  fileName,
		})
		return fmt.Errorf("failed to prepare output path: %w", err)
	}

	// Check for existing file and warn user
	if err := checkExistingFile(filePath); err != nil {
		logger.Warn("File check failed", map[string]interface{}{
			"cmd":      "secret export",
			"filePath": filePath,
			"err":      err.Error(),
		})
		// Continue anyway, just log the warning
	}

	// Get version details
	version, err := getVersionDetails(secretClient, orgName, groupName, envName, versionID, logger)
	if err != nil {
		return err
	}

	// Validate secrets exist
	if err := validateSecretsExist(version, versionID, logger); err != nil {
		return err
	}

	// Generate .env content
	envContent := generateEnvContent(version.Secrets)

	// Write file
	if err := writeEnvFile(filePath, envContent, versionID, logger); err != nil {
		return err
	}

	// Log success
	logExportSuccess(logger, versionID, filePath, version, len(envContent))

	return nil
}

// validateExportParams validates the export command parameters
func validateExportParams(versionID, outputDir, fileName string) error {
	if versionID == "" {
		return fmt.Errorf("version ID is required")
	}

	if outputDir != "" && !isValidPath(outputDir) {
		return fmt.Errorf("invalid output directory path: %s", outputDir)
	}

	if fileName != "" && !isValidFileName(fileName) {
		return fmt.Errorf("invalid filename: %s", fileName)
	}

	return nil
}

// isValidPath checks if a path is valid
func isValidPath(path string) bool {
	// Basic validation - could be enhanced with more checks
	return len(path) > 0 && len(path) < 4096
}

// isValidFileName checks if a filename is valid
func isValidFileName(name string) bool {
	// Basic validation - could be enhanced with more checks
	return len(name) > 0 && len(name) < 255 && !containsInvalidChars(name)
}

// containsInvalidChars checks for invalid characters in filename
func containsInvalidChars(name string) bool {
	invalidChars := []string{"<", ">", ":", "\"", "|", "?", "*", "\\", "/"}
	for _, char := range invalidChars {
		if strings.Contains(name, char) {
			return true
		}
	}
	return false
}

// getDefaultOutputDir returns the default output directory
func getDefaultOutputDir(outputDir string) string {
	if outputDir == "" {
		return "."
	}
	return outputDir
}

// getDefaultFileName returns the default filename
func getDefaultFileName(fileName string) string {
	if fileName == "" {
		return ".env"
	}
	return fileName
}

// prepareOutputPath creates the output directory and returns the full file path
func prepareOutputPath(outputDir, fileName string) (string, error) {
	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Construct full file path
	filePath := filepath.Join(outputDir, fileName)

	// Validate the final path
	if !isValidPath(filePath) {
		return "", fmt.Errorf("invalid file path: %s", filePath)
	}

	return filePath, nil
}

// checkExistingFile checks if file exists and warns user
func checkExistingFile(filePath string) error {
	if _, err := os.Stat(filePath); err == nil {
		fmt.Printf("\n⚠️  File '%s' already exists. It will be overwritten.\n", filePath)
	}
	return nil
}

// getVersionDetails retrieves version details with proper error handling
func getVersionDetails(secretClient secret.SecretClient, orgName, groupName, envName, versionID string, logger *utils.Logger) (*types.SecretVersionDetailResponse, error) {
	version, err := secretClient.GetVersionDetails(orgName, groupName, envName, versionID)
	if err != nil {
		return nil, handleVersionDetailsError(err, versionID, logger)
	}
	return version, nil
}

// handleVersionDetailsError handles errors from getting version details
func handleVersionDetailsError(err error, versionID string, logger *utils.Logger) error {
	switch err {
	case cliErrors.ErrSecretVersionNotFound:
		logger.Warn("Secret version not found during export", map[string]interface{}{
			"cmd":     "secret export",
			"version": versionID,
		})
		return fmt.Errorf("❌ %v", err)
	case cliErrors.ErrSecretNotFound:
		logger.Warn("Secret not found during export", map[string]interface{}{
			"cmd":     "secret export",
			"version": versionID,
		})
		return fmt.Errorf("❌ %v", err)
	case cliErrors.ErrDecryptionFailed:
		logger.Error("Decryption failed during export", err, map[string]interface{}{
			"cmd":     "secret export",
			"version": versionID,
		})
		return fmt.Errorf("❌ %v", err)
	case cliErrors.ErrInternalServer:
		logger.Error("Internal server error during export", err, map[string]interface{}{
			"cmd":     "secret export",
			"version": versionID,
		})
		return fmt.Errorf("❌ %v", err)
	case cliErrors.ErrNotLoggedIn:
		logger.Warn("User not logged in during export", map[string]interface{}{
			"cmd":     "secret export",
			"version": versionID,
		})
		return fmt.Errorf("❌ %v", err)
	case cliErrors.ErrAccessDenied:
		logger.Warn("Access denied during secret export", map[string]interface{}{
			"cmd":     "secret export",
			"version": versionID,
		})
		return fmt.Errorf("%s", err.Error())
	default:
		logger.Error("Failed to get version details", err, map[string]interface{}{
			"cmd":     "secret export",
			"version": versionID,
		})
		return fmt.Errorf("❌ Failed to get version details: %v", err)
	}
}

// validateSecretsExist checks if secrets exist in the version
func validateSecretsExist(version *types.SecretVersionDetailResponse, versionID string, logger *utils.Logger) error {
	if len(version.Secrets) == 0 {
		logger.Info("No secrets found in version for export", map[string]interface{}{
			"cmd":     "secret export",
			"version": versionID,
		})
		return fmt.Errorf("📝 No secrets found in version '%s' to export", versionID)
	}
	return nil
}

// generateEnvContent creates the .env file content
func generateEnvContent(secrets []types.Secret) string {
	var content strings.Builder
	for _, secret := range secrets {
		content.WriteString(fmt.Sprintf("%s=%s\n", secret.Name, secret.Value))
	}
	return content.String()
}

// writeEnvFile writes the .env content to file
func writeEnvFile(filePath, content, versionID string, logger *utils.Logger) error {
	if err := os.WriteFile(filePath, []byte(content), 0600); err != nil {
		logger.Error("Failed to write .env file", err, map[string]interface{}{
			"cmd":      "secret export",
			"version":  versionID,
			"filePath": filePath,
		})
		return fmt.Errorf("❌ Failed to write .env file: %v", err)
	}
	return nil
}

// logExportSuccess logs the successful export
func logExportSuccess(logger *utils.Logger, versionID, filePath string, version *types.SecretVersionDetailResponse, contentSize int) {
	// Print success message
	fmt.Printf("\n✅ Successfully exported secrets to .env file!\n")
	fmt.Printf("   File Path: %s\n", filePath)
	fmt.Printf("   Version ID: %s\n", version.ID)
	fmt.Printf("   Environment: %s\n", version.EnvironmentID)
	fmt.Printf("   Total Secrets: %d\n", len(version.Secrets))
	fmt.Printf("   File Size: %d bytes\n", contentSize)

	// Log success
	logger.Info("Successfully exported secrets to .env file", map[string]interface{}{
		"cmd":         "secret export",
		"version":     versionID,
		"filePath":    filePath,
		"secretCount": len(version.Secrets),
		"fileSize":    contentSize,
	})
}
