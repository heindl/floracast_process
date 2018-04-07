package algolia

import (
	"bitbucket.org/heindl/process/occurrence"
	"bitbucket.org/heindl/process/predictions"
	"bitbucket.org/heindl/process/store"
	"context"
	"fmt"
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
)

const OccurrenceIndex = store.AlgoliaIndexName("Occurrences")
const PredictionIndex = store.AlgoliaIndexName("Predictions")

var occurrenceIndexSettings = algoliasearch.Map{
	"customRanking": []string{
		fmt.Sprintf("desc(%s)", "Occurrences"),
	},
	"searchableAttributes": []string{
		"_geoloc",
	},
}

var predictionIndexSettings = algoliasearch.Map{
	"customRanking": []string{
		fmt.Sprintf("desc(%s)", "Predictions"),
	},
	"searchableAttributes": []string{
		"_geoloc",
	},
}

type GeolocatedIndexRecord struct {
	StandardNameUsageRecord
	Geolocations []*Coords `json:"_geoloc,omitempty"`
}

func IndexOccurrences() {

}

type CoordinateProvider interface {
	Coordinates() (lat float64, lng float64, err error)
}

func ParseCoords(coords ...CoordinateProvider) ([]*Coords, error) {

	res := []*Coords{}

CoordIterator:
	for _, o := range coords {
		lat, lng, err := o.Coordinates()
		if err != nil {
			return nil, err
		}
		c := &Coords{
			Lat:        lat,
			Lng:        lng,
			comparator: fmt.Sprintf("%.4f,%.4f", lat, lng),
		}
		for _, i := range res {
			if i.comparator == c.comparator {
				continue CoordIterator
			}
		}
		res = append(res, c)
	}

	return res, nil
}

type Coords struct {
	Lat        float64 `json:"lat"`
	Lng        float64 `json:"lng"`
	comparator string
}

func (Ω *NameUsageIndexRecord) fetchOccurrenceGeolocations() error {
	occurrences, err := occurrence.FetchFromFireStore(context.Background(), Ω.floraStore, *Ω.ObjectID)
	if err != nil {
		return err
	}
	Ω.Occurrences = len(occurrences)
	providers := []CoordinateProvider{}
	for _, o := range occurrences {
		providers = append(providers, o)
	}
	Ω.occurrenceLocations, err = ParseCoords(providers...)
	if err != nil {
		return err
	}
	return nil
}

func (Ω *NameUsageIndexRecord) fetchPredictionGeolocations() error {
	predictionList, err := predictions.FetchFromFireStore(context.Background(), Ω.floraStore, *Ω.ObjectID)
	if err != nil {
		return err
	}
	Ω.Predictions = len(predictionList)
	providers := []CoordinateProvider{}
	for _, p := range predictionList {
		providers = append(providers, p)
	}
	Ω.predictionLocations, err = ParseCoords(providers...)
	if err != nil {
		return err
	}
	return nil
}
