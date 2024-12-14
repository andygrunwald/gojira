package v2

// RemoteLinkApplication represents remote links application
type RemoteLinkApplication struct {
	Type string `json:"type,omitempty" structs:"type,omitempty"`
	Name string `json:"name,omitempty" structs:"name,omitempty"`
}
