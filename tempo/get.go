package tempo

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/tim-hilt/tempo/cli/tablecomponent"
	"github.com/tim-hilt/tempo/util"
)

func (t *Tempo) GetMonthlyHours() {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 1, -1)
	worklogs := t.Api.FindWorklogsInRange(start.Format(util.DATE_FORMAT), end.Format(util.DATE_FORMAT))

	bookedTimeSeconds := 0
	for _, worklog := range *worklogs {
		bookedTimeSeconds += worklog.DurationSeconds
	}

	hours, minutes := util.Divmod(bookedTimeSeconds/util.SECONDS_IN_MINUTE, util.MINUTES_IN_HOUR)
	fmt.Println("worked " + fmt.Sprint(hours) + " hours and " + fmt.Sprint(minutes) + " minutes in current month")
}

func (t *Tempo) GetTicketsForDay(day string) {
	worklogs := t.Api.FindWorklogsInRange(day, day)
	rows := []table.Row{}
	for _, worklog := range *worklogs {
		hours, minutes := util.Divmod(worklog.DurationSeconds/util.SECONDS_IN_MINUTE, util.MINUTES_IN_HOUR)
		rows = append(rows, table.Row{worklog.Issue.Ticket, worklog.Issue.Description, fmt.Sprint(hours) + "h" + fmt.Sprint(minutes) + "m"})
	}

	columns := []table.Column{
		// TODO: Make this work + more beautiful!
		{Title: "Ticket", Width: util.Max(getMaxChars(rows, 0), len("Ticket"))},
		{Title: "Description", Width: util.Max(getMaxChars(rows, 1), len("Description"))},
		{Title: "Duration", Width: util.Max(getMaxChars(rows, 2), len("Duration"))},
	}
	tablecomponent.Table(columns, rows)
}

// TODO: Make table.Row => string
func getMaxChars(rows []table.Row, index int) int {
	max := 0
	for _, entry := range rows {
		if len(entry[index]) > max {
			max = len(entry[index])
		}
	}
	return max
}
