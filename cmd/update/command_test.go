package update

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/db/mocks"
	"github.com/previousnext/tl-go/internal/model"
)

func TestNewCommand_UpdatesEntryAndPrintsMessage(t *testing.T) {
	mock := &mocks.MockRepository{
		FindTimeEntryFunc: func(id uint) (*model.TimeEntry, error) {
			return &model.TimeEntry{
				Model:       gorm.Model{ID: id},
				Duration:    time.Hour,
				Description: "Old description",
			}, nil
		},
		UpdateTimeEntryFunc: func(entry *model.TimeEntry) error {
			assert.Equal(t, uint(123), entry.ID)
			assert.Equal(t, 3*time.Hour, entry.Duration)
			assert.Equal(t, "Updated description", entry.Description)
			return nil
		},
	}
	cmd := NewCommand(func() db.TimeEntriesInterface { return mock })

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"123", "3h", "Updated description"})

	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Updated time entry with ID 123")
}

func TestNewCommand_EntryNotFound_PrintsNoEntryMessage(t *testing.T) {
	mock := &mocks.MockRepository{
		FindTimeEntryFunc: func(id uint) (*model.TimeEntry, error) {
			return nil, nil
		},
	}
	cmd := NewCommand(func() db.TimeEntriesInterface { return mock })

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"999", "1h"})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "entry not found with ID: 999")
}

func TestNewCommand_InvalidID_ReturnsError(t *testing.T) {
	mock := &mocks.MockRepository{}
	cmd := NewCommand(func() db.TimeEntriesInterface { return mock })

	cmd.SetArgs([]string{"notanumber", "1h"})
	err := cmd.Execute()
	assert.Error(t, err)
}

func TestNewCommand_InvalidDuration_ReturnsError(t *testing.T) {
	mock := &mocks.MockRepository{
		FindTimeEntryFunc: func(id uint) (*model.TimeEntry, error) {
			return &model.TimeEntry{
				Model: gorm.Model{ID: id},
			}, nil
		},
	}
	cmd := NewCommand(func() db.TimeEntriesInterface { return mock })

	cmd.SetArgs([]string{"1", "notaduration"})
	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid duration")
}

func TestNewCommand_RepositoryError_ReturnsError(t *testing.T) {
	mock := &mocks.MockRepository{
		FindTimeEntryFunc: func(id uint) (*model.TimeEntry, error) {
			return &model.TimeEntry{
				Model: gorm.Model{ID: id},
			}, nil
		},
		UpdateTimeEntryFunc: func(entry *model.TimeEntry) error {
			return errors.New("update failed")
		},
	}
	cmd := NewCommand(func() db.TimeEntriesInterface { return mock })

	cmd.SetArgs([]string{"1", "1h"})
	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update failed")
}
