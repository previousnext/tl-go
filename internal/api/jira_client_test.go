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

// RoundTripFunc type to allow custom http.RoundTripper
type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

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
	assert.Contains(t, string(capturedBody), "2024-06-01T10:00:00Z")
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
	assert.Contains(t, payload, worklog.IssueKey)
	assert.Contains(t, payload, worklog.Comment)
	assert.Contains(t, payload, "2024-06-01T10:00:00Z")
	assert.Contains(t, payload, "7200") // 2 hours in seconds
}
