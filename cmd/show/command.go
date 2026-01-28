package show

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/model"
)

func NewCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:                   "show",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Show details of a time entry",
		Long:                  `Show details of a time entry by its ID.`,
		Example: `  # Show details of a time entry with ID 123
  tl show 123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			r := db.NewRepository(viper.GetString("db_file"))
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			entry, err := r.FindTimeEntry(uint(id))
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					fmt.Printf("No entry with ID %d\n", id)
					return nil
				}
				return err
			}

			fmt.Printf("Time Entry ID: %d\n", entry.ID)
			fmt.Printf("Issue Key:     %s\n", entry.IssueKey)
			fmt.Printf("Duration:      %s\n", model.FormatDuration(entry.Duration))
			fmt.Printf("Description:   %s\n", entry.Description)
			fmt.Printf("Created At:    %s\n", model.FormatDateTime(entry.CreatedAt))
			fmt.Printf("Updated At:    %s\n", model.FormatDateTime(entry.UpdatedAt))
			return nil
		},
	}
	return cmd
}
