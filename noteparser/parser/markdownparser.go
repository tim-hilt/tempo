package parser

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/util"
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

type MarkdownParser struct{}

func (m MarkdownParser) findTicketTable(file []byte) ast.Node {
	document := goldmark.New(goldmark.WithExtensions(extension.Table)).
		Parser().Parse(text.NewReader(file))

	var ticketTable ast.Node = nil
	applyOnChildren(document, "Table", func(node ast.Node) error {
		if m.isTicketTable(node, file) {
			ticketTable = node
		}
		return nil
	})

	if ticketTable == nil {
		log.Fatal().Msg("ticket table not found")
	}
	return ticketTable
}

func (m MarkdownParser) getTableHeaders(table ast.Node, file []byte) set.Set[string] {
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

func (m MarkdownParser) isTicketTable(table ast.Node, file []byte) bool {
	headers := m.getTableHeaders(table, file)
	columns := config.GetColumns()
	for _, column := range []string{columns.Tickets, columns.Comments, columns.Durations} {
		if !headers.Contains(column) {
			return false
		}
	}
	return true
}

func (m MarkdownParser) parseTicketEntries(
	ticketTable ast.Node,
	file []byte,
) ([]DailyNoteEntry, error) {
	ticketEntries := []DailyNoteEntry{}

	err := applyOnChildren(ticketTable, "TableRow", func(tableRow ast.Node) error {
		rowVals := []string{}
		applyOnChildren(tableRow, "TableCell", func(tableCell ast.Node) error {
			rowVals = append(rowVals, string(tableCell.Text(file)))
			return nil
		})
		durationMinutes, err := util.CalcDurationMinutes(rowVals[2])

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

func (m MarkdownParser) getDailyNote(day string) ([]byte, error) {
	// TODO: How can I find out the suffix at runtime and choose
	//       the parser accordingly?
	if !strings.HasSuffix(day, ".md") {
		day = day + ".md"
	}
	notesDir := config.GetNotesdir()
	fileWithPath := filepath.Join(notesDir, day)

	file, err := os.ReadFile(fileWithPath)

	if err != nil {
		return nil, err
	}

	return file, nil
}

// This function suffices to satisfy the parser-interface
func (m MarkdownParser) parseDailyNote(day string) ([]DailyNoteEntry, error) {
	dailyNote, err := m.getDailyNote(day)

	if err != nil {
		return nil, err
	}
	ticketTable := m.findTicketTable(dailyNote)
	ticketEntries, err := m.parseTicketEntries(ticketTable, dailyNote)
	if err != nil {
		return nil, err
	}
	return ticketEntries, nil

}
