package v2

//go:generate genny -tag generated -in=../modelhelpers/unmarshalling_generics.go -out=issuefields_generated.go gen "GenericType=IssueFields"

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fastjson"
	"time"
)

var issueFieldsJsonParserPool fastjson.ParserPool

// IssueFields represents single fields of a Jira issue.
// Every Jira issue has several fields attached.
type IssueFields struct {
	// TODO Missing fields
	//      * "workratio": -1,
	//      * "lastViewed": null,
	//      * "environment": null,
	Expand                        string                     `json:"expand,omitempty" structs:"expand,omitempty"`
	Type                          IssueType                  `json:"issuetype,omitempty" structs:"issuetype,omitempty"`
	Project                       Project                    `json:"project,omitempty" structs:"project,omitempty"`
	Resolution                    *Resolution                `json:"resolution,omitempty" structs:"resolution,omitempty"`
	Priority                      *Priority                  `json:"priority,omitempty" structs:"priority,omitempty"`
	ResolutionDate                time.Time                  `json:"resolutiondate,omitempty" structs:"resolutiondate,omitempty"`
	Created                       time.Time                  `json:"created,omitempty" structs:"created,omitempty"`
	DueDate                       time.Time                  `json:"duedate,omitempty" structs:"duedate,omitempty"`
	Watches                       *Watches                   `json:"watches,omitempty" structs:"watches,omitempty"`
	Assignee                      *User                      `json:"assignee,omitempty" structs:"assignee,omitempty"`
	Updated                       time.Time                  `json:"updated,omitempty" structs:"updated,omitempty"`
	Description                   string                     `json:"description,omitempty" structs:"description,omitempty"`
	Summary                       string                     `json:"summary,omitempty" structs:"summary,omitempty"`
	Creator                       *User                      `json:"Creator,omitempty" structs:"Creator,omitempty"`
	Reporter                      *User                      `json:"reporter,omitempty" structs:"reporter,omitempty"`
	Components                    []*Component               `json:"components,omitempty" structs:"components,omitempty"`
	Status                        *Status                    `json:"status,omitempty" structs:"status,omitempty"`
	Progress                      *Progress                  `json:"progress,omitempty" structs:"progress,omitempty"`
	AggregateProgress             *Progress                  `json:"aggregateprogress,omitempty" structs:"aggregateprogress,omitempty"`
	TimeTracking                  *TimeTracking              `json:"timetracking,omitempty" structs:"timetracking,omitempty"`
	TimeSpent                     time.Duration              `json:"timespent,omitempty" structs:"timespent,omitempty"`
	TimeEstimate                  int                        `json:"timeestimate,omitempty" structs:"timeestimate,omitempty"`
	TimeOriginalEstimate          int                        `json:"timeoriginalestimate,omitempty" structs:"timeoriginalestimate,omitempty"`
	Worklog                       *Worklog                   `json:"worklog,omitempty" structs:"worklog,omitempty"`
	IssueLinks                    []*IssueLink               `json:"issuelinks,omitempty" structs:"issuelinks,omitempty"`
	Comments                      *Comments                  `json:"comment,omitempty" structs:"comment,omitempty"`
	FixVersions                   []*FixVersion              `json:"fixVersions,omitempty" structs:"fixVersions,omitempty"`
	AffectsVersions               []*AffectsVersion          `json:"versions,omitempty" structs:"versions,omitempty"`
	Labels                        []string                   `json:"labels,omitempty" structs:"labels,omitempty"`
	Subtasks                      []*Subtasks                `json:"subtasks,omitempty" structs:"subtasks,omitempty"`
	Attachments                   []*Attachment              `json:"attachment,omitempty" structs:"attachment,omitempty"`
	Epic                          *Epic                      `json:"epic,omitempty" structs:"epic,omitempty"`
	Sprint                        *Sprint                    `json:"sprint,omitempty" structs:"sprint,omitempty"`
	Parent                        *Parent                    `json:"parent,omitempty" structs:"parent,omitempty"`
	AggregateTimeOriginalEstimate int                        `json:"aggregatetimeoriginalestimate,omitempty" structs:"aggregatetimeoriginalestimate,omitempty"`
	AggregateTimeSpent            int                        `json:"aggregatetimespent,omitempty" structs:"aggregatetimespent,omitempty"`
	AggregateTimeEstimate         int                        `json:"aggregatetimeestimate,omitempty" structs:"aggregatetimeestimate,omitempty"`
	CustomFields                  map[string]json.RawMessage `json:"-"`
}

func (i *IssueFields) Unmarshal(j string) error {
	return nil
}

func (i *IssueFields) GetCustomFieldByID(id string) interface{} {
	return i.CustomFields[fmt.Sprintf("customfield_%s", id)]
}

// https://developer.atlassian.com/server/jira/platform/jira-rest-api-examples/

// These comments are here to help users identify what type that some of the custom fields
// should ultimately be type asserted to
// There's little reason to actually add more types to the library though

// type CustomFieldDateTime time.Time
// type CustomFieldDatePicker time.Time
// type CustomFieldFreeText string
// type CustomFieldText string
// type CustomFieldURL string
// type CustomFieldNumberField float64

type CustomFieldSelect struct {
	Value string
}

type CustomFieldCascadingSelect struct {
	CustomFieldSelect
	Child CustomFieldSelect
}

type CustomFieldGroupPicker struct {
	Name string
	ID   *int64
}

type CustomFieldMultiGroupPicker []CustomFieldGroupPicker

type CustomFieldMultiSelect []CustomFieldSelect

type CustomFieldUserPicker struct {
	Name string
	ID   *int64
}

type CustomFieldMultiUserPicker []CustomFieldUserPicker

type CustomFieldProjectPicker struct {
	Key string
	ID  *int64
}

type CustomFieldRadioButtons struct {
	Value string
	ID    *int64
}

type CustomFieldSelectList struct {
	Value string
	ID    *int64
}

type CustomFieldSingleVersionPicker struct {
	Name string
}

type CustomFieldVersionPicker []CustomFieldSingleVersionPicker

// MarshalJSON is a custom JSON marshal function for the IssueFields structs.
// It handles Jira custom fields and maps those from / to "Unknowns" key.
func (i *IssueFields) MarshalJSON() ([]byte, error) {
	//type Alias_ IssueFields
	//return json.Marshal(&struct {
	//	*Alias_
	//	Stamp string `json:"stamp"`
	//}{
	//	Alias_: (*Alias_)(i),
	//	Stamp:  d.Stamp.Format("Mon Jan _2"),
	//})

	//m := structs.Map(i)
	//unknowns, okay := m["Unknowns"]
	//if okay {
	//	// if unknowns present, shift all key value from unknown to a level up
	//	for key, value := range unknowns.(tcontainer.MarshalMap) {
	//		m[key] = value
	//	}
	//	delete(m, "Unknowns")
	//}
	return json.Marshal(m)
}
