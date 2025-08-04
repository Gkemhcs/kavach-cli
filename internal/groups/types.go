package groups

// UserGroupClient defines the interface for user group management operations.
// This interface provides methods for creating, retrieving, deleting, and listing user groups.
type UserGroupClient interface {
	// CreateUserGroup creates a new user group within an organization.
	// Returns an error if the group creation fails or if a group with the same name already exists.
	CreateUserGroup(orgName, userGroupName, description string) error

	// GetUserGroupByName retrieves a user group by its name within an organization.
	// Returns the user group data or an error if the group is not found.
	GetUserGroupByName(orgName, userGroupName string) (*UserGroup, error)

	// DeleteUserGroup removes a user group from an organization.
	// Returns an error if the group doesn't exist or if deletion fails.
	DeleteUserGroup(orgName, userGroupName string) error

	// ListUserGroups retrieves all user groups within an organization.
	// Returns a list of user groups or an error if the operation fails.
	ListUserGroups(orgName string) ([]ListGroupsByOrgRow, error)
}

// UserGroupMemberClient defines the interface for user group member management operations.
// This interface provides methods for adding, removing, and listing members of user groups.
type UserGroupMemberClient interface {
	// GetUserGroupByName retrieves a user group by its name within an organization.
	// Returns the user group data or an error if the group is not found.
	GetUserGroupByName(orgName, userGroupName string) (*UserGroup, error)

	// AddGroupMember adds a user to a user group.
	// Returns an error if the user or group doesn't exist, or if the user is already a member.
	AddGroupMember(orgName, userName, userGroupName string) error

	// RemoveGroupMember removes a user from a user group.
	// Returns an error if the user or group doesn't exist, or if the user is not a member.
	RemoveGroupMember(orgName, userName, userGroupName string) error

	// ListUserGroupMembers retrieves all members of a user group.
	// Returns a list of group members or an error if the operation fails.
	ListUserGroupMembers(orgName, userGroupName string) ([]ListGroupMembersRow, error)
}

// CreateUserGroupRequest represents the request payload for creating a new user group.
// Contains the group name and optional description for the new user group.
type CreateUserGroupRequest struct {
	GroupName   string `json:"group_name"`  // Name of the user group (required, must be unique within org)
	Description string `json:"description"` // Optional description of the group's purpose
}

// UserGroup represents a user group entity with its metadata.
// Contains the complete information about a user group including its ID, name, and organization.
type UserGroup struct {
	ID             string `json:"id"`              // Unique identifier for the user group
	Name           string `json:"name"`            // Name of the user group
	OrganizationID string `json:"organization_id"` // ID of the organization this group belongs to
	Description    string `json:"description"`     // Optional description of the group
}

// ListGroupsByOrgRow represents a user group entry in a list view.
// Contains minimal required fields for display purposes in group listings.
type ListGroupsByOrgRow struct {
	ID          string `json:"id"`          // Unique identifier for the user group
	Name        string `json:"name"`        // Name of the user group
	Description string `json:"description"` // Optional description of the group
}

// AddMemberRequest represents the request payload for adding a member to a user group.
// Contains the username of the user to be added to the group.
type AddMemberRequest struct {
	UserName string `json:"user_name"` // GitHub username of the user to add to the group
}

// RemoveMemberRequest represents the request payload for removing a member from a user group.
// Contains the username of the user to be removed from the group.
type RemoveMemberRequest struct {
	UserName string `json:"user_name"` // GitHub username of the user to remove from the group
}

// ListGroupMembersRow represents a group member entry in a list view.
// Contains the user information for display purposes in member listings.
type ListGroupMembersRow struct {
	ID    string `json:"id"`    // Unique identifier for the user
	Name  string `json:"name"`  // GitHub username of the user
	Email string `json:"email"` // Email address of the user
}
