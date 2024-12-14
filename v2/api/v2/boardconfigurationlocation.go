package v2

// BoardConfigurationLocation reference to the container that the board is located in
type BoardConfigurationLocation struct {
	Type string `json:"type"`
	Key  string `json:"key"`
	ID   string `json:"id"`
	Self string `json:"self"`
	Name string `json:"name"`
}
