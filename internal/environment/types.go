package environment

import "github.com/Gkemhcs/kavach-cli/internal/types"

// CreateEnvironmentRequest is the request body for creating an environment.
type CreateEnvironmentRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type EnvironmentClient interface {
	CreateEnvironment(environmentName, secretGroupName, orgName, description string) error
	ListEnvironment(orgName, groupName string) ([]ListEnvironmentsWithMemberRow, error)
	GetEnvironmentbyName(orgName, groupName, environmentName string) (*EnvironmentResponseData, error)
	DeleteEnvironment(orgName, groupName, envName string) error
	GrantRoleBinding(req types.GrantRoleBindingInput) error
	RevokeRoleBinding(req types.RevokeRoleBindingInput) error
}

type EnvironmentGetter interface {
	GetEnvironmentbyName(orgName, groupName, environmentName string) (*EnvironmentResponseData, error)
}

type EnvironmentResponseData struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	SecretGroupID  string  `json:"secret_group_id"`
	Description    *string `json:"description"`
	OrganizationID string  `json:"organization_id"`
}

type ListEnvironmentsWithMemberRow struct {
	EnvironmentID   string `json:"environment_id"`
	Name            string `json:"name"`
	SecretGroupName string `json:"secret_group_name"`
	Role            string `json:"role"`
}
