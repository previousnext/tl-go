package db

import (
	"fmt"

	"github.com/previousnext/tl-go/internal/model"
)

type TimerEntryStorageInterface interface {
	CreateTimerEntry(entry *model.TimerEntry) error
	DeleteTimerEntry(entry *model.TimerEntry) error
	FindLatestActiveTimerEntry() (*model.TimerEntry, error)
	FindLatestPausedTimerEntry() (*model.TimerEntry, error)
	GetTimerEntry() (*model.TimerEntry, error)
	SaveTimerEntry(entry *model.TimerEntry) error
	FindAllTimerEntries() ([]*model.TimerEntry, error)
}

func (r *Repository) CreateTimerEntry(entry *model.TimerEntry) error {
	db := r.openDB()
	return db.Create(entry).Error
}

func (r *Repository) DeleteTimerEntry(entry *model.TimerEntry) error {
	db := r.openDB()
	return db.Delete(entry).Error
}

func (r *Repository) FindLatestActiveTimerEntry() (*model.TimerEntry, error) {
	db := r.openDB()
	var entry model.TimerEntry
	if err := db.Where("paused = ?", false).Order("id desc").First(&entry).Error; err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *Repository) FindLatestPausedTimerEntry() (*model.TimerEntry, error) {
	db := r.openDB()
	var entry model.TimerEntry
	if err := db.Where("paused = ?", true).Order("id desc").First(&entry).Error; err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *Repository) GetTimerEntry() (*model.TimerEntry, error) {
	db := r.openDB()
	var entry model.TimerEntry
	if err := db.Order("id desc").First(&entry).Error; err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *Repository) SaveTimerEntry(entry *model.TimerEntry) error {
	db := r.openDB()
	return db.Save(entry).Error
}

func (r *Repository) FindAllTimerEntries() ([]*model.TimerEntry, error) {
	db := r.openDB()
	var entries []*model.TimerEntry
	if err := db.Find(&entries).Error; err != nil {
		return entries, fmt.Errorf("failed to retrieve time entries: %w", err)
	}
	return entries, nil
}
