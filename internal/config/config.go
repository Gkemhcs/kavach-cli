package config

import (
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

// Config holds all configuration values for the CLI, loaded from environment variables or config files.
type Config struct {
	DeviceCodeURL   string // Device code endpoint
	DeviceTokenURL  string // Device token endpoint
	TokenFilePath   string // Path to credentials file
	BackendEndpoint string // URL of backend
	LogDirPath      string // Log directory path
}

// Load reads configuration from environment variables and returns a Config struct.
// Sets sensible defaults for local development and production.
func Load() *Config {
	viper.SetDefault("KAVACH_DEVICE_CODE_URL", "http://localhost:8080/api/v1/auth/device/code")
	viper.SetDefault("KAVACH_DEVICE_TOKEN_URL", "http://localhost:8080/api/v1/auth/device/token")
	viper.SetDefault("KAVACH_TOKEN_FILE_PATH", "/.kavach/credentials.json")
	viper.SetDefault("KAVACH_LOG_DIR_PATH", "/.kavach/")
	viper.SetDefault("KAVACH_BACKEND_ENDPOINT", "http://localhost:8080/api/v1/")

	// viper.SetDefault("KAVACH_DEVICE_CODE_URL", "https://kavach.gkem.cloud/api/v1/auth/device/code")
	// viper.SetDefault("KAVACH_DEVICE_TOKEN_URL", "https://kavach.gkem.cloud/api/v1/auth/device/token")
	// viper.SetDefault("KAVACH_TOKEN_FILE_PATH", "/.kavach/credentials.json")
	// viper.SetDefault("KAVACH_LOG_DIR_PATH", "/.kavach/")
	// viper.SetDefault("KAVACH_BACKEND_ENDPOINT", "https://kavach.gkem.cloud/api/v1/")

	return &Config{
		DeviceCodeURL:   viper.GetString("KAVACH_DEVICE_CODE_URL"),
		DeviceTokenURL:  viper.GetString("KAVACH_DEVICE_TOKEN_URL"),
		TokenFilePath:   viper.GetString("KAVACH_TOKEN_FILE_PATH"),
		LogDirPath:      viper.GetString("KAVACH_LOG_DIR_PATH"),
		BackendEndpoint: viper.GetString("KAVACH_BACKEND_ENDPOINT"),
	}
}

// CLIConfig holds user-specific CLI settings (org, secretgroup, environment)
type CLIConfig struct {
	Organization string `yaml:"organization,omitempty"`
	SecretGroup  string `yaml:"secretgroup,omitempty"`
	Environment  string `yaml:"environment,omitempty"`
}

// ConfigFilePath returns the path to the config.yaml file in the user's home directory.
func ConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".kavach", "config.yaml"), nil
}

// LoadCLIConfig loads CLI config from ~/.kavach/config.yaml.
// Returns an empty config if the file does not exist or is empty.
func LoadCLIConfig() (*CLIConfig, error) {
	path, err := ConfigFilePath()
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
	path, err := ConfigFilePath()
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
