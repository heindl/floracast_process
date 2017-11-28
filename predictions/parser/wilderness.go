package parser

import (
	"bitbucket.org/heindl/taxa/store"
	"sync"
	"bitbucket.org/heindl/taxa/utils"
	"context"
)

type WildernessAreaFetcher interface{
	GetWildernessArea(cxt context.Context, latitude, longitude float64) (*store.WildernessArea, error)
}

func NewWildernessAreaFetcher(taxastore store.TaxaStore) WildernessAreaFetcher {
	return &wildernessAreaFetcher{
		Store: taxastore,
		WildernessAreas: store.WildernessAreas{},
	}
}

type wildernessAreaFetcher struct{
	Store store.TaxaStore
	sync.Mutex
	WildernessAreas store.WildernessAreas
}

func (Ω *wildernessAreaFetcher) GetWildernessArea(cxt context.Context, latitude, longitude float64) (*store.WildernessArea, error) {

	for _, w := range Ω.WildernessAreas {
		if utils.CoordinatesEqual(w.Centre.Latitude, latitude) && utils.CoordinatesEqual(w.Centre.Longitude, longitude) {
			return &w, nil
		}
	}

	w, err := Ω.Store.ReadWildernessArea(cxt, latitude, longitude)
	if err != nil {
		return nil, err
	}

	Ω.Lock()
	defer Ω.Unlock()

	Ω.WildernessAreas = append(Ω.WildernessAreas, *w)

	return w, nil
}
