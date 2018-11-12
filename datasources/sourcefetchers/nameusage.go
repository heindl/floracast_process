package sourcefetchers

import (
	"github.com/heindl/floracast_process/datasources"
	"github.com/heindl/floracast_process/datasources/gbif"
	"github.com/heindl/floracast_process/datasources/inaturalist"
	"github.com/heindl/floracast_process/datasources/mushroomobserver"
	"github.com/heindl/floracast_process/datasources/natureserve"
	"github.com/heindl/floracast_process/nameusage/nameusage"
	"context"
	"github.com/dropbox/godropbox/errors"
)

// FetchNameUsages returns taxon NameUsages from the requested source.
func FetchNameUsages(ctx context.Context, sourceType datasources.SourceType, scientificNames []string, targetIDs datasources.TargetIDs) ([]nameusage.NameUsage, error) {
	switch sourceType {
	case datasources.TypeGBIF:
		return gbif.FetchNameUsages(ctx, scientificNames, targetIDs)
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
