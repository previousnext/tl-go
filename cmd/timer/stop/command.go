package stop

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/internal/model"
	"github.com/previousnext/tl-go/internal/service"
)

func NewCommand(timerService func() service.TimerEntryServiceInterface, syncService func() service.SyncInterface) *cobra.Command {
	return &cobra.Command{
		Use:   "stop [timer-id]",
		Short: "Stop tracking time and save entry",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var timerID *uint
			if len(args) == 1 {
				parsed, err := strconv.ParseUint(args[0], 10, 64)
				if err != nil {
					_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Invalid timer ID: %s\n", args[0])
					return err
				}
				id := uint(parsed)
				timerID = &id
			}
			entry, err := timerService().StopTimeEntry(timerID)
			if err != nil {
				return err
			}
			// Sync the issue after creating the time entry
			if entry != nil {
				_, _ = syncService().SyncIssue(entry.IssueKey)
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Stopped time entry for %s, started at %s, duration: %s\n", entry.IssueKey, entry.CreatedAt.Local().Format("2006-01-02 15:04:05"), model.FormatDuration(entry.Duration))
			return nil
		},
	}
}
