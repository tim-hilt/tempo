package noteparser

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/tim-hilt/tempo/util"
)

type DailyNoteEntry struct {
	Ticket          string
	Comment         string
	DurationMinutes int
}

func getDailyNote(day string) ([]byte, error) {
	fileName := day + ".md"
	notesDir := util.GetConfigParams().Notesdir
	fileWithPath := filepath.Join(notesDir, fileName)

	file, err := os.ReadFile(fileWithPath)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func calcDurationMinutes(duration string) (int, error) {
	foo := strings.Split(duration, ":")
	hours, err := strconv.Atoi(foo[0])

	if err != nil {
		return -1, err
	}

	minutes, err := strconv.Atoi(foo[1])

	if err != nil {
		return -1, err
	}

	return hours*60 + minutes, nil
}

func ParseDailyNote(day string) ([]DailyNoteEntry, error) {
	dailyNote, err := getDailyNote(day)

	if err != nil {
		return nil, err
	}

	ticketEntries, err := getTickets(dailyNote)

	if err != nil {
		return nil, err
	}

	return ticketEntries, nil
}
