package setup

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/previousnext/tl-go/internal/db"
)

var (
	cmdLong    = `Initialize the tl configuration file and database.`
	cmdExample = `
  # Initialize tl
  tl init`
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "init",
		Args:                  cobra.NoArgs,
		DisableFlagsInUseLine: true,
		Short:                 "Initialize tl",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Initialize the db
			r := db.NewRepository(viper.GetString("db_file"))
			err := r.InitRepository()
			cobra.CheckErr(err)
			return nil
		},
	}

	return cmd
}
