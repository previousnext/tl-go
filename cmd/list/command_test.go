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

func TestNewCommand_PrintsEntriesInTable(t *testing.T) {
	mock := &mocks.MockRepository{
		FindAllTimeEntriesFunc: func(date time.Time) ([]*model.TimeEntry, error) {
			project1 := &model.Project{Name: "Project1"}
			project2 := &model.Project{Name: "Project2"}
			return []*model.TimeEntry{
				{
					Model:    gorm.Model{ID: 1},
					IssueKey: "PNX-1",
					Issue: &model.Issue{
						Summary: "issue1",
						Project: *project1,
					},
					Duration:    2 * time.Hour,
					Description: "Worked on X",
				},
				{
					Model:    gorm.Model{ID: 2},
					IssueKey: "PNX-2",
					Issue: &model.Issue{
						Summary: "issue2",
						Project: *project2,
					},
					Duration:    30 * time.Minute,
					Description: "Reviewed Y",
				},
			}, nil
		},
	}
	cmd := NewCommand(func() db.TimeEntriesInterface { return mock })

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
	assert.Contains(t, output, "2h0m")
	assert.Contains(t, output, "30m")
}
