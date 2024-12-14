package v2

// SprintsList reflects a list of agile sprints
type SprintsList struct {
	MaxResults int      `json:"maxResults" structs:"maxResults"`
	StartAt    int      `json:"startAt" structs:"startAt"`
	Total      int      `json:"total" structs:"total"`
	IsLast     bool     `json:"isLast" structs:"isLast"`
	Values     []Sprint `json:"values" structs:"values"`
}
