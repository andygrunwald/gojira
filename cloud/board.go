package cloud

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// BoardService handles Agile Boards for the Jira instance / API.
//
// Jira API docs: https://docs.atlassian.com/jira-software/REST/server/
type BoardService service

// BoardsList reflects a list of agile boards
type BoardsList struct {
	MaxResults int     `json:"maxResults" structs:"maxResults"`
	StartAt    int     `json:"startAt" structs:"startAt"`
	Total      int     `json:"total" structs:"total"`
	IsLast     bool    `json:"isLast" structs:"isLast"`
	Values     []Board `json:"values" structs:"values"`
}

// Board represents a Jira agile board
type Board struct {
	ID       int           `json:"id,omitempty" structs:"id,omitempty"`
	Self     string        `json:"self,omitempty" structs:"self,omitempty"`
	Name     string        `json:"name,omitempty" structs:"name,omitemtpy"`
	Type     string        `json:"type,omitempty" structs:"type,omitempty"`
	Location BoardLocation `json:"location,omitempty" structs:"location,omitempty"`
	FilterID int           `json:"filterId,omitempty" structs:"filterId,omitempty"`
}

// BoardLocation represents the location of a Jira board
type BoardLocation struct {
	ProjectID      int    `json:"projectId"`
	UserID         int    `json:"userId"`
	UserAccountID  string `json:"userAccountId"`
	DisplayName    string `json:"displayName"`
	ProjectName    string `json:"projectName"`
	ProjectKey     string `json:"projectKey"`
	ProjectTypeKey string `json:"projectTypeKey"`
	Name           string `json:"name"`
}

// BoardListOptions specifies the optional parameters to the BoardService.GetList
type BoardListOptions struct {
	// BoardType filters results to boards of the specified type.
	// Valid values: scrum, kanban.
	BoardType string `url:"type,omitempty"`
	// Name filters results to boards that match or partially match the specified name.
	Name string `url:"name,omitempty"`
	// ProjectKeyOrID filters results to boards that are relevant to a project.
	// Relevance meaning that the JQL filter defined in board contains a reference to a project.
	ProjectKeyOrID string `url:"projectKeyOrId,omitempty"`

	SearchOptions
}

// GetAllSprintsOptions specifies the optional parameters to the BoardService.GetList
type GetAllSprintsOptions struct {
	// State filters results to sprints in the specified states, comma-separate list
	State string `url:"state,omitempty"`

	SearchOptions
}

// SprintsList reflects a list of agile sprints
type SprintsList struct {
	MaxResults int      `json:"maxResults" structs:"maxResults"`
	StartAt    int      `json:"startAt" structs:"startAt"`
	Total      int      `json:"total" structs:"total"`
	IsLast     bool     `json:"isLast" structs:"isLast"`
	Values     []Sprint `json:"values" structs:"values"`
}

// Sprint represents a sprint on Jira agile board
type Sprint struct {
	ID            int        `json:"id" structs:"id"`
	Name          string     `json:"name" structs:"name"`
	CompleteDate  *time.Time `json:"completeDate" structs:"completeDate"`
	EndDate       *time.Time `json:"endDate" structs:"endDate"`
	StartDate     *time.Time `json:"startDate" structs:"startDate"`
	OriginBoardID int        `json:"originBoardId" structs:"originBoardId"`
	Self          string     `json:"self" structs:"self"`
	State         string     `json:"state" structs:"state"`
	Goal          string     `json:"goal,omitempty" structs:"goal"`
}

// BoardConfiguration represents a boardConfiguration of a jira board
type BoardConfiguration struct {
	ID           int                            `json:"id"`
	Name         string                         `json:"name"`
	Self         string                         `json:"self"`
	Location     BoardConfigurationLocation     `json:"location"`
	Filter       BoardConfigurationFilter       `json:"filter"`
	SubQuery     BoardConfigurationSubQuery     `json:"subQuery"`
	ColumnConfig BoardConfigurationColumnConfig `json:"columnConfig"`
}

// BoardConfigurationFilter reference to the filter used by the given board.
type BoardConfigurationFilter struct {
	ID   string `json:"id"`
	Self string `json:"self"`
}

// BoardConfigurationSubQuery  (Kanban only) - JQL subquery used by the given board.
type BoardConfigurationSubQuery struct {
	Query string `json:"query"`
}

// BoardConfigurationLocation reference to the container that the board is located in
type BoardConfigurationLocation struct {
	Type string `json:"type"`
	Key  string `json:"key"`
	ID   string `json:"id"`
	Self string `json:"self"`
	Name string `json:"name"`
}

// BoardConfigurationColumnConfig lists the columns for a given board in the order defined in the column configuration
// with constrainttype (none, issueCount, issueCountExclSubs)
type BoardConfigurationColumnConfig struct {
	Columns        []BoardConfigurationColumn `json:"columns"`
	ConstraintType string                     `json:"constraintType"`
}

// BoardConfigurationColumn lists the name of the board with the statuses that maps to a particular column
type BoardConfigurationColumn struct {
	Name   string                           `json:"name"`
	Status []BoardConfigurationColumnStatus `json:"statuses"`
	Min    int                              `json:"min,omitempty"`
	Max    int                              `json:"max,omitempty"`
}

// BoardConfigurationColumnStatus represents a status in the column configuration
type BoardConfigurationColumnStatus struct {
	ID   string `json:"id"`
	Self string `json:"self"`
}

// GetBoardIssuesOptions specifies the optional parameters to the BoardService.GetIssuesForBacklog
type GetBoardIssuesOptions struct {
	// Jql filters results to sprints in the specified jql query
	Jql string `json:"jql"`

	SearchOptions
}

// IssuesInBoardResult represents a wrapper struct for search result
type IssuesInBoardResult struct {
	Issues []Issue `json:"issues"`
}

// GetAllBoards will returns all boards. This only includes boards that the user has permission to view.
//
// Jira API docs: https://docs.atlassian.com/jira-software/REST/cloud/#agile/1.0/board-getAllBoards
//
// TODO Double check this method if this works as expected, is using the latest API and the response is complete
// This double check effort is done for v2 - Remove this two lines if this is completed.
func (s *BoardService) GetAllBoards(ctx context.Context, opt *BoardListOptions) (*BoardsList, *Response, error) {
	apiEndpoint := "rest/agile/1.0/board"
	url, err := addOptions(apiEndpoint, opt)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	boards := new(BoardsList)
	resp, err := s.client.Do(req, boards)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}

	return boards, resp, err
}

// GetBoard returns the board for the given board ID.
// This board will only be returned if the user has permission to view it.
// Admins without the view permission will see the board as a private one, so will see only a subset of the board's data (board location for instance).
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/software/rest/api-group-board/#api-rest-agile-1-0-board-boardid-get
func (s *BoardService) GetBoard(ctx context.Context, boardID int64) (*Board, *Response, error) {
	apiEndpoint := fmt.Sprintf("rest/agile/1.0/board/%v", boardID)
	req, err := s.client.NewRequest(ctx, http.MethodGet, apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	board := new(Board)
	resp, err := s.client.Do(req, board)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}

	return board, resp, nil
}

// CreateBoard creates a new board. Board name, type and filter Id is required.
// name - Must be less than 255 characters.
// type - Valid values: scrum, kanban
// filterId - Id of a filter that the user has permissions to view.
// Note, if the user does not have the 'Create shared objects' permission and tries to create a shared board, a private
// board will be created instead (remember that board sharing depends on the filter sharing).
//
// Jira API docs: https://docs.atlassian.com/jira-software/REST/cloud/#agile/1.0/board-createBoard
//
// TODO Double check this method if this works as expected, is using the latest API and the response is complete
// This double check effort is done for v2 - Remove this two lines if this is completed.
func (s *BoardService) CreateBoard(ctx context.Context, board *Board) (*Board, *Response, error) {
	apiEndpoint := "rest/agile/1.0/board"
	req, err := s.client.NewRequest(ctx, http.MethodPost, apiEndpoint, board)
	if err != nil {
		return nil, nil, err
	}

	responseBoard := new(Board)
	resp, err := s.client.Do(req, responseBoard)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}

	return responseBoard, resp, nil
}

// DeleteBoard will delete an agile board.
//
// Jira API docs: https://docs.atlassian.com/jira-software/REST/cloud/#agile/1.0/board-deleteBoard
// Caller must close resp.Body
//
// TODO Double check this method if this works as expected, is using the latest API and the response is complete
// This double check effort is done for v2 - Remove this two lines if this is completed.
func (s *BoardService) DeleteBoard(ctx context.Context, boardID int) (*Board, *Response, error) {
	apiEndpoint := fmt.Sprintf("rest/agile/1.0/board/%v", boardID)
	req, err := s.client.NewRequest(ctx, http.MethodDelete, apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		err = NewJiraError(resp, err)
	}
	return nil, resp, err
}

// GetAllSprints returns all sprints from a board, for a given board ID.
// This only includes sprints that the user has permission to view.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/software/rest/api-group-board/#api-rest-agile-1-0-board-boardid-sprint-get
func (s *BoardService) GetAllSprints(ctx context.Context, boardID int64, options *GetAllSprintsOptions) (*SprintsList, *Response, error) {
	apiEndpoint := fmt.Sprintf("rest/agile/1.0/board/%d/sprint", boardID)
	url, err := addOptions(apiEndpoint, options)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(SprintsList)
	resp, err := s.client.Do(req, result)
	if err != nil {
		err = NewJiraError(resp, err)
	}

	return result, resp, err
}

// GetIssuesForBacklog returns all issues from a backlog, for a given board ID.
// This only includes issues that the user has permission to view.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/software/rest/api-group-board/#api-rest-agile-1-0-board-boardid-backlog-get
func (s *BoardService) GetIssuesForBacklog(ctx context.Context, boardID int64, options *GetBoardIssuesOptions) ([]Issue, *Response, error) {
	apiEndpoint := fmt.Sprintf("rest/agile/1.0/board/%d/backlog", boardID)
	url, err := addOptions(apiEndpoint, options)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(IssuesInBoardResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		err = NewJiraError(resp, err)
	}

	return result.Issues, resp, err
}

// GetBoardConfiguration will return a board configuration for a given board Id
// Jira API docs:https://developer.atlassian.com/cloud/jira/software/rest/#api-rest-agile-1-0-board-boardId-configuration-get
//
// TODO Double check this method if this works as expected, is using the latest API and the response is complete
// This double check effort is done for v2 - Remove this two lines if this is completed.
func (s *BoardService) GetBoardConfiguration(ctx context.Context, boardID int) (*BoardConfiguration, *Response, error) {
	apiEndpoint := fmt.Sprintf("rest/agile/1.0/board/%d/configuration", boardID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, apiEndpoint, nil)

	if err != nil {
		return nil, nil, err
	}

	result := new(BoardConfiguration)
	resp, err := s.client.Do(req, result)
	if err != nil {
		err = NewJiraError(resp, err)
	}

	return result, resp, err

}
