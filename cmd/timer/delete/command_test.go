package delete

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/previousnext/tl-go/internal/model"
	"github.com/previousnext/tl-go/internal/service"
	"github.com/previousnext/tl-go/internal/service/mocks"
)

func newTimerServiceFunc(m *mocks.MockTimerEntryService) func() service.TimerEntryServiceInterface {
	return func() service.TimerEntryServiceInterface { return m }
}

func TestTimerDelete_ConfirmDeletes(t *testing.T) {
	var deletedID uint
	mock := &mocks.MockTimerEntryService{
		GetTimerEntryByIDFunc: func(id uint) (*model.TimerEntry, error) {
			return &model.TimerEntry{ID: id, IssueKey: "PNX-123"}, nil
		},
		DeleteTimerEntryFunc: func(id uint) (*model.TimerEntry, error) {
			deletedID = id
			return &model.TimerEntry{ID: id, IssueKey: "PNX-123"}, nil
		},
	}

	cmd := NewCommand(newTimerServiceFunc(mock))
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetIn(strings.NewReader("y\n"))
	cmd.SetArgs([]string{"3"})

	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, uint(3), deletedID)
	assert.Contains(t, buf.String(), "Timer entry has been deleted.")
}

func TestTimerDelete_DeclineAborts(t *testing.T) {
	deleteCalled := false
	mock := &mocks.MockTimerEntryService{
		GetTimerEntryByIDFunc: func(id uint) (*model.TimerEntry, error) {
			return &model.TimerEntry{ID: id, IssueKey: "PNX-123"}, nil
		},
		DeleteTimerEntryFunc: func(id uint) (*model.TimerEntry, error) {
			deleteCalled = true
			return nil, nil
		},
	}

	cmd := NewCommand(newTimerServiceFunc(mock))
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetIn(strings.NewReader("n\n"))
	cmd.SetArgs([]string{"3"})

	err := cmd.Execute()
	assert.NoError(t, err)
	assert.False(t, deleteCalled)
	assert.Contains(t, buf.String(), "Aborted.")
}

func TestTimerDelete_ForceSkipsPrompt(t *testing.T) {
	mock := &mocks.MockTimerEntryService{
		GetTimerEntryByIDFunc: func(id uint) (*model.TimerEntry, error) {
			return &model.TimerEntry{ID: id, IssueKey: "PNX-123"}, nil
		},
		DeleteTimerEntryFunc: func(id uint) (*model.TimerEntry, error) {
			return &model.TimerEntry{ID: id, IssueKey: "PNX-123"}, nil
		},
	}

	cmd := NewCommand(newTimerServiceFunc(mock))
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"3", "--force"})

	err := cmd.Execute()
	assert.NoError(t, err)
	out := buf.String()
	assert.NotContains(t, out, "[y/N]")
	assert.Contains(t, out, "Timer entry has been deleted.")
}

func TestTimerDelete_InvalidID(t *testing.T) {
	mock := &mocks.MockTimerEntryService{}
	cmd := NewCommand(newTimerServiceFunc(mock))
	cmd.SetArgs([]string{"notanumber"})

	err := cmd.Execute()
	assert.Error(t, err)
}

func TestTimerDelete_NotFound(t *testing.T) {
	mock := &mocks.MockTimerEntryService{
		GetTimerEntryByIDFunc: func(id uint) (*model.TimerEntry, error) {
			return nil, assert.AnError
		},
	}
	cmd := NewCommand(newTimerServiceFunc(mock))
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"999"})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timer entry not found")
}
