package inaturalist

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/datasources/providers"
	"bitbucket.org/heindl/process/datasources/wikipedia"
	"context"
	"github.com/dropbox/godropbox/errors"
)

// descriptionProvider is an interface for taxon description data
type descriptionProvider taxon

// MLACitation formats Wikipedia citation if present, and a blank line if not.
func (p *descriptionProvider) Citation() (string, error) {
	if p.WikipediaURL == "" {
		return "", nil
	}
	return wikipedia.MLACitation(p.WikipediaURL)
}

// Text provides the Wikipedia summary.
func (p *descriptionProvider) Text() (string, error) {
	return p.WikipediaSummary, nil
}

// Source provides the SourceType.
func (p *descriptionProvider) SourceType() datasources.SourceType {
	return datasources.TypeGBIF
}

// FetchDescriptions provides a list of Wikipedia summaries.
func FetchDescriptions(ctx context.Context, targetID datasources.TargetID) ([]providers.Description, error) {

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

	if taxa[0].WikipediaSummary == "" {
		return nil, nil
	}

	if taxa[0].WikipediaURL == "" {
		return nil, errors.Newf("Inaturalist taxon [%s] has a WikipediaSummary but not a WikipediaURL", targetID)
	}

	t := descriptionProvider(*taxa[0])

	return []providers.Description{&t}, nil
}
