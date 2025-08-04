package provider

import (
	"encoding/json"
	"fmt"
	"os"
)

// ReadGCPServiceAccountFile reads a GCP service account JSON file and returns the credentials
func ReadGCPServiceAccountFile(filePath string) (map[string]interface{}, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read service account file: %w", err)
	}

	var credentials map[string]interface{}
	if err := json.Unmarshal(file, &credentials); err != nil {
		return nil, fmt.Errorf("failed to parse service account file: %w", err)
	}

	// Validate required fields
	requiredFields := []string{"type", "project_id", "private_key_id", "private_key", "client_email"}
	for _, field := range requiredFields {
		if _, exists := credentials[field]; !exists {
			return nil, fmt.Errorf("missing required field in service account file: %s", field)
		}
	}

	return credentials, nil
}

// ValidateProviderName validates if the provider name is supported
func ValidateProviderName(providerName string) error {
	validProviders := []string{"github", "gcp", "azure"}
	for _, provider := range validProviders {
		if provider == providerName {
			return nil
		}
	}
	return fmt.Errorf("unsupported provider: %s. Supported providers: %v", providerName, validProviders)
}
