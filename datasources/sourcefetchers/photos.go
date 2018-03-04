package sourcefetchers

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/datasources/gbif"
	"bitbucket.org/heindl/process/datasources/inaturalist"
	"context"
	"github.com/dropbox/godropbox/errors"
)

type Photo interface {
	Citation() string
	Thumbnail() string
	Large() string
	Source() datasources.SourceType
}

func FetchPhotos(ctx context.Context, sourceTypeProvider datasources.SourceTypeProvider, targetIDProvider datasources.TargetIDProvider) ([]Photo, error) {

	sourceType, err := sourceTypeProvider()
	if err != nil {
		return nil, err
	}

	targetID, err := targetIDProvider()
	if err != nil {
		return nil, err
	}

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
