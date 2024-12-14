package jira

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestResolutionService_GetList(t *testing.T) {
	setup()
	defer teardown()
	testAPIEdpoint := "/rest/api/2/resolution"

	raw, err := ioutil.ReadFile("./mocks/all_resolutions.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEdpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, testAPIEdpoint)
		fmt.Fprint(w, string(raw))
	})

	resolution, _, err := testClient.Resolution.GetList()
	if resolution == nil {
		t.Error("Expected resolution list. Resolution list is nil")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}
