package v2

// Filter represents a Filter in Jira
type Issueilter struct {
	Self             string        `json:"self"`
	ID               string        `json:"id"`
	Name             string        `json:"name"`
	Description      string        `json:"description"`
	Owner            User          `json:"owner"`
	Jql              string        `json:"jql"`
	ViewURL          string        `json:"viewUrl"`
	SearchURL        string        `json:"searchUrl"`
	Favourite        bool          `json:"favourite"`
	FavouritedCount  int           `json:"favouritedCount"`
	SharePermissions []interface{} `json:"sharePermissions"`
	Subscriptions    struct {
		Size       int           `json:"size"`
		Items      []interface{} `json:"items"`
		MaxResults int           `json:"max-results"`
		StartIndex int           `json:"start-index"`
		EndIndex   int           `json:"end-index"`
	} `json:"subscriptions"`
}
