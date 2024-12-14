package jira

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestPriorityService_GetList(t *testing.T) {
	setup()
	defer teardown()
	testAPIEdpoint := "/rest/api/2/priority"

	raw, err := ioutil.ReadFile("./mocks/all_priorities.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEdpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, testAPIEdpoint)
		fmt.Fprint(w, string(raw))
	})

	priorities, _, err := testClient.Priority.GetList()
	if priorities == nil {
		t.Error("Expected priority list. Priority list is nil")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}
