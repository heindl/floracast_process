package nameusage

import (
	"fmt"
	"github.com/mongodb/mongo-tools/common/json"
	"github.com/elgs/gostrgen"
	"github.com/dropbox/godropbox/errors"
	"strings"
	"sort"
	"bitbucket.org/heindl/taxa/datasources"
	"bitbucket.org/heindl/taxa/utils"
	"time"
	"bitbucket.org/heindl/taxa/occurrences"
)

type CanonicalNameUsage struct {
	id                    string
	canonicalName         *CanonicalName
	sources               nameUsageSourceMap
	occurrenceAggregation *occurrences.OccurrenceAggregation
}

func NewCanonicalNameUsage(src *NameUsageSource) (*CanonicalNameUsage, error) {

	id, err := gostrgen.RandGen(20, gostrgen.Lower|gostrgen.Digit|gostrgen.Upper, "", "")
	if err != nil {
		return nil, errors.Wrap(err, "Could not generate name usage id")
	}

	return &CanonicalNameUsage{
		id:                    id,
		canonicalName:         src.canonicalName,
		occurrenceAggregation: occurrences.NewOccurrenceAggregation(),
		sources: nameUsageSourceMap{
			src.sourceType: map[datasources.DataSourceTargetID]*NameUsageSource{
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

func (Ω *CanonicalNameUsage) TotalOccurrenceCount() int {
	return Ω.sources.totalOccurrenceCount()
}

func (Ω *CanonicalNameUsage) CachedOccurrences() []*occurrences.Occurrence {
	return Ω.occurrenceAggregation.OccurrenceList()
}

func (Ω *CanonicalNameUsage) hasSource(sourceType datasources.DataSourceType, targetID datasources.DataSourceTargetID) bool {

	if _, ok := Ω.sources[sourceType]; !ok {
		return false
	}

	if _, ok := Ω.sources[sourceType][targetID]; !ok {
		return false

		}

	return true
}

func (Ω *CanonicalNameUsage) registerOccurrenceFetch(sourceType datasources.DataSourceType, targetID datasources.DataSourceTargetID, occurrenceCount int) error {

	if _, ok := Ω.sources[sourceType]; !ok {
		return errors.Newf("CanonicalNameUsage [%s] does not contain sourceType [%s]", Ω.canonicalName, sourceType)
	}

	if _, ok := Ω.sources[sourceType][targetID]; !ok {
		return errors.Newf("CanonicalNameUsage [%s] does not contain sourceType [%s] targetID [%s]", Ω.canonicalName, sourceType, targetID)
	}

	Ω.sources[sourceType][targetID].occurrenceCount = Ω.sources[sourceType][targetID].occurrenceCount + occurrenceCount
	Ω.sources[sourceType][targetID].lastFetchedAt = utils.TimePtr(time.Now())

	return nil
}

func (Ω *CanonicalNameUsage) MarshalJSON() ([]byte, error) {
	if Ω == nil {
		return nil, nil
	}
	return json.Marshal(map[string]interface{}{
		"CanonicalName": Ω.canonicalName,
		"Synonyms": Ω.ScientificNames().Strings(),
		"TotalOccurrenceCount": Ω.sources.totalOccurrenceCount(),
		"Sources": Ω.sources,
	})
}

func (Ω *CanonicalNameUsage) AddSources(sources ...*NameUsageSource) error {
	for _, src := range sources {
		Ω.sources.set(src)
	}
	return nil
}

func (Ω *CanonicalNameUsage) AddOccurrence(o *occurrences.Occurrence) error {

	if o == nil {
		return nil
	}

	if !Ω.hasSource(o.SourceType(), o.TargetID()) {
		return errors.Newf("CanonicalNameUsage does not contain source [%s, %s]", o.SourceType(), o.TargetID())
	}

	Ω.sources[o.SourceType()][o.TargetID()].occurrenceCount += 1
	Ω.sources[o.SourceType()][o.TargetID()].lastFetchedAt = utils.TimePtr(time.Now())

	if Ω.occurrenceAggregation == nil {
		Ω.occurrenceAggregation = occurrences.NewOccurrenceAggregation()
	}

	if err := Ω.occurrenceAggregation.AddOccurrence(o); err != nil && !utils.ContainsError(err, occurrences.ErrCollision) {
		return err
	}

	return nil
}

func (Ω *CanonicalNameUsage) Synonyms() CanonicalNames {
	res := CanonicalNames{}
	for _, cn := range Ω.ScientificNames() {
		if !cn.Equals(Ω.canonicalName) {
			res = res.AddToSet(cn)
		}
	}
	return res
}

func (Ω *CanonicalNameUsage) ContainsCanonicalName(cn *CanonicalName) bool {
	for _, targets := range Ω.sources {
		for _, src := range targets {
			if src.canonicalName.Equals(cn) || src.synonyms.Contains(cn) {
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

func (Ω *CanonicalNameUsage) Sources(sourceTypes ...datasources.DataSourceType) []*NameUsageSource {
	res := []*NameUsageSource{}
	for srcType, targets := range Ω.sources {
		if len(sourceTypes) > 0 && !datasources.HasDataSourceType(sourceTypes, srcType) {
			continue
		}
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

func (Ω *CanonicalNameUsage) CanonicalName() *CanonicalName {
	return Ω.canonicalName
}

func (Ω *CanonicalNameUsage) ScientificNames() CanonicalNames {
	res := CanonicalNames{}
	for _, src := range Ω.Sources() {
		res = res.AddToSet(src.canonicalName)
		res = res.AddToSet(src.synonyms...)
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
	for _, aName := range a.ScientificNames() {
		if b.ContainsCanonicalName(aName) {
			//fmt.Println(fmt.Sprintf("Should combine on CanonicalName [%s, %s]",a.CanonicalName(), b.CanonicalName()))
			return true
		}
	}
	if a.sources.intersects(b.sources) {
		//fmt.Println(fmt.Sprintf("Should combine on Sources [%s, %s]", a.CanonicalName(), b.CanonicalName()))
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
	c.sources.merge(a.sources)
	c.sources.merge(b.sources)

	return &c, nil
}