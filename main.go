package main

import (
	"os"

	"github.com/Gkemhcs/kavach-cli/cmd"
	"github.com/Gkemhcs/kavach-cli/internal/config"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
)

// main is the entry point for the Kavach CLI application.
// It loads configuration, initializes the logger, and executes the root command.
func main() {
	cfg := config.Load()
	logger := utils.NewLogger(cfg)
	// Execute the root command and handle any errors
	if err := cmd.Execute(logger, cfg); err != nil {
		logger.Error("Failed to execute command", err, map[string]interface{}{"exit_code": 1})
		os.Exit(1)
	}
}
