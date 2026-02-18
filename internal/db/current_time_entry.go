package db

import (
	"errors"
	"time"

	"github.com/previousnext/tl-go/internal/model"
)

type CurrentTimeEntryStorageInterface interface {
	StartTimeEntry(issueKey string) error
	PauseTimeEntry() error
	StopTimeEntry() (*model.TimeEntry, error)
	GetCurrentTimeEntry() (*model.CurrentTimeEntry, error)
	SaveCurrentTimeEntry(entry *model.CurrentTimeEntry) error
}

func (r *Repository) StartTimeEntry(issueKey string) error {
	db := r.openDB()
	// Stop any existing session first
	var prev model.CurrentTimeEntry
	if err := db.Where("paused = ?", false).First(&prev).Error; err == nil {
		dur := time.Since(prev.StartTime) + prev.Duration
		timeEntry := &model.TimeEntry{
			IssueKey: prev.IssueKey,
			Duration: dur,
			Sent:     false,
		}
		if err := db.Create(timeEntry).Error; err != nil {
			return err
		}
		if err := db.Delete(&prev).Error; err != nil {
			return err
		}
	}
	// Now start new session
	entry := &model.CurrentTimeEntry{
		IssueKey:  issueKey,
		StartTime: time.Now(),
		Paused:    false,
		Duration:  0,
	}
	return db.Create(entry).Error
}

func (r *Repository) PauseTimeEntry() error {
	db := r.openDB()
	var entry model.CurrentTimeEntry
	if err := db.Where("paused = ?", false).First(&entry).Error; err != nil {
		return errors.New("no active time entry to pause")
	}
	entry.Paused = true
	entry.PauseTime = time.Now()
	return db.Save(&entry).Error
}

func (r *Repository) StopTimeEntry() (*model.TimeEntry, error) {
	db := r.openDB()
	var entry model.CurrentTimeEntry
	if err := db.Where("paused = ?", false).First(&entry).Error; err != nil {
		return nil, errors.New("no active time entry to stop")
	}
	dur := time.Since(entry.StartTime) + entry.Duration
	timeEntry := &model.TimeEntry{
		IssueKey: entry.IssueKey,
		Duration: dur,
		Sent:     false,
	}
	if err := db.Create(timeEntry).Error; err != nil {
		return nil, err
	}
	_ = db.Delete(&entry)
	return timeEntry, nil
}

func (r *Repository) GetCurrentTimeEntry() (*model.CurrentTimeEntry, error) {
	db := r.openDB()
	var entry model.CurrentTimeEntry
	if err := db.First(&entry).Error; err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *Repository) SaveCurrentTimeEntry(entry *model.CurrentTimeEntry) error {
	db := r.openDB()
	return db.Save(entry).Error
}
