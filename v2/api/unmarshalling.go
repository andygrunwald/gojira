package api

type Unmarshaller interface {
	Unmarshal(s string) error
}
