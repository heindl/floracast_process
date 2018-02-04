package nameusage

import (
	"fmt"
	"bitbucket.org/heindl/taxa/store"
	"github.com/mongodb/mongo-tools/common/json"
	"github.com/elgs/gostrgen"
	"github.com/dropbox/godropbox/errors"
	"strings"
	"sort"
)

type CanonicalNameUsage struct {
	id string
	canonicalName   *CanonicalName
	sources         nameUsageSourceMap
}

func NewCanonicalNameUsage(src *NameUsageSource) (*CanonicalNameUsage, error) {

	id, err := gostrgen.RandGen(20, gostrgen.Lower|gostrgen.Digit|gostrgen.Upper, "", "")
	if err != nil {
		return nil, errors.Wrap(err, "Could not generate name usage id")
	}

	return &CanonicalNameUsage{
		id: id,
		canonicalName: src.canonicalName,
		sources: nameUsageSourceMap{
			src.sourceType: map[store.DataSourceTargetID]*NameUsageSource{
				src.targetID: src,
			},
		},
	}, nil
}

func (Ω *CanonicalNameUsage) ID() string {
	return Ω.id
}

func (Ω *CanonicalNameUsage) SourceCount() int {
	return Ω.sources.targetIDCount()
}

func (Ω *CanonicalNameUsage) OccurrenceCount() int {
	Ω.sources.totalOccurrenceCount()
}

func (Ω *CanonicalNameUsage) MarshalJSON() ([]byte, error) {
	if Ω == nil {
		return nil, nil
	}
	return json.Marshal(map[string]interface{}{
		"CanonicalName": Ω.canonicalName,
		"Synonyms": Ω.CanonicalNames().Strings(),
		"OccurrenceCount": Ω.sources.totalOccurrenceCount(),
		"Sources": Ω.sources,
	})
}

func (Ω *CanonicalNameUsage) AddSource(src *NameUsageSource) error {
	Ω.sources.set(src)
	return nil
}

func (Ω *CanonicalNameUsage) Synonyms() CanonicalNames {
	res := CanonicalNames{}
	for _, cn := range Ω.CanonicalNames() {
		if !cn.Equals(Ω.canonicalName) {
			res = res.AddToSet(cn)
		}
	}
	return res
}

func (Ω *CanonicalNameUsage) ContainsCanonicalName(name *CanonicalName) bool {
	for _, targets := range Ω.sources {
		for _, src := range targets {
			if src.canonicalName.Equals(name) {
				return true
			}
		}
	}
	return false
}

func (Ω *CanonicalNameUsage) ScientificNameString() (string, error) {
	if strings.TrimSpace(Ω.canonicalName.name) == "" {
		return "", errors.New("ScientificName not found")
	}
	return Ω.canonicalName.name, nil
}

func (Ω *CanonicalNameUsage) Sources() []*NameUsageSource {
	res := []*NameUsageSource{}
	for _, targets := range Ω.sources {
		for _, src := range targets {
			res = append(res, src)
		}
	}
	return res
}

func (Ω *CanonicalNameUsage) ScientificNameReferenceLedger() NameReferenceLedger {
	ledger := NameReferenceLedger{}
	for _, src := range Ω.Sources() {
		ledger = ledger.IncrementName(src.ScientificName(), src.occurrenceCount)
	}
	sort.Sort(ledger)
	return ledger
}

func (Ω *CanonicalNameUsage) CommonNameReferenceLedger() NameReferenceLedger {
	ledger := NameReferenceLedger{}
	for _, src := range Ω.Sources() {
		for _, cn := range src.CommonNames() {
			ledger = ledger.IncrementName(cn, src.occurrenceCount)
		}
	}
	sort.Sort(ledger)
	return ledger
}

// Get the most popular common name from all sources
// Or perhaps the most common among those those with Occurrence Counts?
func (Ω *CanonicalNameUsage) CommonNameString() (string, error) {

	ledger := Ω.CommonNameReferenceLedger()

	if len(ledger) == 0 || ledger[0].Name == "" {
		return "", errors.New("CommonName not found")
	}

	return ledger[0].Name, nil
}

func (Ω *CanonicalNameUsage) CanonicalNames() CanonicalNames {
	res := CanonicalNames{}
	for _, src := range Ω.Sources() {
		res = res.AddToSet(src.canonicalName)
	}
	return res
}

func (Ω *CanonicalNameUsage) AddSynonyms(sources ...*NameUsageSource) error {
	for _, src := range sources {
		Ω.sources.set(src)
	}
	return nil
}


func (a *CanonicalNameUsage) shouldCombine(b *CanonicalNameUsage) bool {
	for _, aName := range a.CanonicalNames() {
		if b.ContainsCanonicalName(aName) {
			fmt.Println()
			return true
		}
	}
	if a.sources.intersects(b.sources) {
		return true
	}
	return false
}

func (a *CanonicalNameUsage) combine(b *CanonicalNameUsage) (*CanonicalNameUsage, error) {
	c := CanonicalNameUsage{}

	// Slow recalculate this but necessary for clean code.
	namesEquivalent := a.canonicalName.Equals(b.canonicalName)
	bNameIsSynonym := a.Synonyms().Contains(b.canonicalName)
	aNameIsSynonym := b.Synonyms().Contains(a.canonicalName)

	aSourceCount := a.SourceCount()
	bSourceCount := b.SourceCount()

	fmt.Println("----Combining----")
	fmt.Println( a.canonicalName.name, b.canonicalName.name)
	fmt.Println("Synonymous", aNameIsSynonym, bNameIsSynonym)
	fmt.Println("SourceCount", aSourceCount, bSourceCount)
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
	c.sources.merge(a.sources)
	c.sources.merge(b.sources)

	fmt.Println("final", c.canonicalName.name)

	return &c, nil
}