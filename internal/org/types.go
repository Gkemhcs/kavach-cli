package org

import (
	"github.com/Gkemhcs/kavach-cli/internal/types"
	"github.com/google/uuid"
)

type OrgClient interface {
	CreateOrganization(name, description string) error
	ListMyOrganizations() ([]ListMembersOfOrganizationRow, error)
	DeleteOrganization(name string) error
	GetOrganizationByName(name string) (*Organization, error)
	RevokeRoleBinding(req types.RevokeRoleBindingInput) error
	GrantRoleBinding(req types.GrantRoleBindingInput) error
}

type OrgGetterClient interface {
	GetOrganizationByName(name string) (*Organization, error)
}

// Organization represents an organization entity returned from the backend.
type Organization struct {
	ID          uuid.UUID `json:"id"`          // Organization unique ID
	Name        string    `json:"name"`        // Organization name
	Description string    `json:"description"` // Organization description
}

// CreateOrgRequest represents the request body for creating an organization.
type CreateOrgRequest struct {
	Name        string `json:"name"`        // Organization name
	Description string `json:"description"` // Organization description
}

// ListMembersOfOrganizationRow represents a row in the organization membership list.
type ListMembersOfOrganizationRow struct {
	OrgID uuid.UUID `json:"id"`       // Organization unique ID
	Name  string    `json:"org_name"` // Organization name
	Role  string    `json:"role"`     // User's role in the organization
}

type APIResponse[T any] struct {
	Success   bool   `json:"success"`
	Data      T      `json:"data,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
	ErrorMsg  string `json:"error_msg,omitempty"`
}

// Update OrgResponse to match backend JSON
// {"org_id":"...","name":"...","created_at":"...","user_id":"...","role":"..."}
type OrgResponse struct {
	OrgID     string `json:"org_id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UserID    string `json:"user_id"`
	Role      string `json:"role"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type GrantRoleBindingRequest struct {
	UserName       string `json:"user_name"`
	GroupName      string `json:"group_name"`
	Role           string `json:"role"`
	ResourceType   string `json:"resource_type"`
	ResourceID     string `json:"resource_id"`
	OrganizationID string `json:"organization_id"`
}

type RevokeRoleBindingRequest struct {
	UserName       string `json:"user_name"`
	GroupName      string `json:"group_name"`
	Role           string `json:"role"`
	ResourceType   string `json:"resource_type"`
	ResourceID     string `json:"resource_id"`
	OrganizationID string `json:"organization_id"`
}
