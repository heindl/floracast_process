package nameusage

import (
	"time"
	"github.com/dropbox/godropbox/errors"
	"bitbucket.org/heindl/process/utils"
	"strings"
	"sync"
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/nameusage/canonicalname"
	"github.com/mongodb/mongo-tools/common/json"
)


type Source interface{
	RegisterOccurrenceFetch(count int) error
	AddCommonNames(names ...string) error
	SourceType() (datasources.SourceType, error)
	TargetID() (datasources.TargetID, error)
	CanonicalName() *canonicalname.CanonicalName
	Synonyms() canonicalname.CanonicalNames
	LastFetchedAt() *time.Time
	CommonNames() []string
	OccurrenceCount() int
	AddSynonym(synonym *canonicalname.CanonicalName) error
	Bytes() ([]byte, error)
}

type source struct {
	mutex              *sync.Mutex
	TaxonomicReference bool                         `json:",omitempty" firestore:",omitempty"`
	SrcType         datasources.SourceType       `json:"-" firestore:"-"`
	TrgtID           datasources.TargetID         `json:"-" firestore:"-"`
	CnnclNm      *canonicalname.CanonicalName `json:"CanonicalName,omitempty" firestore:"CanonicalName,omitempty"`
	Snnms           canonicalname.CanonicalNames `json:"Synonyms,omitempty" firestore:"Synonyms,omitempty"`
	CmmnNms        []string                     `json:"CommonNames,omitempty" firestore:"CommonNames,omitempty"`
	Occurrences    int                          `json:"Occurrences,omitempty" firestore:"Occurrences,omitempty"`
	LastFtchdAt      *time.Time                   `json:"LastFetchedAt,omitempty" firestore:"LastFetchedAt,omitempty"`
}

type Sources []Source

//func (Ω Sources) HasName(æ *canonicalname.CanonicalName) bool {
//	for _, s := range Ω {
//		if s.CanonicalName().Equals(æ) || s.Synonyms().Contains(æ) {
//			return true
//		}
//	}
//	return false
//}
//
//func (Ω *source) MarshalJSON() ([]byte, error) {
//
//	if Ω == nil {
//		return nil, nil
//	}
//
//	m := map[string]interface{}{
//		//"SourceType": Ω.SourceType,
//		//"TargetID": Ω.TargetID,
//		"CanonicalName": Ω.canonicalName,
//		"TotalOccurrenceCount": Ω.occurrenceCount,
//		"LastFetchedAt": Ω.lastFetchedAt,
//		"ModifiedAt": Ω.modifiedAt,
//		"CreatedAt": Ω.createdAt,
//	}
//
//	if Ω.taxonomicReference {
//		m["TaxonomicReference"] = Ω.taxonomicReference
//	}
//
//	if len(Ω.commonNames) > 0 {
//		m["CommonNames"] = Ω.commonNames
//	}
//
//	if len(Ω.synonyms) > 0 {
//		m["Synonyms"] = Ω.synonyms
//	}
//
//	return json.Marshal(m)
//}

func NewSource(sourceType datasources.SourceType, targetID datasources.TargetID, canonicalName *canonicalname.CanonicalName) (Source, error) {

	if !sourceType.Valid() {
		return nil, errors.Newf("Invalid SourceType [%s]", sourceType)
	}

	if !targetID.Valid(sourceType) {
		return nil, errors.Newf("Invalid TargetID [%s, %s]", targetID, canonicalName.ScientificName())
	}

	isTaxonomic := sourceType == datasources.TypeGBIF || sourceType == datasources.TypeINaturalist || sourceType == datasources.TypeNatureServe

	// Verbose expectations.
	//if !isTaxonomic && expectTaxonomicSourceType {
	//	return nil, errors.Newf("Unexpected: source type should create a taxonomic name usage [%s]", SourceType)
	//}
	//
	//if isTaxonomic && !expectTaxonomicSourceType{
	//	return nil, errors.Newf("Unexpected: source type should not be taxonomic [%s]", SourceType)
	//}

	return &source{
		mutex:              &sync.Mutex{},
		TaxonomicReference: isTaxonomic,
		SrcType:         sourceType,
		TrgtID:           targetID,
		CnnclNm:      canonicalName,
		}, nil
}

func (Ω *source) SourceType() (datasources.SourceType, error) {
	if Ω == nil || !Ω.SrcType.Valid() {
		return datasources.SourceType(""), errors.Newf("Invalid SourceType [%s]", Ω.SrcType)
	}
	return Ω.SrcType, nil
}

func (Ω *source) TargetID() (datasources.TargetID, error) {
	srcType, err := Ω.SourceType()
	if err != nil {
		return datasources.TargetID(""), err
	}
	if Ω == nil || !Ω.TrgtID.Valid(srcType) {
		return datasources.TargetID(""), errors.Newf("Invalid TargetID [%s] with SourceType [%s]", Ω.TrgtID, srcType)
	}
	return Ω.TrgtID, nil
}

func (Ω *source) CanonicalName() *canonicalname.CanonicalName{
	return Ω.CnnclNm
}

func (Ω *source) Synonyms() canonicalname.CanonicalNames {
	return Ω.Snnms
}

func (Ω *source) LastFetchedAt() *time.Time {
	return Ω.LastFtchdAt
}

func (Ω *source) CommonNames() []string {
	return Ω.CmmnNms
}

func (Ω *source) OccurrenceCount() int {
	return Ω.Occurrences
}

func (Ω *source) AddSynonym(synonym *canonicalname.CanonicalName) error {
	Ω.Snnms = Ω.Snnms.AddToSet(synonym)
	return nil
}

func (Ω *source) AddCommonNames(names ...string) error {
	Ω.mutex.Lock()
	defer Ω.mutex.Unlock()
	for _, s := range names {
		s = strings.ToLower(strings.TrimSpace(s))
		if s == "" {
			return errors.New("Invalid CommonName: Received empty string")
		}
		Ω.CmmnNms = utils.AddStringToSet(Ω.CmmnNms, s)
	}
	return nil
}

func (Ω *source) RegisterOccurrenceFetch(count int) error {
	Ω.Occurrences += count
	Ω.LastFtchdAt = utils.TimePtr(time.Now())
	return nil
}

func (Ω *source) Bytes() ([]byte, error) {
	return json.Marshal(*Ω)
}