package nameusage

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/nameusage/canonicalname"
	"bitbucket.org/heindl/process/store"
	"bitbucket.org/heindl/process/utils"
	"context"
	"encoding/json"
	"github.com/dropbox/godropbox/errors"
	"sort"
)

// NameUsage represents a combination of DataSources that leads to a Taxon representation.
type NameUsage interface {
	ID() (ID, error)
	Sources(sourceTypes ...datasources.SourceType) (Sources, error)
	Occurrences(sourceTypes ...datasources.SourceType) (int, error)
	HasSource(datasources.SourceType, datasources.TargetID) (bool, error)
	AddSources(sources ...Source) error
	Synonyms() (canonicalname.Names, error)
	AllScientificNames() ([]string, error)
	CanonicalName() *canonicalname.Name
	CommonName() (string, error)
	ShouldCombine(b NameUsage) (bool, error)
	Combine(b NameUsage) (NameUsage, error)
	HasScientificName(s string) (bool, error)
	ScientificNameReferenceLedger() (nameReferenceLedger, error)
	CommonNameReferenceLedger() (nameReferenceLedger, error)
	Upload(context.Context, store.FloraStore) (deletedUsageIDs IDs, err error)
}

type ByCanonicalName []NameUsage

func (a ByCanonicalName) Len() int      { return len(a) }
func (a ByCanonicalName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByCanonicalName) Less(i, j int) bool {
	return a[i].CanonicalName().ScientificName() < a[j].CanonicalName().ScientificName()
}

const storeKeyScientificName = "ScientificNames"

type usage struct {
	id       ID
	Cn       *canonicalname.Name                                         `json:"Name" firestore:"Name"`
	Occrrncs int                                                         `json:"Occurrences,omitempty" firestore:"Occurrences,omitempty"`
	SciNames map[string]bool                                             `json:"ScientificNames,omitempty" firestore:"ScientificNames,omitempty"`
	Srcs     map[datasources.SourceType]map[datasources.TargetID]*source `json:"Sources,omitempty" firestore:"Sources,omitempty"`
}

// NewNameUsage provides a new NameUsage.
func NewNameUsage(src Source) (NameUsage, error) {

	id, err := NewNameUsageID()
	if err != nil {
		return nil, err
	}

	u := usage{
		id: id,
		Cn: src.CanonicalName(),
	}

	if err := u.AddSources(src); err != nil {
		return nil, err
	}

	return &u, nil
}

// FromJSON parses a NameUsage from json.
func FromJSON(id ID, b []byte) (NameUsage, error) {
	if !id.Valid() {
		return nil, errors.Newf("Invalid ID [%s]", id)
	}
	u := usage{}
	if err := json.Unmarshal(b, &u); err != nil {
		return nil, err
	}
	u.id = id

	return &u, nil
}

// ID provides the NameUsageID from a NameUsage.
func (Ω *usage) ID() (ID, error) {
	if !Ω.id.Valid() {
		return ID(""), errors.Newf("Invalid ID [%s]", Ω.id)
	}
	return Ω.id, nil
}

// Sources returns all sources from the NameUsage
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

// Occurrences returns a sum of all occurrences within NameUsage.
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

func (Ω *usage) CanonicalName() *canonicalname.Name {
	return Ω.Cn
}

func (Ω *usage) Synonyms() (canonicalname.Names, error) {
	res := canonicalname.Names{}
	srcs, err := Ω.Sources()
	if err != nil {
		return nil, err
	}
	for _, src := range srcs {
		for _, s := range src.Synonyms().AddToSet(src.CanonicalName()) {
			if Ω.Cn != nil && s != nil && s.Equals(Ω.Cn) {
				continue
			}
			res = res.AddToSet(s)
		}
	}
	return res, nil
}

func (Ω *usage) AllScientificNames() ([]string, error) {
	synonyms, err := Ω.Synonyms()
	if err != nil {
		return nil, err
	}
	return utils.AddStringToSet(synonyms.ScientificNames(), Ω.CanonicalName().ScientificName()), nil
}

func (Ω *usage) ScientificNameReferenceLedger() (nameReferenceLedger, error) {
	srcs, err := Ω.Sources()
	if err != nil {
		return nil, err
	}
	ledger := nameReferenceLedger{}
	for _, src := range srcs {
		ledger = ledger.IncrementName(src.CanonicalName().ScientificName(), src.OccurrenceCount())
		for _, synonym := range src.Synonyms() {
			ledger = ledger.IncrementName(synonym.ScientificName(), 0)
		}
	}
	sort.Sort(ledger)
	return ledger, nil
}

func (Ω *usage) CommonNameReferenceLedger() (nameReferenceLedger, error) {
	ledger := nameReferenceLedger{}
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

func (Ω *usage) ShouldCombine(b NameUsage) (bool, error) {

	aNames, err := Ω.AllScientificNames()
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

	srcs, err := Ω.Sources()
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

// Slow recalculate this but necessary for clean code.
func shouldUseFirstCanonicalNameUsage(a, b NameUsage) (bool, error) {

	if a.CanonicalName().Equals(b.CanonicalName()) {
		return true, nil
	}

	aSynonyms, err := a.Synonyms()
	if err != nil {
		return false, err
	}

	bSynonyms, err := b.Synonyms()
	if err != nil {
		return false, err
	}

	bNameIsSynonym := aSynonyms.Contains(b.CanonicalName())
	aNameIsSynonym := bSynonyms.Contains(a.CanonicalName())

	if bNameIsSynonym != aNameIsSynonym {
		return aNameIsSynonym, nil
	} else if bNameIsSynonym {
		return false, errors.Newf("Warning: found two name usages [%s, %s] that appear to be synonyms for each other.", a.CanonicalName(), b.CanonicalName())
	}

	aSources, err := a.Sources()
	if err != nil {
		return false, err
	}

	bSources, err := b.Sources()
	if err != nil {
		return false, err
	}

	return len(aSources) > len(bSources) || len(aSynonyms) > len(bSynonyms), nil

	//fmt.Println("----Combining----")
	//fmt.Println( a.canonicalName.name, b.canonicalName.name)
	//fmt.Println("Synonymous", aNameIsSynonym, bNameIsSynonym)
	//fmt.Println("SourceCount", aSourceCount, bSourceCount)
	//targetIDsIntersect := a.nameUsageSourceMap.Intersects(b.nameUsageSourceMap)
	//synonymsIntersect := utils.IntersectsStrings(a.Synonyms, b.Synonyms)
}

func (Ω *usage) Combine(b NameUsage) (NameUsage, error) {

	if !Ω.id.Valid() {
		return nil, errors.Newf("Invalid ID when combining NameUsages [%s]", Ω.id)
	}

	c := usage{
		id:   Ω.id,
		Srcs: Ω.Srcs,
	}

	useFirstName, err := shouldUseFirstCanonicalNameUsage(Ω, b)
	if err != nil {
		return nil, err
	}

	if useFirstName {
		c.Cn = Ω.CanonicalName()
	} else {
		c.Cn = b.CanonicalName()
	}

	aSources, err := Ω.Sources()
	if err != nil {
		return nil, err
	}

	bSources, err := b.Sources()
	if err != nil {
		return nil, err
	}

	if err := c.AddSources(aSources...); err != nil {
		return nil, err
	}

	if err := c.AddSources(bSources...); err != nil {
		return nil, err
	}

	return &c, nil
}
