package tempo

import (
	"fmt"
	"sync"

	"github.com/tim-hilt/tempo/noteparser"
	"github.com/tim-hilt/tempo/util"
)

func (t *Tempo) SubmitDay(day string) {
	ticketEntries := noteparser.ParseDailyNote(day)

	// Clean up first
	t.Api.DeleteWorklogs(day)

	// ...then book on clean state
	workedMinutes := 0
	var wg sync.WaitGroup

	for _, ticket := range ticketEntries {
		wg.Add(1)
		go func(ticket noteparser.DailyNoteEntry) {
			defer wg.Done()
			t.Api.CreateWorklog(ticket.Ticket, ticket.Comment, ticket.DurationMinutes*60, day)
			workedMinutes += ticket.DurationMinutes
		}(ticket)
	}

	wg.Wait()

	hours, minutes := util.Divmod(workedMinutes, util.MINUTES_IN_HOUR)
	fmt.Println("successfully logged " + fmt.Sprintf("%02d", hours) + " hours and " + fmt.Sprintf("%02d", minutes) + " minutes on " + day)
}
