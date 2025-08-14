package secretgroup

import (
	"github.com/Gkemhcs/kavach-cli/internal/types"
	"github.com/google/uuid"
)

// CreateSecretGroupRequest represents the request body for creating a secret group.
type CreateSecretGroupRequest struct {
	Name        string `json:"name" binding:"required"` // Secret group name
	Description string `json:"description"`             // Secret group description
}

// SecretGroupResponseData represents the response data for a secret group.
type SecretGroupResponseData struct {
	ID             string `json:"id"`              // Secret group ID
	Name           string `json:"name"`            // Secret group name
	Description    string `json:"description"`     // Secret group description
	OrganizationID string `json:"organization_id"` // Organization ID
}

// SecretGroupClient defines the interface for secret group operations.
type SecretGroupClient interface {
	CreateSecretGroup(name, description, orgName string) error
	ListSecretGroups(orgName string) ([]ListSecretGroupsWithMemberRow, error)
	GetSecretGroupByName(orgName, secretName string) (*SecretGroupResponseData, error)
	DeleteSecretGroupByName(orgName, groupName string) error
	GrantRoleBinding(req types.GrantRoleBindingInput) error
	RevokeRoleBinding(req types.RevokeRoleBindingInput) error
	ListRoleBindings(orgName, groupName string) ([]RoleBinding, error)
}

// SecretGroupGetter defines the interface for retrieving a secret group by name.
type SecretGroupGetter interface {
	GetSecretGroupByName(orgName, secretName string) (*SecretGroupResponseData, error)
}

// ListSecretGroupsWithMemberRow represents a row in the secret group membership list.
type ListSecretGroupsWithMemberRow struct {
	SecretGroupID    uuid.UUID `json:"id"`                // Secret group unique ID
	Name             string    `json:"name"`              // Secret group name
	OrganizationName string    `json:"organization_name"` // Organization name
	Role             string    `json:"role"`              // User's role in the secret group
	InheritedFrom    string    `json:"inherited_from"`    // Source of the role (secret_group, organization, none)
}

// RoleBinding represents a role binding with resolved names
type RoleBinding struct {
	SecretGroupID string  `json:"secret_group_id"`
	Role          string  `json:"role"`
	BindingType   string  `json:"binding_type"` // "direct" or "inherited"
	EntityType    string  `json:"entity_type"`  // "user" or "group"
	EntityID      *string `json:"entity_id"`    // User ID (if entity_type is "user")
	EntityName    *string `json:"entity_name"`  // User name or group name
	GroupID       *string `json:"group_id"`     // Group ID (if entity_type is "group")
	GroupName     *string `json:"group_name"`   // Group name (if entity_type is "group")
	SourceType    string  `json:"source_type"`  // "secret_group" or "organization"
}

// ListRoleBindingsResponse represents the response for listing role bindings
type ListRoleBindingsResponse struct {
	OrganizationID string        `json:"organization_id"`
	SecretGroupID  string        `json:"secret_group_id"`
	Bindings       []RoleBinding `json:"bindings"`
	Count          int           `json:"count"`
}
