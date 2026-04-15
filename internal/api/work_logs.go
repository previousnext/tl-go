package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/previousnext/tl-go/internal/api/types"
)

func (c *JiraClient) AddWorkLog(worklog types.WorklogRecord) error {
	url := c.params.BaseURL + "/rest/api/3/issue/" + worklog.IssueKey + "/worklog"
	bodyBuf, err := generateWorklogPayload(worklog)
	if err != nil {
		return err
	}
	respBody, err := c.doRequest(http.MethodPost, url, bodyBuf)
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer respBody.Close()
	return nil
}

// worklogPayload represents the JSON body sent to the Jira Add Worklog API.
type worklogPayload struct {
	Comment          worklogComment   `json:"comment"`
	Started          string           `json:"started"`
	TimeSpentSeconds uint             `json:"timeSpentSeconds"`
	Properties       []entityProperty `json:"properties,omitempty"`
}

type worklogComment struct {
	Type    string             `json:"type"`
	Version int                `json:"version"`
	Content []worklogParagraph `json:"content"`
}

type worklogParagraph struct {
	Type    string        `json:"type"`
	Content []worklogText `json:"content"`
}

type worklogText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type entityProperty struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type aiTimeSavedPropertyValue struct {
	DurationSeconds uint `json:"durationSeconds"`
}

func generateWorklogPayload(worklog types.WorklogRecord) (*bytes.Buffer, error) {
	payload := worklogPayload{
		Comment: worklogComment{
			Type:    "doc",
			Version: 1,
			Content: []worklogParagraph{
				{
					Type: "paragraph",
					Content: []worklogText{
						{
							Type: "text",
							Text: worklog.Comment,
						},
					},
				},
			},
		},
		Started:          worklog.Started.Format(DateFormat),
		TimeSpentSeconds: uint(worklog.Duration.Seconds()),
	}

	if worklog.AISavedDuration > 0 {
		payload.Properties = []entityProperty{
			{
				Key: "ai-time-saved",
				Value: aiTimeSavedPropertyValue{
					DurationSeconds: uint(worklog.AISavedDuration.Seconds()),
				},
			},
		}
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		return &buf, fmt.Errorf("failed to encode worklog payload: %w", err)
	}

	return &buf, nil
}

// GetUpdatedWorklogIDs returns all worklog IDs updated since the given timestamp (Unix milliseconds).
// It paginates through all pages automatically.
func (c *JiraClient) GetUpdatedWorklogIDs(sinceMillis int64) ([]types.WorklogChange, error) {
	var allChanges []types.WorklogChange

	url := fmt.Sprintf("%s/rest/api/3/worklog/updated?since=%d", c.params.BaseURL, sinceMillis)
	for {
		respBody, err := c.doRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get updated worklogs: %w", err)
		}

		var resp types.UpdatedWorklogsResponse
		if err := json.NewDecoder(respBody).Decode(&resp); err != nil {
			respBody.Close()
			return nil, fmt.Errorf("failed to decode updated worklogs response: %w", err)
		}
		respBody.Close()

		allChanges = append(allChanges, resp.Values...)

		if resp.LastPage {
			break
		}
		url = resp.NextPage
	}

	return allChanges, nil
}

// BulkGetWorklogs fetches full worklog details for a list of worklog IDs.
// The Jira API accepts up to 1000 IDs per request.
func (c *JiraClient) BulkGetWorklogs(ids []int64) ([]types.Worklog, error) {
	var allWorklogs []types.Worklog

	// Process in batches of 1000
	for i := 0; i < len(ids); i += 1000 {
		end := i + 1000
		if end > len(ids) {
			end = len(ids)
		}
		batch := ids[i:end]

		reqBody := struct {
			IDs []int64 `json:"ids"`
		}{IDs: batch}

		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(reqBody); err != nil {
			return nil, fmt.Errorf("failed to encode worklog IDs: %w", err)
		}

		url := c.params.BaseURL + "/rest/api/3/worklog/list"
		respBody, err := c.doRequest(http.MethodPost, url, &buf)
		if err != nil {
			return nil, fmt.Errorf("failed to bulk get worklogs: %w", err)
		}

		var worklogs []types.Worklog
		if err := json.NewDecoder(respBody).Decode(&worklogs); err != nil {
			respBody.Close()
			return nil, fmt.Errorf("failed to decode worklogs response: %w", err)
		}
		respBody.Close()

		allWorklogs = append(allWorklogs, worklogs...)
	}

	return allWorklogs, nil
}

// GetWorklogProperty fetches a single property value from a worklog.
// Returns nil if the property does not exist (404).
func (c *JiraClient) GetWorklogProperty(issueID string, worklogID string, propertyKey string) (*types.EntityPropertyResponse, error) {
	url := fmt.Sprintf("%s/rest/api/3/issue/%s/worklog/%s/properties/%s", c.params.BaseURL, issueID, worklogID, propertyKey)
	respBody, err := c.doRequest(http.MethodGet, url, nil)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get worklog property: %w", err)
	}
	defer respBody.Close()

	var resp types.EntityPropertyResponse
	if err := json.NewDecoder(respBody).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode worklog property response: %w", err)
	}

	return &resp, nil
}
