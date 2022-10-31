package tempo

import (
	"fmt"
	"time"

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
	for _, worklog := range *worklogs {
		hours, minutes := util.Divmod(worklog.DurationSeconds/util.SECONDS_IN_MINUTE, util.MINUTES_IN_HOUR)
		// TODO: Next line as table with bubbletea
		fmt.Println("Ticket: " + worklog.Issue.Ticket + " Description: " + worklog.Issue.Description + " Duration: " + fmt.Sprint(hours) + "h" + fmt.Sprint(minutes) + "m")
	}
}