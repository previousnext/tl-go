package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/previousnext/tl-go/internal/api/types"
)

// JiraClientInterface defines the methods for interacting with Jira API
// See https://developer.atlassian.com/cloud/jira/platform/rest/v3/intro/
type JiraClientInterface interface {
	AddWorkLog(worklog types.WorklogRecord) error
	BulkFetchIssues(issueKeys []string) (BulkFetchIssuesResponse, error)
}

type HttpClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

type JiraClient struct {
	httpClient HttpClientInterface
	params     types.JiraClientParams
}

type ErrorResponse struct {
	ErrorMessages []string `json:"errorMessages"`
}

func NewJiraClient(httpClient HttpClientInterface, params types.JiraClientParams) *JiraClient {
	return &JiraClient{
		httpClient: httpClient,
		params:     params,
	}
}

func (c *JiraClient) doRequest(method, url string, bodyBuf io.Reader) (io.ReadCloser, error) {
	req, err := http.NewRequest(method, url, bodyBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for %s: %w", url, err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.params.Username, c.params.APIToken)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request for %s: %w", url, err)
	}
	if resp.StatusCode >= http.StatusBadRequest {
		var errorResp ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errorResp)
		if err != nil {
			return nil, fmt.Errorf("error decoding JSON: %w", err)
		}
		errMsg := errorResp.ErrorMessages[0]
		return nil, fmt.Errorf("api request failed with status code: [%d] %s", resp.StatusCode, errMsg)
	}
	return resp.Body, nil
}
