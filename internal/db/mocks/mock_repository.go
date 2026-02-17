package mocks

import (
	"time"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/model"
)

// MockRepository implements db.TimeEntriesInterface for testing.
type MockRepository struct {
	db.TimeEntriesInterface
	db.IssueStorageInterface
	Entries                   []*model.TimeEntry
	FindAllTimeEntriesFunc    func(date time.Time) ([]*model.TimeEntry, error)
	FindUnsentTimeEntriesFunc func() ([]*model.TimeEntry, error)
	FindTimeEntryFunc         func(id uint) (*model.TimeEntry, error)
	UpdateTimeEntryFunc       func(entry *model.TimeEntry) error
	GetSummaryByCategoryFunc  func(start, end time.Time) ([]db.CategorySummary, error)
}

var _ db.IssueStorageInterface = (*MockRepository)(nil)
var _ db.TimeEntriesInterface = (*MockRepository)(nil)

func (m *MockRepository) AutoMigrate() error {
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
func (m *MockRepository) FindAllTimeEntries(date time.Time) ([]*model.TimeEntry, error) {
	if m.FindAllTimeEntriesFunc != nil {
		return m.FindAllTimeEntriesFunc(time.Now())
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

func (m *MockRepository) FindUnsentTimeEntries() ([]*model.TimeEntry, error) {
	if m.FindUnsentTimeEntriesFunc != nil {
		return m.FindUnsentTimeEntriesFunc()
	}
	return nil, nil
}

func (m *MockRepository) GetSummaryByCategory(start, end time.Time) ([]db.CategorySummary, error) {
	if m.GetSummaryByCategoryFunc != nil {
		return m.GetSummaryByCategoryFunc(start, end)
	}
	return nil, nil
}

func (m *MockRepository) FindAllIssues() ([]*model.Issue, error) {
	return []*model.Issue{}, nil
}
func (m *MockRepository) CreateIssue(issue *model.Issue) error {
	return nil
}
func (m *MockRepository) DeleteIssueByKey(key string) error {
	return nil
}
func (m *MockRepository) DeleteAllIssues() error {
	return nil
}
func (m *MockRepository) FindIssueByKey(key string) (*model.Issue, error) {
	return nil, nil
}
func (m *MockRepository) FindRecentIssues(limit int) ([]*model.Issue, error) {
	return []*model.Issue{}, nil
}
