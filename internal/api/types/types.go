package types

import (
	"encoding/json"
	"time"
)

type JiraClientParams struct {
	BaseURL  string
	Username string
	APIToken string
}

type WorklogRecord struct {
	IssueKey        string
	Started         time.Time
	Duration        time.Duration
	Comment         string
	AISavedDuration time.Duration
}

// UpdatedWorklogsResponse is the response from GET /rest/api/3/worklog/updated.
type UpdatedWorklogsResponse struct {
	Values   []WorklogChange `json:"values"`
	LastPage bool            `json:"lastPage"`
	NextPage string          `json:"nextPage"`
	Since    int64           `json:"since"`
	Until    int64           `json:"until"`
}

// WorklogChange represents a single changed worklog ID and its update timestamp.
type WorklogChange struct {
	WorklogID   int64 `json:"worklogId"`
	UpdatedTime int64 `json:"updatedTime"`
}

// Worklog represents a full worklog returned from POST /rest/api/3/worklog/list.
type Worklog struct {
	ID               string        `json:"id"`
	IssueID          string        `json:"issueId"`
	Author           WorklogAuthor `json:"author"`
	Started          string        `json:"started"`
	TimeSpentSeconds int           `json:"timeSpentSeconds"`
}

// WorklogAuthor represents the author of a worklog.
type WorklogAuthor struct {
	AccountID   string `json:"accountId"`
	DisplayName string `json:"displayName"`
}

// EntityPropertyResponse is the response from GET .../properties/{key}.
type EntityPropertyResponse struct {
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
}

// AITimeSavedPropertyValue is the value stored in the ai-time-saved worklog property.
type AITimeSavedPropertyValue struct {
	DurationSeconds int `json:"durationSeconds"`
}
