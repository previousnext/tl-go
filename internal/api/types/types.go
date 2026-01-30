package types

import (
	"time"
)

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
