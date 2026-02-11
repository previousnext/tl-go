package unsent

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/model"
	"github.com/previousnext/tl-go/internal/util"
)

var (
	cmdShort   = `List unsent time entries`
	cmdLong    = `List unsent time entries in the database.`
	cmdExample = `
  # List all time entries
  tl unsent`
)

func NewCommand(r func() db.TimeEntriesInterface) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "unsent",
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
				"Created",
				"Key",
				"Project",
				"Summary",
				"Duration",
				"Description",
			}

			var rows [][]string

			for _, entry := range entries {
				rows = append(rows, []string{
					fmt.Sprintf("%d", entry.ID),
					entry.CreatedAt.Format(time.DateOnly),
					entry.IssueKey,
					entry.Issue.Project.Name,
					entry.Issue.Summary,
					model.FormatDuration(entry.Duration),
					entry.Description,
				})
			}

			return util.PrintTable(cmd.OutOrStdout(), header, rows)
		},
	}

	return cmd
}
