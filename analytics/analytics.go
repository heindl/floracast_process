package main

import (
	"github.com/saleswise/errors/errors"
	"context"
	"cloud.google.com/go/storage"
	"fmt"
	"google.golang.org/api/iterator"
	"gopkg.in/tomb.v2"
	"io/ioutil"
	"github.com/mongodb/mongo-tools/common/json"
	"bitbucket.org/heindl/taxa/store"
	"flag"
	"sync"
	"bitbucket.org/heindl/taxa/utils"
)

type PredictionLine struct {
	Probabilities []float64 `json:"probabilities"`
	Classes       []string  `json:"classes"`
	Key           string    `json:"key"`
}

type Analyzer struct {
	Limiter chan struct{}
	Store store.TaxaStore
	Bucket *storage.BucketHandle
	Catcher *Catcher
}

type Catcher struct {
	sync.Mutex
	Predictions map[store.TaxonID]Values `json:",omitempty"`
}

type Values struct{
	Total int `json:""`
	AboveThreshold int `json:""`
	BelowThreshold int `json:""`
}

func (Ω *Catcher) Catch(taxonID store.TaxonID, above, below int) {
	Ω.Lock()
	defer Ω.Unlock()
	if _, ok := Ω.Predictions[taxonID]; !ok {
		Ω.Predictions[taxonID] = Values{}
	}
	v := Ω.Predictions[taxonID]
	v.Total += (above + below)
	v.AboveThreshold += above
	v.BelowThreshold += below
	Ω.Predictions[taxonID] = v
}

const analyticsProcessLimit = 2000

func main() {

	var err error
	occurrences := flag.Bool("occurrences", false, "count predictions for each taxon?")
	predictionDirectory := flag.String("dir", "", "prediction file directory under gs:floracast-models/predictions/")
	flag.Parse()

	f := Analyzer{
		Limiter:                       make(chan struct{}, analyticsProcessLimit),
		Catcher: &Catcher{Predictions: make(map[store.TaxonID]Values)},
	}

	for i := 0; i < analyticsProcessLimit; i++ {
		f.Limiter <- struct{}{}
	}

	f.Store, err = store.NewTaxaStore()
	if err != nil {
		panic(err)
	}

	cxt := context.Background()

	client, err := storage.NewClient(cxt)
	if err != nil {
		// TODO: Handle error.
	}

	f.Bucket = client.Bucket("floracast-datamining")
	if *occurrences {
		if err := f.CountOccurrences(cxt); err != nil {
			panic(err)
		}
	} else {
		if err := f.FetchPredictionAnalysis(cxt, *predictionDirectory); err != nil {
			panic(err)
		}

		fmt.Println(utils.JsonOrSpew(f.Catcher.Predictions))
	}

}

func (Ω *Analyzer) CountOccurrences(cxt context.Context) error {
	taxa, err := Ω.Store.ReadTaxa(cxt)
	if err != nil {
		return err
	}
	res := map[store.TaxonID]int{}

	for _, t := range taxa {

		occurrences, err := Ω.Store.GetOccurrences(cxt, t.ID)
		if err != nil {
			return err
		}
		res[t.ID] = len(occurrences)
	}

	fmt.Println(utils.JsonOrSpew(res))

	return nil
}


func (Ω *Analyzer) FetchPredictionAnalysis(cxt context.Context, predictionsDirectory string) error {
	q := &storage.Query{Prefix: fmt.Sprintf("predictions/%s/", predictionsDirectory), Delimiter: "/"}

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		iter := Ω.Bucket.Objects(cxt, q)
		for {
			o, err := iter.Next()
			if err != nil && err == iterator.Done {
				break
			}
			if err != nil {
				panic(err)
			}
			name := o.Name
			tmb.Go(func() error {
				//fmt.Println("name", name)
				r, err := Ω.Bucket.Object(name).NewReader(cxt)
				if err != nil {
					panic(err)
				}
				defer r.Close()
				b, err := ioutil.ReadAll(r)
				if err != nil {
					return errors.Wrap(err, "could no read predictions file")
				}
				lines := []PredictionLine{}
				if err := json.Unmarshal(b, &lines); err != nil {
					return errors.Wrap(err, "could not unmarshal prediction file")
				}
				for _, _line := range lines {
					line := _line
					tmb.Go(func() error {
						if len(line.Probabilities) == 0 {
							return nil
						}
						//parts := strings.Split(line.Key, "|")
						//latitude, err := strconv.ParseFloat(parts[0], 64)
						//if err != nil {
						//	return errors.Wrap(err, "could not parse latitude")
						//}
						//longitude, err := strconv.ParseFloat(parts[1], 64)
						//if err != nil {
						//	return errors.Wrap(err, "could not parse longitude")
						//}
						//
						//date, err := time.ParseInLocation("20060102", parts[2], time.UTC)
						//if err != nil {
						//	return errors.Wrap(err, "could not parse date")
						//}

						random_index := 0
						found := false
						for i, c := range line.Classes {
							if c == "0" {
								found = true
								random_index = i
							}
						}
						if !found {
							return errors.New("could not find random index value")
						}
						threshold := line.Probabilities[random_index]

						for i := range line.Probabilities {
							if line.Classes[i] == "0" {
								continue
							}
							if line.Probabilities[i] >= threshold {
								Ω.Catcher.Catch(store.TaxonID(line.Classes[i]), 1, 0)
								//<- Ω.Limiter
								//p := store.Prediction{
								//	CreatedAt: utils.TimePtr(time.Now()),
								//	Location: latlng.LatLng{latitude, longitude},
								//	PredictionValue: line.Probabilities[i],
								//	TaxonID: store.TaxonID(line.Classes[i]),
								//	Date: utils.TimePtr(date),
								//	FormattedDate: date.Format("20060102"),
								//	Month: date.Month(),
								//}
								//tmb.Go(func() error {
								//	defer func() {
								//		Ω.Limiter <- struct{}{}
								//	}()
								//	return Ω.Store.SetPrediction(cxt, p)
								//})
							} else {
								Ω.Catcher.Catch(store.TaxonID(line.Classes[i]), 0, 1)
							}
						}
						return nil
					})
				}
				return nil
			})
		}
		return nil
	})
	return tmb.Wait()
}
