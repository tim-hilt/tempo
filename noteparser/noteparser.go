package noteparser

import (
	"errors"
	"fmt"
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

func findTicketTable(file []byte) (ast.Node, error) {
	md := goldmark.New(goldmark.WithExtensions(extension.Table)).Parser().Parse(text.NewReader(file))
	md.Dump(file, 1)
	traverseAst(md)
	node := md

	for {
		node, err := findNext(node, "Table")
		if err != nil {
			return nil, err
		} else if isTicketTable(file, node) {
			return node, nil
		}
	}
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

	return slicesEqual([]string{"Ticket", "Doings", "Time spent"}, headers)
}

func slicesEqual(a, b []string) bool {
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

func findNext(st ast.Node, kind string) (ast.Node, error) {
	if st == nil {
		return nil, errors.New(kind + " not found in ast")
	} else if st.Kind().String() == kind {
		log.Info().Msg("found table")
		return st, nil
	} else if st.HasChildren() {
		return findNext(st.FirstChild(), kind)
	} else if st.NextSibling() == nil {
		return findNext(st.Parent(), kind)
	} else {
		return findNext(st.NextSibling(), kind)
	}
}

// TODO: Adapt findNext to reflect the below func
// TODO: Delete func once the above TODO is finished
func traverseAst(st ast.Node) {
	if st == nil {
		return
	} else if st.HasChildren() {
		fmt.Println(st.Kind())
		traverseAst(st.FirstChild())
	} else if st.NextSibling() != nil {
		fmt.Println(st.Kind())
		traverseAst(st.NextSibling())
	} else {
		fmt.Println(st.Kind())
		st = findNextParentSibling(st)
		traverseAst(st)
	}
}

func findNextParentSibling(st ast.Node) ast.Node {
	if st.Parent() == nil {
		return nil
	} else if st.Parent().NextSibling() != nil {
		return st.Parent().NextSibling()
	}
	return findNextParentSibling(st.Parent())
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
