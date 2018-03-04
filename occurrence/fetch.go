package occurrence

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/datasources/sourcefetchers"
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/terra/ecoregions"
	"bitbucket.org/heindl/process/terra/geo"
	"bitbucket.org/heindl/process/utils"
	"context"
	"github.com/dropbox/godropbox/errors"
	"github.com/golang/glog"
	"gopkg.in/tomb.v2"
	"time"
)

// FetchOccurrences takes sources from NameUsage, fetches occurrences for each, and updates the fetch time.
func FetchOccurrences(ctx context.Context, usage nameusage.NameUsage, limitToCount bool, sourceTypes ...datasources.SourceType) (*Aggregation, error) {

	if len(sourceTypes) == 0 {
		sourceTypes = []datasources.SourceType{
			datasources.TypeINaturalist,
			datasources.TypeGBIF,
			datasources.TypeMushroomObserver,
		}
	}

	srcs, err := usage.Sources(sourceTypes...)
	if err != nil {
		return nil, err
	}

	id, err := usage.ID()
	if err != nil {
		return nil, err
	}

	glog.Infof("Fetching Occurrences for NameUsage [%s, %s] with %d Sources", usage.CanonicalName(), id, len(srcs))

	res := &Aggregation{}

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, 𝝨 := range srcs {
			_src := 𝝨
			tmb.Go(func() error {
				src := _src
				if limitToCount && src.OccurrenceCount() == 0 {
					return nil
				}
				return fetchAndMerge(ctx, src, res)
			})
		}
		return nil
	})
	if err := tmb.Wait(); err != nil {
		return nil, err
	}

	glog.Infof("%d Occurrences Aggregated for NameUsage [%s, %s] with %d Sources", res.Count(), usage.CanonicalName(), id, len(srcs))

	return res, nil

}

func fetchAndMerge(ctx context.Context, src nameusage.Source, parentAggregation *Aggregation) error {

	srcType, err := src.SourceType()
	if err != nil {
		return err
	}

	targetID, err := src.TargetID()
	if err != nil {
		return err
	}

	aggr, err := fetchOccurrencesForTarget(ctx, srcType, targetID, src.LastFetchedAt())
	if err != nil {
		return err
	}

	if err := src.RegisterOccurrenceFetch(aggr.Count()); err != nil {
		return err
	}

	return parentAggregation.Merge(aggr)
}

func fetchOccurrencesForTarget(ctx context.Context, sourceType datasources.SourceType, targetID datasources.TargetID, since *time.Time) (*Aggregation, error) {

	glog.Infof("Fetching Occurrences [%s, %s] since %v", sourceType, targetID, since)

	// Only fetch once a day.
	if since != nil && since.After(time.Now().Add(time.Hour*24*-1)) {
		return nil, nil
	}

	providers, err := sourcefetchers.FetchOccurrences(ctx, sourceType, targetID, since)
	if err != nil {
		return nil, err
	}

	glog.Infof("Received %d Occurrences Providers for Source [%s, %s] Since %v", len(providers), sourceType, targetID, since)

	aggregation := Aggregation{}

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

	glog.Infof("Processed %d Occurrences for NameUsage Source [%s, %s]", aggregation.Count(), sourceType, targetID)

	return &aggregation, nil
}

func parseOccurrenceProvider(sourceType datasources.SourceType, targetID datasources.TargetID, provided sourcefetchers.OccurrenceProvider, aggr *Aggregation) error {
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

	if err := aggr.AddOccurrence(o); err != nil && !utils.ContainsError(err, ErrCollision) {
		return err
	}
	return nil
}
