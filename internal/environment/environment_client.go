package environment

import (
	"encoding/json"
	"fmt"

	"github.com/Gkemhcs/kavach-cli/internal/auth"
	"github.com/Gkemhcs/kavach-cli/internal/client"
	"github.com/Gkemhcs/kavach-cli/internal/config"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/secretgroup"
	"github.com/Gkemhcs/kavach-cli/internal/types"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
)

// NewEnvironmentHTTPClient creates a new EnvironmentHTTPClient with the given logger, config, and group getter.
func NewEnvironmentHTTPClient(logger *utils.Logger, cfg *config.Config, groupGetter secretgroup.SecretGroupGetter) *EnvironmentHTTPClient {
	return &EnvironmentHTTPClient{
		logger:      logger,
		cfg:         cfg,
		groupGetter: groupGetter,
	}
}

// EnvironmentHTTPClient implements EnvironmentClient for making HTTP requests to the backend for environment operations.
type EnvironmentHTTPClient struct {
	logger      *utils.Logger
	cfg         *config.Config
	groupGetter secretgroup.SecretGroupGetter
}

// CreateEnvironment creates a new environment with the given parameters.
// Handles error reporting and logging.
func (c *EnvironmentHTTPClient) CreateEnvironment(environmentName, secretGroupName, orgName, description string) error {
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return err
	}
	group, err := c.groupGetter.GetSecretGroupByName(orgName, secretGroupName)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(CreateEnvironmentRequest{
		Name:        environmentName,
		Description: description,
	})
	payload := client.RequestPayload{
		Method: "POST",
		URL:    fmt.Sprintf("%sorganizations/%s/secret-groups/%s/environments", c.cfg.BackendEndpoint, group.OrganizationID, group.ID),
		Body:   body,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	respBody, err := client.DoAuthenticatedRequest[any](payload, c.logger, c.cfg)
	if err != nil {
		c.logger.Error("Failed to create environment", err, map[string]interface{}{"cmd": "env create", "env": environmentName, "group": secretGroupName, "org": orgName})
		return err
	}
	if !respBody.Success {
		if respBody.ErrorCode == "duplicate_environment" {
			c.logger.Print("Sorry, create with another name as already with this name environment exists in this secretgroup")
			c.logger.Warn("Duplicate environment during create", map[string]interface{}{"cmd": "env create", "env": environmentName, "group": secretGroupName, "org": orgName})
			return cliErrors.ErrDuplicateEnvironment
		}
		c.logger.Error("Failed to create environment", fmt.Errorf(respBody.ErrorMsg), map[string]interface{}{"cmd": "env create", "env": environmentName, "group": secretGroupName, "org": orgName})
		return fmt.Errorf(respBody.ErrorMsg)
	}
	c.logger.Info("Environment created successfully", map[string]interface{}{"cmd": "env create", "env": environmentName, "group": secretGroupName, "org": orgName})
	return nil
}

// ListEnvironment returns a list of environments for the given org and group.
func (c *EnvironmentHTTPClient) ListEnvironment(orgName, groupName string) ([]ListEnvironmentsWithMemberRow, error) {
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return nil, err
	}
	group, err := c.groupGetter.GetSecretGroupByName(orgName, groupName)
	if err != nil {
		return nil, err
	}
	payload := client.RequestPayload{
		Method: "GET",
		URL:    fmt.Sprintf("%sorganizations/%s/secret-groups/%s/environments/my", c.cfg.BackendEndpoint, group.OrganizationID, group.ID),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	respBody, err := client.DoAuthenticatedRequest[[]ListEnvironmentsWithMemberRow](payload, c.logger, c.cfg)
	if err != nil {
		c.logger.Error("Failed to list environments", err, map[string]interface{}{"cmd": "env list", "group": groupName, "org": orgName})
		return nil, err
	}
	if !respBody.Success {
		c.logger.Error("Failed to list environments", fmt.Errorf("API error: %s", respBody.ErrorMsg), map[string]interface{}{"cmd": "env list", "group": groupName, "org": orgName})
		return nil, fmt.Errorf(respBody.ErrorMsg)
	}
	c.logger.Info("Listed environments successfully", map[string]interface{}{"cmd": "env list", "count": len(respBody.Data), "group": groupName, "org": orgName})
	return respBody.Data, nil
}

// GetEnvironmentbyName returns an environment by name for the given org and group.
func (c *EnvironmentHTTPClient) GetEnvironmentbyName(orgName, groupName, environmentName string) (*EnvironmentResponseData, error) {
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return nil, err
	}
	group, err := c.groupGetter.GetSecretGroupByName(orgName, groupName)
	if err != nil {
		return nil, err
	}
	payload := client.RequestPayload{
		Method: "GET",
		URL:    fmt.Sprintf("%sorganizations/%s/secret-groups/%s/environments/by-name/%s", c.cfg.BackendEndpoint, group.OrganizationID, group.ID, environmentName),
		Headers: map[string]string{
			"Content-Type": "aplication/json",
		},
	}
	respBody, err := client.DoAuthenticatedRequest[*EnvironmentResponseData](payload, c.logger, c.cfg)
	if err != nil {
		c.logger.Error("Failed to get environment by name", err, map[string]interface{}{"cmd": "env get", "env": environmentName, "group": groupName, "org": orgName})
		return nil, err
	}
	if !respBody.Success {
		if respBody.ErrorCode == "environment_not_exist" {
			c.logger.Warn("Environment not found during get by name", map[string]interface{}{"cmd": "env get", "env": environmentName, "group": groupName, "org": orgName})
			return nil, cliErrors.ErrEnvironmentNotFound
		}
		c.logger.Error("Failed to get environment by name", fmt.Errorf(respBody.ErrorMsg), map[string]interface{}{"cmd": "env get", "env": environmentName, "group": groupName, "org": orgName})
		return nil, fmt.Errorf(respBody.ErrorMsg)
	}
	c.logger.Info("Environment fetched by name successfully", map[string]interface{}{"cmd": "env get", "env": environmentName, "group": groupName, "org": orgName})
	return respBody.Data, nil
}

// DeleteEnvironment deletes an environment by name for the given org and group.
func (c *EnvironmentHTTPClient) DeleteEnvironment(orgName, groupName, envName string) error {
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return err
	}
	group, err := c.groupGetter.GetSecretGroupByName(orgName, groupName)
	if err != nil {
		c.logger.Error("Failed to get secret group during environment delete", err, map[string]interface{}{"cmd": "env delete", "group": groupName, "org": orgName})
		return err
	}
	env, err := c.GetEnvironmentbyName(orgName, groupName, envName)
	if err != nil {
		c.logger.Error("Failed to get environment during delete", err, map[string]interface{}{"cmd": "env delete", "env": envName, "group": groupName, "org": orgName})
		return err
	}
	payload := client.RequestPayload{
		Method: "DELETE",
		URL:    fmt.Sprintf("%sorganizations/%s/secret-groups/%s/environments/%s", c.cfg.BackendEndpoint, group.OrganizationID, group.ID, env.ID),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	respBody, err := client.DoAuthenticatedRequest[any](payload, c.logger, c.cfg)
	if err != nil {
		c.logger.Error("Failed to delete environment", err, map[string]interface{}{"cmd": "env delete", "env": envName, "group": groupName, "org": orgName})
		return err
	}
	if !respBody.Success {
		c.logger.Warn("Environment not found during delete", map[string]interface{}{"cmd": "env delete", "env": envName, "group": groupName, "org": orgName})
		return cliErrors.ErrEnvironmentNotFound
	}
	c.logger.Info("Environment deleted successfully", map[string]interface{}{"cmd": "env delete", "env": envName, "group": groupName, "org": orgName})
	return nil
}

// GrantRoleBinding grants a role to a user or user group on an environment.
// Handles authentication, environment validation, and provides detailed error messages.
func (c *EnvironmentHTTPClient) GrantRoleBinding(req types.GrantRoleBindingInput) error {
	// Validate authentication
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return err
	}

	// Get environment details for role binding
	env, err := c.GetEnvironmentbyName(req.OrgName, req.SecretGroupName, req.EnvironmentName)
	if err != nil {
		c.logger.Error("Failed to get environment during role binding grant", err, map[string]interface{}{
			"cmd":         "env grant",
			"environment": req.EnvironmentName,
			"secretGroup": req.SecretGroupName,
			"org":         req.OrgName,
		})
		return err
	}

	// Prepare role binding parameters
	params := types.GrantRoleBindingRequest{
		UserName:       req.UserName,
		GroupName:      req.GroupName,
		Role:           req.Role,
		ResourceType:   "environment",
		ResourceID:     env.ID,
		OrganizationID: env.OrganizationID,
		SecretGroupID:  &env.SecretGroupID,
		EnvironmentID:  &env.ID,
	}

	body, err := json.Marshal(params)
	if err != nil {
		c.logger.Error("Failed to marshal role binding request", err, map[string]interface{}{
			"cmd":         "env grant",
			"environment": req.EnvironmentName,
			"org":         req.OrgName,
		})
		return err
	}

	// Prepare and execute request
	payload := client.RequestPayload{
		Method: "POST",
		URL:    fmt.Sprintf("%spermissions/grant", c.cfg.BackendEndpoint),
		Body:   body,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	respBody, err := client.DoAuthenticatedRequest[any](payload, c.logger, c.cfg)
	if err != nil {
		c.logger.Error("Failed to grant environment role binding", err, map[string]interface{}{
			"cmd":         "env grant",
			"environment": req.EnvironmentName,
			"org":         req.OrgName,
		})
		return err
	}

	// Handle response
	if !respBody.Success {
		if respBody.ErrorCode == "duplicate_role_binding" {
			return cliErrors.ErrDuplicateRoleBinding
		}
		c.logger.Error("Failed to grant environment role binding", fmt.Errorf(respBody.ErrorMsg), map[string]interface{}{
			"cmd":         "env grant",
			"environment": req.EnvironmentName,
			"org":         req.OrgName,
		})
		return fmt.Errorf(respBody.ErrorMsg)
	}

	c.logger.Info("Environment role binding granted successfully", map[string]interface{}{
		"cmd":         "env grant",
		"environment": req.EnvironmentName,
		"org":         req.OrgName,
	})
	return nil
}

// RevokeRoleBinding revokes a role from a user or user group on an environment.
// Handles authentication, environment validation, and provides detailed error messages.
func (c *EnvironmentHTTPClient) RevokeRoleBinding(req types.RevokeRoleBindingInput) error {
	// Validate authentication
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return err
	}

	// Get environment details for role binding revocation
	env, err := c.GetEnvironmentbyName(req.OrgName, req.SecretGroupName, req.EnvironmentName)
	if err != nil {
		c.logger.Error("Failed to get environment during role binding revocation", err, map[string]interface{}{
			"cmd":         "env revoke",
			"environment": req.EnvironmentName,
			"secretGroup": req.SecretGroupName,
			"org":         req.OrgName,
		})
		return err
	}

	// Prepare role binding revocation parameters
	params := types.RevokeRoleBindingRequest{
		UserName:       req.UserName,
		GroupName:      req.GroupName,
		Role:           req.Role,
		ResourceType:   "environment",
		ResourceID:     env.ID,
		OrganizationID: env.OrganizationID,
		SecretGroupID:  &env.SecretGroupID,
		EnvironmentID:  &env.ID,
	}

	body, err := json.Marshal(params)
	if err != nil {
		c.logger.Error("Failed to marshal role binding revocation request", err, map[string]interface{}{
			"cmd":         "env revoke",
			"environment": req.EnvironmentName,
			"org":         req.OrgName,
		})
		return err
	}

	// Prepare and execute request
	payload := client.RequestPayload{
		Method: "DELETE",
		URL:    fmt.Sprintf("%spermissions/revoke", c.cfg.BackendEndpoint),
		Body:   body,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	respBody, err := client.DoAuthenticatedRequest[any](payload, c.logger, c.cfg)
	if err != nil {
		c.logger.Error("Failed to revoke environment role binding", err, map[string]interface{}{
			"cmd":         "env revoke",
			"environment": req.EnvironmentName,
			"org":         req.OrgName,
		})
		return err
	}

	// Handle response
	if !respBody.Success {
		if respBody.ErrorCode == "role_binding_not_found" {
			return cliErrors.ErrRoleBindingNotFound
		}
		c.logger.Error("Failed to revoke environment role binding", fmt.Errorf(respBody.ErrorMsg), map[string]interface{}{
			"cmd":         "env revoke",
			"environment": req.EnvironmentName,
			"org":         req.OrgName,
		})
		return fmt.Errorf(respBody.ErrorMsg)
	}

	c.logger.Info("Environment role binding revoked successfully", map[string]interface{}{
		"cmd":         "env revoke",
		"environment": req.EnvironmentName,
		"org":         req.OrgName,
	})
	return nil
}

// ListRoleBindings lists all role bindings for an environment with resolved names.
func (c *EnvironmentHTTPClient) ListRoleBindings(orgName, groupName, envName string) ([]RoleBinding, error) {
	// Validate authentication
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return nil, err
	}

	// Get environment details
	env, err := c.GetEnvironmentbyName(orgName, groupName, envName)
	if err != nil {
		c.logger.Error("Failed to get environment during role bindings listing", err, map[string]interface{}{
			"cmd":   "env list-bindings",
			"org":   orgName,
			"group": groupName,
			"env":   envName,
		})
		return nil, err
	}

	// Prepare and execute request
	payload := client.RequestPayload{
		Method: "GET",
		URL:    fmt.Sprintf("%sorganizations/%s/secret-groups/%s/environments/%s/role-bindings", c.cfg.BackendEndpoint, env.OrganizationID, env.SecretGroupID, env.ID),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	respBody, err := client.DoAuthenticatedRequest[ListRoleBindingsResponse](payload, c.logger, c.cfg)
	if err != nil {
		c.logger.Error("Failed to list environment role bindings", err, map[string]interface{}{
			"cmd":   "env list-bindings",
			"org":   orgName,
			"group": groupName,
			"env":   envName,
		})
		return nil, err
	}

	// Handle response
	if !respBody.Success {
		c.logger.Error("Failed to list environment role bindings", fmt.Errorf(respBody.ErrorMsg), map[string]interface{}{
			"cmd":   "env list-bindings",
			"org":   orgName,
			"group": groupName,
			"env":   envName,
		})
		return nil, cliErrors.HandleRoleBindingListError(respBody.ErrorCode, respBody.ErrorMsg)
	}

	c.logger.Info("Environment role bindings listed successfully", map[string]interface{}{
		"cmd":   "env list-bindings",
		"org":   orgName,
		"group": groupName,
		"env":   envName,
		"count": len(respBody.Data.Bindings),
	})

	return respBody.Data.Bindings, nil
}
