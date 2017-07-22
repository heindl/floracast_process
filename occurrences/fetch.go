package occurrences

import (
	"github.com/omidnikta/logrus"
	"github.com/saleswise/errors/errors"
	"time"
	"github.com/jonboulle/clockwork"
	"gopkg.in/tomb.v2"
	"bitbucket.org/heindl/species/store"
	"strconv"
	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
	"sync"
)

func main() {

	ts, err := store.NewTaxaStore()
	if err != nil {
		panic(err)
	}

	fetcher := NewOccurrenceFetcher(ts, clockwork.NewRealClock())

	if err := fetcher.FetchOccurrences(); err != nil {
		panic(err)
	}

	return

}

const (
	speciesFetchLimit = 100
)

func NewOccurrenceFetcher(store store.TaxaStore, clock clockwork.Clock) OccurrenceFetcher {
	f := fetcher{
		//Log:                           logrus.New(),
		TaxaStore:                     store,
		Limiter:                       make(chan struct{}, speciesFetchLimit),
		Clock:                         clock,
	}
	for i := 0; i < speciesFetchLimit; i++ {
		f.Limiter <- struct{}{}
	}

	return &f
}

type OccurrenceFetcher interface{
	FetchOccurrences() error
}

type fetcher struct {
	Log             *logrus.Logger
	TaxaStore       store.TaxaStore
	Clock           clockwork.Clock
	//NSQProducer     nsqeco.Producer
	// The species limiter sets the number of species whose occurrences will be updated concurrently.
	Limiter chan struct{}
	//geography.BoundaryFetcher
	Occurrences *store.Occurrences
	Schema *store.Schema
}

func (Ω *fetcher) FetchOccurrences() error {

	//modelReprocess := struct{
	//	sync.Mutex
	//	Names []string
	//}{}

	tmb := tomb.Tomb{}
	tmb.Go(func() error {

		taxa, err := Ω.TaxaStore.ReadTaxaFromCanonicalNames(store.CanonicalName("Limenitis camilla"))
		if err != nil {
			return err
		}

		iter := Ω.TaxaStore.NewOccurrenceSchemeIterator(taxa[0].Key)
		i := 0
		for {

			s := store.Scheme{}
			k, err := iter.Next(&s)
			if i > 0 {
				break
			}
			i++
			if err != nil {
				if err != iterator.Done {
					return errors.Wrap(err, "could not fetch occurrence scheme")
				} else {
					break
				}
			}
			if err != nil && err != iterator.Done {
				return errors.Wrap(err, "could not fetch occurrence scheme")
			}
			if err != nil && err == iterator.Done {
				break
			}
			<-Ω.Limiter
			tmb.Go(func() error {
				defer func() {
					Ω.Limiter <- struct{}{}
				}()

				t, err := Ω.TaxaStore.GetTaxon(s.Key.Parent)
				if err != nil {
					return err
				}

				if _, err := Ω.fetchSource(k, t.CanonicalName, t.Rank, s.LastFetchedAt); err != nil {
					return err
				}
				s.LastFetchedAt = Ω.Clock.Now()
				if err := Ω.batchScheme(&s); err != nil {
					return err
				}

				//if !hasNew {
				//	return nil
				//}
				//modelReprocess.Lock()
				//defer modelReprocess.Unlock()
				//modelReprocess.Names = utils.AddStringToSet(modelReprocess.Names, string(spcs.CanonicalName))
				return nil
			})
		}
		return nil
	})

	if err := tmb.Wait(); err != nil {
		return err
	}

	if err := Ω.flushOccurrences(); err != nil {
		return err
	}
	if err := Ω.flushSchema(); err != nil {
		return err
	}

	//for _, m := range modelReprocess.Names {
	//	if err := Ω.NSQProducer.Publish(nsqeco.NSQClassifySpecies, []byte(m)); err != nil {
	//		return errors.Wrap(err, "could not schedule waypoint reprocess for taxon %s")
	//	}
	//}

	return nil

}

var occurrences struct {
	List store.Occurrences
	sync.Mutex
}
func (Ω *fetcher) batchOccurrence(o *store.Occurrence) error {
	occurrences.Lock()
	defer occurrences.Unlock()
	occurrences.List = append(occurrences.List, o)
	if len(occurrences.List) >= 1000 {
		return Ω.flushOccurrences()
	}
	return nil
}
func (Ω *fetcher) flushOccurrences() error {
	if err := Ω.TaxaStore.SetOccurrences(occurrences.List); err != nil {
		return err
	}
	occurrences.List = occurrences.List[:0]
	return nil
}

var schema struct {
	List store.Schema
	sync.Mutex
}
func (Ω *fetcher) batchScheme(s *store.Scheme) error {
	schema.Lock()
	defer schema.Unlock()
	schema.List = append(schema.List, s)
	if len(schema.List) >= 1000 {
		return Ω.flushSchema()
	}
	return nil
}
func (Ω *fetcher) flushSchema() error {
	if err := Ω.TaxaStore.UpdateSchemaLastFetched(schema.List); err != nil {
		return err
	}
	schema.List = schema.List[:0]
	return nil
}

func (Ω *fetcher) fetchSource(schemeKey *datastore.Key, name store.CanonicalName, rank store.TaxonRank, lastFetchedAt time.Time) (bool, error) {
	hasNew := false
	hasBeenFetchedInLastDay := !lastFetchedAt.IsZero() && lastFetchedAt.After(Ω.Clock.Now().Add(time.Hour * -24))
	if hasBeenFetchedInLastDay {
		return false, nil
	}

	schemeSourceID, schemeTarget := store.SchemeKeyName(schemeKey.Name).Parse()
	fetcher, err := newfetcher(schemeSourceID, schemeTarget)
	if err != nil {
		return false, err
	}

	if fetcher == nil {
		//Ω.Log.WithFields(M{logkeys.Source: query.Source}.Fields()).Error("could not find fetcher for data source")
		return false, nil
	}

	occurrences, err := fetcher.Fetch(&lastFetchedAt, Ω.Clock.Now().In(time.UTC))
	if err != nil {
		//Ω.Log.WithFields(M{logkeys.Error: err}.Fields()).Warn("no occurrences found")
		return false, err
	}

	for i := range occurrences {
		if occurrences[i] == nil {
			continue
		}

		occurrences[i].Key.Parent = schemeKey
		if err := Ω.batchOccurrence(occurrences[i]); err != nil {
			return false, errors.Wrap(err, "could not run occurrence in batch")
		}
		hasNew = true
	}

	return hasNew, nil
}

func newfetcher(sourceID store.SchemeSourceID, targetID store.SchemeTargetID) (Fetcher, error) {
	switch sourceID {
	case store.SchemeSourceIDGBIF:
		i, err := strconv.Atoi(string(targetID))
		if err != nil {
			return nil, errors.Wrap(err, "could not cast targetID as string")
		}
		return Fetcher(GBIF(i)), nil
	//case species.SourceTypeEOL:
		// EOL doesn't store occurrences.
		//return nil, nil
	default:
		return nil, errors.New("unsupported source type")
	}
}

type Fetcher interface {
	Fetch(begin *time.Time, end time.Time) (store.Occurrences, error)
}
