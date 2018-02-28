package aggregate

import (
	"bitbucket.org/heindl/process/algolia"
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/nameusage/canonicalname"
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/store"
	"bitbucket.org/heindl/process/taxa"
	"bitbucket.org/heindl/process/utils"
	"context"
	"gopkg.in/tomb.v2"
	"sync"
)

type Aggregate struct {
	list []nameusage.NameUsage
	sync.Mutex
}

type EachFunction func(ctx context.Context, usage nameusage.NameUsage) error

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

type FilterFunction func(usage nameusage.NameUsage) (bool, error)

func (Ω *Aggregate) Filter(shouldFilter FilterFunction) (*Aggregate, error) {
	res := Aggregate{}
	for _, u := range Ω.list {
		should, err := shouldFilter(u)
		if err != nil {
			return nil, err
		}
		if should {
			continue
		}
		if err := res.AddUsage(u); err != nil {
			return nil, err
		}
	}
	return &res, nil
}

func (Ω *Aggregate) Upload(cxt context.Context, florastore store.FloraStore) error {

	return Ω.Each(cxt, func(ctx context.Context, usage nameusage.NameUsage) error {
		deletedUsageIDs, err := usage.Upload(ctx, florastore)
		if err != nil {
			return err
		}
		if err := algolia.UploadNameUsageObjects(ctx, florastore, usage, deletedUsageIDs...); err != nil {
			return err
		}
		return taxa.UploadMaterializedTaxa(ctx, florastore, usage, deletedUsageIDs...)
	})

}

func (Ω *Aggregate) Count() int {
	return len(Ω.list)
}

func (Ω *Aggregate) Occurrences() (int, error) {
	res := 0
	for _, l := range Ω.list {
		i, err := l.Occurrences()
		if err != nil {
			return 0, err
		}
		res += i
	}
	return res, nil
}

func (Ω *Aggregate) ScientificNames() ([]string, error) {
	res := canonicalname.CanonicalNames{}
	for _, l := range Ω.list {
		res = res.AddToSet(l.CanonicalName())
		synonyms, err := l.Synonyms()
		if err != nil {
			return nil, err
		}
		res = res.AddToSet(synonyms...)
	}
	return res.ScientificNames(), nil
}

func (Ω *Aggregate) TargetIDs(sourceTypes ...datasources.SourceType) (datasources.TargetIDs, error) {
	res := datasources.TargetIDs{}
	for _, usage := range Ω.list {
		srcs, err := usage.Sources(sourceTypes...)
		if err != nil {
			return nil, err
		}
		for _, src := range srcs {
			targetID, err := src.TargetID()
			if err != nil {
				return nil, err
			}
			res = append(res, targetID)
		}
	}
	return res, nil
}

func (Ω *Aggregate) AddUsage(usages ...nameusage.NameUsage) error {
	Ω.Lock()
	defer Ω.Unlock()

	Ω.list = append(Ω.list, usages...)

ResetLoop:
	for {
		for i := range Ω.list {
			for k := range Ω.list {
				if k == i {
					continue
				}
				shouldCombine, err := Ω.list[i].ShouldCombine(Ω.list[k])
				if err != nil {
					return err
				}
				if shouldCombine {
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
		break
	}

	return nil
}

func (Ω *Aggregate) HasCanonicalName(name canonicalname.CanonicalName) (bool, error) {
	for i := range Ω.list {
		names, err := Ω.list[i].AllScientificNames()
		if err != nil {
			return false, err
		}
		if utils.ContainsString(names, name.ScientificName()) {
			return true, err
		}
	}
	return false, nil
}
