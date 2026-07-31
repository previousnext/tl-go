package start

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/internal/alias"
	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/model"
	"github.com/previousnext/tl-go/internal/service"
	"github.com/previousnext/tl-go/internal/util"
)

func NewCommand(timerService func() service.TimerEntryServiceInterface, issueStorage func() db.IssueStorageInterface) *cobra.Command {
	return &cobra.Command{
		Use:   "start [issue-key] [description...]",
		Short: "Start tracking time for an issue",
		Args:  cobra.MinimumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				completions, _ := util.CompleteAliasesAndIssueKeys(issueStorage())
				return completions, cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			issueKey := alias.ResolveAlias(args[0])
			var description *string
			if len(args) > 1 {
				desc := strings.Join(args[1:], " ")
				description = &desc
			}
			prev, err := timerService().StartTimeEntry(issueKey, description)
			if err != nil {
				return err
			}
			if prev != nil {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Paused time entry for %s, duration: %s\n", prev.IssueKey, model.FormatDuration(prev.Duration))
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Started time entry for %s\n", issueKey)
			return nil
		},
	}
}
