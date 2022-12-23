package tempo

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/noteparser"
	"github.com/tim-hilt/tempo/util"
	"github.com/tim-hilt/tempo/util/config"
	"golang.org/x/sync/errgroup"
)

func (t *Tempo) SubmitDate(date string) {

	if util.IsFullDate(date) {
		if err := t.submit(date); err != nil {
			log.Fatal().Err(err).Str("day", date).Msg("error when submitting")
		}
	} else if util.IsYearMonth(date) {
		if err := t.submitMonth(); err != nil {
			log.Fatal().Err(err).Str("month", date).Msg("error when submitting")
		}
	}
}

func (t *Tempo) submitMonth() error {
	notesDir := config.GetNotesdir()
	fs, err := os.ReadDir(notesDir)

	if err != nil {
		return errors.New("error when reading directory")
	}

	toSubmit := []string{}

	for _, f := range fs {
		fn := strings.TrimSuffix(f.Name(), ".md")
		if f.Type().IsRegular() && util.IsFullDate(fn) {
			d, err := time.Parse(util.DATE_FORMAT, fn)

			if err != nil {
				return errors.New("error parsing to time.Time, expected format " + util.DATE_FORMAT)
			}

			if !olderThanCurrentMonth(d) {
				toSubmit = append(toSubmit, fn)
			}
		}
	}

	errs, _ := errgroup.WithContext(context.Background())

	for _, d := range toSubmit {
		d := d // Necessary as of https://go.dev/doc/faq#closures_and_goroutines
		errs.Go(func() error {
			if err := t.submit(d); err != nil {
				return err
			}
			return err
		})
	}

	if err := errs.Wait(); err != nil {
		return err
	}

	return nil
}

func olderThanCurrentMonth(day time.Time) bool {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	return day.Before(firstOfMonth)
}

func (t *Tempo) submit(day string) error {

	d, err := time.Parse(util.DATE_FORMAT, day)
	if err != nil {
		return err
	}

	if olderThanCurrentMonth(d) {
		return errors.New("day " + fmt.Sprint(d) + " is older than current month")
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
