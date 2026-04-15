package aits

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/model"
)

var (
	cmdLong    = `Set the AI time saved on a time entry.`
	cmdExample = `
  # Set 1 hour of AI time saved on time entry 42
  tl aits 42 1h

  # Set 30 minutes of AI time saved on time entry 7
  tl aits 7 30m`
)

func NewCommand(r func() db.TimeEntriesInterface) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "aits <id> <duration>",
		Args:                  cobra.ExactArgs(2),
		DisableFlagsInUseLine: true,
		Short:                 "Set AI time saved on a time entry",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid time entry ID: %s", args[0])
			}

			dur, err := time.ParseDuration(args[1])
			if err != nil {
				return fmt.Errorf("invalid duration: %s", args[1])
			}

			entryStorage := r()
			entry, err := entryStorage.FindTimeEntry(uint(id))
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					_, _ = fmt.Fprintf(cmd.OutOrStdout(), "No entry with ID %d\n", id)
					return nil
				}
				return err
			}

			entry.AISavedDuration = dur

			if err := entryStorage.UpdateTimeEntry(entry); err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Set AI time saved to %s on time entry ID %d\n", model.FormatDuration(dur), entry.ID)

			return nil
		},
	}
	return cmd
}
