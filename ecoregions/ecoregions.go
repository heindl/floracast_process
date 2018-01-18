package ecoregions

import (
	"bitbucket.org/heindl/taxa/ecoregions/cache"
	"github.com/saleswise/errors/errors"
	"gopkg.in/tomb.v2"
	"bitbucket.org/heindl/taxa/terra"
)



type EcoRegionsCache []*EcoRegion

type EcoRegion struct {
	EcoName    string  // Ecoregion Name
	EcoCode    EcoCode // This is an alphanumeric code that is similar to eco_ID but a little easier to interpret. The first 2 characters (letters) are the realm the ecoregion is in. The 2nd 2 characters are the biome and the last 2 characters are the ecoregion number.
	EcoNum     EcoNum  // A unique number for each ecoregion within each biome nested within each realm.
	EcoID      EcoID   // This number is created by combining REALM, BIOME, and ECO_NUM, thus creating a unique numeric ID for each ecoregion.
	Biome      Biome
	MultiPolygon terra.MultiPolygon
}

func NewEcoRegionsCache() (EcoRegionsCache, error) {
	res := EcoRegionsCache{}
	for _, cr := range cache.EcoRegionCache {

		nr := EcoRegion{
			EcoCode: EcoCode(cr.EcoCode),
			EcoNum:  EcoNum(cr.EcoNum),
			EcoID:   EcoID(cr.EcoID),
			Biome:   Biome(cr.Biome),
			MultiPolygon: terra.MultiPolygon{},
		}

		//var ok bool
		//if nr.EcoName, ok = EcoIDDefinitions[nr.EcoID]; !ok {
		//	fmt.Println(nr.EcoID)
		//	continue
		//}
		//if _, ok := BiomeDefinitions[nr.Biome]; !ok {
		//	fmt.Println(nr.Biome)
		//	continue
		//}

		for _, geo := range cr.Geometries {
			p, err := terra.NewPolygon(geo)
			if err != nil {
				return nil, errors.Wrap(err, "could not create polygon")
			}
			nr.MultiPolygon = nr.MultiPolygon.PushPolygon(p)
		}

		res = append(res, &nr)
	}
	return res, nil
}

func (Ω EcoRegionsCache) EcoID(lat, lng float64) EcoID {

	var id EcoID

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _i := range Ω {
			i := _i
			_o := Ω[i]
			tmb.Go(func() error {
				o := _o
				if o.MultiPolygon.Contains(lat, lng) {
					id = o.EcoID
				}
				return nil
			})
		}
		return nil
	})

	_ = tmb.Wait()

	return id
}
