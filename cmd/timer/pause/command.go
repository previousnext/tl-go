package pause

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/internal/service"
)

var (
	cmdUse   = "pause"
	cmdShort = "Pause the timer entry"
	cmdLong  = `Pause the timer entry.

This command will pause the timer for the timer entry, if one is in progress.`
)

func NewCommand(timerService func() service.TimerEntryServiceInterface) *cobra.Command {
	return &cobra.Command{
		Use:   cmdUse,
		Short: cmdShort,
		Long:  cmdLong,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := timerService().PauseTimeEntry()
			if err != nil {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "No timer entry to pause.\n")
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Timer entry has been paused.\n")
			return nil
		},
	}
}
