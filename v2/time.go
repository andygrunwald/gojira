package jira

import "time"

// UnmarshalJSON will transform the Jira time into a time.Time
// during the transformation of the Jira JSON response
func TimeUnmarshalJSON(b []byte) (*time.Time, error) {
	// Ignore null, like in the main JSON package.
	if string(b) == "null" {
		return nil, nil
	}
	ti, err := time.Parse("2006-01-02T15:04:05.999-0700", string(b))
	return &ti, err
}

// MarshalJSON will transform the time.Time into a Jira time
// during the creation of a Jira request
func TimeMarshalJSON(t time.Time) ([]byte, error) {
	return []byte(time.Time(t).Format("2006-01-02T15:04:05.000-0700")), nil
}

// UnmarshalJSON will transform the Jira date into a time.Time
// during the transformation of the Jira JSON response
func DateUnmarshalJSON(b []byte) (*time.Time, error) {
	// Ignore null, like in the main JSON package.
	if string(b) == "null" {
		return nil, nil
	}
	ti, err := time.Parse("\"2006-01-02\"", string(b))
	return &ti, err
}

// MarshalJSON will transform the Date object into a short
// date string as Jira expects during the creation of a
// Jira request
func DateMarshalJSON(t time.Time) ([]byte, error) {
	return []byte(t.Format("2006-01-02")), nil
}
