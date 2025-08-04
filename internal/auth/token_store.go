package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
)

// SaveToken saves the user's token data to the credentials file.
// Logs errors and ensures the credentials directory exists.
func SaveToken(tokenConfig TokenData, logger *utils.Logger, cfg *config.Config) error {
	homeDir, _ := os.UserHomeDir()
	path := filepath.Join(homeDir, cfg.TokenFilePath)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		logger.Error("Failed to create credentials directory", err, map[string]interface{}{"cmd": "save-token", "path": path})
		return err
	}
	data := TokenData{
		AccessToken:  tokenConfig.AccessToken,
		RefreshToken: tokenConfig.RefreshToken,
		Name:         tokenConfig.Name,
		Email:        tokenConfig.Email,
	}
	file, err := os.Create(path)
	if err != nil {
		logger.Error("Failed to create credentials file", err, map[string]interface{}{"cmd": "save-token", "path": path})
		return err
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(&data)
}

// LoadToken loads the user's token data from the credentials file.
// Logs errors and returns a user-friendly error if the file is missing or corrupt.
func LoadToken(logger *utils.Logger, cfg *config.Config) (*TokenData, error) {
	homeDir, _ := os.UserHomeDir()
	path := filepath.Join(homeDir, cfg.TokenFilePath)
	file, err := os.Open(path)
	if err != nil {
		logger.Error("Credentials file not found", err, map[string]interface{}{"cmd": "load-token", "path": path})
		return nil, fmt.Errorf("please login first using `kavach login`")
	}
	defer file.Close()
	var data TokenData
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		logger.Error("Corrupt token file", err, map[string]interface{}{"cmd": "load-token", "path": path})
		return nil, fmt.Errorf("corrupt token file: %w", err)
	}
	return &data, nil
}

// Logout deletes the local credentials file, effectively logging out the user.
// Logs the result and handles edge cases (file missing, is directory, etc.).
func Logout(logger *utils.Logger, cfg *config.Config) error {
	homeDir, _ := os.UserHomeDir()
	path := filepath.Join(homeDir, cfg.TokenFilePath)
	logger.Debug("Attempting to delete credentials file", map[string]interface{}{"cmd": "logout", "path": path})
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "❌ You are not logged in.")
		logger.Info("You are not logged in.", map[string]interface{}{"cmd": "logout", "path": path})
		return err
	}
	if err != nil {
		logger.Error("Failed to stat credentials file during logout", err, map[string]interface{}{"cmd": "logout", "path": path})
		return err
	}
	if info.IsDir() {
		logger.Error("Logout failed: path is a directory, not a file", err, map[string]interface{}{"cmd": "logout", "path": path})
		return err
	}
	if err := os.Remove(path); err != nil {
		logger.Error("Failed to delete credentials file during logout", err, map[string]interface{}{"cmd": "logout", "path": path})
		return err
	}
	logger.Info("Successfully logged out. Credentials file deleted.", map[string]interface{}{"cmd": "logout", "path": path})
	return nil
}
