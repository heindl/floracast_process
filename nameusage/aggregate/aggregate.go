package aggregate

import (
	"sync"
	"encoding/json"
	"gopkg.in/tomb.v2"
	"context"
	"bitbucket.org/heindl/processors/nameusage/canonicalname"
	"bitbucket.org/heindl/processors/nameusage/nameusage"
	"bitbucket.org/heindl/processors/datasources"
	"bitbucket.org/heindl/processors/store"
	"bitbucket.org/heindl/processors/algolia"
	"bitbucket.org/heindl/processors/taxa"
)

type Aggregate struct {
	list []*nameusage.NameUsage
	sync.Mutex
}

type EachFunction func(ctx context.Context, usage *nameusage.NameUsage) error

func (Ω *Aggregate) Each(ctx context.Context, handler EachFunction) error {
	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _i := range Ω.list {
			i := _i
			tmb.Go(func() error {
				return handler(ctx, Ω.list[i])
			})
		}
		return nil
	})
	return tmb.Wait()
}


type FilterFunction func(usage *nameusage.NameUsage) bool

func (Ω *Aggregate) Filter(shouldFilter FilterFunction) (*Aggregate, error) {
	res := Aggregate{}
	for _, u := range Ω.list {
		if shouldFilter(u) {
			continue
		}
		if err := res.AddUsage(u); err != nil {
			return nil, err
		}
	}
	return &res, nil
}

func (Ω *Aggregate) Upload(cxt context.Context, florastore store.FloraStore) error {

	return Ω.Each(cxt, func(ctx context.Context, usage *nameusage.NameUsage) error {
		deletedUsageIDs, err := usage.Upload(ctx, florastore)
		if err != nil {
			return err
		}
		if err := algolia.UploadNameUsageObjects(ctx, florastore, usage, deletedUsageIDs); err != nil {
			return err
		}
		return taxa.UploadMaterializedTaxa(ctx, florastore, usage, deletedUsageIDs)
	})

}

func (Ω *Aggregate) ScientificNames() []string {
	res := canonicalname.CanonicalNames{}
	for _, l := range Ω.list {
		res = res.AddToSet(l.CanonicalName())
		res = res.AddToSet(l.Synonyms()...)
	}
	return res.ScientificNames()
}

func (Ω *Aggregate) TargetIDs(sourceTypes ...datasources.SourceType) (res datasources.TargetIDs) {
	for _, usage := range Ω.list {
		for _, src := range usage.Sources(sourceTypes...) {
			res = append(res, src.TargetID())
		}
	}
	return
}

func (Ω *Aggregate) MarshalJSON() ([]byte, error) {
	if Ω == nil {
		return nil, nil
	}
	return json.Marshal(Ω.list)
}

func (Ω *Aggregate) AddUsage(usages ...*nameusage.NameUsage) error {
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
				if Ω.list[i].ShouldCombine(Ω.list[k]) {
					changed = true
					var err error
					Ω.list[i], err = Ω.list[i].Combine(Ω.list[k])
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

func (Ω *Aggregate) HasCanonicalName(name *canonicalname.CanonicalName) bool {
	for i := range Ω.list {
		if Ω.list[i].HasScientificName(name.ScientificName()) {
			return true
		}
	}
	return false
}