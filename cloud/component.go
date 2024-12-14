package cloud

import (
	"context"
	"fmt"
	"net/http"
)

// ComponentService represents project components.
// Use it to get, create, update, and delete project components.
// Also get components for project and get a count of issues by component.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-project-components/#api-group-project-components
type ComponentService service

const (
	AssigneeTypeProjectLead    = "PROJECT_LEAD"
	AssigneeTypeComponentLead  = "COMPONENT_LEAD"
	AssigneeTypeUnassigned     = "UNASSIGNED"
	AssigneeTypeProjectDefault = "PROJECT_DEFAULT"
)

// ComponentCreateOptions are passed to the ComponentService.Create function to create a new Jira component
type ComponentCreateOptions struct {
	// Name: The unique name for the component in the project.
	// Required when creating a component.
	// Optional when updating a component.
	// The maximum length is 255 characters.
	Name string `json:"name,omitempty" structs:"name,omitempty"`

	// Description: The description for the component.
	// Optional when creating or updating a component.
	Description string `json:"description,omitempty" structs:"description,omitempty"`

	// LeadAccountId: The accountId of the component's lead user.
	// The accountId uniquely identifies the user across all Atlassian products.
	// For example, 5b10ac8d82e05b22cc7d4ef5.
	LeadAccountId string `json:"leadAccountId,omitempty" structs:"leadAccountId,omitempty"`

	// AssigneeType: The nominal user type used to determine the assignee for issues created with this component.
	// Can take the following values:
	//	PROJECT_LEAD the assignee to any issues created with this component is nominally the lead for the project the component is in.
	//	COMPONENT_LEAD the assignee to any issues created with this component is nominally the lead for the component.
	//	UNASSIGNED an assignee is not set for issues created with this component.
	//	PROJECT_DEFAULT the assignee to any issues created with this component is nominally the default assignee for the project that the component is in.
	//
	// Default value: PROJECT_DEFAULT.
	// Optional when creating or updating a component.
	AssigneeType string `json:"assigneeType,omitempty" structs:"assigneeType,omitempty"`

	// Project: The key of the project the component is assigned to.
	// Required when creating a component.
	// Can't be updated.
	Project string `json:"project,omitempty" structs:"project,omitempty"`
}

// Create creates a component.
// Use components to provide containers for issues within a project.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-project-components/#api-rest-api-3-component-post
func (s *ComponentService) Create(ctx context.Context, options *ComponentCreateOptions) (*ProjectComponent, *Response, error) {
	apiEndpoint := "rest/api/3/component"
	req, err := s.client.NewRequest(ctx, http.MethodPost, apiEndpoint, options)
	if err != nil {
		return nil, nil, err
	}

	component := new(ProjectComponent)
	resp, err := s.client.Do(req, component)
	if err != nil {
		return nil, resp, NewJiraError(resp, err)
	}

	return component, resp, nil
}

// Get returns a component for the given componentID.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-project-components/#api-rest-api-3-component-id-get
func (s *ComponentService) Get(ctx context.Context, componentID string) (*ProjectComponent, *Response, error) {
	apiEndpoint := fmt.Sprintf("rest/api/3/component/%s", componentID)
	req, err := s.client.NewRequest(ctx, http.MethodGet, apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	component := new(ProjectComponent)
	resp, err := s.client.Do(req, component)
	if err != nil {
		return nil, resp, NewJiraError(resp, err)
	}

	return component, resp, nil
}

// Update updates an existing component.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-project-components/#api-rest-api-3-component-id-put
func (s *ComponentService) Update(ctx context.Context, component *ProjectComponent) (*ProjectComponent, *Response, error) {
	apiEndpoint := fmt.Sprintf("rest/api/3/component/%s", component.ID)
	req, err := s.client.NewRequest(ctx, http.MethodPut, apiEndpoint, component)
	if err != nil {
		return nil, nil, err
	}

	// Returning the same pointer here is pointless, we return the component from the response instead.
	updatedComponent := new(ProjectComponent)
	resp, err := s.client.Do(req, updatedComponent)
	if err != nil {
		return nil, resp, NewJiraError(resp, err)
	}

	return updatedComponent, resp, nil
}

// TODO Add "Delete component" method. See https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-project-components/#api-rest-api-3-component-id-delete

// TODO Add "Get component issues count" method. See https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-project-components/#api-rest-api-3-component-id-relatedissuecounts-get

// TODO Add "Get project components paginated" method. See https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-project-components/#api-rest-api-3-project-projectidorkey-component-get

// TODO Add "Get project components" method. See https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-project-components/#api-rest-api-3-project-projectidorkey-components-get
