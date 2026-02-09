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

func NewCommand(r func() db.TimeEntriesInterface) *cobra.Command {
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
					_, _ = fmt.Fprintf(cmd.OutOrStdout(), "No entry with ID %d\n", id)
					return nil
				}
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "ID:\t\t%d\n", entry.ID)
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Key:\t\t%s\n", entry.IssueKey)
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Summary:\t%s\n", entry.Issue.Summary)
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Duration:\t%s\n", model.FormatDuration(entry.Duration))
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Description:\t%s\n", entry.Description)
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Created At:\t%s\n", model.FormatDateTime(entry.CreatedAt))
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Updated At:\t%s\n", model.FormatDateTime(entry.UpdatedAt))
			return nil
		},
	}
	return cmd
}
