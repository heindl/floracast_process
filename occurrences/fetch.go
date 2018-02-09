package occurrences

import (
	"bitbucket.org/heindl/taxa/datasources/inaturalist"
	"bitbucket.org/heindl/taxa/datasources/mushroomobserver"
	"bitbucket.org/heindl/taxa/datasources"
	"bitbucket.org/heindl/taxa/datasources/gbif"
	"context"
	"time"
	"github.com/dropbox/godropbox/errors"
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
	case datasources.DataSourceTypeGBIF:
		providers, err := gbif.FetchOccurrences(ctx, targetID, since)
		if err != nil {
			return nil, err
		}
		for _, _provider := range providers {
			provided = append(provided, _provider)
		}
	case datasources.DataSourceTypeINaturalist:
		providers, err := inaturalist.FetchOccurrences(ctx, targetID, since)
		if err != nil {
			return nil, err
		}
		for _, _provider := range providers {
			provided = append(provided, _provider)
		}
	case datasources.DataSourceTypeMushroomObserver:
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

	aggregation := OccurrenceAggregation{}

	for _, p := range provided {
		o, err := NewOccurrence(sourceType, targetID, p.SourceOccurrenceID())
		if err != nil {
			return nil, err
		}
		lng, err := p.Lng()
		if err != nil {
			return nil, errors.Wrap(err, "Could not get Occurrence Longitude")
		}
		lat, err := p.Lat()
		if err != nil {
			return nil, errors.Wrap(err, "Could not get Occurrence Latitude")
		}
		if err := o.SetGeospatial(lat, lng, p.DateString(), p.CoordinatesEstimated()); err != nil {
			return nil, errors.Wrap(err, "Could not set Occurrence Geospatial")
		}
		if err := aggregation.AddOccurrence(o); err != nil && err != ErrCollision {
			return nil, err
		}
	}

	return &aggregation, nil
}
