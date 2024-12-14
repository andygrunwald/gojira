package errors

import (
	"bytes"
	"fmt"
	"net/http"
)

const (
	HttpRequestCreationFailureMessageFormat = "failed to create http request: %w"
)

// https://developer.atlassian.com/cloud/jira/platform/rest/v2/#status-codes
// https://developer.atlassian.com/cloud/jira/platform/rest/v3/#status-codes
type Collection struct {
	ErrorMessages []string          `json:"errorMessages"`
	Errors        map[string]string `json:"errors"`
	Status        int32             `json:"status"`
}

type RequestError struct {
	JiraErrorCollection *JiraErrorCollection
	HttpResponse        *http.Response
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("jira error encountered. httpStatus: %s, jiraStatus: %d", e.HttpResponse.Status, e.JiraErrorCollection.Status)
}

// NewJiraError creates a new jira Error
func NewRequestError(resp *http.Response) *RequestError {
	jiraRequestError := &RequestError{}

	return jiraRequestError

	//defer resp.Body.Close()
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//
	//}
	//
	//contentType := resp.Header.Get("Content-Type")
	//if strings.HasPrefix(contentType, "application/json") {
	//	err = json.Unmarshal(body, &jerr)
	//	if err != nil {
	//		httpError = fmt.Errorf("could not parse JSON: %w", httpError.Error())
	//		errors.
	//		return errors.Wrap(err, httpError.Error())
	//	}
	//} else {
	//	if httpError == nil {
	//		return fmt.Errorf("got response status %s:%s", resp.Status, string(body))
	//	}
	//	return errors.Wrap(httpError, fmt.Sprintf("%s: %s", resp.Status, string(body)))
	//}
	//
	//return &jerr
}

// Error is a short string representing the error
func (e *Collection) Error() string {
	if len(e.ErrorMessages) > 0 {
		// return fmt.Sprintf("%v", e.HTTPError)
		return fmt.Sprintf("%s: %v", e.ErrorMessages[0], e.HTTPError)
	}
	if len(e.Errors) > 0 {
		for key, value := range e.Errors {
			return fmt.Sprintf("%s - %s: %v", key, value, e.HTTPError)
		}
	}
	return e.HTTPError.Error()
}

// LongError is a full representation of the error as a string
func (e *Collection) LongError() string {
	var msg bytes.Buffer
	if e.HTTPError != nil {
		msg.WriteString("Original:\n")
		msg.WriteString(e.HTTPError.Error())
		msg.WriteString("\n")
	}
	if len(e.ErrorMessages) > 0 {
		msg.WriteString("Messages:\n")
		for _, v := range e.ErrorMessages {
			msg.WriteString(" - ")
			msg.WriteString(v)
			msg.WriteString("\n")
		}
	}
	if len(e.Errors) > 0 {
		for key, value := range e.Errors {
			msg.WriteString(" - ")
			msg.WriteString(key)
			msg.WriteString(" - ")
			msg.WriteString(value)
			msg.WriteString("\n")
		}
	}
	return msg.String()
}
