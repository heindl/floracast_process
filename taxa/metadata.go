package taxa

import (
	"bitbucket.org/heindl/processors/datasources/sourcefetchers"
	"golang.org/x/net/context"
	"bitbucket.org/heindl/processors/datasources"
	"bitbucket.org/heindl/processors/nameusage/nameusage"
)

type photo struct {
	Citation string `json:""`
	Thumbnail string `json:""`
	Large string `json:""`
	Rank int
}

func photos(ctx context.Context, usage nameusage.NameUsage) ([]*photo, error) {

	photos := []*photo{}

	srcs, err := usage.Sources(datasources.TypeGBIF, datasources.TypeINaturalist)
	if err != nil {
		return nil, err
	}

	for _, src := range srcs {

		srcType, err := src.SourceType()
		if err != nil {
			return nil, err
		}

		targetID, err := src.TargetID()
		if err != nil {
			return nil, err
		}

		fetchedPhotos, err := sourcefetchers.FetchPhotos(ctx, srcType, targetID)
		if err != nil {
			return nil, err
		}
		for _, p := range fetchedPhotos {
			rank := 0
			if p.Source() == datasources.TypeINaturalist {
				rank += 1
			}
			if p.Thumbnail() != "" {
				rank += 1
			}
			if p.Citation() != "" {
				rank += 1
			}
			photos = append(photos, &photo{
				Citation: p.Citation(),
				Large: p.Large(),
				Thumbnail: p.Thumbnail(),
				Rank: rank,
			})
		}
	}

	// TODO: Sort and clip to five here.

	return photos, nil
}

type description struct {
	Citation string `json:""`
	Text string `json:""`
}

func descriptions(ctx context.Context, usage nameusage.NameUsage) ([]*description, error) {

	res := []*description{}

	srcs, err := usage.Sources(datasources.TypeGBIF, datasources.TypeINaturalist)
	if err != nil {
		return nil, err
	}

	for _, src := range srcs {

		srcType, err := src.SourceType()
		if err != nil {
			return nil, err
		}

		targetID, err := src.TargetID()
		if err != nil {
			return nil, err
		}

		fetchedDescriptions, err := sourcefetchers.FetchDescriptions(ctx, srcType, targetID)
		if err != nil {
			return nil, err
		}
		for _, p := range fetchedDescriptions {
			rank := 0
			if p.Source() == datasources.TypeINaturalist {
				rank += 1
			}
			citation, err := p.Citation()
			if err != nil {
				return nil, err
			}
			if citation != "" {
				rank += 1
			}
			text, err := p.Text()
			if err != nil {
				return nil, err
			}
			res = append(res, &description{
				Citation: citation,
				Text: text,
			})
		}
	}

	// TODO: Sort and clip to five here.

	return res, nil
}