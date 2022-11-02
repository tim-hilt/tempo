package tablecomponent

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/tim-hilt/tempo/util"
)

func getMaxChars(rows []table.Row, index int) int {
	max := 0
	for _, entry := range rows {
		if len(entry[index]) > max {
			max = len(entry[index])
		}
	}
	return max
}

func CreateColumns(rows []table.Row, columnTitles []string) []table.Column {
	columns := []table.Column{}
	for i, columnTitle := range columnTitles {
		columns = append(columns, table.Column{Title: columnTitle, Width: util.Max(getMaxChars(rows, i), len(columnTitle))})
	}
	return columns
}
