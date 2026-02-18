package pause

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/internal/db"
)

func NewCommand(currentTimeStorage func() db.CurrentTimeEntryStorageInterface) *cobra.Command {
	return &cobra.Command{
		Use:   "pause",
		Short: "Pause the current time entry",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := currentTimeStorage().PauseTimeEntry()
			if err != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "No current time entry to pause.\n")
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Current time entry has been paused.\n")
			return nil
		},
	}
}
