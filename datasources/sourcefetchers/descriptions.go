package sourcefetchers

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/datasources/gbif"
	"bitbucket.org/heindl/process/datasources/inaturalist"
	"context"
	"github.com/dropbox/godropbox/errors"
)

type Description interface {
	Citation() (string, error)
	Text() (string, error)
	Source() datasources.SourceType
}

func FetchDescriptions(ctx context.Context, sourceType datasources.SourceType, targetID datasources.TargetID) ([]Description, error) {

	res := []Description{}
	switch sourceType {
	case datasources.TypeGBIF:
		descriptions, err := gbif.FetchDescriptions(ctx, targetID)
		if err != nil {
			return nil, err
		}
		for _, d := range descriptions {
			res = append(res, d)
		}
	case datasources.TypeINaturalist:
		descriptions, err := inaturalist.FetchDescriptions(ctx, targetID)
		if err != nil {
			return nil, err
		}
		for _, d := range descriptions {
			res = append(res, d)
		}
	default:
		return nil, errors.Newf("Unsupported SourceType [%s]", sourceType)
	}

	return res, nil

}
