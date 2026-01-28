package mocks

import "github.com/previousnext/tl-go/internal/model"

// MockRepository implements db.RepositoryInterface for testing.
type MockRepository struct {
	Entries                []*model.TimeEntry
	FindAllTimeEntriesFunc func() ([]*model.TimeEntry, error)
	FindTimeEntryFunc      func(id uint) (*model.TimeEntry, error)
	UpdateTimeEntryFunc    func(entry *model.TimeEntry) error
}

func (m *MockRepository) InitRepository() error {
	return nil
}
func (m *MockRepository) CreateTimeEntry(entry *model.TimeEntry) error {
	entry.ID = 42
	m.Entries = append(m.Entries, entry)
	return nil
}
func (m *MockRepository) FindTimeEntry(id uint) (*model.TimeEntry, error) {
	if m.FindTimeEntryFunc != nil {
		return m.FindTimeEntryFunc(id)
	}
	if len(m.Entries) > 0 {
		return m.Entries[len(m.Entries)-1], nil
	}
	return nil, nil
}
func (m *MockRepository) FindAllTimeEntries() ([]*model.TimeEntry, error) {
	if m.FindAllTimeEntriesFunc != nil {
		return m.FindAllTimeEntriesFunc()
	}
	return m.Entries, nil
}
func (m *MockRepository) UpdateTimeEntry(entry *model.TimeEntry) error {
	if m.UpdateTimeEntryFunc != nil {
		return m.UpdateTimeEntryFunc(entry)
	}
	return nil
}
func (m *MockRepository) DeleteTimeEntry(id uint) error {
	if len(m.Entries) > 0 {
		m.Entries = m.Entries[:len(m.Entries)-1]
	}
	return nil
}
