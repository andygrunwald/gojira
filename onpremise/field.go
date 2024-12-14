package onpremise

import (
	"context"
	"fmt"
	"net/http"
)

// FieldService handles fields for the Jira instance / API.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/#api-Field
type FieldService service

// Field represents a field of a Jira issue.
type Field struct {
	ID          string      `json:"id,omitempty" structs:"id,omitempty"`
	Key         string      `json:"key,omitempty" structs:"key,omitempty"`
	Name        string      `json:"name,omitempty" structs:"name,omitempty"`
	Custom      bool        `json:"custom,omitempty" structs:"custom,omitempty"`
	Navigable   bool        `json:"navigable,omitempty" structs:"navigable,omitempty"`
	Searchable  bool        `json:"searchable,omitempty" structs:"searchable,omitempty"`
	ClauseNames []string    `json:"clauseNames,omitempty" structs:"clauseNames,omitempty"`
	Schema      FieldSchema `json:"schema,omitempty" structs:"schema,omitempty"`
}

// FieldSchema represents a schema of a Jira field.
// Documentation: https://developer.atlassian.com/cloud/jira/platform/rest/v2/api-group-issue-fields/#api-rest-api-2-field-get
type FieldSchema struct {
	Type     string `json:"type,omitempty" structs:"type,omitempty"`
	Items    string `json:"items,omitempty" structs:"items,omitempty"`
	Custom   string `json:"custom,omitempty" structs:"custom,omitempty"`
	System   string `json:"system,omitempty" structs:"system,omitempty"`
	CustomID int64  `json:"customId,omitempty" structs:"customId,omitempty"`
}

// GetList gets all fields from Jira
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/#api-api-2-field-get
//
// TODO Double check this method if this works as expected, is using the latest API and the response is complete
// This double check effort is done for v2 - Remove this two lines if this is completed.
func (s *FieldService) GetList(ctx context.Context) ([]Field, *Response, error) {
	apiEndpoint := "rest/api/2/field"
	req, err := s.client.NewRequest(ctx, http.MethodGet, apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	fieldList := []Field{}
	resp, err := s.client.Do(req, &fieldList)
	if err != nil {
		return nil, resp, NewJiraError(resp, err)
	}
	return fieldList, resp, nil
}

// FieldCreateOptions are passed to the Field.Create function to create a new Jira cusotm field
type FieldCreateOptions struct {
	// Name: The name of the custom field, which is displayed in Jira. This is not the unique identifier.
	// Required when creating a field.
	// Optional when updating a field.
	Name string `json:"name,omitempty" structs:"name,omitempty"`

	// Description: The description for the field.
	// Optional when creating or updating a field.
	Description string `json:"description,omitempty" structs:"description,omitempty"`

	// SearcherKey: The searcher defines the way the field is searched in Jira. For example, com.atlassian.jira.plugin.system.customfieldtypes:grouppickersearcher. The search UI (basic search and JQL search) will display different operations and values for the field, based on the field searcher.
	// Can take the following values:
	//	com.atlassian.jira.plugin.system.customfieldtypes:cascadingselectsearcher
	//  com.atlassian.jira.plugin.system.customfieldtypes:daterange
	//  com.atlassian.jira.plugin.system.customfieldtypes:datetimerange
	//  com.atlassian.jira.plugin.system.customfieldtypes:exactnumber
	//  com.atlassian.jira.plugin.system.customfieldtypes:exacttextsearcher
	//  com.atlassian.jira.plugin.system.customfieldtypes:grouppickersearcher
	//  com.atlassian.jira.plugin.system.customfieldtypes:labelsearcher
	//  com.atlassian.jira.plugin.system.customfieldtypes:multiselectsearcher
	//  com.atlassian.jira.plugin.system.customfieldtypes:numberrange
	//  com.atlassian.jira.plugin.system.customfieldtypes:projectsearcher
	//  com.atlassian.jira.plugin.system.customfieldtypes:textsearcher
	//  com.atlassian.jira.plugin.system.customfieldtypes:userpickergroupsearcher
	//  com.atlassian.jira.plugin.system.customfieldtypes:versionsearcher
	//
	// Optional when creating or updating a field.
	SearcherKey string `json:"searcherKey,omitempty" structs:"searcherKey,omitempty"`

	// Project: The type of the custom field.
	// Can take the following values:
	//  com.atlassian.jira.plugin.system.customfieldtypes:cascadingselect
	//  com.atlassian.jira.plugin.system.customfieldtypes:datepicker
	//  com.atlassian.jira.plugin.system.customfieldtypes:datetime
	//  com.atlassian.jira.plugin.system.customfieldtypes:float
	//  com.atlassian.jira.plugin.system.customfieldtypes:grouppicker
	//  com.atlassian.jira.plugin.system.customfieldtypes:importid
	//  com.atlassian.jira.plugin.system.customfieldtypes:labels
	//  com.atlassian.jira.plugin.system.customfieldtypes:multicheckboxes
	//  com.atlassian.jira.plugin.system.customfieldtypes:multigrouppicker
	//  com.atlassian.jira.plugin.system.customfieldtypes:multiversion
	//  com.atlassian.jira.plugin.system.customfieldtypes:project
	//  com.atlassian.jira.plugin.system.customfieldtypes:radiobuttons
	//  com.atlassian.jira.plugin.system.customfieldtypes:readonlyfield
	//  com.atlassian.jira.plugin.system.customfieldtypes:select
	//  com.atlassian.jira.plugin.system.customfieldtypes:textarea
	//  com.atlassian.jira.plugin.system.customfieldtypes:textfield
	//  com.atlassian.jira.plugin.system.customfieldtypes:url
	//  com.atlassian.jira.plugin.system.customfieldtypes:userpicker
	//  com.atlassian.jira.plugin.system.customfieldtypes:version
	//
	// Required when creating a field.
	// Can't be updated.
	Type string `json:"type,omitempty" structs:"type,omitempty"`
}

// Creates a custom field on Jira
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/api-group-issue-fields/#api-rest-api-2-field-post
func (s *FieldService) CreateCustom(ctx context.Context, options *FieldCreateOptions) (*Field, *Response, error) {
	apiEndpoint := "rest/api/2/field"
	req, err := s.client.NewRequest(ctx, http.MethodPost, apiEndpoint, options)
	if err != nil {
		return nil, nil, err
	}

	component := new(Field)
	resp, err := s.client.Do(req, component)
	if err != nil {
		return nil, resp, NewJiraError(resp, err)
	}

	return component, resp, nil
}

// Updates a custom field on Jira
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/api-group-issue-fields/#api-rest-api-2-field-fieldid-put
func (s *FieldService) UpdateCustom(ctx context.Context, fieldId string, options *FieldCreateOptions) (*Response, error) {
	apiEndpoint := fmt.Sprintf("rest/api/2/field/%s", fieldId)
	req, err := s.client.NewRequest(ctx, http.MethodPut, apiEndpoint, options)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, NewJiraError(resp, err)
	}

	return resp, nil
}

func (s *FieldService) DeleteCustom(ctx context.Context, fieldId string) (*Response, error) {
	apiEndpoint := fmt.Sprintf("rest/api/2/field/%s", fieldId)
	req, err := s.client.NewRequest(ctx, http.MethodDelete, apiEndpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if resp.StatusCode != 303 {
		return resp, NewJiraError(resp, err)
	}

	return resp, nil
}
