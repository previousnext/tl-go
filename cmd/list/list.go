package list

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/model"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "list",
		Args:                  cobra.NoArgs,
		DisableFlagsInUseLine: true,
		Short:                 "List all time entries",
		Long:                  "List all time entries in the database.",
		Example: `  # List all time entries
  tl list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			r := db.NewRepository(viper.GetString("db_file"))
			entries, err := r.FindAllTimeEntries()
			if err != nil {
				return err
			}
			if len(entries) == 0 {
				cmd.Println("No time entries found.")
				return nil
			}

			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 1, ' ', 0)
			fmt.Fprintln(w, "ID\tIssue Key\tDuration\tDescription")
			fmt.Fprintln(w, "--\t---------\t--------\t-----------")
			for _, entry := range entries {
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", entry.ID, entry.IssueKey, model.FormatDuration(entry.Duration), entry.Description)
			}
			w.Flush()
			return nil
		},
	}

	return cmd
}
