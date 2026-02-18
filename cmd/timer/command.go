package timer

import (
	"github.com/spf13/cobra"

	continuecmd "github.com/previousnext/tl-go/cmd/timer/continue"
	pausecmd "github.com/previousnext/tl-go/cmd/timer/pause"
	showcmd "github.com/previousnext/tl-go/cmd/timer/show"
	startcmd "github.com/previousnext/tl-go/cmd/timer/start"
	stopcmd "github.com/previousnext/tl-go/cmd/timer/stop"
	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/service"
)

func NewCommand(currentTimeStorage func() db.CurrentTimeEntryStorageInterface, issueStorage func() db.IssueStorageInterface, syncService func() service.SyncInterface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "timer",
		Short: "Start, stop, pause, and show the current timer session.",
		Long:  "Start, stop, pause, and show the current timer session. This is used to track time spent on an issue in real-time.",
	}

	cmd.AddCommand(startcmd.NewCommand(currentTimeStorage, issueStorage))
	cmd.AddCommand(stopcmd.NewCommand(currentTimeStorage, syncService))
	cmd.AddCommand(pausecmd.NewCommand(currentTimeStorage))
	cmd.AddCommand(showcmd.NewCommand(currentTimeStorage))
	cmd.AddCommand(continuecmd.NewCommand(currentTimeStorage))
	return cmd
}
