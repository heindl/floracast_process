package parser

import (
	"bitbucket.org/heindl/taxa/store"
	"sync"
	"bitbucket.org/heindl/taxa/utils"
	"context"
	"fmt"
	"google.golang.org/genproto/googleapis/type/latlng"
	"strings"
)

type WildernessAreaFetcher interface{
	GetWildernessArea(cxt context.Context, latitude, longitude float64) (*store.ProtectedArea, error)
}

func NewWildernessAreaFetcher(taxastore store.TaxaStore) WildernessAreaFetcher {
	return &wildernessAreaFetcher{
		Store: taxastore,
		WildernessAreas: store.ProtectedAreas{},
	}
}

type wildernessAreaFetcher struct{
	Store store.TaxaStore
	sync.Mutex
	WildernessAreas store.ProtectedAreas
}

func (Ω *wildernessAreaFetcher) GetWildernessArea(cxt context.Context, latitude, longitude float64) (*store.ProtectedArea, error) {

	for _, w := range Ω.WildernessAreas {
		if utils.CoordinatesEqual(w.Centre.Latitude, latitude) && utils.CoordinatesEqual(w.Centre.Longitude, longitude) {
			return &w, nil
		}
	}

	w, err := Ω.Store.ReadProtectedAreaByLatLng(cxt, latitude, longitude)
	if err != nil && strings.Contains(err.Error(), "no wilderness area found") {
		fmt.Println("could not find wilderness area", latitude, longitude)
		w = &store.ProtectedArea{
			Centre: latlng.LatLng{latitude, longitude},
		}
	} else if err != nil {
		return nil, err
	}

	Ω.Lock()
	defer Ω.Unlock()

	Ω.WildernessAreas = append(Ω.WildernessAreas, *w)

	return w, nil
}
