package jira

import "context"

// PriorityService handles priorities for the Jira instance / API.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/#api-Priority
type PriorityService struct {
	client *Client
}

// GetListWithContext gets all priorities from Jira
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/#api-api-2-priority-get
func (s *PriorityService) GetListWithContext(ctx context.Context) ([]Priority, *Response, error) {
	apiEndpoint := "rest/api/2/priority"
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	priorityList := []Priority{}
	resp, err := s.client.Do(req, &priorityList)
	if err != nil {
		return nil, resp, NewJiraError(resp, err)
	}
	return priorityList, resp, nil
}

// GetList wraps GetListWithContext using the background context.
func (s *PriorityService) GetList() ([]Priority, *Response, error) {
	return s.GetListWithContext(context.Background())
}
