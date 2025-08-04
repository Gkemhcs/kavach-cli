package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Gkemhcs/kavach-cli/internal/config"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// StatusTokenData represents the non-sensitive fields of the user's token data for status display.
type StatusTokenData struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Organization string `json:"organization,omitempty"`
	Environment  string `json:"environment,omitempty"`
}

// statusCmd is the Cobra command for displaying current login status.
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current login status (without tokens)",
	Long:  "Display the current login status, showing all fields except sensitive tokens.",
	Run: func(cmd *cobra.Command, args []string) {
		logger := cmd.Context().Value(CtxLoggerKey).(*utils.Logger)
		defer logger.Close()
		// Get home directory and config
		homeDir, _ := os.UserHomeDir()
		cfg := cmd.Context().Value(CtxCfgKey).(*config.Config)
		path := filepath.Join(homeDir, cfg.TokenFilePath)
		file, err := os.Open(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "❌ You are not logged in. Please run `kavach login` to log in.")
			logger.Warn("User not logged in when running status command", map[string]interface{}{"cmd": "status", "path": path})
			return
		}
		defer file.Close()
		var raw map[string]interface{}
		if err := json.NewDecoder(file).Decode(&raw); err != nil {
			fmt.Fprintln(os.Stderr, "❌ Failed to parse credentials file. Please try logging in again.")
			logger.Error("Failed to parse credentials file in status command", err, map[string]interface{}{"cmd": "status", "path": path})
			return
		}
		// Remove sensitive fields
		delete(raw, "access_token")
		delete(raw, "refresh_token")
		output, _ := json.MarshalIndent(raw, "", "  ")
		fmt.Println("🔑 Current Login Status:")
		fmt.Println(string(output))
		logger.Info("Displayed current login status", map[string]interface{}{"cmd": "status"})
	},
}

// init registers the status command with the root command.
func init() {
	rootCmd.AddCommand(statusCmd)
}
