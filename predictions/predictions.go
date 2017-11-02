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
	"github.com/montanaflynn/stats"
	"os"
	"path"
	"github.com/elgs/gostrgen"
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
	PredictionsOverTaxon map[store.TaxonID]stats.Float64Data
	PredictionsOverDay map[string]stats.Float64Data
	sync.Mutex
	TempWriteDirectory string
}

func (Ω *PredictionUploader) GetWildernessArea(latitude, longitude float64) (*store.WildernessArea, error) {

	for _, w := range Ω.WildernessAreas {
		if utils.CoordinatesEqual(w.Centre.Latitude, latitude) && utils.CoordinatesEqual(w.Centre.Longitude, longitude) {
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

	random_string, err := gostrgen.RandGen(10, gostrgen.Lower, "", "")
	if err != nil {
		panic(err)
	}

	tmp := path.Join("/tmp/", fmt.Sprintf("predictions-%s", random_string))
	if err := os.Mkdir(tmp, os.ModePerm); err != nil {
		panic(err)
	}

	fmt.Println("TEMP_DIRECTORY", tmp)


	f := PredictionUploader{
		Limiter:                       make(chan struct{}, predictionUploadLimit),
		PredictionsOverDay: make(map[string]stats.Float64Data),
		PredictionsOverTaxon: make(map[store.TaxonID]stats.Float64Data),
		TempWriteDirectory: tmp,
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
		panic(err)
	}

	f.Bucket = client.Bucket("floracast-datamining")

	//if err := f.GatherAnalytics(cxt, *predictionDirectory); err != nil {
	//	panic(err)
	//}
	//
	//values := ByValue{}
	//for k, v := range f.PredictionsOverTaxon {
	//	m, _ := v.Median()
	//	values = append(values, ValueToSort{string(k), len(v), m})
	//}
	//sort.Sort(values)
	//for _, v := range values {
	//	fmt.Println(v.Key, v.Value, v.Decorator)
	//}

	if err := f.FetchUploadPredictions(cxt, *predictionDirectory); err != nil {
		panic(err)
	}

}

type ValueToSort struct {
	Key string
	Value int
	Decorator float64
}

type ByValue []ValueToSort

func (a ByValue) Len() int           { return len(a) }
func (a ByValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByValue) Less(i, j int) bool { return a[i].Value < a[j].Value }



func (Ω *PredictionUploader) GatherAnalytics(cxt context.Context, predictionsDirectory string) error {
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
				predictions, err := Ω.parseBucketObject(cxt, name)
				if err != nil {
					return err
				}
				Ω.Lock()
				defer Ω.Unlock()
				for _, p := range predictions {
					if _, ok := Ω.PredictionsOverTaxon[p.TaxonID]; !ok {
						Ω.PredictionsOverTaxon[p.TaxonID] = stats.Float64Data{}
					}
					Ω.PredictionsOverTaxon[p.TaxonID] = append(Ω.PredictionsOverTaxon[p.TaxonID], p.PredictionValue)
					if _, ok := Ω.PredictionsOverDay[p.FormattedDate]; !ok {
						Ω.PredictionsOverDay[p.FormattedDate] = stats.Float64Data{}
					}
					Ω.PredictionsOverDay[p.FormattedDate] = append(Ω.PredictionsOverDay[p.FormattedDate], p.PredictionValue)
				}
				return nil
			})
		}
		return nil
	})
	return tmb.Wait()
}

func (Ω *PredictionUploader) FetchUploadPredictions(cxt context.Context, predictionsDirectory string) error {

	if Ω.TempWriteDirectory == "" {
		return errors.New("temp write directory required")
	}

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
			fmt.Println("name", o.Name)
			tmb.Go(func() error {
				predictions, err := Ω.parseBucketObject(cxt, name)
				if err != nil {
					return err
				}
				for _, _p := range predictions  {
					p := _p
					//p.PercentileOverAllTaxaPredictionsForDay, err = Ω.PredictionsOverDay[p.FormattedDate].Percentile(p.PredictionValue)
					//if err != nil {
					//	return errors.Wrap(err, "could not calc percentile")
					//}
					//p.PercentileOverAllTaxonPredictions, err = Ω.PredictionsOverTaxon[p.TaxonID].Percentile(p.PredictionValue)
					//if err != nil {
					//	return errors.Wrap(err, "could not calc percentile")
					//}
					<- Ω.Limiter
					tmb.Go(func() error {
						defer func() {
							Ω.Limiter <- struct{}{}
						}()
						if err := Ω.writeTaxaFile(p); err != nil {
							return err
						}
						return Ω.writeTaxonFile(p)
					})
				}
				return nil
			})
		}
		return nil
	})
	return tmb.Wait()
}

func (Ω *PredictionUploader) writeTaxonFile(p store.Prediction) error {
	// First write predictions.
	taxonFile, err := os.OpenFile(path.Join(Ω.TempWriteDirectory, string(p.TaxonID)+".taxon"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer taxonFile.Close()

	line := fmt.Sprintf("%s,%s,%.6f",
		p.FormattedDate,
		p.WildernessAreaID,
		p.PredictionValue)

	line = strings.Replace(line, ".", "|", -1)

	line = fmt.Sprintf("{\"%s\": [%.6f, %.6f]}\n", line, p.Location.Latitude, p.Location.Longitude)

	if _, err := taxonFile.WriteString(line); err != nil {
		return errors.Wrap(err, "could not write taxon file")
	}
	return nil
}
func (Ω *PredictionUploader) writeTaxaFile(p store.Prediction) error {
	taxaFile, err := os.OpenFile(path.Join(Ω.TempWriteDirectory, "/species.txt"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer taxaFile.Close()
	// TODO: There is the potential tha two percentiles could be exactly the same in which case one would override the other, though this seems unlikely
	line := fmt.Sprintf("%s,%s,%.6f,%.6f,%.6f",
		p.TaxonID,
		p.FormattedDate,
		p.PredictionValue,
		p.PercentileOverAllTaxonPredictions,
		p.PercentileOverAllTaxaPredictionsForDay)
	line = strings.Replace(line, ".", "|", -1)
	line = fmt.Sprintf("{\"%s\": [%.6f, %.6f]}\n", line, p.Location.Latitude, p.Location.Longitude)

	if _, err := taxaFile.WriteString(line); err != nil {
		return errors.Wrap(err, "could not write taxa file")
	}
	return nil
}

type predictionCatcher struct{
	sync.Mutex
	Predictions []store.Prediction
}

func (Ω *PredictionUploader) parseBucketObject(cxt context.Context, name string) ([]store.Prediction, error) {

	r, err := Ω.Bucket.Object(name).NewReader(cxt)
	if err != nil {
		panic(err)
	}
	defer r.Close()
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "could no read predictions file")
	}
	catcher := predictionCatcher{}
	lines := []PredictionLine{}
	if err := json.Unmarshal(b, &lines); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal prediction file")
	}
	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, _line := range lines {
			line := _line
			tmb.Go(func() error {
				catcher.Lock()
				defer catcher.Unlock()
				predictions, err := Ω.parseLine(cxt, line)
				if err != nil {
					return err
				}
				catcher.Predictions = append(catcher.Predictions, predictions...)
				return nil
			})
		}
		return nil
	})
	if err := tmb.Wait(); err != nil {
		return nil, err
	}
	return catcher.Predictions, nil
}

func (Ω *PredictionUploader) parseLine(cxt context.Context, line PredictionLine) ([]store.Prediction, error)  {
	if len(line.Probabilities) == 0 {
		return nil, nil
	}
	parts := strings.Split(line.Key, "|")
	latitude, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse latitude")
	}
	longitude, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse longitude")
	}

	wa, err := Ω.GetWildernessArea(latitude, longitude)
	if err != nil {
		return nil, err
	}

	date, err := time.ParseInLocation("20060102", parts[2], time.UTC)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse date")
	}

	random_index := 0
	for i, c := range line.Classes {
		if c == "0" {
			random_index = i
		}
	}
	threshold := line.Probabilities[random_index]

	res := []store.Prediction{}
	for i := range line.Probabilities {
		if line.Classes[i] != "0" && line.Probabilities[i] > threshold {
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
			res = append(res, p)
		}
	}
	return res, nil
}