package usergroup

import (
	"fmt"

	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
)

// handleCommonUserGroupErrors handles common errors that occur in user group operations
func handleCommonUserGroupErrors(err error, operation string, groupName string, logger *utils.Logger) error {
	switch {
	case err == cliErrors.ErrUnReachableBackend:
		fmt.Println("🚨 Oops! Kavach backend is currently down or not responding.")
		fmt.Println("📡 This may be due to server downtime or high request volume.")
		fmt.Println("📩 If this persists, please drop us a message at 👉 **gudikotieswarmani@gmail.com**")
		logger.Error("Backend unreachable during user group operation", err, map[string]interface{}{
			"operation": operation,
			"group":     groupName,
		})
		return nil

	case err == cliErrors.ErrNotLoggedIn:
		fmt.Printf("\n🔒 You are not logged in. Please run 'kavach login' to log in.\n")
		logger.Warn("User not logged in during user group operation", map[string]interface{}{
			"operation": operation,
			"group":     groupName,
		})
		return nil

	case err == cliErrors.ErrAccessDenied:
		fmt.Printf("\n%s\n", err.Error())
		logger.Warn("Access denied during user group operation", map[string]interface{}{
			"operation": operation,
			"group":     groupName,
		})
		return nil

	case err == cliErrors.ErrInvalidToken:
		fmt.Printf("\n🔒 You are not logged in. Please run 'kavach login' to log in.\n")
		logger.Warn("User not logged in during user group operation", map[string]interface{}{
			"operation": operation,
			"group":     groupName,
		})
		return nil
	}

	// Check if the error message contains authentication-related text
	if cliErrors.IsAuthenticationError(err) {
		fmt.Printf("\n🔑 Please login again, the session is expired, unable to authenticate you\n")
		logger.Warn("Authentication error during user group operation", map[string]interface{}{
			"operation": operation,
			"group":     groupName,
		})
		return nil
	}

	// Return the original error if it's not a common error
	return err
}

// handleUserGroupNotFoundError handles user group not found errors
func handleUserGroupNotFoundError(groupName string, logger *utils.Logger) {
	fmt.Printf("\n❌ User group '%s' not found in the current organization.\n", groupName)
	fmt.Println("💡 Use 'kavach user-group list' to see available user groups.")
	logger.Warn("User group not found", map[string]interface{}{
		"group": groupName,
	})
}

// handleUserGroupValidationError handles user group validation errors
func handleUserGroupValidationError(groupName string, field string, logger *utils.Logger) {
	fmt.Printf("\n❌ Invalid user group '%s': %s is required.\n", groupName, field)
	fmt.Println("💡 Please provide a valid value for this field.")
	logger.Warn("User group validation failed", map[string]interface{}{
		"group": groupName,
		"field": field,
	})
}

// handleUserGroupCreationError handles user group creation errors
func handleUserGroupCreationError(groupName string, err error, logger *utils.Logger) {
	fmt.Printf("\n❌ Failed to create user group '%s': %v\n", groupName, err)
	fmt.Println("💡 Make sure the group name is unique and you have create permissions.")
	logger.Error("User group creation failed", err, map[string]interface{}{
		"group": groupName,
	})
}

// handleUserGroupDeletionError handles user group deletion errors
func handleUserGroupDeletionError(groupName string, err error, logger *utils.Logger) {
	fmt.Printf("\n❌ Failed to delete user group '%s': %v\n", groupName, err)
	fmt.Println("💡 Make sure the group is empty and you have delete permissions.")
	logger.Error("User group deletion failed", err, map[string]interface{}{
		"group": groupName,
	})
}

// handleMemberOperationError handles member-related operation errors
func handleMemberOperationError(operation string, groupName string, userName string, err error, logger *utils.Logger) {
	fmt.Printf("\n❌ Failed to %s user '%s' from group '%s': %v\n", operation, userName, groupName, err)
	fmt.Printf("💡 Make sure the user exists and you have permission to %s them.\n", operation)
	logger.Error("Member operation failed", err, map[string]interface{}{
		"operation": operation,
		"group":     groupName,
		"user":      userName,
	})
}

// handleUserNotFoundError handles user not found errors
func handleUserNotFoundError(userName string, logger *utils.Logger) {
	fmt.Printf("\n❌ User '%s' not found in the current organization.\n", userName)
	fmt.Println("💡 Make sure the user exists and you have access to them.")
	logger.Warn("User not found", map[string]interface{}{
		"user": userName,
	})
}

// handleDuplicateMemberError handles duplicate member errors
func handleDuplicateMemberError(groupName string, userName string, logger *utils.Logger) {
	fmt.Printf("\n❌ User '%s' is already a member of group '%s'.\n", userName, groupName)
	fmt.Println("💡 No action needed - the user is already in the group.")
	logger.Warn("Duplicate member operation", map[string]interface{}{
		"group": groupName,
		"user":  userName,
	})
}

// handleMemberNotFoundError handles member not found errors
func handleMemberNotFoundError(groupName string, userName string, logger *utils.Logger) {
	fmt.Printf("\n❌ User '%s' is not a member of group '%s'.\n", userName, groupName)
	fmt.Println("💡 No action needed - the user is not in the group.")
	logger.Warn("Member not found in group", map[string]interface{}{
		"group": groupName,
		"user":  userName,
	})
}

// handlePermissionError handles permission-related errors
func handlePermissionError(operation string, groupName string, err error, logger *utils.Logger) {
	fmt.Printf("\n🚫 Permission denied for operation '%s' on user group '%s'\n", operation, groupName)
	fmt.Println("💡 Contact your organization administrator for access.")
	logger.Warn("Permission denied for user group operation", map[string]interface{}{
		"operation": operation,
		"group":     groupName,
		"error":     err.Error(),
	})
}

// handleOrganizationRequiredError handles missing organization context errors
func handleOrganizationRequiredError(logger *utils.Logger) {
	fmt.Println("\n❌ No organization is currently active.")
	fmt.Println("💡 Use 'kavach org activate <org-name>' to activate an organization first.")
	logger.Warn("No organization active for user group operation")
}
