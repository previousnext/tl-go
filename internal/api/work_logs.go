package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/previousnext/tl-go/internal/api/types"
	"github.com/previousnext/tl-go/internal/util"
)

// worklogPayload represents the JSON structure for Jira worklog API.
type worklogPayload struct {
	Comment          adfDocument `json:"comment"`
	Started          string      `json:"started"`
	TimeSpentSeconds uint        `json:"timeSpentSeconds"`
}

// adfDocument represents an Atlassian Document Format document.
type adfDocument struct {
	Type    string         `json:"type"`
	Version int            `json:"version"`
	Content []adfParagraph `json:"content"`
}

// adfParagraph represents a paragraph in ADF.
type adfParagraph struct {
	Type    string        `json:"type"`
	Content []adfTextNode `json:"content"`
}

// adfTextNode represents a text node in ADF.
type adfTextNode struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

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

func generateWorklogPayload(worklog types.WorklogRecord) (*bytes.Buffer, error) {
	payload := worklogPayload{
		Comment: adfDocument{
			Type:    "doc",
			Version: 1,
			Content: []adfParagraph{
				{
					Type: "paragraph",
					Content: []adfTextNode{
						{
							Type: "text",
							Text: util.SanitizeComment(worklog.Comment),
						},
					},
				},
			},
		},
		Started:          worklog.Started.Format(DateFormat),
		TimeSpentSeconds: uint(worklog.Duration.Seconds()),
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		return &buf, err
	}

	return &buf, nil
}
