package v2

// RemoteLink represents remote links which linked to issues
type IssueRemoteLink struct {
	ID           int                    `json:"id,omitempty" structs:"id,omitempty"`
	Self         string                 `json:"self,omitempty" structs:"self,omitempty"`
	GlobalID     string                 `json:"globalId,omitempty" structs:"globalId,omitempty"`
	Application  *RemoteLinkApplication `json:"application,omitempty" structs:"application,omitempty"`
	Relationship string                 `json:"relationship,omitempty" structs:"relationship,omitempty"`
	Object       *RemoteLinkObject      `json:"object,omitempty" structs:"object,omitempty"`
}
