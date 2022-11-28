package tempo

import (
	"github.com/tim-hilt/tempo/api"
	"github.com/tim-hilt/tempo/noteparser/parser"
)

type Tempo struct {
	Api                   *api.Api
	PreviousTicketEntries []parser.DailyNoteEntry
}

func New(user string, password string) *Tempo {
	api := api.New(user, password)
	tempo := &Tempo{Api: api, PreviousTicketEntries: []parser.DailyNoteEntry{}}
	return tempo
}
