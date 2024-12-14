package v2

// RemoteLinkIcon represents icon displayed next to link
type RemoteLinkIcon struct {
	Url16x16 string `json:"url16x16,omitempty" structs:"url16x16,omitempty"`
	Title    string `json:"title,omitempty" structs:"title,omitempty"`
	Link     string `json:"link,omitempty" structs:"link,omitempty"`
}
