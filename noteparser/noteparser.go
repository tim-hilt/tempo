package noteparser

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/util"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"
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

func findTicketTable(file []byte) ast.Node {
	node := goldmark.New(goldmark.WithExtensions(extension.Table)).Parser().Parse(text.NewReader(file)).FirstChild()

	for node != nil {
		if node.Kind().String() == "Table" {
			if isTicketTable(node, file) {
				return node
			}
		}
		node = node.NextSibling()
	}

	log.Fatal().Msg("ticket table not found")
	return nil
}

func isTicketTable(table ast.Node, file []byte) bool {
	headers := []string{}
	tableRow := table.FirstChild()

	// TODO: Is there a better pattern than the nested while-loops?
	for tableRow != nil {
		if tableRow.Kind().String() == "TableHeader" {
			tableCell := tableRow.FirstChild()
			for tableCell != nil {
				if tableCell.Kind().String() == "TableCell" {
					headers = append(headers, string(tableCell.Text(file)))
				}
				tableCell = tableCell.NextSibling()
			}
		}
		tableRow = tableRow.NextSibling()
	}

	return util.SlicesEqual([]string{"Ticket", "Doings", "Time spent"}, headers)
}

func calcDurationMinutes(duration string) int {
	foo := strings.Split(duration, ":")
	hours, err := strconv.Atoi(foo[0])
	util.HandleErr(err, "error when converting hours-string in duration \""+duration+"\" to int")
	minutes, err := strconv.Atoi(foo[1])
	util.HandleErr(err, "error when converting minutes-string in duration \""+duration+"\" to int")
	return hours*60 + minutes
}

func parseTicketEntries(ticketTable ast.Node, file []byte) []DailyNoteEntry {
	ticketEntries := []DailyNoteEntry{}
	tableRow := ticketTable.FirstChild()

	for tableRow != nil {
		if tableRow.Kind().String() == "TableRow" {
			tableCell := tableRow.FirstChild()
			rowVals := []string{}
			for tableCell != nil {
				if tableCell.Kind().String() == "TableCell" {
					rowVals = append(rowVals, string(tableCell.Text(file)))
				}
				tableCell = tableCell.NextSibling()
			}
			ticketEntries = append(ticketEntries, DailyNoteEntry{Ticket: rowVals[0], Comment: rowVals[1], DurationMinutes: calcDurationMinutes(rowVals[2])})
		}
		tableRow = tableRow.NextSibling()
	}
	return ticketEntries
}

func ParseDailyNote(day string) []DailyNoteEntry {
	dailyNote := getDailyNote(day)
	ticketTable := findTicketTable(dailyNote)
	ticketEntries := parseTicketEntries(ticketTable, dailyNote)
	return ticketEntries
}
