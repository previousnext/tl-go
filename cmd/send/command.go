package send

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/internal/api"
	"github.com/previousnext/tl-go/internal/api/types"
	"github.com/previousnext/tl-go/internal/db"
)

func NewCommand(r func() db.RepositoryInterface, j func() api.JiraClientInterface) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "send",
		Args:                  cobra.NoArgs,
		DisableFlagsInUseLine: true,
		Short:                 "Send time entries to Jira",
		Long:                  "Send all unsent time entries to the configured Jira instance.",
		Example: `  # Send all unsent time entries to Jira
  tl send`,
		RunE: func(cmd *cobra.Command, args []string) error {
			repository := r()
			unsentEntries, err := repository.FindUnsentTimeEntries()
			if err != nil {
				return fmt.Errorf("could not find time entries to send: %v", err)
			}

			jiraClient := j()

			count := len(unsentEntries)
			if count == 0 {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), "No unsent time entries found.")
				return nil
			}

			for _, timeEntry := range unsentEntries {
				worklog := types.WorklogRecord{
					IssueKey: timeEntry.IssueKey,
					Started:  timeEntry.CreatedAt,
					Duration: time.Minute * time.Duration(timeEntry.Duration),
					Comment:  timeEntry.Description,
				}
				err := jiraClient.AddWorkLog(worklog)
				if err != nil {
					return fmt.Errorf("failed to send time entry ID %d to Jira: %v", timeEntry.ID, err)
				}

				timeEntry.Sent = true
				err = repository.UpdateTimeEntry(timeEntry)
				if err != nil {
					return fmt.Errorf("failed to mark time entry ID %d as sent: %v", timeEntry.ID, err)
				}
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Added %d worklogs to Jira\n", count)

			return nil
		},
	}
	return cmd
}
