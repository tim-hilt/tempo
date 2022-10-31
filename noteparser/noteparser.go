package noteparser

import (
	"os"
	"strconv"
	"strings"

	"github.com/tim-hilt/tempo/util"
)

type DailyNoteEntry struct {
	Ticket          string
	Comment         string
	DurationMinutes int
}

func getDailyNote(day string) []string {
	noteName := "/home/tim/notes/daily/" + day + ".md"
	file, err := os.ReadFile(noteName)
	util.HandleErr(err, "error when reading daily note "+noteName)

	fileLines := strings.Split(string(file), "\n")
	return fileLines
}

func findTicketTable(lines []string) []string {
	beginTicketTable := -1
	endTicketTable := -1

	for i, line := range lines {
		if strings.HasPrefix(line, "|") && beginTicketTable == -1 {
			beginTicketTable = i + 2
		} else if !strings.HasPrefix(line, "|") && beginTicketTable != -1 {
			endTicketTable = i
			break
		}
	}

	ticketTable := lines[beginTicketTable:endTicketTable]
	return ticketTable
}

func calcDurationMinutes(duration string) int {
	foo := strings.Split(duration, ":")
	hours, err := strconv.Atoi(foo[0])
	util.HandleErr(err, "error when converting hours-string in duration \""+duration+"\" to int")
	minutes, err := strconv.Atoi(foo[1])
	util.HandleErr(err, "error when converting minutes-string in duration \""+duration+"\" to int")
	return hours*60 + minutes
}

func parseTicketEntries(ticketTable []string) []DailyNoteEntry {
	ticketEntries := []DailyNoteEntry{}

	for _, entry := range ticketTable {
		worklogPieces := strings.Split(entry, "|")
		ticketEntries = append(ticketEntries, DailyNoteEntry{
			Ticket:          strings.TrimSpace(worklogPieces[1]),
			Comment:         strings.TrimSpace(worklogPieces[2]),
			DurationMinutes: calcDurationMinutes(strings.TrimSpace(worklogPieces[3])),
		})
	}

	return ticketEntries
}

func ParseDailyNote(day string) []DailyNoteEntry {
	fileLines := getDailyNote(day)
	ticketTable := findTicketTable(fileLines)
	ticketEntries := parseTicketEntries(ticketTable)
	return ticketEntries
}