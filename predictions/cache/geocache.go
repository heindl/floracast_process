package cache

import (
	"fmt"
	"os"
	"path"
	"strings"

	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/predictions"
	"bitbucket.org/heindl/process/utils"
	"github.com/dropbox/godropbox/errors"
	"github.com/elgs/gostrgen"
	"github.com/kellydunn/golang-geo"
	"github.com/tidwall/buntdb"
	"gopkg.in/tomb.v2"
)

type localGeoCacheWriter struct {
	DB *buntdb.DB
}

func NewLocalGeoCache() (PredictionCache, func() error, error) {

	random_string, err := gostrgen.RandGen(10, gostrgen.Lower, "", "")
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not create temp file random string name")
	}

	tmp := path.Join("/tmp/", fmt.Sprintf("predictions-%s", random_string))
	if err = os.Mkdir(tmp, os.ModePerm); err != nil {
		return nil, nil, errors.Wrap(err, "could not create tmp path")
	}

	db, err := buntdb.Open(path.Join(tmp, "data.db"))
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not open cache")
	}

	fmt.Println("TEMP GEOCACHE", path.Join(tmp, "data.db"))

	c := localGeoCacheWriter{DB: db}

	return &c, c.Close, nil

}

func bbox(lat, lng, radius float64) string {
	//LatLon bounding box: [-112.26 33.51],[-112.18 33.67]
	centre := geo.NewPoint(lat, lng)
	sw := centre.PointAtDistanceAndBearing(radius, 225)
	ne := centre.PointAtDistanceAndBearing(radius, 45)
	return fmt.Sprintf("[%.6f %.6f],[%.6f %.6f]", sw.Lng(), sw.Lat(), ne.Lng(), ne.Lat())
}

func (Ω *localGeoCacheWriter) ReadPredictions(lat, lng, radius float64, qDate string, nameUsageID *nameusage.ID) ([]string, error) {

	res := []string{}

	column, err := newIndexKey(nameUsageID, qDate)
	if err != nil {
		return nil, err
	}

	if err := Ω.DB.View(func(tx *buntdb.Tx) error {
		return tx.Intersects(column, bbox(lat, lng, radius), func(key, val string) bool {
			k := strings.Split(strings.Split(key, ":")[1], ",")
			taxonID := k[0]
			areaID := k[1]
			prediction := k[2]
			scarcity := k[3]
			v := strings.Split(strings.Trim(val, "[]"), " ")
			res = append(res, fmt.Sprintf("%s,%s,%s,%s,%s,%s", taxonID, areaID, prediction, scarcity, v[1], v[0]))
			return true
		})
	}); err != nil {
		return nil, err
	}
	return res, nil
}

func (Ω *localGeoCacheWriter) WritePredictions(prediction_list predictions.Predictions) error {
	//Species:taxon_id,date:pos
	//taxon_id:date,id,prediction:
	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, _prediction := range prediction_list {
			prediction := _prediction
			tmb.Go(func() error {
				transaction, err := newTransaction(prediction)
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

func newTransaction(p predictions.Prediction) (func(*buntdb.Tx) error, error) {
	indexKeys, err := allIndexKeys(p)
	if err != nil {
		return nil, err
	}

	valueLine, err := getValueLine(p)
	if err != nil {
		return nil, err
	}

	lat, lng, err := p.LatLng()
	if err != nil {
		return nil, err
	}

	return func(tx *buntdb.Tx) error {
		if err := ensureIndexes(tx, indexKeys); err != nil {
			return nil
		}
		pos := fmt.Sprintf("[%.6f %.6f]", lng, lat)
		for _, k := range indexKeys {
			if _, _, err := tx.Set(k+":"+valueLine, pos, nil); err != nil {
				return err
			}
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

func getValueLine(p predictions.Prediction) (string, error) {

	usageID, err := p.UsageID()
	if err != nil {
		return "", err
	}

	areaID, err := p.ProtectedArea()
	if err != nil {
		return "", err
	}

	scaledValue, err := p.ScaledPrediction()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s,%s,%.8f:pos", usageID, areaID, scaledValue), nil
}

func allIndexKeys(p predictions.Prediction) ([]string, error) {

	usageID, err := p.UsageID()
	if err != nil {
		return nil, err
	}

	date, err := p.Date()
	if err != nil {
		return nil, err
	}

	k1, err := newIndexKey(&usageID, date)
	if err != nil {
		return nil, err
	}

	k2, err := newIndexKey(nil, date)
	if err != nil {
		return nil, err
	}

	return []string{k1, k2}, nil
}

func newIndexKey(usageID *nameusage.ID, date string) (string, error) {
	if len(date) != 8 {
		return "", errors.Newf("Invalid Prediction Date [%s]", date)
	}

	if usageID != nil {
		if !usageID.Valid() {
			return "", errors.Newf("Invalid UsageID [%s]", usageID)
		}
		return fmt.Sprintf("%s-%s", usageID, date), nil
	}

	return fmt.Sprintf("taxa-%s", date), nil
}

func ensureIndexes(tx *buntdb.Tx, newIndexes []string) error {

	existingIndexes, err := tx.Indexes()
	if err != nil {
		return err
	}

	for _, indx := range newIndexes {
		if !utils.ContainsString(existingIndexes, indx) {
			pattern := fmt.Sprintf("%s:*:pos", indx)
			if err := tx.CreateSpatialIndex(indx, pattern, buntdb.IndexRect); err != nil {
				return err
			}
		}
	}

	return nil
}
