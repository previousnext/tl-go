package list

import (
	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/internal/service"
)

var (
	cmdShort   = `List all time entries`
	cmdLong    = `List all time entries in the database.`
	cmdExample = `
  # List all time entries
  tl list`
)

func NewCommand(currentTimeStorage func() service.TimerEntryStorageInterface) *cobra.Command {
	return &cobra.Command{
		Use:                   "list",
		Args:                  cobra.NoArgs,
		DisableFlagsInUseLine: true,
		Short:                 cmdShort,
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			currentTimeStorage().FindAllTimerEntries()
			return nil
		},
	}
}
