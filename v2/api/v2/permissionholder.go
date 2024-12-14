package v2

type PermissionHolder struct {
	Type      string `json:"type" structs:"type"`
	Parameter string `json:"parameter" structs:"parameter"`
	Expand    string `json:"expand" structs:"expand"`
}
