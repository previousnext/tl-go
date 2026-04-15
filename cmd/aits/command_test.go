package aits

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/previousnext/tl-go/internal/db"
	dbmocks "github.com/previousnext/tl-go/internal/db/mocks"
	"github.com/previousnext/tl-go/internal/model"
)

func TestAits(t *testing.T) {
	entry := &model.TimeEntry{
		Model:    gorm.Model{ID: 42},
		IssueKey: "PNX-123",
		Duration: 2 * time.Hour,
		Issue:    &model.Issue{Project: model.Project{}},
	}
	var updated *model.TimeEntry
	mock := &dbmocks.MockRepository{
		FindTimeEntryFunc: func(id uint) (*model.TimeEntry, error) {
			assert.Equal(t, uint(42), id)
			return entry, nil
		},
		UpdateTimeEntryFunc: func(e *model.TimeEntry) error {
			updated = e
			return nil
		},
	}

	cmd := NewCommand(func() db.TimeEntriesInterface { return mock })

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"42", "1h"})

	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Set AI time saved to 1h on time entry ID 42")
	assert.NotNil(t, updated)
	assert.Equal(t, time.Hour, updated.AISavedDuration)
}

func TestAits_InvalidID_ReturnsError(t *testing.T) {
	mock := &dbmocks.MockRepository{}
	cmd := NewCommand(func() db.TimeEntriesInterface { return mock })

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"notanid", "1h"})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid time entry ID")
}

func TestAits_InvalidDuration_ReturnsError(t *testing.T) {
	mock := &dbmocks.MockRepository{}
	cmd := NewCommand(func() db.TimeEntriesInterface { return mock })

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"42", "notaduration"})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid duration")
}

func TestAits_EntryNotFound(t *testing.T) {
	mock := &dbmocks.MockRepository{
		FindTimeEntryFunc: func(id uint) (*model.TimeEntry, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	cmd := NewCommand(func() db.TimeEntriesInterface { return mock })

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"99", "1h"})

	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "No entry with ID 99")
}
