package cloud

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
)

func TestBoardService_GetAllBoards(t *testing.T) {
	setup()
	defer teardown()
	testapiEndpoint := "/rest/agile/1.0/board"

	raw, err := os.ReadFile("../testing/mock-data/all_boards.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testapiEndpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testRequestURL(t, r, testapiEndpoint)
		fmt.Fprint(w, string(raw))
	})

	projects, _, err := testClient.Board.GetAllBoards(context.Background(), nil)
	if projects == nil {
		t.Error("Expected boards list. Boards list is nil")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}

// Test with params
func TestBoardService_GetAllBoards_WithFilter(t *testing.T) {
	setup()
	defer teardown()
	testapiEndpoint := "/rest/agile/1.0/board"

	raw, err := os.ReadFile("../testing/mock-data/all_boards_filtered.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testapiEndpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testRequestURL(t, r, testapiEndpoint)
		testRequestParams(t, r, map[string]string{"type": "scrum", "name": "Test", "startAt": "1", "maxResults": "10", "projectKeyOrId": "TE"})
		fmt.Fprint(w, string(raw))
	})

	boardsListOptions := &BoardListOptions{
		BoardType:      "scrum",
		Name:           "Test",
		ProjectKeyOrID: "TE",
	}
	boardsListOptions.StartAt = 1
	boardsListOptions.MaxResults = 10

	projects, _, err := testClient.Board.GetAllBoards(context.Background(), boardsListOptions)
	if projects == nil {
		t.Error("Expected boards list. Boards list is nil")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}

func TestBoardService_GetBoard(t *testing.T) {
	setup()
	defer teardown()
	testapiEndpoint := "/rest/agile/1.0/board/1"

	testMux.HandleFunc(testapiEndpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testRequestURL(t, r, testapiEndpoint)
		fmt.Fprint(w, `{"id":4,"self":"https://test.jira.org/rest/agile/1.0/board/1","name":"Test Weekly","type":"scrum"}`)
	})

	board, _, err := testClient.Board.GetBoard(context.Background(), 1)
	if board == nil {
		t.Error("Expected board list. Board list is nil")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}

func TestBoardService_GetBoard_WrongID(t *testing.T) {
	setup()
	defer teardown()
	testAPIEndpoint := "/rest/api/2/board/99999999"

	testMux.HandleFunc(testAPIEndpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testRequestURL(t, r, testAPIEndpoint)
		fmt.Fprint(w, nil)
	})

	board, resp, err := testClient.Board.GetBoard(context.Background(), 99999999)
	if board != nil {
		t.Errorf("Expected nil. Got %s", err)
	}

	if resp.Status == "404" {
		t.Errorf("Expected status 404. Got %s", resp.Status)
	}
	if err == nil {
		t.Errorf("Error given: %s", err)
	}
}

func TestBoardService_CreateBoard(t *testing.T) {
	setup()
	defer teardown()
	testMux.HandleFunc("/rest/agile/1.0/board", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testRequestURL(t, r, "/rest/agile/1.0/board")

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"id":17,"self":"https://test.jira.org/rest/agile/1.0/board/17","name":"Test","type":"kanban"}`)
	})

	b := &Board{
		Name:     "Test",
		Type:     "kanban",
		FilterID: 17,
	}
	issue, _, err := testClient.Board.CreateBoard(context.Background(), b)
	if issue == nil {
		t.Error("Expected board. Board is nil")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}

func TestBoardService_DeleteBoard(t *testing.T) {
	setup()
	defer teardown()
	testMux.HandleFunc("/rest/agile/1.0/board/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		testRequestURL(t, r, "/rest/agile/1.0/board/1")

		w.WriteHeader(http.StatusNoContent)
		fmt.Fprint(w, `{}`)
	})

	_, resp, err := testClient.Board.DeleteBoard(context.Background(), 1)
	if resp.StatusCode != 204 {
		t.Error("Expected board not deleted.")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}

func TestBoardService_GetAllSprints(t *testing.T) {
	setup()
	defer teardown()

	testAPIEndpoint := "/rest/agile/1.0/board/123/sprint"

	raw, err := os.ReadFile("../testing/mock-data/sprints_filtered.json")
	if err != nil {
		t.Error(err.Error())
	}

	testMux.HandleFunc(testAPIEndpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testRequestURL(t, r, testAPIEndpoint)
		fmt.Fprint(w, string(raw))
	})

	sprints, _, err := testClient.Board.GetAllSprints(context.Background(), 123, &GetAllSprintsOptions{State: "active,future"})
	if err != nil {
		t.Errorf("Got error: %v", err)
	}

	if sprints == nil {
		t.Error("Expected sprint list. Got nil.")
		return
	}

	if len(sprints.Values) != 1 {
		t.Errorf("Expected 1 transition. Got %d", len(sprints.Values))
	}
}

func TestBoardService_GetBoardConfigoration(t *testing.T) {
	setup()
	defer teardown()
	testAPIEndpoint := "/rest/agile/1.0/board/35/configuration"

	raw, err := os.ReadFile("../testing/mock-data/board_configuration.json")
	if err != nil {
		t.Error(err.Error())
	}

	testMux.HandleFunc(testAPIEndpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testRequestURL(t, r, testAPIEndpoint)
		fmt.Fprint(w, string(raw))
	})

	boardConfiguration, _, err := testClient.Board.GetBoardConfiguration(context.Background(), 35)
	if err != nil {
		t.Errorf("Got error: %v", err)
	}

	if boardConfiguration == nil {
		t.Error("Expected boardConfiguration. Got nil.")
		return
	}

	if len(boardConfiguration.ColumnConfig.Columns) != 6 {
		t.Errorf("Expected 6 columns. go %d", len(boardConfiguration.ColumnConfig.Columns))
	}

	backlogColumn := boardConfiguration.ColumnConfig.Columns[0]
	if backlogColumn.Min != 5 {
		t.Errorf("Expected a min of 5 issues in backlog. Got %d", backlogColumn.Min)
	}
	if backlogColumn.Max != 30 {
		t.Errorf("Expected a max of 30 issues in backlog. Got %d", backlogColumn.Max)
	}

	inProgressColumn := boardConfiguration.ColumnConfig.Columns[2]
	if inProgressColumn.Min != 0 {
		t.Errorf("Expected a min of 0 issues in progress. Got %d", inProgressColumn.Min)
	}
	if inProgressColumn.Max != 0 {
		t.Errorf("Expected a max of 0 issues in progress. Got %d", inProgressColumn.Max)
	}
}

func TestBoardService_GetIssuesForBacklog(t *testing.T) {
	setup()
	defer teardown()
	testapiEndpoint := "/rest/agile/1.0/board/123/backlog"

	raw, err := os.ReadFile("../testing/mock-data/backlog_in_board.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testapiEndpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testRequestURL(t, r, testapiEndpoint)
		fmt.Fprint(w, string(raw))
	})

	issues, _, err := testClient.Board.GetIssuesForBacklog(context.Background(), 123, nil)
	if err != nil {
		t.Errorf("Error given: %v", err)
	}
	if issues == nil {
		t.Error("Expected issues in sprint list. Issues list is nil")
	}
	if len(issues) != 1 {
		t.Errorf("Expect there to be 1 issue in the sprint, found %v", len(issues))
	}

}
