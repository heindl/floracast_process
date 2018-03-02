package sourcefetchers

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/datasources/gbif"
	"bitbucket.org/heindl/process/datasources/inaturalist"
	"bitbucket.org/heindl/process/datasources/mushroomobserver"
	"context"
	"github.com/dropbox/godropbox/errors"
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

func FetchOccurrences(ctx context.Context, sourceType datasources.SourceType, targetID datasources.TargetID, since *time.Time) ([]OccurrenceProvider, error) {
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
