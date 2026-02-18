package continuecmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/internal/db"
)

func NewCommand(currentTimeStorage func() db.CurrentTimeEntryStorageInterface) *cobra.Command {
	return &cobra.Command{
		Use:   "continue",
		Short: "Continue (un-pause) the current time entry",
		RunE: func(cmd *cobra.Command, args []string) error {
			storage := currentTimeStorage()
			entry, err := storage.GetCurrentTimeEntry()
			if err != nil || entry == nil {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "No paused time entry to continue.\n")
				return err
			}
			if !entry.Paused {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Current time entry is not paused.\n")
				return nil
			}
			// Ensure minimum time period is 1 minute
			pausedDuration := entry.PauseTime.Sub(entry.StartTime)
			if pausedDuration < time.Minute {
				pausedDuration = time.Minute
			}
			entry.Duration += pausedDuration
			entry.Paused = false
			entry.StartTime = time.Now()
			entry.PauseTime = time.Time{}
			if err := storage.SaveCurrentTimeEntry(entry); err != nil {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Failed to continue time entry: %v\n", err)
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Current time entry has been continued.\n")
			return nil
		},
	}
}
