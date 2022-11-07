package noteparser

import (
	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/util"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"
)

type applyFunc func(child ast.Node) error

func applyOnChildren(parent ast.Node, kind string, fun applyFunc) error {
	child := parent.FirstChild()
	for child != nil {
		if child.Kind().String() == kind {
			err := fun(child)
			if err != nil {
				return err
			}
		}
		child = child.NextSibling()
	}
	return nil
}

// TODO: Could implement interface that has to be satisfied in order to
// add parsers for more file-formats
func findTicketTable(file []byte) ast.Node {
	document := goldmark.New(goldmark.WithExtensions(extension.Table)).
		Parser().Parse(text.NewReader(file))

	var ticketTable ast.Node = nil
	applyOnChildren(document, "Table", func(node ast.Node) error {
		if isTicketTable(node, file) {
			ticketTable = node
		}
		return nil
	})

	if ticketTable == nil {
		log.Fatal().Msg("ticket table not found")
	}
	return ticketTable
}

func getTableHeaders(table ast.Node, file []byte) []string {
	headers := []string{}
	applyOnChildren(table, "TableHeader", func(tableRow ast.Node) error {
		applyOnChildren(tableRow, "TableCell", func(tableCell ast.Node) error {
			headers = append(headers, string(tableCell.Text(file)))
			return nil
		})
		return nil
	})
	return headers
}

func isTicketTable(table ast.Node, file []byte) bool {
	headers := getTableHeaders(table, file)
	return util.SlicesEqual([]string{"Ticket", "Doings", "Time spent"}, headers)
}

func parseTicketEntries(ticketTable ast.Node, file []byte) ([]DailyNoteEntry, error) {
	ticketEntries := []DailyNoteEntry{}

	err := applyOnChildren(ticketTable, "TableRow", func(tableRow ast.Node) error {
		rowVals := []string{}
		applyOnChildren(tableRow, "TableCell", func(tableCell ast.Node) error {
			rowVals = append(rowVals, string(tableCell.Text(file)))
			return nil
		})
		durationMinutes, err := calcDurationMinutes(rowVals[3])

		if err != nil {
			return err
		}

		ticketEntries = append(ticketEntries, DailyNoteEntry{Ticket: rowVals[0], Comment: rowVals[1], DurationMinutes: durationMinutes})
		return nil
	})

	if err != nil {
		return nil, err
	}

	return ticketEntries, nil
}

func getTickets(dailyNote []byte) ([]DailyNoteEntry, error) {
	ticketTable := findTicketTable(dailyNote)
	ticketEntries, err := parseTicketEntries(ticketTable, dailyNote)
	if err != nil {
		return nil, err
	}
	return ticketEntries, nil
}
