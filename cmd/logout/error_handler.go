package logout

import (
	"fmt"

	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
)

// handleCommonLogoutErrors handles common errors that occur in logout operations
func handleCommonLogoutErrors(err error, logger *utils.Logger) error {
	switch {
	case err == cliErrors.ErrUnReachableBackend:
		fmt.Println("🚨 Oops! Kavach backend is currently down or not responding.")
		fmt.Println("📡 This may be due to server downtime or high request volume.")
		fmt.Println("📩 If this persists, please drop us a message at 👉 **gudikotieswarmani@gmail.com**")
		logger.Error("Backend unreachable during logout", err, map[string]interface{}{
			"operation": "logout",
		})
		return nil

	case err == cliErrors.ErrNotLoggedIn:
		fmt.Printf("\n🔒 You are not currently logged in to Kavach.\n")
		logger.Warn("User not logged in during logout")
		return nil

	case err == cliErrors.ErrInvalidToken:
		fmt.Printf("\n🔒 Invalid or expired token. Proceeding with local cleanup.\n")
		logger.Warn("Invalid token during logout")
		return nil
	}

	// Check if the error message contains authentication-related text
	if cliErrors.IsAuthenticationError(err) {
		fmt.Printf("\n🔑 Authentication error. Proceeding with local cleanup.\n")
		logger.Warn("Authentication error during logout")
		return nil
	}

	// Return the original error if it's not a common error
	return err
}

// handleTokenRemovalError handles errors when removing stored tokens
func handleTokenRemovalError(err error, logger *utils.Logger) {
	fmt.Println("\n💾 Failed to remove stored authentication token.")
	fmt.Println("💡 You may need to manually delete the token file.")
	logger.Error("Token removal failed", err, map[string]interface{}{
		"operation": "logout",
	})
}

// handleSessionCleanupError handles errors during session cleanup
func handleSessionCleanupError(err error, logger *utils.Logger) {
	fmt.Println("\n🧹 Failed to clean up session data.")
	fmt.Println("💡 Some session data may remain on your system.")
	logger.Error("Session cleanup failed", err, map[string]interface{}{
		"operation": "logout",
	})
}

// handleNotLoggedInError handles when user tries to logout without being logged in
func handleNotLoggedInError(logger *utils.Logger) {
	fmt.Println("\nℹ️  You are not currently logged in to Kavach.")
	fmt.Println("💡 No action needed.")
	logger.Info("Logout attempted when not logged in", map[string]interface{}{
		"operation": "logout",
	})
}

// handleLogoutSuccess handles successful logout
func handleLogoutSuccess(logger *utils.Logger) {
	fmt.Println("\n✅ Successfully logged out of Kavach.")
	fmt.Println("💡 Your authentication tokens have been removed.")
	logger.Info("User logged out successfully", map[string]interface{}{
		"operation": "logout",
	})
}
