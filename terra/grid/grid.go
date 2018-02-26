package grid

import (
	"github.com/golang/geo/s2"
	"bitbucket.org/heindl/process/ecoregions"
	"github.com/paulmach/go.geojson"
	"gopkg.in/tomb.v2"
	"sync"
	"github.com/dropbox/godropbox/errors"
)

type Generator interface{
	SubDivide(g *Bound, level int) (Bounds, error)
}

func NewGridGenerator() (Generator, error) {
	ecoCache, err := ecoregions.NewEcoRegionsCache()
	if err != nil {
		return nil, err
	}
	return &generator{
		ecoRegionCache: ecoCache,
		grids: Bounds{},
	}, nil
}

type generator struct {
	ecoRegionCache *ecoregions.EcoRegionsCache
	grids []*Bound
	sync.Mutex
}

var NorthAmerica = &Bound{
	North: 53.5555501, // Edmonton, Alberta
	West: -137.8424302, // Glacier Bay
	East: -53.1078873, // St. Johns, Newfoundland
	South: 20.6737777, // Guadalajara, Mexico
}

type Bound struct {
	North, South, East, West float64
}

type Bounds []*Bound

func (Ω Bounds) ToGeoJSON() ([]byte, error) {
	if len(Ω) == 0 {
		return nil, errors.New("At least one Bound is required for GeoJSON")
	}

	fc := geojson.NewFeatureCollection()
	for _, b := range Ω {
		f := geojson.NewPolygonFeature(append([][][]float64{}, [][]float64{
			{b.West, b.North},
			{b.East, b.North},
			{b.East, b.South},
			{b.West, b.South},
			{b.West, b.North},
		}))
		fc = fc.AddFeature(f)
	}
	b, err := fc.MarshalJSON()
	if err != nil {
		return nil, errors.Wrap(err, "Could not marshal FeatureCollection")
	}
	return b, nil
}

func (Ω *generator) SubDivide(g *Bound, level int) (Bounds, error) {

	region := s2.Region(s2.RectFromLatLng(s2.LatLngFromDegrees(g.North, g.East)).AddPoint(s2.LatLngFromDegrees(g.South, g.West)))
	regionCoverer := &s2.RegionCoverer{MaxLevel: level, MinLevel: level, MaxCells: 500}

	covering := regionCoverer.Covering(region)

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, _cellID := range covering {
			cellID := _cellID
			tmb.Go(func() error {
				return Ω.parseCell(cellID)
			})
		}
		return nil
	})
	if err := tmb.Wait(); err != nil {
		return nil, err
	}

	return Ω.grids, nil
}

func (Ω *generator) parseCell(cellID s2.CellID) error {
	rb := s2.CellFromCellID(cellID).RectBound()

	n := rb.Hi().Lat.Degrees()
	e := rb.Hi().Lng.Degrees()
	s := rb.Lo().Lat.Degrees()
	w := rb.Lo().Lng.Degrees()

	f := geojson.NewPolygonFeature(append([][][]float64{}, [][]float64{
		{w, n},
		{e, n},
		{e, s},
		{w, s},
		{w, n},
	}))

	ecoRegionTouches := 0

	if _, err := Ω.ecoRegionCache.EcoID(rb.Center().Lat.Degrees(), rb.Center().Lng.Degrees()); err == nil {
		ecoRegionTouches += 1
	}

	for _, loc := range f.Geometry.Polygon[0][:3] {
		_, err := Ω.ecoRegionCache.EcoID(loc[1], loc[0])
		if err == ecoregions.ErrNotFound {
			continue
		}
		if err != nil {
			return errors.Wrapf(err, "Could not get EcoID from location [%f, %f]", loc[0], loc[1])
		}
		ecoRegionTouches += 1
		if ecoRegionTouches > 1 {
			Ω.Lock()
			defer Ω.Unlock()
			Ω.grids = append(Ω.grids, &Bound{
				North: n,
				South: s,
				East: e,
				West: w,
			})
			return nil
		}
	}

	return nil

}
