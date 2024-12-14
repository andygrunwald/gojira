package v2

// FiltersListItem represents a Filter of FiltersList in Jira
type IssueFiltersListItem struct {
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
	Subscriptions    []struct {
		ID   int  `json:"id"`
		User User `json:"user"`
	} `json:"subscriptions"`
}
