package logout

import (
	"github.com/Gkemhcs/kavach-cli/internal/auth"
	"github.com/Gkemhcs/kavach-cli/internal/config"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// NewLogoutCommand creates the logout CLI command.
func NewLogoutCommand(logger *utils.Logger, cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "🚪 Log out and remove local credentials",
		Long: `Log out from Kavach and remove your locally stored credentials.

This command securely removes your authentication token and other credentials
stored locally on your machine. After running this command, you'll need to
run 'kavach login' again to authenticate for future CLI operations.

The logout process:
1. Removes the stored access token from your home directory
2. Clears any cached authentication data
3. Ensures you're completely logged out from the CLI

This is useful when:
- Switching between different GitHub accounts
- Securing your machine before sharing it
- Troubleshooting authentication issues
- Ensuring complete logout for security purposes

Examples:
  kavach logout                    # Log out and remove credentials
  kavach logout --help            # Show detailed help

Note: This only affects your local CLI credentials. It doesn't revoke
your GitHub OAuth application access.`,
		Example: `  kavach logout
  # Log out and remove all local credentials`,
		RunE: func(cmd *cobra.Command, args []string) error {
			defer logger.Close()
			// Call the user-facing logout logic
			return auth.LogoutUser(logger, cfg)
		},
	}
}
