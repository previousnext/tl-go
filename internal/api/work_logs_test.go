package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/previousnext/tl-go/internal/api/types"
)

func TestJiraClient_AddWorkLog(t *testing.T) {
	var capturedRequest *http.Request
	var capturedBody []byte

	// Custom RoundTripper to capture request and return a fake response
	rt := RoundTripFunc(func(req *http.Request) *http.Response {
		capturedRequest = req
		body, _ := io.ReadAll(req.Body)
		capturedBody = body
		return &http.Response{
			StatusCode: http.StatusCreated,
			Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
			Header:     make(http.Header),
		}
	})

	httpClient := &http.Client{Transport: rt}
	jiraClient := NewJiraClient(httpClient, types.JiraClientParams{
		BaseURL:  "https://example.atlassian.net",
		Username: "user",
		APIToken: "token",
	})

	worklog := types.WorklogRecord{
		IssueKey: "PROJ-123",
		Comment:  "Worked on bug fix",
		Started:  time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC),
		Duration: 2 * time.Hour,
	}

	err := jiraClient.AddWorkLog(worklog)
	assert.NoError(t, err)
	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "POST", capturedRequest.Method)
	assert.Contains(t, capturedRequest.URL.Path, "/rest/api/3/issue/PROJ-123/worklog")
	assert.Equal(t, "application/json", capturedRequest.Header.Get("Content-Type"))
	assert.Contains(t, string(capturedBody), "Worked on bug fix")
	assert.Contains(t, string(capturedBody), "2024-06-01T10:00:00.000+0000")
	assert.Contains(t, string(capturedBody), "7200")
}

func TestGenerateWorklogPayload(t *testing.T) {
	worklog := types.WorklogRecord{
		Comment:  "Worked on bug fix",
		Started:  time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC),
		Duration: 2 * time.Hour,
	}

	buf, err := generateWorklogPayload(worklog)
	assert.NoError(t, err)
	payload := buf.String()
	fmt.Println(payload)
	assert.NotEmpty(t, payload)
	assert.Contains(t, payload, worklog.Comment)
	assert.Contains(t, payload, "2024-06-01T10:00:00.000+0000")
	assert.Contains(t, payload, "7200") // 2 hours in seconds
	assert.NotContains(t, payload, "properties")
}

func TestGenerateWorklogPayload_WithAITimeSaved(t *testing.T) {
	worklog := types.WorklogRecord{
		Comment:         "Worked on bug fix",
		Started:         time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC),
		Duration:        2 * time.Hour,
		AISavedDuration: 1 * time.Hour,
	}

	buf, err := generateWorklogPayload(worklog)
	assert.NoError(t, err)
	payload := buf.String()
	fmt.Println(payload)
	assert.NotEmpty(t, payload)
	assert.Contains(t, payload, worklog.Comment)
	assert.Contains(t, payload, "7200")
	assert.Contains(t, payload, `"properties"`)
	assert.Contains(t, payload, `"ai-time-saved"`)
	assert.Contains(t, payload, `"durationSeconds":3600`)
}

func TestGenerateWorklogPayload_SpecialCharsInComment(t *testing.T) {
	worklog := types.WorklogRecord{
		Comment:  `Fixed "bug" with {brackets} & <tags>`,
		Started:  time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC),
		Duration: 30 * time.Minute,
	}

	buf, err := generateWorklogPayload(worklog)
	assert.NoError(t, err)
	payload := buf.String()
	// encoding/json properly escapes special characters
	assert.Contains(t, payload, `Fixed \"bug\" with {brackets} \u0026 \u003ctags\u003e`)
}

func TestJiraClient_GetUpdatedWorklogIDs(t *testing.T) {
	rt := RoundTripFunc(func(req *http.Request) *http.Response {
		assert.Equal(t, "GET", req.Method)
		assert.Contains(t, req.URL.String(), "/rest/api/3/worklog/updated?since=1000")
		body := `{
			"values": [
				{"worklogId": 101, "updatedTime": 1001},
				{"worklogId": 102, "updatedTime": 1002}
			],
			"lastPage": true,
			"since": 1000,
			"until": 1002
		}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(body)),
			Header:     make(http.Header),
		}
	})

	httpClient := &http.Client{Transport: rt}
	jiraClient := NewJiraClient(httpClient, types.JiraClientParams{
		BaseURL:  "https://example.atlassian.net",
		Username: "user",
		APIToken: "token",
	})

	changes, err := jiraClient.GetUpdatedWorklogIDs(1000)
	assert.NoError(t, err)
	assert.Len(t, changes, 2)
	assert.Equal(t, int64(101), changes[0].WorklogID)
	assert.Equal(t, int64(102), changes[1].WorklogID)
}

func TestJiraClient_BulkGetWorklogs(t *testing.T) {
	var capturedBody []byte
	rt := RoundTripFunc(func(req *http.Request) *http.Response {
		assert.Equal(t, "POST", req.Method)
		assert.Contains(t, req.URL.Path, "/rest/api/3/worklog/list")
		capturedBody, _ = io.ReadAll(req.Body)
		body := `[
			{
				"id": "101",
				"issueId": "10001",
				"author": {"accountId": "abc123", "displayName": "Test User"},
				"started": "2024-06-01T10:00:00.000+0000",
				"timeSpentSeconds": 7200
			}
		]`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(body)),
			Header:     make(http.Header),
		}
	})

	httpClient := &http.Client{Transport: rt}
	jiraClient := NewJiraClient(httpClient, types.JiraClientParams{
		BaseURL:  "https://example.atlassian.net",
		Username: "user",
		APIToken: "token",
	})

	worklogs, err := jiraClient.BulkGetWorklogs([]int64{101})
	assert.NoError(t, err)
	assert.Len(t, worklogs, 1)
	assert.Equal(t, "101", worklogs[0].ID)
	assert.Equal(t, "10001", worklogs[0].IssueID)
	assert.Equal(t, "Test User", worklogs[0].Author.DisplayName)
	assert.Equal(t, 7200, worklogs[0].TimeSpentSeconds)
	assert.Contains(t, string(capturedBody), "101")
}

func TestJiraClient_GetWorklogProperty(t *testing.T) {
	rt := RoundTripFunc(func(req *http.Request) *http.Response {
		assert.Equal(t, "GET", req.Method)
		assert.Contains(t, req.URL.Path, "/rest/api/3/issue/10001/worklog/101/properties/ai-time-saved")
		body := `{
			"key": "ai-time-saved",
			"value": {"durationSeconds": 3600}
		}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(body)),
			Header:     make(http.Header),
		}
	})

	httpClient := &http.Client{Transport: rt}
	jiraClient := NewJiraClient(httpClient, types.JiraClientParams{
		BaseURL:  "https://example.atlassian.net",
		Username: "user",
		APIToken: "token",
	})

	prop, err := jiraClient.GetWorklogProperty("10001", "101", "ai-time-saved")
	assert.NoError(t, err)
	assert.NotNil(t, prop)
	assert.Equal(t, "ai-time-saved", prop.Key)
	assert.Contains(t, string(prop.Value), "3600")
}

func TestJiraClient_GetWorklogProperty_NotFound(t *testing.T) {
	rt := RoundTripFunc(func(req *http.Request) *http.Response {
		body := `{"errorMessages": ["not found"]}`
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(bytes.NewBufferString(body)),
			Header:     make(http.Header),
		}
	})

	httpClient := &http.Client{Transport: rt}
	jiraClient := NewJiraClient(httpClient, types.JiraClientParams{
		BaseURL:  "https://example.atlassian.net",
		Username: "user",
		APIToken: "token",
	})

	prop, err := jiraClient.GetWorklogProperty("10001", "101", "ai-time-saved")
	assert.NoError(t, err)
	assert.Nil(t, prop)
}
