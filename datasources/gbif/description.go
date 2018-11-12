package gbif

import (
	"github.com/heindl/floracast_process/datasources"
	"github.com/heindl/floracast_process/datasources/providers"
	"github.com/heindl/floracast_process/utils"
	"context"
	"fmt"
	"github.com/dropbox/godropbox/errors"
)

type description struct {
	Description    string `json:"description"`
	Key            int    `json:"key"`
	Language       string `json:"language"`
	Src            string `json:"source"`
	SourceTaxonKey int    `json:"sourceTaxonKey"`
	Type           string `json:"type"`
	License        string `json:"license"`
}

func (Ω species) fetchDescriptions() (descriptions []description, err error) {
	var offset int
	for {
		var response struct {
			page
			Results []description `json:"results"`
		}
		url := fmt.Sprintf("%s/descriptions?offset=%d&limit=50", Ω.url(), offset)
		if err := utils.RequestJSON(url, &response); err != nil {
			return nil, err
		}
		descriptions = append(descriptions, response.Results...)
		offset += response.Limit
		if response.EndOfRecords {
			break
		}
	}
	return
}

func (p *description) Citation() (string, error) {
	return p.Src, nil
}

func (p *description) Text() (string, error) {
	return p.Description, nil
}

func (p *description) SourceType() datasources.SourceType {
	return datasources.TypeGBIF
}

// FetchDescriptions returns a list of DescriptionProviders.
func FetchDescriptions(ctx context.Context, targetID datasources.TargetID) ([]providers.Description, error) {

	i, err := targetID.ToInt()
	if err != nil {
		return nil, err
	}

	list, err := species(i).fetchDescriptions()
	if err != nil {
		return nil, errors.Wrapf(err, "Could not fetch GBIF media [%s]", targetID)
	}

	descriptions := []providers.Description{}

	// TODO: Check license here for appropriateness.
	// TODO: Look for created time? Doesn't appear to be in all records.
	for _, l := range list {
		if l.Language != "en" {
			continue
		}

		if l.Type == "distribution" {
			continue
		}

		// Not alot of information available on descriptions, so wait until we come across one.
		return nil, errors.Newf("Have description with unknown type [%s]", utils.JsonOrSpew(l))
	}

	return descriptions, nil
}
