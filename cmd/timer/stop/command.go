package stop

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/internal/model"
	"github.com/previousnext/tl-go/internal/service"
)

func NewCommand(timerService func() service.TimerEntryServiceInterface) *cobra.Command {
	var aiTimeSavedStr string

	cmd := &cobra.Command{
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

			var stopOpts []service.StopOptions
			if aiTimeSavedStr != "" {
				aiDur, err := time.ParseDuration(aiTimeSavedStr)
				if err != nil {
					return fmt.Errorf("invalid AI time saved duration: %s", aiTimeSavedStr)
				}
				stopOpts = append(stopOpts, service.StopOptions{AISavedDuration: aiDur})
			}

			entry, err := timerService().StopTimeEntry(timerID, stopOpts...)
			if err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Stopped time entry for %s, started at %s, duration: %s\n", entry.IssueKey, entry.CreatedAt.Local().Format("2006-01-02 15:04:05"), model.FormatDuration(entry.Duration))
			return nil
		},
	}

	cmd.Flags().StringVarP(&aiTimeSavedStr, "ai-time-saved", "a", "", "Duration of time saved by AI (e.g. 1h, 30m)")
	cmd.Flags().StringVar(&aiTimeSavedStr, "aits", "", "Duration of time saved by AI (shorthand for --ai-time-saved)")
	_ = cmd.Flags().MarkHidden("aits")

	return cmd
}
