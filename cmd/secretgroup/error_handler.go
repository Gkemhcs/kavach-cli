package secretgroup

import (
	"fmt"

	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
)

// handleCommonSecretGroupErrors handles common errors that occur in secret group operations
func handleCommonSecretGroupErrors(err error, operation string, groupName string, logger *utils.Logger) error {
	switch {
	case err == cliErrors.ErrUnReachableBackend:
		fmt.Println("🚨 Oops! Kavach backend is currently down or not responding.")
		fmt.Println("📡 This may be due to server downtime or high request volume.")
		fmt.Println("📩 If this persists, please drop us a message at 👉 **gudikotieswarmani@gmail.com**")
		logger.Error("Backend unreachable during secret group operation", err, map[string]interface{}{
			"operation": operation,
			"group":     groupName,
		})
		return nil

	case err == cliErrors.ErrNotLoggedIn:
		fmt.Printf("\n🔒 You are not logged in. Please run 'kavach login' to log in.\n")
		logger.Warn("User not logged in during secret group operation", map[string]interface{}{
			"operation": operation,
			"group":     groupName,
		})
		return nil

	case err == cliErrors.ErrAccessDenied:
		fmt.Printf("\n%s\n", err.Error())
		logger.Warn("Access denied during secret group operation", map[string]interface{}{
			"operation": operation,
			"group":     groupName,
		})
		return nil

	case err == cliErrors.ErrInvalidToken:
		fmt.Printf("\n🔒 You are not logged in. Please run 'kavach login' to log in.\n")
		logger.Warn("User not logged in during secret group operation", map[string]interface{}{
			"operation": operation,
			"group":     groupName,
		})
		return nil
	}

	// Check if the error message contains authentication-related text
	if cliErrors.IsAuthenticationError(err) {
		fmt.Printf("\n🔑 Please login again, the session is expired, unable to authenticate you\n")
		logger.Warn("Authentication error during secret group operation", map[string]interface{}{
			"operation": operation,
			"group":     groupName,
		})
		return nil
	}

	// Return the original error if it's not a common error
	return err
}

// handleSecretGroupNotFoundError handles secret group not found errors
func handleSecretGroupNotFoundError(groupName string, logger *utils.Logger) {
	fmt.Printf("\n❌ Secret group '%s' not found in the current organization.\n", groupName)
	fmt.Println("💡 Use 'kavach group list' to see available secret groups.")
	logger.Warn("Secret group not found", map[string]interface{}{
		"group": groupName,
	})
}

// handleSecretGroupValidationError handles secret group validation errors
func handleSecretGroupValidationError(groupName string, field string, logger *utils.Logger) {
	fmt.Printf("\n❌ Invalid secret group '%s': %s is required.\n", groupName, field)
	fmt.Println("💡 Please provide a valid value for this field.")
	logger.Warn("Secret group validation failed", map[string]interface{}{
		"group": groupName,
		"field": field,
	})
}

// handleSecretGroupCreationError handles secret group creation errors
func handleSecretGroupCreationError(groupName string, err error, logger *utils.Logger) {
	fmt.Printf("\n❌ Failed to create secret group '%s': %v\n", groupName, err)
	fmt.Println("💡 Make sure the group name is unique and you have create permissions.")
	logger.Error("Secret group creation failed", err, map[string]interface{}{
		"group": groupName,
	})
}

// handleSecretGroupActivationError handles secret group activation errors
func handleSecretGroupActivationError(groupName string, err error, logger *utils.Logger) {
	fmt.Printf("\n❌ Failed to activate secret group '%s': %v\n", groupName, err)
	fmt.Println("💡 Make sure you have access to this secret group.")
	logger.Error("Secret group activation failed", err, map[string]interface{}{
		"group": groupName,
	})
}

// handleSecretGroupDeletionError handles secret group deletion errors
func handleSecretGroupDeletionError(groupName string, err error, logger *utils.Logger) {
	fmt.Printf("\n❌ Failed to delete secret group '%s': %v\n", groupName, err)
	fmt.Println("💡 Make sure the group is empty and you have delete permissions.")
	logger.Error("Secret group deletion failed", err, map[string]interface{}{
		"group": groupName,
	})
}

// handlePermissionError handles permission-related errors
func handlePermissionError(operation string, groupName string, err error, logger *utils.Logger) {
	fmt.Printf("\n🚫 Permission denied for operation '%s' on secret group '%s'\n", operation, groupName)
	fmt.Println("💡 Contact your organization administrator for access.")
	logger.Warn("Permission denied for secret group operation", map[string]interface{}{
		"operation": operation,
		"group":     groupName,
		"error":     err.Error(),
	})
}

// handleOrganizationRequiredError handles missing organization context errors
func handleOrganizationRequiredError(logger *utils.Logger) {
	fmt.Println("\n❌ No organization is currently active.")
	fmt.Println("💡 Use 'kavach org activate <org-name>' to activate an organization first.")
	logger.Warn("No organization active for secret group operation")
}

// handleDuplicateSecretGroupError handles duplicate secret group errors
func handleDuplicateSecretGroupError(groupName string, logger *utils.Logger) {
	fmt.Printf("\n❌ Secret group '%s' already exists in this organization.\n", groupName)
	fmt.Println("💡 Please choose a different name for your secret group.")
	logger.Warn("Duplicate secret group name", map[string]interface{}{
		"group": groupName,
	})
}

// handleSecretGroupNotEmptyError handles non-empty secret group deletion errors
func handleSecretGroupNotEmptyError(groupName string, logger *utils.Logger) {
	fmt.Printf("\n❌ Cannot delete secret group '%s' because it contains environments or secrets.\n", groupName)
	fmt.Println("💡 Please delete all environments and secrets first, then try again.")
	logger.Warn("Attempted to delete non-empty secret group", map[string]interface{}{
		"group": groupName,
	})
}
