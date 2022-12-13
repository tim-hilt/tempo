package tempo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/noteparser"
	"github.com/tim-hilt/tempo/util"
	"github.com/tim-hilt/tempo/util/config"
	"golang.org/x/sync/errgroup"
)

// TODO: Don't permit submitting day from last month! Booking already closed
// TODO: Should expect a Time-instance. Not string
func (t *Tempo) SubmitDay(day string) {
	if err := t.submit(day); err != nil {
		log.Fatal().Err(err).Str("day", day).Msg("error when submitting")
	}
}

// TODO: Should expect a Time-instance. Not string
func dayNotOlderThanOneMonth(day string) error {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	d, err := time.Parse(util.DATE_FORMAT, day)

	// TODO: This has to be caught earlier!
	if err != nil {
		return err
	}

	if firstOfMonth.Before(d) {
		return nil
	}

	return errors.New("day " + fmt.Sprint(d) + " is older than current month")
}

// TODO: Don't permit submitting day from last month! Booking already closed
// TODO: Should expect a Time-instance. Not string
func (t *Tempo) submit(day string) error {

	if err := dayNotOlderThanOneMonth(day); err != nil {
		return err
	}

	ticketEntries, err := noteparser.ParseDailyNote(day)

	if err != nil {
		return err
	}

	// TODO: This concept doesn't quite work, when dealing with multiple files.
	// I'd rather need a map [datestring]: []parser.DailyNoteEntry
	// Also: It isn't relevant for submitting a single file. Maybe I can generalize this?
	if cmp.Equal(ticketEntries, t.PreviousTicketEntries) {
		log.Info().Msg("ticketEntries equal. not submitting")
		return nil
	}

	worklogs, err := t.Api.FindWorklogsInRange(day, day)

	if err != nil {
		return err
	}

	if err := t.Api.DeleteWorklogs(worklogs); err != nil {
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

	overtimeMinutes := (hours*util.MINUTES_IN_HOUR + minutes) - (config.GetWorkhours() * util.MINUTES_IN_HOUR)
	overtimeHours, overtimeMinutes := util.Divmod(overtimeMinutes, util.MINUTES_IN_HOUR)
	log.Trace().Int("overtimeHours", overtimeHours).Int("overtimeMinutes", overtimeMinutes).Msg("")

	t.PreviousTicketEntries = ticketEntries

	return nil
}
