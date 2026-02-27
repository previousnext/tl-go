package continuecmd

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/previousnext/tl-go/internal/model"
	"github.com/previousnext/tl-go/internal/service"
)

type mockTimerEntryStorage struct {
	entry   *model.TimerEntry
	retErr  error
	updated *model.TimerEntry
}

func (m *mockTimerEntryStorage) StartTimeEntry(issueKey string) error     { return nil }
func (m *mockTimerEntryStorage) PauseTimeEntry() error                    { return nil }
func (m *mockTimerEntryStorage) StopTimeEntry() (*model.TimeEntry, error) { return nil, nil }
func (m *mockTimerEntryStorage) GetTimerEntry() (*model.TimerEntry, error) {
	return m.entry, m.retErr
}
func (m *mockTimerEntryStorage) SaveTimerEntry(entry *model.TimerEntry) error {
	m.updated = &model.TimerEntry{}
	*m.updated = *entry
	return nil
}
func (m *mockTimerEntryStorage) FindAllTimerEntries() ([]*model.TimerEntry, error) {
	return nil, nil
}

func TestContinueCommand_PausedEntry(t *testing.T) {
	start := time.Date(2026, 2, 18, 10, 0, 0, 0, time.Local)
	pause := start.Add(5 * time.Minute)
	mock := &mockTimerEntryStorage{
		entry: &model.TimerEntry{
			IssueKey:  "PNX-123",
			StartTime: start,
			Paused:    true,
			PauseTime: pause,
			Duration:  0,
		},
	}
	cmd := NewCommand(func() service.TimerEntryStorageInterface { return mock })
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
	mock := &mockTimerEntryStorage{
		entry: &model.TimerEntry{
			IssueKey: "PNX-123",
			Paused:   false,
		},
	}
	cmd := NewCommand(func() service.TimerEntryStorageInterface { return mock })
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})
	assert.NoError(t, cmd.Execute())
	out := buf.String()
	assert.Contains(t, out, "Current time entry is not paused.")
}

func TestContinueCommand_NoEntry(t *testing.T) {
	mock := &mockTimerEntryStorage{entry: nil, retErr: nil}
	cmd := NewCommand(func() service.TimerEntryStorageInterface { return mock })
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})
	_ = cmd.Execute()
	out := buf.String()
	assert.Contains(t, out, "No paused time entry to continue.")
}
