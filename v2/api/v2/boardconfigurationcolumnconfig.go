package v2

// BoardConfigurationColumnConfig lists the columns for a given board in the order defined in the column configuration
// with constrainttype (none, issueCount, issueCountExclSubs)
type BoardConfigurationColumnConfig struct {
	Columns        []BoardConfigurationColumn `json:"columns"`
	ConstraintType string                     `json:"constraintType"`
}
