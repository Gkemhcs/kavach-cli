package org

import (
	"fmt"

	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
)

// handleCommonOrgErrors handles common errors that occur in organization operations
func handleCommonOrgErrors(err error, operation string, orgName string, logger *utils.Logger) error {
	switch {
	case err == cliErrors.ErrUnReachableBackend:
		fmt.Println("🚨 Oops! Kavach backend is currently down or not responding.")
		fmt.Println("📡 This may be due to server downtime or high request volume.")
		fmt.Println("📩 If this persists, please drop us a message at 👉 **gudikotieswarmani@gmail.com**")
		logger.Error("Backend unreachable during org operation", err, map[string]interface{}{
			"operation": operation,
			"org":       orgName,
		})
		return nil

	case err == cliErrors.ErrNotLoggedIn:
		fmt.Printf("\n🔒 You are not logged in. Please run 'kavach login' to log in.\n")
		logger.Warn("User not logged in during org operation", map[string]interface{}{
			"operation": operation,
			"org":       orgName,
		})
		return nil

	case err == cliErrors.ErrDuplicateOrganisation:
		fmt.Printf("\n❌ Organization '%s' already exists. Please choose a different name.\n", orgName)
		logger.Warn("Duplicate organization during operation", map[string]interface{}{
			"operation": operation,
			"org":       orgName,
		})
		return nil

	case err == cliErrors.ErrAccessDenied:
		fmt.Printf("\n%s\n", err.Error())
		logger.Warn("Access denied during org operation", map[string]interface{}{
			"operation": operation,
			"org":       orgName,
		})
		return nil

	case err == cliErrors.ErrInvalidToken:
		fmt.Printf("\n🔒 You are not logged in. Please run 'kavach login' to log in.\n")
		logger.Warn("User not logged in during org operation", map[string]interface{}{
			"operation": operation,
			"org":       orgName,
		})
		return nil
	}

	// Check if the error message contains authentication-related text
	if cliErrors.IsAuthenticationError(err) {
		fmt.Printf("\n🔑 Please login again, the session is expired, unable to authenticate you\n")
		logger.Warn("Authentication error during org operation", map[string]interface{}{
			"operation": operation,
			"org":       orgName,
		})
		return nil
	}

	// Return the original error if it's not a common error
	return err
}

// handleOrgNotFoundError handles organization not found errors
func handleOrgNotFoundError(orgName string, logger *utils.Logger) {
	fmt.Printf("\n❌ Organization '%s' not found or you don't have access to it.\n", orgName)
	fmt.Println("💡 Use 'kavach org list' to see available organizations.")
	logger.Warn("Organization not found", map[string]interface{}{
		"org": orgName,
	})
}

// handleOrgActivationError handles organization activation errors
func handleOrgActivationError(orgName string, err error, logger *utils.Logger) {
	fmt.Printf("\n❌ Failed to activate organization '%s': %v\n", orgName, err)
	fmt.Println("💡 Make sure you have access to this organization.")
	logger.Error("Failed to activate organization", err, map[string]interface{}{
		"org": orgName,
	})
}

// handleOrgDeletionError handles organization deletion errors
func handleOrgDeletionError(orgName string, err error, logger *utils.Logger) {
	fmt.Printf("\n❌ Failed to delete organization '%s': %v\n", orgName, err)
	fmt.Println("💡 Make sure the organization is empty and you have admin privileges.")
	logger.Error("Failed to delete organization", err, map[string]interface{}{
		"org": orgName,
	})
}

// handlePermissionError handles permission-related errors
func handlePermissionError(operation string, orgName string, err error, logger *utils.Logger) {
	fmt.Printf("\n🚫 Permission denied for operation '%s' on organization '%s'\n", operation, orgName)
	fmt.Println("💡 Contact your organization administrator for access.")
	logger.Warn("Permission denied for org operation", map[string]interface{}{
		"operation": operation,
		"org":       orgName,
		"error":     err.Error(),
	})
}
