package jira

import (
	"fmt"
	"net/http"
	"testing"
)

func TestComponentService_Create_Success(t *testing.T) {
	setup()
	defer teardown()
	testMux.HandleFunc("/rest/api/2/component", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testRequestURL(t, r, "/rest/api/2/component")

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{ "self": "http://www.example.com/jira/rest/api/2/component/10000", "id": "10000", "name": "Component 1", "description": "This is a Jira component", "lead": { "self": "http://www.example.com/jira/rest/api/2/user?username=fred", "name": "fred", "avatarUrls": { "48x48": "http://www.example.com/jira/secure/useravatar?size=large&ownerId=fred", "24x24": "http://www.example.com/jira/secure/useravatar?size=small&ownerId=fred", "16x16": "http://www.example.com/jira/secure/useravatar?size=xsmall&ownerId=fred", "32x32": "http://www.example.com/jira/secure/useravatar?size=medium&ownerId=fred" }, "displayName": "Fred F. User", "active": false }, "assigneeType": "PROJECT_LEAD", "assignee": { "self": "http://www.example.com/jira/rest/api/2/user?username=fred", "name": "fred", "avatarUrls": { "48x48": "http://www.example.com/jira/secure/useravatar?size=large&ownerId=fred", "24x24": "http://www.example.com/jira/secure/useravatar?size=small&ownerId=fred", "16x16": "http://www.example.com/jira/secure/useravatar?size=xsmall&ownerId=fred", "32x32": "http://www.example.com/jira/secure/useravatar?size=medium&ownerId=fred" }, "displayName": "Fred F. User", "active": false }, "realAssigneeType": "PROJECT_LEAD", "realAssignee": { "self": "http://www.example.com/jira/rest/api/2/user?username=fred", "name": "fred", "avatarUrls": { "48x48": "http://www.example.com/jira/secure/useravatar?size=large&ownerId=fred", "24x24": "http://www.example.com/jira/secure/useravatar?size=small&ownerId=fred", "16x16": "http://www.example.com/jira/secure/useravatar?size=xsmall&ownerId=fred", "32x32": "http://www.example.com/jira/secure/useravatar?size=medium&ownerId=fred" }, "displayName": "Fred F. User", "active": false }, "isAssigneeTypeValid": false, "project": "HSP", "projectId": 10000 }`)
	})

	component, _, err := testClient.Component.Create(&CreateComponentOptions{
		Name: "foo-bar",
	})
	if component == nil {
		t.Error("Expected component. Component is nil")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}
