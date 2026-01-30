package send

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/previousnext/tl-go/internal/api"
	"github.com/previousnext/tl-go/internal/api/mock"
	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/db/mocks"
	"github.com/previousnext/tl-go/internal/model"
)

func TestListUnsentEntries(t *testing.T) {
	expected := []*model.TimeEntry{
		{Model: gorm.Model{ID: 1}, Description: "Test entry 1"},
		{Model: gorm.Model{ID: 2}, Description: "Test entry 2"},
	}

	mockRepo := &mocks.MockRepository{
		FindUnsentTimeEntriesFunc: func() ([]*model.TimeEntry, error) {
			return expected, nil
		},
	}

	cmd := NewCommand(
		func() db.RepositoryInterface { return mockRepo },
		func() api.JiraClientInterface { return &mock.JiraClient{} },
	)
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()

	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Added 2 worklogs to Jira")
}

func TestListUnsentEntries_NoEntries(t *testing.T) {
	expected := []*model.TimeEntry{}

	mockRepo := &mocks.MockRepository{
		FindUnsentTimeEntriesFunc: func() ([]*model.TimeEntry, error) {
			return expected, nil
		},
	}

	cmd := NewCommand(
		func() db.RepositoryInterface { return mockRepo },
		func() api.JiraClientInterface { return &mock.JiraClient{} },
	)
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()

	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "No unsent time entries found")
}
