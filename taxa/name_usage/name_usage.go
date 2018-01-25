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

func (Ω CanonicalNameUsages) TargetIDs(srcType store.DataSourceID) store.DataSourceTargetIDs {
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

func (Ω CanonicalNameUsages) FirstIndexOfID(src store.DataSourceID, id store.DataSourceTargetID) int {
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
		fmt.Println("should combine on names", a.CanonicalName, b.CanonicalName)
		return true
	}

	sharesSynonyms := utils.IntersectsStrings(a.Synonyms, b.Synonyms)
	if sharesSynonyms {
		fmt.Println("should combine on synonyms", a.CanonicalName, b.CanonicalName)
		return true
	}

	bNameIsSynonym := utils.ContainsString(a.Synonyms, b.CanonicalName)
	aNameIsSynonym := utils.ContainsString(b.Synonyms, a.CanonicalName)
	if bNameIsSynonym || aNameIsSynonym {
		fmt.Println("should combine on one name is the synonym of the other", a.CanonicalName, b.CanonicalName)
		return true
	}

	if a.SourceTargetOccurrenceCount.Intersects(b.SourceTargetOccurrenceCount) {
		fmt.Println("should combine on targetIDs", a.CanonicalName, b.CanonicalName)
		return true
	}

	return false
}

func (a *CanonicalNameUsage) Combine(b CanonicalNameUsage) (CanonicalNameUsage, error) {
	c := CanonicalNameUsage{}

	bNameIsSynonym := utils.ContainsString(a.Synonyms, b.CanonicalName)
	aNameIsSynonym := utils.ContainsString(b.Synonyms, a.CanonicalName)

	if bNameIsSynonym && aNameIsSynonym {
		fmt.Println(fmt.Sprintf("Warning: found two name usages [%s, %s] that appear to be synonyms  for each other.", a.CanonicalName, b.CanonicalName))
	}

	if bNameIsSynonym {
		c.CanonicalName = a.CanonicalName
		c.Synonyms = append(c.Synonyms, b.CanonicalName)
	} else if aNameIsSynonym {
		c.CanonicalName = b.CanonicalName
		c.Synonyms = append(c.Synonyms, a.CanonicalName)
	} else {
		c.CanonicalName = b.CanonicalName
	}

	c.Synonyms = utils.RemoveStringDuplicates(append(c.Synonyms, append(a.Synonyms, b.Synonyms...)...))
	c.Ranks = utils.RemoveStringDuplicates(append(a.Ranks, b.Ranks...))
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
