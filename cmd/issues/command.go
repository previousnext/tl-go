package issues

import (
	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/util"
)

var (
	cmdShort   = "List recent issues"
	cmdLong    = "List recent issues that have been used for time entries. This can be used to quickly find issue keys for adding new time entries."
	flagLimit  = 10
	cmdExample = `
  # Edit time entry with ID 1 to have a duration of 3 hours and a new description
  tl issues`
)

func NewCommand(r func() db.IssueStorageInterface) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "issues",
		Args:                  cobra.NoArgs,
		DisableFlagsInUseLine: true,
		Short:                 cmdShort,
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			issues, err := r().FindRecentIssues(flagLimit)
			if err != nil {
				return err
			}

			headers := []string{
				"Key",
				"Summary",
				"Project",
			}
			rows := make([][]string, len(issues))
			for i, issue := range issues {
				rows[i] = []string{
					issue.Key,
					issue.Summary,
					issue.Project.Name,
				}
			}

			return util.PrintTable(cmd.OutOrStdout(), headers, rows)
		},
	}

	cmd.Flags().IntVarP(&flagLimit, "limit", "l", flagLimit, "Maximum number of issues to fetch")

	return cmd
}
