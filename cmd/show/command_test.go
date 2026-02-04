package show

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/db/mocks"
	"github.com/previousnext/tl-go/internal/model"
)

func TestNewCommand_PrintsEntryDetails(t *testing.T) {
	mock := &mocks.MockRepository{
		FindTimeEntryFunc: func(id uint) (*model.TimeEntry, error) {
			return &model.TimeEntry{
				Model: gorm.Model{
					ID:        id,
					CreatedAt: time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2024, 6, 2, 12, 0, 0, 0, time.UTC),
				},
				IssueKey:    "PNX-42",
				Duration:    90, // minutes
				Description: "Worked on something",
			}, nil
		},
	}
	cmd := NewCommand(func() db.TimeEntriesInterface { return mock })

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"123"})

	err := cmd.Execute()
	assert.NoError(t, err)
	output := buf.String()
	fmt.Print(output)
	assert.Equal(t, `Time Entry ID:	123
Issue Key:	PNX-42
Duration:	1h30m
Description:	Worked on something
Created At:	01 Jun 24 10:00 UTC
Updated At:	02 Jun 24 12:00 UTC
`, output)
}

func TestNewCommand_NotFound_PrintsNoEntryMessage(t *testing.T) {
	mock := &mocks.MockRepository{
		FindTimeEntryFunc: func(id uint) (*model.TimeEntry, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	cmd := NewCommand(func() db.TimeEntriesInterface { return mock })

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"999"})

	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "No entry with ID 999")
}

func TestNewCommand_RepositoryError_ReturnsError(t *testing.T) {
	mock := &mocks.MockRepository{
		FindTimeEntryFunc: func(id uint) (*model.TimeEntry, error) {
			return nil, errors.New("db error")
		},
	}
	cmd := NewCommand(func() db.TimeEntriesInterface { return mock })

	cmd.SetArgs([]string{"1"})
	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
}

func TestNewCommand_InvalidID_ReturnsError(t *testing.T) {
	mock := &mocks.MockRepository{}
	cmd := NewCommand(func() db.TimeEntriesInterface { return mock })

	cmd.SetArgs([]string{"notanumber"})
	err := cmd.Execute()
	assert.Error(t, err)
}
