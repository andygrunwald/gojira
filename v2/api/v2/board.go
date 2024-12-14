package v2

// Board represents a Jira agile board
type Board struct {
	ID       int    `json:"id,omitempty" structs:"id,omitempty"`
	Self     string `json:"self,omitempty" structs:"self,omitempty"`
	Name     string `json:"name,omitempty" structs:"name,omitemtpy"`
	Type     string `json:"type,omitempty" structs:"type,omitempty"`
	FilterID int    `json:"filterId,omitempty" structs:"filterId,omitempty"`
}
