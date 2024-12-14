package jira

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestStatusCategoryService_GetList(t *testing.T) {
	setup()
	defer teardown()
	testAPIEdpoint := "/rest/api/2/statuscategory"

	raw, err := ioutil.ReadFile("./mocks/all_statuscategories.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEdpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, testAPIEdpoint)
		fmt.Fprint(w, string(raw))
	})

	statusCategory, _, err := testClient.StatusCategory.GetList()
	if statusCategory == nil {
		t.Error("Expected statusCategory list. StatusCategory list is nil")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}
