package set

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/internal/alias"
)

var (
	cmdShort   = `Set an issue alias.`
	cmdLong    = `Set an issue alias.`
	cmdExample = `
  # Add an alias
  tl alias set foo PNX-123
`
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "set <alias> [issue-key]",
		Args:                  cobra.ExactArgs(2),
		DisableFlagsInUseLine: true,
		Short:                 cmdShort,
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			aliasName := strings.TrimSpace(args[0])
			key := strings.TrimSpace(args[1])
			storage := alias.NewAliasStorage()
			err := storage.SetAlias(aliasName, key)
			if err != nil {
				return fmt.Errorf("failed to set alias: %v", err)
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Added alias '%s' for value '%s'\n", aliasName, key)
			return nil
		},
	}
	return cmd
}
