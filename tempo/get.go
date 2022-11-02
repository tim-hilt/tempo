package tempo

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/tim-hilt/tempo/cli/tablecomponent"
	"github.com/tim-hilt/tempo/cmd/flags"
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

	columns := tablecomponent.CreateColumns(rows, []string{"Ticket", "Description", "Duration"})
	tablecomponent.Table(columns, rows)
}

func (t *Tempo) GetMonthlyOvertime(month string) {
	start, err := time.Parse(util.MONTH_FORMAT, month)
	end := start.AddDate(0, 1, -1)
	util.HandleErr(err, "Error when parsing "+month+" to time.Time")
	worklogs := t.Api.FindWorklogsInRange(start.Format(util.DATE_FORMAT), end.Format(util.DATE_FORMAT))

	var workedSeconds float64 = 0
	daysWorked := []string{}

	for _, worklog := range *worklogs {
		workedSeconds += float64(worklog.DurationSeconds)

		// Find all days where worklogs were created
		date := worklog.Date[:10]
		if !util.Contains(daysWorked, date) {
			daysWorked = append(daysWorked, date)
		}
	}

	workedHours := workedSeconds / util.MINUTES_IN_HOUR / util.SECONDS_IN_MINUTE
	overtime := workedHours - float64(len(daysWorked)*flags.DailyWorkhours)

	fmt.Println("Overtime for " + month + ": " + fmt.Sprint(overtime) + " hours")
}
