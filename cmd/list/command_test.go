package list

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/db/mocks"
	"github.com/previousnext/tl-go/internal/model"
)

func twoEntries(start, end time.Time) ([]*model.TimeEntry, error) {
	category1 := &model.Category{Model: gorm.Model{ID: 1}, Name: "Category1"}
	category2 := &model.Category{Model: gorm.Model{ID: 2}, Name: "Category2"}
	project1 := model.Project{Name: "Project1", Category: category1}
	project2 := model.Project{Name: "Project2", Category: category2}
	return []*model.TimeEntry{
		{
			Model:    gorm.Model{ID: 1},
			IssueKey: "PNX-1",
			Issue: &model.Issue{
				Summary: "issue1",
				Project: project1,
			},
			Duration:    2 * time.Hour,
			Description: "Worked on X",
		},
		{
			Model:    gorm.Model{ID: 2},
			IssueKey: "PNX-2",
			Issue: &model.Issue{
				Summary: "issue2",
				Project: project2,
			},
			Duration:    30 * time.Minute,
			Description: "Reviewed Y",
		},
	}, nil
}

func TestNewCommand_PrintsEntriesInTable(t *testing.T) {
	mock := &mocks.MockRepository{
		FindTimeEntriesInRangeFunc: twoEntries,
	}
	cmd := NewCommand(func() db.TimeEntriesInterface { return mock })
	cmd.SetArgs([]string{"--output=wide"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	assert.NoError(t, err)
	output := buf.String()
	fmt.Print(output)
	assert.Contains(t, output, "PNX-1")
	assert.Contains(t, output, "PNX-2")
	assert.Contains(t, output, "issue1")
	assert.Contains(t, output, "issue2")
	assert.Contains(t, output, "Worked on X")
	assert.Contains(t, output, "Reviewed Y")
	assert.Contains(t, output, "2h")
	assert.Contains(t, output, "30m")
}

func TestNewCommand_PrintsEntriesWithNilCategory(t *testing.T) {
	mock := &mocks.MockRepository{
		FindTimeEntriesInRangeFunc: func(start, end time.Time) ([]*model.TimeEntry, error) {
			project1 := model.Project{Name: "Project1", Category: nil}
			project2 := model.Project{Name: "Project2", Category: nil}
			return []*model.TimeEntry{
				{
					Model:    gorm.Model{ID: 1},
					IssueKey: "PNX-1",
					Issue: &model.Issue{
						Summary: "issue1",
						Project: project1,
					},
					Duration:    2 * time.Hour,
					Description: "Worked on X",
				},
				{
					Model:    gorm.Model{ID: 2},
					IssueKey: "PNX-2",
					Issue: &model.Issue{
						Summary: "issue2",
						Project: project2,
					},
					Duration:    30 * time.Minute,
					Description: "Reviewed Y",
				},
			}, nil
		},
	}
	cmd := NewCommand(func() db.TimeEntriesInterface { return mock })
	cmd.SetArgs([]string{"--output=wide"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "PNX-1")
	assert.Contains(t, output, "PNX-2")
	assert.Contains(t, output, "issue1")
	assert.Contains(t, output, "issue2")
	assert.Contains(t, output, "Worked on X")
	assert.Contains(t, output, "Reviewed Y")
	assert.Contains(t, output, "2h")
	assert.Contains(t, output, "30m")
}

func TestNewCommand_HumanDateKeywords(t *testing.T) {
	tests := []struct {
		dateFlag  string
		wantLabel string
	}{
		{dateFlag: "today", wantLabel: "today"},
		{dateFlag: "yesterday", wantLabel: "yesterday"},
		{dateFlag: "last week", wantLabel: "last week"},
		{dateFlag: "this week", wantLabel: "this week"},
		{dateFlag: "last month", wantLabel: "last month"},
		{dateFlag: "this month", wantLabel: "this month"},
		{dateFlag: "2026-01-15", wantLabel: "2026-01-15"},
	}

	for _, tt := range tests {
		t.Run(tt.dateFlag, func(t *testing.T) {
			mock := &mocks.MockRepository{
				FindTimeEntriesInRangeFunc: twoEntries,
			}
			cmd := NewCommand(func() db.TimeEntriesInterface { return mock })
			cmd.SetArgs([]string{"--date", tt.dateFlag})

			var buf bytes.Buffer
			cmd.SetOut(&buf)

			err := cmd.Execute()
			assert.NoError(t, err)
			assert.Contains(t, buf.String(), tt.wantLabel)
		})
	}
}

func TestNewCommand_DefaultDateIsToday(t *testing.T) {
	mock := &mocks.MockRepository{
		FindTimeEntriesInRangeFunc: twoEntries,
	}
	cmd := NewCommand(func() db.TimeEntriesInterface { return mock })

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "today")
}

func TestNewCommand_InvalidDate(t *testing.T) {
	mock := &mocks.MockRepository{}
	cmd := NewCommand(func() db.TimeEntriesInterface { return mock })
	cmd.SetArgs([]string{"--date", "not-a-date"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	assert.Error(t, err)
}

func TestNewCommand_NoEntries(t *testing.T) {
	mock := &mocks.MockRepository{
		FindTimeEntriesInRangeFunc: func(start, end time.Time) ([]*model.TimeEntry, error) {
			return []*model.TimeEntry{}, nil
		},
	}
	cmd := NewCommand(func() db.TimeEntriesInterface { return mock })

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "No time entries found")
}
