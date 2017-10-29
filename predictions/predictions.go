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
	"strings"
	"bitbucket.org/heindl/taxa/store"
	"strconv"
	"flag"
	"google.golang.org/genproto/googleapis/type/latlng"
	"time"
	"bitbucket.org/heindl/taxa/utils"
	"sync"
)

type PredictionLine struct {
	Probabilities []float64 `json:"probabilities"`
	Classes       []string  `json:"classes"`
	Key           string    `json:"key"`
}

type PredictionUploader struct {
	Limiter chan struct{}
	Store store.TaxaStore
	Bucket *storage.BucketHandle
	WildernessAreas store.WildernessAreas
	sync.Mutex
}

func (Ω *PredictionUploader) GetWildernessArea(latitude, longitude float64) (*store.WildernessArea, error) {

	for _, w := range Ω.WildernessAreas {
		if w.Centre.Latitude == latitude && w.Centre.Longitude == longitude {
			return &w, nil
		}
	}

	w, err := Ω.Store.ReadWildernessArea(context.Background(), latitude, longitude)
	if err != nil {
		return nil, err
	}

	Ω.Lock()
	defer Ω.Unlock()

	Ω.WildernessAreas = append(Ω.WildernessAreas, *w)

	return w, nil
}

const predictionUploadLimit = 2000

func main() {

	var err error
	predictionDirectory := flag.String("dir", "", "prediction file directory under gs:floracast-models/predictions/")
	flag.Parse()

	f := PredictionUploader{
		Limiter:                       make(chan struct{}, predictionUploadLimit),
	}

	for i := 0; i < predictionUploadLimit; i++ {
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

	if err := f.FetchUploadPredictions(cxt, *predictionDirectory); err != nil {
		panic(err)
	}

}

func (Ω *PredictionUploader) FetchUploadPredictions(cxt context.Context, predictionsDirectory string) error {
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
				fmt.Println("name", name)
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
						parts := strings.Split(line.Key, "|")
						latitude, err := strconv.ParseFloat(parts[0], 64)
						if err != nil {
							return errors.Wrap(err, "could not parse latitude")
						}
						longitude, err := strconv.ParseFloat(parts[1], 64)
						if err != nil {
							return errors.Wrap(err, "could not parse longitude")
						}

						wa, err := Ω.GetWildernessArea(latitude, longitude)
						if err != nil {
							return err
						}

						date, err := time.ParseInLocation("20060102", parts[2], time.UTC)
						if err != nil {
							return errors.Wrap(err, "could not parse date")
						}

						random_index := 0
						for i, c := range line.Classes {
							if c == "0" {
								random_index = i
							}
						}
						threshold := line.Probabilities[random_index]

						for i := range line.Probabilities {
							if line.Classes[i] != "0" && line.Probabilities[i] > threshold {
								<- Ω.Limiter
								p := store.Prediction{
									CreatedAt: utils.TimePtr(time.Now()),
									Location: latlng.LatLng{latitude, longitude},
									PredictionValue: line.Probabilities[i],
									TaxonID: store.TaxonID(line.Classes[i]),
									Date: utils.TimePtr(date),
									FormattedDate: date.Format("20060102"),
									Month: date.Month(),
									WildernessAreaID: wa.ID,
									WildernessAreaName: wa.Name,
								}
								tmb.Go(func() error {
									defer func() {
										Ω.Limiter <- struct{}{}
									}()
									return Ω.Store.SetPrediction(cxt, p)
								})
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
