package cloud

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestComponentService_Create_Success(t *testing.T) {
	setup()
	defer teardown()
	testAPIEndpoint := "/rest/api/3/component"

	testMux.HandleFunc(testAPIEndpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testRequestURL(t, r, testAPIEndpoint)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{ "self": "http://www.example.com/jira/rest/api/2/component/10000", "id": "10000", "name": "Component 1", "description": "This is a Jira component", "lead": { "self": "http://www.example.com/jira/rest/api/2/user?username=fred", "name": "fred", "avatarUrls": { "48x48": "http://www.example.com/jira/secure/useravatar?size=large&ownerId=fred", "24x24": "http://www.example.com/jira/secure/useravatar?size=small&ownerId=fred", "16x16": "http://www.example.com/jira/secure/useravatar?size=xsmall&ownerId=fred", "32x32": "http://www.example.com/jira/secure/useravatar?size=medium&ownerId=fred" }, "displayName": "Fred F. User", "active": false }, "assigneeType": "PROJECT_LEAD", "assignee": { "self": "http://www.example.com/jira/rest/api/2/user?username=fred", "name": "fred", "avatarUrls": { "48x48": "http://www.example.com/jira/secure/useravatar?size=large&ownerId=fred", "24x24": "http://www.example.com/jira/secure/useravatar?size=small&ownerId=fred", "16x16": "http://www.example.com/jira/secure/useravatar?size=xsmall&ownerId=fred", "32x32": "http://www.example.com/jira/secure/useravatar?size=medium&ownerId=fred" }, "displayName": "Fred F. User", "active": false }, "realAssigneeType": "PROJECT_LEAD", "realAssignee": { "self": "http://www.example.com/jira/rest/api/2/user?username=fred", "name": "fred", "avatarUrls": { "48x48": "http://www.example.com/jira/secure/useravatar?size=large&ownerId=fred", "24x24": "http://www.example.com/jira/secure/useravatar?size=small&ownerId=fred", "16x16": "http://www.example.com/jira/secure/useravatar?size=xsmall&ownerId=fred", "32x32": "http://www.example.com/jira/secure/useravatar?size=medium&ownerId=fred" }, "displayName": "Fred F. User", "active": false }, "isAssigneeTypeValid": false, "project": "HSP", "projectId": 10000 }`)
	})

	component, _, err := testClient.Component.Create(context.Background(), &ComponentCreateOptions{
		Name: "foo-bar",
	})
	if component == nil {
		t.Error("Expected component. Component is nil")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}

func TestComponentService_Update_Success(t *testing.T) {
	setup()
	defer teardown()
	testAPIEndpoint := "/rest/api/3/component/42102"

	// Mock simple component with a concrete name
	component := ProjectComponent{}
	component.ID = "42102"
	component.Name = "Brand new component"

	raw_component_put_response, err := os.ReadFile("../testing/mock-data/component_updated.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEndpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testRequestURL(t, r, testAPIEndpoint)
		fmt.Fprint(w, string(raw_component_put_response))
	})

	componentUpdated, _, err := testClient.Component.Update(context.Background(), &component)

	// Now we check, that the received component contains values directly from the server response
	// updatedNameFromServerResponse = "Some Component UPDATED" (see the ../testing/mock-data/component_updated.json file)
	if componentUpdated == nil {
		t.Error("Expected component. Component is nil")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}

	// Check that "UPDATED" is inside the name, see the ../testing/mock-data/component_updated.json file
	if !strings.Contains(componentUpdated.Name, "UPDATED") {
		t.Errorf("Expected 'UPDATED' in name. Got %s", componentUpdated.Name)
	}
}

func TestComponentService_Get(t *testing.T) {
	setup()
	defer teardown()
	testAPIEndpoint := "/rest/api/3/component/42102"

	raw, err := os.ReadFile("../testing/mock-data/component_get.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEndpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, testAPIEndpoint)
		fmt.Fprint(w, string(raw))
	})

	component, _, err := testClient.Component.Get(context.Background(), "42102")
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
	if component == nil {
		t.Error("Expected component. Component is nil")
	}
}

func TestComponentService_Get_NoComponent(t *testing.T) {
	setup()
	defer teardown()
	testAPIEndpoint := "/rest/api/3/component/99999999"

	testMux.HandleFunc(testAPIEndpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, testAPIEndpoint)
		fmt.Fprint(w, nil)
	})

	component, resp, err := testClient.Component.Get(context.Background(), "99999999")

	if component != nil {
		t.Errorf("Expected nil. Got %+v", component)
	}
	if resp.Status == "404" {
		t.Errorf("Expected status 404. Got %s", resp.Status)
	}
	if err == nil {
		t.Error("No error given. Expected one")
	}
}
