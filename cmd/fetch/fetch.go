package fetch

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/internal/service"
)

var (
	cmdShort   = "Fetch issues from Jira."
	cmdLong    = "Fetch issue details from the Jira API and store them in the local database. This command is useful for ensuring that you have the latest issue information available locally."
	cmdExample = `  # Fetch issues from Jira
  tl fetch PNX-123 PNX-456`
)

func NewCommand(s func() service.SyncInterface) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "fetch",
		Args:                  cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 cmdShort,
		Long:                  cmdLong,
		Example:               cmdExample,
		Hidden:                true,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := s().SyncIssues(args)
			if err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Fetched details for %d issues from Jira\n", len(args))
			return nil
		},
	}

	return cmd
}
