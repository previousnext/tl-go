package manpages

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var flagsOutputDir = "out-dir"

func NewCommand(rootCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "manpages",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			outDir, err := cmd.Flags().GetString(flagsOutputDir)
			if err != nil {
				return err
			}
			if err := os.MkdirAll(outDir, 0755); err != nil {
				return err
			}
			head := &doc.GenManHeader{
				Title:   "TL",
				Section: "1",
			}
			return doc.GenManTree(rootCmd, head, filepath.Clean(outDir))
		},
	}
	cmd.Flags().StringP(flagsOutputDir, "o", "manpages", "Output directory for manpages")
	return cmd
}
