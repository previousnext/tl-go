package stop

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/model"
	"github.com/previousnext/tl-go/internal/service"
)

func NewCommand(currentTimeStorage func() db.CurrentTimeEntryStorageInterface, syncService func() service.SyncInterface) *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stop tracking time and save entry",
		RunE: func(cmd *cobra.Command, args []string) error {
			entry, err := currentTimeStorage().StopTimeEntry()
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
