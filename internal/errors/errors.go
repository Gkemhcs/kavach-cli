package errors

import (
	"errors"
	"strings"
)

// Predefined error variables for common CLI error scenarios.
var ErrNotLoggedIn = errors.New("🚫 you are not logged in. Please run 'kavach login'")
var ErrTokenExpired = errors.New("🔑 token expired or invalid. Please log in again")
var ErrInvalidToken = errors.New("❌ internal error. Please login again")
var ErrOrganizationNotFound = errors.New("❌ sorry, the organization was not found")
var ErrSecretGroupNotFound = errors.New("❌ sorry, the secret group was not found")
var ErrEnvironmentNotFound = errors.New("❌ sorry, the environment was not found")
var ErrUserGroupNotFound = errors.New("❌ sorry, the usergroup was not found")
var ErrUserNotFound = errors.New("❌ sorry, the user was not found")
var ErrMemberNotFound = errors.New("❌ sorry, the member was not found")
var ErrRoleBindingNotFound = errors.New("❌ Sorry, the role binding was not found")
var ErrDuplicateOrganisation = errors.New("⚠️ The organization already exists")
var ErrDuplicateSecretGroup = errors.New("⚠️ The secret group already exists")
var ErrDuplicateEnvironment = errors.New("⚠️ The environment already exists")
var ErrDuplicateUserGroup = errors.New("⚠️ The usergroup already exists")
var ErrDuplicateMember = errors.New("⚠️ The member already exists")
var ErrDuplicateRoleBinding = errors.New("⚠️ The role binding already exists")
var ErrConfigFileEmpty = errors.New("⚠️ The config file is empty")
var ErrConfigFileNotFound = errors.New("⚠️ The config file is missing. Please activate an organization or group")
var ErrUnReachableBackend = errors.New("🚨 The backend is not reachable at this time")
var ErrConnectionFailed = errors.New("🌐 Connection to the server failed. Please check if the server is running and try again")
var ErrForeignKeyViolation = errors.New("⚠️ foreign key violation")
var ErrInternalServer = errors.New("🚨 Internal server error. Please try again later")
var ErrSecretVersionNotFound = errors.New("❌ Sorry, the secret version was not found")
var ErrSecretNotFound = errors.New("❌ Sorry, the secret was not found")
var ErrTargetSecretVersionNotFound = errors.New("❌ Sorry, the target secret version was not found")
var ErrEnvironmentsMisMatch = errors.New("❌ Sorry, the target secret version environment is different")
var ErrEncryptionFailed = errors.New("❌ Failed to encrypt secret value")
var ErrDecryptionFailed = errors.New("❌ Failed to decrypt secret value")
var ErrRollbackFailed = errors.New("❌ Failed to rollback to specified version")
var ErrCopySecretsFailed = errors.New("❌ Failed to copy the secrets from previous version to rollback")
var ErrInvalidSecretName = errors.New("❌ Secret name contains invalid characters or is empty")
var ErrSecretValueTooLong = errors.New("❌ Secret value exceeds maximum allowed length")
var ErrEmptySecrets = errors.New("❌ The secrets are empty")
var ErrTooManySecrets = errors.New("❌ Number of secrets exceeds maximum allowed limit")
var ErrAccessDenied = errors.New("🚫 You don't have access to perform this action")

// Provider-specific errors
var ErrProviderCredentialNotFound = errors.New("❌ Provider credential not found")
var ErrProviderCredentialExists = errors.New("⚠️ Provider credential already exists")
var ErrInvalidProviderType = errors.New("❌ Invalid provider type. Supported: github, gcp, azure")
var ErrInvalidProviderData = errors.New("❌ Invalid provider data. Please check your credentials and configuration")
var ErrProviderEncryptionFailed = errors.New("❌ Failed to encrypt provider credentials")
var ErrProviderCredentialCreateFailed = errors.New("❌ Failed to create provider credential")
var ErrProviderCredentialGetFailed = errors.New("❌ Failed to retrieve provider credential")
var ErrProviderCredentialUpdateFailed = errors.New("❌ Failed to update provider credential")
var ErrProviderCredentialDeleteFailed = errors.New("❌ Failed to delete provider credential")
var ErrProviderCredentialListFailed = errors.New("❌ Failed to list provider credentials")
var ErrProviderSyncFailed = errors.New("❌ Failed to sync with provider")
var ErrProviderCredentialValidationFailed = errors.New("❌ Provider credential validation failed. Please check your credentials")
var ErrNoSecretsToSync = errors.New("❌ No secrets found to sync. Please ensure secrets exist in the environment")
var ErrGitHubEnvironmentNotFound = errors.New("❌ GitHub environment specified in config was not found in the repository")
var ErrGitHubEncryptionFailed = errors.New("❌ Failed to encrypt secret for GitHub. Please try again")
var ErrGCPInvalidLocation = errors.New("❌ GCP Secret Manager location specified in config is invalid or not supported")

// IsConnectionError returns true if the error message indicates a network connection error.
func IsConnectionError(errorMsg string) bool {
	if strings.Contains(errorMsg, "connection refused") ||
		strings.Contains(errorMsg, "no such host") ||
		strings.Contains(errorMsg, "dial tcp") ||
		strings.Contains(errorMsg, "EOF") ||
		strings.Contains(errorMsg, "connection reset") ||
		strings.Contains(errorMsg, "broken pipe") ||
		strings.Contains(errorMsg, "network is unreachable") {
		return true
	}
	return false
}
