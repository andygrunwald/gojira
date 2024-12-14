package v2

type Permission struct {
	ID     int    `json:"id" structs:"id"`
	Self   string `json:"expand" structs:"expand"`
	Holder Holder `json:"holder" structs:"holder"`
	Name   string `json:"permission" structs:"permission"`
}
