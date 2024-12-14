package jira

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestFieldService_GetList(t *testing.T) {
	setup()
	defer teardown()
	testAPIEdpoint := "/rest/api/2/field"

	raw, err := ioutil.ReadFile("./mocks/all_fields.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEdpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, testAPIEdpoint)
		fmt.Fprint(w, string(raw))
	})

	fields, _, err := testClient.Field.GetList()
	if fields == nil {
		t.Error("Expected field list. Field list is nil")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}
