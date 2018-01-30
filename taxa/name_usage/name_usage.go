package name_usage

import (
	"fmt"
	"bitbucket.org/heindl/taxa/utils"
	"bitbucket.org/heindl/taxa/store"
	"strings"
)

type CanonicalNameUsages []CanonicalNameUsage

func (Ω CanonicalNameUsages) IndexOfNames(qNames ...string) int {
	for i := range Ω {
		if utils.IntersectsStrings(append(Ω[i].Synonyms, Ω[i].CanonicalName), qNames) {
			return i
		}
	}
	return -1
}

func (Ω CanonicalNameUsages) Names() []string {
	res := []string{}
	for i := range Ω {
		res = utils.AddStringToSet(res, append(Ω[i].Synonyms, Ω[i].CanonicalName)...)
	}
	return res
}

func (Ω CanonicalNameUsages) TargetIDs(srcType store.DataSourceType) store.DataSourceTargetIDs {
	res := store.DataSourceTargetIDs{}
	for i := range Ω {
		res = res.AddToSet(Ω[i].SourceTargetOccurrenceCount.TargetIDs(srcType)...)
	}
	return res
}

func (Ω CanonicalNameUsages) Condense() (CanonicalNameUsages, error) {

ResetLoop:
	for {
		changed := false
		for i := range Ω {
			for k := range Ω {
				if k == i {
					continue
				}
				if Ω[i].ShouldCombine(Ω[k]) {
					changed = true
					var err error
					Ω[i], err = Ω[i].Combine(Ω[k])
					if err != nil {
						return nil, err
					}
					Ω = append(Ω[:k], Ω[k+1:]...)
					continue ResetLoop
				}
			}
		}
		if !changed {
			break
		}
	}

	return Ω, nil
}

func (Ω CanonicalNameUsages) FirstIndexOfName(name string) int {

	n := strings.ToLower(name)

	for i := range Ω {
		if Ω[i].CanonicalName == n {
			return i
		}

		if utils.ContainsString(Ω[i].Synonyms, n) {
			return i
		}
	}
	return -1
}

func (Ω CanonicalNameUsages) FirstIndexOfID(src store.DataSourceType, id store.DataSourceTargetID) int {
	for i := range Ω {
		if _, ok := Ω[i].SourceTargetOccurrenceCount[src]; ok {
			if _, ok := Ω[i].SourceTargetOccurrenceCount[src][id]; ok {
				return i
			}
		}
	}
	return -1
}

func (a *CanonicalNameUsage) ShouldCombine(b CanonicalNameUsage) bool {
	namesEqual := a.CanonicalName == b.CanonicalName && a.CanonicalName != ""
	if namesEqual {
		fmt.Println(fmt.Sprintf("Combining: names [%s, %s] are equivalent", a.CanonicalName, b.CanonicalName))
		return true
	}

	bNameIsSynonym := utils.ContainsString(a.Synonyms, b.CanonicalName)
	if bNameIsSynonym {
		fmt.Println(fmt.Sprintf("Combining: b [%s] is synonym of a [%s]", a.CanonicalName, b.CanonicalName))
		return true
	}

	aNameIsSynonym := utils.ContainsString(b.Synonyms, a.CanonicalName)
	if aNameIsSynonym {
		fmt.Println(fmt.Sprintf("Combining: a [%s] is synonym of b [%s]", a.CanonicalName, b.CanonicalName))
		return true
	}

	sharesSynonyms := utils.IntersectsStrings(a.Synonyms, b.Synonyms)
	if sharesSynonyms {
		fmt.Println(fmt.Sprintf("Combining: Synonyms intersect [%s, %s]", a.CanonicalName, b.CanonicalName))
		return true
	}

	if a.SourceTargetOccurrenceCount.Intersects(b.SourceTargetOccurrenceCount) {
		fmt.Println(fmt.Sprintf("Combining: Sources intersect [%s, %s]", a.CanonicalName, b.CanonicalName))
		return true
	}

	return false
}

func (a *CanonicalNameUsage) Combine(b CanonicalNameUsage) (CanonicalNameUsage, error) {
	c := CanonicalNameUsage{}

	// Slow recalculate this but necessary for clean code.
	namesEquivalent := a.CanonicalName == b.CanonicalName
	bNameIsSynonym := utils.ContainsString(a.Synonyms, b.CanonicalName)
	aNameIsSynonym := utils.ContainsString(b.Synonyms, a.CanonicalName)
	//targetIDsIntersect := a.SourceTargetOccurrenceCount.Intersects(b.SourceTargetOccurrenceCount)
	//synonymsIntersect := utils.IntersectsStrings(a.Synonyms, b.Synonyms)

	aSourceCount := a.CountSources()
	bSourceCount := b.CountSources()

	switch{
	case namesEquivalent:
		c.CanonicalName = b.CanonicalName
	case bNameIsSynonym && aNameIsSynonym:
		fmt.Println(fmt.Sprintf("Warning: found two name usages [%s, %s] that appear to be synonyms for each other.", a.CanonicalName, b.CanonicalName))
		c.CanonicalName = a.CanonicalName
	case bNameIsSynonym:
		c.CanonicalName = b.CanonicalName
	case aNameIsSynonym:
		c.CanonicalName = a.CanonicalName
	case aSourceCount > bSourceCount:
		c.CanonicalName = a.CanonicalName
	case aSourceCount < bSourceCount:
		c.CanonicalName = b.CanonicalName
	case len(a.Synonyms) > len(b.Synonyms):
		c.CanonicalName = a.CanonicalName
	default:
		c.CanonicalName = b.CanonicalName
	}

	c.Synonyms = utils.AddStringToSet(a.Synonyms, b.Synonyms...)
	c.Synonyms = utils.AddStringToSet(c.Synonyms, a.CanonicalName, b.CanonicalName)
	c.Ranks = utils.AddStringToSet(a.Ranks, b.Ranks...)
	b.SourceTargetOccurrenceCount.AddAll(a.SourceTargetOccurrenceCount)
	c.SourceTargetOccurrenceCount = b.SourceTargetOccurrenceCount

	return c, nil
}

func (Ω *CanonicalNameUsage) Valid() bool {
	Ω.Ranks = utils.StringsToLower(Ω.Ranks...)
	if !utils.IntersectsStrings([]string{"species", "form", "subspecies", "variety"}, Ω.Ranks) {
		return false
	}
	return true
}

func (Ω *CanonicalNameUsage) CountSources() int {
	return Ω.SourceTargetOccurrenceCount.TargetIDCount()
}

func (Ω *CanonicalNameUsage) OccurrenceCount() int {
	return Ω.SourceTargetOccurrenceCount.TotalOccurrenceCount()
}