package parser

import (
	"errors"
	"strings"

	"github.com/tim-hilt/tempo/util"
	"github.com/tim-hilt/tempo/util/config"
	"github.com/tim-hilt/tempo/util/set"
)

type MarkdownParser struct{}

func (m MarkdownParser) findTicketTable(file []byte) ([]string, error) {
	lines := strings.Split(string(file), "\n")
	tableStart := -1

	for i := 0; i < len(lines); i++ {
		i, err := findTableStart(lines[i:])

		if err != nil {
			return nil, err
		}

		if isTicketTable(lines[i]) {
			tableStart = i
			break
		}
	}

	if tableStart == -1 {
		return nil, errors.New("ticket table not found")
	}

	tableEnd := findTableEnd(lines, tableStart)

	return lines[tableStart:tableEnd], nil
}

func findTableStart(lines []string) (int, error) {
	for i, l := range lines {
		if isTableLine(l) {
			return i, nil
		}
	}

	return -1, errors.New("no table found")
}

func isTableLine(line string) bool {
	es := strings.Split(line, "|")
	return len(es) > 3
}

func isTicketTable(line string) bool {
	elems := strings.Split(line, "|")
	s := set.New[string]()

	for _, elem := range elems {
		e := strings.TrimSpace(elem)
		if len(e) > 0 {
			s.Add(e)
		}
	}

	cols := config.GetColumns()

	for _, c := range []string{cols.Tickets, cols.Comments, cols.Durations} {
		if !s.Contains(c) {
			return false
		}
	}

	return true
}

func findTableEnd(lines []string, start int) int {
	for l := start; l < len(lines); l++ {
		if !isTableLine(lines[l]) {
			return l
		}
	}
	return len(lines)
}

func getLineEntries(line string) []string {
	elems := strings.Split(line, "|")
	elems = elems[1 : len(elems)-1]

	for i, e := range elems {
		elems[i] = strings.TrimSpace(e)
	}

	return elems
}

type Columns struct {
	TicketColumn   int
	CommentColumn  int
	DurationColumn int
}

func getColumns(header string) (Columns, error) {
	headers := getLineEntries(header)
	cols := config.GetColumns()
	c := Columns{}

	for i, header := range headers {
		if header == cols.Tickets {
			c.TicketColumn = i
		} else if header == cols.Comments {
			c.CommentColumn = i
		} else if header == cols.Durations {
			c.DurationColumn = i
		} else {
			return Columns{}, errors.New("unexpected table-header: " + header)
		}
	}

	return c, nil
}

func (m MarkdownParser) parseTicketEntries(table []string) ([]DailyNoteEntry, error) {

	contentLines := table[2:]
	ticketEntries := []DailyNoteEntry{}

	c, err := getColumns(table[0])

	if err != nil {
		return nil, err
	}

	for _, l := range contentLines {
		entries := getLineEntries(l)

		durationSeconds, err := util.CalcDurationSeconds(entries[c.DurationColumn])

		if err != nil {
			return nil, err
		}

		if durationSeconds == 0 || entries[c.TicketColumn] == "" || entries[c.CommentColumn] == "" {
			// Don't add to ticketEntries if there are missing fields
			continue
		}

		ticketEntries = append(
			ticketEntries,
			DailyNoteEntry{
				Ticket:          entries[c.TicketColumn],
				Comment:         entries[c.CommentColumn],
				DurationSeconds: durationSeconds,
			},
		)
	}

	return ticketEntries, nil
}

// This function suffices to satisfy the parser-interface
func (m MarkdownParser) parseDailyNote(dailyNote []byte) ([]DailyNoteEntry, error) {
	ticketTable, err := m.findTicketTable(dailyNote)

	if err != nil {
		return nil, err
	}

	ticketEntries, err := m.parseTicketEntries(ticketTable)

	if err != nil {
		return nil, err
	}

	return ticketEntries, nil
}
