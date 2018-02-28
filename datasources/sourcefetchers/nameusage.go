package sourcefetchers

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/datasources/gbif"
	"bitbucket.org/heindl/process/datasources/inaturalist"
	"bitbucket.org/heindl/process/datasources/mushroomobserver"
	"bitbucket.org/heindl/process/datasources/natureserve"
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"context"
	"github.com/dropbox/godropbox/errors"
)

func FetchNameUsages(ctx context.Context, sourceType datasources.SourceType, scientificNames []string, targetIDs datasources.TargetIDs) ([]nameusage.NameUsage, error) {
	switch sourceType {
	case datasources.TypeGBIF:
		return gbif.FetchNamesUsages(ctx, scientificNames, targetIDs)
	case datasources.TypeINaturalist:
		return inaturalist.FetchNameUsages(ctx, scientificNames, targetIDs)
	case datasources.TypeNatureServe:
		return natureserve.FetchNameUsages(ctx, scientificNames, targetIDs)
	case datasources.TypeMushroomObserver:
		return mushroomobserver.FetchNameUsages(ctx, scientificNames, targetIDs)
	default:
		return nil, errors.Newf("Unsupported SourceType [%s] for NameUsages", sourceType)
	}
}
