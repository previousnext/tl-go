package show

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/model"
)

func NewCommand(r func() db.RepositoryInterface) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "show",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Show details of a time entry",
		Long:                  `Show details of a time entry by its ID.`,
		Example: `  # Show details of a time entry with ID 123
  tl show 123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			entry, err := r().FindTimeEntry(uint(id))
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					fmt.Fprintf(cmd.OutOrStdout(), "No entry with ID %d\n", id)
					return nil
				}
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Time Entry ID:\t%d\n", entry.ID)
			fmt.Fprintf(cmd.OutOrStdout(), "Issue Key:\t%s\n", entry.IssueKey)
			fmt.Fprintf(cmd.OutOrStdout(), "Duration:\t%s\n", model.FormatDuration(entry.Duration))
			fmt.Fprintf(cmd.OutOrStdout(), "Description:\t%s\n", entry.Description)
			fmt.Fprintf(cmd.OutOrStdout(), "Created At:\t%s\n", model.FormatDateTime(entry.CreatedAt))
			fmt.Fprintf(cmd.OutOrStdout(), "Updated At:\t%s\n", model.FormatDateTime(entry.UpdatedAt))
			return nil
		},
	}
	return cmd
}
