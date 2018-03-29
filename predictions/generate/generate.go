package generate

import (
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/predictions"
	"bitbucket.org/heindl/process/store"
	"bitbucket.org/heindl/process/tfrecords"
	"bitbucket.org/heindl/process/utils"
	"fmt"
	tg "github.com/galeone/tfgo"
	"github.com/saleswise/errors/errors"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"gopkg.in/tomb.v2"
	"path"
	"sort"
	"strings"
	"sync"
)

func GeneratePredictions(ctx context.Context, nameUsageID nameusage.ID, floraStore store.FloraStore, dateRange *DateRange) (list predictions.Predictions, err error) {

	if !nameUsageID.Valid() {
		return nil, errors.Newf("Invalid NameUsageID [%s]", nameUsageID)
	}

	modeller, err := NewModeller(floraStore)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not get model for NameUsage [%s]", nameUsageID)
	}
	defer utils.SafeClose(modeller, &err)

	model, err := modeller.FetchModel(ctx, nameUsageID)
	if err != nil {
		return nil, err
	}

	g := generator{
		nameUsageID: nameUsageID,
		model:       model,
		dateRange:   dateRange,
		floraStore:  floraStore,
	}

	filenames, err := g.fetchLatestProtectedAreaFileNames(ctx)
	if err != nil {
		return nil, err
	}

	res := predictions.Predictions{}
	tmb := tomb.Tomb{}
	lock := sync.Mutex{}
	tmb.Go(func() error {
		for _, _f := range filenames {
			f := _f
			tmb.Go(func() error {
				list, err := g.generatePredictionsFromGCSPath(ctx, f)
				if err != nil {
					return err
				}
				lock.Lock()
				defer lock.Unlock()
				res = append(res, list...)
				return nil
			})
		}
		return nil
	})
	if err := tmb.Wait(); err != nil {
		return nil, err
	}

	fmt.Println("TOTAL READ", g.totalRead)

	return res, nil

}

type DateRange struct {
	Start, End string
}

type generator struct {
	nameUsageID nameusage.ID
	model       *tg.Model
	dateRange   *DateRange
	floraStore  store.FloraStore
	totalRead   int
}

func (Ω *generator) fetchLatestProtectedAreaFileNames(ctx context.Context) ([]string, error) {

	prefix := fmt.Sprintf("protected_areas/")
	names, err := Ω.floraStore.CloudStorageObjectNames(ctx, prefix, ".tfrecords")
	if err != nil {
		return nil, err
	}

	dateFileMap := map[string][]string{}
	for _, n := range names {
		s := strings.Split(n, "/")
		d := s[0]
		f := s[1]
		if _, ok := dateFileMap[d]; !ok {
			dateFileMap[d] = []string{}
		}
		dateFileMap[d] = append(dateFileMap[d], f)
	}

	res := []string{}
	for date, files := range dateFileMap {
		if Ω.dateRange != nil && (Ω.dateRange.Start > date || Ω.dateRange.End < date) {
			continue
		}
		sort.Strings(files)
		res = append(res, path.Join(prefix, date, files[len(files)-1]))
	}

	return res, nil
}

func (Ω *generator) generatePredictionsFromGCSPath(ctx context.Context, gcsPath string) (predictions.Predictions, error) {

	iter, err := tfrecords.NewIterator(ctx, Ω.floraStore, gcsPath)
	if err != nil {
		return nil, err
	}

	date := strings.Split(gcsPath, "/")[1]

	res := predictions.Predictions{}

	for {
		r, err := iter.Next(ctx)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		p, err := Ω.predictionFromTFRecord(r, date)
		if err != nil {
			return nil, err
		}
		if p == nil {
			continue
		}

		res = append(res, p)

	}

	return res, nil

}

func (Ω *generator) predictionFromTFRecord(r tfrecords.Record, date string) (predictions.Prediction, error) {

	Ω.totalRead++

	t, err := r.Tensor()
	if err != nil {
		return nil, err
	}

	results := Ω.model.Exec([]tf.Output{
		Ω.model.Op("dnn/head/predictions/probabilities", 0),
	}, map[tf.Output]*tf.Tensor{
		Ω.model.Op("input_example_tensor", 0): t,
	})

	value := results[0].Value().([][]float32)

	if value[0][0] > 0.5 {
		return nil, nil
	}

	lng, err := r.Longitude()
	if err != nil {
		return nil, err
	}

	lat, err := r.Latitude()
	if err != nil {
		return nil, err
	}

	return predictions.NewPrediction(Ω.nameUsageID, date, lat, lng, float64(value[0][1]))
}
