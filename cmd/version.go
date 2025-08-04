package cmd

import (
	"fmt"
	"time"

	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/Gkemhcs/kavach-cli/internal/version"
	"github.com/spf13/cobra"
)

// versionCmd is the Cobra command for displaying version information.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Long: `Display detailed version information for the Kavach CLI.

This command shows:
  • Semantic version number
  • Build timestamp
  • Git commit hash and branch
  • Go runtime version
  • Target platform (OS/Architecture)

Examples:
  kavach version
  kavach version --json
  kavach version --short`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := cmd.Context().Value(CtxLoggerKey).(*utils.Logger)
		defer logger.Close()
		short, _ := cmd.Flags().GetBool("short")
		json, _ := cmd.Flags().GetBool("json")
		if short {
			fmt.Println("🔖 Version:", version.GetVersion())
			logger.Info("Displayed short version info", map[string]interface{}{"cmd": "version", "short": true})
			return
		}
		if json {
			jsonOutput := fmt.Sprintf(`{
  "version": "%s",
  "versionType": "%s",
  "buildTime": "%s",
  "gitCommit": "%s",
  "gitBranch": "%s",
  "goVersion": "%s",
  "platform": "%s"
}`, version.GetVersion(), version.GetVersionType(), version.GetBuildTime().Format(time.RFC3339), version.GetGitCommit(), version.GetGitBranch(), version.GetGoVersion(), version.GetPlatform())
			fmt.Println(jsonOutput)
			logger.Info("Displayed version info in JSON format", map[string]interface{}{"cmd": "version", "json": true})
			return
		}
		fmt.Println(version.BuildInfo())
		logger.Info("Displayed full version info", map[string]interface{}{"cmd": "version", "short": false, "json": false})
	},
}

// init registers the version command with the root command.
func init() {
	rootCmd.AddCommand(versionCmd)
	// Add flags for different output formats
	versionCmd.Flags().BoolP("short", "s", false, "Display only the version number")
	versionCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
	// Mark flags as mutually exclusive
	versionCmd.MarkFlagsMutuallyExclusive("short", "json")
}
