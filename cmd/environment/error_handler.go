package environment

import (
	"fmt"

	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
)

// handleCommonEnvironmentErrors handles common errors that occur in environment operations
func handleCommonEnvironmentErrors(err error, operation string, envName string, logger *utils.Logger) error {
	switch {
	case err == cliErrors.ErrUnReachableBackend:
		fmt.Println("🚨 Oops! Kavach backend is currently down or not responding.")
		fmt.Println("📡 This may be due to server downtime or high request volume.")
		fmt.Println("📩 If this persists, please drop us a message at 👉 **gudikotieswarmani@gmail.com**")
		logger.Error("Backend unreachable during environment operation", err, map[string]interface{}{
			"operation": operation,
			"env":       envName,
		})
		return nil

	case err == cliErrors.ErrNotLoggedIn:
		fmt.Printf("\n🔒 You are not logged in. Please run 'kavach login' to log in.\n")
		logger.Warn("User not logged in during environment operation", map[string]interface{}{
			"operation": operation,
			"env":       envName,
		})
		return nil

	case err == cliErrors.ErrAccessDenied:
		fmt.Printf("\n%s\n", err.Error())
		logger.Warn("Access denied during environment operation", map[string]interface{}{
			"operation": operation,
			"env":       envName,
		})
		return nil

	case err == cliErrors.ErrInvalidToken:
		fmt.Printf("\n🔒 You are not logged in. Please run 'kavach login' to log in.\n")
		logger.Warn("User not logged in during environment operation", map[string]interface{}{
			"operation": operation,
			"env":       envName,
		})
		return nil
	}

	// Check if the error message contains authentication-related text
	if cliErrors.IsAuthenticationError(err) {
		fmt.Printf("\n🔑 Please login again, the session is expired, unable to authenticate you\n")
		logger.Warn("Authentication error during environment operation", map[string]interface{}{
			"operation": operation,
			"env":       envName,
		})
		return nil
	}

	// Return the original error if it's not a common error
	return err
}

// handleEnvironmentNotFoundError handles environment not found errors
func handleEnvironmentNotFoundError(envName string, logger *utils.Logger) {
	fmt.Printf("\n❌ Environment '%s' not found in the current secret group.\n", envName)
	fmt.Println("💡 Use 'kavach env list' to see available environments.")
	logger.Warn("Environment not found", map[string]interface{}{
		"env": envName,
	})
}

// handleEnvironmentValidationError handles environment validation errors
func handleEnvironmentValidationError(envName string, field string, logger *utils.Logger) {
	fmt.Printf("\n❌ Invalid environment '%s': %s is required.\n", envName, field)
	fmt.Println("💡 Please provide a valid value for this field.")
	logger.Warn("Environment validation failed", map[string]interface{}{
		"env":   envName,
		"field": field,
	})
}

// handleEnvironmentCreationError handles environment creation errors
func handleEnvironmentCreationError(envName string, err error, logger *utils.Logger) {
	fmt.Printf("\n❌ Failed to create environment '%s': %v\n", envName, err)
	fmt.Println("💡 Make sure the environment name is unique and you have create permissions.")
	logger.Error("Environment creation failed", err, map[string]interface{}{
		"env": envName,
	})
}

// handleEnvironmentActivationError handles environment activation errors
func handleEnvironmentActivationError(envName string, err error, logger *utils.Logger) {
	fmt.Printf("\n❌ Failed to activate environment '%s': %v\n", envName, err)
	fmt.Println("💡 Make sure you have access to this environment.")
	logger.Error("Environment activation failed", err, map[string]interface{}{
		"env": envName,
	})
}

// handleEnvironmentDeletionError handles environment deletion errors
func handleEnvironmentDeletionError(envName string, err error, logger *utils.Logger) {
	fmt.Printf("\n❌ Failed to delete environment '%s': %v\n", envName, err)
	fmt.Println("💡 Make sure the environment is empty and you have delete permissions.")
	logger.Error("Environment deletion failed", err, map[string]interface{}{
		"env": envName,
	})
}

// handlePermissionError handles permission-related errors
func handlePermissionError(operation string, envName string, err error, logger *utils.Logger) {
	fmt.Printf("\n🚫 Permission denied for operation '%s' on environment '%s'\n", operation, envName)
	fmt.Println("💡 Contact your environment administrator for access.")
	logger.Warn("Permission denied for environment operation", map[string]interface{}{
		"operation": operation,
		"env":       envName,
		"error":     err.Error(),
	})
}

// handleSecretGroupRequiredError handles missing secret group context errors
func handleSecretGroupRequiredError(logger *utils.Logger) {
	fmt.Println("\n❌ No secret group is currently active.")
	fmt.Println("💡 Use 'kavach group activate <group-name>' to activate a secret group first.")
	logger.Warn("No secret group active for environment operation")
}

// handleOrganizationRequiredError handles missing organization context errors
func handleOrganizationRequiredError(logger *utils.Logger) {
	fmt.Println("\n❌ No organization is currently active.")
	fmt.Println("💡 Use 'kavach org activate <org-name>' to activate an organization first.")
	logger.Warn("No organization active for environment operation")
}
