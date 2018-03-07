package sourcefetchers

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/datasources/gbif"
	"bitbucket.org/heindl/process/datasources/inaturalist"
	"bitbucket.org/heindl/process/datasources/providers"
	"context"
	"github.com/dropbox/godropbox/errors"
)

// FetchPhotos returns photos from the requested source.
func FetchPhotos(ctx context.Context, sourceTypeProvider datasources.SourceTypeProvider, targetIDProvider datasources.TargetIDProvider) ([]providers.Photo, error) {

	sourceType, err := sourceTypeProvider()
	if err != nil {
		return nil, err
	}

	targetID, err := targetIDProvider()
	if err != nil {
		return nil, err
	}
	switch sourceType {
	case datasources.TypeGBIF:
		return gbif.FetchPhotos(ctx, targetID)
	case datasources.TypeINaturalist:
		return inaturalist.FetchPhotos(ctx, targetID)
	default:
		return nil, errors.Newf("Unsupported SourceType [%s]", sourceType)
	}

}
