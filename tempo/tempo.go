package tempo

import (
	"github.com/tim-hilt/tempo/api"
)

type Tempo struct {
	Api *api.Api
}

func New(user string, password string) *Tempo {
	api := api.New(user, password)
	tempo := &Tempo{Api: api}
	return tempo
}
