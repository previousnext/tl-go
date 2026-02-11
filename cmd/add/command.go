package add

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/model"
	"github.com/previousnext/tl-go/internal/service"
)

var (
	cmdLong    = `Add a time entry`
	cmdExample = `
  # Add 2 hours to a project a project with issue ID PNX-123
  tl add PNX-123 2h "Worked on feature X"`
)

func NewCommand(r func() db.TimeEntriesInterface, s func() service.SyncInterface) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "add <key> <duration> [description]",
		Args:                  cobra.MinimumNArgs(2),
		DisableFlagsInUseLine: true,
		Short:                 "Add a time entry",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			dur, err := time.ParseDuration(args[1])
			if err != nil {
				return fmt.Errorf("invalid duration: %s", args[1])
			}
			key := args[0]

			issue, err := s().SyncIssue(key)
			if err != nil {
				return err
			}

			entry := &model.TimeEntry{
				IssueKey: key,
				Issue:    issue,
				Duration: dur,
			}
			if len(args) > 2 {
				entry.Description = args[2]
			}

			err = r().CreateTimeEntry(entry)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Added time entry: ID=%d,", entry.ID)

			return nil
		},
	}

	return cmd
}
