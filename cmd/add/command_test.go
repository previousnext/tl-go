package add

import (
	"bytes"
	"testing"

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
	assert.Contains(t, output, "Added time entry: ID=42")
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
	assert.Contains(t, output, "Added time entry: ID=42")
}
