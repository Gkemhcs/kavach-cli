package secret

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Gkemhcs/kavach-cli/internal/types"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
)

// StagingService manages the staging area for secrets
type StagingService struct {
	logger *utils.Logger
}

// NewStagingService creates a new staging service
func NewStagingService(logger *utils.Logger) *StagingService {
	return &StagingService{
		logger: logger,
	}
}

// getStagingFilePath returns the path to the staging file
func (s *StagingService) getStagingFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	stagingDir := filepath.Join(homeDir, ".kavach", "staging")
	if err := os.MkdirAll(stagingDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create staging directory: %w", err)
	}

	return filepath.Join(stagingDir, "secrets.json"), nil
}

// LoadStagedSecrets loads secrets from the staging area
func (s *StagingService) LoadStagedSecrets() (*StagedSecrets, error) {
	filePath, err := s.getStagingFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &StagedSecrets{
				Secrets: []types.Secret{},
			}, nil
		}
		return nil, fmt.Errorf("failed to read staging file: %w", err)
	}

	// Handle empty file case
	if len(data) == 0 {
		return &StagedSecrets{
			Secrets: []types.Secret{},
		}, nil
	}

	var staged StagedSecrets
	if err := json.Unmarshal(data, &staged); err != nil {
		return nil, fmt.Errorf("failed to unmarshal staging data: %w", err)
	}

	return &staged, nil
}

// SaveStagedSecrets saves secrets to the staging area
func (s *StagingService) SaveStagedSecrets(staged *StagedSecrets) error {
	filePath, err := s.getStagingFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(staged, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal staging data: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write staging file: %w", err)
	}

	return nil
}

// AddSecret adds a secret to the staging area
func (s *StagingService) AddSecret(orgName, groupName, envName, name, value string) error {
	staged, err := s.LoadStagedSecrets()
	if err != nil {
		return err
	}

	// Check if secret already exists and update it
	for i, secret := range staged.Secrets {
		if secret.Name == name {
			staged.Secrets[i].Value = value
			s.logger.Info("Updated existing secret in staging", map[string]interface{}{
				"secret_name": name,
			})
			return s.SaveStagedSecrets(staged)
		}
	}

	// Add new secret
	staged.Secrets = append(staged.Secrets, types.Secret{
		Name:  name,
		Value: value,
	})

	s.logger.Info("Added secret to staging", map[string]interface{}{
		"secret_name": name,
	})

	return s.SaveStagedSecrets(staged)
}

// AddSecretToStaging adds a secret to the staging area without context
// This method is used by the secret add command to store secrets temporarily
func (s *StagingService) AddSecretToStaging(name, value string) error {
	staged, err := s.LoadStagedSecrets()
	if err != nil {
		return err
	}

	// Check if secret already exists and update it
	for i, secret := range staged.Secrets {
		if secret.Name == name {
			staged.Secrets[i].Value = value
			s.logger.Info("Updated existing secret in staging", map[string]interface{}{
				"secret_name": name,
			})
			return s.SaveStagedSecrets(staged)
		}
	}

	// Add new secret
	staged.Secrets = append(staged.Secrets, types.Secret{
		Name:  name,
		Value: value,
	})

	s.logger.Info("Added secret to staging", map[string]interface{}{
		"secret_name": name,
	})

	return s.SaveStagedSecrets(staged)
}

// GetStagedSecrets returns the current staged secrets
func (s *StagingService) GetStagedSecrets() (*StagedSecrets, error) {
	return s.LoadStagedSecrets()
}

// ClearStaging clears the staging area
func (s *StagingService) ClearStaging() error {
	filePath, err := s.getStagingFilePath()
	if err != nil {
		return err
	}

	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove staging file: %w", err)
	}

	s.logger.Info("Cleared staging area")
	return nil
}
