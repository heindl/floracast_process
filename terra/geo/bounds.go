package geo

import (
	"github.com/dropbox/godropbox/errors"
	"github.com/paulmach/go.geojson"
)

type Bound struct {
	North, South, East, West float64
}

func (Ω *Bound) Widen(b *Bound) {
	if Ω == nil {
		*Ω = Bound{}
	}
	if Ω.West == 0 || b.West < Ω.West {
		Ω.West = b.West
	}
	if Ω.South == 0 || b.South < Ω.South {
		Ω.South = b.South
	}
	if Ω.East == 0 || b.East > Ω.East {
		Ω.East = b.East
	}
	if Ω.North == 0 || b.North > Ω.North {
		Ω.North = b.North
	}
}

func (Ω *Bound) WidthMeters() (float64, error) {
	p, err := Ω.Center()
	if err != nil {
		return 0, err
	}
	w, err := NewPoint(p.Latitude(), Ω.West)
	if err != nil {
		return 0, err
	}
	e, err := NewPoint(p.Latitude(), Ω.East)
	if err != nil {
		return 0, err
	}
	return w.DistanceMeters(e), nil
}

func (Ω *Bound) HeightMeters() (float64, error) {
	p, err := Ω.Center()
	if err != nil {
		return 0, err
	}
	n, err := NewPoint(Ω.North, p.Longitude())
	if err != nil {
		return 0, err
	}
	s, err := NewPoint(Ω.South, p.Longitude())
	if err != nil {
		return 0, err
	}
	return n.DistanceMeters(s), nil
}

type Bounds []*Bound

func BoundFromPolygons(polygons [][][]float64) (*Bound, error) {
	res := &Bound{}
	for _, polygon := range polygons {
		b, err := BoundFromPolygon(polygon)
		if err != nil {
			return nil, err
		}
		res.Widen(b)
	}
	return res, nil
}

func BoundFromPolygon(polygon [][]float64) (*Bound, error) {
	b := Bound{}
	for _, coord := range polygon {
		lng := coord[0]
		lat := coord[1]
		if b.West == 0 || lng < b.West {
			b.West = lng
		}
		if b.South == 0 || lat < b.South {
			b.South = lat
		}
		if b.East == 0 || lng > b.East {
			b.East = lng
		}
		if b.North == 0 || lat > b.North {
			b.North = lat
		}
	}
	return &b, nil
}

func (Ω *Bound) Center() (*Point, error) {
	ne, err := NewPoint(Ω.North, Ω.East)
	if err != nil {
		return nil, err
	}
	sw, err := NewPoint(Ω.South, Ω.West)
	if err != nil {
		return nil, err
	}
	return ne.MidPoint(sw)
}

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
