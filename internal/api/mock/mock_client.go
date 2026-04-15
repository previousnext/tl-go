package mock

import (
	"github.com/previousnext/tl-go/internal/api"
	"github.com/previousnext/tl-go/internal/api/types"
)

type JiraClient struct {
	api.JiraClientInterface
}

func (j *JiraClient) AddWorkLog(worklog types.WorklogRecord) error {
	return nil
}

func (j *JiraClient) GetUpdatedWorklogIDs(sinceMillis int64) ([]types.WorklogChange, error) {
	return nil, nil
}

func (j *JiraClient) BulkGetWorklogs(ids []int64) ([]types.Worklog, error) {
	return nil, nil
}

func (j *JiraClient) GetWorklogProperty(issueID string, worklogID string, propertyKey string) (*types.EntityPropertyResponse, error) {
	return nil, nil
}
