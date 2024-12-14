package v2

// Parent represents the parent of a Jira issue, to be used with subtask issue types.
type IssueParent struct {
	ID  string `json:"id,omitempty" structs:"id"`
	Key string `json:"key,omitempty" structs:"key"`
}
