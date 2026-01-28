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
		FindAllTimeEntriesFunc: func() ([]*model.TimeEntry, error) {
			return []*model.TimeEntry{
				{Model: gorm.Model{ID: 1}, IssueKey: "PNX-1", Duration: uint(2 * time.Hour.Minutes()), Description: "Worked on X"},
				{Model: gorm.Model{ID: 2}, IssueKey: "PNX-2", Duration: uint(30 * time.Minute.Minutes()), Description: "Reviewed Y"},
			}, nil
		},
	}
	cmd := NewCommand(func() db.RepositoryInterface { return mock })

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	assert.NoError(t, err)
	output := buf.String()
	fmt.Print(output)
	assert.Contains(t, output, "ID Issue Key Duration Description")
	assert.Contains(t, output, "1  PNX-1     2h0m     Worked on X")
	assert.Contains(t, output, "2  PNX-2     30m      Reviewed Y")
}
