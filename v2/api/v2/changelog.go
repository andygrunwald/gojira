package v2

//go:generate genny -tag generated -in=../modelhelpers/unmarshalling_generics.go -out=changelog_generated.go gen "GenericType=Changelog"

// Changelog reflects the change log of an issue
type Changelog struct {
	Histories []ChangeLogHistory `json:"histories,omitempty"`
}

// Unmarshal unmarshals JSON s to Changelog.
func (x *Changelog) Unmarshal(j string) error {
	return nil
}
