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
}
