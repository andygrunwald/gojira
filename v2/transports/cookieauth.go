package transports

import (
	"bytes"
	"encoding/json"
	"fmt"
	jerrors "github.com/andygrunwald/go-jira/v2/errors"
	"net/http"
	"time"
)

// CookieAuthTransport is an http.RoundTripper that authenticates all requests
// using Jira's cookie-based authentication.
//
// Note that it is generally preferable to use HTTP BASIC authentication with the REST API.
// However, this resource may be used to mimic the behaviour of Jira's log-in page (e.g. to display log-in errors to a user).
//
// Jira API docs: https://docs.atlassian.com/jira/REST/latest/#auth/1/session
type CookieAuth struct {
	username string
	password string
	AuthURL  string

	// SessionObject is the authenticated cookie string.s
	// It's passed in each call to prove the client is authenticated.
	SessionObject []*http.Cookie

	// Transport is the underlying HTTP transport to use when making requests.
	// It will default to http.DefaultTransport if nil.
	Transport http.RoundTripper
}

// RoundTrip adds the session object to the request.
func (t *CookieAuth) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.SessionObject == nil {
		err := t.setSessionObject()
		if err != nil {
			return nil, fmt.Errorf("cookieauth: no session object has been set: %w", err)
		}
	}

	for _, cookie := range t.SessionObject {
		// Don't add an empty value cookie to the request
		if cookie.Value != "" {
			req.AddCookie(cookie)
		}
	}

	return t.transport().RoundTrip(req)
}

// Client returns an *http.Client that makes requests that are authenticated
// using cookie authentication
func (t *CookieAuth) Client() *http.Client {
	return &http.Client{Transport: t}
}

// setSessionObject attempts to authenticate the user and set
// the session object (e.g. cookie)
func (t *CookieAuth) setSessionObject() error {
	req, err := t.buildAuthRequest()
	if err != nil {
		return fmt.Errorf("failed to create auth request: %w", err)
	}

	var authClient = &http.Client{
		Timeout: time.Second * 60,
	}
	resp, err := authClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	t.SessionObject = resp.Cookies()
	return nil
}

// getAuthRequest assembles the request to get the authenticated cookie
func (t *CookieAuth) buildAuthRequest() (*http.Request, error) {
	body := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		t.username,
		t.password,
	}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(body)
	if err != nil {
		return nil, fmt.Errorf("failed to encode body as json: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, t.AuthURL, b)
	if err != nil {
		return nil, fmt.Errorf(jerrors.HttpRequestCreationFailureMessageFormat, err)
	}

	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (t *CookieAuth) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}
	return http.DefaultTransport
}
