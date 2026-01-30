package types

import (
	"net/http"
	"time"
)

type JiraClientInterface interface {
	AddWorkLog(worklog WorklogRecord) error
}

type HttpClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

type JiraClientParams struct {
	BaseURL  string
	Username string
	APIToken string
}

type WorklogRecord struct {
	IssueKey string
	Started  time.Time
	Duration time.Duration
	Comment  string
}
