package tempo

import (
	"fmt"
	"time"

	"github.com/tim-hilt/tempo/noteparser"
	"github.com/tim-hilt/tempo/rest"
	"github.com/tim-hilt/tempo/util"
)

type Tempo struct {
	Api *rest.Api
}

func New(user string, password string) *Tempo {
	api := rest.New(user, password)
	tempo := &Tempo{Api: api}
	return tempo
}

func (t *Tempo) SubmitDay(day string) {
	ticketEntries := noteparser.ParseDailyNote(day)

	// Clean up first
	t.Api.DeleteWorklogs(day)

	// ...then book on clean state
	workedMinutes := 0
	for _, ticket := range ticketEntries {
		t.Api.CreateWorklog(ticket.Ticket, ticket.Comment, ticket.DurationMinutes*60, day)
		workedMinutes += ticket.DurationMinutes
	}

	hours, minutes := util.Divmod(workedMinutes, util.MINUTES_IN_HOUR)
	fmt.Println("successfully logged " + fmt.Sprint(hours) + " hours and " + fmt.Sprint(minutes) + " minutes on " + day)
}

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

// TODO: Fill with life
func WatchNotes() {
	// 1. Watch note-dir
	// 2. On Create/Change...
	// 2.1. Wait 5 minutes
	// 2.2. If 5 minutes passed and no new event happened, Delete all worklogs for that day and CreateWorklog for all the other entries
}
