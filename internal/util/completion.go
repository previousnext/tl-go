package util

import (
	"github.com/previousnext/tl-go/internal/alias"
	"github.com/previousnext/tl-go/internal/db"
)

// CompleteAliasesAndIssueKeys returns a list of aliases and issue keys for shell completion.
func CompleteAliasesAndIssueKeys(issueStorage db.IssueStorageInterface) ([]string, error) {
	aliasStorage := alias.NewAliasStorage()
	aliases, err := aliasStorage.LoadAliases()
	completions := []string{}
	if err == nil {
		for k := range aliases {
			completions = append(completions, k)
		}
	}
	issues, err := issueStorage.FindAllIssues()
	if err == nil {
		for _, issue := range issues {
			completions = append(completions, issue.Key)
		}
	}
	return completions, nil
}
