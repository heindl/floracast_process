package cache

import (
	"fmt"
	"strings"

	"bitbucket.org/heindl/process/predictions"
	"github.com/dropbox/godropbox/errors"
	"github.com/tidwall/buntdb"
	"gopkg.in/tomb.v2"
	"sync"
)

type localGeoCacheWriter struct {
	DB          *buntdb.DB
	Predictions map[string]predictions.Prediction
	sync.Mutex
}

func predictionKey(p predictions.Prediction) (string, error) {
	loc, err := p.ProtectedArea()
	if err != nil {
		return "", err
	}
	date, err := p.Date()
	if err != nil {
		return "", err
	}
	nameUsageID, err := p.NameUsageID()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s+%s+%s", loc, date, nameUsageID), nil
}

// NewLocalGeoCache creates a PredictionCache with additional geoquery methods.
func NewLocalGeoCache(cachePath string) (PredictionCache, func() error, error) {

	//if cachePath == "" {
	//	randomString, err := gostrgen.RandGen(10, gostrgen.Lower, "", "")
	//	if err != nil {
	//		return nil, nil, errors.Wrap(err, "could not create temp file random string name")
	//	}
	//	cachePath = path.Join("/tmp/", fmt.Sprintf("predictions-%s", randomString))
	//	if err = os.Mkdir(cachePath, os.ModePerm); err != nil {
	//		return nil, nil, errors.Wrap(err, "could not create tmp path")
	//	}
	//	cachePath = path.Join(cachePath, "data.db")
	//}

	fmt.Println("opening cache path", cachePath)
	db, err := buntdb.Open(":memory:")
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not open cache")
	}

	c := localGeoCacheWriter{
		DB:          db,
		Predictions: map[string]predictions.Prediction{},
	}

	return &c, c.Close, nil

}

//func bbox(lat, lng, radius float64) string {
//	//LatLon bounding box: [-112.26 33.51],[-112.18 33.67]
//	centre := geo.NewPoint(lat, lng)
//	sw := centre.PointAtDistanceAndBearing(radius, 225)
//	ne := centre.PointAtDistanceAndBearing(radius, 45)
//	return fmt.Sprintf("[%.6f %.6f],[%.6f %.6f]", sw.Lng(), sw.Lat(), ne.Lng(), ne.Lat())
//}

func (Ω *localGeoCacheWriter) ReadPredictions(bounds string) (predictions.Predictions, error) {

	res := predictions.Predictions{}

	if err := Ω.DB.View(func(tx *buntdb.Tx) error {

		return tx.Intersects("taxa", bounds, func(key, val string) bool {

			res = append(res, Ω.Predictions[strings.Split(key, ":")[1]])
			return true
		})
	}); err != nil && err != buntdb.ErrNotFound {
		return nil, err
	}
	return res, nil
}

func (Ω *localGeoCacheWriter) WritePredictions(predictionList predictions.Predictions) error {
	//Species:taxon_id,date:pos
	//taxon_id:date,id,prediction:

	// Set Index:
	if err := Ω.DB.Update(func(tx *buntdb.Tx) error {
		return tx.CreateSpatialIndex("taxa", "taxa:*:pos", buntdb.IndexRect)
	}); err != nil {
		return err
	}

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, _prediction := range predictionList {
			prediction := _prediction
			tmb.Go(func() error {
				transaction, err := Ω.newTransaction(prediction)
				if err != nil {
					return err
				}
				return Ω.DB.Update(transaction)
			})
		}
		return nil
	})
	return nil
}

func (Ω *localGeoCacheWriter) newTransaction(p predictions.Prediction) (func(*buntdb.Tx) error, error) {

	k, err := predictionKey(p)
	if err != nil {
		return nil, err
	}

	lat, lng, err := p.Coordinates()
	if err != nil {
		return nil, err
	}

	Ω.Lock()
	defer Ω.Unlock()

	Ω.Predictions[k] = p

	return func(tx *buntdb.Tx) error {
		if _, _, err := tx.Set(
			fmt.Sprintf("taxa:%s:pos", k),
			fmt.Sprintf("[%.6f %.6f]", lng, lat),
			nil,
		); err != nil {
			return err
		}
		return nil
	}, nil
}

func (Ω *localGeoCacheWriter) Close() error {
	if err := Ω.DB.Close(); err != nil {
		return errors.Wrap(err, "could not close db")
	}
	return nil
}

//func ensureIndexes(tx *buntdb.Tx, indx string) error {
//
//	existingIndexes, err := tx.Indexes()
//	if err != nil {
//		return err
//	}
//
//	if !utils.ContainsString(existingIndexes, indx) {
//		pattern := fmt.Sprintf("%s:*:pos", indx)
//		if err := tx.CreateSpatialIndex(indx, pattern, buntdb.IndexRect); err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
