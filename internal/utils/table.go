package utils

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

// Used for displaying tabular data in a user-friendly format in the CLI.
func RenderTable(headers []string, rows [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header(headers)
	for _, row := range rows {
		table.Append(row)
	}
	table.Render()
}
