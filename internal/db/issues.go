package db

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/previousnext/tl-go/internal/model"
)

type IssueStorageInterface interface {
	FindAllIssues() ([]*model.Issue, error)
	CreateIssue(issue *model.Issue) error
	DeleteIssueByKey(key string) error
	DeleteAllIssues() error
	FindIssueByKey(key string) (*model.Issue, error)
	FindRecentIssues(limit int) ([]*model.Issue, error)
}

func (r *Repository) FindIssueByKey(key string) (*model.Issue, error) {
	db := r.openDB()
	var issue model.Issue
	if err := db.Preload("Project").Where("key = ?", key).First(&issue).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &issue, nil
}

func (r *Repository) FindAllIssues() ([]*model.Issue, error) {
	db := r.openDB()
	var issues []*model.Issue
	if err := db.Preload("Project").Find(&issues).Error; err != nil {
		return nil, err
	}
	return issues, nil
}

func (r *Repository) CreateIssue(issue *model.Issue) error {
	db := r.openDB()
	if err := db.Create(&issue).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteIssueByKey(key string) error {
	db := r.openDB()
	if err := db.Where("key = ?", key).Delete(&model.Issue{}).Error; err != nil {
		return fmt.Errorf("error deleting issue with key %s: %w", key, err)
	}
	return nil
}

func (r *Repository) DeleteAllIssues() error {
	db := r.openDB()
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Issue{}).Error; err != nil {
		return fmt.Errorf("error deleting all issues: %w", err)
	}
	return nil
}

func (r *Repository) FindRecentIssues(limit int) ([]*model.Issue, error) {
	db := r.openDB()
	var issues []*model.Issue
	if err := db.Preload("Project").
		Joins("JOIN time_entries ON time_entries.issue_key = issues.key").
		Group("issues.key").
		Order("MAX(time_entries.created_at) DESC").
		Limit(limit).
		Find(&issues).Error; err != nil {
		return nil, err
	}
	return issues, nil
}
