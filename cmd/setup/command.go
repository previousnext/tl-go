package setup

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/internal/db"
)

var (
	cmdLong    = `Initialize the tl configuration file and database.`
	cmdExample = `
  # Initialize tl
  tl init`
)

func NewCommand(r func() db.RepositoryInterface) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "init",
		Args:                  cobra.NoArgs,
		DisableFlagsInUseLine: true,
		Short:                 "Initialize tl",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := r().InitRepository()
			cobra.CheckErr(err)
			fmt.Fprintln(cmd.OutOrStdout(), "Successfully initialized repository")
			return nil
		},
	}

	return cmd
}
