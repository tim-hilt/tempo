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

func getDailyNote(day string) []byte {
	fileName := day + ".md"
	notesDir := util.GetConfigParams().Notesdir
	fileWithPath := filepath.Join(notesDir, fileName)
	if strings.HasPrefix(fileWithPath, "~") {
		home, err := os.UserHomeDir()
		util.HandleErr(err, "error when searching for users homedir")
		fileWithPath = filepath.Join(home, fileWithPath[1:])
	}

	file, err := os.ReadFile(fileWithPath)
	util.HandleErr(err, "error when reading daily note "+fileWithPath)

	return file
}

func calcDurationMinutes(duration string) int {
	foo := strings.Split(duration, ":")
	hours, err := strconv.Atoi(foo[0])
	util.HandleErr(err, "error when converting hours-string in duration \""+duration+"\" to int")
	minutes, err := strconv.Atoi(foo[1])
	util.HandleErr(err, "error when converting minutes-string in duration \""+duration+"\" to int")
	return hours*60 + minutes
}

func ParseDailyNote(day string) []DailyNoteEntry {
	dailyNote := getDailyNote(day)
	ticketEntries := getTickets(dailyNote)
	return ticketEntries
}
