package parser

import (
	"os"

	"github.com/rs/zerolog/log"
)

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

func DailyNoteEntriesEqual(a []DailyNoteEntry, b []DailyNoteEntry) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func GetTickets(p Parser, filePath string) ([]DailyNoteEntry, error) {
	dailyNote, err := os.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	log.Trace().Str("file", filePath).Msg("parsing daily note")

	return p.parseDailyNote(dailyNote)
}
