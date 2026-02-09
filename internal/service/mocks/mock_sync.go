package mocks

import (
	"github.com/previousnext/tl-go/internal/model"
	"github.com/previousnext/tl-go/internal/service"
)

type MockSync struct {
	service.SyncInterface
	SyncIssueFunc func() (*model.Issue, error)
}

func (s *MockSync) SyncIssue(issueKey string, options ...service.SyncOption) (*model.Issue, error) {
	if s.SyncIssueFunc != nil {
		return s.SyncIssueFunc()
	}
	return nil, nil
}

func (s *MockSync) SyncIssues(issueKeys []string) error {
	return nil
}
