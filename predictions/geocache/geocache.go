package geocache

import (
	"github.com/tidwall/buntdb"
	"fmt"
	"github.com/saleswise/errors/errors"
	"github.com/elgs/gostrgen"
	"path"
	"os"
	"bitbucket.org/heindl/taxa/store"
	"github.com/kellydunn/golang-geo"
	"strings"
)

type CacheWriter struct {
	DB *buntdb.DB
}

func NewCacheWriter(taxa []string) (*CacheWriter, error) {

	random_string, err := gostrgen.RandGen(10, gostrgen.Lower, "", "")
	if err != nil {
		return nil, errors.Wrap(err, "could not create temp file random string name")
	}

	tmp := path.Join(os.TempDir(), fmt.Sprintf("predictions-%s", random_string))
	if err := os.Mkdir(tmp, os.ModePerm); err != nil {
		return nil, errors.Wrap(err, "could not create tmp path")
	}


	db, err := buntdb.Open(path.Join(tmp, "data.db"))
	if err != nil {
		return nil, errors.Wrap(err, "could not open cache")
	}

	fmt.Println("TEMP GEOCACHE", path.Join(tmp, "data.db"))

	if err := db.Update(func(tx *buntdb.Tx) error {
		for _, taxon := range taxa {
			if err := tx.CreateSpatialIndex("taxa","taxa:*:pos", buntdb.IndexRect); err != nil {
				return errors.Wrap(err, "could not create spatial index")
			}
			k := fmt.Sprintf("%s:*:pos", taxon)
			if err := tx.CreateSpatialIndex(string(taxon),k, buntdb.IndexRect); err != nil {
				return errors.Wrap(err, "could not create spatial index")
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &CacheWriter{DB: db}, nil

}

func bbox(lat, lng, radius float64) string {
	//LatLon bounding box: [-112.26 33.51],[-112.18 33.67]
	centre := geo.NewPoint(lat, lng)
	sw := centre.PointAtDistanceAndBearing(radius, 225)
	ne := centre.PointAtDistanceAndBearing(radius, 45)
	return fmt.Sprintf("[%.6f %.6f],[%.6f %.6f]", sw.Lng(), sw.Lat(), ne.Lng(), ne.Lat())
}

func (Ω *CacheWriter) ReadTaxon(taxon store.TaxonID, lat, lng, radius float64) ([]string, error) {
	res := []string{}
	//LatLon bounding box: [-112.26 33.51],[-112.18 33.67]
	Ω.DB.View(func(tx *buntdb.Tx) error {
		return tx.Intersects(string(taxon), bbox(lat, lng, radius), func(key, val string) bool {
			k := strings.Split(strings.Split(key,":")[1], ",")
			v := strings.Split(strings.Trim(val, "[]"), " ")
			res = append(res, fmt.Sprintf("%s,%s,%s,%s,%s", k[0], k[1], k[2], v[1], v[0]))
			return true
		})
	})
	return res, nil
}

func (Ω *CacheWriter) ReadTaxa(lat, lng, radius float64) ([]string, error) {
	res := []string{}
	Ω.DB.View(func(tx *buntdb.Tx) error {
		return tx.Intersects("taxa", bbox(lat, lng, radius), func(key, val string) bool {
			k := strings.Split(strings.Split(key,":")[1], ",")
			v := strings.Split(strings.Trim(val, "[]"), " ")
			res = append(res, fmt.Sprintf("%s,%s,%s,%s", k[0], k[1], v[1], v[0]))
			return true
		})
	})
	return res, nil
}

func (Ω *CacheWriter) WritePredictionLine(p store.Prediction) error {
	//Species:taxon_id,date:pos
	//taxon_id:date,id,prediction:
	if err := Ω.DB.Update(func(tx *buntdb.Tx) error {
		k1 := fmt.Sprintf("%s:%s,%s,%.8f:pos", p.TaxonID, p.FormattedDate, p.WildernessAreaID, p.PredictionValue)
		k2 := fmt.Sprintf("taxa:%s,%s:pos", p.TaxonID, p.FormattedDate)
		pos := fmt.Sprintf("[%.6f %.6f]", p.Location.Longitude, p.Location.Latitude)
		if _, _, err := tx.Set(k1, pos, nil); err != nil {
			return err
		}
		if _, _, err := tx.Set(k2, pos, nil); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "could not update prediction cache")
	}
	return nil
}

func (Ω *CacheWriter) Close() error {
	if err := Ω.DB.Close(); err != nil {
		return errors.Wrap(err, "could not close db")
	}
	return nil
}


