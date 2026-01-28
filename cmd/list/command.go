package list

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/model"
)

func NewCommand(r func() db.RepositoryInterface) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "list",
		Args:                  cobra.NoArgs,
		DisableFlagsInUseLine: true,
		Short:                 "List all time entries",
		Long:                  "List all time entries in the database.",
		Example: `  # List all time entries
  tl list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := r().FindAllTimeEntries()
			if err != nil {
				return err
			}
			if len(entries) == 0 {
				cmd.Println("No time entries found.")
				return nil
			}

			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 1, ' ', 0)
			_, _ = fmt.Fprintln(w, "ID\tIssue Key\tDuration\tDescription")
			_, _ = fmt.Fprintln(w, "--\t---------\t--------\t-----------")
			for _, entry := range entries {
				_, _ = fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", entry.ID, entry.IssueKey, model.FormatDuration(entry.Duration), entry.Description)
			}
			_ = w.Flush()
			return nil
		},
	}

	return cmd
}
