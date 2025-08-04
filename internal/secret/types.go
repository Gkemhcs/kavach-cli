package secret

import "github.com/Gkemhcs/kavach-cli/internal/types"

// SecretClient defines the interface for secret operations
type SecretClient interface {
	CreateVersion(orgName, groupName, envName string, secrets []types.Secret, commitMessage string) (*types.SecretVersionResponse, error)
	ListVersions(orgName, groupName, envName string) ([]types.SecretVersionResponse, error)
	GetVersionDetails(orgName, groupName, envName, versionID string) (*types.SecretVersionDetailResponse, error)
	RollbackToVersion(orgName, groupName, envName, versionID string, commitMessage string) (*types.SecretVersionResponse, error)
	GetVersionDiff(orgName, groupName, envName, fromVersion, toVersion string) (*types.SecretDiffResponse, error)
	SyncSecrets(orgName, groupName, envName, provider, versionID string) (*types.SyncSecretsResponse, error)
}

// StagedSecrets represents the staging area data
type StagedSecrets struct {
	Secrets []types.Secret `json:"secrets"`
}
