package gbif

import (
	"bitbucket.org/heindl/process/datasources"
	"context"
	"bitbucket.org/heindl/process/datasources/gbif/api"
	"github.com/dropbox/godropbox/errors"
	"bitbucket.org/heindl/process/utils"
)

type MaterializedDescription struct {
	attribution string `json:""`
	text string `json:""`
}

func (p *MaterializedDescription) Citation() (string, error) {
	return p.attribution, nil
}

func (p *MaterializedDescription) Text() (string, error) {
	return p.text, nil
}

func (p *MaterializedDescription) Source() datasources.SourceType {
	return datasources.TypeGBIF
}

func FetchDescriptions(ctx context.Context, targetID datasources.TargetID) ([]*MaterializedDescription, error) {

	i, err := targetID.ToInt()
	if err != nil {
		return nil, err
	}

	list, err := api.Species(i).Descriptions()
	if err != nil {
		return nil, errors.Wrapf(err, "Could not fetch GBIF Media [%s]", targetID)
	}

	descriptions := []*MaterializedDescription{}

	// TODO: Check license here for appropriateness.
	// TODO: Look for created time? Doesn't appear to be in all records.
	for _, l := range list {
		if l.Language != "en" {
			continue
		}

		if l.Type == "Distribution" {
			continue
		}

		// Not alot of information available on descriptions, so wait until we come across one.
		return nil, errors.Newf("Have Description with unknown type [%s]", utils.JsonOrSpew(l))
	}

	return descriptions, nil
}
