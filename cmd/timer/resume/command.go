package resume

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/internal/service"
)

var (
	cmdUse   = "resume"
	cmdShort = "Resume the paused timer entry"
	cmdLong  = `Resume the paused timer entry.

This command will resume the timer entry that is currently paused.`
)

func NewCommand(timerService func() service.TimerEntryServiceInterface) *cobra.Command {
	return &cobra.Command{
		Use:   cmdUse + " [timer-id]",
		Short: cmdShort,
		Long:  cmdLong,
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
			err := timerService().ResumeTimerEntry(timerID)
			if err != nil {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\n", err.Error())
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Timer entry has been resumed.\n")
			return nil
		},
	}
}
