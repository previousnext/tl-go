package show

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/previousnext/tl-go/internal/model"
	"github.com/previousnext/tl-go/internal/service"
)

type mockTimerEntryStorage struct {
	entry  *model.TimerEntry
	retErr error
}

func (m *mockTimerEntryStorage) StartTimeEntry(issueKey string) error     { return nil }
func (m *mockTimerEntryStorage) PauseTimeEntry() error                    { return nil }
func (m *mockTimerEntryStorage) StopTimeEntry() (*model.TimeEntry, error) { return nil, nil }
func (m *mockTimerEntryStorage) GetTimerEntry() (*model.TimerEntry, error) {
	return m.entry, m.retErr
}
func (m *mockTimerEntryStorage) SaveTimerEntry(entry *model.TimerEntry) error {
	return nil
}
func (m *mockTimerEntryStorage) FindAllTimerEntries() ([]*model.TimerEntry, error) {
	return nil, nil
}

func TestShowCommand_InProgress(t *testing.T) {
	mock := &mockTimerEntryStorage{
		entry: &model.TimerEntry{
			IssueKey:  "PNX-123",
			StartTime: time.Now(),
			Paused:    false,
			Duration:  0,
		},
	}
	cmd := NewCommand(func() service.TimerEntryStorageInterface { return mock })
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
	mock := &mockTimerEntryStorage{
		entry: &model.TimerEntry{
			IssueKey:  "PNX-456",
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
	assert.Contains(t, out, "PNX-456")
	assert.Contains(t, out, "No")
	assert.Contains(t, out, "5m")
}

func TestShowCommand_NoEntry(t *testing.T) {
	mock := &mockTimerEntryStorage{entry: nil, retErr: nil}
	cmd := NewCommand(func() service.TimerEntryStorageInterface { return mock })
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})
	assert.NoError(t, cmd.Execute())
	out := buf.String()
	assert.Contains(t, out, "No current time entry in progress.")
}
