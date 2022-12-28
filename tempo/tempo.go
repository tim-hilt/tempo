package tempo

import (
	"github.com/tim-hilt/tempo/api"
)

type Tempo struct {
	Api *api.Api
}

func New() *Tempo {
	api := api.New()
	tempo := &Tempo{Api: api}
	return tempo
}
