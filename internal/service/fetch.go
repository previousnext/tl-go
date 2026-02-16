package service

import (
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"

	"github.com/previousnext/tl-go/internal/api"
	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/model"
)

type SyncInterface interface {
	SyncIssue(issueKey string, options ...SyncOption) (*model.Issue, error)
	SyncIssues(issueKeys []string) error
}

type FetchService struct {
	jiraClient    func() api.JiraClientInterface
	issuesStorage func() db.IssueStorageInterface
}

type SyncOptions struct {
	Force bool
}

type SyncOption func(*SyncOptions)

func NewSync(issueStorage func() db.IssueStorageInterface, jiraClient func() api.JiraClientInterface) SyncInterface {
	return &FetchService{
		issuesStorage: issueStorage,
		jiraClient:    jiraClient,
	}
}

func WithForce(force bool) SyncOption {
	return func(opts *SyncOptions) {
		opts.Force = force
	}
}

func (f *FetchService) SyncIssue(key string, options ...SyncOption) (*model.Issue, error) {
	force := false
	if len(options) > 0 {
		opts := &SyncOptions{
			Force: false,
		}
		for _, opt := range options {
			opt(opts)
		}
		if opts.Force {
			force = true
		}
	}

	// Check if the issue already exists in the database.
	issue, err := f.issuesStorage().FindIssueByKey(key)
	if err != nil {
		return nil, fmt.Errorf("error checking issue in database: %w", err)
	}
	if issue != nil && !force {
		return issue, nil // Issue already exists in the database. No need to update unless force is true.
	}

	// Check if the issue exists in Jira.
	issueResp, err := f.jiraClient().FetchIssue(key)
	if err != nil {
		if errors.Is(err, api.ErrNotFound) {
			return nil, fmt.Errorf("issue with key %s not found in Jira", key)
		}
		return nil, fmt.Errorf("error fetching issue from Jira: %w", err)
	}

	// Create a new issue.
	issue, err = f.doCreateIssue(issueResp)
	if err != nil {
		return nil, fmt.Errorf("error creating issue in database: %w", err)
	}
	return issue, nil
}

func (f *FetchService) SyncIssues(issueKeys []string) error {
	issuesResponse, err := f.jiraClient().BulkFetchIssues(issueKeys)
	if err != nil {
		return fmt.Errorf("error fetching issues from Jira: %w", err)
	}
	for _, issueResp := range issuesResponse.Issues {
		_, err = f.doCreateIssue(issueResp)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *FetchService) doCreateIssue(issueResp api.IssueResponse) (*model.Issue, error) {
	issueID, err := strconv.Atoi(issueResp.ID)
	if err != nil {
		return nil, fmt.Errorf("error converting issue ID to integer: %w", err)
	}

	projectID, err := strconv.ParseUint(issueResp.Fields.Project.ID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error converting project ID to uint: %w", err)
	}

	categoryID, err := strconv.ParseUint(issueResp.Fields.Project.ProjectCategory.ID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error converting category ID to uint: %w", err)
	}

	catID := uint(categoryID)
	project := model.Project{
		Model: gorm.Model{
			ID: uint(projectID),
		},
		Key:        issueResp.Fields.Project.Key,
		Name:       issueResp.Fields.Project.Name,
		CategoryID: &catID,
		Category: &model.Category{
			Model: gorm.Model{
				ID: uint(categoryID),
			},
			Name: issueResp.Fields.Project.ProjectCategory.Name,
		},
	}

	issue := &model.Issue{
		Model: gorm.Model{
			ID: uint(issueID),
		},
		Key:       issueResp.Key,
		Summary:   issueResp.Fields.Summary,
		ProjectID: uint(projectID),
		Project:   project,
	}
	if err := f.issuesStorage().CreateIssue(issue); err != nil {
		return issue, fmt.Errorf("error saving issue to database: %w", err)
	}

	return issue, nil
}
