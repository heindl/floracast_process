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

	f.Bucket = client.Bucket("floracast-models")

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
			tmb.Go(func() error {
				name := o.Name
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

						date, err := time.ParseInLocation("20060102", parts[2], time.UTC)
						if err != nil {
							return errors.Wrap(err, "could not parse date")
						}
						threshold := line.Probabilities[0]
						probabilities := line.Probabilities[1:]
						classes := line.Classes[1:]
						for i := range probabilities {
							if probabilities[i] > threshold {
								<- Ω.Limiter
								p := store.Prediction{
									CreatedAt: utils.TimePtr(time.Now()),
									Location: latlng.LatLng{latitude, longitude},
									PredictionValue: probabilities[i],
									TaxonID: store.TaxonID(classes[i]),
									Date: utils.TimePtr(date),
									FormattedDate: date.Format("20060102"),
									Month: date.Month(),
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
