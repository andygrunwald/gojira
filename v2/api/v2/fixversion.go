package v2

import "time"

// FixVersion represents a software release in which an issue is fixed.
type FixVersion struct {
	Self            string    `json:"self,omitempty" structs:"self,omitempty"`
	ID              string    `json:"id,omitempty" structs:"id,omitempty"`
	Name            string    `json:"name,omitempty" structs:"name,omitempty"`
	Description     string    `json:"description,omitempty" structs:"description,omitempty"`
	Archived        *bool     `json:"archived,omitempty" structs:"archived,omitempty"`
	Released        *bool     `json:"released,omitempty" structs:"released,omitempty"`
	ReleaseDate     time.Time `json:"releaseDate,omitempty" structs:"releaseDate,omitempty"`
	UserReleaseDate time.Time `json:"userReleaseDate,omitempty" structs:"userReleaseDate,omitempty"`
	ProjectID       int       `json:"projectId,omitempty" structs:"projectId,omitempty"` // Unlike other IDs, this is returned as a number
	StartDate       string    `json:"startDate,omitempty" structs:"startDate,omitempty"`
}
