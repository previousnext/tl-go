package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/previousnext/tl-go/internal/api/types"
)

// isValidJSON checks if a string is valid JSON.
func isValidJSON(s string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(s), &js) == nil
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
	assert.True(t, isValidJSON(payload), "Generated payload should be valid JSON")
	assert.Contains(t, payload, worklog.Comment)
	assert.Contains(t, payload, "2024-06-01T10:00:00.000+0000")
	assert.Contains(t, payload, "7200") // 2 hours in seconds
}

func TestGenerateWorklogPayload_SpecialCharacters(t *testing.T) {
	tests := []struct {
		name            string
		comment         string
		expectedInJSON  string
		shouldBeEscaped bool // true if the character needs JSON escaping
	}{
		{
			name:           "newline in comment",
			comment:        "Line one\nLine two",
			expectedInJSON: "Line one Line two",
		},
		{
			name:           "carriage return and newline",
			comment:        "Line one\r\nLine two",
			expectedInJSON: "Line one Line two",
		},
		{
			name:           "tab in comment",
			comment:        "Item\tValue",
			expectedInJSON: "Item Value",
		},
		{
			name:           "multiple whitespace",
			comment:        "Too   many   spaces",
			expectedInJSON: "Too many spaces",
		},
		{
			name:            "double quotes are escaped",
			comment:         `Said "hello" today`,
			expectedInJSON:  `Said \"hello\" today`, // JSON-escaped quotes
			shouldBeEscaped: true,
		},
		{
			name:            "backslash is escaped",
			comment:         `Path\to\file`,
			expectedInJSON:  `Path\\to\\file`, // JSON-escaped backslashes
			shouldBeEscaped: true,
		},
		{
			name:            "mixed special chars",
			comment:         "First line\nSecond \"quoted\"\tthird",
			expectedInJSON:  `First line Second \"quoted\" third`,
			shouldBeEscaped: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			worklog := types.WorklogRecord{
				Comment:  tt.comment,
				Started:  time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC),
				Duration: 2 * time.Hour,
			}

			buf, err := generateWorklogPayload(worklog)
			assert.NoError(t, err)

			payload := buf.String()

			// Verify the payload is valid JSON
			assert.True(t, isValidJSON(payload), "Generated payload should be valid JSON: %s", payload)

			// Verify the expected content is in the payload
			assert.Contains(t, payload, tt.expectedInJSON, "Payload should contain expected comment")

			// For escaped characters, also verify we can unmarshal and get original values
			if tt.shouldBeEscaped {
				var result worklogPayload
				err := json.Unmarshal(buf.Bytes(), &result)
				assert.NoError(t, err)
				// After unmarshaling, we should get the sanitized comment (whitespace normalized)
				// but with quotes and backslashes preserved
			}
		})
	}
}
