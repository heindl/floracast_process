package occurrences

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/datasources/gbif"
	"bitbucket.org/heindl/process/datasources/inaturalist"
	"bitbucket.org/heindl/process/datasources/mushroomobserver"
	"bitbucket.org/heindl/process/terra/ecoregions"
	"bitbucket.org/heindl/process/terra/geo"
	"bitbucket.org/heindl/process/utils"
	"context"
	"github.com/dropbox/godropbox/errors"
	"gopkg.in/tomb.v2"
	"time"
)

// OccurrenceProvider is a standard interface for sources that fetch occurrences.
type OccurrenceProvider interface {
	Lat() (float64, error)
	Lng() (float64, error)
	DateString() string
	CoordinatesEstimated() bool
	SourceOccurrenceID() string
}

// FetchOccurrences from datasource with given SourceType and TargetID.
func FetchOccurrences(ctx context.Context, sourceType datasources.SourceType, targetID datasources.TargetID, since *time.Time) (*OccurrenceAggregation, error) {

	// Only fetch once a day.
	if since != nil && since.After(time.Now().Add(time.Hour*24*-1)) {
		return nil, nil
	}

	providers, err := fetchOccurrencesFromSource(ctx, sourceType, targetID, since)
	if err != nil {
		return nil, err
	}

	aggregation := OccurrenceAggregation{}

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, 𝝨 := range providers {
			provided := 𝝨
			tmb.Go(func() error {
				return parseOccurrenceProvider(sourceType, targetID, provided, &aggregation)
			})
		}
		return nil
	})
	if err := tmb.Wait(); err != nil {
		return nil, err
	}

	return &aggregation, nil
}

func fetchOccurrencesFromSource(ctx context.Context, sourceType datasources.SourceType, targetID datasources.TargetID, since *time.Time) ([]OccurrenceProvider, error) {
	res := []OccurrenceProvider{}
	switch sourceType {
	case datasources.TypeGBIF:
		gvn, err := gbif.FetchOccurrences(ctx, targetID, since)
		for i := range gvn {
			res = append(res, gvn[i])
		}
		return res, err
	case datasources.TypeINaturalist:
		gvn, err := inaturalist.FetchOccurrences(ctx, targetID, since)
		for i := range gvn {
			res = append(res, gvn[i])
		}
		return res, err
	case datasources.TypeMushroomObserver:
		gvn, err := mushroomobserver.FetchOccurrences(ctx, targetID, since)
		for i := range gvn {
			res = append(res, gvn[i])
		}
		return res, err
	default:
		return nil, errors.Newf("Unsupported SourceType [%s]", sourceType)
	}
}

func parseOccurrenceProvider(sourceType datasources.SourceType, targetID datasources.TargetID, provided OccurrenceProvider, aggr *OccurrenceAggregation) error {
	o, err := NewOccurrence(sourceType, targetID, provided.SourceOccurrenceID())
	if err != nil {
		return err
	}

	lat, latErr := provided.Lat()
	lng, lngErr := provided.Lng()
	if latErr != nil || lngErr != nil {
		return errors.Wrap(err, "Invalid Coordinate")
	}

	err = o.SetGeoSpatial(lat, lng, provided.DateString(), provided.CoordinatesEstimated())
	if utils.ContainsError(err, geo.ErrInvalidCoordinates) ||
		utils.ContainsError(err, ErrInvalidDate) ||
		utils.ContainsError(err, ecoregions.ErrNotFound) {
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "Invalid Occurrence GeoSpatial")
	}

	err = aggr.AddOccurrence(o)
	if err != nil && !utils.ContainsError(err, ErrCollision) {
		return err
	}
	return nil
}
