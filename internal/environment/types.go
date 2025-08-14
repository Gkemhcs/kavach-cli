package environment

import "github.com/Gkemhcs/kavach-cli/internal/types"

// CreateEnvironmentRequest is the request body for creating an environment.
type CreateEnvironmentRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// EnvironmentClient defines the interface for environment operations
type EnvironmentClient interface {
	CreateEnvironment(environmentName, secretGroupName, orgName, description string) error
	ListEnvironment(orgName, groupName string) ([]ListEnvironmentsWithMemberRow, error)
	GetEnvironmentbyName(orgName, groupName, environmentName string) (*EnvironmentResponseData, error)
	DeleteEnvironment(orgName, groupName, envName string) error
	GrantRoleBinding(req types.GrantRoleBindingInput) error
	RevokeRoleBinding(req types.RevokeRoleBindingInput) error
	ListRoleBindings(orgName, groupName, envName string) ([]RoleBinding, error)
}

// EnvironmentGetter defines the interface for retrieving environment data
type EnvironmentGetter interface {
	GetEnvironmentbyName(orgName, groupName, environmentName string) (*EnvironmentResponseData, error)
}

// EnvironmentResponseData represents the response data for environment operations
type EnvironmentResponseData struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	SecretGroupID  string  `json:"secret_group_id"`
	Description    *string `json:"description"`
	OrganizationID string  `json:"organization_id"`
}

// ListEnvironmentsWithMemberRow represents environment data with member information
type ListEnvironmentsWithMemberRow struct {
	EnvironmentID   string `json:"environment_id"`
	Name            string `json:"name"`
	SecretGroupName string `json:"secret_group_name"`
	Role            string `json:"role"`
	InheritedFrom   string `json:"inherited_from"` // Source of the role (environment, secret_group, none)
}

// RoleBinding represents a single role binding with resolved names
type RoleBinding struct {
	EntityID       string `json:"entity_id"`
	EntityName     string `json:"entity_name"`
	EntityType     string `json:"entity_type"` // "user" or "group"
	GroupID        string `json:"group_id"`
	GroupName      string `json:"group_name"`
	Role           string `json:"role"`
	BindingType    string `json:"binding_type"` // "direct" or "inherited"
	SourceType     string `json:"source_type"`  // "organization", "secret_group", or "environment"
	InheritedFrom  string `json:"inherited_from"`
	OrganizationID string `json:"organization_id"`
	SecretGroupID  string `json:"secret_group_id"`
	EnvironmentID  string `json:"environment_id"`
}

// ListRoleBindingsResponse represents the response for listing role bindings
type ListRoleBindingsResponse struct {
	OrganizationID string        `json:"organization_id"`
	SecretGroupID  string        `json:"secret_group_id"`
	EnvironmentID  string        `json:"environment_id"`
	Bindings       []RoleBinding `json:"bindings"`
	Count          int           `json:"count"`
}
