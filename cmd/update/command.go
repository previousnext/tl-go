package update

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/model"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "update",
		Args:                  cobra.MinimumNArgs(2),
		DisableFlagsInUseLine: true,
		Short:                 "Update a time entry",
		Long:                  "Update a time entry by its ID.",
		Example: `  # Update a time entry with ID 123 to have a duration of 3h and a new description
  tl update 123 3h "Updated description"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			r := db.NewRepository(viper.GetString("db_file"))
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			entry, err := r.FindTimeEntry(uint(id))
			if err != nil {
				return err
			}
			if entry == nil {
				fmt.Printf("No entry with ID %d\n", id)
				return nil
			}

			dur, err := model.ParseDuration(args[1])
			if err != nil {
				return fmt.Errorf("invalid duration: %s", args[1])
			}
			entry.Duration = dur
			if len(args) > 2 {
				entry.Description = args[2]
			}

			err = r.UpdateTimeEntry(entry)
			if err != nil {
				return err
			}

			fmt.Printf("Updated time entry with ID %d\n", id)

			return nil
		},
	}
	return cmd
}
