package jira

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	errors2 "github.com/andygrunwald/go-jira/v2/errors"
	"github.com/valyala/fastjson"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// httpClient defines an interface for an http.Client implementation so that alternative
// http Clients can be passed in for making requests
type HttpClient interface {
	Do(request *http.Request) (response *http.Response, err error)
}

// A Client manages communication with the Jira API.
type Client struct {
	// HTTP client used to communicate with the API.
	HttpClient HttpClient
	parserPool fastjson.ParserPool
	arenaPool  fastjson.ArenaPool

	// Base URL for API requests.
	baseURL *url.URL

	// Session storage if the user authenticates with a Session cookie
	session *Session

	// Services used for talking to different parts of the Jira API.
	Authentication   *AuthenticationService
	Issue            *IssueService
	Project          *ProjectService
	Board            *BoardService
	Sprint           *SprintService
	User             *UserService
	Group            *GroupService
	Version          *VersionService
	Priority         *PriorityService
	Field            *FieldService
	Component        *ComponentService
	Resolution       *ResolutionService
	StatusCategory   *StatusCategoryService
	Filter           *FilterService
	Role             *RoleService
	PermissionScheme *PermissionSchemeService
	Status           *StatusService
	IssueLinkType    *IssueLinkTypeService
}

type ClientOption func(c *Client)

type HttpRequestOption func(r *http.Request) error

// NewClient returns a new Jira API client.
// To use API methods which require authentication you can follow the preferred solution and
// provide an http.Client that will perform the authentication for you with OAuth and HTTP Basic (such as that provided by the golang.org/x/oauth2 library).
// As an alternative you can use Session Cookie based authentication provided by this package as well.
// See https://docs.atlassian.com/jira/REST/latest/#authentication
func NewClient(baseURL string, options ...ClientOption) (*Client, error) {
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		HttpClient: http.DefaultClient,
		parserPool: fastjson.ParserPool{},
		arenaPool:  fastjson.ArenaPool{},
		baseURL:    parsedBaseURL,
	}
	c.Authentication = &AuthenticationService{client: c}
	c.Issue = &IssueService{client: c}
	c.Project = &ProjectService{client: c}
	c.Board = &BoardService{client: c}
	c.Sprint = &SprintService{client: c}
	c.User = &UserService{client: c}
	c.Group = &GroupService{client: c}
	c.Version = &VersionService{client: c}
	c.Priority = &PriorityService{client: c}
	c.Field = &FieldService{client: c}
	c.Component = &ComponentService{client: c}
	c.Resolution = &ResolutionService{client: c}
	c.StatusCategory = &StatusCategoryService{client: c}
	c.Filter = &FilterService{client: c}
	c.Role = &RoleService{client: c}
	c.PermissionScheme = &PermissionSchemeService{client: c}
	c.Status = &StatusService{client: c}
	c.IssueLinkType = &IssueLinkTypeService{client: c}

	for _, o := range options {
		o(c)
	}

	return c, nil
}

// NewRawRequestWithContext creates an API request.
// A relative URL can be provided in urlStr, in which case it is resolved relative to the baseURL of the Client.
// Allows using an optional native io.Reader for sourcing the request body.
func (c *Client) NewRawRequest(ctx context.Context, method, urlStr string, body io.Reader, options ...HttpRequestOption) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	// Relative URLs should be specified without a preceding slash since baseURL will have the trailing slash
	rel.Path = strings.TrimLeft(rel.Path, "/")

	u := c.baseURL.ResolveReference(rel)

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		// since this method is a wrapper around net/http.NewRequestWithContext
		// wrapping the returned error here would likely be redundant
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	switch c.Authentication.authType {
	case authTypeSession:
		// Set session cookie if there is one
		if c.session != nil {
			for _, cookie := range c.session.Cookies {
				req.AddCookie(cookie)
			}
		}
	case authTypeBasic:
		// Set basic auth information
		if c.Authentication.username != "" {
			req.SetBasicAuth(c.Authentication.username, c.Authentication.password)
		}
	}

	for _, opt := range options {
		err := opt(req)
		if err != nil {
			return req, err
		}
	}

	return req, nil
}

// NewRequestWithContext creates an API request.
// A relative URL can be provided in urlStr, in which case it is resolved relative to the baseURL of the Client.
// If specified, the value pointed to by body is JSON encoded and included as the request body.
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}, options ...HttpRequestOption) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, fmt.Errorf("failed to encode json body: %w", err)
		}
	}

	return c.NewRawRequest(ctx, method, urlStr, buf, options...)
}

// addOptions adds the parameters in opt as URL query parameters to s.  opt
// must be a struct whose fields may contain "url" tags.
//func addOptions(s string, opt interface{}) (string, error) {
//	v := reflect.ValueOf(opt)
//	if v.Kind() == reflect.Ptr && v.IsNil() {
//		return s, nil
//	}
//
//	u, err := url.Parse(s)
//	if err != nil {
//		return s, err
//	}
//
//	qs, err := query.Values(opt)
//	if err != nil {
//		return s, err
//	}
//
//	u.RawQuery = qs.Encode()
//	return u.String(), nil
//}

// Do sends an API request and returns the API response.
// The API response is JSON decoded and returned as a fastjson.Value when possible
// Value may be nil when
// * request failed to send
// * no response received
// * body empty
// * response content type was not application/json
func (c *Client) Do(req *http.Request) (value *fastjson.Value, resp *http.Response, err error) {
	resp, err = c.HttpClient.Do(req)
	if err != nil {
		return nil, resp, fmt.Errorf("error making http request: %w", err)
	}

	if resp == nil {
		return nil, resp, errors.New("no response returned")
	}

	// read the body, but set it back on the http response to be read later
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to read body: %w", err)
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	if len(body) > 0 && resp.Header.Get("Content-Type") == "application/json" {
		parser := c.parserPool.Get()
		defer c.parserPool.Put(parser)

		value, err := parser.ParseBytes(body)
		if err != nil {
			return value, resp, fmt.Errorf("failed to parse body: %w", err)
		}
	}

	if c := resp.StatusCode; !(200 <= c && c <= 299) {
		return value, resp, errors2.NewJiraRequestError(resp)
	}

	return value, resp, nil
}

// GetBaseURL will return you the Base URL.
// This is the same URL as in the NewClient constructor
func (c *Client) GetBaseURL() url.URL {
	return *c.baseURL
}

// Response represents Jira API response. It wraps http.Response returned from
// API and provides information about paging.
//type Response struct {
//	*http.Response
//
//	StartAt    int
//	MaxResults int
//	Total      int
//}

//func NewResponse(r *http.Response, v interface{}) *Response {
//	resp := &Response{Response: r}
//
//	if v == nil {
//		return resp
//	}
//
//	switch value := v.(type) {
//	case *SearchResult:
//		resp.StartAt = value.StartAt
//		resp.MaxResults = value.MaxResults
//		resp.Total = value.Total
//	case *GroupMembersResult:
//		resp.StartAt = value.StartAt
//		resp.MaxResults = value.MaxResults
//		resp.Total = value.Total
//	}
//	return resp
//}

// cloneRequest returns a clone of the provided *http.Request.
// The clone is a shallow copy of the struct and its Header map.
//func cloneRequest(r *http.Request) *http.Request {
//	// shallow copy of the struct
//	r2 := new(http.Request)
//	*r2 = *r
//	// deep copy of the Header
//	r2.Header = make(http.Header, len(r.Header))
//	for k, s := range r.Header {
//		r2.Header[k] = append([]string(nil), s...)
//	}
//	return r2
//}
