package org

import (
	"encoding/json"
	"fmt"

	"github.com/Gkemhcs/kavach-cli/internal/auth"
	"github.com/Gkemhcs/kavach-cli/internal/client"
	"github.com/Gkemhcs/kavach-cli/internal/config"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/types"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
)

// OrgHttpClient implements OrgClient for making HTTP requests to the backend for organization operations.
type OrgHttpClient struct {
	cfg    *config.Config
	logger *utils.Logger
}

// NewOrgHttpClient creates a new OrgHttpClient with the given config and logger.
func NewOrgHttpClient(cfg *config.Config, logger *utils.Logger) *OrgHttpClient {
	return &OrgHttpClient{
		cfg:    cfg,
		logger: logger,
	}
}

// CreateOrganization creates a new organization with the given name and description.
// Handles error reporting and logging.
func (c *OrgHttpClient) CreateOrganization(name, description string) error {
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return err
	}
	body, _ := json.Marshal(CreateOrgRequest{Name: name, Description: description})
	orgBaseUrl := c.cfg.BackendEndpoint + "organizations/"
	requestPayload := client.RequestPayload{
		Method: "POST",
		URL:    orgBaseUrl,
		Body:   body,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	respBody, err := client.DoAuthenticatedRequest[any](requestPayload, c.logger, c.cfg)
	if err != nil {
		if err == cliErrors.ErrNotLoggedIn {
			c.logger.Print("You are not logged in. Please run `kavach login` to log in.")
			c.logger.Warn("User not logged in during org create", map[string]interface{}{"cmd": "org create", "org": name})
			return nil
		}
		c.logger.Error("Failed to create organization", err, map[string]interface{}{"cmd": "org create", "org": name})
		return err
	}
	if !respBody.Success {
		if respBody.ErrorCode == "duplicate_organization" {
			c.logger.Print("Sorry, create with another name as already with this name organization exists.")
			c.logger.Warn("Duplicate organization during create", map[string]interface{}{"cmd": "org create", "org": name})
			return cliErrors.ErrDuplicateOrganisation
		}
		c.logger.Error("Failed to create organization", fmt.Errorf(respBody.ErrorMsg), map[string]interface{}{"cmd": "org create", "org": name})
		return fmt.Errorf(respBody.ErrorMsg)
	}
	c.logger.Info("Organization created successfully", map[string]interface{}{"cmd": "org create", "org": name})
	return nil
}

// ListMyOrganizations returns a list of organizations the user is a member of.
func (c *OrgHttpClient) ListMyOrganizations() ([]ListMembersOfOrganizationRow, error) {
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return nil, err
	}
	requestPayload := client.RequestPayload{
		Method: "GET",
		URL:    c.cfg.BackendEndpoint + "organizations/my",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	respBody, err := client.DoAuthenticatedRequest[[]ListMembersOfOrganizationRow](requestPayload, c.logger, c.cfg)
	if err != nil {
		if err == cliErrors.ErrNotLoggedIn {
			c.logger.Print("You are not logged in. Please run `kavach login` to log in.")
			c.logger.Warn("User not logged in during org list", map[string]interface{}{"cmd": "org list"})
			return nil, err
		}
		c.logger.Error("Failed to list organizations", err, map[string]interface{}{"cmd": "org list"})
		return nil, err
	}
	orgList := respBody.Data
	c.logger.Info("Listed organizations successfully", map[string]interface{}{"cmd": "org list", "count": len(orgList)})
	return orgList, nil
}

// DeleteOrganization deletes an organization by name.
// Handles error reporting and logging.
func (c *OrgHttpClient) DeleteOrganization(name string) error {
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return err
	}
	org, err := c.GetOrganizationByName(name)
	if err != nil {
		return err
	}

	payload := client.RequestPayload{
		Method: "DELETE",
		URL:    fmt.Sprintf("%sorganizations/%s", c.cfg.BackendEndpoint, org.ID),
		Headers: map[string]string{
			"Content-Type": "applications/json",
		},
	}
	deleteRespBody, err := client.DoAuthenticatedRequest[any](payload, c.logger, c.cfg)
	if err != nil {
		c.logger.Print("Something went wrong, please try again.")
		c.logger.Error("Failed to delete organization", err, map[string]interface{}{"cmd": "org delete", "org": name})
		return err
	}
	if !deleteRespBody.Success {

		if deleteRespBody.ErrorCode == "foreign_key_constraint_violation" {
			return cliErrors.ErrForeignKeyViolation
		}

		if deleteRespBody.ErrorCode == "organisation_not_exist" {
			c.logger.Warn("Organization not found during delete", map[string]interface{}{"cmd": "org delete", "org": name})
			return cliErrors.ErrOrganizationNotFound
		}
		c.logger.Error("Failed to delete organization", fmt.Errorf(deleteRespBody.ErrorMsg), map[string]interface{}{"cmd": "org delete", "org": name})
		return fmt.Errorf(deleteRespBody.ErrorMsg)
	}
	c.logger.Info("Organization deleted successfully", map[string]interface{}{"cmd": "org delete", "org": name})
	return nil
}

// GetOrganizationByName returns an organization by name.
// Handles error reporting and logging.
func (c *OrgHttpClient) GetOrganizationByName(name string) (*Organization, error) {
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return nil, err
	}
	payload := client.RequestPayload{
		Method: "GET",
		URL:    c.cfg.BackendEndpoint + "organizations/by-name/" + name,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	respBody, err := client.DoAuthenticatedRequest[Organization](payload, c.logger, c.cfg)
	if err != nil {
		c.logger.Error("Failed to get organization by name", err, map[string]interface{}{"cmd": "org get", "org": name})
		return nil, err
	}
	if !respBody.Success {
		if respBody.ErrorCode == "organisation_not_exist" {
			c.logger.Warn("Organization not found during get by name", map[string]interface{}{"cmd": "org get", "org": name})
			return nil, cliErrors.ErrOrganizationNotFound
		}
		c.logger.Print(respBody.ErrorMsg)
		c.logger.Error("Failed to get organization by name", fmt.Errorf(respBody.ErrorMsg), map[string]interface{}{"cmd": "org get", "org": name})
	}
	c.logger.Info("Organization fetched by name successfully", map[string]interface{}{"cmd": "org get", "org": name})
	return &respBody.Data, nil
}

// GrantRoleBinding grants a role to a user or user group on an organization.
// Handles authentication, organization validation, and provides detailed error messages.
func (c *OrgHttpClient) GrantRoleBinding(req types.GrantRoleBindingInput) error {
	// Validate authentication
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return err
	}

	// Get organization details for role binding
	org, err := c.GetOrganizationByName(req.OrgName)
	if err != nil {
		c.logger.Error("Failed to get organization during role binding grant", err, map[string]interface{}{
			"cmd": "org grant",
			"org": req.OrgName,
		})
		return err
	}

	// Prepare role binding parameters
	body, err := json.Marshal(GrantRoleBindingRequest{
		UserName:       req.UserName,
		GroupName:      req.GroupName,
		Role:           req.Role,
		ResourceType:   "organization",
		ResourceID:     org.ID.String(),
		OrganizationID: org.ID.String(),
	})
	if err != nil {
		c.logger.Error("Failed to marshal role binding request", err, map[string]interface{}{
			"cmd": "org grant",
			"org": req.OrgName,
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
		c.logger.Error("Failed to grant organization role binding", err, map[string]interface{}{
			"cmd": "org grant",
			"org": req.OrgName,
		})
		return err
	}

	// Handle response
	if !respBody.Success {
		if respBody.ErrorCode == "duplicate_role_binding" {
			return cliErrors.ErrDuplicateRoleBinding
		}
		c.logger.Error("Failed to grant organization role binding", fmt.Errorf(respBody.ErrorMsg), map[string]interface{}{
			"cmd": "org grant",
			"org": req.OrgName,
		})
		return fmt.Errorf(respBody.ErrorMsg)
	}

	c.logger.Info("Organization role binding granted successfully", map[string]interface{}{
		"cmd": "org grant",
		"org": req.OrgName,
	})
	return nil
}

// RevokeRoleBinding revokes a role from a user or user group on an organization.
// Handles authentication, organization validation, and provides detailed error messages.
func (c *OrgHttpClient) RevokeRoleBinding(req types.RevokeRoleBindingInput) error {
	// Validate authentication
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return err
	}

	// Get organization details for role binding revocation
	org, err := c.GetOrganizationByName(req.OrgName)
	if err != nil {
		c.logger.Error("Failed to get organization during role binding revocation", err, map[string]interface{}{
			"cmd": "org revoke",
			"org": req.OrgName,
		})
		return err
	}

	// Prepare role binding revocation parameters
	body, err := json.Marshal(RevokeRoleBindingRequest{
		UserName:       req.UserName,
		GroupName:      req.GroupName,
		Role:           req.Role,
		ResourceType:   "organization",
		ResourceID:     org.ID.String(),
		OrganizationID: org.ID.String(),
	})
	if err != nil {
		c.logger.Error("Failed to marshal role binding revocation request", err, map[string]interface{}{
			"cmd": "org revoke",
			"org": req.OrgName,
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
		c.logger.Error("Failed to revoke organization role binding", err, map[string]interface{}{
			"cmd": "org revoke",
			"org": req.OrgName,
		})
		return err
	}

	// Handle response
	if !respBody.Success {
		if respBody.ErrorCode == "role_binding_not_found" {
			return cliErrors.ErrRoleBindingNotFound
		}
		c.logger.Error("Failed to revoke organization role binding", fmt.Errorf(respBody.ErrorMsg), map[string]interface{}{
			"cmd": "org revoke",
			"org": req.OrgName,
		})
		return fmt.Errorf(respBody.ErrorMsg)
	}

	c.logger.Info("Organization role binding revoked successfully", map[string]interface{}{
		"cmd": "org revoke",
		"org": req.OrgName,
	})
	return nil
}
