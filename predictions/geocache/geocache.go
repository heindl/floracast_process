package geocache

import (
	"bitbucket.org/heindl/taxa/store"
	"bitbucket.org/heindl/taxa/utils"
	"fmt"
	"github.com/elgs/gostrgen"
	"github.com/kellydunn/golang-geo"
	"github.com/saleswise/errors/errors"
	"github.com/tidwall/buntdb"
	"os"
	"path"
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

	tmp := path.Join("/tmp/", fmt.Sprintf("predictions-%s", random_string))
	if err := os.Mkdir(tmp, os.ModePerm); err != nil {
		return nil, errors.Wrap(err, "could not create tmp path")
	}

	db, err := buntdb.Open(path.Join(tmp, "data.db"))
	if err != nil {
		return nil, errors.Wrap(err, "could not open cache")
	}

	fmt.Println("TEMP GEOCACHE", path.Join(tmp, "data.db"))

	return &CacheWriter{DB: db}, nil

}

func ensureIndex(tx *buntdb.Tx, taxon, date string) (string, error) {

	if len(date) != 8 {
		return "", errors.New("invalid date")
	}

	if taxon == "" {
		taxon = "taxa"
	}

	index := taxon + "-" + date

	indexes, err := tx.Indexes()
	if err != nil {
		return "", err
	}

	if !utils.Contains(indexes, index) {
		pattern := index + ":*:pos"
		if err := tx.CreateSpatialIndex(index, pattern, buntdb.IndexRect); err != nil {
			return "", err
		}
	}

	return index, nil
}

func bbox(lat, lng, radius float64) string {
	//LatLon bounding box: [-112.26 33.51],[-112.18 33.67]
	centre := geo.NewPoint(lat, lng)
	sw := centre.PointAtDistanceAndBearing(radius, 225)
	ne := centre.PointAtDistanceAndBearing(radius, 45)
	return fmt.Sprintf("[%.6f %.6f],[%.6f %.6f]", sw.Lng(), sw.Lat(), ne.Lng(), ne.Lat())
}

func datesMatch(qDate, tDate string) bool {
	if qDate == "" {
		return true
	}
	if len(qDate) == 2 && qDate == tDate[4:6] {
		return true
	}
	if len(qDate) == 8 && qDate == tDate {
		return true
	}
	return false
}

func (Ω *CacheWriter) ReadTaxa(lat, lng, radius float64, qDate string, taxon string) ([]string, error) {

	fmt.Println("reading", lat, lng, radius, qDate, taxon)
	res := []string{}

	if len(qDate) != 8 {
		return nil, errors.New("invalid date")
	}

	column := "taxa-" + qDate
	if strings.TrimSpace(taxon) != "" {
		column = taxon + "-" + qDate
	}

	fmt.Println("column", column)

	Ω.DB.View(func(tx *buntdb.Tx) error {
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
	})
	return res, nil
}

func (Ω *CacheWriter) WritePredictionLine(p *store.Prediction) error {
	//Species:taxon_id,date:pos
	//taxon_id:date,id,prediction:

	if err := Ω.DB.Update(func(tx *buntdb.Tx) error {

		k1, err := ensureIndex(tx, string(p.TaxonID), p.FormattedDate)
		if err != nil {
			return err
		}

		k2, err := ensureIndex(tx, "taxa", p.FormattedDate)
		if err != nil {
			return err
		}

		v := fmt.Sprintf("%s,%s,%.8f,%.8f:pos", p.TaxonID, p.WildernessAreaID, p.ScaledPredictionValue, p.ScarcityValue)
		pos := fmt.Sprintf("[%.6f %.6f]", p.Location.Longitude, p.Location.Latitude)
		fmt.Println(k1+":"+v, pos)
		fmt.Println(k2+":"+v, pos)
		if _, _, err := tx.Set(k1+":"+v, pos, nil); err != nil {
			return err
		}
		if _, _, err := tx.Set(k2+":"+v, pos, nil); err != nil {
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
