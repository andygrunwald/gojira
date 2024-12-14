package v2

// ChangelogHistory reflects one single changelog history entry
type ChangeLogHistory struct {
	Id      string           `json:"id" structs:"id"`
	Author  User             `json:"author" structs:"author"`
	Created string           `json:"created" structs:"created"`
	Items   []ChangeLogItems `json:"items" structs:"items"`
}
