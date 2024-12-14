package jira

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestPermissionSchemeService_GetList(t *testing.T) {
	setup()
	defer teardown()
	testAPIEndpoint := "/rest/api/3/permissionscheme"

	raw, err := ioutil.ReadFile("./mocks/all_permissionschemes.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEndpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, testAPIEndpoint)
		fmt.Fprint(w, string(raw))
	})

	permissionScheme, _, err := testClient.PermissionScheme.GetList()
	if err != nil {
		t.Errorf("Error given: %v", err)
	}
	if permissionScheme == nil {
		t.Error("Expected permissionScheme list. PermissionScheme list is nil")
		return
	}
	if len(permissionScheme.PermissionSchemes) != 2 {
		t.Errorf("Expected %d permissionSchemes but got %d", 2, len(permissionScheme.PermissionSchemes))
	}
}

func TestPermissionSchemeService_GetList_NoList(t *testing.T) {
	setup()
	defer teardown()
	testAPIEndpoint := "/rest/api/3/permissionscheme"

	raw, err := ioutil.ReadFile("./mocks/no_permissionschemes.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEndpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, testAPIEndpoint)
		fmt.Fprint(w, string(raw))
	})

	permissionScheme, _, err := testClient.PermissionScheme.GetList()
	if permissionScheme != nil {
		t.Errorf("Expected permissionScheme list has %d entries but should be nil", len(permissionScheme.PermissionSchemes))
	}
	if err == nil {
		t.Errorf("No error given")
	}
}

func TestPermissionSchemeService_Get(t *testing.T) {
	setup()
	defer teardown()
	testAPIEdpoint := "/rest/api/3/permissionscheme/10100"
	raw, err := ioutil.ReadFile("./mocks/permissionscheme.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEdpoint, func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, "GET")
		testRequestURL(t, request, testAPIEdpoint)
		fmt.Fprint(writer, string(raw))
	})

	permissionScheme, _, err := testClient.PermissionScheme.Get(10100)
	if permissionScheme == nil {
		t.Errorf("Expected permissionscheme, got nil")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}

func TestPermissionSchemeService_Get_NoScheme(t *testing.T) {
	setup()
	defer teardown()
	testAPIEdpoint := "/rest/api/3/permissionscheme/99999"
	raw, err := ioutil.ReadFile("./mocks/no_permissionscheme.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEdpoint, func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, "GET")
		testRequestURL(t, request, testAPIEdpoint)
		fmt.Fprint(writer, string(raw))
	})

	permissionScheme, _, err := testClient.PermissionScheme.Get(99999)
	if permissionScheme != nil {
		t.Errorf("Expected nil, got permissionschme %v", permissionScheme)
	}
	if err == nil {
		t.Errorf("No error given")
	}
}
