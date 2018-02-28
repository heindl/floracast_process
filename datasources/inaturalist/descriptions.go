package inaturalist

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/datasources/wikipedia"
	"context"
	"github.com/dropbox/godropbox/errors"
)

type DescriptionProvider Taxon

func (p *DescriptionProvider) Citation() (string, error) {
	if p.WikipediaURL == "" {
		return "", nil
	}
	return wikipedia.Citation(p.WikipediaURL)
}

func (p *DescriptionProvider) Text() (string, error) {
	return p.WikipediaSummary, nil
}

func (p *DescriptionProvider) Source() datasources.SourceType {
	return datasources.TypeGBIF
}

func FetchDescriptions(ctx context.Context, targetID datasources.TargetID) ([]*DescriptionProvider, error) {

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

	if taxa[0].WikipediaSummary == "" {
		return nil, nil
	}

	if taxa[0].WikipediaURL == "" {
		return nil, errors.Newf("Inaturalist Taxon [%s] has a WikipediaSummary but not a WikipediaURL", targetID)
	}

	t := DescriptionProvider(*taxa[0])

	return []*DescriptionProvider{&t}, nil
}
