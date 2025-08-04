package groups

import (
	"encoding/json"
	"fmt"

	"github.com/Gkemhcs/kavach-cli/internal/auth"
	"github.com/Gkemhcs/kavach-cli/internal/client"
	"github.com/Gkemhcs/kavach-cli/internal/config"
	cliErrors "github.com/Gkemhcs/kavach-cli/internal/errors"
	"github.com/Gkemhcs/kavach-cli/internal/org"
	"github.com/Gkemhcs/kavach-cli/internal/utils"
)

// NewUsergroupHttpClient creates a new UserGroupHttpClient instance with the provided dependencies.
// This client handles HTTP communication with the backend for user group operations.
func NewUsergroupHttpClient(orgGetter org.OrgGetterClient, logger *utils.Logger, cfg *config.Config) *UserGroupHttpClient {
	return &UserGroupHttpClient{
		orgGetter,
		logger,
		cfg,
	}
}

// UserGroupHttpClient implements UserGroupClient for making HTTP requests to the backend for user group operations.
// It provides methods for creating, managing, and listing user groups within organizations.
type UserGroupHttpClient struct {
	orgGetter org.OrgGetterClient
	logger    *utils.Logger
	cfg       *config.Config
}

// CreateUserGroup creates a new user group within an organization.
// Handles authentication, organization validation, and error reporting with detailed logging.
func (c *UserGroupHttpClient) CreateUserGroup(orgName, userGroupName, description string) error {
	// Validate authentication
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return err
	}

	// Get organization details
	org, err := c.orgGetter.GetOrganizationByName(orgName)
	if err != nil {
		c.logger.Error("Failed to get organization during user group creation", err, map[string]interface{}{
			"cmd":       "user-group create",
			"userGroup": userGroupName,
			"org":       orgName,
		})
		return err
	}

	// Prepare request payload
	body, _ := json.Marshal(CreateUserGroupRequest{
		GroupName:   userGroupName,
		Description: description,
	})

	payload := client.RequestPayload{
		Method: "POST",
		URL:    fmt.Sprintf("%sorganizations/%s/user-groups/", c.cfg.BackendEndpoint, org.ID),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: body,
	}

	// Execute the request
	respBody, err := client.DoAuthenticatedRequest[any](payload, c.logger, c.cfg)
	if err != nil {
		if err == cliErrors.ErrNotLoggedIn {
			c.logger.Print("🔒 You are not logged in. Please run `kavach login` to authenticate.")
			c.logger.Warn("User not logged in during user group creation", map[string]interface{}{
				"cmd":       "user-group create",
				"userGroup": userGroupName,
				"org":       orgName,
			})
			return nil
		}
		c.logger.Error("Failed to create user group", err, map[string]interface{}{
			"cmd":       "user-group create",
			"userGroup": userGroupName,
			"org":       orgName,
		})
		return err
	}

	// Handle response
	if !respBody.Success {
		if respBody.ErrorCode == "duplicate_user_group" {
			c.logger.Print("⚠️  A user group with this name already exists in the organization. Please choose a different name.")
			c.logger.Warn("Duplicate user group during creation", map[string]interface{}{
				"cmd":       "user-group create",
				"userGroup": userGroupName,
				"org":       orgName,
			})
			return cliErrors.ErrDuplicateUserGroup
		}
		c.logger.Error("Failed to create user group", fmt.Errorf(respBody.ErrorMsg), map[string]interface{}{
			"cmd":       "user-group create",
			"userGroup": userGroupName,
			"org":       orgName,
		})
		return fmt.Errorf(respBody.ErrorMsg)
	}

	c.logger.Info("User group created successfully", map[string]interface{}{
		"cmd":       "user-group create",
		"userGroup": userGroupName,
		"org":       orgName,
	})
	return nil
}

// GetUserGroupByName retrieves a user group by its name within an organization.
// Handles authentication, organization validation, and provides detailed error messages.
func (c *UserGroupHttpClient) GetUserGroupByName(orgName, userGroupName string) (*UserGroup, error) {
	// Validate authentication
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return nil, err
	}

	// Get organization details
	org, err := c.orgGetter.GetOrganizationByName(orgName)
	if err != nil {
		c.logger.Error("Failed to get organization during user group retrieval", err, map[string]interface{}{
			"cmd":       "user-group get",
			"userGroup": userGroupName,
			"org":       orgName,
		})
		return nil, err
	}

	// Prepare request payload
	payload := client.RequestPayload{
		Method: "GET",
		URL:    fmt.Sprintf("%sorganizations/%s/user-groups/by-name", c.cfg.BackendEndpoint, org.ID),
		QueryParams: map[string]string{
			"name": userGroupName,
		},
	}

	// Execute the request
	respBody, err := client.DoAuthenticatedRequest[UserGroup](payload, c.logger, c.cfg)
	if err != nil {
		if err == cliErrors.ErrNotLoggedIn {
			c.logger.Print("🔒 You are not logged in. Please run `kavach login` to authenticate.")
			c.logger.Warn("User not logged in during user group retrieval", map[string]interface{}{
				"cmd":       "user-group get",
				"userGroup": userGroupName,
				"org":       orgName,
			})
			return nil, cliErrors.ErrNotLoggedIn
		}
		c.logger.Error("Failed to get user group by name", err, map[string]interface{}{
			"cmd":       "user-group get",
			"userGroup": userGroupName,
			"org":       orgName,
		})
		return nil, err
	}

	// Handle response
	if !respBody.Success {
		if respBody.ErrorCode == "user_group_not_exist" {
			c.logger.Print("🚨 User group not found in the organization. Please verify the group name.")
			return nil, cliErrors.ErrUserGroupNotFound
		}
		c.logger.Error("Failed to get user group by name", fmt.Errorf(respBody.ErrorMsg), map[string]interface{}{
			"cmd":       "user-group get",
			"userGroup": userGroupName,
			"org":       orgName,
		})
		return nil, fmt.Errorf(respBody.ErrorMsg)
	}

	return &respBody.Data, nil
}

// DeleteUserGroup removes a user group from an organization.
// Handles authentication, organization validation, and provides detailed error messages.
func (c *UserGroupHttpClient) DeleteUserGroup(orgName, userGroupName string) error {
	// Validate authentication
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return err
	}

	// Get organization details
	org, err := c.orgGetter.GetOrganizationByName(orgName)
	if err != nil {
		c.logger.Error("Failed to get organization during user group deletion", err, map[string]interface{}{
			"cmd":       "user-group delete",
			"userGroup": userGroupName,
			"org":       orgName,
		})
		return err
	}

	// Get user group details for deletion
	userGroup, err := c.GetUserGroupByName(orgName, userGroupName)
	if err != nil {
		c.logger.Error("Failed to get user group during deletion", err, map[string]interface{}{
			"cmd":       "user-group delete",
			"userGroup": userGroupName,
			"org":       orgName,
		})
		return err
	}

	// Prepare request payload
	payload := client.RequestPayload{
		Method: "DELETE",
		URL:    fmt.Sprintf("%sorganizations/%s/user-groups/%s", c.cfg.BackendEndpoint, org.ID, userGroup.ID),
	}

	// Execute the request
	respBody, err := client.DoAuthenticatedRequest[any](payload, c.logger, c.cfg)
	if err != nil {
		if err == cliErrors.ErrNotLoggedIn {
			c.logger.Print("🔒 You are not logged in. Please run `kavach login` to authenticate.")
			c.logger.Warn("User not logged in during user group deletion", map[string]interface{}{
				"cmd":       "user-group delete",
				"userGroup": userGroupName,
				"org":       orgName,
			})
			return cliErrors.ErrNotLoggedIn
		}
		c.logger.Error("Failed to delete user group", err, map[string]interface{}{
			"cmd":       "user-group delete",
			"userGroup": userGroupName,
			"org":       orgName,
		})
		return err
	}

	// Handle response
	if !respBody.Success {
		if respBody.ErrorCode == "user_group_not_exist" {
			c.logger.Print("🚨 User group not found in the organization. Please verify the group name.")
			return cliErrors.ErrUserGroupNotFound
		}
	}

	c.logger.Info("User group deleted successfully", map[string]interface{}{
		"cmd":       "user-group delete",
		"userGroup": userGroupName,
		"org":       orgName,
	})
	return nil
}

// ListUserGroups retrieves all user groups within an organization.
// Handles authentication, organization validation, and provides detailed error messages.
func (c *UserGroupHttpClient) ListUserGroups(orgName string) ([]ListGroupsByOrgRow, error) {
	// Validate authentication
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return nil, err
	}

	// Get organization details
	org, err := c.orgGetter.GetOrganizationByName(orgName)
	if err != nil {
		c.logger.Error("Failed to get organization during user group listing", err, map[string]interface{}{
			"cmd": "user-group list",
			"org": orgName,
		})
		return nil, err
	}

	// Prepare request payload
	payload := client.RequestPayload{
		Method: "GET",
		URL:    fmt.Sprintf("%sorganizations/%s/user-groups/", c.cfg.BackendEndpoint, org.ID),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	// Execute the request
	respBody, err := client.DoAuthenticatedRequest[[]ListGroupsByOrgRow](payload, c.logger, c.cfg)
	if err != nil {
		if err == cliErrors.ErrNotLoggedIn {
			c.logger.Print("🔒 You are not logged in. Please run `kavach login` to authenticate.")
			c.logger.Warn("User not logged in during user group listing", map[string]interface{}{
				"cmd": "user-group list",
				"org": orgName,
			})
			return nil, cliErrors.ErrNotLoggedIn
		}
		c.logger.Error("Failed to list user groups", err, map[string]interface{}{
			"cmd": "user-group list",
			"org": orgName,
		})
		return nil, err
	}

	// Handle response
	if !respBody.Success {
		return nil, fmt.Errorf("API error: %s", respBody.ErrorMsg)
	}

	c.logger.Info("User groups listed successfully", map[string]interface{}{
		"cmd":   "user-group list",
		"count": len(respBody.Data),
		"org":   orgName,
	})
	return respBody.Data, nil
}

// AddGroupMember adds a user to a user group.
// Handles authentication, organization validation, and provides detailed error messages.
func (c *UserGroupHttpClient) AddGroupMember(orgName, userName, userGroupName string) error {
	// Validate authentication
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return err
	}

	// Get organization details
	org, err := c.orgGetter.GetOrganizationByName(orgName)
	if err != nil {
		c.logger.Error("Failed to get organization during member addition", err, map[string]interface{}{
			"cmd":       "user-group members add",
			"userGroup": userGroupName,
			"user":      userName,
			"org":       orgName,
		})
		return err
	}

	// Get user group details
	userGroup, err := c.GetUserGroupByName(orgName, userGroupName)
	if err != nil {
		c.logger.Error("Failed to get user group during member addition", err, map[string]interface{}{
			"cmd":       "user-group members add",
			"userGroup": userGroupName,
			"user":      userName,
			"org":       orgName,
		})
		return err
	}

	// Prepare request payload
	body, _ := json.Marshal(AddMemberRequest{
		UserName: userName,
	})

	payload := client.RequestPayload{
		Method: "POST",
		URL:    fmt.Sprintf("%sorganizations/%s/user-groups/%s/members", c.cfg.BackendEndpoint, org.ID, userGroup.ID),
		Body:   body,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	// Execute the request
	respBody, err := client.DoAuthenticatedRequest[any](payload, c.logger, c.cfg)
	if err != nil {
		if err == cliErrors.ErrNotLoggedIn {
			c.logger.Print("🔒 You are not logged in. Please run `kavach login` to authenticate.")
			c.logger.Warn("User not logged in during member addition", map[string]interface{}{
				"cmd":       "user-group members add",
				"userGroup": userGroupName,
				"user":      userName,
				"org":       orgName,
			})
			return cliErrors.ErrNotLoggedIn
		}
		c.logger.Error("Failed to add user group member", err, map[string]interface{}{
			"cmd":       "user-group members add",
			"userGroup": userGroupName,
			"user":      userName,
			"org":       orgName,
		})
		return err
	}

	// Handle response
	if !respBody.Success {
		if respBody.ErrorCode == "duplicate_user_group_membership" {
			c.logger.Print("⚠️  User is already a member of this group.")
			c.logger.Warn("Duplicate user group membership during addition", map[string]interface{}{
				"cmd":       "user-group members add",
				"userGroup": userGroupName,
				"user":      userName,
				"org":       orgName,
			})
			return cliErrors.ErrDuplicateMember
		}
		if respBody.ErrorCode == "user_not_exist" {
			c.logger.Print("🚨 User not found. Please verify the GitHub username.")
			return cliErrors.ErrUserNotFound
		}
		return fmt.Errorf("API error: %s", respBody.ErrorMsg)
	}

	c.logger.Info("User group member added successfully", map[string]interface{}{
		"cmd":       "user-group members add",
		"userGroup": userGroupName,
		"user":      userName,
		"org":       orgName,
	})
	return nil
}

// RemoveGroupMember removes a user from a user group.
// Handles authentication, organization validation, and provides detailed error messages.
func (c *UserGroupHttpClient) RemoveGroupMember(orgName, userName, userGroupName string) error {
	// Validate authentication
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return err
	}

	// Get organization details
	org, err := c.orgGetter.GetOrganizationByName(orgName)
	if err != nil {
		c.logger.Error("Failed to get organization during member removal", err, map[string]interface{}{
			"cmd":       "user-group members remove",
			"userGroup": userGroupName,
			"user":      userName,
			"org":       orgName,
		})
		return err
	}

	// Get user group details
	userGroup, err := c.GetUserGroupByName(orgName, userGroupName)
	if err != nil {
		c.logger.Error("Failed to get user group during member removal", err, map[string]interface{}{
			"cmd":       "user-group members remove",
			"userGroup": userGroupName,
			"user":      userName,
			"org":       orgName,
		})
		return err
	}

	// Prepare request payload
	body, _ := json.Marshal(RemoveMemberRequest{
		UserName: userName,
	})

	payload := client.RequestPayload{
		Method: "DELETE",
		URL:    fmt.Sprintf("%sorganizations/%s/user-groups/%s/members", c.cfg.BackendEndpoint, org.ID, userGroup.ID),
		Body:   body,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	// Execute the request
	respBody, err := client.DoAuthenticatedRequest[any](payload, c.logger, c.cfg)
	if err != nil {
		if err == cliErrors.ErrNotLoggedIn {
			c.logger.Print("🔒 You are not logged in. Please run `kavach login` to authenticate.")
			c.logger.Warn("User not logged in during member removal", map[string]interface{}{
				"cmd":       "user-group members remove",
				"userGroup": userGroupName,
				"user":      userName,
				"org":       orgName,
			})
			return cliErrors.ErrNotLoggedIn
		}
		c.logger.Error("Failed to remove user group member", err, map[string]interface{}{
			"cmd":       "user-group members remove",
			"userGroup": userGroupName,
			"user":      userName,
			"org":       orgName,
		})
		return err
	}

	// Handle response
	if !respBody.Success {
		if respBody.ErrorCode == "user_membership_not_exist" {
			c.logger.Print("⚠️  User is not a member of this group.")
			return cliErrors.ErrMemberNotFound
		}
		if respBody.ErrorCode == "user_not_exist" {
			c.logger.Print("🚨 User not found. Please verify the GitHub username.")
			return cliErrors.ErrUserNotFound
		}
		return fmt.Errorf("API error: %s", respBody.ErrorMsg)
	}

	c.logger.Info("User group member removed successfully", map[string]interface{}{
		"cmd":       "user-group members remove",
		"userGroup": userGroupName,
		"user":      userName,
		"org":       orgName,
	})
	return nil
}

// ListUserGroupMembers retrieves all members of a user group.
// Handles authentication, organization validation, and provides detailed error messages.
func (c *UserGroupHttpClient) ListUserGroupMembers(orgName, userGroupName string) ([]ListGroupMembersRow, error) {
	// Validate authentication
	if err := auth.RequireAuthMiddleware(c.logger, c.cfg); err != nil {
		return nil, err
	}

	// Get organization details
	org, err := c.orgGetter.GetOrganizationByName(orgName)
	if err != nil {
		c.logger.Error("Failed to get organization during member listing", err, map[string]interface{}{
			"cmd":       "user-group members list",
			"userGroup": userGroupName,
			"org":       orgName,
		})
		return nil, err
	}

	// Get user group details
	userGroup, err := c.GetUserGroupByName(orgName, userGroupName)
	if err != nil {
		c.logger.Error("Failed to get user group during member listing", err, map[string]interface{}{
			"cmd":       "user-group members list",
			"userGroup": userGroupName,
			"org":       orgName,
		})
		return nil, err
	}

	// Prepare request payload
	payload := client.RequestPayload{
		Method: "GET",
		URL:    fmt.Sprintf("%sorganizations/%s/user-groups/%s/members", c.cfg.BackendEndpoint, org.ID, userGroup.ID),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	// Execute the request
	respBody, err := client.DoAuthenticatedRequest[[]ListGroupMembersRow](payload, c.logger, c.cfg)
	if err != nil {
		if err == cliErrors.ErrNotLoggedIn {
			c.logger.Print("🔒 You are not logged in. Please run `kavach login` to authenticate.")
			c.logger.Warn("User not logged in during member listing", map[string]interface{}{
				"cmd":       "user-group members list",
				"userGroup": userGroupName,
				"org":       orgName,
			})
			return nil, cliErrors.ErrNotLoggedIn
		}
		c.logger.Error("Failed to list user group members", err, map[string]interface{}{
			"cmd":       "user-group members list",
			"userGroup": userGroupName,
			"org":       orgName,
		})
		return nil, err
	}

	// Handle response
	if !respBody.Success {
		return nil, fmt.Errorf("API error: %s", respBody.ErrorMsg)
	}

	c.logger.Info("User group members listed successfully", map[string]interface{}{
		"cmd":       "user-group members list",
		"count":     len(respBody.Data),
		"userGroup": userGroupName,
		"org":       orgName,
	})
	return respBody.Data, nil
}
