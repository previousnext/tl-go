package summary

import (
	"bytes"
	"testing"
	"time"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/db/mocks"
	"github.com/stretchr/testify/assert"
)

func TestSummaryCommand_PrintsTable(t *testing.T) {
	mock := &mocks.MockRepository{}
	mock.GetSummaryByCategoryFunc = func(start, end time.Time) ([]db.CategorySummary, error) {
		return []db.CategorySummary{
			{CategoryName: "Billable", Duration: 2 * time.Hour, Percentage: 66.7},
			{CategoryName: "Non Billable", Duration: 1 * time.Hour, Percentage: 33.3},
		}, nil
	}
	cmd := NewCommand(func() db.TimeEntriesInterface { return mock })

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	assert.NoError(t, err)
	output := buf.String()

	// Check for table headers
	assert.Contains(t, output, "Category")
	assert.Contains(t, output, "Total Time")
	assert.Contains(t, output, "Percentage")
	// Check for category rows
	assert.Contains(t, output, "Billable")
	assert.Contains(t, output, "Non Billable")
	// Check for formatted durations
	assert.Contains(t, output, "2h0m")
	assert.Contains(t, output, "1h0m")
	// Check for percentages
	assert.Contains(t, output, "66.7%")
	assert.Contains(t, output, "33.3%")
	// Check for total
	assert.Contains(t, output, "Total")
	assert.Contains(t, output, "3h0m")
	// Check for summary line
	assert.Contains(t, output, "Summary of time spent per category")
}

func TestSummaryCommand_NoResults(t *testing.T) {
	mock := &mocks.MockRepository{}
	mock.GetSummaryByCategoryFunc = func(start, end time.Time) ([]db.CategorySummary, error) {
		return nil, nil
	}
	cmd := NewCommand(func() db.TimeEntriesInterface { return mock })

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "No time entries found for the period")
}
