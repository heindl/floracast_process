package sourcefetchers

import (
	"github.com/heindl/floracast_process/datasources"
	"github.com/heindl/floracast_process/datasources/gbif"
	"github.com/heindl/floracast_process/datasources/inaturalist"
	"github.com/heindl/floracast_process/datasources/mushroomobserver"
	"github.com/heindl/floracast_process/datasources/providers"
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
		return inaturalist.FetchOccurrences(ctx, targetID, since, true)
	case datasources.TypeMushroomObserver:
		return mushroomobserver.FetchOccurrences(ctx, targetID, since)
	default:
		return nil, errors.Newf("Unsupported SourceType [%s]", sourceType)
	}
}
