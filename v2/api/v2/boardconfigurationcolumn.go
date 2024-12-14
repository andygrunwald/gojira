package v2

// BoardConfigurationColumn lists the name of the board with the statuses that maps to a particular column
type BoardConfigurationColumn struct {
	Name   string                           `json:"name"`
	Status []BoardConfigurationColumnStatus `json:"statuses"`
}
