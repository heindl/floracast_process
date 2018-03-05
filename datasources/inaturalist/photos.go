package inaturalist

import (
	"bitbucket.org/heindl/process/datasources"
	"context"
	"github.com/dropbox/godropbox/errors"
)

type photo struct {
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

func (p *photo) Citation() string {
	return p.Attribution
}

func (p *photo) Thumbnail() string {
	return p.SmallURL
}

func (p *photo) Large() string {
	return p.LargeURL
}

func (p *photo) Source() datasources.SourceType {
	return datasources.TypeINaturalist
}

func FetchPhotos(ctx context.Context, targetID datasources.TargetID) ([]*photo, error) {
	taxa, err := newTaxaFetcher(ctx, false, false).FetchTaxa(taxonIDFromTargetID(targetID))
	if err != nil {
		return nil, err
	}
	if len(taxa) == 0 {
		return nil, errors.Newf("INaturalist taxon [%s] not found", targetID)
	}
	if len(taxa) > 1 {
		return nil, errors.Newf("Multiple INaturalist taxon found for TargetID [%s]", targetID)
	}

	res := []*photo{}
	for _, taxonPhoto := range taxa[0].TaxonPhotos {
		res = append(res, &taxonPhoto.Photo)
	}

	return res, nil
}
