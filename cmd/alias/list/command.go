package list

import (
	"bytes"
	"fmt"

	"github.com/aquasecurity/table"
	"github.com/jwalton/gchalk"
	"github.com/spf13/cobra"

	"github.com/previousnext/tl-go/internal/alias"
	"github.com/previousnext/tl-go/internal/util"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "list",
		DisableFlagsInUseLine: true,
		Short:                 "List all aliases.",
		Long:                  "List all command aliases.",
		RunE: func(cmd *cobra.Command, args []string) error {
			storage := alias.NewAliasStorage()
			aliases, err := storage.LoadAliases()
			if err != nil {
				return err
			}

			if len(aliases) == 0 {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), "No aliases found.")
				return nil
			}

			headers := []string{
				"Alias",
				"Value",
			}
			var b bytes.Buffer
			t := table.New(&b)
			util.ApplyTableFormatting(t)

			var formattedHeaders []string
			for _, h := range headers {
				formattedHeaders = append(formattedHeaders, gchalk.WithHex(util.HexOrange).Bold(h))
			}
			t.AddHeaders(formattedHeaders...)

			for aliasName, value := range aliases {
				t.AddRow(aliasName, value)
			}

			t.Render()

			_, err = fmt.Fprintf(cmd.OutOrStdout(), "\n%s\n", b.String())
			if err != nil {
				return fmt.Errorf("failed to print table: %w", err)
			}

			return nil
		},
	}
	return cmd
}
