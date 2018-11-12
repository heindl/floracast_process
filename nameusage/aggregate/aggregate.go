package aggregate

import (
	"github.com/heindl/floracast_process/algolia"
	"github.com/heindl/floracast_process/datasources"
	"github.com/heindl/floracast_process/nameusage/canonicalname"
	"github.com/heindl/floracast_process/nameusage/nameusage"
	"github.com/heindl/floracast_process/store"
	"github.com/heindl/floracast_process/taxa"
	"context"
	"sort"
	"sync"
)

// Aggregate is a structure for grouping and combining NameUsages.
type Aggregate struct {
	list []nameusage.NameUsage
	sync.Mutex
}

// EachFunction is a callback for iterating over an aggregation.
type EachFunction func(ctx context.Context, usage nameusage.NameUsage) error

// Each is a helper for iterating over an aggregation.
func (Ω *Aggregate) Each(ctx context.Context, cb EachFunction) error {
	//tmb := tomb.Tomb{}
	//tmb.Go(func() error {
	for 𝝨 := range Ω.list {
		_i := 𝝨
		//tmb.Go(func() error {
		i := _i
		return cb(ctx, Ω.list[i])
		//})
	}
	return nil
	//})
	//return tmb.Wait()
}

// FilterFunction is a callback for iterating over aggregation.
type FilterFunction func(usage nameusage.NameUsage) (bool, error)

// Filter is a helper for filtering an aggregation.
func (Ω *Aggregate) Filter(shouldFilter FilterFunction) (*Aggregate, error) {
	res := Aggregate{}
	for _, 𝝨 := range Ω.list {
		u := 𝝨
		should, err := shouldFilter(u)
		if err != nil {
			return nil, err
		}
		if should {
			continue
		}
		if err := res.Append(u); err != nil {
			return nil, err
		}
	}
	return &res, nil
}

// Upload saves a NameUsage to FireStore, deletes old records, creates/uploads
// Algolia objects, and materializes/uploads a Taxon object.
func (Ω *Aggregate) Upload(cxt context.Context, florastore store.FloraStore) error {

	return Ω.Each(cxt, func(ctx context.Context, usage nameusage.NameUsage) error {
		nameUsageID, err := usage.ID()
		if err != nil {
			return err
		}
		_, err = usage.Upload(ctx, florastore)
		if err != nil {
			return err
		}

		// TODO: Ignore deleted ids for now as the matching logic requires updates.
		if err := taxa.UploadMaterializedTaxon(ctx, florastore, usage); err != nil {
			return err
		}
		// TODO: Ignore deleted ids for now as the matching logic requires updates.
		//for _, idToDelete := range deletedUsageIDs {
		//	if err := algolia.DeleteNameUsage(ctx, florastore, idToDelete); err != nil {
		//		return err
		//	}
		//}
		return algolia.IndexNameUsage(ctx, florastore, nameUsageID)
	})

}

// Count returns the number of NameUsage objects in the aggregate.
func (Ω *Aggregate) Count() int {
	return len(Ω.list)
}

//func (Ω *Aggregate) occurrences() (int, error) {
//	res := 0
//	for _, l := range Ω.list {
//		i, err := l.Occurrences()
//		if err != nil {
//			return 0, err
//		}
//		res += i
//	}
//	return res, nil
//}

// ScientificNames returns a list of all names in all sources in all NameUsages.
func (Ω *Aggregate) ScientificNames() ([]string, error) {
	res := canonicalname.Names{}
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

// TargetIDs returns a list of all TargetIDs in given sources from all NameUsages.
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

// Append adds a new usage to the list but does not reduce.
func (Ω *Aggregate) Append(usages ...nameusage.NameUsage) error {
	Ω.list = append(Ω.list, usages...)
	return nil
}

// Reduce combines all name usages that share scientific names.
func (Ω *Aggregate) Reduce() error {

	Ω.Lock()
	defer Ω.Unlock()

ResetLoop:
	for {

		sort.Sort(nameusage.ByCanonicalName(Ω.list))
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

// Add adds a new NameUsage to an aggregate and reduces.
func (Ω *Aggregate) Combine(usages ...nameusage.NameUsage) error {
	if err := Ω.Append(usages...); err != nil {
		return err
	}
	return Ω.Reduce()
}

//func (Ω *Aggregate) hasCanonicalName(name canonicalname.Name) (bool, error) {
//	for i := range Ω.list {
//		names, err := Ω.list[i].AllScientificNames()
//		if err != nil {
//			return false, err
//		}
//		if utils.ContainsString(names, name.ScientificName()) {
//			return true, err
//		}
//	}
//	return false, nil
//}
