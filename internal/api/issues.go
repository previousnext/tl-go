package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"
)

type BulkFetchIssuesResponse struct {
	Issues []struct {
		ID     string `json:"id"`
		Key    string `json:"key"`
		Fields struct {
			Summary string `json:"summary"`
			Project struct {
				ID              string `json:"id"`
				Key             string `json:"key"`
				Name            string `json:"name"`
				ProjectCategory struct {
					Description string `json:"description"`
				} `json:"projectCategory"`
			} `json:"project"`
		} `json:"fields"`
	} `json:"issues"`
}

func (c *JiraClient) BulkFetchIssues(issueKeys []string) (BulkFetchIssuesResponse, error) {
	var issuesResp BulkFetchIssuesResponse
	url := c.params.BaseURL + "/rest/api/3/issue/bulkfetch"
	reqBody, err := generateBulkFetchIssuesBody(issueKeys)
	if err != nil {
		return issuesResp, err
	}
	respBody, err := c.doPostRequest(url, reqBody)
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

func generateBulkFetchIssuesBody(issueKeys []string) (bytes.Buffer, error) {
	issuesTmpl := `{
  "issueIdsOrKeys": {{.issueKeys | toJSON }},
  "fields": ["summary", "project"]
}`
	var buf bytes.Buffer

	t, err := template.New("issues").Funcs(template.FuncMap{
		"toJSON": jsonMarshal,
	}).Parse(issuesTmpl)
	if err != nil {
		return buf, fmt.Errorf("failed to parse body template: %w", err)
	}

	if err := t.Execute(&buf, issueKeys); err != nil {
		return buf, fmt.Errorf("failed to execute body template: %w", err)
	}

	return buf, nil
}

func jsonMarshal(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
