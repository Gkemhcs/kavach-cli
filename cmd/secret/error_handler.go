package secret

import (
	"fmt"

	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
)

// handleCommonSecretErrors handles common errors that occur in secret operations
func handleCommonSecretErrors(err error, operation string, secretName string, logger *utils.Logger) error {
	switch {
	case err == cliErrors.ErrUnReachableBackend:
		fmt.Println("🚨 Oops! Kavach backend is currently down or not responding.")
		fmt.Println("📡 This may be due to server downtime or high request volume.")
		fmt.Println("📩 If this persists, please drop us a message at 👉 **gudikotieswarmani@gmail.com**")
		logger.Error("Backend unreachable during secret operation", err, map[string]interface{}{
			"operation": operation,
			"secret":    secretName,
		})
		return nil

	case err == cliErrors.ErrNotLoggedIn:
		fmt.Printf("\n🔒 You are not logged in. Please run 'kavach login' to log in.\n")
		logger.Warn("User not logged in during secret operation", map[string]interface{}{
			"operation": operation,
			"secret":    secretName,
		})
		return nil

	case err == cliErrors.ErrAccessDenied:
		fmt.Printf("\n%s\n", err.Error())
		logger.Warn("Access denied during secret operation", map[string]interface{}{
			"operation": operation,
			"secret":    secretName,
		})
		return nil

	case err == cliErrors.ErrInvalidToken:
		fmt.Printf("\n🔒 You are not logged in. Please run 'kavach login' to log in.\n")
		logger.Warn("User not logged in during secret operation", map[string]interface{}{
			"operation": operation,
			"secret":    secretName,
		})
		return nil
	}

	// Check if the error message contains authentication-related text
	if cliErrors.IsAuthenticationError(err) {
		fmt.Printf("\n🔑 Please login again, the session is expired, unable to authenticate you\n")
		logger.Warn("Authentication error during secret operation", map[string]interface{}{
			"operation": operation,
			"secret":    secretName,
		})
		return nil
	}

	// Return the original error if it's not a common error
	return err
}

// handleSecretNotFoundError handles secret not found errors
func handleSecretNotFoundError(secretName string, logger *utils.Logger) {
	fmt.Printf("\n❌ Secret '%s' not found in the current environment.\n", secretName)
	fmt.Println("💡 Use 'kavach secret list' to see available secrets.")
	logger.Warn("Secret not found", map[string]interface{}{
		"secret": secretName,
	})
}

// handleSecretValidationError handles secret validation errors
func handleSecretValidationError(secretName string, field string, logger *utils.Logger) {
	fmt.Printf("\n❌ Invalid secret '%s': %s is required.\n", secretName, field)
	fmt.Println("💡 Please provide a valid value for this field.")
	logger.Warn("Secret validation failed", map[string]interface{}{
		"secret": secretName,
		"field":  field,
	})
}

// handleStagingError handles staging-related errors
func handleStagingError(operation string, secretName string, err error, logger *utils.Logger) {
	fmt.Printf("\n❌ Failed to %s secret '%s' in staging: %v\n", operation, secretName, err)
	fmt.Println("💡 Make sure the staging area is accessible and you have proper permissions.")
	logger.Error("Staging operation failed", err, map[string]interface{}{
		"operation": operation,
		"secret":    secretName,
	})
}

// handleCommitError handles commit-related errors
func handleCommitError(secretName string, err error, logger *utils.Logger) {
	fmt.Printf("\n❌ Failed to commit secret '%s': %v\n", secretName, err)
	fmt.Println("💡 Make sure you have write permissions in the current environment.")
	logger.Error("Secret commit failed", err, map[string]interface{}{
		"secret": secretName,
	})
}

// handleSyncError handles sync-related errors
func handleSyncError(provider string, err error, logger *utils.Logger) {
	fmt.Printf("\n❌ Failed to sync secrets to %s: %v\n", provider, err)
	fmt.Println("💡 Check your provider credentials and permissions.")
	logger.Error("Secret sync failed", err, map[string]interface{}{
		"provider": provider,
	})
}

// handleExportError handles export-related errors
func handleExportError(secretName string, err error, logger *utils.Logger) {
	fmt.Printf("\n❌ Failed to export secret '%s': %v\n", secretName, err)
	fmt.Println("💡 Make sure you have read permissions for this secret.")
	logger.Error("Secret export failed", err, map[string]interface{}{
		"secret": secretName,
	})
}

// handleRollbackError handles rollback-related errors
func handleRollbackError(secretName string, version string, err error, logger *utils.Logger) {
	fmt.Printf("\n❌ Failed to rollback secret '%s' to version '%s': %v\n", secretName, version, err)
	fmt.Println("💡 Make sure the version exists and you have rollback permissions.")
	logger.Error("Secret rollback failed", err, map[string]interface{}{
		"secret":  secretName,
		"version": version,
	})
}

// handleDiffError handles diff-related errors
func handleDiffError(secretName string, version1 string, version2 string, err error, logger *utils.Logger) {
	fmt.Printf("\n❌ Failed to compare versions '%s' and '%s' for secret '%s': %v\n", version1, version2, secretName, err)
	fmt.Println("💡 Make sure both versions exist and you have read permissions.")
	logger.Error("Secret diff failed", err, map[string]interface{}{
		"secret":   secretName,
		"version1": version1,
		"version2": version2,
	})
}
