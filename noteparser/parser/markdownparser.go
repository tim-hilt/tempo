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
	entries := []string{}

	for _, e := range elems {
		e = strings.TrimSpace(e)
		if len(e) > 0 {
			entries = append(entries, e)
		}
	}

	return entries
}

func (m MarkdownParser) parseTicketEntries(ticketTable []string) ([]DailyNoteEntry, error) {
	ticketTable = ticketTable[2:]
	ticketEntries := []DailyNoteEntry{}

	for _, l := range ticketTable {
		entries := getLineEntries(l)

		durationSeconds, err := util.CalcDurationSeconds(entries[2])

		if err != nil {
			return nil, err
		}

		if durationSeconds == 0 {
			// Don't add to ticketEntries if no duration
			continue
		}

		ticketEntries = append(
			ticketEntries,
			DailyNoteEntry{
				Ticket:          entries[0],
				Comment:         entries[1],
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
