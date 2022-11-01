package tablecomponent

import (
	"os"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/util"
	"github.com/tim-hilt/tempo/util/logging"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	logging.SetLoglevel()
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder())

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, tea.Quit
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func Table(columns []table.Column, rows []table.Row) {

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithHeight(len(rows)),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true)
	s.Selected = lipgloss.NewStyle()
	t.SetStyles(s)

	m := model{t}
	err := tea.NewProgram(m).Start()
	util.HandleErr(err, "error running table")
}
