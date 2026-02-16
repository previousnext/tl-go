package alias

import (
	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/cmd/alias/delete"
	"github.com/previousnext/tl-go/cmd/alias/list"
	"github.com/previousnext/tl-go/cmd/alias/set"
)

var (
	cmdShort   = `Add, list, and delete command aliases.`
	cmdLong    = `Add, list, and delete command aliases.`
	cmdExample = `
  # Set an alias
  tl alias set foo PNX-123

  # Delete an alias
  tl alias delete foo

  # List all aliases
  tl alias list
`
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "alias [command]",
		DisableFlagsInUseLine: true,
		Short:                 cmdShort,
		Long:                  cmdLong,
		Example:               cmdExample,
	}

	cmd.AddCommand(set.NewCommand())
	cmd.AddCommand(delete.NewCommand())
	cmd.AddCommand(list.NewCommand())

	return cmd
}
