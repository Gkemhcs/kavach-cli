package provider

import "github.com/Gkemhcs/kavach-cli/internal/types"

// ProviderClient defines the interface for provider operations
type ProviderClient interface {
	CreateProviderCredential(orgName, groupName, envName, providerName string, credentials, config map[string]interface{}) (*types.ProviderCredentialResponse, error)
	GetProviderCredential(orgName, groupName, envName, providerName string) (*types.ProviderCredentialResponse, error)
	ListProviderCredentials(orgName, groupName, envName string) ([]types.ProviderCredentialResponse, error)
	UpdateProviderCredential(orgName, groupName, envName, providerName string, credentials, config map[string]interface{}) (*types.ProviderCredentialResponse, error)
	DeleteProviderCredential(orgName, groupName, envName, providerName string) error
}

// ProviderType represents supported secret sync providers
type ProviderType string

const (
	ProviderGitHub ProviderType = "github"
	ProviderGCP    ProviderType = "gcp"
	ProviderAzure  ProviderType = "azure"
)

// CreateProviderCredentialRequest represents the request to create a new provider credential
type CreateProviderCredentialRequest struct {
	Provider    ProviderType           `json:"provider" binding:"required"`
	Credentials map[string]interface{} `json:"credentials" binding:"required"`
	Config      map[string]interface{} `json:"config" binding:"required"`
}

// UpdateProviderCredentialRequest represents the request to update an existing provider credential
type UpdateProviderCredentialRequest struct {
	Credentials map[string]interface{} `json:"credentials" binding:"required"`
	Config      map[string]interface{} `json:"config" binding:"required"`
}

// GitHubCredentials represents GitHub Personal Access Token credentials
type GitHubCredentials struct {
	Token string `json:"token" binding:"required"`
}

// GitHubConfig represents GitHub-specific configuration
type GitHubConfig struct {
	Owner            string `json:"owner" binding:"required"`
	Repository       string `json:"repository" binding:"required"`
	Environment      string `json:"environment,omitempty"`
	SecretVisibility string `json:"secret_visibility,omitempty"` // all, selected, private
}

// GCPCredentials represents GCP Service Account key credentials
type GCPCredentials struct {
	Type                    string `json:"type" binding:"required"`
	ProjectID               string `json:"project_id" binding:"required"`
	PrivateKeyID            string `json:"private_key_id" binding:"required"`
	PrivateKey              string `json:"private_key" binding:"required"`
	ClientEmail             string `json:"client_email" binding:"required"`
	ClientID                string `json:"client_id,omitempty"`
	AuthURI                 string `json:"auth_uri,omitempty"`
	TokenURI                string `json:"token_uri,omitempty"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url,omitempty"`
	ClientX509CertURL       string `json:"client_x509_cert_url,omitempty"`
}

// GCPConfig represents GCP-specific configuration
type GCPConfig struct {
	ProjectID             string `json:"project_id" binding:"required"`
	SecretManagerLocation string `json:"secret_manager_location,omitempty"`
	Prefix                string `json:"prefix,omitempty"`
	Replication           string `json:"replication,omitempty"`
	DisableOlderVersions  bool   `json:"disable_older_versions,omitempty"`
}

// AzureCredentials represents Azure Service Principal credentials
type AzureCredentials struct {
	TenantID     string `json:"tenant_id" binding:"required"`
	ClientID     string `json:"client_id" binding:"required"`
	ClientSecret string `json:"client_secret" binding:"required"`
}

// AzureConfig represents Azure-specific configuration
type AzureConfig struct {
	SubscriptionID string `json:"subscription_id" binding:"required"`
	ResourceGroup  string `json:"resource_group" binding:"required"`
	KeyVaultName   string `json:"key_vault_name" binding:"required"`
	Prefix         string `json:"prefix,omitempty"`
}
