package issues

import (
	"bytes"
	"fmt"

	"github.com/aquasecurity/table"
	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/util"
)

var (
	cmdShort   = "List recent issues"
	cmdLong    = "List recent issues that have been used for time entries. This can be used to quickly find issue keys for adding new time entries."
	flagLimit  = 20
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

			if len(issues) == 0 {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), "No recent issues found")
				return nil
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Showing %d recent issues:\n", len(issues))

			var b bytes.Buffer
			t := table.New(&b)

			util.ApplyDefaultTableFormatting(t)

			headers := []string{
				"Key",
				"Summary",
				"Project",
				"Category",
			}

			var formattedHeaders []string
			for _, h := range headers {
				formattedHeaders = append(formattedHeaders, util.ApplyHeaderFormatting(h))
			}
			t.SetHeaders(formattedHeaders...)

			for _, issue := range issues {
				categoryName := ""
				if issue.Project.Category != nil {
					categoryName = issue.Project.Category.Name
				}
				t.AddRow(issue.Key,
					issue.Summary,
					issue.Project.Name,
					categoryName,
				)
			}

			t.Render()

			_, err = fmt.Fprintf(cmd.OutOrStdout(), "\n%s\n", b.String())
			if err != nil {
				return fmt.Errorf("failed to print table: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().IntVarP(&flagLimit, "limit", "l", flagLimit, "Maximum number of issues to fetch")

	return cmd
}
