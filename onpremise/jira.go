package onpremise

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"

	"github.com/google/go-querystring/query"
)

const (
	ClientVersion = "2.0.0"

	defaultUserAgent = "go-jira" + "/" + ClientVersion
)

// A Client manages communication with the Jira API.
type Client struct {
	clientMu sync.Mutex   // clientMu protects the client during calls that modify it.
	client   *http.Client // HTTP client used to communicate with the API.

	// Base URL for API requests.
	// Should be set to a domain endpoint of the Jira instance.
	// BaseURL should always be specified with a trailing slash.
	BaseURL *url.URL

	// User agent used when communicating with the Jira API.
	UserAgent string

	// Session storage if the user authenticates with a Session cookie
	// TODO Needed in Cloud and/or onpremise?
	session *Session

	// Reuse a single struct instead of allocating one for each service on the heap.
	common service

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
	Organization     *OrganizationService
	ServiceDesk      *ServiceDeskService
	Customer         *CustomerService
	Request          *RequestService
}

// service is the base structure to bundle API services
// under a sub-struct.
type service struct {
	client *Client
}

// Client returns the http.Client used by this Jira client.
func (c *Client) Client() *http.Client {
	c.clientMu.Lock()
	defer c.clientMu.Unlock()
	clientCopy := *c.client
	return &clientCopy
}

// NewClient returns a new Jira API client with provided base URL (often is your Jira hostname)
// If a nil httpClient is provided, a new http.Client will be used.
// To use API methods which require authentication, provide an http.Client that will perform the authentication for you (such as that provided by the golang.org/x/oauth2 library).
// baseURL is the HTTP endpoint of your Jira instance and should always be specified with a trailing slash.
func NewClient(baseURL string, httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	baseEndpoint, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	// ensure the baseURL contains a trailing slash so that all paths are preserved in later calls
	if !strings.HasSuffix(baseEndpoint.Path, "/") {
		baseEndpoint.Path += "/"
	}

	c := &Client{
		client:    httpClient,
		BaseURL:   baseEndpoint,
		UserAgent: defaultUserAgent,
	}
	c.common.client = c

	// TODO Check if the authentication service is still needed (because of the transports)
	c.Authentication = &AuthenticationService{client: c}
	c.Issue = (*IssueService)(&c.common)
	c.Project = (*ProjectService)(&c.common)
	c.Board = (*BoardService)(&c.common)
	c.Sprint = (*SprintService)(&c.common)
	c.User = (*UserService)(&c.common)
	c.Group = (*GroupService)(&c.common)
	c.Version = (*VersionService)(&c.common)
	c.Priority = (*PriorityService)(&c.common)
	c.Field = (*FieldService)(&c.common)
	c.Component = (*ComponentService)(&c.common)
	c.Resolution = (*ResolutionService)(&c.common)
	c.StatusCategory = (*StatusCategoryService)(&c.common)
	c.Filter = (*FilterService)(&c.common)
	c.Role = (*RoleService)(&c.common)
	c.PermissionScheme = (*PermissionSchemeService)(&c.common)
	c.Status = (*StatusService)(&c.common)
	c.IssueLinkType = (*IssueLinkTypeService)(&c.common)
	c.Organization = (*OrganizationService)(&c.common)
	c.ServiceDesk = (*ServiceDeskService)(&c.common)
	c.Customer = (*CustomerService)(&c.common)
	c.Request = (*RequestService)(&c.common)

	return c, nil
}

// TODO Do we need it?
// NewRawRequest creates an API request.
// A relative URL can be provided in urlStr, in which case it is resolved relative to the baseURL of the Client.
// Allows using an optional native io.Reader for sourcing the request body.
func (c *Client) NewRawRequest(ctx context.Context, method, urlStr string, body io.Reader) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	// Relative URLs should be specified without a preceding slash since baseURL will have the trailing slash
	rel.Path = strings.TrimLeft(rel.Path, "/")

	u := c.BaseURL.ResolveReference(rel)

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Set authentication information
	if c.Authentication.authType == authTypeSession {
		// Set session cookie if there is one
		if c.session != nil {
			for _, cookie := range c.session.Cookies {
				req.AddCookie(cookie)
			}
		}
	} else if c.Authentication.authType == authTypeBasic {
		// Set basic auth information
		if c.Authentication.username != "" {
			req.SetBasicAuth(c.Authentication.username, c.Authentication.password)
		}
	}

	return req, nil
}

// NewRequest creates an API request.
// A relative URL can be provided in urlStr, in which case it is resolved relative to the BaseURL of the Client.
// If specified, the value pointed to by body is JSON encoded and included as the request body.
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	// Relative URLs should be specified without a preceding slash since BaseURL will have the trailing slash
	rel.Path = strings.TrimLeft(rel.Path, "/")

	u := c.BaseURL.ResolveReference(rel)

	// TODO This part is the difference between NewRawRequestWithContext
	// Check if we can get this working in one function
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Set authentication information
	if c.Authentication.authType == authTypeSession {
		// Set session cookie if there is one
		if c.session != nil {
			for _, cookie := range c.session.Cookies {
				req.AddCookie(cookie)
			}
		}
	} else if c.Authentication.authType == authTypeBasic {
		// Set basic auth information
		if c.Authentication.username != "" {
			req.SetBasicAuth(c.Authentication.username, c.Authentication.password)
		}
	}

	return req, nil
}

// addOptions adds the parameters in opts as URL query parameters to s. opts
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opts interface{}) (string, error) {
	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opts)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// NewMultiPartRequest creates an API request including a multi-part file.
// A relative URL can be provided in urlStr, in which case it is resolved relative to the baseURL of the Client.
// If specified, the value pointed to by buf is a multipart form.
func (c *Client) NewMultiPartRequest(ctx context.Context, method, urlStr string, buf *bytes.Buffer) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	// Relative URLs should be specified without a preceding slash since baseURL will have the trailing slash
	rel.Path = strings.TrimLeft(rel.Path, "/")

	u := c.BaseURL.ResolveReference(rel)

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	// Set required headers
	req.Header.Set("X-Atlassian-Token", "nocheck")

	// Set authentication information
	if c.Authentication.authType == authTypeSession {
		// Set session cookie if there is one
		if c.session != nil {
			for _, cookie := range c.session.Cookies {
				req.AddCookie(cookie)
			}
		}
	} else if c.Authentication.authType == authTypeBasic {
		// Set basic auth information
		if c.Authentication.username != "" {
			req.SetBasicAuth(c.Authentication.username, c.Authentication.password)
		}
	}

	return req, nil
}

// Do sends an API request and returns the API response.
// The API response is JSON decoded and stored in the value pointed to by v, or returned as an error if an API error has occurred.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	httpResp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	err = CheckResponse(httpResp)
	if err != nil {
		// Even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return newResponse(httpResp, nil), err
	}

	if v != nil {
		// Open a NewDecoder and defer closing the reader only if there is a provided interface to decode to
		body, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			return newResponse(httpResp, nil), err
		}
		err = json.Unmarshal(body, v)
	}

	resp := newResponse(httpResp, v)
	return resp, err
}

// CheckResponse checks the API response for errors, and returns them if present.
// A response is considered an error if it has a status code outside the 200 range.
// The caller is responsible to analyze the response body.
// The body can contain JSON (if the error is intended) or xml (sometimes Jira just failes).
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	err := fmt.Errorf("request failed. Please analyze the request body for more details. Status code: %d", r.StatusCode)
	return err
}

// Response represents Jira API response. It wraps http.Response returned from
// API and provides information about paging.
type Response struct {
	*http.Response

	StartAt    int
	MaxResults int
	Total      int
}

func newResponse(r *http.Response, v interface{}) *Response {
	resp := &Response{Response: r}
	resp.populatePageValues(v)
	return resp
}

// Sets paging values if response json was parsed to searchResult type
// (can be extended with other types if they also need paging info)
func (r *Response) populatePageValues(v interface{}) {
	switch value := v.(type) {
	case *searchResult:
		r.StartAt = value.StartAt
		r.MaxResults = value.MaxResults
		r.Total = value.Total
	case *groupMembersResult:
		r.StartAt = value.StartAt
		r.MaxResults = value.MaxResults
		r.Total = value.Total
	}
}
