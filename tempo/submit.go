package tempo

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/noteparser"
	"github.com/tim-hilt/tempo/util"
	"golang.org/x/sync/errgroup"
)

func (t Tempo) SubmitDay(day string) {
	if err := t.submit(day); err != nil {
		log.Fatal().Err(err).Str("day", day).Msg("error when submitting")
	}
}

func (t Tempo) submit(note string) error {
	ticketEntries, err := noteparser.ParseDailyNote(note)

	if err != nil {
		return err
	}

	if err := t.Api.DeleteWorklogs(note); err != nil {
		return err
	}

	errs, _ := errgroup.WithContext(context.Background())

	workedMinutes := 0

	for _, ticket := range ticketEntries {
		ticket := ticket // Necessary as of https://go.dev/doc/faq#closures_and_goroutines
		errs.Go(func() error {
			if err := t.Api.CreateWorklog(ticket.Ticket, ticket.Comment, ticket.DurationMinutes*60, note); err != nil {
				return err
			}
			workedMinutes += ticket.DurationMinutes
			return nil
		})
	}

	if err := errs.Wait(); err != nil {
		return err
	}

	hours, minutes := util.Divmod(workedMinutes, util.MINUTES_IN_HOUR)
	log.Info().Int("hours", hours).Int("minutes", minutes).Msg("successfully logged")
	return nil
}
