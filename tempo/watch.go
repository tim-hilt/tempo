package tempo

import (
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/bep/debounce"
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/noteparser"
	"github.com/tim-hilt/tempo/util"
	"github.com/tim-hilt/tempo/util/config"
	"github.com/tim-hilt/tempo/util/set"
)

var changedFiles = set.New[string]()

func (t *Tempo) WatchNotes() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal().Err(err).Msg("error when creating watcher")
	}
	defer watcher.Close()
	var wg sync.WaitGroup
	wg.Add(1)
	go t.watchLoop(watcher, &wg)

	notesDir := config.GetNotesdir()
	addDirs(watcher, []string{notesDir})

	wg.Wait()
	log.Fatal().Msg("waitgroup is finished")
}

func (t Tempo) watchLoop(watcher *fsnotify.Watcher, wg *sync.WaitGroup) {
	defer wg.Done()
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
				log.Info().Str("file", modifiedFile).Msg("file modification")
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
			log.Fatal().Err(err).Str("watchDir", dir).Msg("error when adding dir to watcher")
		}
		log.Info().Str("watchDir", dir).Msg("watching directory")
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
	log.Info().
		Strs("files", changedFiles.Items()).
		Msg("creating worklogs for files")
	for _, note := range changedFiles.Items() {
		note = filepath.Base(note)
		if err := t.submit(note); err != nil {
			log.Error().
				Err(err).
				Str("file", note).
				Strs("files", changedFiles.Items()).
				Msg("error when submitting tickets")
			return
		}
	}
	log.Info().
		Strs("files", changedFiles.Items()).
		Msg("finished creating worklogs")
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
	var wg sync.WaitGroup

	for _, ticket := range ticketEntries {
		wg.Add(1)
		go func(ticket noteparser.DailyNoteEntry) {
			defer wg.Done()
			if err := t.Api.CreateWorklog(ticket.Ticket, ticket.Comment, ticket.DurationMinutes*60, note); err != nil {
				log.Error().Err(err).Msg("error when creating worklog")
			}
			workedMinutes += ticket.DurationMinutes
		}(ticket)
	}

	wg.Wait()

	hours, minutes := util.Divmod(workedMinutes, util.MINUTES_IN_HOUR)
	log.Info().Int("hours", hours).Int("minutes", minutes).Msg("successfully logged")
	return nil
}
