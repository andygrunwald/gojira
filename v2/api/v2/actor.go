package v2

// Actor represents a Jira actor
type Actor struct {
	ID          int        `json:"id" structs:"id"`
	DisplayName string     `json:"displayName" structs:"displayName"`
	Type        string     `json:"type" structs:"type"`
	Name        string     `json:"name" structs:"name"`
	AvatarURL   string     `json:"avatarUrl" structs:"avatarUrl"`
	ActorUser   *ActorUser `json:"actorUser" structs:"actoruser"`
}
