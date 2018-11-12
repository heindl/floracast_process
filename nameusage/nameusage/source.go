package nameusage

import (
	"github.com/heindl/floracast_process/datasources"
	"github.com/heindl/floracast_process/nameusage/canonicalname"
	"github.com/heindl/floracast_process/utils"
	"encoding/json"
	"github.com/dropbox/godropbox/errors"
	"strings"
	"sync"
	"time"
)

// Source represents a NameUsage source that provides taxonomical or occurrence information.
type Source interface {
	RegisterOccurrenceFetch(count int) error
	AddCommonNames(names ...string) error
	SourceType() (datasources.SourceType, error)
	TargetID() (datasources.TargetID, error)
	CanonicalName() *canonicalname.Name
	Synonyms() canonicalname.Names
	LastFetchedAt() *time.Time
	CommonNames() []string
	OccurrenceCount() int
	AddSynonym(synonym *canonicalname.Name) error
	Bytes() ([]byte, error)
}

type source struct {
	mutex              *sync.Mutex
	TaxonomicReference bool                   `json:",omitempty" firestore:",omitempty"`
	SrcType            datasources.SourceType `json:"-" firestore:"-"`
	TrgtID             datasources.TargetID   `json:"-" firestore:"-"`
	CnnclNm            *canonicalname.Name    `json:"Name,omitempty" firestore:"Name,omitempty"`
	Snnms              canonicalname.Names    `json:"Synonyms,omitempty" firestore:"Synonyms,omitempty"`
	CmmnNms            []string               `json:"CommonNames,omitempty" firestore:"CommonNames,omitempty"`
	Occurrences        int                    `json:"Occurrences,omitempty" firestore:"Occurrences,omitempty"`
	LastFtchdAt        *time.Time             `json:"LastFetchedAt,omitempty" firestore:"LastFetchedAt,omitempty"`
}

// Sources is a list of sources.
type Sources []Source

//func (Ω Sources) HasName(æ *canonicalname.Name) bool {
//	for _, s := range Ω {
//		if s.Name().Equals(æ) || s.Synonyms().Contains(æ) {
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
//		"Name": Ω.canonicalName,
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

// NewSource creates a new NameUsageSource.
func NewSource(sourceType datasources.SourceType, targetID datasources.TargetID, canonicalName *canonicalname.Name) (Source, error) {

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
		SrcType:            sourceType,
		TrgtID:             targetID,
		CnnclNm:            canonicalName,
	}, nil
}

// SourceType validates and returns the SourceType
func (Ω *source) SourceType() (datasources.SourceType, error) {
	if Ω == nil || !Ω.SrcType.Valid() {
		return datasources.SourceType(""), errors.Newf("Invalid SourceType [%s]", Ω.SrcType)
	}
	return Ω.SrcType, nil
}

// TargetID validates and returns the TargetID
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

func (Ω *source) CanonicalName() *canonicalname.Name {
	return Ω.CnnclNm
}

func (Ω *source) Synonyms() canonicalname.Names {
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

func (Ω *source) AddSynonym(synonym *canonicalname.Name) error {
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
