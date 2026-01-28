package create

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/model"
)

var (
	cmdLong    = `Add a time entry`
	cmdExample = `
  # Add 2 hours to a project a project with issue ID PNX-123
  tl add PNX-123 2h "Worked on feature X"`
)

func NewCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:                   "add <issue_key> <duration> [description]",
		Args:                  cobra.MinimumNArgs(2),
		DisableFlagsInUseLine: true,
		Short:                 "Add a time entry",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			dur, err := model.ParseDuration(args[1])
			if err != nil {
				return fmt.Errorf("invalid duration: %s", args[1])
			}
			entry := &model.TimeEntry{
				IssueKey: args[0],
				Duration: dur,
			}
			if len(args) > 2 {
				entry.Description = args[2]
			}

			r := db.NewRepository(viper.GetString("db_file"))
			err = r.CreateTimeEntry(entry)
			if err != nil {
				return err
			}

			fmt.Printf("Added time entry: ID=%d,", entry.ID)

			return nil
		},
	}

	return cmd
}
