package types

type GrantRoleBindingRequest struct {
	UserName       string  `json:"user_name"`
	GroupName      string  `json:"group_name"`
	Role           string  `json:"role"`
	ResourceType   string  `json:"resource_type"`
	ResourceID     string  `json:"resource_id"`
	OrganizationID string  `json:"organization_id"`
	SecretGroupID  *string `json:"secret_group_id"`
	EnvironmentID  *string `json:"environment_id"`
}

type RevokeRoleBindingRequest struct {
	UserName       string  `json:"user_name"`
	GroupName      string  `json:"group_name"`
	Role           string  `json:"role"`
	ResourceType   string  `json:"resource_type"`
	ResourceID     string  `json:"resource_id"`
	OrganizationID string  `json:"organization_id"`
	SecretGroupID  *string `json:"secret_group_id"`
	EnvironmentID  *string `json:"environment_id"`
}

type GrantRoleBindingInput struct {
	UserName        string
	GroupName       string
	Role            string
	OrgName         string
	SecretGroupName string
	EnvironmentName string
}

type RevokeRoleBindingInput struct {
	UserName        string
	GroupName       string
	Role            string
	OrgName         string
	SecretGroupName string
	EnvironmentName string
}

// Secret-related types
type Secret struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type CreateSecretVersionRequest struct {
	Secrets       []Secret `json:"secrets"`
	CommitMessage string   `json:"commit_message"`
}

type SecretVersionResponse struct {
	ID            string `json:"id"`
	EnvironmentID string `json:"environment_id"`
	CommitMessage string `json:"commit_message"`
	CreatedAt     string `json:"created_at"`
	SecretCount   int    `json:"secret_count"`
}

type SecretVersionDetailResponse struct {
	ID            string   `json:"id"`
	EnvironmentID string   `json:"environment_id"`
	CommitMessage string   `json:"commit_message"`
	CreatedAt     string   `json:"created_at"`
	Secrets       []Secret `json:"secrets"`
}

type SecretDiffChange struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	OldValue string `json:"old_value,omitempty"`
	NewValue string `json:"new_value,omitempty"`
}

type SecretDiffResponse struct {
	FromVersion string             `json:"from_version"`
	ToVersion   string             `json:"to_version"`
	Changes     []SecretDiffChange `json:"changes"`
}

type RollbackRequest struct {
	VersionID     string `json:"version_id"`
	CommitMessage string `json:"commit_message"`
}

// Sync-related types
type SyncResult struct {
	Name    string `json:"name"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type SyncSecretsResponse struct {
	Provider    string       `json:"provider"`
	Status      string       `json:"status"`
	Message     string       `json:"message"`
	SyncedCount int          `json:"synced_count"`
	FailedCount int          `json:"failed_count"`
	TotalCount  int          `json:"total_count"`
	Results     []SyncResult `json:"results"`
	Errors      []string     `json:"errors,omitempty"`
	SyncedAt    string       `json:"synced_at"`
}

// Provider-related types
type ProviderCredentialResponse struct {
	ID            string                 `json:"id"`
	EnvironmentID string                 `json:"environment_id"`
	Provider      string                 `json:"provider"`
	Credentials   map[string]interface{} `json:"credentials,omitempty"`
	Config        map[string]interface{} `json:"config"`
	CreatedAt     string                 `json:"created_at"`
	UpdatedAt     string                 `json:"updated_at"`
}
