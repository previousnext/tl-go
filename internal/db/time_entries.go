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
	GetSummaryByCategory(start time.Time, end time.Time) ([]CategorySummary, error)
}

type CategorySummary struct {
	CategoryName string
	Duration     time.Duration
	Percentage   float64
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
	if err := db.Preload("Issue.Project.Category").First(&entry, id).Error; err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *Repository) FindAllTimeEntries(date time.Time) ([]*model.TimeEntry, error) {
	start, end := getStartAndEndOfDay(date)
	db := r.openDB()
	var entries []*model.TimeEntry
	if err := db.Preload("Issue.Project.Category").Where("created_at BETWEEN ? AND ?", start, end).Find(&entries).Error; err != nil {
		return nil, err
	}
	return entries, nil
}

func (r *Repository) FindUnsentTimeEntries() ([]*model.TimeEntry, error) {
	db := r.openDB()
	var entries []*model.TimeEntry
	if err := db.Preload("Issue.Project.Category").Where("sent = ?", false).Find(&entries).Error; err != nil {
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

func (r *Repository) GetSummaryByCategory(start time.Time, end time.Time) ([]CategorySummary, error) {
	db := r.openDB()
	var summaries []struct {
		CategoryName string
		Duration     time.Duration
	}
	if err := db.Model(&model.TimeEntry{}).
		Select("COALESCE(categories.name, 'None') as category_name, SUM(duration) as duration").
		Joins("JOIN issues ON time_entries.issue_key = issues.key").
		Joins("JOIN projects ON issues.project_id = projects.id").
		Joins("LEFT JOIN categories ON projects.category_id = categories.id").
		Where("time_entries.created_at BETWEEN ? AND ?", start, end).
		Group("category_name").
		Scan(&summaries).Error; err != nil {
		return nil, err
	}

	var totalDuration time.Duration
	for _, summary := range summaries {
		totalDuration += summary.Duration
	}

	var results []CategorySummary
	for _, summary := range summaries {
		percent := 0.0
		if totalDuration > 0 {
			percent = float64(summary.Duration) / float64(totalDuration) * 100
		}
		results = append(results, CategorySummary{
			CategoryName: summary.CategoryName,
			Duration:     summary.Duration,
			Percentage:   percent,
		})
	}
	return results, nil
}
