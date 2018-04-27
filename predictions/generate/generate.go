package generate

import (
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/predictions"
	"bitbucket.org/heindl/process/store"
	"bitbucket.org/heindl/process/tfrecords"
	"bitbucket.org/heindl/process/utils"
	"context"
	"github.com/Jeffail/tunny"
	tg "github.com/galeone/tfgo"
	"github.com/golang/glog"
	"github.com/saleswise/errors/errors"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"google.golang.org/api/iterator"
	"io"
	"runtime"
	"time"
)

type DateRange struct {
	Start, End string
}

type generator struct {
	nameUsageID nameusage.ID
	dateRange   *DateRange
	floraStore  store.FloraStore
	collection  predictions.Collection
	totalRead   int
	classifier  FloraClassifier
	batchCount  int
}

func GeneratePredictions(
	ctx context.Context,
	nameUsageID nameusage.ID,
	floraStore store.FloraStore,
	dateRange *DateRange,
	modelPath string,
	protectedAreaGlob string,
) (list predictions.Collection, err error) {

	if !nameUsageID.Valid() {
		return nil, errors.Newf("Invalid NameUsageID [%s]", nameUsageID)
	}

	classifier, err := NewFloraClassifier(ctx, floraStore, nameUsageID, modelPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not get model for NameUsage [%s]", nameUsageID)
	}
	defer utils.SafeClose(classifier, &err)

	collection, err := predictions.NewCollection(nameUsageID, floraStore)
	if err != nil {
		return nil, err
	}

	if protectedAreaGlob == "" {
		return nil, errors.New("Missing Protected Area Path")
	}

	g := generator{
		nameUsageID: nameUsageID,
		dateRange:   dateRange,
		floraStore:  floraStore,
		collection:  collection,
		classifier:  classifier,
	}

	processStartTime := time.Now()

	runtime.GOMAXPROCS(8)

	pool := tunny.New(8, g.worker)

	iter, err := tfrecords.NewLocalIterator(ctx, protectedAreaGlob)
	if err != nil {
		return nil, err
	}
	batch := tfrecords.Records{}
	_batch := tfrecords.Records{}

	for {
		r, err := iter.Next(ctx)
		if err == iterator.Done || err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		batch = append(batch, r)
		if len(batch) >= 4000 {
			_batch, batch = batch[0:], batch[:0]
			if er := pool.Process(_batch); er != nil {
				return nil, er.(error)
			}
		}
	}
	if er := pool.Process(batch); er != nil {
		return nil, er.(error)
	}

	for pool.QueueLength() > 0 {
		time.Sleep(1 * time.Second)
	}

	glog.Infof("Completed Parsing %d ProtectedAreas into %d Positive Predictions in %f Minutes", g.totalRead, g.collection.Count(), time.Now().Sub(processStartTime).Minutes())

	return g.collection, nil

	//filenames, err := g.fetchLatestProtectedAreaFileNames(ctx)
	//if err != nil {
	//	return nil, err
	//}
	//
	//tmb := tomb.Tomb{}
	//tmb.Go(func() error {
	//	for _, _f := range filenames {
	//		f := _f
	//		tmb.Go(func() error {
	//			return g.generatePredictionsFromGCSPath(ctx, f)
	//		})
	//	}
	//	return nil
	//})
	//if err := tmb.Wait(); err != nil {
	//	return nil, err
	//}
	//
	//glog.Infof(
	//	"Prediction Generation Complete for NameUsage [%s], with %d Positive Results from %d Total Predictions.",
	//	nameUsageID,
	//	g.collection.Count(),
	//	g.totalRead,
	//)

}

func (Ω *generator) handleRecord(r tfrecords.Record, model *tg.Model) error {
	Ω.totalRead++

	t, err := r.Tensor()
	if err != nil {
		return err
	}

	results := model.Exec([]tf.Output{
		model.Op("dnn/head/predictions/probabilities", 0),
	}, map[tf.Output]*tf.Tensor{
		model.Op("input_example_tensor", 0): t,
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

	date, err := r.Date()
	if err != nil {
		return err
	}

	return Ω.collection.Add(lat, lng, date, float64(value[0][1]))
}

func (Ω *generator) worker() tunny.Worker {
	m, err := Ω.classifier.NewClassifierInstance(context.Background())
	if err != nil {
		panic(err)
	}

	probability_op := m.Op("dnn/head/predictions/probabilities", 0)
	input_op := m.Op("input_example_tensor", 0)

	return &worker{
		generator:      Ω,
		model:          m,
		probability_op: &probability_op,
		input_op:       &input_op,
		lmtr:           utils.NewLimiter(1),
	}
}

type worker struct {
	generator      *generator
	model          *tg.Model
	probability_op *tf.Output
	input_op       *tf.Output
	lmtr           utils.Limiter
}

// Process will synchronously perform a job and return the result.
func (Ω *worker) Process(i interface{}) interface{} {
	defer Ω.lmtr.Release()
	Ω.generator.batchCount++

	records := i.(tfrecords.Records)

	glog.Infof("Processing TFRecord Batch [%d] with %d Records", Ω.generator.batchCount, len(records))

	if len(records) == 0 {
		return nil
	}

	Ω.generator.totalRead += len(records)

	t, err := records.Tensor()
	if err != nil {
		return err
	}

	output := Ω.model.Exec([]tf.Output{
		*Ω.probability_op,
	}, map[tf.Output]*tf.Tensor{
		*Ω.input_op: t,
	})

	results := output[0].Value().([][]float32)

	for i, r := range records {

		unlikely := results[i][0]
		likely := results[i][1]

		if unlikely >= 0.5 {
			continue
		}

		lng, err := r.Longitude()
		if err != nil {
			return err
		}

		lat, err := r.Latitude()
		if err != nil {
			return err
		}

		date, err := r.Date()
		if err != nil {
			return err
		}

		if err := Ω.generator.collection.Add(lat, lng, date, float64(likely)); err != nil {
			return err
		}
	}
	return nil
}

// BlockUntilReady is called before each job is processed and must block the
// calling goroutine until the Worker is ready to process the next job.
func (Ω *worker) BlockUntilReady() {
	Ω.lmtr.Wait()
}

// Interrupt is called when a job is cancelled. The worker is responsible
// for unblocking the Process implementation.
func (Ω *worker) Interrupt() {
	return
}

// Terminate is called when a Worker is removed from the processing pool
// and is responsible for cleaning up any held resources.
func (Ω *worker) Terminate() {
	return
}

//func (Ω *generator) fetchLatestProtectedAreaFileNames(ctx context.Context) ([]string, error) {
//
//	prefix := fmt.Sprintf("protected_areas/")
//	names, err := Ω.floraStore.CloudStorageObjectNames(ctx, prefix, ".tfrecords")
//	if err != nil {
//		return nil, err
//	}
//
//	dateFileMap := map[string][]string{}
//	for _, n := range names {
//		s := strings.Split(n, "/")
//		d := s[0]
//		f := s[1]
//		if _, ok := dateFileMap[d]; !ok {
//			dateFileMap[d] = []string{}
//		}
//		dateFileMap[d] = append(dateFileMap[d], f)
//	}
//
//	res := []string{}
//	for date, files := range dateFileMap {
//		if Ω.dateRange != nil && (Ω.dateRange.Start > date || Ω.dateRange.End < date) {
//			continue
//		}
//		sort.Strings(files)
//		res = append(res, path.Join(prefix, date, files[len(files)-1]))
//	}
//
//	return res, nil
//}
//
//func (Ω *generator) generatePredictionsFromGCSPath(ctx context.Context, gcsPath string) error {
//
//	iter, err := tfrecords.NewIterator(ctx, Ω.floraStore, gcsPath)
//	if err != nil {
//		return err
//	}
//
//	for {
//		r, err := iter.Next(ctx)
//		if err == iterator.Done {
//			break
//		}
//		if err != nil {
//			return err
//		}
//
//		if err := Ω.handleRecord(r, Ω.model); err != nil {
//			return err
//		}
//
//	}
//
//	return nil
//
//}
