package name_usage

import (
	"time"
	"bitbucket.org/heindl/taxa/store"
	"github.com/dropbox/godropbox/errors"
	"bitbucket.org/heindl/taxa/utils"
)

type NameUsageSource struct {
	TaxonomicReference bool `json:""`
	SourceType store.DataSourceType `json:",omitempty"`
	TargetID store.DataSourceTargetID `json:",omitempty"`
	CanonicalName store.CanonicalName `json:",omitempty"`
	Synonyms []string `json:",omitempty"`
	CommonNames []string `json:",omitempty"`
	Descriptions []Description `json:",omitempty"`
	Photos []Photo `json:",omitempty"`
	OccurrenceCount int `json:",omitempty"`
	LastFetchedAt *time.Time `json:",omitempty"`
	ModifiedAt *time.Time `json:",omitempty"`
	CreatedAt *time.Time `json:",omitempty"`
}

type Description struct {
	Citation string `json:",omitempty"`
	Text string `json:",omitempty"`
}

type Photo struct {
	Citation string `json:",omitempty"`
	Source string `json:",omitempty"`
}

func NewNameUsageSource(sourceType store.DataSourceType, targetID store.DataSourceTargetID, canonicalName store.CanonicalName, isTaxonomic bool) (*NameUsageSource, error) {

	if !sourceType.Valid() {
		return nil, errors.Newf("Invalid source type [%s]", sourceType)
	}

	if !targetID.Valid() {
		return nil, errors.Newf("Invalid target id [%s]", targetID)
	}

	if err := canonicalName.Validate(); err != nil {
		return nil, err
	}

	// Verbose expectations.
	if (!isTaxonomic) {
		if sourceType == store.DataSourceTypeGBIF || sourceType == store.DataSourceTypeINaturalist || sourceType == store.DataSourceTypeNatureServe {
			return nil, errors.Newf("Unexpected: source type should create a taxonomic name usage [%s]", sourceType)
		}
	} else {
		return nil, errors.Newf("Unexpected: source type should not be taxonomic [%s]", sourceType)
	}

	return &NameUsageSource{
		TaxonomicReference: isTaxonomic,
		SourceType: sourceType,
		TargetID: targetID,
		CanonicalName: canonicalName,
		CreatedAt: utils.TimePtr(time.Now()),
		ModifiedAt: utils.TimePtr(time.Now()),
		}, nil
}