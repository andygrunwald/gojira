package jira

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/andygrunwald/go-jira/v2/errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	// AssigneeAutomatic represents the value of the "Assignee: Automatic" of Jira
	AssigneeAutomatic = "-1"
)

// IssueService handles Issues for the Jira instance / API.
//
// Jira API docs: https://docs.atlassian.com/jira/REST/latest/#api/2/issue
type IssueService struct {
	client *Client
}

// UpdateQueryOptions specifies the optional parameters to the Edit issue
type UpdateQueryOptions struct {
	NotifyUsers            bool `url:"notifyUsers,omitempty"`
	OverrideScreenSecurity bool `url:"overrideScreenSecurity,omitempty"`
	OverrideEditableFlag   bool `url:"overrideEditableFlag,omitempty"`
}

// UnmarshalJSON is a custom JSON marshal function for the IssueFields structs.
// It handles Jira custom fields and maps those from / to "Unknowns" key.
func (i *IssueFields) UnmarshalJSON(data []byte) error {


	// We have to unmarshal into a map in order to capture the custom fields
	// TODO performance might be able to be increased by unmarshalling only once
	tempMap := make(map[string]interface{})
	err := json.Unmarshal(data, &tempMap)
	if err != nil {
		return fmt.Errorf("failed to unmarshal json: %w", err)
	}

	for k, jsonv := range tempMap {
		//if strings.HasPrefix(k, "customfield_") {
		//
		//}

		// https://golang.org/pkg/encoding/json/#Unmarshal
		switch _ := jsonv.(type) {
		case map[string]interface{}:
			possibleTypes := []type{
			CustomFieldSelect
		}

			if v, ok := jsonv.(CustomFieldSelect); ok {
				issueFields.CustomFields[k] = v
			}
		}
	}

	return nil

}























// CreateTransitionPayload is used for creating new issue transitions
type CreateTransitionPayload struct {
	Transition TransitionPayload       `json:"transition" structs:"transition"`
	Fields     TransitionPayloadFields `json:"fields" structs:"fields"`
}

// TransitionPayload represents the request payload of Transition calls like DoTransition
type TransitionPayload struct {
	ID string `json:"id" structs:"id"`
}

// TransitionPayloadFields represents the fields that can be set when executing a transition
type TransitionPayloadFields struct {
	Resolution *Resolution `json:"resolution,omitempty" structs:"resolution,omitempty"`
}

// Option represents an option value in a SelectList or MultiSelect
// custom issue field
type Option struct {
	Value string `json:"value" structs:"value"`
}





type EntityProperty struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}





// SearchOptions specifies the optional parameters to various List methods that
// support pagination.
// Pagination is used for the Jira REST APIs to conserve server resources and limit
// response size for resources that return potentially large collection of items.
// A request to a pages API will result in a values array wrapped in a JSON object with some paging metadata
// Default Pagination options
type SearchOptions struct {
	// StartAt: The starting index of the returned projects. Base index: 0.
	StartAt int `url:"startAt,omitempty"`
	// MaxResults: The maximum number of projects to return per page. Default: 50.
	MaxResults int `url:"maxResults"`
	// Expand: Expand specific sections in the returned issues
	Expand string `url:"expand,omitempty"`
	Fields []string
	// ValidateQuery: The validateQuery param offers control over whether to validate and how strictly to treat the validation. Default: strict.
	ValidateQuery string `url:"validateQuery,omitempty"`
}

// searchResult is only a small wrapper around the Search (with JQL) method
// to be able to parse the results
type SearchResult struct {
	Issues     []Issue `json:"issues" structs:"issues"`
	StartAt    int     `json:"startAt" structs:"startAt"`
	MaxResults int     `json:"maxResults" structs:"maxResults"`
	Total      int     `json:"total" structs:"total"`
}

// GetWorklogsQueryOptions specifies the optional parameters for the Get Worklogs method
type GetWorklogsQueryOptions struct {
	StartAt    int64  `url:"startAt,omitempty"`
	MaxResults int32  `url:"maxResults,omitempty"`
	Expand     string `url:"expand,omitempty"`
}

type AddWorklogQueryOptions struct {
	NotifyUsers          *bool  `url:"notifyUsers,omitempty"`
	AdjustEstimate       string `url:"adjustEstimate,omitempty"`
	NewEstimate          string `url:"newEstimate,omitempty"`
	ReduceBy             string `url:"reduceBy,omitempty"`
	Expand               string `url:"expand,omitempty"`
	OverrideEditableFlag *bool  `url:"overrideEditableFlag,omitempty"`
}

// CustomFields represents custom fields of Jira
// This can heavily differ between Jira instances
type CustomFields map[string]string


// GetQueryOptions specifies the optional parameters for the Get Issue methods
type GetQueryOptions struct {
	// Fields is the list of fields to return for the issue. By default, all fields are returned.
	Fields string `url:"fields,omitempty"`
	Expand string `url:"expand,omitempty"`
	// Properties is the list of properties to return for the issue. By default no properties are returned.
	Properties string `url:"properties,omitempty"`
	// FieldsByKeys if true then fields in issues will be referenced by keys instead of ids
	FieldsByKeys  bool   `url:"fieldsByKeys,omitempty"`
	UpdateHistory bool   `url:"updateHistory,omitempty"`
	ProjectKeys   string `url:"projectKeys,omitempty"`
}

type GetOptions struct {
	Context        context.Context
	RequestOptions []HttpRequestOption
	QueryOptions   *GetQueryOptions
}

type GetOption func(opt *GetOptions)

// GetWithContext returns a full representation of the issue for the given issue key.
// Jira will attempt to identify the issue by the issueIdOrKey path parameter.
// This can be an issue id, or an issue key.
// If the issue cannot be found via an exact match, Jira will also look for the issue in a case-insensitive way, or by looking to see if the issue was moved.
//
// The given options will be appended to the query string
//
// Jira API docs: https://docs.atlassian.com/jira/REST/latest/#api/2/issue-getIssue

// Get wraps GetWithContext using the background context.
func (s *IssueService) Get(issueID string, options ...GetOption) (*Issue, *Response, error) {
	o := &GetOptions{
		Context:      context.Background(),
		QueryOptions: &GetQueryOptions{},
	}

	for _, opt := range options {
		opt(o)
	}

	apiEndpoint := fmt.Sprintf("rest/api/2/issue/%s", issueID)
	req, err := s.client.NewRequest(o.Context, http.MethodGet, apiEndpoint, nil)
	if err != nil {
		return nil, nil, fmt.Errorf(errors.HttpRequestCreationFailureMessageFormat, err)
	}

	q, err := query.Values(o.QueryOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create http query: %w", err)
	}
	req.URL.RawQuery = q.Encode()

	issue := new(Issue)
	value, resp, err := s.client.Do(req)
	return issue, resp, err
}

type DownloadAttachmentOptions struct {
	Context context.Context
}

type DownloadAttachmentOption func(opt *DownloadAttachmentOptions)

// DownloadAttachmentWithContext returns a Response of an attachment for a given attachmentID.
// The attachment is in the Response.Body of the response.
// This is an io.ReadCloser.
// The caller should close the resp.Body.
func (s *IssueService) DownloadAttachment(attachmentID string, options ...DownloadAttachmentOption) (*Response, error) {
	o := &DownloadAttachmentOptions{
		Context: context.Background(),
	}

	for _, opt := range options {
		opt(o)
	}

	apiEndpoint := fmt.Sprintf("secure/attachment/%s/", attachmentID)
	req, err := s.client.NewRequestWithContext(o.Context, http.MethodGet, apiEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf(errors.HttpRequestCreationFailureMessageFormat, err)
	}

	return s.client.Do(req, nil)
}

type PostAttachmentOptions struct {
	Context context.Context
}

type PostAttachmentOption func(opt *PostAttachmentOptions)

// PostAttachmentWithContext uploads r (io.Reader) as an attachment to a given issueID
func (s *IssueService) PostAttachmentWithContext(issueID string, r io.Reader, attachmentName string, options ...PostAttachmentOption) (*[]Attachment, *Response, error) {
	o := &PostAttachmentOptions{
		Context: context.Background(),
	}

	for _, opt := range options {
		opt(o)
	}

	apiEndpoint := fmt.Sprintf("rest/api/2/issue/%s/attachments", issueID)

	b := new(bytes.Buffer)
	writer := multipart.NewWriter(b)

	fw, err := writer.CreateFormFile("file", attachmentName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create form file: %w", err)
	}

	if r != nil {
		if _, err = io.Copy(fw, r); err != nil {
			return nil, nil, fmt.Errorf("failed to copy form file contents to reader: %w", err)
		}
	}
	writer.Close()

	req, err := s.client.NewMultiPartRequestWithContext(o.Context, http.MethodPost, apiEndpoint, b)
	if err != nil {
		return nil, nil, fmt.Errorf(errors.HttpRequestCreationFailureMessageFormat, err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// PostAttachment response returns a JSON array (as multiple attachments can be posted)
	attachment := new([]Attachment)
	resp, err := s.client.Do(req, attachment)
	return attachment, resp, err
}

type DeleteAttachmentOptions struct {
	Context context.Context
}

type DeleteAttachmentOption func(opt *DeleteAttachmentOptions)

// DeleteAttachmentWithContext deletes an attachment of a given attachmentID
func (s *IssueService) DeleteAttachmentWithContext(attachmentID string, options ...DeleteAttachmentOption) (*Response, error) {
	o := &DeleteAttachmentOptions{
		Context: context.Background(),
	}

	for _, opt := range options {
		opt(o)
	}

	apiEndpoint := fmt.Sprintf("rest/api/2/attachment/%s", attachmentID)

	req, err := s.client.NewRequestWithContext(o.Context, http.MethodDelete, apiEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf(errors.HttpRequestCreationFailureMessageFormat, err)
	}

	return s.client.Do(req, nil)
}

type GetWorklogsOptions struct {
	Context            context.Context
	HttpRequestOptions []func(r *http.Request) error
}

type GetWorklogsOption func(opt *GetWorklogsOptions)

// GetWorklogsWithContext gets all the worklogs for an issue.
// This method is especially important if you need to read all the worklogs, not just the first page.
//
// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue/{issueIdOrKey}/worklog-getIssueWorklog
func (s *IssueService) GetWorklogs(issueID string, options ...GetWorklogsOption) (*Worklog, *Response, error) {
	o := &GetWorklogsOptions{
		Context: context.Background(),
	}

	for _, opt := range options {
		opt(o)
	}

	apiEndpoint := fmt.Sprintf("rest/api/2/issue/%s/worklog", issueID)

	req, err := s.client.NewRequestWithContext(o.Context, http.MethodGet, apiEndpoint, nil)
	if err != nil {
		return nil, nil, fmt.Errorf(errors.HttpRequestCreationFailureMessageFormat, err)
	}

	for _, reqOpt := range o.HttpRequestOptions {
		err = reqOpt(req)
		if err != nil {
			return nil, nil, err
		}
	}

	v := new(Worklog)
	resp, err := s.client.Do(req, v)
	return v, resp, err
}

// Applies query options to http request.
// This helper is meant to be used with all "QueryOptions" structs.
func WithQueryOptions(options interface{}) func(*http.Request) error {
	q, err := query.Values(options)
	if err != nil {
		return func(*http.Request) error {
			return err
		}
	}

	return func(r *http.Request) error {
		r.URL.RawQuery = q.Encode()
		return nil
	}
}

// CreateWithContext creates an issue or a sub-task from a JSON representation.
// Creating a sub-task is similar to creating a regular issue, with two important differences:
// The issueType field must correspond to a sub-task issue type and you must provide a parent field in the issue create request containing the id or key of the parent issue.
//
// Jira API docs: https://docs.atlassian.com/jira/REST/latest/#api/2/issue-createIssues
func (s *IssueService) CreateWithContext(ctx context.Context, issue *Issue) (*Issue, *Response, error) {
	apiEndpoint := "rest/api/2/issue"
	req, err := s.client.NewRequestWithContext(ctx, "POST", apiEndpoint, issue)
	if err != nil {
		return nil, nil, err
	}
	resp, err := s.client.Do(req, nil)
	if err != nil {
		// incase of error return the resp for further inspection
		return nil, resp, err
	}

	responseIssue := new(Issue)
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, fmt.Errorf("could not read the returned data")
	}
	err = json.Unmarshal(data, responseIssue)
	if err != nil {
		return nil, resp, fmt.Errorf("could not unmarshall the data into struct")
	}
	return responseIssue, resp, nil
}

// Create wraps CreateWithContext using the background context.
func (s *IssueService) Create(issue *Issue) (*Issue, *Response, error) {
	return s.CreateWithContext(context.Background(), issue)
}

// UpdateWithOptionsWithContext updates an issue from a JSON representation,
// while also specifying query params. The issue is found by key.
//
// Jira API docs: https://docs.atlassian.com/jira/REST/cloud/#api/2/issue-editIssue
func (s *IssueService) UpdateWithOptionsWithContext(ctx context.Context, issue *Issue, opts *UpdateQueryOptions) (*Issue, *Response, error) {
	apiEndpoint := fmt.Sprintf("rest/api/2/issue/%v", issue.Key)
	url, err := addOptions(apiEndpoint, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequestWithContext(ctx, "PUT", url, issue)
	if err != nil {
		return nil, nil, err
	}
	resp, err := s.client.Do(req, nil)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}

	// This is just to follow the rest of the API's convention of returning an issue.
	// Returning the same pointer here is pointless, so we return a copy instead.
	ret := *issue
	return &ret, resp, nil
}

// UpdateWithOptions wraps UpdateWithOptionsWithContext using the background context.
func (s *IssueService) UpdateWithOptions(issue *Issue, opts *UpdateQueryOptions) (*Issue, *Response, error) {
	return s.UpdateWithOptionsWithContext(context.Background(), issue, opts)
}

// UpdateWithContext updates an issue from a JSON representation. The issue is found by key.
//
// Jira API docs: https://docs.atlassian.com/jira/REST/cloud/#api/2/issue-editIssue
func (s *IssueService) UpdateWithContext(ctx context.Context, issue *Issue) (*Issue, *Response, error) {
	return s.UpdateWithOptions(issue, nil)
}

// Update wraps UpdateWithContext using the background context.
func (s *IssueService) Update(issue *Issue) (*Issue, *Response, error) {
	return s.UpdateWithContext(context.Background(), issue)
}

// UpdateIssueWithContext updates an issue from a JSON representation. The issue is found by key.
//
// https://docs.atlassian.com/jira/REST/7.4.0/#api/2/issue-editIssue
func (s *IssueService) UpdateIssueWithContext(ctx context.Context, jiraID string, data map[string]interface{}) (*Response, error) {
	apiEndpoint := fmt.Sprintf("rest/api/2/issue/%v", jiraID)
	req, err := s.client.NewRequestWithContext(ctx, "PUT", apiEndpoint, data)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	// This is just to follow the rest of the API's convention of returning an issue.
	// Returning the same pointer here is pointless, so we return a copy instead.
	return resp, nil
}

// UpdateIssue wraps UpdateIssueWithContext using the background context.
func (s *IssueService) UpdateIssue(jiraID string, data map[string]interface{}) (*Response, error) {
	return s.UpdateIssueWithContext(context.Background(), jiraID, data)
}

// AddCommentWithContext adds a new comment to issueID.
//
// Jira API docs: https://docs.atlassian.com/jira/REST/latest/#api/2/issue-addComment
func (s *IssueService) AddCommentWithContext(ctx context.Context, issueID string, comment *Comment) (*Comment, *Response, error) {
	apiEndpoint := fmt.Sprintf("rest/api/2/issue/%s/comment", issueID)
	req, err := s.client.NewRequestWithContext(ctx, "POST", apiEndpoint, comment)
	if err != nil {
		return nil, nil, err
	}

	responseComment := new(Comment)
	resp, err := s.client.Do(req, responseComment)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}

	return responseComment, resp, nil
}

// AddComment wraps AddCommentWithContext using the background context.
func (s *IssueService) AddComment(issueID string, comment *Comment) (*Comment, *Response, error) {
	return s.AddCommentWithContext(context.Background(), issueID, comment)
}

// UpdateCommentWithContext updates the body of a comment, identified by comment.ID, on the issueID.
//
// Jira API docs: https://docs.atlassian.com/jira/REST/cloud/#api/2/issue/{issueIdOrKey}/comment-updateComment
func (s *IssueService) UpdateCommentWithContext(ctx context.Context, issueID string, comment *Comment) (*Comment, *Response, error) {
	reqBody := struct {
		Body string `json:"body"`
	}{
		Body: comment.Body,
	}
	apiEndpoint := fmt.Sprintf("rest/api/2/issue/%s/comment/%s", issueID, comment.ID)
	req, err := s.client.NewRequestWithContext(ctx, "PUT", apiEndpoint, reqBody)
	if err != nil {
		return nil, nil, err
	}

	responseComment := new(Comment)
	resp, err := s.client.Do(req, responseComment)
	if err != nil {
		return nil, resp, err
	}

	return responseComment, resp, nil
}

// UpdateComment wraps UpdateCommentWithContext using the background context.
func (s *IssueService) UpdateComment(issueID string, comment *Comment) (*Comment, *Response, error) {
	return s.UpdateCommentWithContext(context.Background(), issueID, comment)
}

// DeleteCommentWithContext Deletes a comment from an issueID.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v3/#api-api-3-issue-issueIdOrKey-comment-id-delete
func (s *IssueService) DeleteCommentWithContext(ctx context.Context, issueID, commentID string) error {
	apiEndpoint := fmt.Sprintf("rest/api/2/issue/%s/comment/%s", issueID, commentID)
	req, err := s.client.NewRequestWithContext(ctx, "DELETE", apiEndpoint, nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return jerr
	}

	return nil
}

// DeleteComment wraps DeleteCommentWithContext using the background context.
func (s *IssueService) DeleteComment(issueID, commentID string) error {
	return s.DeleteCommentWithContext(context.Background(), issueID, commentID)
}

// AddWorklogRecordWithContext adds a new worklog record to issueID.
//
// https://developer.atlassian.com/cloud/jira/platform/rest/#api-api-2-issue-issueIdOrKey-worklog-post
func (s *IssueService) AddWorklogRecordWithContext(ctx context.Context, issueID string, record *WorklogRecord, options ...func(*http.Request) error) (*WorklogRecord, *Response, error) {
	apiEndpoint := fmt.Sprintf("rest/api/2/issue/%s/worklog", issueID)
	req, err := s.client.NewRequestWithContext(ctx, "POST", apiEndpoint, record)
	if err != nil {
		return nil, nil, err
	}

	for _, option := range options {
		err = option(req)
		if err != nil {
			return nil, nil, err
		}
	}

	responseRecord := new(WorklogRecord)
	resp, err := s.client.Do(req, responseRecord)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}

	return responseRecord, resp, nil
}

// AddWorklogRecord wraps AddWorklogRecordWithContext using the background context.
func (s *IssueService) AddWorklogRecord(issueID string, record *WorklogRecord, options ...func(*http.Request) error) (*WorklogRecord, *Response, error) {
	return s.AddWorklogRecordWithContext(context.Background(), issueID, record, options...)
}

// UpdateWorklogRecordWithContext updates a worklog record.
//
// https://docs.atlassian.com/software/jira/docs/api/REST/7.1.2/#api/2/issue-updateWorklog
func (s *IssueService) UpdateWorklogRecordWithContext(ctx context.Context, issueID, worklogID string, record *WorklogRecord, options ...func(*http.Request) error) (*WorklogRecord, *Response, error) {
	apiEndpoint := fmt.Sprintf("rest/api/2/issue/%s/worklog/%s", issueID, worklogID)
	req, err := s.client.NewRequestWithContext(ctx, "PUT", apiEndpoint, record)
	if err != nil {
		return nil, nil, err
	}

	for _, option := range options {
		err = option(req)
		if err != nil {
			return nil, nil, err
		}
	}

	responseRecord := new(WorklogRecord)
	resp, err := s.client.Do(req, responseRecord)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}

	return responseRecord, resp, nil
}

// UpdateWorklogRecord wraps UpdateWorklogRecordWithContext using the background context.
func (s *IssueService) UpdateWorklogRecord(issueID, worklogID string, record *WorklogRecord, options ...func(*http.Request) error) (*WorklogRecord, *Response, error) {
	return s.UpdateWorklogRecordWithContext(context.Background(), issueID, worklogID, record, options...)
}

// AddLinkWithContext adds a link between two issues.
//
// Jira API docs: https://docs.atlassian.com/jira/REST/latest/#api/2/issueLink
func (s *IssueService) AddLinkWithContext(ctx context.Context, issueLink *IssueLink) (*Response, error) {
	apiEndpoint := "rest/api/2/issueLink"
	req, err := s.client.NewRequestWithContext(ctx, "POST", apiEndpoint, issueLink)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		err = NewJiraError(resp, err)
	}

	return resp, err
}

// AddLink wraps AddLinkWithContext using the background context.
func (s *IssueService) AddLink(issueLink *IssueLink) (*Response, error) {
	return s.AddLinkWithContext(context.Background(), issueLink)
}

// SearchWithContext will search for tickets according to the jql
//
// Jira API docs: https://developer.atlassian.com/jiradev/jira-apis/jira-rest-apis/jira-rest-api-tutorials/jira-rest-api-example-query-issues
func (s *IssueService) SearchWithContext(ctx context.Context, jql string, options *SearchOptions) ([]Issue, *Response, error) {
	u := url.URL{
		Path: "rest/api/2/search",
	}
	uv := url.Values{}
	if jql != "" {
		uv.Add("jql", jql)
	}

	if options != nil {
		if options.StartAt != 0 {
			uv.Add("startAt", strconv.Itoa(options.StartAt))
		}
		if options.MaxResults != 0 {
			uv.Add("maxResults", strconv.Itoa(options.MaxResults))
		}
		if options.Expand != "" {
			uv.Add("expand", options.Expand)
		}
		if strings.Join(options.Fields, ",") != "" {
			uv.Add("fields", strings.Join(options.Fields, ","))
		}
		if options.ValidateQuery != "" {
			uv.Add("validateQuery", options.ValidateQuery)
		}
	}

	u.RawQuery = uv.Encode()

	req, err := s.client.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return []Issue{}, nil, err
	}

	v := new(searchResult)
	resp, err := s.client.Do(req, v)
	if err != nil {
		err = NewJiraError(resp, err)
	}
	return v.Issues, resp, err
}

// Search wraps SearchWithContext using the background context.
func (s *IssueService) Search(jql string, options *SearchOptions) ([]Issue, *Response, error) {
	return s.SearchWithContext(context.Background(), jql, options)
}

// SearchPagesWithContext will get issues from all pages in a search
//
// Jira API docs: https://developer.atlassian.com/jiradev/jira-apis/jira-rest-apis/jira-rest-api-tutorials/jira-rest-api-example-query-issues
func (s *IssueService) SearchPagesWithContext(ctx context.Context, jql string, options *SearchOptions, f func(Issue) error) error {
	if options == nil {
		options = &SearchOptions{
			StartAt:    0,
			MaxResults: 50,
		}
	}

	if options.MaxResults == 0 {
		options.MaxResults = 50
	}

	issues, resp, err := s.Search(jql, options)
	if err != nil {
		return err
	}

	if len(issues) == 0 {
		return nil
	}

	for {
		for _, issue := range issues {
			err = f(issue)
			if err != nil {
				return err
			}
		}

		if resp.StartAt+resp.MaxResults >= resp.Total {
			return nil
		}

		options.StartAt += resp.MaxResults
		issues, resp, err = s.Search(jql, options)
		if err != nil {
			return err
		}
	}
}

// SearchPages wraps SearchPagesWithContext using the background context.
func (s *IssueService) SearchPages(jql string, options *SearchOptions, f func(Issue) error) error {
	return s.SearchPagesWithContext(context.Background(), jql, options, f)
}

// GetCustomFieldsWithContext returns a map of customfield_* keys with string values
func (s *IssueService) GetCustomFieldsWithContext(ctx context.Context, issueID string) (CustomFields, *Response, error) {
	apiEndpoint := fmt.Sprintf("rest/api/2/issue/%s", issueID)
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	issue := new(map[string]interface{})
	resp, err := s.client.Do(req, issue)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}

	m := *issue
	f := m["fields"]
	cf := make(CustomFields)
	if f == nil {
		return cf, resp, nil
	}

	if rec, ok := f.(map[string]interface{}); ok {
		for key, val := range rec {
			if strings.Contains(key, "customfield") {
				if valMap, ok := val.(map[string]interface{}); ok {
					if v, ok := valMap["value"]; ok {
						val = v
					}
				}
				cf[key] = fmt.Sprint(val)
			}
		}
	}
	return cf, resp, nil
}

// GetCustomFields wraps GetCustomFieldsWithContext using the background context.
func (s *IssueService) GetCustomFields(issueID string) (CustomFields, *Response, error) {
	return s.GetCustomFieldsWithContext(context.Background(), issueID)
}

// GetTransitionsWithContext gets a list of the transitions possible for this issue by the current user,
// along with fields that are required and their types.
//
// Jira API docs: https://docs.atlassian.com/jira/REST/latest/#api/2/issue-getTransitions
func (s *IssueService) GetTransitionsWithContext(ctx context.Context, id string) ([]Transition, *Response, error) {
	apiEndpoint := fmt.Sprintf("rest/api/2/issue/%s/transitions?expand=transitions.fields", id)
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(transitionResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		err = NewJiraError(resp, err)
	}
	return result.Transitions, resp, err
}

// GetTransitions wraps GetTransitionsWithContext using the background context.
func (s *IssueService) GetTransitions(id string) ([]Transition, *Response, error) {
	return s.GetTransitionsWithContext(context.Background(), id)
}

// DoTransitionWithContext performs a transition on an issue.
// When performing the transition you can update or set other issue fields.
//
// Jira API docs: https://docs.atlassian.com/jira/REST/latest/#api/2/issue-doTransition
func (s *IssueService) DoTransitionWithContext(ctx context.Context, ticketID, transitionID string) (*Response, error) {
	payload := CreateTransitionPayload{
		Transition: TransitionPayload{
			ID: transitionID,
		},
	}
	return s.DoTransitionWithPayloadWithContext(ctx, ticketID, payload)
}

// DoTransition wraps DoTransitionWithContext using the background context.
func (s *IssueService) DoTransition(ticketID, transitionID string) (*Response, error) {
	return s.DoTransitionWithContext(context.Background(), ticketID, transitionID)
}

// DoTransitionWithPayloadWithContext performs a transition on an issue using any payload.
// When performing the transition you can update or set other issue fields.
//
// Jira API docs: https://docs.atlassian.com/jira/REST/latest/#api/2/issue-doTransition
func (s *IssueService) DoTransitionWithPayloadWithContext(ctx context.Context, ticketID, payload interface{}) (*Response, error) {
	apiEndpoint := fmt.Sprintf("rest/api/2/issue/%s/transitions", ticketID)

	req, err := s.client.NewRequestWithContext(ctx, "POST", apiEndpoint, payload)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		err = NewJiraError(resp, err)
	}

	return resp, err
}

// DoTransitionWithPayload wraps DoTransitionWithPayloadWithContext using the background context.
func (s *IssueService) DoTransitionWithPayload(ticketID, payload interface{}) (*Response, error) {
	return s.DoTransitionWithPayloadWithContext(context.Background(), ticketID, payload)
}

// InitIssueWithMetaAndFields returns Issue with with values from fieldsConfig properly set.
//  * metaProject should contain metaInformation about the project where the issue should be created.
//  * metaIssuetype is the MetaInformation about the Issuetype that needs to be created.
//  * fieldsConfig is a key->value pair where key represents the name of the field as seen in the UI
//		And value is the string value for that particular key.
// Note: This method doesn't verify that the fieldsConfig is complete with mandatory fields. The fieldsConfig is
//		 supposed to be already verified with MetaIssueType.CheckCompleteAndAvailable. It will however return
//		 error if the key is not found.
//		 All values will be packed into Unknowns. This is much convenient. If the struct fields needs to be
//		 configured as well, marshalling and unmarshalling will set the proper fields.
func InitIssueWithMetaAndFields(metaProject *MetaProject, metaIssuetype *MetaIssueType, fieldsConfig map[string]string) (*Issue, error) {
	issue := new(Issue)
	issueFields := new(IssueFields)
	//issueFields.Unknowns = tcontainer.NewMarshalMap()

	// map the field names the User presented to jira's internal key
	allFields, _ := metaIssuetype.GetAllFields()
	for key, value := range fieldsConfig {
		jiraKey, found := allFields[key]
		if !found {
			return nil, fmt.Errorf("key %s is not found in the list of fields", key)
		}

		valueType, err := metaIssuetype.Fields.String(jiraKey + "/schema/type")
		if err != nil {
			return nil, err
		}
		switch valueType {
		case "array":
			elemType, err := metaIssuetype.Fields.String(jiraKey + "/schema/items")
			if err != nil {
				return nil, err
			}
			switch elemType {
			case "component":
				issueFields.Unknowns[jiraKey] = []Component{{Name: value}}
			case "option":
				issueFields.Unknowns[jiraKey] = []map[string]string{{"value": value}}
			default:
				issueFields.Unknowns[jiraKey] = []string{value}
			}
		case "string":
			issueFields.Unknowns[jiraKey] = value
		case "date":
			issueFields.Unknowns[jiraKey] = value
		case "datetime":
			issueFields.Unknowns[jiraKey] = value
		case "any":
			// Treat any as string
			issueFields.Unknowns[jiraKey] = value
		case "project":
			issueFields.Unknowns[jiraKey] = Project{
				Name: metaProject.Name,
				ID:   metaProject.Id,
			}
		case "priority":
			issueFields.Unknowns[jiraKey] = Priority{Name: value}
		case "user":
			issueFields.Unknowns[jiraKey] = User{
				Name: value,
			}
		case "issuetype":
			issueFields.Unknowns[jiraKey] = IssueType{
				Name: value,
			}
		case "option":
			issueFields.Unknowns[jiraKey] = Option{
				Value: value,
			}
		default:
			return nil, fmt.Errorf("unknown issue type encountered: %s for %s", valueType, key)
		}
	}

	issue.Fields = issueFields

	return issue, nil
}

// DeleteWithContext will delete a specified issue.
func (s *IssueService) DeleteWithContext(ctx context.Context, issueID string) (*Response, error) {
	apiEndpoint := fmt.Sprintf("rest/api/2/issue/%s", issueID)

	// to enable deletion of subtasks; without this, the request will fail if the issue has subtasks
	deletePayload := make(map[string]interface{})
	deletePayload["deleteSubtasks"] = "true"
	content, _ := json.Marshal(deletePayload)

	req, err := s.client.NewRequestWithContext(ctx, "DELETE", apiEndpoint, content)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	return resp, err
}

// Delete wraps DeleteWithContext using the background context.
func (s *IssueService) Delete(issueID string) (*Response, error) {
	return s.DeleteWithContext(context.Background(), issueID)
}

// GetWatchersWithContext wil return all the users watching/observing the given issue
//
// Jira API docs: https://docs.atlassian.com/software/jira/docs/api/REST/latest/#api/2/issue-getIssueWatchers
func (s *IssueService) GetWatchersWithContext(ctx context.Context, issueID string) (*[]User, *Response, error) {
	watchesAPIEndpoint := fmt.Sprintf("rest/api/2/issue/%s/watchers", issueID)

	req, err := s.client.NewRequestWithContext(ctx, "GET", watchesAPIEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	watches := new(Watches)
	resp, err := s.client.Do(req, watches)
	if err != nil {
		return nil, nil, NewJiraError(resp, err)
	}

	result := []User{}
	for _, watcher := range watches.Watchers {
		var user *User
		if watcher.AccountID != "" {
			user, resp, err = s.client.User.GetByAccountID(watcher.AccountID)
			if err != nil {
				return nil, resp, NewJiraError(resp, err)
			}
		}
		result = append(result, *user)
	}

	return &result, resp, nil
}

// GetWatchers wraps GetWatchersWithContext using the background context.
func (s *IssueService) GetWatchers(issueID string) (*[]User, *Response, error) {
	return s.GetWatchersWithContext(context.Background(), issueID)
}

// AddWatcherWithContext adds watcher to the given issue
//
// Jira API docs: https://docs.atlassian.com/software/jira/docs/api/REST/latest/#api/2/issue-addWatcher
func (s *IssueService) AddWatcherWithContext(ctx context.Context, issueID string, userName string) (*Response, error) {
	apiEndPoint := fmt.Sprintf("rest/api/2/issue/%s/watchers", issueID)

	req, err := s.client.NewRequestWithContext(ctx, "POST", apiEndPoint, userName)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		err = NewJiraError(resp, err)
	}

	return resp, err
}

// AddWatcher wraps AddWatcherWithContext using the background context.
func (s *IssueService) AddWatcher(issueID string, userName string) (*Response, error) {
	return s.AddWatcherWithContext(context.Background(), issueID, userName)
}

// RemoveWatcherWithContext removes given user from given issue
//
// Jira API docs: https://docs.atlassian.com/software/jira/docs/api/REST/latest/#api/2/issue-removeWatcher
func (s *IssueService) RemoveWatcherWithContext(ctx context.Context, issueID string, userName string) (*Response, error) {
	apiEndPoint := fmt.Sprintf("rest/api/2/issue/%s/watchers", issueID)

	req, err := s.client.NewRequestWithContext(ctx, "DELETE", apiEndPoint, userName)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		err = NewJiraError(resp, err)
	}

	return resp, err
}

// RemoveWatcher wraps RemoveWatcherWithContext using the background context.
func (s *IssueService) RemoveWatcher(issueID string, userName string) (*Response, error) {
	return s.RemoveWatcherWithContext(context.Background(), issueID, userName)
}

// UpdateAssigneeWithContext updates the user assigned to work on the given issue
//
// Jira API docs: https://docs.atlassian.com/software/jira/docs/api/REST/7.10.2/#api/2/issue-assign
func (s *IssueService) UpdateAssigneeWithContext(ctx context.Context, issueID string, assignee *User) (*Response, error) {
	apiEndPoint := fmt.Sprintf("rest/api/2/issue/%s/assignee", issueID)

	req, err := s.client.NewRequestWithContext(ctx, "PUT", apiEndPoint, assignee)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		err = NewJiraError(resp, err)
	}

	return resp, err
}

// UpdateAssignee wraps UpdateAssigneeWithContext using the background context.
func (s *IssueService) UpdateAssignee(issueID string, assignee *User) (*Response, error) {
	return s.UpdateAssigneeWithContext(context.Background(), issueID, assignee)
}

func (c ChangelogHistory) CreatedTime() (time.Time, error) {
	var t time.Time
	// Ignore null
	if string(c.Created) == "null" {
		return t, nil
	}
	t, err := time.Parse("2006-01-02T15:04:05.999-0700", c.Created)
	return t, err
}

// GetRemoteLinksWithContext gets remote issue links on the issue.
//
// Jira API docs: https://docs.atlassian.com/jira/REST/latest/#api/2/issue-getRemoteIssueLinks
func (s *IssueService) GetRemoteLinksWithContext(ctx context.Context, id string) (*[]RemoteLink, *Response, error) {
	apiEndpoint := fmt.Sprintf("rest/api/2/issue/%s/remotelink", id)
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new([]RemoteLink)
	resp, err := s.client.Do(req, result)
	if err != nil {
		err = NewJiraError(resp, err)
	}
	return result, resp, err
}

// GetRemoteLinks wraps GetRemoteLinksWithContext using the background context.
func (s *IssueService) GetRemoteLinks(id string) (*[]RemoteLink, *Response, error) {
	return s.GetRemoteLinksWithContext(context.Background(), id)
}

// AddRemoteLinkWithContext adds a remote link to issueID.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/#api-rest-api-2-issue-issueIdOrKey-remotelink-post
func (s *IssueService) AddRemoteLinkWithContext(ctx context.Context, issueID string, remotelink *RemoteLink) (*RemoteLink, *Response, error) {
	apiEndpoint := fmt.Sprintf("rest/api/2/issue/%s/remotelink", issueID)
	req, err := s.client.NewRequestWithContext(ctx, "POST", apiEndpoint, remotelink)
	if err != nil {
		return nil, nil, err
	}

	responseRemotelink := new(RemoteLink)
	resp, err := s.client.Do(req, responseRemotelink)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}

	return responseRemotelink, resp, nil
}

// AddRemoteLink wraps AddRemoteLinkWithContext using the background context.
func (s *IssueService) AddRemoteLink(issueID string, remotelink *RemoteLink) (*RemoteLink, *Response, error) {
	return s.AddRemoteLinkWithContext(context.Background(), issueID, remotelink)
}
