package nameusage

import (
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"sort"
	"bitbucket.org/heindl/processors/datasources"
	"bitbucket.org/heindl/processors/nameusage/canonicalname"
	"bitbucket.org/heindl/processors/utils"
	"context"
	"bitbucket.org/heindl/processors/store"
	"encoding/json"
)

type NameUsage interface {
	ID() (NameUsageID, error)
	Sources(sourceTypes ...datasources.SourceType) (Sources, error)
	Occurrences(sourceTypes ...datasources.SourceType) (int, error)
	HasSource(datasources.SourceType, datasources.TargetID) (bool, error)
	AddSources(sources ...Source) error
	Synonyms() (canonicalname.CanonicalNames, error)
	AllScientificNames() ([]string, error)
	CanonicalName() *canonicalname.CanonicalName
	CommonName() (string, error)
	ShouldCombine(b NameUsage) (bool, error)
	Combine(b NameUsage) (NameUsage, error)
	HasScientificName(s string) (bool, error)
	ScientificNameReferenceLedger() (NameReferenceLedger, error)
	CommonNameReferenceLedger() (NameReferenceLedger, error)
	Upload(context.Context, store.FloraStore) (deletedUsageIDs NameUsageIDs, err error)
}

const storeKeyScientificName = "ScientificNames"

type usage struct {
	Id                    NameUsageID`json:"-" firestore:"-"`
	Cn         *canonicalname.CanonicalName `json:"CanonicalName" firestore:"CanonicalName"`
	Occrrncs int `json:"Occurrences,omitempty" firestore:"Occurrences,omitempty"`
	SciNames         map[string]bool `json:"ScientificNames,omitempty" firestore:"ScientificNames,omitempty"`
	Srcs             map[datasources.SourceType]map[datasources.TargetID]*source `json:"Sources,omitempty" firestore:"Sources,omitempty"`
}

func NewNameUsage(src Source) (NameUsage, error) {

	id, err := newNameUsageID()
	if err != nil {
		return nil, err
	}

	u := usage{
		Id:                    id,
		Cn:         src.CanonicalName(),
	}

	if err := u.AddSources(src); err != nil {
		return nil, err
	}

	return &u, nil
}

func NameUsageFromJSON(id NameUsageID, b []byte) (NameUsage, error) {
	if !id.Valid() {
		return nil, errors.Newf("Invalid NameUsageID [%s]", id)
	}
	u := usage{}
	if err := json.Unmarshal(b, &u); err != nil {
		return nil, err
	}
	u.Id = id

	return &u, nil
}

func (Ω *usage) ID() (NameUsageID, error) {
	if !Ω.Id.Valid() {
		return NameUsageID(""), errors.Newf("Invalid NameUsageID [%s]", Ω.Id)
	}
	return Ω.Id, nil
}

func (Ω *usage) Sources(sourceTypes ...datasources.SourceType) (Sources, error) {
	res := Sources{}
	for srcType, targets := range Ω.Srcs {

		if !srcType.Valid() {
			return nil, errors.Newf("Invalid SourceType [%s]", srcType)
		}

		if len(sourceTypes) > 0 && !datasources.HasDataSourceType(sourceTypes, srcType) {
			continue
		}

		for targetID, src := range targets {
			if !targetID.Valid(srcType) {
				return nil, errors.Newf("Invalid TargetID [%s] with SourceType [%s]", targetID, srcType)
			}
			src.TrgtID = targetID
			src.SrcType = srcType
			res = append(res, src)
		}

	}
	return res, nil
}

func (Ω *usage) Occurrences(sourceTypes ...datasources.SourceType) (int, error) {

	srcs, err := Ω.Sources(sourceTypes...)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, src := range srcs {
		count = count + src.OccurrenceCount()
	}
	return count, nil
}

func (Ω *usage) HasScientificName(s string) (bool, error) {
	names, err := Ω.AllScientificNames()
	if err != nil {
		return false, err
	}
	return utils.ContainsString(names, s), nil
}

func (Ω *usage) HasSource(sourceType datasources.SourceType, targetID datasources.TargetID) (bool, error) {

	if !sourceType.Valid() {
		return false, errors.Newf("Invalid SourceType [%s]", sourceType)
	}

	if !targetID.Valid(sourceType) {
		return false, errors.Newf("Invalid TargetID [%s] from SourceType [%s]", targetID, sourceType)
	}

	if _, ok := Ω.Srcs[sourceType]; !ok {
		return false, nil
	}

	if _, ok := Ω.Srcs[sourceType][targetID]; !ok {
		return false, nil
	}

	return true, nil
}

func (Ω *usage) AddSources(sources ...Source) error {

	if Ω.Srcs == nil {
		Ω.Srcs = map[datasources.SourceType]map[datasources.TargetID]*source{}
	}

	for _, src := range sources {

		srcType, err := src.SourceType()
		if err != nil {
			return err
		}

		targetID, err := src.TargetID()
		if err != nil {
			return err
		}

		if _, ok := Ω.Srcs[srcType]; !ok {
			Ω.Srcs[srcType] = map[datasources.TargetID]*source{}
		}

		b, err := src.Bytes()
		if err != nil {
			return err
		}
		s := source{}
		if err := json.Unmarshal(b, &s); err != nil {
			return errors.Wrap(err, "Could not Unmarshal Source")
		}

		Ω.Srcs[srcType][targetID] = &s

	}
	return nil
}

func (Ω *usage) CanonicalName() *canonicalname.CanonicalName {
	return Ω.Cn
}


func (Ω *usage) Synonyms() (canonicalname.CanonicalNames, error) {
	res := canonicalname.CanonicalNames{}
	srcs, err := Ω.Sources()
	if err != nil {
		return nil, err
	}
	for _, src := range srcs {
		for _, s := range src.Synonyms().AddToSet(src.CanonicalName()) {
			if s.Equals(Ω.Cn) {
				continue
			}
			res = res.AddToSet(s)
		}
	}
	return res, nil
}

func (Ω *usage) AllScientificNames() ([]string, error){
	synonyms, err := Ω.Synonyms()
	if err != nil {
		return nil, err
	}
	return utils.AddStringToSet(synonyms.ScientificNames(), Ω.CanonicalName().ScientificName()), nil
}

func (Ω *usage) ScientificNameReferenceLedger() (NameReferenceLedger, error) {
	srcs, err := Ω.Sources()
	if err != nil {
		return nil, err
	}
	ledger := NameReferenceLedger{}
	for _, src := range srcs {
		ledger = ledger.IncrementName(src.CanonicalName().ScientificName(), src.OccurrenceCount())
		for _, synonym := range src.Synonyms() {
			ledger = ledger.IncrementName(synonym.ScientificName(), 0)
		}
	}
	sort.Sort(ledger)
	return ledger, nil
}

func (Ω *usage) CommonNameReferenceLedger() (NameReferenceLedger, error) {
	ledger := NameReferenceLedger{}
	srcs, err := Ω.Sources()
	if err != nil {
		return nil, err
	}
	for _, src := range srcs {
		for _, cn := range src.CommonNames() {
			ledger = ledger.IncrementName(cn, src.OccurrenceCount())
		}
	}
	return ledger, nil
}

// Get the most popular common name from all sources
// Or perhaps the most common among those those with Occurrence Counts?
func (Ω *usage) CommonName() (string, error) {

	ledger, err := Ω.CommonNameReferenceLedger()
	if err != nil {
		return "", err
	}

	if len(ledger) == 0 || ledger[0].Name == "" {
		return "", errors.New("CommonName not found")
	}

	return ledger[0].Name, nil
}

func (a *usage) ShouldCombine(b NameUsage) (bool, error) {

	aNames, err := a.AllScientificNames()
	if err != nil {
		return false, err
	}

	bNames, err := b.AllScientificNames()
	if err != nil {
		return false, err
	}

	if utils.IntersectsStrings(aNames, bNames) {
		return true, nil
	}

	srcs, err := a.Sources()
	if err != nil {
		return false, err
	}

	for _, src := range srcs {

		srcType, err := src.SourceType()
		if err != nil {
			return false, err
		}

		targetID, err := src.TargetID()
		if err != nil {
			return false, err
		}

		hasSource, err := b.HasSource(srcType, targetID)
		if err != nil {
			return false, err
		}
		if hasSource {
			return true, nil
		}
	}

	return false, nil
}

//func (a *usage) MarshalJSON() ([]byte, error) {
//	return json.Marshal(a)
//}

func (a *usage) Combine(b NameUsage) (NameUsage, error) {

	if !a.Id.Valid() {
		return nil, errors.Newf("Invalid ID when combining NameUsages [%s]", a.Id)
	}

	c := usage{
		Id: a.Id,
	}

	// Slow recalculate this but necessary for clean code.
	namesEquivalent := a.Cn.Equals(b.CanonicalName())

	aSynonyms, err := a.Synonyms()
	if err != nil {
		return nil, err
	}

	bSynonyms, err := b.Synonyms()
	if err != nil {
		return nil, err
	}

	bNameIsSynonym := aSynonyms.Contains(b.CanonicalName())
	aNameIsSynonym := bSynonyms.Contains(a.CanonicalName())

	aSrcs, err := a.Sources()
	if err != nil {
		return nil, err
	}

	bSrcs, err := b.Sources()
	if err != nil {
		return nil, err
	}

	aSourceCount := len(aSrcs)
	bSourceCount := len(bSrcs)

	//fmt.Println("----Combining----")
	//fmt.Println( a.canonicalName.name, b.canonicalName.name)
	//fmt.Println("Synonymous", aNameIsSynonym, bNameIsSynonym)
	//fmt.Println("SourceCount", aSourceCount, bSourceCount)
	//targetIDsIntersect := a.nameUsageSourceMap.Intersects(b.nameUsageSourceMap)
	//synonymsIntersect := utils.IntersectsStrings(a.Synonyms, b.Synonyms)

	if bNameIsSynonym && aNameIsSynonym {
		fmt.Println(fmt.Sprintf("Warning: found two name usages [%s, %s] that appear to be synonyms for each other.", a.CanonicalName(), b.CanonicalName()))
		return nil, nil
	}

	switch{
	case namesEquivalent:
		c.Cn = b.CanonicalName()
	case bNameIsSynonym && !aNameIsSynonym:
		c.Cn = a.CanonicalName()
	case aNameIsSynonym && !bNameIsSynonym:
		c.Cn = b.CanonicalName()
	case aSourceCount > bSourceCount:
		c.Cn = a.CanonicalName()
	case aSourceCount < bSourceCount:
		c.Cn = b.CanonicalName()
	case len(aSynonyms) > len(bSynonyms):
		c.Cn = a.CanonicalName()
	default:
		c.Cn = b.CanonicalName()
	}

	if err := c.AddSources(aSrcs...); err != nil {
		return nil, err
	}
	if err := c.AddSources(bSrcs...); err != nil {
		return nil, err
	}

	return &c, nil
}