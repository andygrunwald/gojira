package v2

// RemoteLinkObject represents remote link object itself
type RemoteLinkObject struct {
	URL     string            `json:"url,omitempty" structs:"url,omitempty"`
	Title   string            `json:"title,omitempty" structs:"title,omitempty"`
	Summary string            `json:"summary,omitempty" structs:"summary,omitempty"`
	Icon    *RemoteLinkIcon   `json:"icon,omitempty" structs:"icon,omitempty"`
	Status  *RemoteLinkStatus `json:"status,omitempty" structs:"status,omitempty"`
}
