package store

import (
	"testing"
	"github.com/jonboulle/clockwork"
	"time"
)

func NewMockStore(t *testing.T) TaxaStore {
	ds, err := NewMockFirestore()
	if err != nil {
		panic(err)
	}
	return TaxaStore(&store{clockwork.NewFakeClockAt(time.Now()), ds})
}
