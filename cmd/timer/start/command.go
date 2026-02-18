package start

import (
	"github.com/previousnext/tl-go/internal/alias"
	"github.com/previousnext/tl-go/internal/db"
	"github.com/previousnext/tl-go/internal/util"
	"github.com/spf13/cobra"
)

func NewCommand(currentTimeStorage func() db.CurrentTimeEntryStorageInterface, issueStorage func() db.IssueStorageInterface) *cobra.Command {
	return &cobra.Command{
		Use:   "start [issue-key]",
		Short: "Start tracking time for an issue",
		Args:  cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				completions, _ := util.CompleteAliasesAndIssueKeys(issueStorage())
				return completions, cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			issueKey := alias.ResolveAlias(args[0])
			return currentTimeStorage().StartTimeEntry(issueKey)
		},
	}
}
