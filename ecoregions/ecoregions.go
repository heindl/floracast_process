package ecoregions

import (
	"bitbucket.org/heindl/taxa/ecoregions/generated_cache"
	"github.com/saleswise/errors/errors"
	"github.com/tidwall/tile38/geojson"
	"gopkg.in/tomb.v2"
)

type EcoRegionsCache []EcoRegion

type EcoRegion struct {
	EcoName    string  // Ecoregion Name
	EcoCode    EcoCode // This is an alphanumeric code that is similar to eco_ID but a little easier to interpret. The first 2 characters (letters) are the realm the ecoregion is in. The 2nd 2 characters are the biome and the last 2 characters are the ecoregion number.
	EcoNum     EcoNum  // A unique number for each ecoregion within each biome nested within each realm.
	EcoID      EcoID   // This number is created by combining REALM, BIOME, and ECO_NUM, thus creating a unique numeric ID for each ecoregion.
	Biome      Biome
	GeoObjects []geojson.Object
}

func NewEcoRegionsCache() (EcoRegionsCache, error) {
	res := EcoRegionsCache{}
	for _, cr := range cache.EcoRegionCache {

		nr := EcoRegion{
			EcoCode: EcoCode(cr.EcoCode),
			EcoNum:  EcoNum(cr.EcoNum),
			EcoID:   EcoID(cr.EcoID),
			Biome:   Biome(cr.Biome),
		}

		var ok bool
		if nr.EcoName, ok = EcoIDDefinitions[nr.EcoID]; !ok {
			continue
		}
		if _, ok := BiomeDefinitions[nr.Biome]; !ok {
			continue
		}

		for _, geoStr := range cr.Geometries {
			geoObj, err := geojson.ObjectJSON(geoStr)
			if err != nil {
				return nil, errors.Wrapf(err, "could not parse geojson string for region[%s] in biome[%s]", cr.EcoCode, cr.Biome)
			}
			nr.GeoObjects = append(nr.GeoObjects, geoObj)
		}

		res = append(res, nr)
	}
	return res, nil
}

func (Ω EcoRegionsCache) EcoID(lat, lng float64) EcoID {
	p := geojson.New2DPoint(lng, lat)

	var id EcoID

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _i := range Ω {
			i := _i
			_o := Ω[i]
			tmb.Go(func() error {
				o := _o
				for _, geoObj := range o.GeoObjects {
					if p.Within(geoObj) {
						id = o.EcoID
						tmb.Kill(nil)
					}
				}
				return nil
			})
		}
		return nil
	})

	_ = tmb.Wait()

	return id
}
