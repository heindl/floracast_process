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

func photos(ctx context.Context, usage *nameusage.NameUsage) ([]*photo, error) {

	photos := []*photo{}
	for _, src := range usage.Sources(datasources.TypeGBIF, datasources.TypeINaturalist) {
		fetchedPhotos, err := sourcefetchers.FetchPhotos(ctx, src.SourceType(), src.TargetID())
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

func descriptions(ctx context.Context, usage *nameusage.NameUsage) ([]*description, error) {

	res := []*description{}
	for _, src := range usage.Sources(datasources.TypeGBIF, datasources.TypeINaturalist) {
		fetchedDescriptions, err := sourcefetchers.FetchDescriptions(ctx, src.SourceType(), src.TargetID())
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