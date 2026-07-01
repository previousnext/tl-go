package mocks

import (
	"github.com/previousnext/tl-go/internal/model"
	"github.com/previousnext/tl-go/internal/service"
)

// MockTimerEntryService is a test double for service.TimerEntryServiceInterface.
// Only the methods exercised by tests need a corresponding func field set;
// unset methods fall back to the embedded interface and will panic if called.
type MockTimerEntryService struct {
	service.TimerEntryServiceInterface

	GetTimerEntryByIDFunc func(id uint) (*model.TimerEntry, error)
	DeleteTimerEntryFunc  func(id uint) (*model.TimerEntry, error)
}

func (m *MockTimerEntryService) GetTimerEntryByID(id uint) (*model.TimerEntry, error) {
	if m.GetTimerEntryByIDFunc != nil {
		return m.GetTimerEntryByIDFunc(id)
	}
	return nil, nil
}

func (m *MockTimerEntryService) DeleteTimerEntry(id uint) (*model.TimerEntry, error) {
	if m.DeleteTimerEntryFunc != nil {
		return m.DeleteTimerEntryFunc(id)
	}
	return nil, nil
}
