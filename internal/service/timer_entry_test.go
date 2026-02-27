package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/model"
)

type mockTimerEntryStorage struct {
	entry       *model.TimerEntry
	pausedEntry *model.TimerEntry
	deleted     bool
}

func (m *mockTimerEntryStorage) CreateTimerEntry(entry *model.TimerEntry) error {
	m.entry = entry
	return nil
}

func (m *mockTimerEntryStorage) DeleteTimerEntry(entry *model.TimerEntry) error {
	m.deleted = true
	return nil
}

func (m *mockTimerEntryStorage) FindLatestActiveTimerEntry() (*model.TimerEntry, error) {
	if m.entry == nil || m.entry.Paused {
		return nil, gorm.ErrRecordNotFound
	}
	return m.entry, nil
}

func (m *mockTimerEntryStorage) FindLatestPausedTimerEntry() (*model.TimerEntry, error) {
	if m.entry == nil || !m.entry.Paused {
		return nil, gorm.ErrRecordNotFound
	}
	return m.entry, nil
}

func (m *mockTimerEntryStorage) GetTimerEntry() (*model.TimerEntry, error) {
	if m.entry == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return m.entry, nil
}

func (m *mockTimerEntryStorage) SaveTimerEntry(entry *model.TimerEntry) error {
	m.pausedEntry = &model.TimerEntry{}
	*m.pausedEntry = *entry
	m.entry = entry
	return nil
}

func (m *mockTimerEntryStorage) FindAllTimerEntries() ([]*model.TimerEntry, error) {
	if m.entry == nil {
		return nil, nil
	}
	return []*model.TimerEntry{m.entry}, nil
}

func (m *mockTimerEntryStorage) FindTimerEntryByID(id uint) (*model.TimerEntry, error) {
	if m.entry == nil || m.entry.ID != id {
		return nil, gorm.ErrRecordNotFound
	}
	return m.entry, nil
}

type mockTimeEntriesStorage struct {
	created *model.TimeEntry
}

func (m *mockTimeEntriesStorage) CreateTimeEntry(entry *model.TimeEntry) error {
	m.created = entry
	return nil
}

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

func TestTimerEntryService_TimerWorkflow(t *testing.T) {
	start := time.Date(2026, 2, 27, 9, 0, 0, 0, time.Local)
	pause := start.Add(5 * time.Minute)
	resume := pause.Add(30 * time.Minute)
	stop := resume.Add(10 * time.Minute)
	nowTimes := []time.Time{start, pause, resume, stop}
	idx := 0

	mockTimer := &mockTimerEntryStorage{}
	mockTimeEntries := &mockTimeEntriesStorage{}
	service := NewTimerEntryService(mockTimer, mockTimeEntries)
	service.now = func() time.Time {
		current := nowTimes[idx]
		idx++
		return current
	}

	assert.NoError(t, service.StartTimeEntry("PNX-123", nil))
	assert.NoError(t, service.PauseTimeEntry())
	assert.NoError(t, service.ResumeTimerEntry(nil))
	timeEntry, err := service.StopTimeEntry(nil)
	assert.NoError(t, err)
	assert.NotNil(t, timeEntry)

	assert.Equal(t, "PNX-123", timeEntry.IssueKey)
	assert.Equal(t, 15*time.Minute, timeEntry.Duration)
	assert.True(t, mockTimer.deleted)
	assert.NotNil(t, mockTimeEntries.created)
}

func TestTimerEntryService_StopRoundsToQuarterHour(t *testing.T) {
	start := time.Date(2026, 2, 27, 9, 0, 0, 0, time.Local)
	stop := start.Add(16 * time.Minute)
	nowTimes := []time.Time{start, stop}
	idx := 0

	mockTimer := &mockTimerEntryStorage{}
	mockTimeEntries := &mockTimeEntriesStorage{}
	service := NewTimerEntryService(mockTimer, mockTimeEntries)
	service.now = func() time.Time {
		current := nowTimes[idx]
		idx++
		return current
	}

	assert.NoError(t, service.StartTimeEntry("PNX-456", nil))
	entry, err := service.StopTimeEntry(nil)
	assert.NoError(t, err)
	assert.NotNil(t, entry)
	assert.NotNil(t, mockTimeEntries.created)
	assert.Equal(t, 30*time.Minute, mockTimeEntries.created.Duration)
}

func TestTimerEntryService_StopByID(t *testing.T) {
	start := time.Date(2026, 2, 27, 9, 0, 0, 0, time.Local)
	stop := start.Add(7 * time.Minute)
	nowTimes := []time.Time{stop}
	idx := 0

	mockTimer := &mockTimerEntryStorage{
		entry: &model.TimerEntry{
			ID:             99,
			IssueKey:       "PNX-777",
			StartTime:      start,
			LastActiveTime: start,
			Paused:         false,
			Duration:       0,
		},
	}
	mockTimeEntries := &mockTimeEntriesStorage{}
	service := NewTimerEntryService(mockTimer, mockTimeEntries)
	service.now = func() time.Time {
		current := nowTimes[idx]
		idx++
		return current
	}

	entry, err := service.StopTimeEntry(func() *uint { v := uint(99); return &v }())
	assert.NoError(t, err)
	assert.NotNil(t, entry)
	assert.Equal(t, "PNX-777", entry.IssueKey)
	assert.NotNil(t, mockTimeEntries.created)
}

func TestTimerEntryService_StopByID_NotFound(t *testing.T) {
	mockTimer := &mockTimerEntryStorage{entry: nil}
	mockTimeEntries := &mockTimeEntriesStorage{}
	service := NewTimerEntryService(mockTimer, mockTimeEntries)

	entry, err := service.StopTimeEntry(func() *uint { v := uint(123); return &v }())
	assert.Nil(t, entry)
	assert.Error(t, err)
	assert.Equal(t, "timer entry not found", err.Error())
}

func TestTimerEntryService_OnlyOneActiveTimer(t *testing.T) {
	start1 := time.Date(2026, 2, 27, 9, 0, 0, 0, time.Local)
	start2 := start1.Add(2 * time.Minute)
	nowTimes := []time.Time{start1, start2}
	idx := 0

	mockTimer := &mockTimerEntryStorage{}
	mockTimeEntries := &mockTimeEntriesStorage{}
	service := NewTimerEntryService(mockTimer, mockTimeEntries)
	service.now = func() time.Time {
		current := nowTimes[idx]
		idx++
		return current
	}

	assert.NoError(t, service.StartTimeEntry("PNX-111", nil))
	assert.False(t, mockTimer.entry.Paused)
	assert.NoError(t, service.StartTimeEntry("PNX-222", nil))
	assert.NotNil(t, mockTimer.pausedEntry)
	assert.True(t, mockTimer.pausedEntry.Paused)
	assert.Equal(t, "PNX-111", mockTimer.pausedEntry.IssueKey)
	assert.Equal(t, "PNX-222", mockTimer.entry.IssueKey)
	assert.False(t, mockTimer.entry.Paused)
}
