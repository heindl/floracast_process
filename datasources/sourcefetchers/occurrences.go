package sourcefetchers

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/datasources/gbif"
	"bitbucket.org/heindl/process/datasources/inaturalist"
	"bitbucket.org/heindl/process/datasources/mushroomobserver"
	"bitbucket.org/heindl/process/datasources/providers"
	"context"
	"github.com/dropbox/godropbox/errors"
	"time"
)

// FetchOccurrences returns occurrences from the requested source.
func FetchOccurrences(ctx context.Context, sourceType datasources.SourceType, targetID datasources.TargetID, since *time.Time) ([]providers.Occurrence, error) {
	switch sourceType {
	case datasources.TypeGBIF:
		return gbif.FetchOccurrences(ctx, targetID, since)
	case datasources.TypeINaturalist:
		return inaturalist.FetchOccurrences(ctx, targetID, since)
	case datasources.TypeMushroomObserver:
		return mushroomobserver.FetchOccurrences(ctx, targetID, since)
	default:
		return nil, errors.Newf("Unsupported SourceType [%s]", sourceType)
	}
}
