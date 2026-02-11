package edit

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"github.com/previousnext/tl-go/internal/db"
)

var (
	cmdLong    = `Edit a time entry`
	cmdExample = `
  # Edit time entry with ID 1 to have a duration of 3 hours and a new description
  tl edit 1 --duration 3h --description "Updated description"`
)

func NewCommand(r func() db.TimeEntriesInterface) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "edit <id> [flags]",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Edit a time entry",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get the time entry ID from the arguments
			i, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid time entry ID: %s", args[0])
			}
			entryID := uint(i)

			entryStorage := r()
			timeEntry, err := entryStorage.FindTimeEntry(entryID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					_, _ = fmt.Fprintf(cmd.OutOrStderr(), "No entry with ID %d\n", entryID)
					return nil
				}
				return err
			}

			// Get the new duration and description from the flags
			durStr, _ := cmd.Flags().GetString("dur")
			if durStr != "" {
				dur, err := time.ParseDuration(durStr)
				if err != nil {
					return fmt.Errorf("invalid duration: %s", durStr)
				}
				timeEntry.Duration = dur
			}
			desc, _ := cmd.Flags().GetString("desc")
			if desc != "" {
				timeEntry.Description = desc
			}
			// If no changes were specified, print an error message and return
			if durStr == "" && desc == "" {
				_, _ = fmt.Fprintln(cmd.OutOrStderr(), "No changes specified. Use --dur and/or --desc to specify changes.")
				return nil
			}

			// Update the time entry in the database
			err = entryStorage.UpdateTimeEntry(timeEntry)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Updated time entry ID %d\n", timeEntry.ID)

			return nil
		},
	}

	cmd.Flags().StringP("dur", "d", "", "New duration (e.g. 2h30m)")
	cmd.Flags().StringP("desc", "s", "", "New description")

	return cmd
}
