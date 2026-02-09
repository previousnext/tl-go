package update

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/internal/db"
)

func NewCommand(r func() db.TimeEntriesInterface) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "update",
		Args:                  cobra.MinimumNArgs(2),
		DisableFlagsInUseLine: true,
		Short:                 "Update a time entry",
		Long:                  "Update a time entry by its ID.",
		Example: `  # Update a time entry with ID 123 to have a duration of 3h and a new description
  tl update 123 3h "Updated description"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("entry ID should be a positive integer, got: %s", args[0])
			}
			repository := r()
			entry, err := repository.FindTimeEntry(uint(id))
			if err != nil {
				return err
			}
			if entry == nil {
				return fmt.Errorf("entry not found with ID: %d", id)
			}

			dur, err := time.ParseDuration(args[1])
			if err != nil {
				return fmt.Errorf("invalid duration: %s", args[1])
			}
			entry.Duration = dur
			if len(args) > 2 {
				entry.Description = args[2]
			}

			err = repository.UpdateTimeEntry(entry)
			if err != nil {
				return fmt.Errorf("update failed: %w", err)
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Updated time entry with ID %d\n", id)

			return nil
		},
	}
	return cmd
}
