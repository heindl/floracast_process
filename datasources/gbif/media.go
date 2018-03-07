package gbif

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/datasources/providers"
	"bitbucket.org/heindl/process/utils"
	"context"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"strings"
)

type media struct {
	TaxonKey       int    `json:"taxonKey"`
	Format         string `json:"format"`
	Identifier     string `json:"identifier"`
	License        string `json:"license"`
	Publisher      string `json:"publisher"`
	References     string `json:"references"`
	Src            string `json:"source"`
	SourceTaxonKey int    `json:"sourceTaxonKey"`
	Title          string `json:"title"`
	Type           string `json:"type"`
	RightsHolder   string `json:"rightsHolder"` // From occurrence search
	Created        string `json:"created"`      // From occurrence search
	Creator        string `json:"creator"`      // From occurrence search
}

func (Ω species) fetchMedia() ([]*media, error) {
	var offset int
	res := []*media{}
	for {
		var response struct {
			page
			Results []*media `json:"results"`
		}
		url := fmt.Sprintf("%s/media?offset=%d&limit=300", Ω.url(), offset)
		if err := utils.RequestJSON(url, &response); err != nil {
			return nil, err
		}
		res = append(res, response.Results...)
		offset += response.Limit
		if response.EndOfRecords {
			break
		}
	}
	return res, nil
}

func (Ω *media) Citation() string {
	attr := Ω.Creator + ", " + Ω.Src
	return strings.TrimRight(strings.TrimSpace(attr), ",")
}

func (Ω *media) Thumbnail() string {
	return ""
}

func (Ω *media) Large() string {
	return Ω.Identifier
}

func (Ω *media) SourceType() datasources.SourceType {
	return datasources.TypeGBIF
}

// FetchPhotos returns an interface list of PhotoProviders.
func FetchPhotos(_ context.Context, targetID datasources.TargetID) ([]providers.Photo, error) {

	i, err := targetID.ToInt()
	if err != nil {
		return nil, err
	}

	list, err := species(i).fetchMedia()
	if err != nil {
		return nil, errors.Wrapf(err, "Could not fetch GBIF media [%s]", targetID)
	}

	photos := []providers.Photo{}

	// TODO: Check license here for appropriateness.
	// TODO: Look for created time? Doesn't appear to be in all records.
	for _, _m := range list {
		m := _m
		if m.Type != "StillImage" {
			continue
		}
		if m.Identifier == "" {
			continue
		}

		photos = append(photos, m)
	}

	return photos, nil
}
