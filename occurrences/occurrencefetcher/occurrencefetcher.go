package occurrencefetcher

import (
	"bitbucket.org/heindl/taxa/occurrences"
	"bitbucket.org/heindl/taxa/datasources/inaturalist"
	"bitbucket.org/heindl/taxa/datasources/mushroomobserver"
	"bitbucket.org/heindl/taxa/datasources"
	"bitbucket.org/heindl/taxa/datasources/gbif"
	"context"
	"time"
	"github.com/dropbox/godropbox/errors"
)

type OccurrenceFetcher func(context.Context, datasources.DataSourceTargetID, *time.Time) (*occurrences.OccurrenceAggregation, error)

func FetchOccurrences(ctx context.Context, sourceType datasources.DataSourceType, targetID datasources.DataSourceTargetID, since *time.Time) (*occurrences.OccurrenceAggregation, error) {

	// Only fetch once a day.
	if since != nil && since.After(time.Now().Add(time.Hour * 24 * -1)) {
		return nil, nil
	}

	var list *occurrences.OccurrenceAggregation
	var err error
	switch sourceType {
	case datasources.DataSourceTypeGBIF:
		list, err = gbif.FetchOccurrences(ctx, targetID, since)
	case datasources.DataSourceTypeINaturalist:
		list, err = inaturalist.FetchOccurrences(ctx, targetID, since)
	case datasources.DataSourceTypeMushroomObserver:
		list, err = mushroomobserver.FetchOccurrences(ctx, targetID, since)
	default:
		return nil, errors.Newf("Unsupported SourceType [%s]", sourceType)
	}

	return list, err
}
