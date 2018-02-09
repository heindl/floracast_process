package materialtaxa

import (
	"bitbucket.org/heindl/taxa/datasources/fetchers"
	"golang.org/x/net/context"
	"bitbucket.org/heindl/taxa/datasources"
)

type photo struct {
	Citation string `json:""`
	Thumbnail string `json:""`
	Large string `json:""`
	Rank int `json:""`
}

func (Ω *NameUsage) photos(ctx context.Context) ([]*photo, error) {

	photos := []*photo{}
	for _, src := range Ω.Sources(datasources.DataSourceTypeGBIF, datasources.DataSourceTypeINaturalist) {
		fetchedPhotos, err := fetchers.FetchPhotos(ctx, src.SourceType(), src.TargetID())
		if err != nil {
			return nil, err
		}
		for _, p := range fetchedPhotos {
			rank := 0
			if p.Source() == datasources.DataSourceTypeINaturalist {
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