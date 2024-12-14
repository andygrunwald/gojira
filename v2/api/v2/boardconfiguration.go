package v2

// BoardConfiguration represents a boardConfiguration of a jira board
type BoardConfiguration struct {
	ID           int                            `json:"id"`
	Name         string                         `json:"name"`
	Self         string                         `json:"self"`
	Location     BoardConfigurationLocation     `json:"location"`
	Filter       BoardConfigurationFilter       `json:"filter"`
	SubQuery     BoardConfigurationSubQuery     `json:"subQuery"`
	ColumnConfig BoardConfigurationColumnConfig `json:"columnConfig"`
}
