package parser

// Parser implements an interface that has to be satisfied in order to
// implement the note-parsing-capabilities that this program needs
type Parser interface {
	parseDailyNote(string) ([]DailyNoteEntry, error)
}

type DailyNoteEntry struct {
	Ticket          string
	Comment         string
	DurationMinutes int
}

func GetTickets(p Parser, day string) ([]DailyNoteEntry, error) {
	return p.parseDailyNote(day)
}
