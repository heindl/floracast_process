package ecoregions

import (
	"os"
	"bufio"
	"github.com/buger/jsonparser"
	"io/ioutil"
	geojson "github.com/tidwall/tile38/geojson"
	"gopkg.in/tomb.v2"
	"github.com/saleswise/errors/errors"
	"sync"
)

type EcoRegions interface{
	PointWithin(lat, lng float64) (string, string, error)
}

func NewEcoRegions(wwf_geojson string) (EcoRegions, error) {
	r := &regions{}
	if err := r.Load(wwf_geojson); err != nil {
		return nil, err
	}
	return EcoRegions(r), nil
}

type regions struct {
	tmb     tomb.Tomb
	bBox    geojson.BBox
	sync.Mutex
	objects []geoObject
}

type geoObject struct {
	key string
	obj geojson.Object
	name string
}

func (Ω *regions) Load(wwf_geojson string) error {
	// Read huge json file using decoder.
	f, err := os.Open(wwf_geojson)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(bufio.NewReader(f))
	if err != nil {
		return err
	}

	Ω.tmb = tomb.Tomb{}
	Ω.bBox = geojson.New2DBBox(-124.848974, 24.396308, -66.885444, 49.384358)

	Ω.tmb.Go(func() error {
		if _, err := jsonparser.ArrayEach(b, Ω.parse, "features"); err != nil {
			return errors.Wrap(err, "could not parse features array")
		}
		return nil
	})

	return Ω.tmb.Wait()
}

func (Ω *regions) PointWithin(lat, lng float64) (name string, key string, err error) {
	p := geojson.New2DPoint(lng, lat)
	for _, o := range Ω.objects {
		if p.Within(o.obj) {
			return o.name, o.key, nil
		}
	}
	return "", "", nil
}

func(Ω *regions) parse(value []byte, dataType jsonparser.ValueType, offset int, err error) {
	Ω.tmb.Go(func() error {
		if err != nil {
			return err
		}
		obj, err := geojson.ObjectJSON(string(value))
		if err != nil {
			return errors.Wrap(err, "could not parse geojson")
		}
		if !obj.IntersectsBBox(Ω.bBox) {
			return nil
		}

		name, _, _, err := jsonparser.Get(value, "properties", "ECO_NAME")
		if err != nil {
			return errors.Wrap(err, "could not read ECO_NUM property")
		}
		key, _, _, err := jsonparser.Get(value, "properties", "ECO_SYM")
		if err != nil {
			return errors.Wrap(err, "could not read ECO_NUM property")
		}
		Ω.Lock()
		defer Ω.Unlock()
		Ω.objects = append(Ω.objects, geoObject{name: string(name), key: string(key), obj: obj})

		return nil
	})
}