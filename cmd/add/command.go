package add

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/internal/alias"
	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/model"
	"github.com/previousnext/tl-go/internal/service"
)

var (
	cmdLong    = `Add a time entry`
	cmdExample = `
  # Add 2 hours to a project a project with issue ID PNX-123
  tl add PNX-123 2h "Worked on feature X"`
	date time.Time
)

func NewCommand(r func() db.TimeEntriesInterface, s func() service.SyncInterface) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "add <key> <time> [description] [flags]",
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

			// Check if this is an alias for an issue key and resolve it if necessary
			aliasStorage := alias.NewAliasStorage()
			aliases, err := aliasStorage.LoadAliases()
			if err != nil {
				return fmt.Errorf("error loading aliases: %w", err)
			}

			if resolvedKey, ok := aliases[key]; ok {
				key = resolvedKey
			}

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
			entry.CreatedAt = date

			err = r().CreateTimeEntry(entry)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Added time entry: ID=%d,", entry.ID)

			return nil
		},
	}

	timeFormats := []string{
		time.DateOnly,
	}
	cmd.Flags().TimeVarP(&date, "date", "d", time.Now(), timeFormats, "Date for the entry.")

	return cmd
}
