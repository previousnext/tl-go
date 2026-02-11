package fetch

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/internal/service"
)

func NewCommand(s func() service.SyncInterface) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "fetch",
		Args:                  cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Fetch issues from Jira",
		Long:                  "Fetch issue details from the Jira API and store them in the local database.",
		Example: `  # Fetch issues from Jira
  tl fetch PNX-123 PNX-456`,
		Hidden: true,
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
