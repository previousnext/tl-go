package pause

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/internal/service"
)

var (
	cmdUse   = "pause"
	cmdShort = "Pause the current time entry"
	cmdLong  = `Pause the current time entry.

This command will pause the timer for the current time entry, if one is in progress.`
)

func NewCommand(currentTimeStorage func() service.TimerEntryStorageInterface) *cobra.Command {
	return &cobra.Command{
		Use:   cmdUse,
		Short: cmdShort,
		Long:  cmdLong,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := currentTimeStorage().PauseTimeEntry()
			if err != nil {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "No current time entry to pause.\n")
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Current time entry has been paused.\n")
			return nil
		},
	}
}
