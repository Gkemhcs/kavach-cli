package provider

import (
	"fmt"

	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
)

// handleCommonProviderErrors handles common errors that occur in provider operations
func handleCommonProviderErrors(err error, operation string, providerName string, logger *utils.Logger) error {
	switch {
	case err == cliErrors.ErrUnReachableBackend:
		fmt.Println("🚨 Oops! Kavach backend is currently down or not responding.")
		fmt.Println("📡 This may be due to server downtime or high request volume.")
		fmt.Println("📩 If this persists, please drop us a message at 👉 **gudikotieswarmani@gmail.com**")
		logger.Error("Backend unreachable during provider operation", err, map[string]interface{}{
			"operation": operation,
			"provider":  providerName,
		})
		return nil

	case err == cliErrors.ErrNotLoggedIn:
		fmt.Printf("\n🔒 You are not logged in. Please run 'kavach login' to log in.\n")
		logger.Warn("User not logged in during provider operation", map[string]interface{}{
			"operation": operation,
			"provider":  providerName,
		})
		return nil

	case err == cliErrors.ErrAccessDenied:
		fmt.Printf("\n%s\n", err.Error())
		logger.Warn("Access denied during provider operation", map[string]interface{}{
			"operation": operation,
			"provider":  providerName,
		})
		return nil

	case err == cliErrors.ErrInvalidToken:
		fmt.Printf("\n🔒 You are not logged in. Please run 'kavach login' to log in.\n")
		logger.Warn("User not logged in during provider operation", map[string]interface{}{
			"operation": operation,
			"provider":  providerName,
		})
		return nil
	}

	// Check if the error message contains authentication-related text
	if cliErrors.IsAuthenticationError(err) {
		fmt.Printf("\n🔑 Please login again, the session is expired, unable to authenticate you\n")
		logger.Warn("Authentication error during provider operation", map[string]interface{}{
			"operation": operation,
			"provider":  providerName,
		})
		return nil
	}

	// Return the original error if it's not a common error
	return err
}

// handleProviderNotFoundError handles provider not found errors
func handleProviderNotFoundError(providerName string, logger *utils.Logger) {
	fmt.Printf("\n❌ Provider '%s' not found in the current environment.\n", providerName)
	fmt.Println("💡 Use 'kavach provider list' to see available providers.")
	logger.Warn("Provider not found", map[string]interface{}{
		"provider": providerName,
	})
}

// handleProviderValidationError handles provider validation errors
func handleProviderValidationError(providerName string, field string, logger *utils.Logger) {
	fmt.Printf("\n❌ Invalid provider '%s': %s is required.\n", providerName, field)
	fmt.Println("💡 Please provide a valid value for this field.")
	logger.Warn("Provider validation failed", map[string]interface{}{
		"provider": providerName,
		"field":    field,
	})
}

// handleProviderConfigurationError handles provider configuration errors
func handleProviderConfigurationError(providerName string, err error, logger *utils.Logger) {
	fmt.Printf("\n❌ Failed to configure provider '%s': %v\n", providerName, err)
	fmt.Println("💡 Check your credentials and permissions for this provider.")
	logger.Error("Provider configuration failed", err, map[string]interface{}{
		"provider": providerName,
	})
}

// handleProviderUpdateError handles provider update errors
func handleProviderUpdateError(providerName string, err error, logger *utils.Logger) {
	fmt.Printf("\n❌ Failed to update provider '%s': %v\n", providerName, err)
	fmt.Println("💡 Make sure you have update permissions for this provider.")
	logger.Error("Provider update failed", err, map[string]interface{}{
		"provider": providerName,
	})
}

// handleProviderDeletionError handles provider deletion errors
func handleProviderDeletionError(providerName string, err error, logger *utils.Logger) {
	fmt.Printf("\n❌ Failed to delete provider '%s': %v\n", providerName, err)
	fmt.Println("💡 Make sure you have delete permissions and no active syncs.")
	logger.Error("Provider deletion failed", err, map[string]interface{}{
		"provider": providerName,
	})
}

// handleProviderSyncError handles provider sync errors
func handleProviderSyncError(providerName string, err error, logger *utils.Logger) {
	fmt.Printf("\n❌ Failed to sync with provider '%s': %v\n", providerName, err)
	fmt.Println("💡 Check your provider credentials and network connectivity.")
	logger.Error("Provider sync failed", err, map[string]interface{}{
		"provider": providerName,
	})
}

// handleCredentialError handles credential-related errors
func handleCredentialError(providerName string, field string, logger *utils.Logger) {
	fmt.Printf("\n❌ Invalid credentials for provider '%s': %s is missing or invalid.\n", providerName, field)
	fmt.Println("💡 Please check your credentials and try again.")
	logger.Warn("Provider credential error", map[string]interface{}{
		"provider": providerName,
		"field":    field,
	})
}

// handleDuplicateProviderError handles duplicate provider errors
func handleDuplicateProviderError(providerName string, logger *utils.Logger) {
	fmt.Printf("\n❌ Provider '%s' is already configured in this environment.\n", providerName)
	fmt.Println("💡 Use 'kavach provider update' to modify existing configuration.")
	logger.Warn("Duplicate provider configuration", map[string]interface{}{
		"provider": providerName,
	})
}

// handlePermissionError handles permission-related errors
func handlePermissionError(operation string, providerName string, err error, logger *utils.Logger) {
	fmt.Printf("\n🚫 Permission denied for operation '%s' on provider '%s'\n", operation, providerName)
	fmt.Println("💡 Contact your environment administrator for access.")
	logger.Warn("Permission denied for provider operation", map[string]interface{}{
		"operation": operation,
		"provider":  providerName,
		"error":     err.Error(),
	})
}

// handleEnvironmentRequiredError handles missing environment context errors
func handleEnvironmentRequiredError(logger *utils.Logger) {
	fmt.Println("\n❌ No environment is currently active.")
	fmt.Println("💡 Use 'kavach env activate <env-name>' to activate an environment first.")
	logger.Warn("No environment active for provider operation")
}

// handleSecretGroupRequiredError handles missing secret group context errors
func handleSecretGroupRequiredError(logger *utils.Logger) {
	fmt.Println("\n❌ No secret group is currently active.")
	fmt.Println("💡 Use 'kavach group activate <group-name>' to activate a secret group first.")
	logger.Warn("No secret group active for provider operation")
}

// handleOrganizationRequiredError handles missing organization context errors
func handleOrganizationRequiredError(logger *utils.Logger) {
	fmt.Println("\n❌ No organization is currently active.")
	fmt.Println("💡 Use 'kavach org activate <org-name>' to activate an organization first.")
	logger.Warn("No organization active for provider operation")
}
