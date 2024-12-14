package cloud

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestUserService_Get_Success(t *testing.T) {
	setup()
	defer teardown()
	testMux.HandleFunc("/rest/api/2/user", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testRequestURL(t, r, "/rest/api/2/user?accountId=000000000000000000000000")

		fmt.Fprint(w, `{"self":"http://www.example.com/jira/rest/api/2/user?username=fred","key":"fred",
        "name":"fred","emailAddress":"fred@example.com","avatarUrls":{"48x48":"http://www.example.com/jira/secure/useravatar?size=large&ownerId=fred",
        "24x24":"http://www.example.com/jira/secure/useravatar?size=small&ownerId=fred","16x16":"http://www.example.com/jira/secure/useravatar?size=xsmall&ownerId=fred",
        "32x32":"http://www.example.com/jira/secure/useravatar?size=medium&ownerId=fred"},"displayName":"Fred F. User","active":true,"timeZone":"Australia/Sydney","groups":{"size":3,"items":[
        {"name":"jira-user","self":"http://www.example.com/jira/rest/api/2/group?groupname=jira-user"},{"name":"jira-admin",
        "self":"http://www.example.com/jira/rest/api/2/group?groupname=jira-admin"},{"name":"important","self":"http://www.example.com/jira/rest/api/2/group?groupname=important"
        }]},"applicationRoles":{"size":1,"items":[]},"expand":"groups,applicationRoles"}`)
	})

	if user, _, err := testClient.User.Get(context.Background(), "000000000000000000000000"); err != nil {
		t.Errorf("Error given: %s", err)
	} else if user == nil {
		t.Error("Expected user. User is nil")
	}
}

func TestUserService_GetByAccountID_Success(t *testing.T) {
	setup()
	defer teardown()
	testMux.HandleFunc("/rest/api/2/user", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testRequestURL(t, r, "/rest/api/2/user?accountId=000000000000000000000000")

		fmt.Fprint(w, `{"self":"http://www.example.com/jira/rest/api/2/user?accountId=000000000000000000000000","accountId": "000000000000000000000000",
        "name":"fred","emailAddress":"fred@example.com","avatarUrls":{"48x48":"http://www.example.com/jira/secure/useravatar?size=large&ownerId=fred",
        "24x24":"http://www.example.com/jira/secure/useravatar?size=small&ownerId=fred","16x16":"http://www.example.com/jira/secure/useravatar?size=xsmall&ownerId=fred",
        "32x32":"http://www.example.com/jira/secure/useravatar?size=medium&ownerId=fred"},"displayName":"Fred F. User","active":true,"timeZone":"Australia/Sydney","groups":{"size":3,"items":[
        {"name":"jira-user","self":"http://www.example.com/jira/rest/api/2/group?groupname=jira-user"},{"name":"jira-admin",
        "self":"http://www.example.com/jira/rest/api/2/group?groupname=jira-admin"},{"name":"important","self":"http://www.example.com/jira/rest/api/2/group?groupname=important"
        }]},"applicationRoles":{"size":1,"items":[]},"expand":"groups,applicationRoles"}`)
	})

	if user, _, err := testClient.User.GetByAccountID(context.Background(), "000000000000000000000000"); err != nil {
		t.Errorf("Error given: %s", err)
	} else if user == nil {
		t.Error("Expected user. User is nil")
	}
}

func TestUserService_Create(t *testing.T) {
	setup()
	defer teardown()
	testMux.HandleFunc("/rest/api/2/user", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testRequestURL(t, r, "/rest/api/2/user")

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `
		{
			"name": "charlie",
			"password": "abracadabra",
			"emailAddress": "charlie@atlassian.com",
			"displayName": "Charlie of Atlassian",
			"applicationRoles": {
				"size": 1,
				"max-results": 1,
				"items": [{
					"key": "jira-software",
					"groups": [
						"jira-software-users",
						"jira-testers"
					],
					"name": "Jira Software",
					"defaultGroups": [
						"jira-software-users"
					],
					"selectedByDefault": false,
					"defined": false,
					"numberOfSeats": 10,
					"remainingSeats": 5,
					"userCount": 5,
					"userCountDescription": "5 developers",
					"hasUnlimitedSeats": false,
					"platform": false
				}]
			},
			"groups": {
				"size": 2,
				"max-results": 2,
				"items": [{
						"name": "jira-core",
						"self": "jira-core"
					},
					{
						"name": "jira-test",
						"self": "jira-test"
					}
				]
			}
		}
		`)
	})

	u := &User{
		Name:         "charlie",
		Password:     "abracadabra",
		EmailAddress: "charlie@atlassian.com",
		DisplayName:  "Charlie of Atlassian",
		Groups: UserGroups{
			Size: 2,
			Items: []UserGroup{
				{
					Name: "jira-core",
					Self: "jira-core",
				},
				{
					Name: "jira-test",
					Self: "jira-test",
				},
			},
		},
		ApplicationRoles: ApplicationRoles{
			Size: 1,
			Items: []ApplicationRole{
				{
					Key:                  "jira-software",
					Groups:               []string{"jira-software-users", "jira-testers"},
					Name:                 "Jira Software",
					DefaultGroups:        []string{"jira-software-users"},
					SelectedByDefault:    false,
					Defined:              false,
					NumberOfSeats:        10,
					RemainingSeats:       5,
					UserCount:            5,
					UserCountDescription: "5 developers",
					HasUnlimitedSeats:    false,
					Platform:             false,
				},
			},
		},
	}

	if user, _, err := testClient.User.Create(context.Background(), u); err != nil {
		t.Errorf("Error given: %s", err)
	} else if user == nil {
		t.Error("Expected user. User is nil")
	}
}

func TestUserService_Delete(t *testing.T) {
	setup()
	defer teardown()
	testMux.HandleFunc("/rest/api/2/user", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		testRequestURL(t, r, "/rest/api/2/user?accountId=000000000000000000000000")

		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := testClient.User.Delete(context.Background(), "000000000000000000000000")
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
		testMethod(t, r, http.MethodGet)
		testRequestURL(t, r, "/rest/api/2/user/groups?accountId=000000000000000000000000")

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `[{"name":"jira-software-users","self":"http://www.example.com/jira/rest/api/2/user?accountId=000000000000000000000000"}]`)
	})

	if groups, _, err := testClient.User.GetGroups(context.Background(), "000000000000000000000000"); err != nil {
		t.Errorf("Error given: %s", err)
	} else if groups == nil {
		t.Error("Expected user groups. []UserGroup is nil")
	}
}

func TestUserService_GetCurrentUser(t *testing.T) {
	setup()
	defer teardown()
	testMux.HandleFunc("/rest/api/3/myself", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testRequestURL(t, r, "/rest/api/3/myself")

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"self": "https://your-domain.atlassian.net/rest/api/3/user?accountId=5b10a2844c20165700ede21g",
			"key": "",
			"accountId": "5b10a2844c20165700ede21g",
			"accountType": "atlassian",
			"name": "",
			"emailAddress": "mia@example.com",
			"avatarUrls": {
			  "48x48": "https://avatar-management--avatars.server-location.prod.public.atl-paas.net/initials/MK-5.png?size=48&s=48",
			  "24x24": "https://avatar-management--avatars.server-location.prod.public.atl-paas.net/initials/MK-5.png?size=24&s=24",
			  "16x16": "https://avatar-management--avatars.server-location.prod.public.atl-paas.net/initials/MK-5.png?size=16&s=16",
			  "32x32": "https://avatar-management--avatars.server-location.prod.public.atl-paas.net/initials/MK-5.png?size=32&s=32"
			},
			"displayName": "Mia Krystof",
			"active": true,
			"timeZone": "Australia/Sydney",
			"groups": {
			  "size": 3,
			  "items": []
			},
			"applicationRoles": {
			  "size": 1,
			  "items": []
			}
		  }`)
	})

	if user, _, err := testClient.User.GetCurrentUser(context.Background()); err != nil {
		t.Errorf("Error given: %s", err)

	} else if user == nil {
		t.Error("Expected user groups. []UserGroup is nil")

	} else if user.EmailAddress != "mia@example.com" || !user.Active || user.DisplayName != "Mia Krystof" {
		t.Errorf("User JSON deserialized incorrectly")
	}
}

func TestUserService_Find_Success(t *testing.T) {
	setup()
	defer teardown()
	testMux.HandleFunc("/rest/api/2/user/search", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testRequestURL(t, r, "/rest/api/2/user/search?query=fred@example.com")

		fmt.Fprint(w, `[{"self":"http://www.example.com/jira/rest/api/2/user?accountId=000000000000000000000000","key":"fred",
        "name":"fred","emailAddress":"fred@example.com","avatarUrls":{"48x48":"http://www.example.com/jira/secure/useravatar?size=large&ownerId=fred",
        "24x24":"http://www.example.com/jira/secure/useravatar?size=small&ownerId=fred","16x16":"http://www.example.com/jira/secure/useravatar?size=xsmall&ownerId=fred",
        "32x32":"http://www.example.com/jira/secure/useravatar?size=medium&ownerId=fred"},"displayName":"Fred F. User","active":true,"timeZone":"Australia/Sydney","groups":{"size":3,"items":[
        {"name":"jira-user","self":"http://www.example.com/jira/rest/api/2/group?groupname=jira-user"},{"name":"jira-admin",
        "self":"http://www.example.com/jira/rest/api/2/group?groupname=jira-admin"},{"name":"important","self":"http://www.example.com/jira/rest/api/2/group?groupname=important"
        }]},"applicationRoles":{"size":1,"items":[]},"expand":"groups,applicationRoles"}]`)
	})

	if user, _, err := testClient.User.Find(context.Background(), "fred@example.com"); err != nil {
		t.Errorf("Error given: %s", err)
	} else if user == nil {
		t.Error("Expected user. User is nil")
	}
}

func TestUserService_Find_SuccessParams(t *testing.T) {
	setup()
	defer teardown()
	testMux.HandleFunc("/rest/api/2/user/search", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testRequestURL(t, r, "/rest/api/2/user/search?query=fred@example.com&startAt=100&maxResults=1000")

		fmt.Fprint(w, `[{"self":"http://www.example.com/jira/rest/api/2/user?query=fred","key":"fred",
        "name":"fred","emailAddress":"fred@example.com","avatarUrls":{"48x48":"http://www.example.com/jira/secure/useravatar?size=large&ownerId=fred",
        "24x24":"http://www.example.com/jira/secure/useravatar?size=small&ownerId=fred","16x16":"http://www.example.com/jira/secure/useravatar?size=xsmall&ownerId=fred",
        "32x32":"http://www.example.com/jira/secure/useravatar?size=medium&ownerId=fred"},"displayName":"Fred F. User","active":true,"timeZone":"Australia/Sydney","groups":{"size":3,"items":[
        {"name":"jira-user","self":"http://www.example.com/jira/rest/api/2/group?groupname=jira-user"},{"name":"jira-admin",
        "self":"http://www.example.com/jira/rest/api/2/group?groupname=jira-admin"},{"name":"important","self":"http://www.example.com/jira/rest/api/2/group?groupname=important"
        }]},"applicationRoles":{"size":1,"items":[]},"expand":"groups,applicationRoles"}]`)
	})

	if user, _, err := testClient.User.Find(context.Background(), "fred@example.com", WithStartAt(100), WithMaxResults(1000)); err != nil {
		t.Errorf("Error given: %s", err)
	} else if user == nil {
		t.Error("Expected user. User is nil")
	}
}

func TestUserService_FindAssignable_Success(t *testing.T) {
	setup()
	defer teardown()
	testMux.HandleFunc("/rest/api/2/user/assignable/search", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testRequestURL(t, r, "/rest/api/2/user/assignable/search?username=fred@example.com&project=DAPHNE&issueKey=DAPHNE-11&maxResults=100&startAt=0")

		fmt.Fprint(w, `[{"self":"http://www.example.com/jira/rest/api/2/user?accountId=000000000000000000000000","key":"fred",
        "name":"fred","emailAddress":"fred@example.com","avatarUrls":{"48x48":"http://www.example.com/jira/secure/useravatar?size=large&ownerId=fred",
        "24x24":"http://www.example.com/jira/secure/useravatar?size=small&ownerId=fred","16x16":"http://www.example.com/jira/secure/useravatar?size=xsmall&ownerId=fred",
        "32x32":"http://www.example.com/jira/secure/useravatar?size=medium&ownerId=fred"},"displayName":"Fred F. User","active":true,"timeZone":"Australia/Sydney","groups":{"size":3,"items":[
        {"name":"jira-user","self":"http://www.example.com/jira/rest/api/2/group?groupname=jira-user"},{"name":"jira-admin",
        "self":"http://www.example.com/jira/rest/api/2/group?groupname=jira-admin"},{"name":"important","self":"http://www.example.com/jira/rest/api/2/group?groupname=important"
        }]},"applicationRoles":{"size":1,"items":[]},"expand":"groups,applicationRoles"}]`)
	})

	if user, _, err := testClient.User.FindAssignable(context.Background(), WithUsername("fred@example.com"), WithProject("DAPHNE"), WithIssueKey("DAPHNE-11"), WithMaxResults(100), WithStartAt(0)); err != nil {
		t.Errorf("Error given: %s", err)
	} else if user == nil {
		t.Error("Expected user. User is nil")
	}
}
