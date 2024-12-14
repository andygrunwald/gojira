package v2

type IssueFieldSchema struct {
	Type   string `json:"type,omitempty" structs:"type,omitempty"`
	System string `json:"system,omitempty" structs:"system,omitempty"`
}
