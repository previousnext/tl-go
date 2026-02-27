package list

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/internal/model"
	"github.com/previousnext/tl-go/internal/service"
	"github.com/previousnext/tl-go/internal/util"
)

var (
	cmdShort   = `List all timer entries`
	cmdLong    = `List all timer entries in the database.`
	cmdExample = `
  # List all timer entries
  tl timer list`
)

func NewCommand(timerService func() service.TimerEntryServiceInterface) *cobra.Command {
	return &cobra.Command{
		Use:                   "list",
		Args:                  cobra.NoArgs,
		DisableFlagsInUseLine: true,
		Short:                 cmdShort,
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := timerService().FindAllTimerEntries()
			if err != nil {
				return err
			}
			if len(entries) == 0 {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), "No timer entries found.")
				return nil
			}

			headers := []string{
				"ID",
				"Issue",
				"Started",
				"Duration",
				"Paused",
				"Description",
			}

			var rows [][]string
			totalDuration := time.Duration(0)

			for _, entry := range entries {
				dur := entry.Duration
				if !entry.Paused {
					lastActive := entry.LastActiveTime
					if lastActive.IsZero() {
						lastActive = entry.StartTime
					}
					dur += time.Since(lastActive)
				}

				description := ""
				if entry.Description != nil {
					description = *entry.Description
				}

				rows = append(rows, []string{
					fmt.Sprintf("%d", entry.ID),
					entry.IssueKey,
					entry.StartTime.Local().Format("2006-01-02 15:04:05"),
					model.FormatDuration(dur),
					fmt.Sprintf("%t", entry.Paused),
					description,
				})
				totalDuration += dur
			}

			footer := []string{
				"",
				"",
				util.ApplyHeaderFormatting("Total"),
				util.ApplyHeaderFormatting(model.FormatDuration(totalDuration)),
				"",
				"",
			}

			if err := util.PrintTable(cmd.OutOrStdout(), headers, rows, footer); err != nil {
				return fmt.Errorf("error printing table: %w", err)
			}
			return nil
		},
	}
}
