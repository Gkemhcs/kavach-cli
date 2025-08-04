package login

import (
	"fmt"

	"github.com/Gkemhcs/kavach-cli/internal/auth"
	"github.com/Gkemhcs/kavach-cli/internal/config"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewLoginCommand returns a Cobra command for authenticating with GitHub and saving the access token.
// Handles user login flow, error reporting, and logging.
func NewLoginCommand(logger *utils.Logger, cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "🔐 Authenticate with GitHub and save access token",
		Long: `Authenticate with GitHub using OAuth device flow and save your access token for CLI use.

This command initiates a secure OAuth authentication flow with GitHub. You'll be prompted
to complete device authentication in your browser by visiting a GitHub URL and entering
a device code.

The authentication process:
1. Generates a device code and authorization URL
2. Opens your browser to complete GitHub authorization
3. Waits for you to authorize the application
4. Saves the access token locally for future CLI operations

Your credentials are stored securely in your home directory and are used automatically
for all subsequent CLI commands that require authentication.

Examples:
  kavach login                    # Start GitHub authentication
  kavach login --help            # Show detailed help

Note: You only need to run this command once. The token will be automatically
refreshed when needed.`,
		Example: `  kavach login
  # Start GitHub OAuth authentication flow
  # Follow the prompts to complete authentication in your browser`,
		RunE: func(cmd *cobra.Command, args []string) error {
			defer logger.Close()
			fmt.Println("🔐 Starting login process...🔃")
			logger.Info("Login command started", map[string]interface{}{"cmd": "login"})

			err := auth.LoginWithGitHub(logger, cfg)
			if err != nil {
				if err == cliErrors.ErrUnReachableBackend {
					fmt.Println(err)
					fmt.Println("🚨 Oops! Kavach backend is currently down or not responding.")
					fmt.Println("📡 This may be due to server downtime or high request volume.")
					fmt.Println("📩 If this persists, please drop us a message at 👉 **gudikotieswarmani@gmail.com**")
					logger.Error("Backend unreachable during login", err, map[string]interface{}{"cmd": "login"})
					return nil // or exit gracefully
				}
				if err.Error() == "device_authorization_timeout" {
					logger.Warn("Login timed out after 2 minutes. User did not complete authentication in time.", map[string]interface{}{"cmd": "login"})
					logger.Print("Login timed out after 2 minutes. Please try again and complete authentication in time.")
					return nil
				}
				if err.Error() == "device_authorization_timeout_backend" {
					logger.Warn("Login failed due to a backend timeout during device authorization.", map[string]interface{}{"cmd": "login"})
					logger.Print("Login failed due to a backend timeout. Please login again, there were issues with the device authorization.")
					return nil
				}
				logger.Error("Login failed with unexpected error", err, map[string]interface{}{"cmd": "login"})
				return err
			}
			fmt.Println("✔️ Successfully logged in to GitHub!")
			logger.Info("User successfully logged in", map[string]interface{}{"cmd": "login"})
			return nil
		},
	}
}
