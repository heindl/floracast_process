package sourcefetchers

import (
	"bitbucket.org/heindl/processors/datasources/inaturalist"
	"bitbucket.org/heindl/processors/datasources"
	"bitbucket.org/heindl/processors/datasources/gbif"
	"context"
	"github.com/dropbox/godropbox/errors"
)

type Photo interface {
	Citation() string
	Thumbnail() string
	Large() string
	Source() datasources.SourceType
}

func FetchPhotos(ctx context.Context, sourceType datasources.SourceType, targetID datasources.TargetID) ([]Photo, error) {

	res := []Photo{}
	switch sourceType {
	case datasources.TypeGBIF:
		photos, err := gbif.FetchPhotos(ctx, targetID)
		if err != nil {
			return nil, err
		}
		for _, o := range photos {
			res = append(res, o)
		}
	case datasources.TypeINaturalist:
		photos, err := inaturalist.FetchPhotos(ctx, targetID)
		if err != nil {
			return nil, err
		}
		for _, o := range photos {
			res = append(res, o)
		}
	default:
		return nil, errors.Newf("Unsupported SourceType [%s]", sourceType)
	}

	return res, nil

}
