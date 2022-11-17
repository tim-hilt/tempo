package noteparser

import "github.com/tim-hilt/tempo/noteparser/parser"

func ParseDailyNote(day string) ([]parser.DailyNoteEntry, error) {
	// Change this if you want a different parser
	ticketEntries, err := parser.GetTickets(parser.MarkdownParser{}, day)

	if err != nil {
		return nil, err
	}

	return ticketEntries, nil
}
