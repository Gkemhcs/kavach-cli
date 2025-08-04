package secret

import (
	"encoding/json"
	"fmt"

	"github.com/Gkemhcs/kavach-cli/internal/auth"
	"github.com/Gkemhcs/kavach-cli/internal/client"
	"github.com/Gkemhcs/kavach-cli/internal/config"
	"github.com/Gkemhcs/kavach-cli/internal/environment"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/org"
	"github.com/Gkemhcs/kavach-cli/internal/secretgroup"
	"github.com/Gkemhcs/kavach-cli/internal/types"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
)

// SecretHttpClient implements SecretClient for making HTTP requests to the backend
type SecretHttpClient struct {
	cfg               *config.Config
	logger            *utils.Logger
	orgGetter         org.OrgGetterClient
	secretGroupGetter secretgroup.SecretGroupGetter
	envGetter         environment.EnvironmentGetter
}

// NewSecretHttpClient creates a new SecretHttpClient with the given config and logger
func NewSecretHttpClient(cfg *config.Config,
	logger *utils.Logger, orgGetter org.OrgGetterClient,
	secretGroupGetter secretgroup.SecretGroupGetter,
	envGetter environment.EnvironmentGetter) *SecretHttpClient {
	return &SecretHttpClient{
		cfg:               cfg,
		logger:            logger,
		orgGetter:         orgGetter,
		secretGroupGetter: secretGroupGetter,
		envGetter:         envGetter,
	}
}

// CreateVersion creates a new version of secrets for an environment
func (c *SecretHttpClient) CreateVersion(orgName, groupName, envName string, secrets []types.Secret, commitMessage string) (*types.SecretVersionResponse, error) {
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return nil, err
	}
	org, err := c.orgGetter.GetOrganizationByName(orgName)
	if err != nil {
		return nil, err
	}
	secretGroup, err := c.secretGroupGetter.GetSecretGroupByName(orgName, groupName)
	if err != nil {
		return nil, err
	}
	env, err := c.envGetter.GetEnvironmentbyName(orgName, groupName, envName)
	if err != nil {
		return nil, err
	}

	body, _ := json.Marshal(types.CreateSecretVersionRequest{
		Secrets:       secrets,
		CommitMessage: commitMessage,
	})

	url := fmt.Sprintf("%sorganizations/%s/secret-groups/%s/environments/%s/secrets/",
		c.cfg.BackendEndpoint, org.ID, secretGroup.ID, env.ID)
	requestPayload := client.RequestPayload{
		Method: "POST",
		URL:    url,
		Body:   body,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	respBody, err := client.DoAuthenticatedRequest[types.SecretVersionResponse](requestPayload, c.logger, c.cfg)
	if err != nil {
		if err == cliErrors.ErrNotLoggedIn {
			c.logger.Print("You are not logged in. Please run `kavach login` to log in.")
			c.logger.Warn("User not logged in during secret create", map[string]interface{}{"cmd": "secret create", "env": envName})
			return nil, err
		}
		c.logger.Error("Failed to create secret version", err, map[string]interface{}{"cmd": "secret create", "env": envName})
		return nil, err
	}

	if !respBody.Success {
		switch respBody.ErrorCode {
		case "empty_secrets":
			return nil, cliErrors.ErrEmptySecrets
		case "too_many_secrets":
			return nil, cliErrors.ErrTooManySecrets
		case "invalid_secret_name":
			return nil, cliErrors.ErrInvalidSecretName
		case "secret_value_too_long":
			return nil, cliErrors.ErrSecretValueTooLong
		case "encryption_failed":
			return nil, cliErrors.ErrEncryptionFailed
		default:
			c.logger.Error("Failed to create secret version", fmt.Errorf(respBody.ErrorMsg), map[string]interface{}{"cmd": "secret create", "env": envName})
			return nil, fmt.Errorf(respBody.ErrorMsg)
		}
	}

	c.logger.Info("Secret version created successfully", map[string]interface{}{"cmd": "secret create", "env": envName, "version_id": respBody.Data.ID})
	return &respBody.Data, nil
}

// ListVersions lists all versions for an environment
func (c *SecretHttpClient) ListVersions(orgName, groupName, envName string) ([]types.SecretVersionResponse, error) {
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return nil, err
	}
	org, err := c.orgGetter.GetOrganizationByName(orgName)
	if err != nil {
		return nil, err
	}
	secretGroup, err := c.secretGroupGetter.GetSecretGroupByName(orgName, groupName)
	if err != nil {
		return nil, err
	}
	env, err := c.envGetter.GetEnvironmentbyName(orgName, groupName, envName)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%sorganizations/%s/secret-groups/%s/environments/%s/secrets/versions",
		c.cfg.BackendEndpoint, org.ID, secretGroup.ID, env.ID)
	requestPayload := client.RequestPayload{
		Method: "GET",
		URL:    url,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	respBody, err := client.DoAuthenticatedRequest[[]types.SecretVersionResponse](requestPayload, c.logger, c.cfg)
	if err != nil {
		if err == cliErrors.ErrNotLoggedIn {
			c.logger.Print("You are not logged in. Please run `kavach login` to log in.")
			c.logger.Warn("User not logged in during secret list", map[string]interface{}{"cmd": "secret list", "env": envName})
			return nil, err
		}
		c.logger.Error("Failed to list secret versions", err, map[string]interface{}{"cmd": "secret list", "env": envName})
		return nil, err
	}

	if !respBody.Success {
		switch respBody.ErrorCode {
		case "environment_not_exist":
			return nil, cliErrors.ErrEnvironmentNotFound
		case "secretgroup_not_exist":
			return nil, cliErrors.ErrSecretGroupNotFound
		case "organisation_not_exist":
			return nil, cliErrors.ErrOrganizationNotFound
		case "internal_error":
			return nil, cliErrors.ErrInternalServer
		default:
			c.logger.Error("Failed to list secret versions", fmt.Errorf(respBody.ErrorMsg), map[string]interface{}{"cmd": "secret list", "env": envName})
			return nil, fmt.Errorf(respBody.ErrorMsg)
		}
	}

	c.logger.Info("Secret versions listed successfully", map[string]interface{}{"cmd": "secret list", "env": envName, "count": len(respBody.Data)})
	return respBody.Data, nil
}

// GetVersionDetails gets detailed information about a specific version including secrets
func (c *SecretHttpClient) GetVersionDetails(orgName, groupName, envName, versionID string) (*types.SecretVersionDetailResponse, error) {
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return nil, err
	}

	org, err := c.orgGetter.GetOrganizationByName(orgName)
	if err != nil {
		return nil, err
	}
	secretGroup, err := c.secretGroupGetter.GetSecretGroupByName(orgName, groupName)
	if err != nil {
		return nil, err
	}
	env, err := c.envGetter.GetEnvironmentbyName(orgName, groupName, envName)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%sorganizations/%s/secret-groups/%s/environments/%s/secrets/versions/%s",
		c.cfg.BackendEndpoint, org.ID, secretGroup.ID, env.ID, versionID)

	requestPayload := client.RequestPayload{
		Method: "GET",
		URL:    url,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	respBody, err := client.DoAuthenticatedRequest[types.SecretVersionDetailResponse](requestPayload, c.logger, c.cfg)
	if err != nil {
		if err == cliErrors.ErrNotLoggedIn {
			c.logger.Print("You are not logged in. Please run `kavach login` to log in.")
			c.logger.Warn("User not logged in during secret details", map[string]interface{}{"cmd": "secret details", "version": versionID})
			return nil, err
		}
		c.logger.Error("Failed to get version details", err, map[string]interface{}{"cmd": "secret details", "version": versionID})
		return nil, err
	}

	if !respBody.Success {
		switch respBody.ErrorCode {
		case "secret_version_not_found":
			return nil, cliErrors.ErrSecretVersionNotFound
		case "secret_not_found":
			return nil, cliErrors.ErrSecretNotFound
		case "decryption_failed":
			return nil, cliErrors.ErrDecryptionFailed
		case "internal_error":
			return nil, cliErrors.ErrInternalServer
		default:
			c.logger.Error("Failed to get version details", fmt.Errorf(respBody.ErrorMsg), map[string]interface{}{"cmd": "secret details", "version": versionID})
			return nil, fmt.Errorf(respBody.ErrorMsg)
		}
	}

	c.logger.Info("Version details retrieved successfully", map[string]interface{}{"cmd": "secret details", "version": versionID})
	return &respBody.Data, nil
}

// RollbackToVersion creates a new version by copying secrets from a previous version
func (c *SecretHttpClient) RollbackToVersion(orgName, groupName, envName, versionID string, commitMessage string) (*types.SecretVersionResponse, error) {
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return nil, err
	}

	org, err := c.orgGetter.GetOrganizationByName(orgName)
	if err != nil {
		return nil, err
	}
	secretGroup, err := c.secretGroupGetter.GetSecretGroupByName(orgName, groupName)
	if err != nil {
		return nil, err
	}
	env, err := c.envGetter.GetEnvironmentbyName(orgName, groupName, envName)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%sorganizations/%s/secret-groups/%s/environments/%s/secrets/rollback",
		c.cfg.BackendEndpoint, org.ID, secretGroup.ID, env.ID)
	body, _ := json.Marshal(types.RollbackRequest{
		VersionID:     versionID,
		CommitMessage: commitMessage,
	})

	requestPayload := client.RequestPayload{
		Method: "POST",
		URL:    url,
		Body:   body,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	respBody, err := client.DoAuthenticatedRequest[types.SecretVersionResponse](requestPayload, c.logger, c.cfg)
	if err != nil {
		if err == cliErrors.ErrNotLoggedIn {
			c.logger.Print("You are not logged in. Please run `kavach login` to log in.")
			c.logger.Warn("User not logged in during secret rollback", map[string]interface{}{"cmd": "secret rollback", "env": envName})
			return nil, err
		}
		c.logger.Error("Failed to rollback to version", err, map[string]interface{}{"cmd": "secret rollback", "env": envName})
		return nil, err
	}

	if !respBody.Success {
		switch respBody.ErrorCode {
		case "environment_not_exist":
			return nil, cliErrors.ErrEnvironmentNotFound
		case "target_secret_version_not_found":
			return nil, cliErrors.ErrTargetSecretVersionNotFound
		case "environment_mismatch":
			return nil, cliErrors.ErrEnvironmentsMisMatch
		case "rollback_failed":
			return nil, cliErrors.ErrRollbackFailed
		case "secret_copy_failed":
			return nil, cliErrors.ErrCopySecretsFailed
		case "encryption_failed":
			return nil, cliErrors.ErrEncryptionFailed
		case "internal_error":
			return nil, cliErrors.ErrInternalServer
		default:
			c.logger.Error("Failed to rollback to version", fmt.Errorf(respBody.ErrorMsg), map[string]interface{}{"cmd": "secret rollback", "env": envName})
			return nil, fmt.Errorf(respBody.ErrorMsg)
		}
	}

	c.logger.Info("Successfully rolled back to version", map[string]interface{}{"cmd": "secret rollback", "env": envName, "new_version": respBody.Data.ID})
	return &respBody.Data, nil
}

// GetVersionDiff gets the differences between two versions
func (c *SecretHttpClient) GetVersionDiff(orgName, groupName, envName, fromVersion, toVersion string) (*types.SecretDiffResponse, error) {
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return nil, err
	}

	org, err := c.orgGetter.GetOrganizationByName(orgName)
	if err != nil {
		return nil, err
	}
	secretGroup, err := c.secretGroupGetter.GetSecretGroupByName(orgName, groupName)
	if err != nil {
		return nil, err
	}
	env, err := c.envGetter.GetEnvironmentbyName(orgName, groupName, envName)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%sorganizations/%s/secret-groups/%s/environments/%s/secrets/diff?from=%s&to=%s",
		c.cfg.BackendEndpoint, org.ID, secretGroup.ID, env.ID, fromVersion, toVersion)

	requestPayload := client.RequestPayload{
		Method: "GET",
		URL:    url,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	respBody, err := client.DoAuthenticatedRequest[types.SecretDiffResponse](requestPayload, c.logger, c.cfg)
	if err != nil {
		if err == cliErrors.ErrNotLoggedIn {
			c.logger.Print("You are not logged in. Please run `kavach login` to log in.")
			c.logger.Warn("User not logged in during secret diff", map[string]interface{}{"cmd": "secret diff"})
			return nil, err
		}
		c.logger.Error("Failed to get version diff", err, map[string]interface{}{"cmd": "secret diff"})
		return nil, err
	}

	if !respBody.Success {
		switch respBody.ErrorCode {
		case "secret_version_not_found":
			return nil, cliErrors.ErrSecretVersionNotFound
		case "decryption_failed":
			return nil, cliErrors.ErrDecryptionFailed
		case "internal_error":
			return nil, cliErrors.ErrInternalServer
		default:
			c.logger.Error("Failed to get version diff", fmt.Errorf(respBody.ErrorMsg), map[string]interface{}{"cmd": "secret diff"})
			return nil, fmt.Errorf(respBody.ErrorMsg)
		}
	}

	c.logger.Info("Version diff retrieved successfully", map[string]interface{}{"cmd": "secret diff", "changes": len(respBody.Data.Changes)})
	return &respBody.Data, nil
}

// SyncSecrets syncs secrets to a provider
func (c *SecretHttpClient) SyncSecrets(orgName, groupName, envName, provider, versionID string) (*types.SyncSecretsResponse, error) {
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return nil, err
	}

	org, err := c.orgGetter.GetOrganizationByName(orgName)
	if err != nil {
		return nil, err
	}
	secretGroup, err := c.secretGroupGetter.GetSecretGroupByName(orgName, groupName)
	if err != nil {
		return nil, err
	}
	env, err := c.envGetter.GetEnvironmentbyName(orgName, groupName, envName)
	if err != nil {
		return nil, err
	}

	// Prepare request body
	requestBody := map[string]interface{}{
		"provider": provider,
	}
	if versionID != "" {
		requestBody["version_id"] = versionID
	}

	body, _ := json.Marshal(requestBody)

	url := fmt.Sprintf("%sorganizations/%s/secret-groups/%s/environments/%s/secrets/sync",
		c.cfg.BackendEndpoint, org.ID, secretGroup.ID, env.ID)
	requestPayload := client.RequestPayload{
		Method: "POST",
		URL:    url,
		Body:   body,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	respBody, err := client.DoAuthenticatedRequest[types.SyncSecretsResponse](requestPayload, c.logger, c.cfg)
	if err != nil {
		if err == cliErrors.ErrNotLoggedIn {
			c.logger.Print("You are not logged in. Please run `kavach login` to log in.")
			c.logger.Warn("User not logged in during secret sync", map[string]interface{}{"cmd": "secret sync", "env": envName, "provider": provider})
			return nil, err
		}
		c.logger.Error("Failed to sync secrets", err, map[string]interface{}{"cmd": "secret sync", "env": envName, "provider": provider})
		return nil, err
	}

	if !respBody.Success {
		switch respBody.ErrorCode {
		case "invalid_provider_type":
			return nil, cliErrors.ErrInvalidProviderType
		case "provider_credential_not_found":
			return nil, cliErrors.ErrProviderCredentialNotFound
		case "provider_credential_validation_failed":
			return nil, cliErrors.ErrProviderCredentialValidationFailed
		case "provider_sync_failed":
			return nil, cliErrors.ErrProviderSyncFailed
		case "no_secrets_to_sync":
			return nil, cliErrors.ErrNoSecretsToSync
		case "github_environment_not_found":
			return nil, cliErrors.ErrGitHubEnvironmentNotFound
		case "github_encryption_failed":
			return nil, cliErrors.ErrGitHubEncryptionFailed
		case "gcp_invalid_location":
			return nil, cliErrors.ErrGCPInvalidLocation
		case "secret_version_not_found":
			return nil, cliErrors.ErrSecretVersionNotFound
		case "decryption_failed":
			return nil, cliErrors.ErrDecryptionFailed
		case "internal_error":
			return nil, cliErrors.ErrInternalServer
		default:
			c.logger.Error("Failed to sync secrets", fmt.Errorf(respBody.ErrorMsg), map[string]interface{}{"cmd": "secret sync", "env": envName, "provider": provider})
			return nil, fmt.Errorf(respBody.ErrorMsg)
		}
	}

	c.logger.Info("Successfully synced secrets", map[string]interface{}{
		"cmd":          "secret sync",
		"env":          envName,
		"provider":     provider,
		"synced_count": respBody.Data.SyncedCount,
		"failed_count": respBody.Data.FailedCount,
		"total_count":  respBody.Data.TotalCount,
	})
	return &respBody.Data, nil
}
