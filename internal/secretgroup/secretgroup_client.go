package secretgroup

import (
	"encoding/json"
	"fmt"

	"github.com/Gkemhcs/kavach-cli/internal/auth"
	"github.com/Gkemhcs/kavach-cli/internal/client"
	"github.com/Gkemhcs/kavach-cli/internal/config"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/org"
	"github.com/Gkemhcs/kavach-cli/internal/types"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
)

// NewSecretGroupHttpClient creates a new SecretGroupHttpClient with the given config, logger, and org getter.
func NewSecretGroupHttpClient(cfg *config.Config, logger *utils.Logger, orgGetter org.OrgGetterClient) *SecretGroupHttpClient {
	return &SecretGroupHttpClient{
		cfg:       cfg,
		logger:    logger,
		orgGetter: orgGetter,
	}
}

// SecretGroupHttpClient implements SecretGroupClient for making HTTP requests to the backend for secret group operations.
type SecretGroupHttpClient struct {
	cfg       *config.Config
	logger    *utils.Logger
	orgGetter org.OrgGetterClient
}

// CreateSecretGroup creates a new secret group with the given parameters.
// Handles error reporting and logging.
func (c *SecretGroupHttpClient) CreateSecretGroup(name, description, orgName string) error {
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return err
	}
	org, err := c.orgGetter.GetOrganizationByName(orgName)
	if err != nil {
		return err
	}
	req := CreateSecretGroupRequest{
		Name:        name,
		Description: description,
	}
	body, _ := json.Marshal(req)
	payload := client.RequestPayload{
		Method: "POST",
		URL:    fmt.Sprintf("%sorganizations/%s/secret-groups/", c.cfg.BackendEndpoint, org.ID),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: body,
	}
	respBody, err := client.DoAuthenticatedRequest[any](payload, c.logger, c.cfg)
	if err != nil {
		c.logger.Error("Failed to create secret group", err, map[string]interface{}{"cmd": "group create", "group": name, "org": orgName})
		return err
	}
	if !respBody.Success {
		if respBody.ErrorCode == "duplicate_secret_group" {
			c.logger.Print("Sorry, create with another name as already with this name secret group exists in this organization")
			c.logger.Warn("Duplicate secret group during create", map[string]interface{}{"cmd": "group create", "group": name, "org": orgName})
			return cliErrors.ErrDuplicateSecretGroup
		}
		c.logger.Error("Failed to create secret group", fmt.Errorf(respBody.ErrorMsg), map[string]interface{}{"cmd": "group create", "group": name, "org": orgName})
		return fmt.Errorf(respBody.ErrorMsg)
	}
	c.logger.Info("Secret group created successfully", map[string]interface{}{"cmd": "group create", "group": name, "org": orgName})
	return nil
}

// ListSecretGroups returns a list of secret groups for the given organization.
func (c *SecretGroupHttpClient) ListSecretGroups(orgName string) ([]ListSecretGroupsWithMemberRow, error) {
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return nil, err
	}
	org, err := c.orgGetter.GetOrganizationByName(orgName)

	if err != nil {
		return nil, err
	}
	payload := client.RequestPayload{
		Method: "GET",
		URL:    fmt.Sprintf("%sorganizations/%s/secret-groups/my", c.cfg.BackendEndpoint, org.ID),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	respBody, err := client.DoAuthenticatedRequest[[]ListSecretGroupsWithMemberRow](payload, c.logger, c.cfg)
	if err != nil {
		c.logger.Error("Failed to list secret groups", err, map[string]interface{}{"cmd": "group list", "org": orgName})
		return nil, err
	}
	c.logger.Info("Listed secret groups successfully", map[string]interface{}{"cmd": "group list", "count": len(respBody.Data), "org": orgName})
	return respBody.Data, nil
}

// GetSecretGroupByName returns a secret group by name for the given organization.
func (c *SecretGroupHttpClient) GetSecretGroupByName(orgName, groupName string) (*SecretGroupResponseData, error) {
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return nil, err
	}
	org, err := c.orgGetter.GetOrganizationByName(orgName)

	if err != nil {
		return nil, err
	}
	payload := client.RequestPayload{
		Method: "GET",
		URL:    fmt.Sprintf("%sorganizations/%s/secret-groups/by-name/%s", c.cfg.BackendEndpoint, org.ID, groupName),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	respBody, err := client.DoAuthenticatedRequest[SecretGroupResponseData](payload, c.logger, c.cfg)
	if err != nil {
		c.logger.Error("Failed to get secret group by name", err, map[string]interface{}{"cmd": "group get", "group": groupName, "org": orgName})
		return nil, err
	}
	if !respBody.Success {
		if respBody.ErrorCode == "secretgroup_not_exist" {
			c.logger.Warn("Secret group not found during get by name", map[string]interface{}{"cmd": "group get", "group": groupName, "org": orgName})
			return nil, cliErrors.ErrSecretGroupNotFound
		}
		c.logger.Error("Failed to get secret group by name", fmt.Errorf(respBody.ErrorMsg), map[string]interface{}{"cmd": "group get", "group": groupName, "org": orgName})
		return nil, fmt.Errorf(respBody.ErrorMsg)
	}
	c.logger.Info("Secret group fetched by name successfully", map[string]interface{}{"cmd": "group get", "group": groupName, "org": orgName})
	return &respBody.Data, nil
}

// DeleteSecretGroupByName deletes a secret group by name for the given organization.
func (c *SecretGroupHttpClient) DeleteSecretGroupByName(orgName, groupName string) error {
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return err
	}
	group, err := c.GetSecretGroupByName(orgName, groupName)
	if err != nil {
		c.logger.Error("Failed to get secret group during delete", err, map[string]interface{}{"cmd": "group delete", "group": groupName, "org": orgName})
		return err
	}
	payload := client.RequestPayload{
		Method: "DELETE",
		URL:    fmt.Sprintf("%sorganizations/%s/secret-groups/%s", c.cfg.BackendEndpoint, group.OrganizationID, group.ID),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	respBody, err := client.DoAuthenticatedRequest[any](payload, c.logger, c.cfg)
	if !respBody.Success {
		if respBody.ErrorCode == "foreign_key_constraint_violation" {
			return cliErrors.ErrForeignKeyViolation
		}
		if respBody.ErrorCode == "organisation_not_exist" {
			c.logger.Warn("Organization not found during delete", map[string]interface{}{"cmd": "group delete", "org": orgName})
			return cliErrors.ErrOrganizationNotFound
		}
		if respBody.ErrorCode == "secretgroup_not_exist" {
			c.logger.Warn("secret group not found during delete", map[string]interface{}{"cmd": "group delete", "group": groupName})
			return cliErrors.ErrOrganizationNotFound
		}
	}
	if err != nil {
		c.logger.Error("Failed to delete secret group", err, map[string]interface{}{"cmd": "group delete", "group": groupName, "org": orgName})
		return err
	}
	c.logger.Info("Secret group deleted successfully", map[string]interface{}{"cmd": "group delete", "group": groupName, "org": orgName})
	return nil
}

// GrantRoleBinding grants a role to a user or user group on a secret group.
// Handles authentication, secret group validation, and provides detailed error messages.
func (c *SecretGroupHttpClient) GrantRoleBinding(req types.GrantRoleBindingInput) error {
	// Validate authentication
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return err
	}

	// Get secret group details for role binding
	secret_group, err := c.GetSecretGroupByName(req.OrgName, req.SecretGroupName)
	if err != nil {
		c.logger.Error("Failed to get secret group during role binding grant", err, map[string]interface{}{
			"cmd":         "group grant",
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
		ResourceType:   "secret_group",
		ResourceID:     secret_group.ID,
		OrganizationID: secret_group.OrganizationID,
		SecretGroupID:  &secret_group.ID,
	}

	body, err := json.Marshal(params)
	if err != nil {
		c.logger.Error("Failed to marshal role binding request", err, map[string]interface{}{
			"cmd":         "group grant",
			"secretGroup": req.SecretGroupName,
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
		c.logger.Error("Failed to grant secret group role binding", err, map[string]interface{}{
			"cmd":         "group grant",
			"secretGroup": req.SecretGroupName,
			"org":         req.OrgName,
		})
		return err
	}

	// Handle response
	if !respBody.Success {

		if respBody.ErrorCode == "duplicate_role_binding" {
			return cliErrors.ErrDuplicateRoleBinding
		}
		c.logger.Error("Failed to grant secret group role binding", fmt.Errorf(respBody.ErrorMsg), map[string]interface{}{
			"cmd":         "group grant",
			"secretGroup": req.SecretGroupName,
			"org":         req.OrgName,
		})
		return fmt.Errorf(respBody.ErrorMsg)
	}

	c.logger.Info("Secret group role binding granted successfully", map[string]interface{}{
		"cmd":         "group grant",
		"secretGroup": req.SecretGroupName,
		"org":         req.OrgName,
	})
	return nil
}

// RevokeRoleBinding revokes a role from a user or user group on a secret group.
// Handles authentication, secret group validation, and provides detailed error messages.
func (c *SecretGroupHttpClient) RevokeRoleBinding(req types.RevokeRoleBindingInput) error {
	// Validate authentication
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return err
	}

	// Get secret group details for role binding revocation
	secret_group, err := c.GetSecretGroupByName(req.OrgName, req.SecretGroupName)
	if err != nil {
		c.logger.Error("Failed to get secret group during role binding revocation", err, map[string]interface{}{
			"cmd":         "group revoke",
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
		ResourceType:   "secret_group",
		ResourceID:     secret_group.ID,
		OrganizationID: secret_group.OrganizationID,
		SecretGroupID:  &secret_group.ID,
	}

	body, err := json.Marshal(params)
	if err != nil {
		c.logger.Error("Failed to marshal role binding revocation request", err, map[string]interface{}{
			"cmd":         "group revoke",
			"secretGroup": req.SecretGroupName,
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
		c.logger.Error("Failed to revoke secret group role binding", err, map[string]interface{}{
			"cmd":         "group revoke",
			"secretGroup": req.SecretGroupName,
			"org":         req.OrgName,
		})
		return err
	}

	// Handle response
	if !respBody.Success {
		if respBody.ErrorCode == "role_binding_not_found" {
			return cliErrors.ErrRoleBindingNotFound
		}
		c.logger.Error("Failed to revoke secret group role binding", fmt.Errorf(respBody.ErrorMsg), map[string]interface{}{
			"cmd":         "group revoke",
			"secretGroup": req.SecretGroupName,
			"org":         req.OrgName,
		})
		return fmt.Errorf(respBody.ErrorMsg)
	}

	c.logger.Info("Secret group role binding revoked successfully", map[string]interface{}{
		"cmd":         "group revoke",
		"secretGroup": req.SecretGroupName,
		"org":         req.OrgName,
	})
	return nil
}
