package provider

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

// ProviderHttpClient implements ProviderClient for making HTTP requests to the backend
type ProviderHttpClient struct {
	cfg               *config.Config
	logger            *utils.Logger
	orgGetter         org.OrgGetterClient
	secretGroupGetter secretgroup.SecretGroupGetter
	envGetter         environment.EnvironmentGetter
}

// NewProviderHttpClient creates a new ProviderHttpClient with the given config and logger
func NewProviderHttpClient(cfg *config.Config,
	logger *utils.Logger, orgGetter org.OrgGetterClient,
	secretGroupGetter secretgroup.SecretGroupGetter,
	envGetter environment.EnvironmentGetter) *ProviderHttpClient {
	return &ProviderHttpClient{
		cfg:               cfg,
		logger:            logger,
		orgGetter:         orgGetter,
		secretGroupGetter: secretGroupGetter,
		envGetter:         envGetter,
	}
}

// CreateProviderCredential creates a new provider credential for an environment
func (c *ProviderHttpClient) CreateProviderCredential(orgName, groupName, envName, providerName string, credentials, config map[string]interface{}) (*types.ProviderCredentialResponse, error) {
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

	body, _ := json.Marshal(CreateProviderCredentialRequest{
		Provider:    ProviderType(providerName),
		Credentials: credentials,
		Config:      config,
	})

	url := fmt.Sprintf("%sorganizations/%s/secret-groups/%s/environments/%s/providers/credentials",
		c.cfg.BackendEndpoint, org.ID, secretGroup.ID, env.ID)
	requestPayload := client.RequestPayload{
		Method: "POST",
		URL:    url,
		Body:   body,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	respBody, err := client.DoAuthenticatedRequest[types.ProviderCredentialResponse](requestPayload, c.logger, c.cfg)
	if err != nil {
		if err == cliErrors.ErrNotLoggedIn {
			c.logger.Print("You are not logged in. Please run `kavach login` to log in.")
			c.logger.Warn("User not logged in during provider create", map[string]interface{}{"cmd": "provider create", "env": envName})
			return nil, err
		}
		c.logger.Error("Failed to create provider credential", err, map[string]interface{}{"cmd": "provider create", "env": envName})
		return nil, err
	}

	if !respBody.Success {
		switch respBody.ErrorCode {
		case "invalid_provider_type":
			return nil, cliErrors.ErrInvalidProviderType
		case "invalid_provider_data":
			return nil, cliErrors.ErrInvalidProviderData
		case "provider_credential_exists":
			return nil, cliErrors.ErrProviderCredentialExists
		case "provider_encryption_failed":
			return nil, cliErrors.ErrProviderEncryptionFailed
		case "provider_credential_create_failed":
			return nil, cliErrors.ErrProviderCredentialCreateFailed
		default:
			c.logger.Error("Failed to create provider credential", fmt.Errorf(respBody.ErrorMsg), map[string]interface{}{"cmd": "provider create", "env": envName})
			return nil, fmt.Errorf(respBody.ErrorMsg)
		}
	}

	return &respBody.Data, nil
}

// GetProviderCredential retrieves a provider credential by environment ID and provider
func (c *ProviderHttpClient) GetProviderCredential(orgName, groupName, envName, providerName string) (*types.ProviderCredentialResponse, error) {
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

	url := fmt.Sprintf("%sorganizations/%s/secret-groups/%s/environments/%s/providers/credentials/%s",
		c.cfg.BackendEndpoint, org.ID, secretGroup.ID, env.ID, providerName)
	requestPayload := client.RequestPayload{
		Method: "GET",
		URL:    url,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	respBody, err := client.DoAuthenticatedRequest[types.ProviderCredentialResponse](requestPayload, c.logger, c.cfg)
	if err != nil {
		if err == cliErrors.ErrNotLoggedIn {
			c.logger.Print("You are not logged in. Please run `kavach login` to log in.")
			c.logger.Warn("User not logged in during provider get", map[string]interface{}{"cmd": "provider get", "env": envName})
			return nil, err
		}
		c.logger.Error("Failed to get provider credential", err, map[string]interface{}{"cmd": "provider get", "env": envName})
		return nil, err
	}

	if !respBody.Success {
		switch respBody.ErrorCode {
		case "provider_credential_not_found":
			return nil, cliErrors.ErrProviderCredentialNotFound
		case "provider_credential_get_failed":
			return nil, cliErrors.ErrProviderCredentialGetFailed
		default:
			c.logger.Error("Failed to get provider credential", fmt.Errorf(respBody.ErrorMsg), map[string]interface{}{"cmd": "provider get", "env": envName})
			return nil, fmt.Errorf(respBody.ErrorMsg)
		}
	}

	return &respBody.Data, nil
}

// ListProviderCredentials lists all provider credentials for an environment
func (c *ProviderHttpClient) ListProviderCredentials(orgName, groupName, envName string) ([]types.ProviderCredentialResponse, error) {
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

	url := fmt.Sprintf("%sorganizations/%s/secret-groups/%s/environments/%s/providers/credentials",
		c.cfg.BackendEndpoint, org.ID, secretGroup.ID, env.ID)
	requestPayload := client.RequestPayload{
		Method: "GET",
		URL:    url,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	respBody, err := client.DoAuthenticatedRequest[[]types.ProviderCredentialResponse](requestPayload, c.logger, c.cfg)
	if err != nil {
		if err == cliErrors.ErrNotLoggedIn {
			c.logger.Print("You are not logged in. Please run `kavach login` to log in.")
			c.logger.Warn("User not logged in during provider list", map[string]interface{}{"cmd": "provider list", "env": envName})
			return nil, err
		}
		c.logger.Error("Failed to list provider credentials", err, map[string]interface{}{"cmd": "provider list", "env": envName})
		return nil, err
	}

	if !respBody.Success {
		switch respBody.ErrorCode {
		case "provider_credential_list_failed":
			return nil, cliErrors.ErrProviderCredentialListFailed
		default:
			c.logger.Error("Failed to list provider credentials", fmt.Errorf(respBody.ErrorMsg), map[string]interface{}{"cmd": "provider list", "env": envName})
			return nil, fmt.Errorf(respBody.ErrorMsg)
		}
	}

	return respBody.Data, nil
}

// UpdateProviderCredential updates an existing provider credential
func (c *ProviderHttpClient) UpdateProviderCredential(orgName, groupName, envName, providerName string, credentials, config map[string]interface{}) (*types.ProviderCredentialResponse, error) {
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

	body, _ := json.Marshal(UpdateProviderCredentialRequest{
		Credentials: credentials,
		Config:      config,
	})

	url := fmt.Sprintf("%sorganizations/%s/secret-groups/%s/environments/%s/providers/credentials/%s",
		c.cfg.BackendEndpoint, org.ID, secretGroup.ID, env.ID, providerName)
	requestPayload := client.RequestPayload{
		Method: "PUT",
		URL:    url,
		Body:   body,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	respBody, err := client.DoAuthenticatedRequest[types.ProviderCredentialResponse](requestPayload, c.logger, c.cfg)
	if err != nil {
		if err == cliErrors.ErrNotLoggedIn {
			c.logger.Print("You are not logged in. Please run `kavach login` to log in.")
			c.logger.Warn("User not logged in during provider update", map[string]interface{}{"cmd": "provider update", "env": envName})
			return nil, err
		}
		c.logger.Error("Failed to update provider credential", err, map[string]interface{}{"cmd": "provider update", "env": envName})
		return nil, err
	}

	if !respBody.Success {
		switch respBody.ErrorCode {
		case "provider_credential_not_found":
			return nil, cliErrors.ErrProviderCredentialNotFound
		case "invalid_provider_type":
			return nil, cliErrors.ErrInvalidProviderType
		case "invalid_provider_data":
			return nil, cliErrors.ErrInvalidProviderData
		case "provider_encryption_failed":
			return nil, cliErrors.ErrProviderEncryptionFailed
		case "provider_credential_update_failed":
			return nil, cliErrors.ErrProviderCredentialUpdateFailed
		default:
			c.logger.Error("Failed to update provider credential", fmt.Errorf(respBody.ErrorMsg), map[string]interface{}{"cmd": "provider update", "env": envName})
			return nil, fmt.Errorf(respBody.ErrorMsg)
		}
	}

	return &respBody.Data, nil
}

// DeleteProviderCredential deletes a provider credential
func (c *ProviderHttpClient) DeleteProviderCredential(orgName, groupName, envName, providerName string) error {
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return err
	}

	org, err := c.orgGetter.GetOrganizationByName(orgName)
	if err != nil {
		return err
	}
	secretGroup, err := c.secretGroupGetter.GetSecretGroupByName(orgName, groupName)
	if err != nil {
		return err
	}
	env, err := c.envGetter.GetEnvironmentbyName(orgName, groupName, envName)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%sorganizations/%s/secret-groups/%s/environments/%s/providers/credentials/%s",
		c.cfg.BackendEndpoint, org.ID, secretGroup.ID, env.ID, providerName)
	requestPayload := client.RequestPayload{
		Method: "DELETE",
		URL:    url,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	respBody, err := client.DoAuthenticatedRequest[map[string]interface{}](requestPayload, c.logger, c.cfg)
	if err != nil {
		if err == cliErrors.ErrNotLoggedIn {
			c.logger.Print("You are not logged in. Please run `kavach login` to log in.")
			c.logger.Warn("User not logged in during provider delete", map[string]interface{}{"cmd": "provider delete", "env": envName})
			return err
		}
		c.logger.Error("Failed to delete provider credential", err, map[string]interface{}{"cmd": "provider delete", "env": envName})
		return err
	}

	if !respBody.Success {
		switch respBody.ErrorCode {
		case "provider_credential_delete_failed":
			return cliErrors.ErrProviderCredentialDeleteFailed
		default:
			c.logger.Error("Failed to delete provider credential", fmt.Errorf(respBody.ErrorMsg), map[string]interface{}{"cmd": "provider delete", "env": envName})
			return fmt.Errorf(respBody.ErrorMsg)
		}
	}

	return nil
}
