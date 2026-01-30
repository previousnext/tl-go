package api

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/previousnext/tl-go/internal/api/types"
)

type JiraClientInterface interface {
	AddWorkLog(worklog types.WorklogRecord) error
}

type HttpClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

type JiraClient struct {
	httpClient HttpClientInterface
	params     types.JiraClientParams
}

func NewJiraClient(httpClient HttpClientInterface, params types.JiraClientParams) *JiraClient {
	return &JiraClient{
		httpClient: httpClient,
		params:     params,
	}
}

func (c *JiraClient) AddWorkLog(worklog types.WorklogRecord) error {
	url := c.params.BaseURL + "/rest/api/3/issue/" + worklog.IssueKey + "/worklog"
	bodyBuf, err := generateWorklogPayload(worklog)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, url, &bodyBuf)
	if err != nil {
		return fmt.Errorf("failed to create request for %s: %w", url, err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.params.Username, c.params.APIToken)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request for %s: %w", url, err)
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create worklog for issue %s: received status code %d", worklog.IssueKey, resp.StatusCode)
	}

	return nil
}

func generateWorklogPayload(worklog types.WorklogRecord) (bytes.Buffer, error) {
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
  "timeSpentSeconds": {{ .timeSpentSeconds }},
}`
	var buf bytes.Buffer

	t, err := template.New("payload").Parse(payloadTmpl)
	if err != nil {
		return buf, fmt.Errorf("failed to parse body template: %w", err)
	}

	data := map[string]interface{}{
		"comment":          worklog.Comment,
		"started":          worklog.Started.Format(time.RFC3339),
		"timeSpentSeconds": uint(worklog.Duration.Seconds()),
	}

	if err := t.Execute(&buf, data); err != nil {
		return buf, fmt.Errorf("failed to execute body template: %w", err)
	}

	return buf, nil
}
