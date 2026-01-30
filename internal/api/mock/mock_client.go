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
