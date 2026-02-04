package db

import "github.com/previousnext/tl-go/internal/model"

type IssuesInterface interface {
	FindAllIssues() ([]*model.Issue, error)
	CreateIssue(issue *model.Issue) error
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
