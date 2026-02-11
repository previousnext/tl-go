package api

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"

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

func generateWorklogPayload(worklog types.WorklogRecord) (*bytes.Buffer, error) {
	payloadTmpl := `{
  "comment": {
    "content": [
      {
        "content": [
          {
            "text": "{{ .comment }}",
            "type": "text"
          }
        ],
        "type": "paragraph"
      }
    ],
    "type": "doc",
    "version": 1
  },
  "started": "{{ .started }}",
  "timeSpentSeconds": {{ .timeSpentSeconds }}
}`
	var buf bytes.Buffer

	t, err := template.New("payload").Parse(payloadTmpl)
	if err != nil {
		return &buf, fmt.Errorf("failed to parse body template: %w", err)
	}

	data := map[string]interface{}{
		"comment":          worklog.Comment,
		"started":          worklog.Started.Format(DateFormat),
		"timeSpentSeconds": uint(worklog.Duration.Seconds()),
	}

	if err := t.Execute(&buf, data); err != nil {
		return &buf, fmt.Errorf("failed to execute body template: %w", err)
	}

	return &buf, nil
}
