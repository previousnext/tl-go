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
	flagOutput = "table"
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
				"Time",
				"Description",
			}
			if flagOutput == "wide" {
				header = append(header, "Summary", "Project")
			}

			var rows [][]string

			totalDuration := time.Duration(0)
			for _, entry := range entries {
				row := []string{
					fmt.Sprintf("%d", entry.ID),
					entry.CreatedAt.Format(time.DateOnly),
					entry.IssueKey,
					model.FormatDuration(entry.Duration),
					entry.Description,
				}
				if flagOutput == "wide" {
					row = append(row, entry.Issue.Summary, entry.Issue.Project.Name)
				}
				rows = append(rows, row)
				totalDuration += entry.Duration
			}

			footer := []string{
				"",
				"",
				util.ApplyHeaderFormatting("Total"),
				util.ApplyHeaderFormatting(model.FormatDuration(totalDuration)),
				"",
			}
			if flagOutput == "wide" {
				footer = append(footer, "", "")
			}

			err = util.PrintTable(cmd.OutOrStdout(), header, rows, footer)
			if err != nil {
				return fmt.Errorf("error printing table: %w", err)
			}
			util.PrintIssueLinks(cmd.OutOrStdout(), entries)
			
			return nil
		},
	}

	cmd.Flags().StringVarP(&flagOutput, "output", "o", flagOutput, "Output format (table, wide). Defaults to table.")

	return cmd
}
