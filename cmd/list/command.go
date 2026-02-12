package list

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/model"
	"github.com/previousnext/tl-go/internal/util"
)

var (
	cmdShort   = `List all time entries`
	cmdLong    = `List all time entries in the database.`
	flagDate   = ""
	cmdExample = `
  # List all time entries
  tl list`
)

func NewCommand(r func() db.TimeEntriesInterface) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "list",
		Args:                  cobra.NoArgs,
		DisableFlagsInUseLine: true,
		Short:                 cmdShort,
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {

			// Default to today's date if no date flag is provided
			d := time.Now()
			dateOutput := "today"
			if flagDate != "" {
				var err error
				d, err = time.ParseInLocation(time.DateOnly, flagDate, time.Local)
				if err != nil {
					return fmt.Errorf("invalid d format: %s. Expected YYYY-MM-DD", flagDate)
				}
				dateOutput = flagDate
			}

			entries, err := r().FindAllTimeEntries(d)
			if err != nil {
				return err
			}

			if len(entries) == 0 {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "No time entries found for %s.\n", dateOutput)
				return nil
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Found %d time entries for %s.\n", len(entries), dateOutput)

			header := []string{
				"ID",
				"Issue",
				"Project",
				"Summary",
				"Time",
				"Description",
				"Sent",
			}

			var rows [][]string

			totalDuration := time.Duration(0)

			for _, entry := range entries {
				rows = append(rows, []string{
					fmt.Sprintf("%d", entry.ID),
					entry.IssueKey,
					entry.Issue.Project.Name,
					entry.Issue.Summary,
					model.FormatDuration(entry.Duration),
					entry.Description,
					util.FormatBool(entry.Sent),
				})
				totalDuration += entry.Duration
			}

			footer := []string{
				"",
				"",
				"",
				"Total",
				model.FormatDuration(totalDuration),
				"",
				"",
			}

			return util.PrintTable(cmd.OutOrStdout(), header, rows, footer)
		},
	}

	cmd.Flags().StringVarP(&flagDate, "date", "d", "", "List time entries created on a specific date (YYYY-MM-DD)")
	return cmd
}
