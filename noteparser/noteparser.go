package noteparser

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/viper"
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
	notesDir := viper.GetString("notesDir")
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

func findTicketTable(file []byte) (ast.Node, error) {
	md := goldmark.New(goldmark.WithExtensions(extension.Table)).Parser().Parse(text.NewReader(file))
	node := md.FirstChild()

	for node != nil {
		if node.Kind().String() == "Table" && isTicketTable(file, node) {
			return node, nil
		}
		node = node.NextSibling()
	}

	return nil, errors.New("ticket table not found")
}

func isTicketTable(file []byte, table ast.Node) bool {
	tableRow := table.FirstChild()
	headers := []string{}

	for tableRow != nil {
		if tableRow.Kind().String() == "TableHeader" {
			tableCell := tableRow.FirstChild()
			for tableCell != nil {
				headers = append(headers, string(tableCell.Text(file)))
				tableCell = tableCell.NextSibling()
			}
		}
		tableRow = tableRow.NextSibling()
	}

	return hasColumnHeaders([]string{"Ticket", "Doings", "Time spent"}, headers)
}

func hasColumnHeaders(a, b []string) bool {
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
			ticket := string(tableRow.FirstChild().Text(file))
			doings := string(tableRow.FirstChild().NextSibling().Text(file))
			duration := string(tableRow.FirstChild().NextSibling().NextSibling().Text(file))

			ticketEntries = append(ticketEntries, DailyNoteEntry{Ticket: ticket, Comment: doings, DurationMinutes: calcDurationMinutes(duration)})
		}
		tableRow = tableRow.NextSibling()
	}

	return ticketEntries
}

func ParseDailyNote(day string) []DailyNoteEntry {
	dailyNote := getDailyNote(day)
	ticketTable, err := findTicketTable(dailyNote)
	util.HandleErr(err, "error when looking for ticket-table")
	ticketEntries := parseTicketEntries(ticketTable, dailyNote)
	return ticketEntries
}
