package nameusage

import (
	"time"
	"bitbucket.org/heindl/taxa/store"
	"github.com/dropbox/godropbox/errors"
	"bitbucket.org/heindl/taxa/utils"
	"strings"
	"sync"
	"encoding/json"
)

type NameUsageSource struct {
	mutex *sync.Mutex
	taxonomicReference bool
	sourceType store.DataSourceType
	targetID store.DataSourceTargetID
	canonicalName *CanonicalName
	//synonyms CanonicalNames
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
		"OccurrenceCount": Ω.occurrenceCount,
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

func NewNameUsageSource(sourceType store.DataSourceType, targetID store.DataSourceTargetID, canonicalName *CanonicalName, isTaxonomic bool) (*NameUsageSource, error) {

	if !sourceType.Valid() {
		return nil, errors.Newf("Invalid source type [%s]", sourceType)
	}

	if !targetID.Valid() {
		return nil, errors.Newf("Invalid target id [%s]", targetID)
	}

	expectTaxonomicSourceType := sourceType == store.DataSourceTypeGBIF || sourceType == store.DataSourceTypeINaturalist || sourceType == store.DataSourceTypeNatureServe

	// Verbose expectations.
	if !isTaxonomic && expectTaxonomicSourceType {
		return nil, errors.Newf("Unexpected: source type should create a taxonomic name usage [%s]", sourceType)
	}

	if isTaxonomic && !expectTaxonomicSourceType{
		return nil, errors.Newf("Unexpected: source type should not be taxonomic [%s]", sourceType)
	}

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

func (Ω *NameUsageSource) ScientificName() string {
	return Ω.canonicalName.name
}

func (Ω *NameUsageSource) CommonNames() []string {
	return Ω.commonNames
}

func (Ω *NameUsageSource) SetCommonNames(names ...string) error {
	res := []string{}
	for _, s := range names {
		s = strings.ToLower(strings.TrimSpace(s))
		if s == "" {
			return errors.New("Invalid Synonym: Received empty string")
		}
		res = append(res, s)
	}
	Ω.mutex.Lock()
	defer Ω.mutex.Unlock()
	Ω.commonNames = utils.AddStringToSet(Ω.commonNames, res...)
	return nil
}

func (Ω *NameUsageSource) SetOccurrenceCount(count int) error {
	Ω.occurrenceCount = count
	return nil
}