package continuecmd

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/model"
)

type mockCurrentTimeEntryStorage struct {
	entry   *model.CurrentTimeEntry
	retErr  error
	updated *model.CurrentTimeEntry
}

func (m *mockCurrentTimeEntryStorage) StartTimeEntry(issueKey string) error     { return nil }
func (m *mockCurrentTimeEntryStorage) PauseTimeEntry() error                    { return nil }
func (m *mockCurrentTimeEntryStorage) StopTimeEntry() (*model.TimeEntry, error) { return nil, nil }
func (m *mockCurrentTimeEntryStorage) GetCurrentTimeEntry() (*model.CurrentTimeEntry, error) {
	return m.entry, m.retErr
}
func (m *mockCurrentTimeEntryStorage) SaveCurrentTimeEntry(entry *model.CurrentTimeEntry) error {
	m.updated = &model.CurrentTimeEntry{}
	*m.updated = *entry
	return nil
}

func TestContinueCommand_PausedEntry(t *testing.T) {
	start := time.Date(2026, 2, 18, 10, 0, 0, 0, time.Local)
	pause := start.Add(5 * time.Minute)
	mock := &mockCurrentTimeEntryStorage{
		entry: &model.CurrentTimeEntry{
			IssueKey:  "PNX-123",
			StartTime: start,
			Paused:    true,
			PauseTime: pause,
			Duration:  0,
		},
	}
	cmd := NewCommand(func() db.CurrentTimeEntryStorageInterface { return mock })
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})
	assert.NoError(t, cmd.Execute())
	out := buf.String()
	assert.Contains(t, out, "Current time entry has been continued.")
	assert.NotNil(t, mock.updated)
	assert.False(t, mock.updated.Paused)
	assert.Equal(t, 5*time.Minute, mock.updated.Duration)
	assert.True(t, mock.updated.StartTime.After(pause))
}

func TestContinueCommand_NoPausedEntry(t *testing.T) {
	mock := &mockCurrentTimeEntryStorage{
		entry: &model.CurrentTimeEntry{
			IssueKey: "PNX-123",
			Paused:   false,
		},
	}
	cmd := NewCommand(func() db.CurrentTimeEntryStorageInterface { return mock })
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})
	assert.NoError(t, cmd.Execute())
	out := buf.String()
	assert.Contains(t, out, "Current time entry is not paused.")
}

func TestContinueCommand_NoEntry(t *testing.T) {
	mock := &mockCurrentTimeEntryStorage{entry: nil, retErr: nil}
	cmd := NewCommand(func() db.CurrentTimeEntryStorageInterface { return mock })
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})
	_ = cmd.Execute()
	out := buf.String()
	assert.Contains(t, out, "No paused time entry to continue.")
}
