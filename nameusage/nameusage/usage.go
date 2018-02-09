package nameusage

import (
	"fmt"
	"github.com/mongodb/mongo-tools/common/json"
	"github.com/elgs/gostrgen"
	"github.com/dropbox/godropbox/errors"
	"sort"
	"bitbucket.org/heindl/taxa/datasources"
	"bitbucket.org/heindl/taxa/occurrences"
	"bitbucket.org/heindl/taxa/nameusage/canonicalname"
	"bitbucket.org/heindl/taxa/nameusage/nameusagesource"
	"bitbucket.org/heindl/taxa/utils"
)


type NameUsage struct {
	id                    string
	canonicalName         *canonicalname.CanonicalName
	sources               nameUsageSourceMap
	occurrenceAggregation *occurrences.OccurrenceAggregation
}

type nameUsageSourceMap map[datasources.SourceType]map[datasources.TargetID]*nameusagesource.Source

func NewNameUsage(src *nameusagesource.Source) (*NameUsage, error) {

	id, err := gostrgen.RandGen(20, gostrgen.Lower|gostrgen.Digit|gostrgen.Upper, "", "")
	if err != nil {
		return nil, errors.Wrap(err, "Could not generate name usage id")
	}

	u := NameUsage{
		id:                    id,
		canonicalName:         src.CanonicalName(),
		occurrenceAggregation: occurrences.NewOccurrenceAggregation(),
	}

	if err := u.AddSources(src); err != nil {
		return nil, err
	}

	return &u, nil
}

func (Ω *NameUsage) ID() string {
	return Ω.id
}

func (Ω *NameUsage) Sources(sourceTypes ...datasources.SourceType) (res nameusagesource.Sources) {
	for srcType, targets := range Ω.sources {
		if len(sourceTypes) > 0 && !datasources.HasDataSourceType(sourceTypes, srcType) {
			continue
		}
		for _, src := range targets {
			res = append(res, src)
		}
	}
	return
}

func (Ω *NameUsage) TotalOccurrenceCount(sourceTypes ...datasources.SourceType) int {
	count := 0
	for _, src := range Ω.Sources(sourceTypes...) {
		count += src.OccurrenceCount()
	}
	return count
}

func (Ω *NameUsage) hasSource(sourceType datasources.SourceType, targetID datasources.TargetID) bool {

	if _, ok := Ω.sources[sourceType]; !ok {
		return false
	}

	if _, ok := Ω.sources[sourceType][targetID]; !ok {
		return false
}

	return true
}

func (Ω *NameUsage) MarshalJSON() ([]byte, error) {
	if Ω == nil {
		return nil, nil
	}
	return json.Marshal(map[string]interface{}{
		"CanonicalName": Ω.CanonicalName().ScientificName(),
		"Synonyms": Ω.Synonyms().ScientificNames(),
		"TotalOccurrenceCount": Ω.TotalOccurrenceCount(),
		"Sources": Ω.sources,
	})
}

func (Ω *NameUsage) AddSources(sources ...*nameusagesource.Source) error {
	for _, src := range sources {

		if _, ok := Ω.sources[src.SourceType()]; !ok {
			Ω.sources[src.SourceType()] = map[datasources.TargetID]*nameusagesource.Source{}
		}
		Ω.sources[src.SourceType()][src.TargetID()] = src

	}
	return nil
}

func (Ω *NameUsage) HasScientificName(name string) bool {
	names := utils.AddStringToSet(Ω.Synonyms().ScientificNames(), Ω.CanonicalName().ScientificName())
	return utils.ContainsString(names, name)
}

func (Ω *NameUsage) Synonyms() canonicalname.CanonicalNames {
	res := canonicalname.CanonicalNames{}
	for _, src := range Ω.Sources() {
		for _, s := range src.Synonyms().AddToSet(src.CanonicalName()) {
			if s.Equals(Ω.canonicalName) {
				continue
			}
			res = res.AddToSet(s)
		}
	}
	return res
}

func (Ω *NameUsage) ScientificNameReferenceLedger() NameReferenceLedger {
	ledger := NameReferenceLedger{}
	for _, src := range Ω.Sources() {
		ledger = ledger.IncrementName(src.CanonicalName().ScientificName(), src.OccurrenceCount())
	}
	sort.Sort(ledger)
	return ledger
}

func (Ω *NameUsage) CommonNameReferenceLedger() NameReferenceLedger {
	ledger := NameReferenceLedger{}
	for _, src := range Ω.Sources() {
		for _, cn := range src.CommonNames() {
			ledger = ledger.IncrementName(cn, src.OccurrenceCount())
		}
	}
	sort.Sort(ledger)
	return ledger
}

// Get the most popular common name from all sources
// Or perhaps the most common among those those with Occurrence Counts?
func (Ω *NameUsage) CommonName() (string, error) {

	ledger := Ω.CommonNameReferenceLedger()

	if len(ledger) == 0 || ledger[0].Name == "" {
		return "", errors.New("CommonName not found")
	}

	return ledger[0].Name, nil
}

func (Ω *NameUsage) CanonicalName() *canonicalname.CanonicalName {
	return Ω.canonicalName
}

func (a *NameUsage) ShouldCombine(b *NameUsage) bool {

	for _, s := range append(b.Synonyms().ScientificNames(), b.CanonicalName().ScientificName()) {
		if a.HasScientificName(s) {
			return true
		}
	}

	for _, src := range a.Sources() {
		if b.hasSource(src.SourceType(), src.TargetID()) {
			return true
		}
	}

	return false
}

func (a *NameUsage) Combine(b *NameUsage) (*NameUsage, error) {
	c := NameUsage{}

	// Slow recalculate this but necessary for clean code.
	namesEquivalent := a.canonicalName.Equals(b.canonicalName)
	bNameIsSynonym := a.Synonyms().Contains(b.canonicalName)
	aNameIsSynonym := b.Synonyms().Contains(a.canonicalName)

	aSourceCount := len(a.Sources())
	bSourceCount := len(b.Sources())

	//fmt.Println("----Combining----")
	//fmt.Println( a.canonicalName.name, b.canonicalName.name)
	//fmt.Println("Synonymous", aNameIsSynonym, bNameIsSynonym)
	//fmt.Println("SourceCount", aSourceCount, bSourceCount)
	//targetIDsIntersect := a.nameUsageSourceMap.Intersects(b.nameUsageSourceMap)
	//synonymsIntersect := utils.IntersectsStrings(a.Synonyms, b.Synonyms)

	if bNameIsSynonym && aNameIsSynonym {
		fmt.Println(fmt.Sprintf("Warning: found two name usages [%s, %s] that appear to be synonyms for each other.", a.canonicalName, b.canonicalName))
		return nil, nil
		}

	switch{
	case namesEquivalent:
		c.canonicalName = b.canonicalName
	case bNameIsSynonym && !aNameIsSynonym:
		c.canonicalName = a.canonicalName
	case aNameIsSynonym && !bNameIsSynonym:
		c.canonicalName = b.canonicalName
	case aSourceCount > bSourceCount:
		c.canonicalName = a.canonicalName
	case aSourceCount < bSourceCount:
		c.canonicalName = b.canonicalName
	case len(a.Synonyms()) > len(b.Synonyms()):
		c.canonicalName = a.canonicalName
	default:
		c.canonicalName = b.canonicalName
	}

	c.sources = nameUsageSourceMap{}

	if err := c.AddSources(a.Sources()...); err != nil {
		return nil, err
	}
	if err := c.AddSources(b.Sources()...); err != nil {
		return nil, err
	}

	return &c, nil
}