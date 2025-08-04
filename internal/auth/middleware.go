package auth

import (
	"os"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
)

// RequireAuthMiddleware checks if the user is authenticated by verifying the credentials file exists.
// Logs an error and returns ErrNotLoggedIn if not authenticated.
func RequireAuthMiddleware(logger *utils.Logger, cfg *config.Config) error {
	homeDir, _ := os.UserHomeDir()
	path := homeDir + "/" + cfg.TokenFilePath
	if _, err := os.Stat(path); os.IsNotExist(err) {
		logger.Error("You are not logged in. Please run `kavach login` to log in.", err, map[string]interface{}{"cmd": "auth-middleware", "path": path})
		return cliErrors.ErrNotLoggedIn
	}
	return nil
}
