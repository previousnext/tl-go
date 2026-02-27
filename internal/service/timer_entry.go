package service

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/model"
)

type TimerEntryStorageInterface interface {
	StartTimeEntry(issueKey string) error
	PauseTimeEntry() error
	StopTimeEntry() (*model.TimeEntry, error)
	GetTimerEntry() (*model.TimerEntry, error)
	SaveTimerEntry(entry *model.TimerEntry) error
	FindAllTimerEntries() ([]*model.TimerEntry, error)
}

type TimerEntryService struct {
	timerEntryStorage db.TimerEntryStorageInterface
	timeEntryStorage  db.TimeEntriesInterface
	now               func() time.Time
}

func NewTimerEntryService(timerEntryStorage db.TimerEntryStorageInterface, timeEntryStorage db.TimeEntriesInterface) *TimerEntryService {
	return &TimerEntryService{
		timerEntryStorage: timerEntryStorage,
		timeEntryStorage:  timeEntryStorage,
		now:               time.Now,
	}
}

func (s *TimerEntryService) StartTimeEntry(issueKey string) error {
	now := s.now()
	prev, err := s.timerEntryStorage.FindLatestActiveTimerEntry()
	if err == nil && prev != nil {
		lastActive := prev.LastActiveTime
		if lastActive.IsZero() {
			lastActive = prev.StartTime
		}
		prev.Duration += now.Sub(lastActive)
		prev.Paused = true
		prev.PauseTime = now
		prev.LastActiveTime = now
		if err := s.timerEntryStorage.SaveTimerEntry(prev); err != nil {
			return err
		}
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	entry := &model.TimerEntry{
		IssueKey:       issueKey,
		StartTime:      now,
		LastActiveTime: now,
		Paused:         false,
		Duration:       0,
	}
	return s.timerEntryStorage.CreateTimerEntry(entry)
}

func (s *TimerEntryService) PauseTimeEntry() error {
	now := s.now()
	entry, err := s.timerEntryStorage.FindLatestActiveTimerEntry()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("no active time entry to pause")
		}
		return err
	}
	lastActive := entry.LastActiveTime
	if lastActive.IsZero() {
		lastActive = entry.StartTime
	}
	entry.Duration += now.Sub(lastActive)
	entry.Paused = true
	entry.PauseTime = now
	entry.LastActiveTime = now
	return s.timerEntryStorage.SaveTimerEntry(entry)
}

func (s *TimerEntryService) StopTimeEntry() (*model.TimeEntry, error) {
	now := s.now()
	entry, err := s.timerEntryStorage.FindLatestActiveTimerEntry()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			entry, err = s.timerEntryStorage.FindLatestPausedTimerEntry()
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, errors.New("no current time entry to stop")
				}
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	dur := entry.Duration
	if !entry.Paused {
		lastActive := entry.LastActiveTime
		if lastActive.IsZero() {
			lastActive = entry.StartTime
		}
		dur += now.Sub(lastActive)
	}

	timeEntry := &model.TimeEntry{
		IssueKey: entry.IssueKey,
		Duration: dur,
		Sent:     false,
	}
	if err := s.timeEntryStorage.CreateTimeEntry(timeEntry); err != nil {
		return nil, err
	}
	if err := s.timerEntryStorage.DeleteTimerEntry(entry); err != nil {
		return nil, err
	}
	return timeEntry, nil
}

func (s *TimerEntryService) GetTimerEntry() (*model.TimerEntry, error) {
	return s.timerEntryStorage.GetTimerEntry()
}

func (s *TimerEntryService) SaveTimerEntry(entry *model.TimerEntry) error {
	return s.timerEntryStorage.SaveTimerEntry(entry)
}

func (s *TimerEntryService) FindAllTimerEntries() ([]*model.TimerEntry, error) {
	return s.timerEntryStorage.FindAllTimerEntries()
}
