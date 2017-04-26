package ec2go

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

func newTableWriter() *tablewriter.Table {
	return tablewriter.NewWriter(os.Stdout)
}

func setTable(table *tablewriter.Table, header []string, body [][]string) *tablewriter.Table {
	table.SetHeader(header)

	for _, v := range body {
		table.Append(v)
	}
	return table
}

func displayTable(table *tablewriter.Table) {
	table.Render()
}
