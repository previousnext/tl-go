package timer

import (
	"github.com/spf13/cobra"

	listcmd "github.com/previousnext/tl-go/cmd/timer/list"
	pausecmd "github.com/previousnext/tl-go/cmd/timer/pause"
	resumecmd "github.com/previousnext/tl-go/cmd/timer/resume"
	startcmd "github.com/previousnext/tl-go/cmd/timer/start"
	stopcmd "github.com/previousnext/tl-go/cmd/timer/stop"
	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/service"
)

func NewCommand(timerService func() service.TimerEntryServiceInterface, issueStorage func() db.IssueStorageInterface, syncService func() service.SyncInterface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "timer",
		Short: "Start, stop, pause, resume, and list timer entries.",
		Long:  "Start, stop, pause, resume, and list timer entries. This is used to track time spent on an issue in real-time.",
	}

	cmd.AddCommand(startcmd.NewCommand(timerService, issueStorage))
	cmd.AddCommand(stopcmd.NewCommand(timerService, syncService))
	cmd.AddCommand(pausecmd.NewCommand(timerService))
	cmd.AddCommand(resumecmd.NewCommand(timerService))
	cmd.AddCommand(listcmd.NewCommand(timerService))
	return cmd
}
