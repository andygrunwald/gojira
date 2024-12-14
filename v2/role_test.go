package jira

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestRoleService_GetList_NoList(t *testing.T) {
	setup()
	defer teardown()
	testAPIEndpoint := "/rest/api/3/role"

	raw, err := ioutil.ReadFile("./mocks/no_roles.json")
	if err != nil {
		t.Error(err.Error())
	}

	testMux.HandleFunc(testAPIEndpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, testAPIEndpoint)
		fmt.Fprint(w, string(raw))
	})

	roles, _, err := testClient.Role.GetList()
	if roles != nil {
		t.Errorf("Expected role list has %d entries but should be nil", len(*roles))
	}
	if err == nil {
		t.Errorf("No error given")
	}
}

func TestRoleService_GetList(t *testing.T) {
	setup()
	defer teardown()
	testAPIEndpoint := "/rest/api/3/role"

	raw, err := ioutil.ReadFile("./mocks/all_roles.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEndpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, testAPIEndpoint)
		fmt.Fprint(w, string(raw))
	})

	roles, _, err := testClient.Role.GetList()
	if err != nil {
		t.Errorf("Error given: %v", err)
	}
	if roles == nil {
		t.Error("Expected role list. Role list is nil")
		return
	}
	if len(*roles) != 2 {
		t.Errorf("Expected %d roles but got %d", 2, len(*roles))
	}
}

func TestRoleService_Get_NoRole(t *testing.T) {
	setup()
	defer teardown()
	testAPIEdpoint := "/rest/api/3/role/99999"
	raw, err := ioutil.ReadFile("./mocks/no_role.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEdpoint, func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, "GET")
		testRequestURL(t, request, testAPIEdpoint)
		fmt.Fprint(writer, string(raw))
	})

	role, _, err := testClient.Role.Get(99999)
	if role != nil {
		t.Errorf("Expected nil, got role %v", role)
	}
	if err == nil {
		t.Errorf("No error given")
	}
}

func TestRoleService_Get(t *testing.T) {
	setup()
	defer teardown()
	testAPIEdpoint := "/rest/api/3/role/10002"
	raw, err := ioutil.ReadFile("./mocks/role.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEdpoint, func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, "GET")
		testRequestURL(t, request, testAPIEdpoint)
		fmt.Fprint(writer, string(raw))
	})

	role, _, err := testClient.Role.Get(10002)
	if role == nil {
		t.Errorf("Expected Role, got nil")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}
