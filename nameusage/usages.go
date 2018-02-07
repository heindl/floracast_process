package nameusage

import (
	"sync"
	"encoding/json"
	"bitbucket.org/heindl/taxa/datasources"
	"gopkg.in/tomb.v2"
	"context"
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

type ForEachModifierFunction func(ctx context.Context, usage *CanonicalNameUsage) error

func (Ω *AggregateNameUsages) ForEach(ctx context.Context, modify ForEachModifierFunction) error {
	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _i := range Ω.list {
			i := _i
			tmb.Go(func() error {
				return modify(ctx, Ω.list[i])
			})
		}
		return nil
	})
	return tmb.Wait()
}


type ShouldFilterUsageFunction func(usage *CanonicalNameUsage) bool

func (Ω *AggregateNameUsages) Filter(shouldFilter ShouldFilterUsageFunction) (*AggregateNameUsages, error) {
	res := AggregateNameUsages{}
	for _, u := range Ω.list {
		if shouldFilter(u) {
			continue
		}
		if err := res.AddUsages(u); err != nil {
			return nil, err
		}
	}
	return &res, nil
}

func (Ω *AggregateNameUsages) CanonicalNames() CanonicalNames {
	res := CanonicalNames{}
	for _, l := range Ω.list {
		res = res.AddToSet(l.ScientificNames()...)
	}
	return res
}

func (Ω *AggregateNameUsages) MarshalJSON() ([]byte, error) {
	if Ω == nil {
		return nil, nil
	}
	return json.Marshal(Ω.list)
}

func (Ω *AggregateNameUsages) CombineWith(usage *AggregateNameUsages) error {
	return Ω.AddUsages(usage.list...)

}

func (Ω *AggregateNameUsages) AddUsages(usages ...*CanonicalNameUsage) error {
	Ω.Lock()
	defer Ω.Unlock()

	Ω.list = append(Ω.list, usages...)

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


func (Ω *AggregateNameUsages) TargetIDs(srcType datasources.DataSourceType) datasources.DataSourceTargetIDs {
	res := datasources.DataSourceTargetIDs{}
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

func (Ω AggregateNameUsages) FirstIndexOfID(src datasources.DataSourceType, id datasources.DataSourceTargetID) int {
	for i := range Ω.list {
		if _, hasType := Ω.list[i].sources[src]; hasType {
			if _, hasTarget := Ω.list[i].sources[src][id]; hasTarget {
				return i
			}
		}
	}
	return -1
}
