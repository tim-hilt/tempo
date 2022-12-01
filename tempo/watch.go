package tempo

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/bep/debounce"
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/util"
	"github.com/tim-hilt/tempo/util/config"
	"golang.org/x/sync/errgroup"
)

var (
	changedFile = util.NO_CHANGED_FILES
	watcher     *fsnotify.Watcher
)

func init() {
	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("error when creating watcher")
	}
}

func (t *Tempo) WatchNotes() {
	defer watcher.Close()
	errs, _ := errgroup.WithContext(context.Background())
	errs.Go(t.watchLoop)

	notesDir := config.GetNotesdir()
	addDirs(watcher, []string{notesDir})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Info().
				Str("signal", sig.String()).
				Msg("Application was quit")
			os.Exit(0)
		}
	}()

	if err := errs.Wait(); err != nil {
		log.Fatal().
			Err(err).
			Msg("main loop exited with error")
	}
}

func (t *Tempo) watchLoop() error {
	debounceDuration := 15 * time.Second
	debounced := debounce.New(debounceDuration)
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return errors.New("event not ok")
			}
			modifiedFile := event.Name
			if event.Has(fsnotify.Write) && isDailyNote(modifiedFile) {
				log.Trace().
					Str("file", modifiedFile).
					Msg("file modification")
				// use debounce if file not set or same file edited
				if changedFile == util.NO_CHANGED_FILES || changedFile == modifiedFile {
					changedFile = modifiedFile
					// TODO: Theoretically we could detect changed ticket-Tables here directly
					// and abort the submit
					log.Trace().
						Str("lastModified", changedFile).
						Str("duration", debounceDuration.String()).
						Msg("submitting file in")
					debounced(t.submitChanged)
				} else { // submit last file directly otherwise
					log.Trace().
						Str("lastModified", changedFile).
						Str("newlyModified", modifiedFile).
						Msg("debounce interrupted. submitting last modified file immediately")
					t.submitChanged()
					changedFile = modifiedFile
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return errors.New("error not ok")
			}
			log.Error().
				Err(err).
				Msg("error on channel \"errors\"")
		}
	}
}

func addDirs(watcher *fsnotify.Watcher, dirs []string) {
	for _, dir := range dirs {
		err := watcher.Add(dir)
		if err != nil {
			log.Fatal().
				Err(err).
				Str("watchDir", dir).
				Msg("error when adding dir to watcher")
		}
		log.Info().
			Str("watchDir", dir).
			Msg("watching directory")
	}
}

func isDailyNote(file string) bool {
	file = filepath.Base(file)
	if match, _ := regexp.MatchString(`\d{4}-\d{2}-\d{2}.md`, file); match {
		return true
	}
	return false
}

// TODO: Don't permit submitting day from last month! Booking already closed
func (t *Tempo) submitChanged() {
	if changedFile == util.NO_CHANGED_FILES {
		log.Info().Msg("no changed file")
		return
	}

	// Format file with path to date
	date := filepath.Base(changedFile)
	date = strings.TrimSuffix(date, ".md")

	log.Info().Str("file", changedFile).Msg("creating worklogs")

	if err := t.submit(date); err != nil {
		log.Error().
			Err(err).
			Str("file", changedFile).
			Msg("error when submitting tickets")
		return
	}

	log.Info().
		Str("file", changedFile).
		Msg("finished creating worklogs")

	changedFile = util.NO_CHANGED_FILES
}
