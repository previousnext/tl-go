package delete

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/previousnext/tl-go/internal/db"
)

func NewCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:                   "delete <time_entry_id>",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Delete a time entry",
		Long:                  "Delete a time entry by its ID.",
		Example: `  # Delete a time entry with ID 123
  tl delete 123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			r := db.NewRepository(viper.GetString("db_file"))
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			err = r.DeleteTimeEntry(uint(id))
			if err != nil {
				return err
			}
			fmt.Printf("Deleted time entry with ID %d\n", id)
			return nil
		},
	}
	return cmd
}
