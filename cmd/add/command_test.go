package add

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/previousnext/tl-go/internal/db"
	dbmocks "github.com/previousnext/tl-go/internal/db/mocks"
	"github.com/previousnext/tl-go/internal/service"
	servicemocks "github.com/previousnext/tl-go/internal/service/mocks"
)

func TestAdd(t *testing.T) {
	cmd := NewCommand(
		func() db.TimeEntriesInterface { return &dbmocks.MockRepository{} },
		func() service.SyncInterface { return &servicemocks.MockSync{} },
		func() db.IssueStorageInterface { return &dbmocks.MockRepository{} },
	)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"PNX-123", "2h", "Worked on feature X"})

	err := cmd.Execute()
	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "Added time entry: ID=42\n")
}

func TestAdd_InvalidDuration_ReturnsError(t *testing.T) {
	cmd := NewCommand(
		func() db.TimeEntriesInterface { return &dbmocks.MockRepository{} },
		func() service.SyncInterface { return &servicemocks.MockSync{} },
		func() db.IssueStorageInterface { return &dbmocks.MockRepository{} },
	)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"PNX-123", "notaduration"})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid duration")
}

func TestAdd_DateFlag_LiteralDate(t *testing.T) {
	repo := &dbmocks.MockRepository{}
	cmd := NewCommand(
		func() db.TimeEntriesInterface { return repo },
		func() service.SyncInterface { return &servicemocks.MockSync{} },
		func() db.IssueStorageInterface { return &dbmocks.MockRepository{} },
	)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"PNX-123", "2h", "Worked on feature X", "--date", "2026-07-30"})

	err := cmd.Execute()
	assert.NoError(t, err)

	if assert.Len(t, repo.Entries, 1) {
		want := time.Date(2026, 7, 30, 0, 0, 0, 0, time.Local)
		assert.True(t, repo.Entries[0].CreatedAt.Equal(want), "got %v, want %v", repo.Entries[0].CreatedAt, want)
	}
}

func TestAdd_DateFlag_Keyword(t *testing.T) {
	repo := &dbmocks.MockRepository{}
	cmd := NewCommand(
		func() db.TimeEntriesInterface { return repo },
		func() service.SyncInterface { return &servicemocks.MockSync{} },
		func() db.IssueStorageInterface { return &dbmocks.MockRepository{} },
	)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"PNX-123", "2h", "Worked on feature X", "--date", "yesterday"})

	err := cmd.Execute()
	assert.NoError(t, err)

	if assert.Len(t, repo.Entries, 1) {
		yesterday := time.Now().AddDate(0, 0, -1)
		want := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, time.Local)
		assert.True(t, repo.Entries[0].CreatedAt.Equal(want), "got %v, want %v", repo.Entries[0].CreatedAt, want)
	}
}

func TestAdd_DateFlag_Invalid_ReturnsError(t *testing.T) {
	cmd := NewCommand(
		func() db.TimeEntriesInterface { return &dbmocks.MockRepository{} },
		func() service.SyncInterface { return &servicemocks.MockSync{} },
		func() db.IssueStorageInterface { return &dbmocks.MockRepository{} },
	)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"PNX-123", "2h", "Worked on feature X", "--date", "not-a-date"})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unrecognised date")
}

func TestAdd_NoDescription(t *testing.T) {
	cmd := NewCommand(
		func() db.TimeEntriesInterface { return &dbmocks.MockRepository{} },
		func() service.SyncInterface { return &servicemocks.MockSync{} },
		func() db.IssueStorageInterface { return &dbmocks.MockRepository{} },
	)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"PNX-123", "2h"})

	err := cmd.Execute()
	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "Added time entry: ID=42\n")
}
