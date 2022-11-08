package tempo

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"time"

	"github.com/bep/debounce"
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/util/config"
	"github.com/tim-hilt/tempo/util/set"
	"golang.org/x/sync/errgroup"
)

var (
	changedFiles = set.New[string]()
	watcher      *fsnotify.Watcher
)

func init() {
	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal().Err(err).Msg("error when creating watcher")
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
			log.Info().Str("signal", sig.String()).Msg("Application was quit")
			os.Exit(0)
		}
	}()

	if err := errs.Wait(); err != nil {
		log.Fatal().Err(err).Msg("main loop exited with error")
	}
}

func (t Tempo) watchLoop() error {
	debounced := debounce.New(5 * time.Minute)
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return errors.New("event not ok")
			}
			modifiedFile := event.Name
			if event.Has(fsnotify.Write) && isDailyNote(modifiedFile) {
				log.Info().Str("file", modifiedFile).Msg("file modification")
				changedFiles.Add(modifiedFile)
				debounced(t.submitChanged)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return errors.New("error not ok")
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
