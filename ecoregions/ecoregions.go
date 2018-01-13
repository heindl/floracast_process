package ecoregions

import (
	"bitbucket.org/heindl/taxa/ecoregions/generated_cache"
	"github.com/tidwall/tile38/geojson"
	"github.com/saleswise/errors/errors"
	"fmt"
	"bitbucket.org/heindl/taxa/utils"
)

type EcoRegions []EcoRegion

type EcoRegion struct {
	EcoName string // Ecoregion Name
	EcoCode EcoCode // This is an alphanumeric code that is similar to eco_ID but a little easier to interpret. The first 2 characters (letters) are the realm the ecoregion is in. The 2nd 2 characters are the biome and the last 2 characters are the ecoregion number.
	EcoNum int // A unique number for each ecoregion within each biome nested within each realm.
	EcoID EcoID // This number is created by combining REALM, BIOME, and ECO_NUM, thus creating a unique numeric ID for each ecoregion.
	Biome Biome
	GeoHashes []string
	GeoObject geojson.Object
}

func NewEcoRegionsCache() (EcoRegions, error) {
	res := EcoRegions{}
	for _, cr := range cache.EcoRegionCache {

		nr := EcoRegion{
			EcoCode: EcoCode(cr.EcoCode),
			EcoNum: int(cr.EcoNum),
			EcoID: EcoID(cr.EcoID),
			Biome: Biome(cr.Biome),
		}

		var ok bool
		if nr.EcoName, ok = EcoIDDefinitions[nr.EcoID]; !ok {
			continue
		}
		if _, ok := BiomeDefinitions[nr.Biome]; !ok {
			continue
		}

		var err error
		nr.GeoObject, err = geojson.ObjectJSON(cr.GeoJsonString)
		if err != nil {
			return nil, errors.Wrapf(err, "could not parse geojson string for region[%s] in biome[%s]", cr.EcoCode, cr.Biome)
		}

		res = append(res, nr)
	}
	return res, nil
}


func (Ω EcoRegions) HasPoint(lat, lng float64) (*EcoRegion, error) {
	p := geojson.New2DPoint(lng, lat)
	for i := range Ω {
		if p.Within(Ω[i].GeoObject) {
			return &Ω[i], nil
		}
	}
	return nil, nil
}