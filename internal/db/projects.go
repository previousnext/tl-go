package db

import "github.com/previousnext/tl-go/internal/model"

type ProjectsInterface interface {
	FindAllProjects() ([]*model.Project, error)
	CreateProject(project *model.Project) error
}

func (r *Repository) FindAllProjects() ([]*model.Project, error) {
	db := r.openDB()
	var projects []*model.Project
	if err := db.Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *Repository) CreateProject(project *model.Project) error {
	db := r.openDB()
	if err := db.Create(&project).Error; err != nil {
		return err
	}
	return nil
}
