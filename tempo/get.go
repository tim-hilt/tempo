package tempo

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/cli/tablecomponent"
	"github.com/tim-hilt/tempo/rest"
	"github.com/tim-hilt/tempo/util"
	"github.com/tim-hilt/tempo/util/config"
)

func (t Tempo) boilerPlate(month string) *[]rest.SearchWorklogsResult {
	start, err := time.Parse(util.MONTH_FORMAT, month)

	if err != nil {
		log.Fatal().Err(err).Str("date", month).Msg("parsing error")
	}
	end := start.AddDate(0, 1, -1)
	worklogs, err := t.Api.FindWorklogsInRange(
		start.Format(util.DATE_FORMAT),
		end.Format(util.DATE_FORMAT),
	)

	if err != nil {
		log.Fatal().Err(err).Msg("error when searching for worklogs")
	}
	return worklogs
}

func (t Tempo) GetMonthlyHours(month string) {
	worklogs := t.boilerPlate(month)

	bookedTimeSeconds := 0
	for _, worklog := range *worklogs {
		bookedTimeSeconds += worklog.DurationSeconds
	}

	hours, minutes := util.Divmod(bookedTimeSeconds/util.SECONDS_IN_MINUTE, util.MINUTES_IN_HOUR)
	fmt.Println("Worked hours for " + month + ": " +
		fmt.Sprintf("%02d", hours) + ":" + fmt.Sprintf("%02d", minutes))
}

func (t Tempo) GetWorklogsForDay(day string) {
	worklogs, err := t.Api.FindWorklogsInRange(day, day)

	if err != nil {
		log.Fatal().Err(err).Msg("error when searching for worklogs")
	}

	rows := []table.Row{}
	totalSeconds := 0
	for _, worklog := range *worklogs {
		totalSeconds += worklog.DurationSeconds
		hours, minutes := util.Divmod(
			worklog.DurationSeconds/util.SECONDS_IN_MINUTE,
			util.MINUTES_IN_HOUR,
		)
		rows = append(rows, table.Row{worklog.Issue.Ticket, worklog.Description,
			fmt.Sprintf("%02d", hours) + ":" + fmt.Sprintf("%02d", minutes)})
	}

	totalHours, totalMinutes := util.Divmod(
		totalSeconds/util.SECONDS_IN_MINUTE,
		util.MINUTES_IN_HOUR,
	)
	rows = append(
		rows,
		table.Row{
			"",
			"Sum",
			fmt.Sprintf("%02d", totalHours) + ":" + fmt.Sprintf("%02d", totalMinutes),
		},
	)

	columns := tablecomponent.CreateColumns(rows, []string{"Ticket", "Description", "Duration"})
	if err := tablecomponent.Table(columns, rows); err != nil {
		log.Fatal().Err(err).Msg("error when creating table")
	}
}

// TODO: Doesn't work for days that are used to bring down overtime
func (t Tempo) GetMonthlyOvertime(month string) {
	worklogs := t.boilerPlate(month)

	var workedSeconds float64 = 0
	daysWorked := []string{}

	for _, worklog := range *worklogs {
		workedSeconds += float64(worklog.DurationSeconds)

		// Find all days where worklogs were created
		date := worklog.DateTime[:10]
		if !util.Contains(daysWorked, date) {
			daysWorked = append(daysWorked, date)
		}
	}

	workedHours := workedSeconds / util.MINUTES_IN_HOUR / util.SECONDS_IN_MINUTE
	dailyWorkhours := config.GetWorkhours()
	overtime := workedHours - float64(len(daysWorked)*dailyWorkhours)

	fmt.Println("Overtime for " + month + ": " + fmt.Sprint(overtime) + " hours")
}

func (t Tempo) GetWorklogsForTicket(ticket string) {
	worklogs, err := t.Api.FindWorklogsForTicket(ticket)

	if err != nil {
		log.Fatal().Str("ticket", ticket).Msg("error when searching worklogs")
	}

	rows := []table.Row{}
	totalSeconds := 0
	for _, worklog := range *worklogs {
		totalSeconds += worklog.DurationSeconds
		hours, minutes := util.Divmod(
			worklog.DurationSeconds/util.SECONDS_IN_MINUTE,
			util.MINUTES_IN_HOUR,
		)
		rows = append(
			rows,
			table.Row{worklog.DateTime[:10], worklog.Description,
				fmt.Sprintf("%02d", hours) + ":" + fmt.Sprintf("%02d", minutes)},
		)
	}

	totalHours, totalMinutes := util.Divmod(
		totalSeconds/util.SECONDS_IN_MINUTE,
		util.MINUTES_IN_HOUR,
	)
	rows = append(
		rows,
		table.Row{
			"",
			"Sum",
			fmt.Sprintf("%02d", totalHours) + ":" + fmt.Sprintf("%02d", totalMinutes),
		},
	)

	columns := tablecomponent.CreateColumns(
		rows,
		// TODO: Find a way to make the length of this array fit the rows
		[]string{"Date", "Description", "Duration"},
	)
	if err := tablecomponent.Table(columns, rows); err != nil {
		log.Fatal().Err(err).Msg("error when creating table")
	}
}
