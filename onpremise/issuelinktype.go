package onpremise

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// IssueLinkTypeService handles issue link types for the Jira instance / API.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/#api-group-Issue-link-types
type IssueLinkTypeService service

// GetList gets all of the issue link types from Jira.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/#api-rest-api-2-issueLinkType-get
func (s *IssueLinkTypeService) GetList(ctx context.Context) ([]IssueLinkType, *Response, error) {
	apiEndpoint := "rest/api/2/issueLinkType"
	req, err := s.client.NewRequest(ctx, http.MethodGet, apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	var linkTypeList map[string]([]IssueLinkType)
	resp, err := s.client.Do(req, &linkTypeList)
	if err != nil {
		return nil, resp, NewJiraError(resp, err)
	}
	return linkTypeList["issueLinkTypes"], resp, nil
}

// Get gets info of a specific issue link type from Jira.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/#api-rest-api-2-issueLinkType-issueLinkTypeId-get
//
// TODO Double check this method if this works as expected, is using the latest API and the response is complete
// This double check effort is done for v2 - Remove this two lines if this is completed.
func (s *IssueLinkTypeService) Get(ctx context.Context, ID string) (*IssueLinkType, *Response, error) {
	apiEndPoint := fmt.Sprintf("rest/api/2/issueLinkType/%s", ID)
	req, err := s.client.NewRequest(ctx, http.MethodGet, apiEndPoint, nil)
	if err != nil {
		return nil, nil, err
	}

	linkType := new(IssueLinkType)
	resp, err := s.client.Do(req, linkType)
	if err != nil {
		return nil, resp, NewJiraError(resp, err)
	}
	return linkType, resp, nil
}

// Create creates an issue link type in Jira.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/#api-rest-api-2-issueLinkType-post
//
// TODO Double check this method if this works as expected, is using the latest API and the response is complete
// This double check effort is done for v2 - Remove this two lines if this is completed.
func (s *IssueLinkTypeService) Create(ctx context.Context, linkType *IssueLinkType) (*IssueLinkType, *Response, error) {
	apiEndpoint := "/rest/api/2/issueLinkType"
	req, err := s.client.NewRequest(ctx, http.MethodPost, apiEndpoint, linkType)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()

	responseLinkType := new(IssueLinkType)
	err = json.NewDecoder(resp.Body).Decode(&responseLinkType)
	if err != nil {
		return nil, resp, err
	}

	return linkType, resp, nil
}

// Update updates an issue link type.  The issue is found by key.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/#api-rest-api-2-issueLinkType-issueLinkTypeId-put
// Caller must close resp.Body
//
// TODO Double check this method if this works as expected, is using the latest API and the response is complete
// This double check effort is done for v2 - Remove this two lines if this is completed.
func (s *IssueLinkTypeService) Update(ctx context.Context, linkType *IssueLinkType) (*IssueLinkType, *Response, error) {
	apiEndpoint := fmt.Sprintf("rest/api/2/issueLinkType/%s", linkType.ID)
	req, err := s.client.NewRequest(ctx, http.MethodPut, apiEndpoint, linkType)
	if err != nil {
		return nil, nil, err
	}
	resp, err := s.client.Do(req, nil)
	if err != nil {
		return nil, resp, NewJiraError(resp, err)
	}
	ret := *linkType
	return &ret, resp, nil
}

// Delete deletes an issue link type based on provided ID.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/#api-rest-api-2-issueLinkType-issueLinkTypeId-delete
// Caller must close resp.Body
//
// TODO Double check this method if this works as expected, is using the latest API and the response is complete
// This double check effort is done for v2 - Remove this two lines if this is completed.
func (s *IssueLinkTypeService) Delete(ctx context.Context, ID string) (*Response, error) {
	apiEndpoint := fmt.Sprintf("rest/api/2/issueLinkType/%s", ID)
	req, err := s.client.NewRequest(ctx, http.MethodDelete, apiEndpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	return resp, err
}
