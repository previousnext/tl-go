package setup

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/db/mocks"
)

func TestNewCommand_PrintsSuccessMessage(t *testing.T) {
	mock := &mocks.MockRepository{}

	cmd := NewCommand(func() db.RepositoryInterface { return mock })

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Successfully initialized repository")
}
