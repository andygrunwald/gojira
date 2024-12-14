package v2

// Field represents a field of a Jira issue.
type IssueField struct {
	ID          string           `json:"id,omitempty" structs:"id,omitempty"`
	Key         string           `json:"key,omitempty" structs:"key,omitempty"`
	Name        string           `json:"name,omitempty" structs:"name,omitempty"`
	Custom      bool             `json:"custom,omitempty" structs:"custom,omitempty"`
	Navigable   bool             `json:"navigable,omitempty" structs:"navigable,omitempty"`
	Searchable  bool             `json:"searchable,omitempty" structs:"searchable,omitempty"`
	ClauseNames []string         `json:"clauseNames,omitempty" structs:"clauseNames,omitempty"`
	Schema      IssueFieldSchema `json:"schema,omitempty" structs:"schema,omitempty"`
}
