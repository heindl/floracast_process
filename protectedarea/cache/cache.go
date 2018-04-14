package cache

import (
	"bitbucket.org/heindl/process/protectedarea"
	"bitbucket.org/heindl/process/store"
	"bitbucket.org/heindl/process/terra/geoembed"
	"context"
	"sync"
)

// ProtectedAreaCache fetches and stores ProtectedAreas locally in case the same record
// is requested multiple times by a process.
type ProtectedAreaCache interface {
	GetProtectedArea(cxt context.Context, latitude, longitude float64) (protectedarea.ProtectedArea, error)
}

// NewProtectedAreaCache creates a new ProtectedAreaCache
func NewProtectedAreaCache(florastore store.FloraStore) (ProtectedAreaCache, error) {
	return &cache{
		florastore: florastore,
		areas:      map[geoembed.S2Key]protectedarea.ProtectedArea{},
	}, nil
}

type cache struct {
	florastore store.FloraStore
	sync.Mutex
	areas map[geoembed.S2Key]protectedarea.ProtectedArea
}

func (Ω *cache) GetProtectedArea(cxt context.Context, latitude, longitude float64) (protectedarea.ProtectedArea, error) {

	coordKey, err := geoembed.NewS2Key(latitude, longitude)
	if err != nil {
		return nil, err
	}

	Ω.Lock()
	if v, ok := Ω.areas[coordKey]; ok {
		Ω.Unlock()
		return v, nil
	}
	Ω.Unlock()

	area, err := protectedarea.FetchOne(cxt, Ω.florastore, coordKey)
	if err != nil {
		return nil, err
	}

	Ω.Lock()
	defer Ω.Unlock()

	Ω.areas[coordKey] = area

	return area, nil
}
