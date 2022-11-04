package noteparser

import (
	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/util"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"
)

func applyOnChildren(parent ast.Node, kind string, fun func(child ast.Node)) {
	child := parent.FirstChild()
	for child != nil {
		if child.Kind().String() == kind {
			fun(child)
		}
		child = child.NextSibling()
	}
}

// TODO: Could implement interface that has to be satisfied in order to
// add parsers for more file-formats
func findTicketTable(file []byte) ast.Node {
	document := goldmark.New(goldmark.WithExtensions(extension.Table)).
		Parser().Parse(text.NewReader(file))

	var ticketTable ast.Node = nil
	applyOnChildren(document, "Table", func(node ast.Node) {
		if isTicketTable(node, file) {
			ticketTable = node
		}
	})

	if ticketTable == nil {
		log.Fatal().Msg("ticket table not found")
	}
	return ticketTable
}

func getTableHeaders(table ast.Node, file []byte) []string {
	headers := []string{}
	applyOnChildren(table, "TableHeader", func(tableRow ast.Node) {
		applyOnChildren(tableRow, "TableCell", func(tableCell ast.Node) {
			headers = append(headers, string(tableCell.Text(file)))
		})
	})
	return headers
}

func isTicketTable(table ast.Node, file []byte) bool {
	headers := getTableHeaders(table, file)
	return util.SlicesEqual([]string{"Ticket", "Doings", "Time spent"}, headers)
}

func parseTicketEntries(ticketTable ast.Node, file []byte) []DailyNoteEntry {
	ticketEntries := []DailyNoteEntry{}

	applyOnChildren(ticketTable, "TableRow", func(tableRow ast.Node) {
		rowVals := []string{}
		applyOnChildren(tableRow, "TableCell", func(tableCell ast.Node) {
			rowVals = append(rowVals, string(tableCell.Text(file)))
		})
		ticketEntries = append(ticketEntries, DailyNoteEntry{Ticket: rowVals[0], Comment: rowVals[1], DurationMinutes: calcDurationMinutes(rowVals[2])})
	})

	return ticketEntries
}

func getTickets(dailyNote []byte) []DailyNoteEntry {
	ticketTable := findTicketTable(dailyNote)
	ticketEntries := parseTicketEntries(ticketTable, dailyNote)
	return ticketEntries
}
