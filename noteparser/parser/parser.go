package parser

import "os"

// Parser implements an interface that has to be satisfied in order to
// implement the note-parsing-capabilities that this program needs
type Parser interface {
	parseDailyNote([]byte) ([]DailyNoteEntry, error)
}

type DailyNoteEntry struct {
	Ticket          string
	Comment         string
	DurationSeconds int
}

func GetTickets(p Parser, filePath string) ([]DailyNoteEntry, error) {
	dailyNote, err := os.ReadFile(filePath)

	if err != nil {
		return nil, err
	}
	return p.parseDailyNote(dailyNote)
}
