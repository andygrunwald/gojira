package v2

// Role represents a Jira product role
type Role struct {
	Self        string   `json:"self" structs:"self"`
	Name        string   `json:"name" structs:"name"`
	ID          int      `json:"id" structs:"id"`
	Description string   `json:"description" structs:"description"`
	Actors      []*Actor `json:"actors" structs:"actors"`
}
