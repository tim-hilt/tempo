package noteparser

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/tim-hilt/tempo/noteparser/parser"
	"github.com/tim-hilt/tempo/util/config"
)

// TODO: Make sure that "day" is always just a date in the form 2022-11-23
func ParseDailyNote(day string) ([]parser.DailyNoteEntry, error) {
	notesDir := config.GetNotesdir()
	fileWithPath := filepath.Join(notesDir, day)

	matches, err := filepath.Glob(fileWithPath + "*")

	if err != nil {
		return nil, err
	}

	if len(matches) != 1 {
		return nil, errors.New("more than one daily note found")
	}

	dailyNote := matches[0]
	splitAtDots := strings.Split(dailyNote, ".")

	if len(splitAtDots) != 2 {
		return nil, errors.New("unsupported filename: " + dailyNote)
	}

	fileEnding := splitAtDots[1]

	if fileEnding == "md" {
		ticketEntries, err := parser.GetTickets(parser.MarkdownParser{}, dailyNote)
		if err != nil {
			return nil, err
		}
		return ticketEntries, nil
		// TODO: Add new file-formats here
	} else {
		return nil, errors.New("file-format \"" + fileEnding + "\" not supported")
	}
}
