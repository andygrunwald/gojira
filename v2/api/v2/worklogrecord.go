package v2

import "time"

// WorklogRecord represents one entry of a Worklog
type WorklogRecord struct {
	Self             string           `json:"self,omitempty" structs:"self,omitempty"`
	Author           *User            `json:"author,omitempty" structs:"author,omitempty"`
	UpdateAuthor     *User            `json:"updateAuthor,omitempty" structs:"updateAuthor,omitempty"`
	Comment          string           `json:"comment,omitempty" structs:"comment,omitempty"`
	Created          *time.Time       `json:"created,omitempty" structs:"created,omitempty"`
	Updated          *time.Time       `json:"updated,omitempty" structs:"updated,omitempty"`
	Started          *time.Time       `json:"started,omitempty" structs:"started,omitempty"`
	TimeSpent        string           `json:"timeSpent,omitempty" structs:"timeSpent,omitempty"`
	TimeSpentSeconds int              `json:"timeSpentSeconds,omitempty" structs:"timeSpentSeconds,omitempty"`
	ID               string           `json:"id,omitempty" structs:"id,omitempty"`
	IssueID          string           `json:"issueId,omitempty" structs:"issueId,omitempty"`
	Properties       []EntityProperty `json:"properties,omitempty"`
}
