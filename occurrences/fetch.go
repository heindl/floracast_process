package occurrences

import (
	"bitbucket.org/heindl/processors/datasources/inaturalist"
	"bitbucket.org/heindl/processors/datasources/mushroomobserver"
	"bitbucket.org/heindl/processors/datasources"
	"bitbucket.org/heindl/processors/datasources/gbif"
	"context"
	"time"
	"github.com/dropbox/godropbox/errors"
	"bitbucket.org/heindl/processors/utils"
	"bitbucket.org/heindl/processors/geofeatures"
	"gopkg.in/tomb.v2"
	"fmt"
	"bitbucket.org/heindl/processors/ecoregions"
)

type OccurrenceProvider interface {
	Lat() (float64, error)
	Lng() (float64, error)
	DateString() string
	CoordinatesEstimated() bool
	SourceOccurrenceID() string
}

func FetchOccurrences(ctx context.Context, sourceType datasources.SourceType, targetID datasources.TargetID, since *time.Time) (*OccurrenceAggregation, error) {

	// Only fetch once a day.
	if since != nil && since.After(time.Now().Add(time.Hour * 24 * -1)) {
		return nil, nil
	}

	provided := []OccurrenceProvider{}

	switch sourceType {
	case datasources.TypeGBIF:
		providers, err := gbif.FetchOccurrences(ctx, targetID, since)
		if err != nil {
			return nil, err
		}
		for _, _provider := range providers {
			provided = append(provided, _provider)
		}
	case datasources.TypeINaturalist:
		providers, err := inaturalist.FetchOccurrences(ctx, targetID, since)
		if err != nil {
			return nil, err
		}
		for _, _provider := range providers {
			provided = append(provided, _provider)
		}
	case datasources.TypeMushroomObserver:
		providers, err := mushroomobserver.FetchOccurrences(ctx, targetID, since)
		if err != nil {
			return nil, err
		}
		for _, _provider := range providers {
			provided = append(provided, _provider)
		}
	default:
		return nil, errors.Newf("Unsupported SourceType [%s]", sourceType)
	}

	if len(provided) == 0 {
		return nil, nil
	}

	aggregation := OccurrenceAggregation{}

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, _p := range provided {
			p := _p
			tmb.Go(func() error {
				o, err := NewOccurrence(sourceType, targetID, p.SourceOccurrenceID())
				if err != nil {
					return err
				}
				lng, err := p.Lng()
				if err != nil {
					return errors.Wrap(err, "Could not get Occurrence Longitude")
				}
				lat, err := p.Lat()
				if err != nil {
					return errors.Wrap(err, "Could not get Occurrence Latitude")
				}

				err = o.SetGeospatial(lat, lng, p.DateString(), p.CoordinatesEstimated())
				if err != nil && utils.ContainsError(err, geofeatures.ErrInvalidCoordinate) {
					//fmt.Println(fmt.Sprintf("Invalid Coordinate [%.4f, %.4f] from SourceType [%s, %s]", lat, lng, sourceType, targetID))
					return nil
				}
				if err != nil && utils.ContainsError(err, ErrInvalidDate) {
					return nil
				}
				if err != nil && utils.ContainsError(err, ecoregions.ErrNotFound) {
					fmt.Println(fmt.Sprintf("Invalid EcoRegion [%.4f, %.4f]", lat, lng))
					return nil
				}
				if err != nil {
					return errors.Wrap(err, "Could not set Occurrence Geospatial")
				}
				err = aggregation.AddOccurrence(o)
				if err != nil && err == ErrCollision {
					return nil
				}
				if err != nil {
					return err
				}
				return nil
			})
		}
		return nil
	})
	if err := tmb.Wait(); err != nil {
		return nil, err
	}

	return &aggregation, nil
}
