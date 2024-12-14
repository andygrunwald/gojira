package v2

// IssueComponent represents a "component" of a Jira issue.
// Components can be user defined in every Jira instance.
type IssueComponent struct {
	Self string `json:"self,omitempty" structs:"self,omitempty"`
	ID   string `json:"id,omitempty" structs:"id,omitempty"`
	Name string `json:"name,omitempty" structs:"name,omitempty"`
}
