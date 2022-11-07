package tempo

import (
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/noteparser"
	"github.com/tim-hilt/tempo/util"
)

func (t Tempo) SubmitDay(day string) {
	ticketEntries, err := noteparser.ParseDailyNote(day)

	if err != nil {
		log.Fatal().Err(err).Msg("error when parsing daily note")
	}

	// Clean up first
	if err := t.Api.DeleteWorklogs(day); err != nil {
		log.Fatal().Err(err).Msg("error when deleting worklogs")
	}

	// ...then book on clean state
	workedMinutes := 0
	var wg sync.WaitGroup

	for _, ticket := range ticketEntries {
		wg.Add(1)
		go func(ticket noteparser.DailyNoteEntry) {
			defer wg.Done()
			if err := t.Api.CreateWorklog(ticket.Ticket, ticket.Comment, ticket.DurationMinutes*60, day); err != nil {
				log.Fatal().Err(err).Msg("error whin creating worklog")
			}
			workedMinutes += ticket.DurationMinutes
		}(ticket)
	}

	wg.Wait()

	hours, minutes := util.Divmod(workedMinutes, util.MINUTES_IN_HOUR)
	fmt.Println("successfully logged " + fmt.Sprintf("%02d", hours) + " hours and " +
		fmt.Sprintf("%02d", minutes) + " minutes on " + day)
}
