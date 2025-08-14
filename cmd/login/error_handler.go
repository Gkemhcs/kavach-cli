package login

import (
	"fmt"

	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
)

// handleCommonLoginErrors handles common errors that occur in login operations
func handleCommonLoginErrors(err error, logger *utils.Logger) error {
	switch {
	case err == cliErrors.ErrUnReachableBackend:
		fmt.Println("🚨 Oops! Kavach backend is currently down or not responding.")
		fmt.Println("📡 This may be due to server downtime or high request volume.")
		fmt.Println("📩 If this persists, please drop us a message at 👉 **gudikotieswarmani@gmail.com**")
		logger.Error("Backend unreachable during login", err, map[string]interface{}{
			"operation": "login",
		})
		return nil

	case err == cliErrors.ErrInvalidToken:
		fmt.Printf("\n🔒 Invalid or expired token. Please try logging in again.\n")
		logger.Warn("Invalid token during login")
		return nil
	}

	// Check if the error message contains authentication-related text
	if cliErrors.IsAuthenticationError(err) {
		fmt.Printf("\n🔑 Authentication failed. Please check your credentials and try again.\n")
		logger.Warn("Authentication error during login")
		return nil
	}

	// Return the original error if it's not a common error
	return err
}

// handleNetworkError handles network-related errors during login
func handleNetworkError(err error, logger *utils.Logger) {
	fmt.Println("\n🌐 Network error occurred during login.")
	fmt.Println("💡 Please check your internet connection and try again.")
	logger.Error("Network error during login", err, map[string]interface{}{
		"operation": "login",
	})
}

// handleGitHubError handles GitHub-specific errors during OAuth flow
func handleGitHubError(err error, logger *utils.Logger) {
	fmt.Println("\n🐙 GitHub authentication error occurred.")
	fmt.Println("💡 Please check your GitHub account and try again.")
	fmt.Println("📧 If the problem persists, contact support.")
	logger.Error("GitHub authentication error", err, map[string]interface{}{
		"operation": "login",
		"provider":  "github",
	})
}

// handleDeviceCodeError handles device code flow errors
func handleDeviceCodeError(err error, logger *utils.Logger) {
	fmt.Println("\n📱 Device code authentication failed.")
	fmt.Println("💡 Please try the login process again.")
	logger.Error("Device code authentication failed", err, map[string]interface{}{
		"operation": "login",
	})
}

// handleTokenStorageError handles token storage errors
func handleTokenStorageError(err error, logger *utils.Logger) {
	fmt.Println("\n💾 Failed to save authentication token.")
	fmt.Println("💡 Please check your file permissions and try again.")
	logger.Error("Token storage failed", err, map[string]interface{}{
		"operation": "login",
	})
}

// handleAlreadyLoggedInError handles when user is already logged in
func handleAlreadyLoggedInError(logger *utils.Logger) {
	fmt.Println("\n✅ You are already logged in to Kavach.")
	fmt.Println("💡 Use 'kavach logout' if you want to log out and log in with a different account.")
	logger.Info("User already logged in", map[string]interface{}{
		"operation": "login",
	})
}

// handleLoginTimeoutError handles login timeout errors
func handleLoginTimeoutError(logger *utils.Logger) {
	fmt.Println("\n⏰ Login process timed out.")
	fmt.Println("💡 Please try logging in again.")
	logger.Warn("Login timeout", map[string]interface{}{
		"operation": "login",
	})
}
