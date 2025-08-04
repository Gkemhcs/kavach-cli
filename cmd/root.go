package cmd

import (
	"context"

	"github.com/Gkemhcs/kavach-cli/cmd/environment"
	"github.com/Gkemhcs/kavach-cli/cmd/login"
	"github.com/Gkemhcs/kavach-cli/cmd/logout"
	"github.com/Gkemhcs/kavach-cli/cmd/org"
	"github.com/Gkemhcs/kavach-cli/cmd/provider"
	"github.com/Gkemhcs/kavach-cli/cmd/secret"
	"github.com/Gkemhcs/kavach-cli/cmd/secretgroup"
	usergroup "github.com/Gkemhcs/kavach-cli/cmd/user-group"
	"github.com/Gkemhcs/kavach-cli/internal/config"
	envClient "github.com/Gkemhcs/kavach-cli/internal/environment"
	"github.com/Gkemhcs/kavach-cli/internal/groups"
	orgClient "github.com/Gkemhcs/kavach-cli/internal/org"
	providerClient "github.com/Gkemhcs/kavach-cli/internal/provider"
	secretClient "github.com/Gkemhcs/kavach-cli/internal/secret"
	groupClient "github.com/Gkemhcs/kavach-cli/internal/secretgroup"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
	"github.com/spf13/cobra"
)

// Logger is the global logger instance for the CLI.
var Logger *utils.Logger

// rootCmd is the base command for the CLI.
var rootCmd = &cobra.Command{
	Use:   "kavach",
	Short: "Enterprise-grade secret management and synchronization CLI",
	Long: `Kavach - Enterprise Secret Management Platform

A comprehensive CLI tool for secure secret management, synchronization, and distribution
across cloud providers and environments. Kavach provides enterprise-grade security
features including encryption, access control, and seamless integration
with major cloud platforms.

Features:
  • Secure secret storage and encryption
  • Multi-cloud provider synchronization
  • Role-based access control (RBAC)
  • Organization and group management
  • Environment-specific secret deployment

Examples:
  kavach secret create --name api-key --value "secret123"
  kavach secret sync --provider aws --region us-west-2
  kavach org list
  kavach group create --name developers

For more information, visit: https://github.com/Gkemhcs/kavach-backend`,
}

type contextKey string

const (
	// CtxLoggerKey is the context key for the logger instance
	CtxLoggerKey contextKey = "logger"
	// CtxCfgKey is the context key for the config instance
	CtxCfgKey contextKey = "cfg"
)

// Execute sets up the CLI context, registers all commands, and runs the root command.
// It ensures logger and config are available in all subcommands via context.
func Execute(logger *utils.Logger, cfg *config.Config) error {
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Set config and logger in context for all commands
		ctx := context.WithValue(cmd.Context(), CtxCfgKey, cfg)
		ctx = context.WithValue(ctx, CtxLoggerKey, logger)
		cmd.SetContext(ctx)
		logger.Debug("PersistentPreRunE: context set for command", map[string]interface{}{"cmd": cmd.Name()})
		return nil
	}
	// Initialize clients with logger and config
	orgClient := orgClient.NewOrgHttpClient(cfg, logger)
	groupClient := groupClient.NewSecretGroupHttpClient(cfg, logger, orgClient)
	envClient := envClient.NewEnvironmentHttpClient(logger, cfg, groupClient)
	userGroupClient := groups.NewUsergroupHttpClient(orgClient, logger, cfg)
	secretClient := secretClient.NewSecretHttpClient(cfg, logger, orgClient, groupClient, envClient)
	providerClient := providerClient.NewProviderHttpClient(cfg, logger, orgClient, groupClient, envClient)

	// Register all domain commands
	rootCmd.AddCommand(
		login.NewLoginCommand(logger, cfg),
		logout.NewLogoutCommand(logger, cfg),
		org.NewOrgCommand(logger, cfg, orgClient),
		secretgroup.NewSecretGroupCommand(logger, cfg, groupClient),
		environment.NewEnvironmentCommand(logger, cfg, envClient),
		usergroup.NewUserGroupCommand(logger, cfg, userGroupClient, userGroupClient),
		secret.NewSecretCommand(logger, cfg, secretClient),
		provider.NewProviderCommand(logger, cfg, providerClient),
	)
	logger.Info("CLI initialized and commands registered", map[string]interface{}{"cmd": "root"})
	return rootCmd.Execute()
}
