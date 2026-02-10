package list

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/model"
	"github.com/previousnext/tl-go/internal/util"
)

func NewCommand(r func() db.TimeEntriesInterface) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "list",
		Args:                  cobra.NoArgs,
		DisableFlagsInUseLine: true,
		Short:                 "List all time entries",
		Long:                  "List all time entries in the database.",
		Example: `  # List all time entries
  tl list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := r().FindAllTimeEntries()
			if err != nil {
				return err
			}
			if len(entries) == 0 {
				cmd.Println("No time entries found.")
				return nil
			}

			header := []string{
				"ID",
				"Key",
				"Project",
				"Summary",
				"Duration",
				"Description",
				"Sent",
			}

			var rows [][]string

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
			}

			return util.PrintTable(cmd.OutOrStdout(), header, rows)
		},
	}

	return cmd
}
