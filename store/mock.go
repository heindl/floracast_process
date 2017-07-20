package store

import (
	"testing"
	"github.com/jonboulle/clockwork"
	"time"
	"bitbucket.org/heindl/provision/dseco"
)

func NewMockStore(t *testing.T) TaxaStore {
	ds, err := dseco.NewMockDatastore()
	if err != nil {
		panic(err)
	}
	return TaxaStore(&store{clockwork.NewFakeClockAt(time.Now()), ds})
}
