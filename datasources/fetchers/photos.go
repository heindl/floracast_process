package fetchers

import (
	"bitbucket.org/heindl/taxa/datasources/inaturalist"
	"bitbucket.org/heindl/taxa/datasources"
	"bitbucket.org/heindl/taxa/datasources/gbif"
	"context"
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
	case datasources.DataSourceTypeGBIF:
		photos, err := gbif.FetchPhotos(ctx, targetID)
		if err != nil {
			return nil, err
		}
		for _, o := range photos {
			res = append(res, o)
		}
	case datasources.DataSourceTypeINaturalist:
		photos, err := inaturalist.FetchPhotos(ctx, targetID)
		if err != nil {
			return nil, err
		}
		for _, o := range photos {
			res = append(res, o)
		}
	default:
		return res, nil
	}

	return res, nil

}
