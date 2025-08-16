package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// Config holds all configuration values for the CLI, loaded from environment variables or config files.
type Config struct {
	DeviceCodeURL   string // Device code endpoint
	DeviceTokenURL  string // Device token endpoint
	TokenFilePath   string // Path to credentials file
	BackendEndpoint string // URL of backend
	LogDirPath      string // Log directory path

	// Logging configuration
	LogMaxSize    int  // Maximum size of log file in MB
	LogMaxBackups int  // Maximum number of old log files to keep
	LogMaxAge     int  // Maximum age of log files in days
	LogCompress   bool // Whether to compress old log files
}

// loadEnvFiles attempts to load .env files from multiple locations
// Returns true if any .env file was loaded successfully
func loadEnvFiles() bool {
	// List of potential .env file locations (in order of priority)
	envLocations := []string{
		".env",                                   // Current directory
		".env.local",                             // Local overrides
		".env.development",                       // Development environment
		".env.production",                        // Production environment
		filepath.Join(os.Getenv("HOME"), ".env"), // User home directory
		filepath.Join(os.Getenv("HOME"), ".kavach", ".env"), // Kavach config directory
	}

	loaded := false
	for _, location := range envLocations {
		if location == "" {
			continue
		}

		// Check if file exists before trying to load
		if _, err := os.Stat(location); err == nil {
			if err := godotenv.Load(location); err == nil {
				loaded = true
			}
			// Silently ignore errors and continue
		}
	}

	return loaded
}

// Load reads configuration from environment variables and returns a Config struct.
// Sets sensible defaults for local development and production.
// Never panics - gracefully handles missing .env files and uses defaults.
func Load() *Config {
	// Try to load .env files (optional - won't panic if not found)
	loadEnvFiles()

	// Set defaults for all configuration values
	// These will be used if no environment variables are set

	// Production defaults (secure by default)
	viper.SetDefault("KAVACH_DEVICE_CODE_URL", "https://kavach.gkem.cloud/api/v1/auth/device/code")
	viper.SetDefault("KAVACH_DEVICE_TOKEN_URL", "https://kavach.gkem.cloud/api/v1/auth/device/token")
	viper.SetDefault("KAVACH_BACKEND_ENDPOINT", "https://kavach.gkem.cloud/api/v1/")

	// File path defaults
	viper.SetDefault("KAVACH_TOKEN_FILE_PATH", "/.kavach/credentials.json")
	viper.SetDefault("KAVACH_LOG_DIR_PATH", "/.kavach/")

	// Logging configuration defaults
	viper.SetDefault("KAVACH_LOG_MAX_SIZE", 1)     // 1 MB
	viper.SetDefault("KAVACH_LOG_MAX_BACKUPS", 3)  // 3 backup files
	viper.SetDefault("KAVACH_LOG_MAX_AGE", 28)     // 28 days
	viper.SetDefault("KAVACH_LOG_COMPRESS", false) // Don't compress by default

	// Also set defaults for non-prefixed versions (backward compatibility)
	viper.SetDefault("DEVICE_CODE_URL", viper.GetString("KAVACH_DEVICE_CODE_URL"))
	viper.SetDefault("DEVICE_TOKEN_URL", viper.GetString("KAVACH_DEVICE_TOKEN_URL"))
	viper.SetDefault("BACKEND_ENDPOINT", viper.GetString("KAVACH_BACKEND_ENDPOINT"))
	viper.SetDefault("TOKEN_FILE_PATH", viper.GetString("KAVACH_TOKEN_FILE_PATH"))
	viper.SetDefault("LOG_DIR_PATH", viper.GetString("KAVACH_LOG_DIR_PATH"))

	// Create config with fallback logic
	config := &Config{
		DeviceCodeURL:   getStringWithFallback("KAVACH_DEVICE_CODE_URL", "DEVICE_CODE_URL"),
		DeviceTokenURL:  getStringWithFallback("KAVACH_DEVICE_TOKEN_URL", "DEVICE_TOKEN_URL"),
		TokenFilePath:   getStringWithFallback("KAVACH_TOKEN_FILE_PATH", "TOKEN_FILE_PATH"),
		LogDirPath:      getStringWithFallback("KAVACH_LOG_DIR_PATH", "LOG_DIR_PATH"),
		BackendEndpoint: getStringWithFallback("KAVACH_BACKEND_ENDPOINT", "BACKEND_ENDPOINT"),

		// Logging configuration
		LogMaxSize:    viper.GetInt("KAVACH_LOG_MAX_SIZE"),
		LogMaxBackups: viper.GetInt("KAVACH_LOG_MAX_BACKUPS"),
		LogMaxAge:     viper.GetInt("KAVACH_LOG_MAX_AGE"),
		LogCompress:   viper.GetBool("KAVACH_LOG_COMPRESS"),
	}

	// Configuration loaded silently (removed verbose output)

	return config
}

// getStringWithFallback gets a string value with fallback support
// Tries KAVACH_* first, then falls back to non-prefixed version
func getStringWithFallback(primaryKey, fallbackKey string) string {
	value := viper.GetString(primaryKey)
	if value == "" {
		value = viper.GetString(fallbackKey)
	}
	return value
}

// GetEnvFilePath returns the path to the .env file in the current directory
func GetEnvFilePath() string {
	return ".env"
}

// CreateDefaultEnvFile creates a default .env file if it doesn't exist
func CreateDefaultEnvFile() error {
	envPath := GetEnvFilePath()

	// Check if .env already exists
	if _, err := os.Stat(envPath); err == nil {
		return nil // .env already exists
	}

	// Create default .env content
	defaultContent := `# Kavach CLI Configuration
# This file contains environment-specific configuration
# All values are optional and have sensible defaults

# Authentication endpoints
KAVACH_DEVICE_CODE_URL=https://kavach.gkem.cloud/api/v1/auth/device/code
KAVACH_DEVICE_TOKEN_URL=https://kavach.gkem.cloud/api/v1/auth/device/token

# Backend configuration
KAVACH_BACKEND_ENDPOINT=https://kavach.gkem.cloud/api/v1/

# File paths
KAVACH_TOKEN_FILE_PATH=/.kavach/credentials.json
KAVACH_LOG_DIR_PATH=/.kavach/

# Logging configuration
KAVACH_LOG_MAX_SIZE=1
KAVACH_LOG_MAX_BACKUPS=3
KAVACH_LOG_MAX_AGE=28
KAVACH_LOG_COMPRESS=false

# Development overrides (uncomment to use)
# KAVACH_DEVICE_CODE_URL=http://localhost:8080/api/v1/auth/device/code
# KAVACH_DEVICE_TOKEN_URL=http://localhost:8080/api/v1/auth/device/token
# KAVACH_BACKEND_ENDPOINT=http://localhost:8080/api/v1/
`

	// Write the file
	if err := os.WriteFile(envPath, []byte(defaultContent), 0644); err != nil {
		return fmt.Errorf("failed to create .env file: %w", err)
	}

	fmt.Printf("✅ Created default .env file: %s\n", envPath)
	return nil
}

// ValidateConfig validates the configuration and returns any errors
func (c *Config) ValidateConfig() []string {
	var errors []string

	// Check required fields
	if c.BackendEndpoint == "" {
		errors = append(errors, "Backend endpoint is required")
	}
	if c.DeviceCodeURL == "" {
		errors = append(errors, "Device code URL is required")
	}
	if c.DeviceTokenURL == "" {
		errors = append(errors, "Device token URL is required")
	}

	// Validate URLs
	if c.BackendEndpoint != "" && !strings.HasPrefix(c.BackendEndpoint, "http") {
		errors = append(errors, "Backend endpoint must be a valid HTTP/HTTPS URL")
	}
	if c.DeviceCodeURL != "" && !strings.HasPrefix(c.DeviceCodeURL, "http") {
		errors = append(errors, "Device code URL must be a valid HTTP/HTTPS URL")
	}
	if c.DeviceTokenURL != "" && !strings.HasPrefix(c.DeviceTokenURL, "http") {
		errors = append(errors, "Device token URL must be a valid HTTP/HTTPS URL")
	}

	// Validate logging configuration
	if c.LogMaxSize < 1 {
		errors = append(errors, "Log max size must be at least 1 MB")
	}
	if c.LogMaxBackups < 0 {
		errors = append(errors, "Log max backups cannot be negative")
	}
	if c.LogMaxAge < 0 {
		errors = append(errors, "Log max age cannot be negative")
	}

	return errors
}

// CLIConfig holds user-specific CLI settings (org, secretgroup, environment)
type CLIConfig struct {
	Organization string `yaml:"organization,omitempty"`
	SecretGroup  string `yaml:"secretgroup,omitempty"`
	Environment  string `yaml:"environment,omitempty"`
}

// FilePath returns the path to the config.yaml file in the user's home directory.
func FilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".kavach", "config.yaml"), nil
}

// LoadCLIConfig loads CLI config from ~/.kavach/config.yaml.
// Returns an empty config if the file does not exist or is empty.
func LoadCLIConfig() (*CLIConfig, error) {
	path, err := FilePath()
	if err != nil {
		return &CLIConfig{}, nil // just return empty config
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &CLIConfig{}, nil // config doesn't exist, return default
		}
		return nil, err
	}
	if len(data) == 0 {
		return &CLIConfig{}, nil // empty file, return empty config
	}
	var cfg CLIConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	// Optional: treat config with all fields empty as "empty config"
	if cfg.Organization == "" && cfg.SecretGroup == "" && cfg.Environment == "" {
		return &CLIConfig{}, nil
	}
	return &cfg, nil
}

// SaveCLIConfig saves CLI config to ~/.kavach/config.yaml.
func SaveCLIConfig(cfg *CLIConfig) error {
	path, err := FilePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := yaml.NewEncoder(file)
	defer enc.Close()
	return enc.Encode(cfg)
}
