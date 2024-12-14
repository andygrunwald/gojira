package jira

import "context"

// StatusService handles staties for the Jira instance / API.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/#api-group-Workflow-statuses
type StatusService struct {
	client *Client
}

// GetAllStatusesWithContext returns a list of all statuses associated with workflows.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/#api-rest-api-2-status-get
func (s *StatusService) GetAllStatusesWithContext(ctx context.Context) ([]Status, *Response, error) {
	apiEndpoint := "rest/api/2/status"
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)

	if err != nil {
		return nil, nil, err
	}

	statusList := []Status{}
	resp, err := s.client.Do(req, &statusList)
	if err != nil {
		return nil, resp, NewJiraError(resp, err)
	}

	return statusList, resp, nil
}

// GetAllStatuses wraps GetAllStatusesWithContext using the background context.
func (s *StatusService) GetAllStatuses() ([]Status, *Response, error) {
	return s.GetAllStatusesWithContext(context.Background())
}
