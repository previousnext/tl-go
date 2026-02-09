package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/template"
)

type BulkFetchIssuesResponse struct {
	Issues []IssueResponse `json:"issues"`
}

type IssueResponse struct {
	ID     string      `json:"id"`
	Key    string      `json:"key"`
	Fields IssueFields `json:"fields"`
}

type IssueFields struct {
	Summary string `json:"summary"`
	Project struct {
		ID              string `json:"id"`
		Key             string `json:"key"`
		Name            string `json:"name"`
		ProjectCategory struct {
			Description string `json:"description"`
		} `json:"projectCategory"`
	} `json:"project"`
}

func (c *JiraClient) FetchIssue(issueKey string) (IssueResponse, error) {
	var issuesResp IssueResponse
	url := c.params.BaseURL + "/rest/api/3/issue/" + issueKey + "?fields=summary,project"
	respBody, err := c.doRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return issuesResp, err
	}
	//nolint:errcheck
	defer respBody.Close()
	err = json.NewDecoder(respBody).Decode(&issuesResp)
	if err != nil {
		return issuesResp, fmt.Errorf("error decoding JSON: %w", err)
	}

	return issuesResp, nil

}

func (c *JiraClient) BulkFetchIssues(issueKeys []string) (BulkFetchIssuesResponse, error) {
	var issuesResp BulkFetchIssuesResponse
	url := c.params.BaseURL + "/rest/api/3/issue/bulkfetch"
	reqBody, err := generateBulkFetchIssuesBody(issueKeys)
	if err != nil {
		return issuesResp, err
	}
	respBody, err := c.doRequest(http.MethodPost, url, reqBody)
	if err != nil {
		return issuesResp, err
	}
	//nolint:errcheck
	defer respBody.Close()
	err = json.NewDecoder(respBody).Decode(&issuesResp)
	if err != nil {
		return issuesResp, fmt.Errorf("error decoding JSON: %w", err)
	}

	return issuesResp, nil
}

func generateBulkFetchIssuesBody(issueKeys []string) (io.Reader, error) {
	issuesTmpl := `{
  "issueIdsOrKeys": {{.issueKeys | toJSON }},
  "fields": ["summary", "project"]
}`
	var buf bytes.Buffer

	t, err := template.New("issues").Funcs(template.FuncMap{
		"toJSON": jsonMarshal,
	}).Parse(issuesTmpl)
	if err != nil {
		return &buf, fmt.Errorf("failed to parse body template: %w", err)
	}

	data := map[string]interface{}{
		"issueKeys": issueKeys,
	}
	if err := t.Execute(&buf, data); err != nil {
		return &buf, fmt.Errorf("failed to execute body template: %w", err)
	}

	return &buf, nil
}

func jsonMarshal(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
