package stop

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/model"
	"github.com/previousnext/tl-go/internal/service"
)

type mockCurrentTimeEntryStorage struct {
	entry  *model.TimeEntry
	retErr error
}

func (m *mockCurrentTimeEntryStorage) StopTimeEntry() (*model.TimeEntry, error) {
	return m.entry, m.retErr
}

// Implement other methods to satisfy the interface (no-op)
func (m *mockCurrentTimeEntryStorage) StartTimeEntry(issueKey string) error { return nil }
func (m *mockCurrentTimeEntryStorage) PauseTimeEntry() error                { return nil }
func (m *mockCurrentTimeEntryStorage) GetCurrentTimeEntry() (*model.CurrentTimeEntry, error) {
	return nil, nil
}
func (m *mockCurrentTimeEntryStorage) SaveCurrentTimeEntry(entry *model.CurrentTimeEntry) error {
	return nil
}

type mockSyncService struct {
	calledWith string
}

func (m *mockSyncService) SyncIssue(issueKey string, options ...service.SyncOption) (*model.Issue, error) {
	m.calledWith = issueKey
	return &model.Issue{Key: issueKey}, nil
}
func (m *mockSyncService) SyncIssues(issueKeys []string) error { return nil }

func TestStopCommand_CallsSyncIssue(t *testing.T) {
	mockStorage := &mockCurrentTimeEntryStorage{
		entry: &model.TimeEntry{
			Model:    gorm.Model{CreatedAt: time.Now()},
			IssueKey: "PNX-123",
			Duration: 10 * time.Minute,
		},
	}
	mockSync := &mockSyncService{}
	cmd := NewCommand(
		func() db.CurrentTimeEntryStorageInterface { return mockStorage },
		func() service.SyncInterface { return mockSync },
	)
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})
	assert.NoError(t, cmd.Execute())
	assert.Equal(t, "PNX-123", mockSync.calledWith)
	out := buf.String()
	assert.Contains(t, out, "Stopped time entry for PNX-123")
}
