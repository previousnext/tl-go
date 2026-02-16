package migrate

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/internal/db"
)

var (
	cmdShort   = "Migrate the database schema to the latest version."
	cmdLong    = "Migrate the database schema to the latest version. This command should be run after updating tl-go to ensure that the database schema is compatible with the new version."
	cmdExample = `  # Migrate the database schema
  tl migrate`
)

func NewCommand(d func() db.RepositoryInterface) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "migrate",
		Args:                  cobra.NoArgs,
		DisableFlagsInUseLine: true,
		Short:                 cmdShort,
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := d().AutoMigrate()
			if err != nil {
				return fmt.Errorf("error migrating database: %w", err)
			}
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Database schema migrated successfully.")
			return nil
		},
	}
	return cmd

}
