package main

import (
	"os"
	"bitbucket.org/heindl/taxa/store"
	"googlemaps.github.io/maps"
	"github.com/saleswise/errors/errors"
	"gopkg.in/tomb.v2"
	"context"
)

var EPSILON float64 = 0.00001

func coordinateEquals(a, b float64) bool {
	if ((a - b) < EPSILON && (b - a) < EPSILON) {
		return true
	}
	return false
}

func setElevations(occurrences store.Occurrences) error {

	if len(occurrences) == 0 {
		return nil
	}

	mc, err := maps.NewClient(maps.WithAPIKey(os.Getenv("FLORACAST_GOOGLE_MAPS_API_KEY")))
	if err != nil {
		return errors.Wrap(err, "could not get google maps client")
	}

	start := 0
	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for {
			end := start + 100
			if len(occurrences) <= end {
				end = len(occurrences)
			}
			_list := occurrences[start:end]
			tmb.Go(func() error {
				list := _list
				locations := make([]maps.LatLng, len(list))
				// Gather lat/lng pairs for elevation fetch.
				for i, o := range list {
					locations[i] = maps.LatLng{o.Location.GetLatitude(), o.Location.GetLongitude()}
				}
				res, err := mc.Elevation(context.Background(), &maps.ElevationRequest{Locations: locations})
				if err != nil {
					return errors.Wrap(err, "could not fetch elevations")
				}
			Occurrences:
				for i := range list {
					for _, r := range res {
						if !coordinateEquals(list[i].Location.GetLatitude(), r.Location.Lat) {
							continue
						}
						if !coordinateEquals(list[i].Location.GetLongitude(), r.Location.Lng) {
							continue
						}
						list[i].Elevation = r.Elevation
						continue Occurrences
					}
				}
				return nil
			})
			start = end
			if start >= len(occurrences) {
				break
			}
		}
		return nil
	})
	return tmb.Wait()
}