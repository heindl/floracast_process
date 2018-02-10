package inaturalist

import (
	"context"
	"bitbucket.org/heindl/processors/datasources"
	"github.com/dropbox/godropbox/errors"
)

type Photo struct {
	OriginalURL        string        `json:"original_url"`
	Flags              []interface{} `json:"flags"`
	Type               string        `json:"type"`
	URL                string        `json:"url"`
	SquareURL          string        `json:"square_url"`
	NativePageURL      string        `json:"native_page_url"`
	NativePhotoID      string        `json:"native_photo_id"`
	SmallURL           string        `json:"small_url"`
	Attribution        string        `json:"attribution"`
	MediumURL          string        `json:"medium_url"`
	ID                 int           `json:"id"`
	LicenseCode        string        `json:"license_code"`
	OriginalDimensions interface{}   `json:"original_dimensions"`
	LargeURL           string        `json:"large_url"`
}

func (p *Photo) Citation() string {
	return p.Attribution
}

func (p *Photo) Thumbnail() string {
	return p.SmallURL
}

func (p *Photo) Large() string {
	return p.LargeURL
}

func (p *Photo) Source() datasources.SourceType {
	return datasources.TypeINaturalist
}

func FetchPhotos(ctx context.Context, targetID datasources.TargetID) ([]*Photo, error) {
	taxa, err := NewTaxaFetcher(ctx, false, false).FetchTaxa(TaxonIDFromTargetID(targetID))
	if err != nil {
		return nil, err
	}
	if len(taxa) == 0 {
		return nil, errors.Newf("INaturalist Taxon [%s] not found", targetID)
	}
	if len(taxa) > 1 {
		return nil, errors.Newf("Multiple INaturalist Taxon found for TargetID [%s]", targetID)
	}

	res := []*Photo{}
	for _, taxonPhoto := range taxa[0].TaxonPhotos {
		res = append(res, &taxonPhoto.Photo)
	}

	return res, nil
}
