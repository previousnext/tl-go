package send

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"github.com/previousnext/tl-go/internal/api"
	"github.com/previousnext/tl-go/internal/api/types"
	"github.com/previousnext/tl-go/internal/db"
)

func NewCommand(r func() db.TimeEntriesInterface, j func() api.JiraClientInterface) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "send [entry-id]",
		Args:                  cobra.MaximumNArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Send time entries to Jira",
		Long:                  "Send all unsent time entries to the configured Jira instance, or resend a specific time entry by ID.",
		Example: `  # Send all unsent time entries to Jira
  tl send

  # Resend a specific time entry by ID
  tl send 42`,
		RunE: func(cmd *cobra.Command, args []string) error {
			repository := r()
			jiraClient := j()

			// If an entry ID is provided, resend that specific entry
			if len(args) == 1 {
				id, err := strconv.ParseUint(args[0], 10, 64)
				if err != nil {
					return fmt.Errorf("invalid entry ID: %s", args[0])
				}
				entryID := uint(id)

				timeEntry, err := repository.FindTimeEntry(entryID)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						_, _ = fmt.Fprintf(cmd.OutOrStdout(), "No entry with ID %d\n", entryID)
						return nil
					}
					return fmt.Errorf("could not find time entry with ID %d: %w", entryID, err)
				}

				worklog := types.WorklogRecord{
					IssueKey: timeEntry.IssueKey,
					Started:  timeEntry.CreatedAt,
					Duration: timeEntry.Duration,
					Comment:  timeEntry.Description,
				}
				err = jiraClient.AddWorkLog(worklog)
				if err != nil {
					return fmt.Errorf("failed to send time entry ID %d to Jira: %w", timeEntry.ID, err)
				}

				timeEntry.Sent = true
				err = repository.UpdateTimeEntry(timeEntry)
				if err != nil {
					return fmt.Errorf("failed to mark time entry ID %d as sent: %w", timeEntry.ID, err)
				}

				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Sent time entry ID %d to Jira\n", entryID)
				return nil
			}

			// Otherwise, send all unsent entries
			unsentEntries, err := repository.FindUnsentTimeEntries()
			if err != nil {
				return fmt.Errorf("could not find time entries to send: %v", err)
			}

			count := len(unsentEntries)
			if count == 0 {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), "No unsent time entries found.")
				return nil
			}

			for _, timeEntry := range unsentEntries {
				worklog := types.WorklogRecord{
					IssueKey: timeEntry.IssueKey,
					Started:  timeEntry.CreatedAt,
					Duration: timeEntry.Duration,
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
