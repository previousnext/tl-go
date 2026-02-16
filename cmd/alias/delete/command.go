package delete

import (
	"fmt"

	"github.com/previousnext/tl-go/internal/alias"
	"github.com/spf13/cobra"
)

var (
	cmdShort   = `Delete an alias.`
	cmdLong    = `Delete an issue alias.`
	cmdExample = `
  # Delete an alias
  tl alias delete foo
`
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "delete <alias>",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 cmdShort,
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			storage := alias.NewAliasStorage()
			aliases, err := storage.LoadAliases()
			if err != nil {
				return err
			}
			aliasName := args[0]
			if _, exists := aliases[aliasName]; !exists {
				return fmt.Errorf("alias %s does not exist", aliasName)
			}
			delete(aliases, aliasName)
			err = storage.SaveAliases(aliases)
			if err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Deleted alias '%s'\n", aliasName)
			return nil
		},
	}
	return cmd
}
