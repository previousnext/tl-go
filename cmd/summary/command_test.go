package summary

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/db/mocks"
)

func summaryMock(fn func(start, end time.Time) ([]db.CategorySummary, error)) *mocks.MockRepository {
	mock := &mocks.MockRepository{}
	mock.GetSummaryByCategoryFunc = fn
	return mock
}

func twoCategories(start, end time.Time) ([]db.CategorySummary, error) {
	return []db.CategorySummary{
		{CategoryName: "Billable", Duration: 2 * time.Hour, Percentage: 66.7},
		{CategoryName: "Non Billable", Duration: 1 * time.Hour, Percentage: 33.3},
	}, nil
}

func TestSummaryCommand_PrintsTable(t *testing.T) {
	cmd := NewCommand(func() db.TimeEntriesInterface { return summaryMock(twoCategories) })

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	assert.NoError(t, err)
	output := buf.String()

	assert.Contains(t, output, "Category")
	assert.Contains(t, output, "Total Time")
	assert.Contains(t, output, "Percentage")
	assert.Contains(t, output, "Billable")
	assert.Contains(t, output, "Non Billable")
	assert.Contains(t, output, "2h")
	assert.Contains(t, output, "1h")
	assert.Contains(t, output, "66.7%")
	assert.Contains(t, output, "33.3%")
	assert.Contains(t, output, "Total")
	assert.Contains(t, output, "3h")
	assert.Contains(t, output, "Summary of time spent per category")
}

func TestSummaryCommand_NoResults(t *testing.T) {
	cmd := NewCommand(func() db.TimeEntriesInterface {
		return summaryMock(func(start, end time.Time) ([]db.CategorySummary, error) {
			return nil, nil
		})
	})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "No time entries found for the period")
}

func TestSummaryCommand_DateKeywords(t *testing.T) {
	keywords := []string{
		"today",
		"yesterday",
		"this week",
		"last week",
		"this month",
		"last month",
	}

	for _, kw := range keywords {
		t.Run(kw, func(t *testing.T) {
			var capturedStart, capturedEnd time.Time
			cmd := NewCommand(func() db.TimeEntriesInterface {
				return summaryMock(func(start, end time.Time) ([]db.CategorySummary, error) {
					capturedStart = start
					capturedEnd = end
					return twoCategories(start, end)
				})
			})
			cmd.SetArgs([]string{"--date", kw})

			var buf bytes.Buffer
			cmd.SetOut(&buf)

			err := cmd.Execute()
			assert.NoError(t, err)
			assert.False(t, capturedStart.IsZero(), "start should be set")
			assert.False(t, capturedEnd.IsZero(), "end should be set")
			assert.True(t, capturedStart.Before(capturedEnd), "start should be before end")
			assert.Contains(t, buf.String(), kw)
		})
	}
}

func TestSummaryCommand_DateFlag_YYYYMMDD(t *testing.T) {
	var capturedStart, capturedEnd time.Time
	cmd := NewCommand(func() db.TimeEntriesInterface {
		return summaryMock(func(start, end time.Time) ([]db.CategorySummary, error) {
			capturedStart = start
			capturedEnd = end
			return twoCategories(start, end)
		})
	})
	cmd.SetArgs([]string{"--date", "2026-01-15"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, 2026, capturedStart.Year())
	assert.Equal(t, time.January, capturedStart.Month())
	assert.Equal(t, 15, capturedStart.Day())
	assert.Equal(t, 15, capturedEnd.Day())
}

func TestSummaryCommand_DateFlag_InvalidDate(t *testing.T) {
	cmd := NewCommand(func() db.TimeEntriesInterface { return summaryMock(twoCategories) })
	cmd.SetArgs([]string{"--date", "not-a-date"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	assert.Error(t, err)
}

func TestSummaryCommand_DateFlag_ConflictsWithStart(t *testing.T) {
	cmd := NewCommand(func() db.TimeEntriesInterface { return summaryMock(twoCategories) })
	cmd.SetArgs([]string{"--date", "today", "--start", "2026-01-01"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	assert.Error(t, err)
}

func TestSummaryCommand_DateFlag_ConflictsWithEnd(t *testing.T) {
	cmd := NewCommand(func() db.TimeEntriesInterface { return summaryMock(twoCategories) })
	cmd.SetArgs([]string{"--date", "today", "--end", "2026-01-31"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	assert.Error(t, err)
}

func TestSummaryCommand_StartAndEndFlags(t *testing.T) {
	var capturedStart, capturedEnd time.Time
	cmd := NewCommand(func() db.TimeEntriesInterface {
		return summaryMock(func(start, end time.Time) ([]db.CategorySummary, error) {
			capturedStart = start
			capturedEnd = end
			return twoCategories(start, end)
		})
	})
	cmd.SetArgs([]string{"--start", "2026-01-01", "--end", "2026-01-31"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, 1, capturedStart.Day())
	assert.Equal(t, time.January, capturedStart.Month())
	assert.Equal(t, 31, capturedEnd.Day())
}
