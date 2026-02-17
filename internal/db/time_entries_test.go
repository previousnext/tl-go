package db

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/previousnext/tl-go/internal/model"
)

func setupTestRepo(t *testing.T) *Repository {
	tmpFile, err := os.CreateTemp("", "test_*.db")
	assert.NoError(t, err)
	assert.NoError(t, tmpFile.Close())

	repo := NewRepository(tmpFile.Name())
	err = repo.AutoMigrate()
	assert.NoError(t, err)

	return repo
}

func cleanupTestRepo(_ *testing.T, repo *Repository) {
	_ = os.Remove(repo.dbPath)
}

func TestGetSummaryByCategory(t *testing.T) {
	repo := setupTestRepo(t)
	defer cleanupTestRepo(t, repo)

	category := &model.Category{
		Name: "Billable",
	}
	db := repo.openDB()
	assert.NoError(t, db.Create(category).Error)

	project := &model.Project{
		Key:        "TEST",
		Name:       "Test Project",
		CategoryID: &category.ID,
		Category:   category,
	}
	assert.NoError(t, db.Create(project).Error)

	issue := &model.Issue{
		Key:       "TEST-1",
		Summary:   "Test Issue",
		ProjectID: project.ID,
		Project:   *project,
	}
	assert.NoError(t, db.Create(issue).Error)

	start := time.Now().Add(-24 * time.Hour)
	end := time.Now().Add(24 * time.Hour)

	entry1 := &model.TimeEntry{
		Model: gorm.Model{
			CreatedAt: time.Now(),
		},
		IssueKey:    "TEST-1",
		Duration:    2 * time.Hour,
		Description: "Work 1",
	}
	assert.NoError(t, db.Create(entry1).Error)

	entry2 := &model.TimeEntry{
		Model: gorm.Model{
			CreatedAt: time.Now(),
		},
		IssueKey:    "TEST-1",
		Duration:    3 * time.Hour,
		Description: "Work 2",
	}
	assert.NoError(t, db.Create(entry2).Error)

	summaries, err := repo.GetSummaryByCategory(start, end)

	assert.NoError(t, err)
	assert.Len(t, summaries, 1)
	assert.Equal(t, "Billable", summaries[0].CategoryName)
	assert.Equal(t, 5*time.Hour, summaries[0].Duration)
	assert.InDelta(t, 100.0, summaries[0].Percentage, 0.01)
}

func TestGetSummaryByCategory_NoResults(t *testing.T) {
	repo := setupTestRepo(t)
	defer cleanupTestRepo(t, repo)

	start := time.Now().Add(-24 * time.Hour)
	end := time.Now().Add(24 * time.Hour)

	summaries, err := repo.GetSummaryByCategory(start, end)

	assert.NoError(t, err)
	assert.Len(t, summaries, 0)
}

func TestGetSummaryByCategory_OutsideDateRange(t *testing.T) {
	repo := setupTestRepo(t)
	defer cleanupTestRepo(t, repo)

	category := &model.Category{
		Name: "Billable",
	}
	db := repo.openDB()
	assert.NoError(t, db.Create(category).Error)

	project := &model.Project{
		Key:        "TEST",
		Name:       "Test Project",
		CategoryID: &category.ID,
		Category:   category,
	}
	assert.NoError(t, db.Create(project).Error)

	issue := &model.Issue{
		Key:       "TEST-1",
		Summary:   "Test Issue",
		ProjectID: project.ID,
		Project:   *project,
	}
	assert.NoError(t, db.Create(issue).Error)

	entry := &model.TimeEntry{
		Model: gorm.Model{
			CreatedAt: time.Now().Add(-48 * time.Hour),
		},
		IssueKey:    "TEST-1",
		Duration:    2 * time.Hour,
		Description: "Work",
	}
	assert.NoError(t, db.Create(entry).Error)

	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()

	summaries, err := repo.GetSummaryByCategory(start, end)

	assert.NoError(t, err)
	assert.Len(t, summaries, 0)
}
