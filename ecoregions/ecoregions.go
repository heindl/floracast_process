package ecoregions

import (
	"bitbucket.org/heindl/process/terra"
	"sync"
	"github.com/dropbox/godropbox/errors"
	"gopkg.in/tomb.v2"
)

type EcoRegionsCache struct {
	//fc *terra.FeatureCollection
	fc []*terra.Feature
}

type EcoRegion struct {
	Realm   Realm
	EcoName string  // Ecoregion Name
	EcoCode EcoCode // This is an alphanumeric code that is similar to eco_ID but a little easier to interpret. The first 2 characters (letters) are the realm the ecoregion is in. The 2nd 2 characters are the biome and the last 2 characters are the ecoregion number.
	EcoNum  EcoNum  // A unique number for each ecoregion within each biome nested within each realm.
	EcoID   EcoID   // This number is created by combining realm, BIOME, and ECO_NUM, thus creating a unique numeric ID for each ecoregion.
	Biome        Biome
	MultiPolygon terra.MultiPolygon
}

func NewEcoRegionsCache() (*EcoRegionsCache, error) {

	fc_holder := []*terra.Feature{}
	m := sync.Mutex{}

	callback := func(f *terra.Feature) error {
		m.Lock()
		defer m.Unlock()
		fc_holder = append(fc_holder, f)
		return nil
	}

	if err := terra.ParseGeoJSONFeatureCollection([]byte(ecoregions_geojson), callback); err != nil {
		return nil, err
	}
	//
	//fc := terra.FeatureCollection{}
	//if err := fc.Append(fc_holder...); err != nil {
	//	return nil, err
	//}

	return &EcoRegionsCache{fc_holder}, nil
}


var SignalDone = errors.New("EcoRegion Found")
var ErrNotFound = errors.New("EcoRegion Not Found")
func (Ω *EcoRegionsCache) EcoID(lat, lng float64) (EcoID, error) {

	var id EcoID

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, _f := range Ω.fc {
			f := _f
			tmb.Go(func() error {
				if f.Contains(lat, lng) {
					i, err := f.GetPropertyInt("ECO_ID")
					if err != nil {
						return errors.Wrapf(err, "Could not get ECO_ID property [%.4f, %.4f]", lat, lng)
					}
					id = EcoID(i)
					tmb.Kill(SignalDone)
				}
				return nil
			})
		}
		return nil
	})

	if err := tmb.Wait(); err != nil && err != SignalDone {
		return EcoID(0), err
	}

	if !id.Valid() {
		return EcoID(0), ErrNotFound
	}

	return id, nil
}
