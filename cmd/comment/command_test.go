package comment

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/db/mocks"
	"github.com/previousnext/tl-go/internal/model"
)

func newRepoFunc(m *mocks.MockRepository) func() db.TimeEntriesInterface {
	return func() db.TimeEntriesInterface { return m }
}

func TestComment_UpdatesEntries(t *testing.T) {
	updated := map[uint]string{}
	e1 := &model.TimeEntry{IssueKey: "PNX-1", Duration: time.Hour}
	e1.ID = 1
	e2 := &model.TimeEntry{IssueKey: "PNX-2", Duration: time.Hour}
	e2.ID = 2
	mock := &mocks.MockRepository{
		FindUnsentTimeEntriesWithoutDescriptionFunc: func() ([]*model.TimeEntry, error) {
			return []*model.TimeEntry{e1, e2}, nil
		},
		UpdateTimeEntryFunc: func(e *model.TimeEntry) error {
			updated[e.ID] = e.Description
			return nil
		},
	}

	cmd := NewCommand(newRepoFunc(mock))
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetIn(strings.NewReader("desc1\ndesc2\n"))

	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, "desc1", updated[1])
	assert.Equal(t, "desc2", updated[2])
	assert.Contains(t, buf.String(), "Updated 2 of 2 entries.")
}

func TestComment_EmptyLineSkips(t *testing.T) {
	updates := 0
	mock := &mocks.MockRepository{
		FindUnsentTimeEntriesWithoutDescriptionFunc: func() ([]*model.TimeEntry, error) {
			return []*model.TimeEntry{
				{IssueKey: "PNX-1", Duration: time.Hour},
				{IssueKey: "PNX-2", Duration: time.Hour},
			}, nil
		},
		UpdateTimeEntryFunc: func(e *model.TimeEntry) error {
			updates++
			return nil
		},
	}

	cmd := NewCommand(newRepoFunc(mock))
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetIn(strings.NewReader("\ndesc2\n"))

	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, 1, updates)
	assert.Contains(t, buf.String(), "Updated 1 of 2 entries.")
}

func TestComment_QuitStops(t *testing.T) {
	updates := 0
	mock := &mocks.MockRepository{
		FindUnsentTimeEntriesWithoutDescriptionFunc: func() ([]*model.TimeEntry, error) {
			return []*model.TimeEntry{
				{IssueKey: "PNX-1", Duration: time.Hour},
				{IssueKey: "PNX-2", Duration: time.Hour},
			}, nil
		},
		UpdateTimeEntryFunc: func(e *model.TimeEntry) error {
			updates++
			return nil
		},
	}

	cmd := NewCommand(newRepoFunc(mock))
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetIn(strings.NewReader("q\n"))

	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, 0, updates)
	assert.Contains(t, buf.String(), "Updated 0 of 2 entries.")
}

func TestComment_NoCandidates(t *testing.T) {
	mock := &mocks.MockRepository{
		FindUnsentTimeEntriesWithoutDescriptionFunc: func() ([]*model.TimeEntry, error) {
			return nil, nil
		},
	}

	cmd := NewCommand(newRepoFunc(mock))
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetIn(strings.NewReader(""))

	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "No entries without a description found.")
}

func TestComment_DateFiltersOutOfRange(t *testing.T) {
	updates := 0
	inRange := time.Now()
	outOfRange := time.Now().AddDate(0, 0, -30)
	mock := &mocks.MockRepository{
		FindUnsentTimeEntriesWithoutDescriptionFunc: func() ([]*model.TimeEntry, error) {
			e1 := &model.TimeEntry{IssueKey: "PNX-1", Duration: time.Hour}
			e1.CreatedAt = inRange
			e2 := &model.TimeEntry{IssueKey: "PNX-2", Duration: time.Hour}
			e2.CreatedAt = outOfRange
			return []*model.TimeEntry{e1, e2}, nil
		},
		UpdateTimeEntryFunc: func(e *model.TimeEntry) error {
			updates++
			return nil
		},
	}

	cmd := NewCommand(newRepoFunc(mock))
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetIn(strings.NewReader("desc1\n"))
	cmd.SetArgs([]string{"--date", "today"})

	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, 1, updates)
	assert.Contains(t, buf.String(), "Updated 1 of 1 entries.")
}
