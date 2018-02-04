package nameusage

import (
	"sync"
	"encoding/json"
	"bitbucket.org/heindl/taxa/store"
)

type AggregateNameUsages struct {
	list CanonicalNameUsages
	sync.Mutex
}

type CanonicalNameUsages []*CanonicalNameUsage

func (a CanonicalNameUsages) Len() int           { return len(a) }
func (a CanonicalNameUsages) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CanonicalNameUsages) Less(i, j int) bool { return a[i].canonicalName.name < a[j].canonicalName.name }

func (Ω *AggregateNameUsages) Count() int {
	return len(Ω.list)
}

func (Ω *AggregateNameUsages) NameStrings() []string {
	res := []string{}
	for _, l := range Ω.list {
		res = append(res, l.canonicalName.name)
	}
	return res
}

func (Ω *AggregateNameUsages) MarshalJSON() ([]byte, error) {
	if Ω == nil {
		return nil, nil
	}
	return json.Marshal(Ω.list)
}

func (Ω *AggregateNameUsages) Add(usage *CanonicalNameUsage) error {
	Ω.Lock()
	defer Ω.Unlock()

	//foundCanonicalName := false
	//for i := range Ω.list {
	//	if Ω.list[i].canonicalName.Equals(usage.canonicalName) {
	//		foundCanonicalName = true
	//		var err error
	//		Ω.list[i], err = Ω.list[i].combine(usage)
	//		if err != nil {
	//			return err
	//		}
	//	}
	//}

	//if !foundCanonicalName {
	Ω.list = append(Ω.list, usage)
	//}

	//sort.Sort(Ω.list)

ResetLoop:
	for {
		changed := false
		for i := range Ω.list {
			for k := range Ω.list {
				if k == i {
					continue
				}
				if Ω.list[i].shouldCombine(Ω.list[k]) {
					changed = true
					var err error
					Ω.list[i], err = Ω.list[i].combine(Ω.list[k])
					if err != nil {
						return err
					}
					Ω.list = append(Ω.list[:k], Ω.list[k+1:]...)
					continue ResetLoop
				}
			}
		}
		if !changed {
			break
		}
	}

	return nil
}


func (Ω *AggregateNameUsages) TargetIDs(srcType store.DataSourceType) store.DataSourceTargetIDs {
	res := store.DataSourceTargetIDs{}
	for i := range Ω.list {
		res = res.AddToSet(Ω.list[i].sources.targetIDs(srcType)...)
	}
	return res
}

func (Ω AggregateNameUsages) FirstIndexOfName(name *CanonicalName) int {
	for i := range Ω.list {
		if Ω.list[i].ContainsCanonicalName(name) {
			return i
		}
	}
	return -1
}

func (Ω AggregateNameUsages) FirstIndexOfID(src store.DataSourceType, id store.DataSourceTargetID) int {
	for i := range Ω.list {
		if _, hasType := Ω.list[i].sources[src]; hasType {
			if _, hasTarget := Ω.list[i].sources[src][id]; hasTarget {
				return i
			}
		}
	}
	return -1
}
