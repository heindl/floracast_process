package nameusage

import (
	"time"
	"github.com/dropbox/godropbox/errors"
	"bitbucket.org/heindl/taxa/utils"
	"strings"
	"sync"
	"encoding/json"
	"bitbucket.org/heindl/taxa/datasources"
)

type NameUsageSource struct {
	mutex *sync.Mutex
	taxonomicReference bool
	sourceType datasources.DataSourceType
	targetID datasources.DataSourceTargetID
	canonicalName *CanonicalName
	synonyms CanonicalNames
	commonNames []string
	descriptions []description
	photos []photo
	occurrenceCount int
	lastFetchedAt *time.Time
	modifiedAt *time.Time
	createdAt *time.Time
}

func (Ω *NameUsageSource) MarshalJSON() ([]byte, error) {

	if Ω == nil {
		return nil, nil
	}

	m := map[string]interface{}{
		//"SourceType": Ω.sourceType,
		//"TargetID": Ω.targetID,
		"CanonicalName": Ω.canonicalName,
		"TotalOccurrenceCount": Ω.occurrenceCount,
		"LastFetchedAt": Ω.lastFetchedAt,
		"ModifiedAt": Ω.modifiedAt,
		"CreatedAt": Ω.createdAt,
	}

	if Ω.taxonomicReference {
		m["TaxonomicReference"] = Ω.taxonomicReference
	}

	if len(Ω.commonNames) > 0 {
		m["CommonNames"] = Ω.commonNames
	}

	if len(Ω.synonyms) > 0 {
		m["Synonyms"] = Ω.synonyms
	}

	return json.Marshal(m)
}

type description struct {
	citation string
	text string
}

type photo struct {
	citation string
	source string
}

func NewNameUsageSource(sourceType datasources.DataSourceType, targetID datasources.DataSourceTargetID, canonicalName *CanonicalName) (*NameUsageSource, error) {

	if !sourceType.Valid() {
		return nil, errors.Newf("Invalid SourceType [%s]", sourceType)
	}

	if !targetID.Valid(sourceType) {
		return nil, errors.Newf("Invalid TargetID [%s, %s]", targetID, canonicalName.ScientificName())
	}

	isTaxonomic := sourceType == datasources.DataSourceTypeGBIF || sourceType == datasources.DataSourceTypeINaturalist || sourceType == datasources.DataSourceTypeNatureServe

	// Verbose expectations.
	//if !isTaxonomic && expectTaxonomicSourceType {
	//	return nil, errors.Newf("Unexpected: source type should create a taxonomic name usage [%s]", sourceType)
	//}
	//
	//if isTaxonomic && !expectTaxonomicSourceType{
	//	return nil, errors.Newf("Unexpected: source type should not be taxonomic [%s]", sourceType)
	//}

	return &NameUsageSource{
		mutex: &sync.Mutex{},
		taxonomicReference: isTaxonomic,
		sourceType: sourceType,
		targetID: targetID,
		canonicalName: canonicalName,
		createdAt: utils.TimePtr(time.Now()),
		modifiedAt: utils.TimePtr(time.Now()),
		}, nil
}

func (Ω *NameUsageSource) SourceType() datasources.DataSourceType {
	return Ω.sourceType
}

func (Ω *NameUsageSource) TargetID() datasources.DataSourceTargetID {
	return Ω.targetID
}

func (Ω *NameUsageSource) CanonicalName() *CanonicalName{
	return Ω.canonicalName
}

func (Ω *NameUsageSource) LastFetchedAt() *time.Time {
	return Ω.lastFetchedAt
}

func (Ω *NameUsageSource) ScientificName() string {
	return Ω.canonicalName.name
}

func (Ω *NameUsageSource) CommonNames() []string {
	return Ω.commonNames
}

func (Ω *NameUsageSource) AddSynonym(synonym *CanonicalName) error {
	if Ω.canonicalName.Equals(synonym) {
		return nil
	}
	Ω.synonyms = Ω.synonyms.AddToSet(synonym)
	return nil
}

func (Ω *NameUsageSource) AddCommonNames(names ...string) error {
	Ω.mutex.Lock()
	defer Ω.mutex.Unlock()
	for _, s := range names {
		s = strings.ToLower(strings.TrimSpace(s))
		if s == "" {
			return errors.New("Invalid CommonName: Received empty string")
		}
		Ω.commonNames = utils.AddStringToSet(Ω.commonNames, s)
	}
	return nil
}

func (Ω *NameUsageSource) SetOccurrenceCount(count int) error {
	Ω.occurrenceCount = count
	return nil
}