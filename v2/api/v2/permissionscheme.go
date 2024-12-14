package v2

// PermissionScheme represents the permission scheme for the project
type PermissionScheme struct {
	Expand      string       `json:"expand" structs:"expand,omitempty"`
	Self        string       `json:"self" structs:"self,omitempty"`
	ID          int          `json:"id" structs:"id,omitempty"`
	Name        string       `json:"name" structs:"name,omitempty"`
	Description string       `json:"description" structs:"description,omitempty"`
	Permissions []Permission `json:"permissions" structs:"permissions,omitempty"`
}
