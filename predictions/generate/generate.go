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
)

func GeneratePredictions(ctx context.Context, nameUsageID nameusage.ID, floraStore store.FloraStore, dateRange *DateRange) (list predictions.Collection, err error) {

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

	collection, err := predictions.NewCollection(nameUsageID, floraStore)
	if err != nil {
		return nil, err
	}

	g := generator{
		nameUsageID: nameUsageID,
		model:       model,
		dateRange:   dateRange,
		floraStore:  floraStore,
		collection:  collection,
	}

	filenames, err := g.fetchLatestProtectedAreaFileNames(ctx)
	if err != nil {
		return nil, err
	}

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, _f := range filenames {
			f := _f
			tmb.Go(func() error {
				return g.generatePredictionsFromGCSPath(ctx, f)
			})
		}
		return nil
	})
	if err := tmb.Wait(); err != nil {
		return nil, err
	}

	fmt.Println("TOTAL READ", g.totalRead)

	return g.collection, nil

}

type DateRange struct {
	Start, End string
}

type generator struct {
	nameUsageID nameusage.ID
	model       *tg.Model
	dateRange   *DateRange
	floraStore  store.FloraStore
	collection  predictions.Collection
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

func (Ω *generator) generatePredictionsFromGCSPath(ctx context.Context, gcsPath string) error {

	iter, err := tfrecords.NewIterator(ctx, Ω.floraStore, gcsPath)
	if err != nil {
		return err
	}

	date := utils.FormattedDate(strings.Split(gcsPath, "/")[1])

	for {
		r, err := iter.Next(ctx)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		if err := Ω.predictTFRecord(r, date); err != nil {
			return err
		}

	}

	return nil

}

func (Ω *generator) predictTFRecord(r tfrecords.Record, date utils.FormattedDate) error {

	Ω.totalRead++

	t, err := r.Tensor()
	if err != nil {
		return err
	}

	results := Ω.model.Exec([]tf.Output{
		Ω.model.Op("dnn/head/predictions/probabilities", 0),
	}, map[tf.Output]*tf.Tensor{
		Ω.model.Op("input_example_tensor", 0): t,
	})

	value := results[0].Value().([][]float32)

	if value[0][0] > 0.5 {
		return nil
	}

	lng, err := r.Longitude()
	if err != nil {
		return err
	}

	lat, err := r.Latitude()
	if err != nil {
		return err
	}

	return Ω.collection.Add(lat, lng, date, float64(value[0][1]))
}
