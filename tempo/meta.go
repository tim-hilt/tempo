package tempo

import (
	"github.com/tim-hilt/tempo/rest"
)

type Tempo struct {
	Api *rest.Api
}

func New(user string, password string) *Tempo {
	api := rest.New(user, password)
	tempo := &Tempo{Api: api}
	return tempo
}
