package errors

import (
	"errors"
	"fmt"
	"github.com/andygrunwald/go-jira/v2"
	"net/http"
	"strings"
	"testing"
)

func TestError_NewJiraError(t *testing.T) {
	jira.setup()
	defer jira.teardown()

	jira.testMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"errorMessages":["Issue does not exist or you do not have permission to see it."],"errors":{}}`)
	})

	req, _ := jira.testClient.NewRequest("GET", "/", nil)
	resp, _ := jira.testClient.Do(req, nil)

	err := NewJiraError(resp, errors.New("Original http error"))
	if err, ok := err.(*Error); !ok {
		t.Errorf("Expected jira Error. Got %s", err.Error())
	}

	if !strings.Contains(err.Error(), "Issue does not exist") {
		t.Errorf("Expected issue message. Got: %s", err.Error())
	}
}

func TestError_NoResponse(t *testing.T) {
	err := NewJiraError(nil, errors.New("Original http error"))

	msg := err.Error()
	if !strings.Contains(msg, "Original http error") {
		t.Errorf("Expected the original error message: Got\n%s\n", msg)
	}

	if !strings.Contains(msg, "No response") {
		t.Errorf("Expected the 'No response' error message: Got\n%s\n", msg)
	}
}

func TestError_NoJSON(t *testing.T) {
	jira.setup()
	defer jira.teardown()

	jira.testMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `Original message body`)
	})

	req, _ := jira.testClient.NewRequest("GET", "/", nil)
	resp, _ := jira.testClient.Do(req, nil)

	err := NewJiraError(resp, errors.New("Original http error"))
	msg := err.Error()

	if !strings.Contains(msg, "200 OK: Original message body: Original http error") {
		t.Errorf("Expected the HTTP status: Got\n%s\n", msg)
	}
}

func TestError_Unauthorized_NilError(t *testing.T) {
	jira.setup()
	defer jira.teardown()

	jira.testMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, `User is not authorized`)
	})

	req, _ := jira.testClient.NewRequest("GET", "/", nil)
	resp, _ := jira.testClient.Do(req, nil)

	err := NewJiraError(resp, nil)
	msg := err.Error()
	if !strings.Contains(msg, "401 Unauthorized:User is not authorized") {
		t.Errorf("Expected Unauthorized HTTP status: Got\n%s\n", msg)
	}
}

func TestError_BadJSON(t *testing.T) {
	jira.setup()
	defer jira.teardown()

	jira.testMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `<html>Not JSON</html>`)
	})

	req, _ := jira.testClient.NewRequest("GET", "/", nil)
	resp, _ := jira.testClient.Do(req, nil)

	err := NewJiraError(resp, errors.New("Original http error"))
	msg := err.Error()

	if !strings.Contains(msg, "could not parse JSON") {
		t.Errorf("Expected the 'could not parse JSON' error message: Got\n%s\n", msg)
	}
}

func TestError_NilOriginalMessage(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Expected an error message. Got a panic (%v)", r)
		}
	}()

	msgErr := &Error{
		HTTPError:     nil,
		ErrorMessages: []string{"Issue does not exist"},
		Errors: map[string]string{
			"issuetype": "issue type is required",
			"title":     "title is required",
		},
	}

	_ = msgErr.Error()
}

func TestError_NilOriginalMessageLongError(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Expected an error message. Got a panic (%v)", r)
		}
	}()

	msgErr := &Error{
		HTTPError:     nil,
		ErrorMessages: []string{"Issue does not exist"},
		Errors: map[string]string{
			"issuetype": "issue type is required",
			"title":     "title is required",
		},
	}

	_ = msgErr.LongError()
}

func TestError_ShortMessage(t *testing.T) {
	msgErr := &Error{
		HTTPError:     errors.New("Original http error"),
		ErrorMessages: []string{"Issue does not exist"},
		Errors: map[string]string{
			"issuetype": "issue type is required",
			"title":     "title is required",
		},
	}

	mapErr := &Error{
		HTTPError:     errors.New("Original http error"),
		ErrorMessages: nil,
		Errors: map[string]string{
			"issuetype": "issue type is required",
			"title":     "title is required",
		},
	}

	noErr := &Error{
		HTTPError:     errors.New("Original http error"),
		ErrorMessages: nil,
		Errors:        nil,
	}

	err := msgErr.Error()
	if err != "Issue does not exist: Original http error" {
		t.Errorf("Expected short message. Got %s", err)
	}

	err = mapErr.Error()
	if !(strings.Contains(err, "issue type is required") || strings.Contains(err, "title is required")) {
		t.Errorf("Expected short message. Got %s", err)
	}

	err = noErr.Error()
	if err != "Original http error" {
		t.Errorf("Expected original error message. Got %s", err)
	}
}

func TestError_LongMessage(t *testing.T) {
	longError := &Error{
		HTTPError:     errors.New("Original http error"),
		ErrorMessages: []string{"Issue does not exist."},
		Errors: map[string]string{
			"issuetype": "issue type is required",
			"title":     "title is required",
		},
	}

	msg := longError.LongError()
	if !strings.Contains(msg, "Original http error") {
		t.Errorf("Expected the error message: Got\n%s\n", msg)
	}

	if !strings.Contains(msg, "Issue does not exist") {
		t.Errorf("Expected the error message: Got\n%s\n", msg)
	}

	if !strings.Contains(msg, "title - title is required") {
		t.Errorf("Expected the error map: Got\n%s\n", msg)
	}
}
