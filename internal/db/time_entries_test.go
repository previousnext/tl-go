package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/previousnext/tl-go/internal/model"
)

func setupTestRepo(t *testing.T) *Repository {
	dir := t.TempDir()
	dbPath := dir + "/test.db"

	repo := NewRepository(dbPath)
	err := repo.AutoMigrate()
	assert.NoError(t, err)

	return repo
}

// seedTestIssue creates the standard Category -> Project -> Issue chain used by
// the summary tests and returns the created issue.
func seedTestIssue(t *testing.T, db *gorm.DB) *model.Issue {
	t.Helper()

	category := &model.Category{Name: "Billable"}
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

	return issue
}

// createTestEntry persists a TimeEntry, defaulting CreatedAt to now when unset.
func createTestEntry(t *testing.T, db *gorm.DB, e model.TimeEntry) {
	t.Helper()

	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now()
	}
	assert.NoError(t, db.Create(&e).Error)
}

func TestGetSummaryByCategory(t *testing.T) {
	repo := setupTestRepo(t)
	db := repo.openDB()

	issue := seedTestIssue(t, db)

	start := time.Now().Add(-24 * time.Hour)
	end := time.Now().Add(24 * time.Hour)

	createTestEntry(t, db, model.TimeEntry{
		IssueKey:    "TEST-1",
		IssueID:     issue.ID,
		Duration:    2 * time.Hour,
		Description: "Work 1",
	})
	createTestEntry(t, db, model.TimeEntry{
		IssueKey:    "TEST-1",
		IssueID:     issue.ID,
		Duration:    3 * time.Hour,
		Description: "Work 2",
	})

	summaries, err := repo.GetSummaryByCategory(start, end)

	assert.NoError(t, err)
	assert.Len(t, summaries, 1)
	assert.Equal(t, "Billable", summaries[0].CategoryName)
	assert.Equal(t, 5*time.Hour, summaries[0].Duration)
	assert.InDelta(t, 100.0, summaries[0].Percentage, 0.01)
}

func TestGetSummaryByCategory_NoResults(t *testing.T) {
	repo := setupTestRepo(t)

	start := time.Now().Add(-24 * time.Hour)
	end := time.Now().Add(24 * time.Hour)

	summaries, err := repo.GetSummaryByCategory(start, end)

	assert.NoError(t, err)
	assert.Len(t, summaries, 0)
}

func TestGetSummaryByCategory_OutsideDateRange(t *testing.T) {
	repo := setupTestRepo(t)
	db := repo.openDB()

	issue := seedTestIssue(t, db)

	createTestEntry(t, db, model.TimeEntry{
		Model:       gorm.Model{CreatedAt: time.Now().Add(-48 * time.Hour)},
		IssueKey:    "TEST-1",
		IssueID:     issue.ID,
		Duration:    2 * time.Hour,
		Description: "Work",
	})

	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()

	summaries, err := repo.GetSummaryByCategory(start, end)

	assert.NoError(t, err)
	assert.Len(t, summaries, 0)
}

func TestGetSummaryByCategory_KeyCaseMismatch(t *testing.T) {
	repo := setupTestRepo(t)
	db := repo.openDB()

	issue := seedTestIssue(t, db)

	start := time.Now().Add(-24 * time.Hour)
	end := time.Now().Add(24 * time.Hour)

	// The stored issue_key casing differs from issues.key, but issue_id links
	// the entry correctly. The summary should still resolve the category.
	createTestEntry(t, db, model.TimeEntry{
		IssueKey:    "test-1",
		IssueID:     issue.ID,
		Duration:    2 * time.Hour,
		Description: "Work",
	})

	summaries, err := repo.GetSummaryByCategory(start, end)

	assert.NoError(t, err)
	assert.Len(t, summaries, 1)
	assert.Equal(t, "Billable", summaries[0].CategoryName)
	assert.Equal(t, 2*time.Hour, summaries[0].Duration)
}

func TestGetSummaryByCategory_OrphanEntryNone(t *testing.T) {
	repo := setupTestRepo(t)
	db := repo.openDB()

	start := time.Now().Add(-24 * time.Hour)
	end := time.Now().Add(24 * time.Hour)

	// An entry with no linked issue should roll up under the "None" category
	// rather than being dropped entirely.
	createTestEntry(t, db, model.TimeEntry{
		IssueKey:    "GONE-1",
		Duration:    90 * time.Minute,
		Description: "Orphan work",
	})

	summaries, err := repo.GetSummaryByCategory(start, end)

	assert.NoError(t, err)
	assert.Len(t, summaries, 1)
	assert.Equal(t, "None", summaries[0].CategoryName)
	assert.Equal(t, 90*time.Minute, summaries[0].Duration)
}
