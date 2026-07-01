package comment

import (
	"bufio"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/model"
	"github.com/previousnext/tl-go/internal/util"
)

var (
	cmdShort = "Add descriptions to un-commented time entries"
	cmdLong  = `Cycle through unsent time entries that have no description and prompt
for one for each, so you don't have to edit each entry by ID.

Press Enter to skip an entry, or type "q" to stop.`
	cmdExample = `
  # Comment all unsent entries without a description
  tl comment

  # Only comment entries from today
  tl comment --date today`
	flagDate = ""
)

func NewCommand(r func() db.TimeEntriesInterface) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "comment",
		Args:                  cobra.NoArgs,
		DisableFlagsInUseLine: true,
		Short:                 cmdShort,
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			storage := r()

			entries, err := storage.FindUnsentTimeEntriesWithoutDescription()
			if err != nil {
				return err
			}

			// Optionally restrict to a date range.
			if flagDate != "" {
				start, end, _, err := util.ParseHumanDate(flagDate, time.Now())
				if err != nil {
					return err
				}
				entries = filterByRange(entries, start, end)
			}

			if len(entries) == 0 {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), "No entries without a description found.")
				return nil
			}

			reader := bufio.NewReader(cmd.InOrStdin())
			updated := 0

			for _, entry := range entries {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Entry %d | %s | %s | %s\n",
					entry.ID,
					entry.CreatedAt.Local().Format(time.DateOnly),
					entry.IssueKey,
					model.FormatDuration(entry.Duration),
				)
				_, _ = fmt.Fprint(cmd.OutOrStdout(), "Description [skip=Enter, quit=q]: ")

				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)

				if strings.EqualFold(input, "q") {
					break
				}
				if input == "" {
					continue
				}

				entry.Description = input
				if err := storage.UpdateTimeEntry(entry); err != nil {
					return err
				}
				updated++
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Updated time entry %d.\n", entry.ID)
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Updated %d of %d entries.\n", updated, len(entries))
			return nil
		},
	}

	cmd.Flags().StringVarP(&flagDate, "date", "d", "", "Only comment entries within a date range (YYYY-MM-DD or 'today', 'yesterday', 'this week', 'last week', 'this month', 'last month')")

	return cmd
}

func filterByRange(entries []*model.TimeEntry, start, end time.Time) []*model.TimeEntry {
	var filtered []*model.TimeEntry
	for _, entry := range entries {
		created := entry.CreatedAt
		if (created.Equal(start) || created.After(start)) && (created.Equal(end) || created.Before(end)) {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}
