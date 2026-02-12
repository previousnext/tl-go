package review

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/model"
	"github.com/previousnext/tl-go/internal/util"
)

var (
	cmdShort   = `Review unsent time entries`
	cmdLong    = `Review unsent time entries in the database.`
	cmdExample = `
  # Review unsent time entries
  tl review`
)

func NewCommand(r func() db.TimeEntriesInterface) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "review",
		Args:                  cobra.NoArgs,
		DisableFlagsInUseLine: true,
		Short:                 cmdShort,
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {

			entries, err := r().FindUnsentTimeEntries()
			if err != nil {
				return err
			}

			if len(entries) == 0 {
				cmd.Println("No unsent time entries found.")
				return nil
			}
			header := []string{
				"ID",
				"Date",
				"Issue",
				"Summary",
				"Project",
				"Time",
				"Description",
			}

			var rows [][]string

			totalDuration := time.Duration(0)
			for _, entry := range entries {
				rows = append(rows, []string{
					fmt.Sprintf("%d", entry.ID),
					entry.CreatedAt.Format(time.DateOnly),
					entry.IssueKey,
					entry.Issue.Summary,
					entry.Issue.Project.Name,
					model.FormatDuration(entry.Duration),
					entry.Description,
				})
				totalDuration += entry.Duration
			}
			footer := []string{
				"",
				"",
				"",
				"",
				"Total",
				model.FormatDuration(totalDuration),
				"",
			}

			return util.PrintTable(cmd.OutOrStdout(), header, rows, footer)
		},
	}

	return cmd
}
