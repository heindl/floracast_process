package grid

import (
	"bitbucket.org/heindl/process/terra/ecoregions"
	"bitbucket.org/heindl/process/terra/ecoregions/cache"
	"bitbucket.org/heindl/process/terra/geo"
	"bitbucket.org/heindl/process/utils"
	"github.com/dropbox/godropbox/errors"
	"github.com/golang/geo/s2"
	"github.com/paulmach/go.geojson"
	"gopkg.in/tomb.v2"
	"sync"
)

type Generator interface {
	SubDivide(g *geo.Bound, level int) (geo.Bounds, error)
}

func NewGridGenerator() (Generator, error) {
	return &generator{
		grids: geo.Bounds{},
	}, nil
}

type generator struct {
	grids []*geo.Bound
	sync.Mutex
}

var NorthAmerica = &geo.Bound{
	North: 53.5555501,   // Edmonton, Alberta
	West:  -137.8424302, // Glacier Bay
	East:  -53.1078873,  // St. Johns, Newfoundland
	South: 20.6737777,   // Guadalajara, Mexico
}

func (Ω *generator) SubDivide(g *geo.Bound, level int) (geo.Bounds, error) {

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

	if _, err := cache.FetchEcologicalRegion(rb.Center().Lat.Degrees(), rb.Center().Lng.Degrees()); err == nil {
		ecoRegionTouches += 1
	}

	for _, loc := range f.Geometry.Polygon[0][:3] {
		_, err := cache.FetchEcologicalRegion(loc[1], loc[0])
		if utils.ContainsError(err, ecoregions.ErrNotFound) {
			continue
		}
		if err != nil {
			return errors.Wrapf(err, "Could not get ecoID from location [%f, %f]", loc[0], loc[1])
		}
		ecoRegionTouches += 1
		if ecoRegionTouches > 1 {
			Ω.Lock()
			defer Ω.Unlock()
			Ω.grids = append(Ω.grids, &geo.Bound{
				North: n,
				South: s,
				East:  e,
				West:  w,
			})
			return nil
		}
	}

	return nil

}
