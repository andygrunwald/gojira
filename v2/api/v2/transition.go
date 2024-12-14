package v2

//go:generate genny -tag generated -in=../modelhelpers/unmarshalling_value_generics.go -out=transition_generated.go gen "GenericType=Transition"

// Transition represents an issue transition in Jira
type Transition struct {
	ID     string                     `json:"id" structs:"id"`
	Name   string                     `json:"name" structs:"name"`
	To     Status                     `json:"to" structs:"status"`
	Fields map[string]TransitionField `json:"fields" structs:"fields"`
}

// Unmarshal unmarshals JSON s to Changelog.
func UnmarshalTransition(j string) (Transition, error) {
	var t Transition

	return t, nil
}
