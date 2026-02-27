package stop

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/previousnext/tl-go/internal/model"
	"github.com/previousnext/tl-go/internal/service"
)

type mockTimerEntryStorage struct {
	entry  *model.TimeEntry
	retErr error
}

func (m *mockTimerEntryStorage) StopTimeEntry() (*model.TimeEntry, error) {
	return m.entry, m.retErr
}

// Implement other methods to satisfy the interface (no-op)
func (m *mockTimerEntryStorage) StartTimeEntry(issueKey string) error { return nil }
func (m *mockTimerEntryStorage) PauseTimeEntry() error                { return nil }
func (m *mockTimerEntryStorage) GetTimerEntry() (*model.TimerEntry, error) {
	return nil, nil
}
func (m *mockTimerEntryStorage) SaveTimerEntry(entry *model.TimerEntry) error {
	return nil
}
func (m *mockTimerEntryStorage) FindAllTimerEntries() ([]*model.TimerEntry, error) {
	return nil, nil
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
	mockStorage := &mockTimerEntryStorage{
		entry: &model.TimeEntry{
			Model:    gorm.Model{CreatedAt: time.Now()},
			IssueKey: "PNX-123",
			Duration: 10 * time.Minute,
		},
	}
	mockSync := &mockSyncService{}
	cmd := NewCommand(
		func() service.TimerEntryStorageInterface { return mockStorage },
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
