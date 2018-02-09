package aggregate

import (
	"sync"
	"encoding/json"
	"gopkg.in/tomb.v2"
	"context"
	"bitbucket.org/heindl/taxa/nameusage/canonicalname"
	"bitbucket.org/heindl/taxa/nameusage/nameusage"
	"bitbucket.org/heindl/taxa/datasources"
	"bitbucket.org/heindl/taxa/occurrences"
)

type Aggregate struct {
	list []*nameusage.NameUsage
	sync.Mutex
}

//type CanonicalNameUsages []*NameUsage

//func (a CanonicalNameUsages) Len() int           { return len(a) }
//func (a CanonicalNameUsages) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
//func (a CanonicalNameUsages) Less(i, j int) bool { return a[i].canonicalName.name < a[j].canonicalName.name }


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

func (Ω *Aggregate) ScientificNames() []string {
	res := canonicalname.CanonicalNames{}
	for _, l := range Ω.list {
		res = res.AddToSet(l.CanonicalName())
		res = res.AddToSet(l.Synonyms()...)
	}
	return res.ScientificNames()
}

func (Ω *Aggregate) TargetIDs(sourceTypes ...datasources.SourceType) (res datasources.DataSourceTargetIDs) {
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

type OccurrenceFetcher func(context.Context, datasources.SourceType, datasources.TargetID, *time.Time) (*occurrences.OccurrenceAggregation, error)

func (Ω *Aggregate) FetchOccurrences(ctx context.Context, oFetcher OccurrenceFetcher) error {

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, _src := range Ω.Sources() {
			src := _src
			tmb.Go(func() error {
				aggr, err := oFetcher(ctx, src.SourceType(), src.TargetID(), src.LastFetchedAt())
				if err != nil {
					return err
				}
				if err := src.RegisterOccurrenceFetch(aggr.Count()); err != nil {
					return err
				}
				if Ω.occurrenceAggregation == nil {
					Ω.occurrenceAggregation = aggr
					return nil
				}
				return Ω.occurrenceAggregation.Merge(aggr)
			})
		}
		return nil
	})
	return tmb.Wait()
}

func (Ω *NameUsage) UploadCachedOccurrences(cxt context.Context, firestoreClient *firestore.Client) error {
	if Ω.occurrenceAggregation == nil {
		return nil
	}
	return Ω.occurrenceAggregation.Upsert(cxt, firestoreClient)
}