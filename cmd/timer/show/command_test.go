package show

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/model"
	svc "github.com/previousnext/tl-go/internal/service"
)

type mockTimerEntryStorage struct {
	entry  *model.TimerEntry
	retErr error
}

func (m *mockTimerEntryStorage) CreateTimerEntry(entry *model.TimerEntry) error { return nil }
func (m *mockTimerEntryStorage) DeleteTimerEntry(entry *model.TimerEntry) error { return nil }
func (m *mockTimerEntryStorage) FindLatestActiveTimerEntry() (*model.TimerEntry, error) {
	return nil, nil
}
func (m *mockTimerEntryStorage) FindLatestPausedTimerEntry() (*model.TimerEntry, error) {
	return nil, nil
}
func (m *mockTimerEntryStorage) GetTimerEntry() (*model.TimerEntry, error)         { return m.entry, m.retErr }
func (m *mockTimerEntryStorage) SaveTimerEntry(entry *model.TimerEntry) error      { return nil }
func (m *mockTimerEntryStorage) FindAllTimerEntries() ([]*model.TimerEntry, error) { return nil, nil }

type mockTimeEntriesStorage struct{}

func (m *mockTimeEntriesStorage) CreateTimeEntry(entry *model.TimeEntry) error    { return nil }
func (m *mockTimeEntriesStorage) FindTimeEntry(id uint) (*model.TimeEntry, error) { return nil, nil }
func (m *mockTimeEntriesStorage) FindAllTimeEntries(date time.Time) ([]*model.TimeEntry, error) {
	return nil, nil
}
func (m *mockTimeEntriesStorage) FindUnsentTimeEntries() ([]*model.TimeEntry, error) { return nil, nil }
func (m *mockTimeEntriesStorage) FindUniqueIssueKeys() ([]string, error)             { return nil, nil }
func (m *mockTimeEntriesStorage) UpdateTimeEntry(entry *model.TimeEntry) error       { return nil }
func (m *mockTimeEntriesStorage) DeleteTimeEntry(id uint) error                      { return nil }
func (m *mockTimeEntriesStorage) GetSummaryByCategory(start time.Time, end time.Time) ([]db.CategorySummary, error) {
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
	timerService := svc.NewTimerEntryService(mock, &mockTimeEntriesStorage{})
	cmd := NewCommand(func() svc.TimerEntryServiceInterface { return timerService })
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
	timerService := svc.NewTimerEntryService(mock, &mockTimeEntriesStorage{})
	cmd := NewCommand(func() svc.TimerEntryServiceInterface { return timerService })
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
	timerService := svc.NewTimerEntryService(mock, &mockTimeEntriesStorage{})
	cmd := NewCommand(func() svc.TimerEntryServiceInterface { return timerService })
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})
	assert.NoError(t, cmd.Execute())
	out := buf.String()
	assert.Contains(t, out, "No timer entry in progress.")
}
