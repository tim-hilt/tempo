package noteparser

import (
	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/util/config"
	"github.com/tim-hilt/tempo/util/set"
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

func getTableHeaders(table ast.Node, file []byte) set.Set[string] {
	headers := set.New[string]()
	applyOnChildren(table, "TableHeader", func(tableRow ast.Node) error {
		applyOnChildren(tableRow, "TableCell", func(tableCell ast.Node) error {
			headers.Add(string(tableCell.Text(file)))
			return nil
		})
		return nil
	})
	return headers
}

func isTicketTable(table ast.Node, file []byte) bool {
	headers := getTableHeaders(table, file)
	columns := config.GetColumns()
	for _, column := range []string{columns.Tickets, columns.Comments, columns.Durations} {
		if !headers.Contains(column) {
			return false
		}
	}
	return true
}

func parseTicketEntries(ticketTable ast.Node, file []byte) ([]DailyNoteEntry, error) {
	ticketEntries := []DailyNoteEntry{}

	err := applyOnChildren(ticketTable, "TableRow", func(tableRow ast.Node) error {
		rowVals := []string{}
		applyOnChildren(tableRow, "TableCell", func(tableCell ast.Node) error {
			rowVals = append(rowVals, string(tableCell.Text(file)))
			return nil
		})
		durationMinutes, err := calcDurationMinutes(rowVals[2])

		if durationMinutes == 0 {
			// Don't add to ticketEntries if no duration
			return nil
		}

		if err != nil {
			return err
		}

		ticketEntries = append(
			ticketEntries,
			DailyNoteEntry{
				Ticket:          rowVals[0],
				Comment:         rowVals[1],
				DurationMinutes: durationMinutes,
			},
		)
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
