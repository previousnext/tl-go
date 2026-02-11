package db

import (
	"time"

	"github.com/previousnext/tl-go/internal/model"
)

type TimeEntriesInterface interface {
	CreateTimeEntry(entry *model.TimeEntry) error
	FindTimeEntry(id uint) (*model.TimeEntry, error)
	FindAllTimeEntries(date time.Time) ([]*model.TimeEntry, error)
	FindUnsentTimeEntries() ([]*model.TimeEntry, error)
	FindUniqueIssueKeys() ([]string, error)
	UpdateTimeEntry(entry *model.TimeEntry) error
	DeleteTimeEntry(id uint) error
}

func (r *Repository) CreateTimeEntry(entry *model.TimeEntry) error {
	db := r.openDB()
	if err := db.Create(&entry).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) FindTimeEntry(id uint) (*model.TimeEntry, error) {
	db := r.openDB()
	var entry model.TimeEntry
	if err := db.Preload("Issue.Project").First(&entry, id).Error; err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *Repository) FindAllTimeEntries(date time.Time) ([]*model.TimeEntry, error) {
	start, end := getStartAndEndOfDay(date)
	db := r.openDB()
	var entries []*model.TimeEntry
	if err := db.Preload("Issue.Project").Where("created_at BETWEEN ? AND ?", start, end).Find(&entries).Error; err != nil {
		return nil, err
	}
	return entries, nil
}

func (r *Repository) FindUnsentTimeEntries() ([]*model.TimeEntry, error) {
	db := r.openDB()
	var entries []*model.TimeEntry
	if err := db.Where("sent = ?", false).Find(&entries).Error; err != nil {
		return nil, err
	}
	return entries, nil
}

func (r *Repository) FindUniqueIssueKeys() ([]string, error) {
	db := r.openDB()
	var issueKeys []string
	if err := db.Model(&model.TimeEntry{}).Distinct().Pluck("issue_key", &issueKeys).Error; err != nil {
		return nil, err
	}
	return issueKeys, nil
}

func (r *Repository) UpdateTimeEntry(entry *model.TimeEntry) error {
	db := r.openDB()
	return db.Save(entry).Error
}

func (r *Repository) DeleteTimeEntry(id uint) error {
	db := r.openDB()
	return db.Delete(&model.TimeEntry{}, id).Error
}
