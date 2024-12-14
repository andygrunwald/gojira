package v2

// RemoteLinkStatus if the link is a resolvable object (issue, epic) - the structure represent its status
type RemoteLinkStatus struct {
	Resolved bool
	Icon     *RemoteLinkIcon
}
