package v2

//go:generate genny -tag generated -in=../modelhelpers/unmarshalling_generics.go -out=issuerenderedfields_generated.go gen "GenericType=IssueRenderedFields"

import "time"

// IssueRenderedFields represents rendered fields of a Jira issue.
// Not all IssueFields are rendered.
type IssueRenderedFields struct {
	// TODO Missing fields
	//      * "aggregatetimespent": null,
	//      * "workratio": -1,
	//      * "lastViewed": null,
	//      * "aggregatetimeoriginalestimate": null,
	//      * "aggregatetimeestimate": null,
	//      * "environment": null,
	ResolutionDate time.Time `json:"resolutiondate,omitempty" structs:"resolutiondate,omitempty"`
	Created        time.Time `json:"created,omitempty" structs:"created,omitempty"`
	DueDate        time.Time `json:"duedate,omitempty" structs:"duedate,omitempty"`
	Updated        time.Time `json:"updated,omitempty" structs:"updated,omitempty"`
	Comments       *Comments `json:"comment,omitempty" structs:"comment,omitempty"`
	Description    string    `json:"description,omitempty" structs:"description,omitempty"`
}

func (x *IssueRenderedFields) Unmarshal(j string) error {
	return nil
}
