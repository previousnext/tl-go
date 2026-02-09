package service

import (
	"fmt"
	"strconv"

	"gorm.io/gorm"

	"github.com/previousnext/tl-go/internal/api"
	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/model"
)

type FetchServiceInterface interface {
	FetchIssues(issueKeys []string) error
}

type FetchService struct {
	jiraClient    func() api.JiraClientInterface
	issuesStorage func() db.IssueStorageInterface
}

func NewFetchService(issueStorage func() db.IssueStorageInterface, jiraClient func() api.JiraClientInterface) FetchServiceInterface {
	return &FetchService{
		issuesStorage: issueStorage,
		jiraClient:    jiraClient,
	}
}

func (f *FetchService) FetchIssue(issueKey string) error {
	issue, err := f.issuesStorage().FindIssueByKey(issueKey)
	if err != nil {
		return fmt.Errorf("error checking issue in database: %w", err)
	}
	if issue != nil {
		return nil // Issue already exists in the database, no need to fetch
	}

	return f.FetchIssues([]string{issueKey})
}

func (f *FetchService) FetchIssues(issueKeys []string) error {
	issuesResponse, err := f.jiraClient().BulkFetchIssues(issueKeys)
	if err != nil {
		return fmt.Errorf("error fetching issues from Jira: %w", err)
	}
	for _, issueResp := range issuesResponse.Issues {
		issueID, err := strconv.Atoi(issueResp.ID)
		if err != nil {
			return fmt.Errorf("error converting issue ID to integer: %w", err)
		}
		projectID, err := strconv.ParseUint(issueResp.Fields.Project.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("error converting project ID to uint: %w", err)
		}
		project := model.Project{
			Model: gorm.Model{
				ID: uint(projectID),
			},
			Key:  issueResp.Fields.Project.Key,
			Name: issueResp.Fields.Project.Name,
		}
		issue := model.Issue{
			Model: gorm.Model{
				ID: uint(issueID),
			},
			Key:       issueResp.Key,
			Summary:   issueResp.Fields.Summary,
			ProjectID: uint(projectID),
			Project:   project,
		}
		if err := f.issuesStorage().CreateIssue(&issue); err != nil {
			return fmt.Errorf("error saving issue to database: %w", err)
		}
	}

	return nil
}
