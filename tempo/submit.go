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

func (t Tempo) submit(day string) error {
	worklogs, err := t.Api.FindWorklogsInRange(day, day)

	if err != nil {
		return err
	}

	// TODO: One way of submitting only changed ticket-entries would be to keep a copy of the last state of the ticketEntries and then diff their contents. Afterwards:
	// - Deleted entries are deleted
	// - Changed entries are deleted and re-submitted
	// - New entries are submitted
	// I still need to clarify how I can detect a changed ticket though
	if err := t.Api.DeleteWorklogs(worklogs); err != nil {
		return err
	}

	ticketEntries, err := noteparser.ParseDailyNote(day)

	if err != nil {
		return err
	}

	workedSeconds := 0
	errs, _ := errgroup.WithContext(context.Background())

	for _, ticket := range ticketEntries {
		ticket := ticket // Necessary as of https://go.dev/doc/faq#closures_and_goroutines
		errs.Go(func() error {
			if err := t.Api.CreateWorklog(ticket.Ticket, ticket.Comment, ticket.DurationSeconds, day); err != nil {
				return err
			}
			workedSeconds += ticket.DurationSeconds
			return nil
		})
	}

	if err := errs.Wait(); err != nil {
		return err
	}

	hours, minutes := util.Divmod(workedSeconds/util.SECONDS_IN_MINUTE, util.MINUTES_IN_HOUR)
	log.Info().Int("hours", hours).Int("minutes", minutes).Msg("successfully logged")
	return nil
}
