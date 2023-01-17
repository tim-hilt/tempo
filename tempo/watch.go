package tempo

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/bep/debounce"
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
	"github.com/tim-hilt/tempo/noteparser"
	"github.com/tim-hilt/tempo/noteparser/parser"
	"github.com/tim-hilt/tempo/util"
	"github.com/tim-hilt/tempo/util/config"
	"golang.org/x/sync/errgroup"
)

type ChangedTickets struct {
	Tickets   []parser.DailyNoteEntry
	Submitted bool
}

var (
	changedFiles = make(map[string]ChangedTickets)
	watcher      *fsnotify.Watcher
	mut          sync.Mutex
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

			log.Trace().
				Str("file", event.Name).
				Msg("file modification")

			modifiedFile := strings.TrimSuffix(filepath.Base(event.Name), ".md")

			if event.Has(fsnotify.Write) && isDailyNote(modifiedFile) {
				f, _ := os.ReadFile(event.Name)
				if len(f) < 1 {
					// Some applications (e.g. Obsidian) seem to clear the
					// file first before they are actually storing it
					continue
				}

				log.Trace().Str("dailyNote", modifiedFile).Msg("daily note modification")

				d, err := time.Parse(util.DATE_FORMAT, modifiedFile)
				if err != nil {
					log.Error().
						Err(err).
						Str("file", modifiedFile).
						Str("expectedFormat", util.DATE_FORMAT).
						Msg("unexpected format")
					continue
				}

				if util.FromPreviousMonths(d) {
					log.Warn().
						Str("file", modifiedFile).
						Msg("file from previous months. Won't submit")
					continue
				}

				ticketEntries, err := noteparser.ParseDailyNote(modifiedFile)

				if err != nil {
					log.Error().Err(err).Str("file", modifiedFile).Msg("error when parsing")
					continue
				}

				prevTicketEntries := changedFiles[modifiedFile].Tickets

				if !parser.DailyNoteEntriesEqual(ticketEntries, prevTicketEntries) &&
					!changedFiles[modifiedFile].Submitted {
					mut.Lock()
					changedFiles[modifiedFile] = ChangedTickets{
						Tickets:   ticketEntries,
						Submitted: false,
					}
					mut.Unlock()

					log.Trace().
						Strs("changedFiles", util.GetKeys(changedFiles)).
						Str("duration", debounceDuration.String()).
						Msg("submitting file in")
					debounced(t.submitChanged)
				} else {
					log.Trace().Str("file", modifiedFile).Int("size", len(f)).Msg("ticket entries equal")
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
	match, _ := regexp.MatchString(`\d{4}-\d{2}-\d{2}`, file)
	return match
}

func (t *Tempo) submitChanged() {
	mut.Lock()
	defer mut.Unlock()

	var wg sync.WaitGroup
	wg.Add(len(changedFiles))

	for changedFile := range changedFiles {
		changedFile := changedFile
		go func() {
			defer wg.Done()
			log.Info().Str("file", changedFile).Msg("creating worklogs")

			if err := t.submit(changedFile); err != nil {
				log.Error().
					Err(err).
					Str("file", changedFile).
					Msg("error when submitting tickets")
				return
			}

			log.Info().
				Str("file", changedFile).
				Msg("finished creating worklogs")
			tickets := changedFiles[changedFile].Tickets
			changedFiles[changedFile] = ChangedTickets{Tickets: tickets, Submitted: true}
		}()

	}

	wg.Wait()
}
