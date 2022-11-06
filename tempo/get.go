package tempo

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/cli/tablecomponent"
	"github.com/tim-hilt/tempo/util"
)

// TODO: Could also pass month as arg as in overtime-func below
func (t *Tempo) GetMonthlyHours() {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 1, -1)
	worklogs, err := t.Api.FindWorklogsInRange(start.Format(util.DATE_FORMAT), end.Format(util.DATE_FORMAT))

	if err != nil {
		log.Fatal().Err(err).Msg("error when searching for worklogs")
	}

	bookedTimeSeconds := 0
	for _, worklog := range *worklogs {
		bookedTimeSeconds += worklog.DurationSeconds
	}

	hours, minutes := util.Divmod(bookedTimeSeconds/util.SECONDS_IN_MINUTE, util.MINUTES_IN_HOUR)
	fmt.Println("Worked hours for " + start.Format(util.MONTH_FORMAT) + ": " +
		fmt.Sprintf("%02d", hours) + "." + fmt.Sprintf("%02d", minutes))
}

func (t *Tempo) GetTicketsForDay(day string) {
	worklogs, err := t.Api.FindWorklogsInRange(day, day)

	if err != nil {
		log.Fatal().Err(err).Msg("error when searching for worklogs")
	}

	rows := []table.Row{}
	for _, worklog := range *worklogs {
		hours, minutes := util.Divmod(worklog.DurationSeconds/util.SECONDS_IN_MINUTE, util.MINUTES_IN_HOUR)
		rows = append(rows, table.Row{worklog.Issue.Ticket, worklog.Issue.Description,
			fmt.Sprintf("%02d", hours) + "h" + fmt.Sprintf("%02d", minutes) + "m"})
	}

	columns := tablecomponent.CreateColumns(rows, []string{"Ticket", "Description", "Duration"})
	if err := tablecomponent.Table(columns, rows); err != nil {
		log.Fatal().Err(err).Msg("error when creating table")
	}
}

// TODO: Doesn't work for days that are used to bring down overtime
func (t *Tempo) GetMonthlyOvertime(month string) {
	start, err := time.Parse(util.MONTH_FORMAT, month)

	if err != nil {
		log.Fatal().Err(err).Msg("error when parsing date: " + month)
	}

	end := start.AddDate(0, 1, -1)
	worklogs, err := t.Api.FindWorklogsInRange(start.Format(util.DATE_FORMAT), end.Format(util.DATE_FORMAT))

	if err != nil {
		log.Fatal().Err(err).Msg("error when searching for worklogs")
	}

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
	dailyWorkhours := util.GetConfigParams().DailyWorkhours
	overtime := workedHours - float64(len(daysWorked)*dailyWorkhours)

	fmt.Println("Overtime for " + month + ": " + fmt.Sprint(overtime) + " hours")
}
