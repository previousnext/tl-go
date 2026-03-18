package summary

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
	cmdShort   = "Show a summary of time spent per project category"
	cmdLong    = "Show a summary of time spent per project category for a given period."
	cmdExample = `
  # Show a summary for the current week (default)
  tl summary

  # Show a summary using human-friendly date keywords
  tl summary --date "last week"
  tl summary --date "this month"
  tl summary --date yesterday

  # Show a summary for a specific date range
  tl summary --start 2026-01-01 --end 2026-01-31`
	flagDate  = ""
	flagStart = ""
	flagEnd   = ""
)

func NewCommand(r func() db.TimeEntriesInterface) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "summary",
		Args:                  cobra.NoArgs,
		DisableFlagsInUseLine: true,
		Short:                 cmdShort,
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			var start, end time.Time
			var label string
			var err error

			if flagDate != "" {
				start, end, label, err = util.ParseHumanDate(flagDate, time.Now())
				if err != nil {
					return err
				}
			} else {
				// Resolve --start
				if flagStart == "" {
					start, end, label, _ = util.ParseHumanDate("this week", time.Now())
				} else {
					start, err = time.ParseInLocation("2006-01-02", flagStart, time.Local)
					if err != nil {
						return err
					}
					start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
				}

				// Resolve --end (only when not already set by the this-week default)
				if flagStart != "" || flagEnd != "" {
					if flagEnd == "" {
						// Default to a 7-day inclusive window starting at --start
						end = start.AddDate(0, 0, 6)
					} else {
						end, err = time.ParseInLocation("2006-01-02", flagEnd, time.Local)
						if err != nil {
							return err
						}
					}
					end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, end.Location())
					label = fmt.Sprintf("%s to %s", start.Format("2006-01-02"), end.Format("2006-01-02"))
				}
			}

			if start.After(end) {
				return fmt.Errorf("start date (%s) must be before or equal to end date (%s)", start.Format("2006-01-02"), end.Format("2006-01-02"))
			}

			summaries, err := r().GetSummaryByCategory(start, end)
			if err != nil {
				return err
			}
			if len(summaries) == 0 {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "No time entries found for the period %s to %s.\n", start.Format("2006-01-02"), end.Format("2006-01-02"))
				return nil
			}

			var b bytes.Buffer
			t := table.New(&b)
			util.ApplyDefaultTableFormatting(t)
			headers := []string{"Category", "Total Time", "Percentage"}
			var formattedHeaders []string
			for _, h := range headers {
				formattedHeaders = append(formattedHeaders, util.ApplyHeaderFormatting(h))
			}
			t.SetHeaders(formattedHeaders...)

			var total time.Duration
			for _, s := range summaries {
				total += s.Duration
			}

			for _, s := range summaries {
				t.AddRow(s.CategoryName, model.FormatDuration(s.Duration), fmt.Sprintf("%.1f%%", s.Percentage))
			}

			footer := []string{util.ApplyHeaderFormatting("Total"), model.FormatDuration(total), ""}
			t.SetFooters(footer...)

			t.SetFooterAlignment(table.AlignRight, table.AlignLeft, table.AlignLeft)

			t.Render()
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Summary of time spent per category from %s:\n", label)
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\n", b.String())
			return nil
		},
	}

	cmd.Flags().StringVar(&flagDate, "date", "", "Date range keyword ('today', 'yesterday', 'this week', 'last week', 'this month', 'last month') or YYYY-MM-DD")
	cmd.Flags().StringVar(&flagStart, "start", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&flagEnd, "end", "", "End date (YYYY-MM-DD)")

	cmd.MarkFlagsMutuallyExclusive("date", "start")
	cmd.MarkFlagsMutuallyExclusive("date", "end")

	return cmd
}
