package v2

// These constants are the keys of the default Jira status categories
const (
	IssueStatusCategoryComplete   = "done"
	IssueStatusCategoryInProgress = "indeterminate"
	IssueStatusCategoryToDo       = "new"
	IssueStatusCategoryUndefined  = "undefined"
)

// StatusCategory represents the category a status belongs to.
// Those categories can be user defined in every Jira instance.
type StatusCategory struct {
	Self      string `json:"self" structs:"self"`
	ID        int    `json:"id" structs:"id"`
	Name      string `json:"name" structs:"name"`
	Key       string `json:"key" structs:"key"`
	ColorName string `json:"colorName" structs:"colorName"`
}
