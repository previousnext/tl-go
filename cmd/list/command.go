package list

import (
	"bytes"
	"fmt"
	"time"

	"github.com/aquasecurity/table"
	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/model"
	"github.com/previousnext/tl-go/internal/util"
)

var (
	cmdShort   = `List all time entries`
	cmdLong    = `List all time entries in the database.`
	flagDate   = ""
	flagOutput = "table"
	cmdExample = `
  # List all time entries for today
  tl list

  # List entries using human-friendly date keywords
  tl list --date today
  tl list --date yesterday
  tl list --date "last week"
  tl list --date "this month"

  # List entries for a specific date
  tl list --date 2026-01-15`
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

			input := flagDate
			if input == "" {
				input = "today"
			}

			start, end, dateOutput, err := util.ParseHumanDate(input, time.Now())
			if err != nil {
				return err
			}

			entries, err := r().FindTimeEntriesInRange(start, end)
			if err != nil {
				return err
			}

			if len(entries) == 0 {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "No time entries found for %s.\n", dateOutput)
				return nil
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Found %d time entries for %s.\n", len(entries), dateOutput)

			var b bytes.Buffer
			t := table.New(&b)

			util.ApplyDefaultTableFormatting(t)

			headers := []string{
				"ID",
				"Issue",
				"Cat",
				"Time",
				"AI Saved",
				"Description",
				"Sent",
			}
			if flagOutput == "wide" {
				headers = append(headers, "Project", "Summary")
			}

			var formattedHeaders []string
			for _, h := range headers {
				formattedHeaders = append(formattedHeaders, util.ApplyHeaderFormatting(h))
			}

			t.SetHeaders(formattedHeaders...)

			// Keep track of the total duration for all entries to display in the footer.
			totalDuration := time.Duration(0)

			for _, entry := range entries {
				categoryName := ""
				if entry.Issue.Project.Category != nil {
					categoryName = entry.Issue.Project.Category.Name
				}
				row := []string{
					fmt.Sprintf("%d", entry.ID),
					entry.IssueKey, // Only show plain issue key in table
					util.AbbreviateProjectCategory(categoryName),
					model.FormatDuration(entry.Duration),
					model.FormatDuration(entry.AISavedDuration),
					entry.Description,
					util.FormatBool(entry.Sent),
				}
				if flagOutput == "wide" {
					row = append(row, entry.Issue.Project.Name, entry.Issue.Summary)
				}
				t.AddRow(row...)
				totalDuration += entry.Duration
			}

			footer := []string{
				"",
				"",
				util.ApplyHeaderFormatting("Total"),
				util.ApplyHeaderFormatting(model.FormatDuration(totalDuration)),
				"",
			}

			t.SetFooters(footer...)

			t.SetFooterAlignment(
				table.AlignLeft,
				table.AlignLeft,
				table.AlignLeft,
				table.AlignRight,
				table.AlignLeft,
				table.AlignRight,
				table.AlignLeft,
				table.AlignLeft,
			)

			t.Render()

			_, err = fmt.Fprintf(cmd.OutOrStdout(), "\n%s\n", b.String())
			if err != nil {
				return fmt.Errorf("failed to print table: %w", err)
			}
			util.PrintIssueLinks(cmd.OutOrStdout(), entries)
			return nil
		},
	}

	cmd.Flags().StringVarP(&flagDate, "date", "d", "", "Date to list entries for (YYYY-MM-DD or 'today', 'yesterday', 'last week', 'this week', 'last month', 'this month')")
	cmd.Flags().StringVarP(&flagOutput, "output", "o", flagOutput, "Output format (table,wide).")

	return cmd
}
