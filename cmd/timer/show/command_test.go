package show

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/model"
)

type mockCurrentTimeEntryStorage struct {
	entry  *model.CurrentTimeEntry
	retErr error
}

func (m *mockCurrentTimeEntryStorage) StartTimeEntry(issueKey string) error     { return nil }
func (m *mockCurrentTimeEntryStorage) PauseTimeEntry() error                    { return nil }
func (m *mockCurrentTimeEntryStorage) StopTimeEntry() (*model.TimeEntry, error) { return nil, nil }
func (m *mockCurrentTimeEntryStorage) GetCurrentTimeEntry() (*model.CurrentTimeEntry, error) {
	return m.entry, m.retErr
}
func (m *mockCurrentTimeEntryStorage) SaveCurrentTimeEntry(entry *model.CurrentTimeEntry) error {
	return nil
}

func TestShowCommand_InProgress(t *testing.T) {
	mock := &mockCurrentTimeEntryStorage{
		entry: &model.CurrentTimeEntry{
			IssueKey:  "PNX-123",
			StartTime: time.Now(),
			Paused:    false,
			Duration:  0,
		},
	}
	cmd := NewCommand(func() db.CurrentTimeEntryStorageInterface { return mock })
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})
	assert.NoError(t, cmd.Execute())
	out := buf.String()
	assert.Contains(t, out, "PNX-123")
	assert.Contains(t, out, "Yes")
	assert.Contains(t, out, "s")
}

func TestShowCommand_Paused(t *testing.T) {
	start := time.Now()
	pause := start.Add(5 * time.Minute)
	mock := &mockCurrentTimeEntryStorage{
		entry: &model.CurrentTimeEntry{
			IssueKey:  "PNX-456",
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
	assert.Contains(t, out, "PNX-456")
	assert.Contains(t, out, "No")
	assert.Contains(t, out, "5m")
}

func TestShowCommand_NoEntry(t *testing.T) {
	mock := &mockCurrentTimeEntryStorage{entry: nil, retErr: nil}
	cmd := NewCommand(func() db.CurrentTimeEntryStorageInterface { return mock })
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})
	assert.NoError(t, cmd.Execute())
	out := buf.String()
	assert.Contains(t, out, "No current time entry in progress.")
}
