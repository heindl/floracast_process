package sourcefetchers

import (
	"context"
	"bitbucket.org/heindl/processors/datasources"
	"bitbucket.org/heindl/processors/nameusage/nameusage"
	"bitbucket.org/heindl/processors/datasources/inaturalist"
	"bitbucket.org/heindl/processors/datasources/gbif"
	"bitbucket.org/heindl/processors/datasources/natureserve"
	"github.com/dropbox/godropbox/errors"
)

func FetchNameUsages(ctx context.Context, sourceType datasources.SourceType, scientificNames []string, targetIDs datasources.TargetIDs) ([]*nameusage.NameUsage, error) {
	switch sourceType {
	case datasources.TypeGBIF:
		return gbif.FetchNamesUsages(ctx, scientificNames, targetIDs)
	case datasources.TypeINaturalist:
		return inaturalist.FetchNameUsages(ctx, scientificNames, targetIDs)
	case datasources.TypeNatureServe:
		return natureserve.FetchNameUsages(ctx, scientificNames, targetIDs)
	case datasources.TypeMushroomObserver:
		return natureserve.FetchNameUsages(ctx, scientificNames, targetIDs)
	default:
		return nil, errors.Newf("Unsupported SourceType [%s] for NameUsages", sourceType)
	}
}