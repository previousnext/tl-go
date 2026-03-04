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
		func() db.TimeEntriesInterface { return mockRepo },
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
		func() db.TimeEntriesInterface { return mockRepo },
		func() api.JiraClientInterface { return &mock.JiraClient{} },
	)
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()

	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "No unsent time entries found")
}

func TestSendEntryByID(t *testing.T) {
	expectedEntry := &model.TimeEntry{
		Model:       gorm.Model{ID: 42},
		Description: "Test entry by ID",
		IssueKey:    "TEST-123",
		Duration:    3600,
		Sent:        false,
	}

	var updatedEntry *model.TimeEntry

	mockRepo := &mocks.MockRepository{
		FindTimeEntryFunc: func(id uint) (*model.TimeEntry, error) {
			if id == 42 {
				return expectedEntry, nil
			}
			return nil, gorm.ErrRecordNotFound
		},
		UpdateTimeEntryFunc: func(entry *model.TimeEntry) error {
			updatedEntry = entry
			return nil
		},
	}

	cmd := NewCommand(
		func() db.TimeEntriesInterface { return mockRepo },
		func() api.JiraClientInterface { return &mock.JiraClient{} },
	)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"42"})

	err := cmd.Execute()

	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Resent time entry ID 42 to Jira")
	assert.NotNil(t, updatedEntry)
	assert.True(t, updatedEntry.Sent)
}
