package cloud

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/go-querystring/query"
)

// ServiceDeskService handles ServiceDesk for the Jira instance / API.
type ServiceDeskService service

// ServiceDeskOrganizationDTO is a DTO for ServiceDesk organizations
type ServiceDeskOrganizationDTO struct {
	OrganizationID int `json:"organizationId,omitempty" structs:"organizationId,omitempty"`
}

// GetOrganizations returns a list of
// all organizations associated with a service desk.
//
// https://developer.atlassian.com/cloud/jira/service-desk/rest/api-group-organization/#api-rest-servicedeskapi-servicedesk-servicedeskid-organization-get
//
// TODO Double check this method if this works as expected, is using the latest API and the response is complete
// This double check effort is done for v2 - Remove this two lines if this is completed.
func (s *ServiceDeskService) GetOrganizations(ctx context.Context, serviceDeskID interface{}, start int, limit int, accountID string) (*PagedDTO, *Response, error) {
	apiEndPoint := fmt.Sprintf("rest/servicedeskapi/servicedesk/%v/organization?start=%d&limit=%d", serviceDeskID, start, limit)
	if accountID != "" {
		apiEndPoint += fmt.Sprintf("&accountId=%s", accountID)
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, apiEndPoint, nil)
	req.Header.Set("Accept", "application/json")

	if err != nil {
		return nil, nil, err
	}

	orgs := new(PagedDTO)
	resp, err := s.client.Do(req, &orgs)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}

	return orgs, resp, nil
}

// AddOrganization adds an organization to
// a service desk. If the organization ID is already
// associated with the service desk, no change is made
// and the resource returns a 204 success code.
//
// https://developer.atlassian.com/cloud/jira/service-desk/rest/api-group-organization/#api-rest-servicedeskapi-servicedesk-servicedeskid-organization-post
// Caller must close resp.Body
//
// TODO Double check this method if this works as expected, is using the latest API and the response is complete
// This double check effort is done for v2 - Remove this two lines if this is completed.
func (s *ServiceDeskService) AddOrganization(ctx context.Context, serviceDeskID interface{}, organizationID int) (*Response, error) {
	apiEndPoint := fmt.Sprintf("rest/servicedeskapi/servicedesk/%v/organization", serviceDeskID)

	organization := ServiceDeskOrganizationDTO{
		OrganizationID: organizationID,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, apiEndPoint, organization)

	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return resp, jerr
	}

	return resp, nil
}

// RemoveOrganization removes an organization
// from a service desk. If the organization ID does not
// match an organization associated with the service desk,
// no change is made and the resource returns a 204 success code.
//
// https://developer.atlassian.com/cloud/jira/service-desk/rest/api-group-organization/#api-rest-servicedeskapi-servicedesk-servicedeskid-organization-delete
// Caller must close resp.Body
//
// TODO Double check this method if this works as expected, is using the latest API and the response is complete
// This double check effort is done for v2 - Remove this two lines if this is completed.
func (s *ServiceDeskService) RemoveOrganization(ctx context.Context, serviceDeskID interface{}, organizationID int) (*Response, error) {
	apiEndPoint := fmt.Sprintf("rest/servicedeskapi/servicedesk/%v/organization", serviceDeskID)

	organization := ServiceDeskOrganizationDTO{
		OrganizationID: organizationID,
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, apiEndPoint, organization)

	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return resp, jerr
	}

	return resp, nil
}

// AddCustomers adds customers to the given service desk.
//
// https://developer.atlassian.com/cloud/jira/service-desk/rest/api-group-servicedesk/#api-rest-servicedeskapi-servicedesk-servicedeskid-customer-post
//
// TODO Double check this method if this works as expected, is using the latest API and the response is complete
// This double check effort is done for v2 - Remove this two lines if this is completed.
func (s *ServiceDeskService) AddCustomers(ctx context.Context, serviceDeskID interface{}, acountIDs ...string) (*Response, error) {
	apiEndpoint := fmt.Sprintf("rest/servicedeskapi/servicedesk/%v/customer", serviceDeskID)

	payload := struct {
		AccountIDs []string `json:"accountIds"`
	}{
		AccountIDs: acountIDs,
	}
	req, err := s.client.NewRequest(ctx, http.MethodPost, apiEndpoint, payload)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, NewJiraError(resp, err)
	}

	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	return resp, nil
}

// RemoveCustomers removes customers to the given service desk.
//
// https://developer.atlassian.com/cloud/jira/service-desk/rest/api-group-servicedesk/#api-rest-servicedeskapi-servicedesk-servicedeskid-customer-delete
//
// TODO Double check this method if this works as expected, is using the latest API and the response is complete
// This double check effort is done for v2 - Remove this two lines if this is completed.
func (s *ServiceDeskService) RemoveCustomers(ctx context.Context, serviceDeskID interface{}, acountIDs ...string) (*Response, error) {
	apiEndpoint := fmt.Sprintf("rest/servicedeskapi/servicedesk/%v/customer", serviceDeskID)

	payload := struct {
		AccountIDs []string `json:"accountIDs"`
	}{
		AccountIDs: acountIDs,
	}
	req, err := s.client.NewRequest(ctx, http.MethodDelete, apiEndpoint, payload)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, NewJiraError(resp, err)
	}

	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	return resp, nil
}

// ListCustomers lists customers for a ServiceDesk.
//
// https://developer.atlassian.com/cloud/jira/service-desk/rest/api-group-servicedesk/#api-rest-servicedeskapi-servicedesk-servicedeskid-customer-get
//
// TODO Double check this method if this works as expected, is using the latest API and the response is complete
// This double check effort is done for v2 - Remove this two lines if this is completed.
func (s *ServiceDeskService) ListCustomers(ctx context.Context, serviceDeskID interface{}, options *CustomerListOptions) (*CustomerList, *Response, error) {
	apiEndpoint := fmt.Sprintf("rest/servicedeskapi/servicedesk/%v/customer", serviceDeskID)
	req, err := s.client.NewRequest(ctx, http.MethodGet, apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	// this is an experiemntal endpoint
	req.Header.Set("X-ExperimentalApi", "opt-in")

	if options != nil {
		q, err := query.Values(options)
		if err != nil {
			return nil, nil, err
		}
		req.URL.RawQuery = q.Encode()
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return nil, resp, NewJiraError(resp, err)
	}
	defer resp.Body.Close()

	customerList := new(CustomerList)
	if err := json.NewDecoder(resp.Body).Decode(customerList); err != nil {
		return nil, resp, fmt.Errorf("could not unmarshall the data into struct")
	}

	return customerList, resp, nil
}

// ListComments lists comments for a ServiceDesk issue.
//
// https://docs.atlassian.com/jira-servicedesk/REST/3.6.2/#servicedeskapi/request/{issueIdOrKey}/comment-getRequestComments
func (s *ServiceDeskService) ListComments(ctx context.Context, issueKeyOrID string, params ServiceDeskCommentsListParams) (*ServiceDeskCommentsList, *Response, error) {
	apiEndPoint := fmt.Sprintf("rest/servicedeskapi/request/%s/comment", issueKeyOrID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, apiEndPoint, nil)
	req.Header.Set("Accept", "application/json")

	query := url.Values{}
	query.Add("start", strconv.Itoa(params.Start))
	query.Add("limit", strconv.Itoa(params.Limit))

	if params.HidePublic {
		query.Add("public", "false")
	} else {
		query.Add("public", "true")
	}

	if params.HideInternal {
		query.Add("internal", "false")
	} else {
		query.Add("internal", "true")
	}

	// Add parameters to url
	req.URL.RawQuery = query.Encode()

	if err != nil {
		return nil, nil, err
	}

	orgs := new(ServiceDeskCommentsList)
	resp, err := s.client.Do(req, &orgs)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}
	defer resp.Body.Close()

	return orgs, resp, nil
}

type ServiceDeskCommentsListParams struct {
	// Specifies whether to hide public comments or not. Default: false.
	HidePublic bool
	// Specifies whether to hide internal comments or not. Default: false.
	HideInternal bool
	// The starting index of the returned comments. Base index: 0. See the Pagination(https://docs.atlassian.com/jira-servicedesk/REST/3.6.2/#pagination) section for more details.
	Start int
	// The maximum number of comments to return per page. Default: 50. See the Pagination(https://docs.atlassian.com/jira-servicedesk/REST/3.6.2/#pagination) section for more details.
	Limit int
}

type ServiceDeskComment struct {
	ID     string `json:"id"`
	Body   string `json:"body"`
	Public bool   `json:"public"`
	Author struct {
		Name         string `json:"name"`
		Key          string `json:"key"`
		EmailAddress string `json:"emailAddress"`
		DisplayName  string `json:"displayName"`
		Active       bool   `json:"active"`
		TimeZone     string `json:"timeZone"`
		Links        struct {
			JiraRest   string `json:"jiraRest"`
			AvatarUrls struct {
				Four8X48  string `json:"48x48"`
				Two4X24   string `json:"24x24"`
				One6X16   string `json:"16x16"`
				Three2X32 string `json:"32x32"`
			} `json:"avatarUrls"`
			Self string `json:"self"`
		} `json:"_links"`
	} `json:"author"`
	Created struct {
		Iso8601     string `json:"iso8601"`
		Jira        string `json:"jira"`
		Friendly    string `json:"friendly"`
		EpochMillis int64  `json:"epochMillis"`
	} `json:"created"`
	Links struct {
		Self string `json:"self"`
	} `json:"_links"`
}

type ServiceDeskCommentsList struct {
	Expands    []interface{} `json:"_expands"`
	Size       int           `json:"size"`
	Start      int           `json:"start"`
	Limit      int           `json:"limit"`
	IsLastPage bool          `json:"isLastPage"`
	Links      struct {
		Base    string `json:"base"`
		Context string `json:"context"`
		Next    string `json:"next"`
		Prev    string `json:"prev"`
	} `json:"_links"`
	Values []ServiceDeskComment `json:"values"`
}
