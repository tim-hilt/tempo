package tempo

import (
	"fmt"

	"github.com/tim-hilt/tempo/noteparser"
	"github.com/tim-hilt/tempo/util"
)

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
