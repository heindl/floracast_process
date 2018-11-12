package sourcefetchers

import (
	"github.com/heindl/floracast_process/datasources"
	"github.com/heindl/floracast_process/datasources/gbif"
	"github.com/heindl/floracast_process/datasources/inaturalist"
	"github.com/heindl/floracast_process/datasources/providers"
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
