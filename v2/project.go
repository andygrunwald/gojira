package jira

import (
	"context"
	"fmt"

	"github.com/google/go-querystring/query"
)

// ProjectService handles projects for the Jira instance / API.
//
// Jira API docs: https://docs.atlassian.com/jira/REST/latest/#api/2/project
type ProjectService struct {
	client *Client
}

// GetListWithContext gets all projects form Jira
//
// Jira API docs: https://docs.atlassian.com/jira/REST/latest/#api/2/project-getAllProjects
func (s *ProjectService) GetListWithContext(ctx context.Context) (*ProjectList, *Response, error) {
	return s.ListWithOptionsWithContext(ctx, &GetQueryOptions{})
}

// GetList wraps GetListWithContext using the background context.
func (s *ProjectService) GetList() (*ProjectList, *Response, error) {
	return s.GetListWithContext(context.Background())
}

// ListWithOptionsWithContext gets all projects form Jira with optional query params, like &GetQueryOptions{Expand: "issueTypes"} to get
// a list of all projects and their supported issuetypes
//
// Jira API docs: https://docs.atlassian.com/jira/REST/latest/#api/2/project-getAllProjects
func (s *ProjectService) ListWithOptionsWithContext(ctx context.Context, options *GetQueryOptions) (*ProjectList, *Response, error) {
	apiEndpoint := "rest/api/2/project"
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	if options != nil {
		q, err := query.Values(options)
		if err != nil {
			return nil, nil, err
		}
		req.URL.RawQuery = q.Encode()
	}

	projectList := new(ProjectList)
	resp, err := s.client.Do(req, projectList)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}

	return projectList, resp, nil
}

// ListWithOptions wraps ListWithOptionsWithContext using the background context.
func (s *ProjectService) ListWithOptions(options *GetQueryOptions) (*ProjectList, *Response, error) {
	return s.ListWithOptionsWithContext(context.Background(), options)
}

// GetWithContext returns a full representation of the project for the given issue key.
// Jira will attempt to identify the project by the projectIdOrKey path parameter.
// This can be an project id, or an project key.
//
// Jira API docs: https://docs.atlassian.com/jira/REST/latest/#api/2/project-getProject
func (s *ProjectService) GetWithContext(ctx context.Context, projectID string) (*Project, *Response, error) {
	apiEndpoint := fmt.Sprintf("rest/api/2/project/%s", projectID)
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	project := new(Project)
	resp, err := s.client.Do(req, project)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}

	return project, resp, nil
}

// Get wraps GetWithContext using the background context.
func (s *ProjectService) Get(projectID string) (*Project, *Response, error) {
	return s.GetWithContext(context.Background(), projectID)
}

// GetPermissionSchemeWithContext returns a full representation of the permission scheme for the project
// Jira will attempt to identify the project by the projectIdOrKey path parameter.
// This can be an project id, or an project key.
//
// Jira API docs: https://docs.atlassian.com/jira/REST/latest/#api/2/project-getProject
func (s *ProjectService) GetPermissionSchemeWithContext(ctx context.Context, projectID string) (*PermissionScheme, *Response, error) {
	apiEndpoint := fmt.Sprintf("/rest/api/2/project/%s/permissionscheme", projectID)
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	ps := new(PermissionScheme)
	resp, err := s.client.Do(req, ps)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}

	return ps, resp, nil
}

// GetPermissionScheme wraps GetPermissionSchemeWithContext using the background context.
func (s *ProjectService) GetPermissionScheme(projectID string) (*PermissionScheme, *Response, error) {
	return s.GetPermissionSchemeWithContext(context.Background(), projectID)
}
