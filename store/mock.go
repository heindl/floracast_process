package store

import (
	"bitbucket.org/heindl/mgoeco"
	"testing"
)

func NewMockStore(t *testing.T) SpeciesStore {
	server, m := mgoeco.TestMongo(t)
	return SpeciesStore(&store{server, m})
}
