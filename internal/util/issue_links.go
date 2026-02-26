package util

import (
	"fmt"
	"io"
	"sort"

	"github.com/spf13/viper"

	"github.com/previousnext/tl-go/internal/model"
)

func PrintIssueLinks(w io.Writer, entries []*model.TimeEntry) {
	jiraBaseURL := viper.GetString("jira_base_url")
	if jiraBaseURL == "" || len(entries) == 0 {
		return
	}
	uniqueKeys := make(map[string]struct{})
	for _, entry := range entries {
		uniqueKeys[entry.IssueKey] = struct{}{}
	}
	sortedKeys := make([]string, 0, len(uniqueKeys))
	for issueKey := range uniqueKeys {
		sortedKeys = append(sortedKeys, issueKey)
	}
	sort.Strings(sortedKeys)
	_, _ = fmt.Fprintln(w, ApplyHeaderFormatting("Issue Links:"))
	for _, issueKey := range sortedKeys {
		_, _ = fmt.Fprintf(w, " ∙ %s/browse/%s\n", jiraBaseURL, issueKey)
	}
}
