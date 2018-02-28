package main

//
//import (
//	"bitbucket.org/heindl/process/ecoregions"
//	"bitbucket.org/heindl/process/store"
//	"bitbucket.org/heindl/process/utils"
//	"flag"
//	"fmt"
//	"bitbucket.org/heindl/process/datasources/gbif"
//	"github.com/jonboulle/clockwork"
//	"golang.org/x/net/context"
//	"google.golang.org/genproto/googleapis/type/latlng"
//	"gopkg.in/tomb.v2"
//	"strconv"
//	"strings"
//	"time"
//)
//
//func main() {
//
//	taxa := flag.String("taxa", "", "view taxon information")
//	flag.Parse()
//
//	ts, err := store.NewTaxaStore()
//	if err != nil {
//		panic(err)
//	}
//
//	fetcher, err := NewOccurrenceFetcher(ts, "/Users/m/Downloads/wwf_terr_ecos_oRn.json", clockwork.NewRealClock())
//	if err != nil {
//		panic(err)
//	}
//
//	if err := fetcher.FetchOccurrences(strings.Split(*taxa, ",")); err != nil {
//		panic(err)
//	}
//
//	return
//
//}
//
//const (
//	occurrenceFetchLimit = 100
//)
//
//func NewOccurrenceFetcher(store store.TaxaStore, eco_region_file string, clock clockwork.Clock) (OccurrenceFetcher, error) {
//	f := occurrenceFetcher{
//		//Log:                           logrus.New(),
//		TaxaStore: store,
//		Limiter:   make(chan struct{}, occurrenceFetchLimit),
//		Clock:     clock,
//	}
//	for i := 0; i < occurrenceFetchLimit; i++ {
//		f.Limiter <- struct{}{}
//	}
//
//	var err error
//	f.EcoRegionCache, err = ecoregions.NewEcoRegionCache(eco_region_file)
//	if err != nil {
//		return nil, err
//	}
//
//	return OccurrenceFetcher(&f), nil
//}
//
//type OccurrenceFetcher interface {
//	FetchOccurrences(taxa []string) error
//}
//
//type occurrenceFetcher struct {
//	TaxaStore store.TaxaStore
//	Clock     clockwork.Clock
//	//NSQProducer     nsqeco.Producer
//	// The species limiter sets the number of species whose occurrences will be updated concurrently.
//	Limiter chan struct{}
//	//geography.BoundaryFetcher
//	Occurrences    *datasou.Occurrences
//	Schema         *store.DataSources
//	EcoRegionCache ecoregions.EcoRegionCache
//}
//
//func (Ω *occurrenceFetcher) FetchOccurrences(taxa_ids []string) error {
//
//	//modelReprocess := struct{
//	//	sync.Mutex
//	//	Names []string
//	//}{}
//
//	tmb := tomb.Tomb{}
//	cxt := context.Background()
//	tmb.Go(func() error {
//
//		taxa := store.Taxa{}
//		var err error
//		if len(taxa_ids) > 0 {
//			for _, taxon_id := range taxa_ids {
//				if !store.INaturalistTaxonID(taxon_id).Valid() {
//					return errors.Newf("invalid taxon id[%s]", taxon_id)
//				}
//				taxon, err := Ω.TaxaStore.ReadTaxon(cxt, store.INaturalistTaxonID(taxon_id))
//				if err != nil {
//					return err
//				}
//				taxa = append(taxa, *taxon)
//			}
//		} else {
//			taxa, err = Ω.TaxaStore.ReadSpecies(cxt)
//			if err != nil {
//				return err
//			}
//		}
//
//		fmt.Printf("Gathering occurrences for %d taxa.\n", len(taxa))
//
//		for _, _taxon := range taxa {
//			taxon := _taxon
//			tmb.Go(func() error {
//				dataSources, err := Ω.TaxaStore.GetOccurrenceDataSources(cxt, taxon.ID)
//				if err != nil {
//					return err
//				}
//				for _, _dataSource := range dataSources {
//					dataSource := _dataSource
//					tmb.Go(func() error {
//						<-Ω.Limiter
//						defer func() {
//							Ω.Limiter <- struct{}{}
//						}()
//
//						occurrences, err := Ω.fetchSourceData(dataSource)
//						if err != nil {
//							return err
//						}
//
//						//if err := setElevations(occurrences); err != nil {
//						//	return err
//						//}
//
//						for _, _o := range occurrences {
//							o := _o
//
//							// TODO: Note that there appears to be a problem
//							// with concurrent transactions on the same INaturalistTaxonID field.
//							// Therefore ignore concurrency for that taxon for now on the eco region update,
//							// which appears to repair it.
//							//<-Ω.Limiter
//							//tmb.Go(func() error {
//							//	defer func() {
//							//		Ω.Limiter <- struct{}{}
//							//	}()
//							isNewOccurrence, err := Ω.TaxaStore.UpsertOccurrence(cxt, o)
//							if err != nil {
//								return err
//							}
//							if isNewOccurrence {
//								if err := Ω.TaxaStore.IncrementTaxonEcoRegion(cxt, o.TaxonID, o.EcoRegion); err != nil {
//									return err
//								}
//							}
//							return nil
//							//})
//						}
//
//						if err := Ω.TaxaStore.UpdateDataSourceLastFetched(cxt, dataSource); err != nil {
//							return err
//						}
//
//						return nil
//					})
//				}
//				return nil
//			})
//		}
//
//		return nil
//
//	})
//
//	if err := tmb.Wait(); err != nil {
//		return err
//	}
//
//	return nil
//
//}
//
//func (Ω *occurrenceFetcher) fetchSourceData(src store.DataSource) (store.Occurrences, error) {
//
//	hasBeenFetchedInLastDay := src.LastFetchedAt != nil && !src.LastFetchedAt.IsZero() && src.LastFetchedAt.After(Ω.Clock.Now().AddOccurrence(time.Hour*-24))
//	if hasBeenFetchedInLastDay {
//		return nil, nil
//	}
//
//	f, err := newfetcher(src.SourceID, src.TargetID)
//	if err != nil {
//		return nil, err
//	}
//
//	if f == nil {
//		//Ω.Log.WithFields(M{logkeys.Source: query.Source}.Fields()).Error("could not find fetcher for data source")
//		return nil, nil
//	}
//
//	occurrences, err := f.Fetch(src.LastFetchedAt, Ω.Clock.Now().In(time.UTC))
//	if err != nil {
//		//Ω.Log.WithFields(M{logkeys.Error: err}.Fields()).Warn("no occurrences found")
//		return nil, err
//	}
//
//	res := store.Occurrences{}
//
//	for _, o := range occurrences {
//		if o.OccurrenceID == "" {
//			continue
//		}
//
//		_, ecoRegionKey, err := Ω.EcoRegionCache.PointWithin(o.Location.GetLatitude(), o.Location.GetLongitude())
//		if err != nil {
//			return nil, err
//		}
//		// If no key is found, the point is likely not within the continental United States, or something is broken.
//		// TODO: About three percent of a test batch fell outside of all eco-regions. They appeared to be on costal areas for which the defined regions may not extend.
//		if ecoRegionKey == "" {
//			//fmt.Println("no ecoregion key", o.Location.GetLatitude(), o.Location.GetLongitude())
//			continue
//		}
//
//		o.EcoRegion = ecoRegionKey
//		o.DataSourceID = src.SourceID
//		o.TaxonID = src.TaxonID
//		res = append(res, o)
//	}
//
//	return res, nil
//}
//
//func newfetcher(sourceID store.SourceType, targetID store.DataSourceTargetID) (Fetcher, error) {
//	switch sourceID {
//	case store.DataSourceTypeGBIF:
//		i, err := strconv.Atoi(string(targetID))
//		if err != nil {
//			return nil, errors.Wrap(err, "could not cast targetID as string")
//		}
//		return Fetcher(GBIF(i)), nil
//	//case species.SourceTypeEOL:
//	// EOL doesn't store occurrences.
//	//return nil, nil
//	default:
//		return nil, errors.New("unsupported source type")
//	}
//}
//
//type Fetcher interface {
//	Fetch(begin *time.Time, end time.Time) (store.Occurrences, error)
//}
//
//type GBIF int
//
//func (this GBIF) Fetch(begin *time.Time, end time.Time) (store.Occurrences, error) {
//
//	if begin == nil || begin.IsZero() {
//		begin = utils.TimePtr(time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC))
//	}
//
//	q := gbif.OccurrenceSearchQuery{
//		TaxonKey:           int(this),
//		LastInterpreted:    fmt.Sprintf("%s,%s", begin.Format("2006-01-02"), end.Format("2006-01-02")),
//		HasCoordinate:      true,
//		HasGeospatialIssue: false,
//	}
//
//	results, err := gbif.Occurrences(q)
//	if err != nil {
//		return nil, errors.Wrap(err, "could not request occurrences")
//	}
//
//	if len(results) == 0 {
//		return nil, nil
//	}
//
//	res := store.Occurrences{}
//	for _, r := range results {
//		o := this.parse(r)
//		if o == nil {
//			continue
//		}
//		res = append(res, *o)
//	}
//
//	return res, nil
//}
//
//func (this GBIF) parse(o gbif.Occurrence) *store.Occurrence {
//	// Note that the OccurrenceID, which I originally used, appears to be incomplete, duplicated, and missing in some cases.
//
//	if o.GbifID == "" {
//		return nil
//	}
//	if o.EventDate == nil || o.EventDate.Time.IsZero() {
//		// TODO: Consider reporting malformed occurrence error.
//		return nil
//	}
//	if o.DecimalLatitude == 0 || o.DecimalLongitude == 0 {
//		return nil
//	}
//
//	return &store.Occurrence{
//		TargetID:      o.GbifID,
//		OccurrenceID:  o.OccurrenceID,
//		Location:      latlng.LatLng{o.DecimalLatitude, o.DecimalLongitude},
//		Date:          utils.TimePtr(o.EventDate.Time),
//		FormattedDate: o.EventDate.Time.Format("20060102"),
//		Month:         o.EventDate.Month(),
//		RecordedBy:    o.RecordedBy,
//		References:    o.References,
//		CreatedAt:     utils.TimePtr(time.Now()),
//		CountryCode:   o.CountryCode,
//	}
//}
