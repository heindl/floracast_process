package protected_areas

import (
	"context"
	"bitbucket.org/heindl/process/store"
	"sync"
	"bitbucket.org/heindl/process/geofeatures"
)

type ProtectedAreaCache interface {
	GetProtectedArea(cxt context.Context, latitude, longitude float64) (*ProtectedArea, error)
}

func NewProtectedAreaCache(florastore store.FloraStore) (ProtectedAreaCache, error) {
	return &cache{
		florastore: florastore,
		areas:      map[geofeatures.CoordinateKey]*ProtectedArea{},
	}, nil
}

type cache struct {
	florastore store.FloraStore
	sync.Mutex
	areas      map[geofeatures.CoordinateKey]*ProtectedArea
}

func (Ω *cache) GetProtectedArea(cxt context.Context, latitude, longitude float64) (*ProtectedArea, error) {

	coordKey, err := geofeatures.NewCoordinateKey(latitude, longitude)
	if err != nil {
		return nil, err
	}

	Ω.Lock()
	if v, ok := Ω.areas[coordKey]; ok {
		Ω.Unlock()
		return v, nil
	}
	Ω.Unlock()

	area, err := FetchOne(cxt, Ω.florastore, coordKey)
	if err != nil {
		return nil, err
	}

	Ω.Lock()
	defer Ω.Unlock()

	Ω.areas[coordKey] = area

	return area, nil
}

