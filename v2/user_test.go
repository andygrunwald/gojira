package jira

import (
	"fmt"
	"net/http"
	"testing"
)

func TestUserService_Get_Success(t *testing.T) {
	setup()
	defer teardown()
	testMux.HandleFunc("/rest/api/2/user", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, "/rest/api/2/user?accountId=000000000000000000000000")

		fmt.Fprint(w, `{"self":"http://www.example.com/jira/rest/api/2/user?username=fred","key":"fred",
        "name":"fred","emailAddress":"fred@example.com","avatarUrls":{"48x48":"http://www.example.com/jira/secure/useravatar?size=large&ownerId=fred",
        "24x24":"http://www.example.com/jira/secure/useravatar?size=small&ownerId=fred","16x16":"http://www.example.com/jira/secure/useravatar?size=xsmall&ownerId=fred",
        "32x32":"http://www.example.com/jira/secure/useravatar?size=medium&ownerId=fred"},"displayName":"Fred F. User","active":true,"timeZone":"Australia/Sydney","groups":{"size":3,"items":[
        {"name":"jira-user","self":"http://www.example.com/jira/rest/api/2/group?groupname=jira-user"},{"name":"jira-admin",
        "self":"http://www.example.com/jira/rest/api/2/group?groupname=jira-admin"},{"name":"important","self":"http://www.example.com/jira/rest/api/2/group?groupname=important"
        }]},"applicationRoles":{"size":1,"items":[]},"expand":"groups,applicationRoles"}`)
	})

	if user, _, err := testClient.User.Get("000000000000000000000000"); err != nil {
		t.Errorf("Error given: %s", err)
	} else if user == nil {
		t.Error("Expected user. User is nil")
	}
}

func TestUserService_GetByAccountID_Success(t *testing.T) {
	setup()
	defer teardown()
	testMux.HandleFunc("/rest/api/2/user", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, "/rest/api/2/user?accountId=000000000000000000000000")

		fmt.Fprint(w, `{"self":"http://www.example.com/jira/rest/api/2/user?accountId=000000000000000000000000","accountId": "000000000000000000000000",
        "name":"fred","emailAddress":"fred@example.com","avatarUrls":{"48x48":"http://www.example.com/jira/secure/useravatar?size=large&ownerId=fred",
        "24x24":"http://www.example.com/jira/secure/useravatar?size=small&ownerId=fred","16x16":"http://www.example.com/jira/secure/useravatar?size=xsmall&ownerId=fred",
        "32x32":"http://www.example.com/jira/secure/useravatar?size=medium&ownerId=fred"},"displayName":"Fred F. User","active":true,"timeZone":"Australia/Sydney","groups":{"size":3,"items":[
        {"name":"jira-user","self":"http://www.example.com/jira/rest/api/2/group?groupname=jira-user"},{"name":"jira-admin",
        "self":"http://www.example.com/jira/rest/api/2/group?groupname=jira-admin"},{"name":"important","self":"http://www.example.com/jira/rest/api/2/group?groupname=important"
        }]},"applicationRoles":{"size":1,"items":[]},"expand":"groups,applicationRoles"}`)
	})

	if user, _, err := testClient.User.GetByAccountID("000000000000000000000000"); err != nil {
		t.Errorf("Error given: %s", err)
	} else if user == nil {
		t.Error("Expected user. User is nil")
	}
}

func TestUserService_Create(t *testing.T) {
	setup()
	defer teardown()
	testMux.HandleFunc("/rest/api/2/user", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testRequestURL(t, r, "/rest/api/2/user")

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"name":"charlie","password":"abracadabra","emailAddress":"charlie@atlassian.com",
        "displayName":"Charlie of Atlassian","applicationKeys":["jira-core"]}`)
	})

	u := &User{
		Name:            "charlie",
		Password:        "abracadabra",
		EmailAddress:    "charlie@atlassian.com",
		DisplayName:     "Charlie of Atlassian",
		ApplicationKeys: []string{"jira-core"},
	}

	if user, _, err := testClient.User.Create(u); err != nil {
		t.Errorf("Error given: %s", err)
	} else if user == nil {
		t.Error("Expected user. User is nil")
	}
}

func TestUserService_Delete(t *testing.T) {
	setup()
	defer teardown()
	testMux.HandleFunc("/rest/api/2/user", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testRequestURL(t, r, "/rest/api/2/user?accountId=000000000000000000000000")

		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := testClient.User.Delete("000000000000000000000000")
	if err != nil {
		t.Errorf("Error given: %s", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Wrong status code: %d. Expected %d", resp.StatusCode, http.StatusNoContent)
	}
}

func TestUserService_GetGroups(t *testing.T) {
	setup()
	defer teardown()
	testMux.HandleFunc("/rest/api/2/user/groups", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, "/rest/api/2/user/groups?accountId=000000000000000000000000")

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `[{"name":"jira-software-users","self":"http://www.example.com/jira/rest/api/2/user?accountId=000000000000000000000000"}]`)
	})

	if groups, _, err := testClient.User.GetGroups("000000000000000000000000"); err != nil {
		t.Errorf("Error given: %s", err)
	} else if groups == nil {
		t.Error("Expected user groups. []UserGroup is nil")
	}
}

func TestUserService_GetSelf(t *testing.T) {
	setup()
	defer teardown()
	testMux.HandleFunc("/rest/api/2/myself", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, "/rest/api/2/myself")

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"self":"http://www.example.com/jira/rest/api/2/user?accountId=000000000000000000000000","key":"fred",
        "name":"fred","emailAddress":"fred@example.com","avatarUrls":{"48x48":"http://www.example.com/jira/secure/useravatar?size=large&ownerId=fred",
        "24x24":"http://www.example.com/jira/secure/useravatar?size=small&ownerId=fred","16x16":"http://www.example.com/jira/secure/useravatar?size=xsmall&ownerId=fred",
        "32x32":"http://www.example.com/jira/secure/useravatar?size=medium&ownerId=fred"},"displayName":"Fred F. User","active":true,"timeZone":"Australia/Sydney","groups":{"size":3,"items":[
        {"name":"jira-user","self":"http://www.example.com/jira/rest/api/2/group?groupname=jira-user"},{"name":"jira-admin",
        "self":"http://www.example.com/jira/rest/api/2/group?groupname=jira-admin"},{"name":"important","self":"http://www.example.com/jira/rest/api/2/group?groupname=important"
        }]},"applicationRoles":{"size":1,"items":[]},"expand":"groups,applicationRoles"}`)
	})

	if user, _, err := testClient.User.GetSelf(); err != nil {
		t.Errorf("Error given: %s", err)
	} else if user == nil {
		t.Error("Expected user groups. []UserGroup is nil")
	} else if user.Name != "fred" ||
		!user.Active ||
		user.DisplayName != "Fred F. User" {
		t.Errorf("User JSON deserialized incorrectly")
	}
}

func TestUserService_Find_Success(t *testing.T) {
	setup()
	defer teardown()
	testMux.HandleFunc("/rest/api/2/user/search", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, "/rest/api/2/user/search?query=fred@example.com")

		fmt.Fprint(w, `[{"self":"http://www.example.com/jira/rest/api/2/user?accountId=000000000000000000000000","key":"fred",
        "name":"fred","emailAddress":"fred@example.com","avatarUrls":{"48x48":"http://www.example.com/jira/secure/useravatar?size=large&ownerId=fred",
        "24x24":"http://www.example.com/jira/secure/useravatar?size=small&ownerId=fred","16x16":"http://www.example.com/jira/secure/useravatar?size=xsmall&ownerId=fred",
        "32x32":"http://www.example.com/jira/secure/useravatar?size=medium&ownerId=fred"},"displayName":"Fred F. User","active":true,"timeZone":"Australia/Sydney","groups":{"size":3,"items":[
        {"name":"jira-user","self":"http://www.example.com/jira/rest/api/2/group?groupname=jira-user"},{"name":"jira-admin",
        "self":"http://www.example.com/jira/rest/api/2/group?groupname=jira-admin"},{"name":"important","self":"http://www.example.com/jira/rest/api/2/group?groupname=important"
        }]},"applicationRoles":{"size":1,"items":[]},"expand":"groups,applicationRoles"}]`)
	})

	if user, _, err := testClient.User.Find("fred@example.com"); err != nil {
		t.Errorf("Error given: %s", err)
	} else if user == nil {
		t.Error("Expected user. User is nil")
	}
}

func TestUserService_Find_SuccessParams(t *testing.T) {
	setup()
	defer teardown()
	testMux.HandleFunc("/rest/api/2/user/search", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, "/rest/api/2/user/search?query=fred@example.com&startAt=100&maxResults=1000")

		fmt.Fprint(w, `[{"self":"http://www.example.com/jira/rest/api/2/user?query=fred","key":"fred",
        "name":"fred","emailAddress":"fred@example.com","avatarUrls":{"48x48":"http://www.example.com/jira/secure/useravatar?size=large&ownerId=fred",
        "24x24":"http://www.example.com/jira/secure/useravatar?size=small&ownerId=fred","16x16":"http://www.example.com/jira/secure/useravatar?size=xsmall&ownerId=fred",
        "32x32":"http://www.example.com/jira/secure/useravatar?size=medium&ownerId=fred"},"displayName":"Fred F. User","active":true,"timeZone":"Australia/Sydney","groups":{"size":3,"items":[
        {"name":"jira-user","self":"http://www.example.com/jira/rest/api/2/group?groupname=jira-user"},{"name":"jira-admin",
        "self":"http://www.example.com/jira/rest/api/2/group?groupname=jira-admin"},{"name":"important","self":"http://www.example.com/jira/rest/api/2/group?groupname=important"
        }]},"applicationRoles":{"size":1,"items":[]},"expand":"groups,applicationRoles"}]`)
	})

	if user, _, err := testClient.User.Find("fred@example.com", WithStartAt(100), WithMaxResults(1000)); err != nil {
		t.Errorf("Error given: %s", err)
	} else if user == nil {
		t.Error("Expected user. User is nil")
	}
}
