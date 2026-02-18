package db

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/previousnext/tl-go/internal/model"
)

func setupRepo(t *testing.T) *Repository {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.sqlite")
	repo := NewRepository(dbPath)
	db := repo.openDB()
	assert.NoError(t, db.AutoMigrate(&model.CurrentTimeEntry{}, &model.TimeEntry{}))
	return repo
}

func TestStartTimeEntry(t *testing.T) {
	repo := setupRepo(t)
	// Start first entry
	assert.NoError(t, repo.StartTimeEntry("ISSUE-1"))
	// Start second entry, should stop previous
	time.Sleep(10 * time.Millisecond)
	assert.NoError(t, repo.StartTimeEntry("ISSUE-2"))
	// Check that only one current entry exists
	db := repo.openDB()
	var current model.CurrentTimeEntry
	err := db.Where("paused = ?", false).First(&current).Error
	assert.NoError(t, err)
	assert.Equal(t, "ISSUE-2", current.IssueKey)
	// Check that a completed entry exists for ISSUE-1
	var entries []model.TimeEntry
	db.Find(&entries)
	assert.Len(t, entries, 1)
	assert.Equal(t, "ISSUE-1", entries[0].IssueKey)
	assert.True(t, entries[0].Duration > 0)
}

func TestPauseTimeEntry(t *testing.T) {
	repo := setupRepo(t)
	assert.NoError(t, repo.StartTimeEntry("ISSUE-1"))
	assert.NoError(t, repo.PauseTimeEntry())
	// Check paused state
	db := repo.openDB()
	var current model.CurrentTimeEntry
	assert.NoError(t, db.First(&current).Error)
	assert.True(t, current.Paused)
}

func TestPauseTimeEntry_NoActive(t *testing.T) {
	repo := setupRepo(t)
	err := repo.PauseTimeEntry()
	assert.Error(t, err)
	assert.Equal(t, "no active time entry to pause", err.Error())
}

func TestStopTimeEntry(t *testing.T) {
	repo := setupRepo(t)
	assert.NoError(t, repo.StartTimeEntry("ISSUE-1"))
	time.Sleep(10 * time.Millisecond)
	entry, err := repo.StopTimeEntry()
	assert.NoError(t, err)
	assert.Equal(t, "ISSUE-1", entry.IssueKey)
	assert.True(t, entry.Duration > 0)
	// Should be no current entry
	db := repo.openDB()
	var current model.CurrentTimeEntry
	assert.Error(t, db.First(&current).Error)
}

func TestStopTimeEntry_NoActive(t *testing.T) {
	repo := setupRepo(t)
	entry, err := repo.StopTimeEntry()
	assert.Nil(t, entry)
	assert.Error(t, err)
	assert.Equal(t, "no active time entry to stop", err.Error())
}

func TestPauseContinueStopSequence(t *testing.T) {
	repo := setupRepo(t)
	// Start at t0
	t0 := time.Now()
	assert.NoError(t, repo.StartTimeEntry("PNX-123"))
	// Simulate 5 minutes of work
	db := repo.openDB()
	var current model.CurrentTimeEntry
	assert.NoError(t, db.First(&current).Error)
	current.StartTime = t0
	assert.NoError(t, db.Save(&current).Error)
	// Pause after 5 minutes
	t1 := t0.Add(5 * time.Minute)
	current = model.CurrentTimeEntry{}
	assert.NoError(t, db.First(&current).Error)
	current.PauseTime = t1
	current.Paused = true
	assert.NoError(t, db.Save(&current).Error)
	// Continue after 30 minutes
	t2 := t1.Add(30 * time.Minute)
	current = model.CurrentTimeEntry{}
	assert.NoError(t, db.First(&current).Error)
	// Simulate continue logic
	current.Duration += current.PauseTime.Sub(current.StartTime)
	current.Paused = false
	current.StartTime = t2
	current.PauseTime = time.Time{}
	assert.NoError(t, db.Save(&current).Error)
	// Stop after 10 more minutes
	t3 := t2.Add(10 * time.Minute)
	current = model.CurrentTimeEntry{}
	assert.NoError(t, db.First(&current).Error)
	dur := t3.Sub(current.StartTime) + current.Duration
	// Simulate stop logic
	timeEntry := &model.TimeEntry{
		IssueKey: current.IssueKey,
		Duration: dur,
		Sent:     false,
	}
	assert.NoError(t, db.Create(timeEntry).Error)
	assert.NoError(t, db.Delete(&current).Error)
	// Validate the total duration is 15 minutes
	var entries []model.TimeEntry
	db.Find(&entries)
	assert.Len(t, entries, 1)
	assert.Equal(t, "PNX-123", entries[0].IssueKey)
	assert.InDelta(t, 15*60, entries[0].Duration.Seconds(), 1, "Duration should be 15 minutes")
}
