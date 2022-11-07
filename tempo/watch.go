package tempo

import (
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/bep/debounce"
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/util"
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
	go t.watchLoop(watcher, &wg)

	notesDir := util.GetConfigParams().Notesdir
	addDirs(watcher, []string{notesDir})

	wg.Wait()
}

func (t Tempo) watchLoop(watcher *fsnotify.Watcher, wg *sync.WaitGroup) {
	wg.Add(1)
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
		t.SubmitDay(note)
	}
	changedFiles.Reset()
}
