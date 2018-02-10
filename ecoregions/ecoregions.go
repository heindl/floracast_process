package ecoregions

import (
	"bitbucket.org/heindl/processors/terra"
	"gopkg.in/tomb.v2"
	"sync"
)

type EcoRegionsCache struct {
	fc *terra.FeatureCollection
}

type EcoRegion struct {
	EcoName      string  // Ecoregion Name
	EcoCode      EcoCode // This is an alphanumeric code that is similar to eco_ID but a little easier to interpret. The first 2 characters (letters) are the realm the ecoregion is in. The 2nd 2 characters are the biome and the last 2 characters are the ecoregion number.
	EcoNum       EcoNum  // A unique number for each ecoregion within each biome nested within each realm.
	EcoID        EcoID   // This number is created by combining REALM, BIOME, and ECO_NUM, thus creating a unique numeric ID for each ecoregion.
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

	fc := terra.FeatureCollection{}
	if err := fc.Append(fc_holder...); err != nil {
		return nil, err
	}

	return &EcoRegionsCache{&fc}, nil
}

func (Ω *EcoRegionsCache) EcoID(lat, lng float64) EcoID {

	var id EcoID

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, _feature := range Ω.fc.Features() {
			feature := _feature
			tmb.Go(func() error {
				if feature.Contains(lat, lng) {
					b, _ := feature.GetPropertyInt("ECO_ID")
					id = EcoID(b)
				}
				return nil
			})
		}
		return nil
	})

	_ = tmb.Wait()

	return id
}
