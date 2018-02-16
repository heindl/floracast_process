package geofeatures

import (
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"bitbucket.org/heindl/processors/utils"
	"strings"
	"google.golang.org/genproto/googleapis/type/latlng"
)

func elevationKey(lat, lng float64) string {
	return fmt.Sprintf("%.4f,%.4f", lat, lng)
}

func (Ω *geoFeaturesProcessor) elevationQueueStatus(lat, lng float64) (queued, fetched bool) {
	Ω.Lock()
	defer Ω.Unlock()
	k := elevationKey(lat, lng)
	_, queued = Ω.elevationsQueued[k]
	fetched = false
	if e, ok := Ω.elevationsFetched[k]; ok && e != nil {
		fetched = true
	}
	return
}

func (Ω *geoFeaturesProcessor) queueElevation(lat, lng float64) error {
	queued, fetched := Ω.elevationQueueStatus(lat, lng)
	if queued || fetched {
		return nil
	}
	k := elevationKey(lat, lng)
	Ω.Lock()
	Ω.elevationsQueued[k] = &latlng.LatLng{Latitude: lat, Longitude: lng}
	Ω.Unlock()
	if len(Ω.elevationsQueued) >= 10 {
		return Ω.flushElevations()
	}
	return nil
}

func (Ω *geoFeaturesProcessor) getElevation(lat, lng float64) (*float64, error) {
	queued, fetched := Ω.elevationQueueStatus(lat, lng)

	if !fetched && queued {
		return nil, nil
	}
	if !fetched && !queued {
		return nil, errors.Newf("Trying to get coordinates neither fetched or queued [%s] with %d previously fetched", elevationKey(lat, lng), len(Ω.elevationsFetched))
	}
	if e, ok := Ω.elevationsFetched[elevationKey(lat, lng)]; !ok || e == nil {
		return nil, nil
	} else {
		return e, nil
	}
}

func (Ω *geoFeaturesProcessor) flushElevations() error {
	Ω.Lock()
	defer Ω.Unlock()

	locs := []string{}
	for _, k := range Ω.elevationsQueued {
		locs = append(locs, elevationKey(k.GetLatitude(), k.GetLongitude()))
	}
	if len(locs) == 0 {
		return nil
	}

	var res struct {
		Results []struct {
			Lat float64 `json:"latitude"`
			Lng float64 `json:"longitude"`
			Elevation float64 `json:"elevation"` // Meters
		} `json:"results"`
	}

	if err := utils.RequestJSON("https://api.open-elevation.com/api/v1/lookup?locations=" + strings.Join(locs, "|"), &res); err != nil {
		return errors.Wrap(err, "Could not fetch elevation api")
	}

	//resolvedElevations, err := Ω.mapClient.elevation(context.Background(), &eleReq)
	//if err != nil {
	//	return errors.Wrap(err, "could not fetch elevations")
	//}
	//for _, e := range resolvedElevations {
	//	k := elevationKey(e.Location.Lat, e.Location.Lng)
	//	Ω.elevationsFetched[k] = e.elevation
	//}
	for _, r := range res.Results {
		k := elevationKey(r.Lat, r.Lng)
		if r.Elevation == 0 {
			fmt.Println(fmt.Sprintf("Elevation value is 0 for key [%s], so may not resolved. Will add it regardless.", k))
		}
		Ω.elevationsFetched[k] = utils.FloatPtr(r.Elevation)
	}

	//fmt.Println("Flushing Elevations", len(Ω.elevationsQueued), len(Ω.elevationsFetched), len(locs), len(resolvedElevations))
	Ω.elevationsQueued = map[string]*latlng.LatLng{}
	return nil
}
//
//
//func (Ω *geoFeaturesProcessor) processElevationBatch(cxt context.Context, locations ...PredictableLocation) error {
//
//
//	// First batch and fetch elevations
//	eleReq := maps.ElevationRequest{Locations: []maps.LatLng{}}
//	// Gather lat/lng pairs for elevation fetch.
//	for _, o := range locations {
//		eleReq.Locations = append(eleReq.Locations, maps.LatLng{o.Lat(), o.Lng()})
//	}
//
//	//resolvedElevations := []maps.ElevationResult{}
//
//	//bkf := backoff.NewExponentialBackOff()
//	//bkf.InitialInterval = time.Second * 1
//	//ticker := backoff.NewTicker(bkf)
//	//for _ = range ticker.C {
//
//
//		resolvedElevations, err := Ω.mapClient.elevation(cxt, &eleReq)
//		//if err != nil && strings.Contains(err.Error(), "TLS handshake timeout") {
//		//	fmt.Println("TLS handshake timeout encountered. Backing off ...")
//		//	continue
//		//}
//		//if err != nil && strings.Contains(err.Error(), "DATA_NOT_AVAILABLE") {
//		//	fmt.Println("DATA_NOT_AVAILABLE", len(resolvedElevations))
//		//	ticker.Stop()
//		//	break
//		//}
//		if err != nil {
//			//ticker.Stop()
//			return errors.Wrap(err, "could not fetch elevations")
//		}
//
//	//	ticker.Stop()
//	//	break
//	//}
//
//	tmb := tomb.Tomb{}
//	tmb.Go(func() error {
//		for _, _loc := range locations {
//			loc := _loc
//			tmb.Go(func() error {
//				for _, _r := range resolvedElevations {
//					r := _r
//					if !utils.CoordinatesEqual(loc.Lat(), r.Location.Lat) {
//						continue
//					}
//					if !utils.CoordinatesEqual(loc.Lng(), r.Location.Lng) {
//						continue
//					}
//					return loc.SetElevation(r.elevation)
//				}
//				return errors.Newf("elevation not found: %.5f, %.5f", loc.Lat(), loc.Lng())
//			})
//		}
//		return nil
//	})
//
//	return tmb.Wait()
//
//}

