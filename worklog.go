package jira

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

// WorklogBean that returns worklog api uses this time format.
// Watch the Returns section: http://developer.tempo.io/doc/timesheets/api/rest/latest/#848933329
const TTWorklogTimeFormat = "2006-01-02T15:04:05.000"

// Date format for api request
// http://developer.tempo.io/doc/timesheets/api/rest/latest/#848933329
const TTWorklogDateFormat = "2006-01-02"

// TempoTimesheetsWorklogService handles tempo timesheets worklogs in JIRA rest API.
// See http://developer.tempo.io/doc/timesheets/api/rest/latest/#848933329
type TempoTimesheetsWorklogService struct {
	client *Client
}

// Worklog request options
type TTWorkLogOptions struct {
	Project  string         `url:"projectKey"`
	Username string         `url:"username"`
	DateFrom *TTWorklogDate `url:"dateFrom"`
	DateTo   *TTWorklogDate `url:"dateTo"`
}

// WorkLog time custom format
type TTWorklogTime struct {
	time.Time
}

// WorkLog date custom format
type TTWorklogDate struct {
	time.Time
}

// TimeSheet WorkLog represents one entry of a WorkLog
type TTWorkLog struct {
	Self             string         `json:"self,omitempty" structs:"self,omitempty"`
	ID               int64          `json:"id,omitempty" structs:"id,omitempty"`
	Issue            *TTIssue       `json:"issue,omitempty" structs:"issue,omitempty"`
	Author           *User          `json:"author,omitempty" structs:"author,omitempty"`
	TimeSpent        string         `json:"timeSpent,omitempty" structs:"timeSpent,omitempty"`
	TimeSpentSeconds int            `json:"timeSpentSeconds,omitempty" structs:"timeSpentSeconds,omitempty"`
	Comment          string         `json:"comment,omitempty" structs:"comment,omitempty"`
	Created          *TTWorklogTime `json:"dateCreated,omitempty" structs:"dateCreated,omitempty"`
	Updated          *TTWorklogTime `json:"dateUpdated,omitempty" structs:"dateUpdated,omitempty"`
	UpdateAuthor     *User          `json:"updateAuthor,omitempty" structs:"updateAuthor,omitempty"`
	Started          *TTWorklogTime `json:"dateStarted,omitempty" structs:"dateStarted,omitempty"`
}

// TimeSheet Issue object represents a JIRA issue.
type TTIssue struct {
	Self                     string     `json:"self,omitempty" structs:"self,omitempty"`
	ID                       int64      `json:"id,omitempty" structs:"id,omitempty"`
	ProjectID                int64      `json:"projectId,omitempty" structs:"projectId,omitempty"`
	Key                      string     `json:"key,omitempty" structs:"key,omitempty"`
	RemainingEstimateSeconds int64      `json:"remainingEstimateSeconds,omitempty" structs:"remainingEstimateSeconds,omitempty"`
	IssueType                *IssueType `json:"issueType" structs:"issueType,omitempty"`
	Summary                  string     `json:"summary,omitempty" structs:"summary,omitempty"`
}

// UnmarshalJSON will transform the TimeSheet WorkLog time into a time.Time
// during the transformation of the JSON response
func (t *TTWorklogTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	// Ignore null, like in the main JSON package.
	if s == "null" {
		t.Time = time.Time{}
		return
	}

	t.Time, err = time.Parse(TTWorklogTimeFormat, s)
	return
}

// UnmarshalJSON will transform the TimeSheet WorkLog date into a time.Time
// during the transformation of the JSON response
func (d *TTWorklogDate) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	// Ignore null, like in the main JSON package.
	if s == "null" {
		d.Time = time.Time{}
		return
	}

	d.Time, err = time.Parse(TTWorklogDateFormat, s)
	if err != nil {
		return err
	}
	return
}

// GetWorkLogs returns worklogs for a user on date range
func (w *TempoTimesheetsWorklogService) GetWorkLogs(params *url.Values) ([]TTWorkLog, *Response, error) {
	apiEndpoint := fmt.Sprintf("/rest/tempo-timesheets/3/worklogs?%s", params.Encode())
	req, err := w.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	var result []TTWorkLog
	resp, err := w.client.Do(req, &result)
	if err != nil {
		err = NewJiraError(resp, err)
	}

	return result, resp, err
}
