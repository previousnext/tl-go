package show

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/model"
	"github.com/previousnext/tl-go/internal/util"
)

func NewCommand(r func() db.TimeEntriesInterface) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "show <id>",
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

			headers := []string{
				"Field",
				"Value",
			}
			rows := [][]string{
				{"ID", fmt.Sprintf("%d", entry.ID)},
				{"Key", entry.IssueKey},
				{"Summary", entry.Issue.Summary},
				{"Project", entry.Issue.Project.Name},
				{"Duration", model.FormatDuration(entry.Duration)},
				{"Description", entry.Description},
				{"Created At", model.FormatDateTime(entry.CreatedAt)},
				{"Updated At", model.FormatDateTime(entry.UpdatedAt)},
				{"Sent To Jira", util.FormatBool(entry.Sent)},
			}

			var footer []string

			return util.PrintTable(cmd.OutOrStdout(), headers, rows, footer)
		},
	}
	return cmd
}
