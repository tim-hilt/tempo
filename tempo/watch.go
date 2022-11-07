package tempo

import (
	"context"
	"path/filepath"
	"regexp"
	"time"

	"github.com/bep/debounce"
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/noteparser"
	"github.com/tim-hilt/tempo/util"
	"github.com/tim-hilt/tempo/util/set"
	"golang.org/x/sync/errgroup"
)

var changedFiles = set.New[string]()

func (t *Tempo) WatchNotes() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal().Err(err).Msg("error when creating watcher")
	}
	defer watcher.Close()
	go t.watchLoop(watcher)

	notesDir := util.GetConfigParams().Notesdir
	addDirs(watcher, []string{notesDir})

	<-make(chan struct{})
	log.Fatal().Msg("main goroutine unblocked")
}

func (t Tempo) watchLoop(watcher *fsnotify.Watcher) {
	debounced := debounce.New(5 * time.Minute)
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				log.Fatal().Msg("event not ok!")
				return
			}
			modifiedFile := event.Name
			if event.Has(fsnotify.Write) && isDailyNote(modifiedFile) {
				log.Info().Msg("modified file:" + modifiedFile)
				changedFiles.Add(modifiedFile)
				debounced(t.submitChanged)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				log.Fatal().Msg("error not ok!")
				return
			}
			log.Error().Err(err).Msg("error on channel \"errors\"")
		}
	}
}

func addDirs(watcher *fsnotify.Watcher, dirs []string) {
	for _, dir := range dirs {
		err := watcher.Add(dir)
		if err != nil {
			log.Fatal().Err(err).Msg("error when adding " + dir + " to watcher")
		}
		log.Info().Msg("watching directory " + dir)
	}
}

func isDailyNote(file string) bool {
	file = filepath.Base(file)
	if match, _ := regexp.MatchString(`\d{4}-\d{2}-\d{2}.md`, file); match {
		return true
	}
	return false
}

func (t Tempo) submitChanged() {
	log.Info().Msg("Creating worklogs for the following files: " + changedFiles.String())
	for _, note := range changedFiles.Items() {
		note = filepath.Base(note)
		if err := t.submit(note); err != nil {
			log.Error().Err(err).Msg("error when submitting tickets on " + note)
			return
		}
	}
	log.Info().Msg("Finished creating worklogs")
	changedFiles.Reset()
}

func (t Tempo) submit(note string) error {
	ticketEntries, err := noteparser.ParseDailyNote(note)

	if err != nil {
		return err
	}
	if err := t.Api.DeleteWorklogs(note); err != nil {
		return err
	}

	workedMinutes := 0
	errs, _ := errgroup.WithContext(context.Background())

	for _, ticket := range ticketEntries {
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
