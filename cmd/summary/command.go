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
	cmdLong    = "Show a summary of time spent per project category for the current week."
	cmdExample = `
  # Show a summary of time spent per project category for the current week
  tl summary`
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
			var err error
			if flagStart == "" {
				// Default to the start of the current week if no start date is provided
				start = startOfCurrentWeek()
			} else {
				start, err = time.ParseInLocation("2006-01-02", flagStart, time.Local)
				if err != nil {
					return err
				}
			}
			// Set start to midnight
			start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())

			if flagEnd == "" {
				// Default to 7 days after the start date if no end date is provided
				end = start.AddDate(0, 0, 7)
			} else {
				end, err = time.ParseInLocation("2006-01-02", flagEnd, time.Local)
				if err != nil {
					return err
				}
			}
			// Set end to end of day to include the entire end date
			end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, end.Location())

			// Validate start is before or equal to end
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
			util.ApplyTableFormatting(t)
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
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Summary of time spent per category from %s to %s:\n", start.Format("2006-01-02"), end.Format("2006-01-02"))
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\n", b.String())
			return nil
		},
	}

	cmd.Flags().StringVar(&flagStart, "start", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&flagEnd, "end", "", "End date (YYYY-MM-DD)")
	return cmd
}

func startOfCurrentWeek() time.Time {
	now := time.Now().In(time.Local)
	daysToMonday := (int(now.Weekday()) + 6) % 7
	return now.AddDate(0, 0, -daysToMonday)
}
