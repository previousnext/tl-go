package delete

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/db/mocks"
)

func TestNewCommand_DeletesEntryAndPrintsMessage(t *testing.T) {
	mock := &mocks.MockRepository{}
	cmd := NewCommand(func() db.TimeEntriesInterface { return mock })

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"123"})

	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Deleted time entry with ID 123")
}

func TestNewCommand_InvalidID_ReturnsError(t *testing.T) {
	mock := &mocks.MockRepository{}
	cmd := NewCommand(func() db.TimeEntriesInterface { return mock })

	cmd.SetArgs([]string{"notanumber"})
	err := cmd.Execute()
	assert.Error(t, err)
}
