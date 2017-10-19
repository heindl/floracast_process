package main

import (
	"time"
	"github.com/jonboulle/clockwork"
	"gopkg.in/tomb.v2"
	"bitbucket.org/heindl/taxa/store"
	"bitbucket.org/heindl/taxa/utils"
	"github.com/heindl/gbif"
	"github.com/saleswise/errors/errors"
	"strconv"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/genproto/googleapis/type/latlng"
	"os"
	"googlemaps.github.io/maps"
	"bitbucket.org/heindl/taxa/ecoregions"
)

func main() {

	ts, err := store.NewTaxaStore()
	if err != nil {
		panic(err)
	}

	fetcher, err := NewOccurrenceFetcher(ts, "/Users/m/Downloads/wwf_terr_ecos_oRn.json", clockwork.NewRealClock())
	if (err != nil) {
		panic(err)
	}

	if err := fetcher.FetchOccurrences(); err != nil {
		panic(err)
	}

	return

}

const (
	occurrenceFetchLimit = 1000
)

func NewOccurrenceFetcher(store store.TaxaStore, eco_region_file string, clock clockwork.Clock) (OccurrenceFetcher, error) {
	f := occurrenceFetcher{
		//Log:                           logrus.New(),
		TaxaStore:                     store,
		Limiter:                       make(chan struct{}, occurrenceFetchLimit),
		Clock:                         clock,
	}
	for i := 0; i < occurrenceFetchLimit; i++ {
		f.Limiter <- struct{}{}
	}

	var err error
	f.EcoRegionCache, err = ecoregions.NewEcoRegionCache(eco_region_file)
	if err != nil {
		return nil, err
	}

	return OccurrenceFetcher(&f), nil
}

type OccurrenceFetcher interface{
	FetchOccurrences() error
}

type occurrenceFetcher struct {
	TaxaStore       store.TaxaStore
	Clock           clockwork.Clock
	//NSQProducer     nsqeco.Producer
	// The species limiter sets the number of species whose occurrences will be updated concurrently.
	Limiter chan struct{}
	//geography.BoundaryFetcher
	Occurrences *store.Occurrences
	Schema *store.DataSources
	EcoRegionCache ecoregions.EcoRegionCache
}

func (Ω *occurrenceFetcher) FetchOccurrences() error {

	//modelReprocess := struct{
	//	sync.Mutex
	//	Names []string
	//}{}

	tmb := tomb.Tomb{}
	cxt := context.Background()
	tmb.Go(func() error {

		taxa, err := Ω.TaxaStore.ReadSpecies(cxt)
		if err != nil {
			return err
		}

		fmt.Printf("Gathering occurrences for %d taxa.\n", len(taxa))

		for _, _taxon := range taxa {
			taxon := _taxon
			tmb.Go(func() error {
				dataSources, err := Ω.TaxaStore.GetOccurrenceDataSources(cxt, taxon.ID)
				if err != nil {
					return err
				}
				for _, _dataSource := range dataSources {
					dataSource := _dataSource
					tmb.Go(func() error {

						occurrences, err := Ω.fetchSourceData(dataSource)
						if err != nil {
							return err
						}
						if err := setElevations(occurrences); err != nil {
							return err
						}

						for _, _o := range occurrences {
							o := _o
							// TODO: Note that there appears to be a problem
							// with concurrent transactions on the same TaxonID field.
							// Therefore ignore concurrency for that taxon for now on the eco region update,
							// which appears to repair it.
							//<-Ω.Limiter
							//tmb.Go(func() error {
							//	defer func() {
							//		Ω.Limiter <- struct{}{}
							//	}()
								isNewOccurrence, err := Ω.TaxaStore.UpsertOccurrence(cxt, o)
								if err != nil {
									return err
								}
								if isNewOccurrence {
									if err := Ω.TaxaStore.IncrementTaxonEcoRegion(cxt, o.TaxonID, o.EcoRegion); err != nil {
										return err
									}
								}
								//return nil
							//})
						}

						if err := Ω.TaxaStore.UpdateDataSourceLastFetched(cxt, dataSource); err != nil {
							return err
						}

						return nil
					})
				}
				return nil
			})
		}

		return nil


	})

	if err := tmb.Wait(); err != nil {
		return err
	}

	return nil

}

func (Ω *occurrenceFetcher) fetchSourceData(src store.DataSource) (store.Occurrences, error) {

	hasBeenFetchedInLastDay := src.LastFetchedAt != nil && !src.LastFetchedAt.IsZero() && src.LastFetchedAt.After(Ω.Clock.Now().Add(time.Hour * -24))
	if hasBeenFetchedInLastDay {
		return nil, nil
	}

	f, err := newfetcher(src.SourceID, src.TargetID)
	if err != nil {
		return nil, err
	}

	if f == nil {
		//Ω.Log.WithFields(M{logkeys.Source: query.Source}.Fields()).Error("could not find fetcher for data source")
		return nil, nil
	}

	occurrences, err := f.Fetch(src.LastFetchedAt, Ω.Clock.Now().In(time.UTC))
	if err != nil {
		//Ω.Log.WithFields(M{logkeys.Error: err}.Fields()).Warn("no occurrences found")
		return nil, err
	}

	res := store.Occurrences{}

	for _, o := range occurrences {
		if o.OccurrenceID == "" {
			continue
		}

		_, ecoRegionKey, err := Ω.EcoRegionCache.PointWithin(o.Location.GetLatitude(), o.Location.GetLongitude())
		if err != nil {
			return nil, err
		}
		// If no key is found, the point is likely not within the continental United States, or something is broken.
		if ecoRegionKey == "" {
			continue
		}
		o.EcoRegion = ecoRegionKey
		o.DataSourceID = src.SourceID
		o.TaxonID = src.TaxonID
		res = append(res, o)
	}

	return res, nil
}

func newfetcher(sourceID store.DataSourceID, targetID store.DataSourceTargetID) (Fetcher, error) {
	switch sourceID {
	case store.DataSourceIDGBIF:
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


type GBIF int

func (this GBIF) Fetch(begin *time.Time, end time.Time) (store.Occurrences, error) {

	if begin == nil || begin.IsZero() {
		begin = utils.TimePtr(time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC))
	}

	q := gbif.OccurrenceSearchQuery{
		TaxonKey:          int(this),
		LastInterpreted:    fmt.Sprintf("%s,%s", begin.Format("2006-01-02"), end.Format("2006-01-02")),
		HasCoordinate:      true,
		HasGeospatialIssue: false,
	}

	results, err := gbif.Occurrences(q)
	if err != nil {
		return nil, errors.Wrap(err, "could not request occurrences")
	}

	if len(results) == 0 {
		return nil, nil
	}

	res := store.Occurrences{}
	for _, r := range results {
		o := this.parse(r)
		if o == nil {
			continue
		}
		res = append(res, *o)
	}

	return res, nil
}

func (this GBIF) parse(o gbif.Occurrence) *store.Occurrence {
	// Note that the OccurrenceID, which I originally used, appears to be incomplete, duplicated, and missing in some cases.

	if o.GbifID == "" {
		return nil
	}
	if o.EventDate == nil || o.EventDate.Time.IsZero() {
		// TODO: Consider reporting malformed occurrence error.
		return nil
	}
	if o.DecimalLatitude == 0 || o.DecimalLongitude == 0 {
		return nil
	}

	return &store.Occurrence{
		TargetID:      o.GbifID,
		OccurrenceID:  o.OccurrenceID,
		Location:      latlng.LatLng{o.DecimalLatitude, o.DecimalLongitude},
		Date:          utils.TimePtr(o.EventDate.Time),
		FormattedDate: o.EventDate.Time.Format("20060102"),
		Month:         o.EventDate.Month(),
		RecordedBy:    o.RecordedBy,
		References:    o.References,
		CreatedAt:     utils.TimePtr(time.Now()),
	}
}

var EPSILON float64 = 0.00001

func coordinateEquals(a, b float64) bool {
	if ((a - b) < EPSILON && (b - a) < EPSILON) {
		return true
	}
	return false
}

func setElevations(occurrences store.Occurrences) error {

	if len(occurrences) == 0 {
		return nil
	}

	mc, err := maps.NewClient(maps.WithAPIKey(os.Getenv("FLORACAST_GOOGLE_MAPS_API_KEY")))
	if err != nil {
		return errors.Wrap(err, "could not get google maps client")
	}

	start := 0
	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for {
			end := start + 100
			if len(occurrences) <= end {
				end = len(occurrences)
			}
			_list := occurrences[start:end]
			tmb.Go(func() error {
				list := _list
				locations := make([]maps.LatLng, len(list))
				// Gather lat/lng pairs for elevation fetch.
				for i, o := range list {
					locations[i] = maps.LatLng{o.Location.GetLatitude(), o.Location.GetLongitude()}
				}
				res, err := mc.Elevation(context.Background(), &maps.ElevationRequest{Locations: locations})
				if err != nil {
					return errors.Wrap(err, "could not fetch elevations")
				}
			Occurrences:
				for i := range list {
					for _, r := range res {
						if !coordinateEquals(list[i].Location.GetLatitude(), r.Location.Lat) {
							continue
						}
						if !coordinateEquals(list[i].Location.GetLongitude(), r.Location.Lng) {
							continue
						}
						list[i].Elevation = r.Elevation
						continue Occurrences
					}
				}
				return nil
			})
			start = end
			if start >= len(occurrences) {
				break
			}
		}
		return nil
	})
	return tmb.Wait()
}

