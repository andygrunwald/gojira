package jira

import "context"

// StatusCategoryService handles status categories for the Jira instance / API.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/#api-Statuscategory
type StatusCategoryService struct {
	client *Client
}

// GetListWithContext gets all status categories from Jira
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/#api-api-2-statuscategory-get
func (s *StatusCategoryService) GetListWithContext(ctx context.Context) ([]StatusCategory, *Response, error) {
	apiEndpoint := "rest/api/2/statuscategory"
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	statusCategoryList := []StatusCategory{}
	resp, err := s.client.Do(req, &statusCategoryList)
	if err != nil {
		return nil, resp, NewJiraError(resp, err)
	}
	return statusCategoryList, resp, nil
}

// GetList wraps GetListWithContext using the background context.
func (s *StatusCategoryService) GetList() ([]StatusCategory, *Response, error) {
	return s.GetListWithContext(context.Background())
}
