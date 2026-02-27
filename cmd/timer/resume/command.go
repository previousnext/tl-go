package resume

import (
	"fmt"

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
		Use:   cmdUse,
		Short: cmdShort,
		Long:  cmdLong,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := timerService().ResumeTimerEntry()
			if err != nil {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\n", err.Error())
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Timer entry has been resumed.\n")
			return nil
		},
	}
}
