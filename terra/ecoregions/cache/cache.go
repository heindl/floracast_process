package cache

import (
	"bitbucket.org/heindl/process/terra/ecoregions"
	"bitbucket.org/heindl/process/terra/geo"
	"github.com/dropbox/godropbox/errors"
	"gopkg.in/tomb.v2"
	"sync"
)

var regionCache []*geo.Feature

func init() {

	m := sync.Mutex{}

	regionCache = []*geo.Feature{}

	callback := func(f *geo.Feature) error {
		m.Lock()
		defer m.Unlock()
		regionCache = append(regionCache, f)
		return nil
	}

	if err := geo.ParseGeoJSONFeatureCollection([]byte(ecoregionsGeoJson), callback); err != nil {
		panic(err)
	}
}

// FetchEcologicalRegion returns a region that contains a coordinate from in-memory cache
// If coordinate does not fall within region, return ErrNotFound
func FetchEcologicalRegion(lat, lng float64) (*ecoregions.Region, error) {

	var region *ecoregions.Region

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, 𝝨 := range regionCache {
			f := 𝝨
			tmb.Go(func() error {
				if f.Contains(lat, lng) {
					i, err := f.GetPropertyInt("ECO_ID")
					if err != nil {
						return errors.Wrapf(err, "Could not get ECO_ID property [%.4f, %.4f]", lat, lng)
					}
					region, err = ecoregions.NewRegion(i)
					if err != nil {
						return err
					}
					tmb.Kill(nil)
				}
				return nil
			})
		}
		return nil
	})

	if err := tmb.Wait(); err != nil {
		return nil, err
	}

	if region == nil {
		return nil, ecoregions.ErrNotFound
	}

	return region, nil
}
