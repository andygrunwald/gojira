package jira

import "context"

// FieldService handles fields for the Jira instance / API.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/#api-Field
type FieldService struct {
	client *Client
}

// GetListWithContext gets all fields from Jira
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/#api-api-2-field-get
func (s *FieldService) GetListWithContext(ctx context.Context) ([]Field, *Response, error) {
	apiEndpoint := "rest/api/2/field"
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	fieldList := []Field{}
	v, resp, err := s.client.Do(req)

	if err != nil {
		return nil, resp, NewJiraError(resp, err)
	}
	return fieldList, resp, nil
}

// GetList wraps GetListWithContext using the background context.
func (s *FieldService) GetList() ([]Field, *Response, error) {
	return s.GetListWithContext(context.Background())
}
