package nameusagesource

import (
	"time"
	"github.com/dropbox/godropbox/errors"
	"bitbucket.org/heindl/processors/utils"
	"strings"
	"sync"
	"encoding/json"
	"bitbucket.org/heindl/processors/datasources"
	"bitbucket.org/heindl/processors/nameusage/canonicalname"
)

type Source struct {
	mutex *sync.Mutex
	taxonomicReference bool
	sourceType datasources.SourceType
	targetID datasources.TargetID
	canonicalName *canonicalname.CanonicalName
	synonyms canonicalname.CanonicalNames
	commonNames []string
	occurrenceCount int
	lastFetchedAt *time.Time
	modifiedAt *time.Time
	createdAt *time.Time
}

type Sources []*Source

func (Ω Sources) HasName(æ *canonicalname.CanonicalName) bool {
	for _, s := range Ω {
		if s.CanonicalName().Equals(æ) || s.Synonyms().Contains(æ) {
			return true
		}
	}
	return false
}

func (Ω *Source) MarshalJSON() ([]byte, error) {

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

func NewSource(sourceType datasources.SourceType, targetID datasources.TargetID, canonicalName *canonicalname.CanonicalName) (*Source, error) {

	if !sourceType.Valid() {
		return nil, errors.Newf("Invalid SourceType [%s]", sourceType)
	}

	if !targetID.Valid(sourceType) {
		return nil, errors.Newf("Invalid TargetID [%s, %s]", targetID, canonicalName.ScientificName())
	}

	isTaxonomic := sourceType == datasources.TypeGBIF || sourceType == datasources.TypeINaturalist || sourceType == datasources.TypeNatureServe

	// Verbose expectations.
	//if !isTaxonomic && expectTaxonomicSourceType {
	//	return nil, errors.Newf("Unexpected: source type should create a taxonomic name usage [%s]", sourceType)
	//}
	//
	//if isTaxonomic && !expectTaxonomicSourceType{
	//	return nil, errors.Newf("Unexpected: source type should not be taxonomic [%s]", sourceType)
	//}

	return &Source{
		mutex: &sync.Mutex{},
		taxonomicReference: isTaxonomic,
		sourceType: sourceType,
		targetID: targetID,
		canonicalName: canonicalName,
		createdAt: utils.TimePtr(time.Now()),
		modifiedAt: utils.TimePtr(time.Now()),
		}, nil
}

func (Ω *Source) SourceType() datasources.SourceType {
	return Ω.sourceType
}

func (Ω *Source) TargetID() datasources.TargetID {
	return Ω.targetID
}

func (Ω *Source) CanonicalName() *canonicalname.CanonicalName{
	return Ω.canonicalName
}

func (Ω *Source) Synonyms() canonicalname.CanonicalNames {
	return Ω.synonyms
}

func (Ω *Source) LastFetchedAt() *time.Time {
	return Ω.lastFetchedAt
}

func (Ω *Source) CommonNames() []string {
	return Ω.commonNames
}

func (Ω *Source) OccurrenceCount() int {
	return Ω.occurrenceCount
}

func (Ω *Source) AddSynonym(synonym *canonicalname.CanonicalName) error {
	if Ω.canonicalName.Equals(synonym) {
		return nil
	}
	Ω.synonyms = Ω.synonyms.AddToSet(synonym)
	return nil
}

func (Ω *Source) AddCommonNames(names ...string) error {
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

func (Ω *Source) RegisterOccurrenceFetch(count int) error {
	Ω.occurrenceCount += count
	Ω.lastFetchedAt = utils.TimePtr(time.Now())
	return nil
}