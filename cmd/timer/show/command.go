package show

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/model"
	"github.com/previousnext/tl-go/internal/util"
)

func NewCommand(currentTimeStorage func() db.CurrentTimeEntryStorageInterface) *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show the current in-progress time entry",
		RunE: func(cmd *cobra.Command, args []string) error {
			storage := currentTimeStorage()
			entry, err := storage.GetCurrentTimeEntry()
			if err != nil || entry == nil {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), "No current time entry in progress.")
				return nil
			}

			headers := []string{"Issue Key", "Started", "Duration", "In Progress"}
			var duration time.Duration
			var inProgress string
			if entry.Paused {
				duration = entry.PauseTime.Sub(entry.StartTime)
				inProgress = "No"
			} else {
				duration = time.Since(entry.StartTime)
				inProgress = "Yes"
			}
			// Format duration to minimum unit of minutes and round to nearest minute
			if duration < time.Minute {
				duration = time.Minute
			} else {
				duration = duration.Round(time.Minute)
			}
			rows := [][]string{
				{
					entry.IssueKey,
					entry.StartTime.Local().Format("2006-01-02 15:04:05"),
					model.FormatDuration(duration),
					inProgress,
				},
			}
			if err := util.PrintTable(cmd.OutOrStdout(), headers, rows, nil); err != nil {
				return err
			}
			return nil
		},
	}
}
