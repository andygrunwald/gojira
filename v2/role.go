package jira

import (
	"context"
	"fmt"
)

// RoleService handles roles for the Jira instance / API.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v3/#api-group-Role
type RoleService struct {
	client *Client
}

// GetListWithContext returns a list of all available project roles
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v3/#api-api-3-role-get
func (s *RoleService) GetListWithContext(ctx context.Context) (*[]Role, *Response, error) {
	apiEndpoint := "rest/api/3/role"
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}
	roles := new([]Role)
	resp, err := s.client.Do(req, roles)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}
	return roles, resp, err
}

// GetList wraps GetListWithContext using the background context.
func (s *RoleService) GetList() (*[]Role, *Response, error) {
	return s.GetListWithContext(context.Background())
}

// GetWithContext retreives a single Role from Jira
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v3/#api-api-3-role-id-get
func (s *RoleService) GetWithContext(ctx context.Context, roleID int) (*Role, *Response, error) {
	apiEndpoint := fmt.Sprintf("rest/api/3/role/%d", roleID)
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}
	role := new(Role)
	resp, err := s.client.Do(req, role)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}
	if role.Self == "" {
		return nil, resp, fmt.Errorf("no role with ID %d found", roleID)
	}

	return role, resp, err
}

// Get wraps GetWithContext using the background context.
func (s *RoleService) Get(roleID int) (*Role, *Response, error) {
	return s.GetWithContext(context.Background(), roleID)
}
