package store

import (
	"bitbucket.org/heindl/provision/mgoeco"
	"testing"
	"github.com/jonboulle/clockwork"
	"time"
)

func NewMockStore(t *testing.T) SpeciesStore {
	server, m := mgoeco.TestMongo(t)
	return SpeciesStore(&store{server, m, clockwork.NewFakeClockAt(time.Now())})
}
