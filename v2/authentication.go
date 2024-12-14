package jira

import (
	"context"
	"fmt"
	"net/http"
)

const (
	// HTTP Basic Authentication
	authTypeBasic = iota + 1
	// HTTP Session Authentication
	authTypeSession
)

// AuthenticationService handles authentication for the Jira instance / API.
//
// Jira API docs: https://docs.atlassian.com/jira/REST/latest/#authentication
type AuthenticationService struct {
	client *Client

	// Authentication type
	authType int

	// Basic auth username
	username string

	// Basic auth password
	password string
}

// Authenticated reports if the current Client has authentication details for Jira
func (s *AuthenticationService) Authenticated() bool {
	switch s.authType {
	case authTypeSession:
		return s.client.session != nil
	case authTypeBasic:
		return s.username != ""
	}

	return false
}

// GetCurrentUserWithContext gets the details of the current user.
//
// Jira API docs: https://docs.atlassian.com/jira/REST/latest/#auth/1/session
func (s *AuthenticationService) GetCurrentUser(ctx context.Context) (*Session, error) {
	if s.authType != authTypeSession || s.client.session == nil {
		return nil, fmt.Errorf("no user is authenticated yet")
	}

	apiEndpoint := "rest/auth/1/session"
	req, err := s.client.NewRequest(ctx, http.MethodGet, apiEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request for getting user info : %s", err)
	}

	v, _, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info : %s", err)
	}

	session := new(Session)
	err = session.Unmarshal(v.String())
	if err != nil {
		return nil, fmt.Errorf("could not unmarshall received user info : %s", err)
	}

	return session, nil
}

// Session represents a Session JSON response by the Jira API.
type Session struct {
	Self    string `json:"self,omitempty"`
	Name    string `json:"name,omitempty"`
	Session struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"session,omitempty"`
	LoginInfo struct {
		FailedLoginCount    int    `json:"failedLoginCount"`
		LoginCount          int    `json:"loginCount"`
		LastFailedLoginTime string `json:"lastFailedLoginTime"`
		PreviousLoginTime   string `json:"previousLoginTime"`
	} `json:"loginInfo"`
	Cookies []*http.Cookie
}

func (s *Session) Unmarshal(j string) error {
	return nil
}
