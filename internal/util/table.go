package util

import (
	"bytes"
	"fmt"
	"io"

	"github.com/aquasecurity/table"
	"github.com/jwalton/gchalk"
)

func PrepareTable(t *table.Table) {
	t.SetHeaderStyle(table.StyleBold)
	t.SetLineStyle(table.StyleBrightBlack)
	t.SetDividers(table.UnicodeRoundedDividers)

	t.SetAvailableWidth(80)
	t.SetColumnMaxWidth(80)
}

// PrintTable a table to the console.
func PrintTable(w io.Writer, headers []string, rows [][]string, footer []string) error {
	var b bytes.Buffer

	t := table.New(&b)

	PrepareTable(t)
	
	var formattedHeaders []string

	for _, h := range headers {
		formattedHeaders = append(formattedHeaders, gchalk.WithHex(HexOrange).Bold(h))
	}

	t.SetHeaders(formattedHeaders...)

	for _, row := range rows {
		t.AddRow(row...)
	}

	t.SetFooters(footer...)

	t.Render()

	_, err := fmt.Fprintf(w, "\n%s\n", b.String())
	if err != nil {
		return fmt.Errorf("failed to print table: %w", err)
	}

	return nil
}

func FormatBool(b bool) string {
	if b {
		return "✅"
	}
	return "❌"
}
