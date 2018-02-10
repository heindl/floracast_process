package gbif

import (
	"bitbucket.org/heindl/processors/datasources"
	"context"
	"bitbucket.org/heindl/processors/datasources/gbif/api"
	"github.com/dropbox/godropbox/errors"
)

type MaterializedPhoto struct {
	Link string `json:""`
	Attribution string `json:""`
}

func (p *MaterializedPhoto) Citation() string {
	return p.Attribution
}

func (p *MaterializedPhoto) Thumbnail() string {
	return ""
}

func (p *MaterializedPhoto) Large() string {
	return p.Link
}

func (p *MaterializedPhoto) Source() datasources.SourceType {
	return datasources.TypeGBIF
}

func FetchPhotos(ctx context.Context, targetID datasources.TargetID) ([]*MaterializedPhoto, error) {

	i, err := targetID.ToInt()
	if err != nil {
		return nil, err
	}

	list, err := api.Species(i).Media()
	if err != nil {
		return nil, errors.Wrapf(err, "Could not fetch GBIF Media [%s]", targetID)
	}

	photos := []*MaterializedPhoto{}

	// TODO: Check license here for appropriateness.
	// TODO: Look for created time? Doesn't appear to be in all records.
	for _, l := range list {
		if l.Type != "StillImage" {
			continue
		}
		if l.Identifier == "" {
			continue
		}

		p := MaterializedPhoto{
			Link: l.Identifier,
		}

		p.Attribution = l.Creator
		if l.Source != "" {
			p.Attribution = p.Attribution + ", " + l.Source
		}

		photos = append(photos, &p)
	}

	return photos, nil
}
