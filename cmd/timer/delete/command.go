package delete

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/internal/service"
)

var (
	cmdShort   = "Delete a timer entry without saving it"
	cmdLong    = `Delete a timer entry by its ID without saving it as a time entry.

This lets you remove a timer directly from the timer list without having to
stop it first.`
	cmdExample = `
  # Delete the timer entry with ID 3
  tl timer delete 3

  # Delete without a confirmation prompt
  tl timer delete 3 --force`
)

func NewCommand(timerService func() service.TimerEntryServiceInterface) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:                   "delete <timer-id>",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 cmdShort,
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			parsed, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("timer ID must be a positive integer: %s", args[0])
			}
			id := uint(parsed)

			entry, err := timerService().GetTimerEntryByID(id)
			if err != nil {
				return fmt.Errorf("timer entry not found: %d", id)
			}

			if !force {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Delete timer entry %d for %s (started %s)? [y/N]: ",
					entry.ID, entry.IssueKey, entry.StartTime.Local().Format("2006-01-02 15:04:05"))

				reader := bufio.NewReader(cmd.InOrStdin())
				response, _ := reader.ReadString('\n')
				response = strings.ToLower(strings.TrimSpace(response))
				if response != "y" && response != "yes" {
					_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Aborted.")
					return nil
				}
			}

			if _, err := timerService().DeleteTimerEntry(id); err != nil {
				return err
			}
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Timer entry has been deleted.")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Delete without a confirmation prompt")

	return cmd
}
