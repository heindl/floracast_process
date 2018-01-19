package store

import (
	"bitbucket.org/heindl/taxa/ecoregions"
	"bitbucket.org/heindl/taxa/terra"
	"bitbucket.org/heindl/taxa/utils"
	"context"
	"github.com/dropbox/godropbox/errors"
	"googlemaps.github.io/maps"
	"os"
	"gopkg.in/tomb.v2"
	"fmt"
	"github.com/cenkalti/backoff"
	"time"
	"strings"
)

type PredictableLocation interface {
	Lat() float64
	Lng() float64
	SetGeoFeatures(*GeoFeatures)
}

type GeoFeatures struct {
	Latitude  float64           `firestore:"" json:""`
	Longitude float64           `firestore:"" json:""`
	Biome     ecoregions.Biome  `firestore:"" json:""`
	EcoNum    ecoregions.EcoNum `firestore:"" json:""`
	S2Tokens  []string          `firestore:"" json:""` // Ordered from 0 to 10.
	Elevation *float64          `firestore:"" json:""`
}

func (Ω *GeoFeatures) Valid() bool {

	if !Ω.Biome.Valid() || !Ω.EcoNum.Valid() {
		return false
	}

	if Ω.Latitude == 0 || Ω.Longitude == 0 {
		return false
	}

	if len(Ω.S2Tokens) == 0 {
		return false
	}

	if Ω.Elevation == nil {
		return false
	}

	return true
}

func NewGeoFeaturesProcessor() (*GeoFeaturesProcessor, error) {

	gfp := &GeoFeaturesProcessor{}

	var err error
	gfp.mapClient, err = maps.NewClient(maps.WithAPIKey(os.Getenv("FLORACAST_GOOGLE_MAPS_API_KEY")))
	if err != nil {
		return nil, errors.Wrap(err, "could not get google maps client")
	}

	gfp.ecoRegionCache, err = ecoregions.NewEcoRegionsCache()
	if err != nil {
		return nil, err
	}

	return gfp, nil
}

type GeoFeaturesProcessor struct {
	mapClient      *maps.Client
	ecoRegionCache *ecoregions.EcoRegionsCache
}

func (Ω *GeoFeaturesProcessor) ProcessLocations(cxt context.Context, locations ...PredictableLocation) error {

	if len(locations) == 0 {
		return nil
	}

	// Officially 512 locations per request is the max.
	if len(locations) > 500 {
		return errors.New("Too many locations for which to fetch elevation. Max 500.")
	}

	// First batch and fetch elevations
	elevationRequest := maps.ElevationRequest{Locations: make([]maps.LatLng, len(locations))}
	// Gather lat/lng pairs for elevation fetch.
	for i, o := range locations {
		elevationRequest.Locations[i] = maps.LatLng{o.Lat(), o.Lng()}
	}

	resolvedElevations := []maps.ElevationResult{}

	bkf := backoff.NewExponentialBackOff()
	bkf.InitialInterval = time.Second * 1
	ticker := backoff.NewTicker(bkf)
	for _ = range ticker.C {

		var err error
		resolvedElevations, err = Ω.mapClient.Elevation(cxt, &elevationRequest)
		if err != nil && strings.Contains(err.Error(), "TLS handshake timeout") {
			fmt.Println("TLS handshake timeout encountered. Backing off ...")
			continue
		}
		if err != nil && strings.Contains(err.Error(), "DATA_NOT_AVAILABLE") {
			fmt.Println("DATA_NOT_AVAILABLE", len(resolvedElevations))
			ticker.Stop()
			break
		}
		if err != nil {
			ticker.Stop()
			return errors.Wrap(err, "could not fetch elevations")
		}

		ticker.Stop()
		break
	}

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, _loc := range locations {
			loc := _loc
			tmb.Go(func() error {
				for _, _r := range resolvedElevations {
					r := _r
					if !utils.CoordinatesEqual(loc.Lat(), r.Location.Lat) {
						continue
					}
					if !utils.CoordinatesEqual(loc.Lng(), r.Location.Lng) {
						continue
					}
					Ω.setFeatures(loc, r.Elevation)
					return nil
				}
				return errors.Newf("Elevation not found: %.5f, %.5f", loc.Lat(), loc.Lng())
			})
		}
		return nil
	})

	return tmb.Wait()


}

func (Ω *GeoFeaturesProcessor) setFeatures(loc PredictableLocation, elevation float64) error {

	ecoID := Ω.ecoRegionCache.EcoID(loc.Lat(), loc.Lng())
	if !ecoID.Valid() {
		return nil // Don't fail entire batch for one. Rely on external validity check.
	}

	loc.SetGeoFeatures(&GeoFeatures{
		Latitude:  loc.Lat(),
		Longitude: loc.Lng(),
		Biome:     ecoID.Biome(),
		EcoNum:    ecoID.EcoNum(),
		Elevation: utils.FloatPtr(elevation),
		S2Tokens:  terra.NewPoint(loc.Lat(), loc.Lng()).S2TokenArray(),
	})

	return nil
}
