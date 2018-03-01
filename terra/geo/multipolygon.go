package geo

import (
	"bytes"
	"github.com/golang/geo/s2"
)

type MultiPolygon []*s2.Polygon

func (Ω MultiPolygon) Empty() bool {
	for _, polygon := range Ω {
		for _, l := range polygon.Loops() {
			if len(l.Vertices()) > 0 {
				return false
			}
		}
	}
	return true
}

func (Ω MultiPolygon) Area() float64 {
	area := 0.0
	for _, p := range Ω {
		for _, l := range p.Loops() {
			if l.IsHole() || len(l.Vertices()) == 0 {
				continue
			}
			// // Multiply by radius of earth (km) square
			// (4π)r2
			area += area_to_kilometers(l.Area())
			break
		}
	}
	return area
}

func area_to_kilometers(a float64) float64 {
	return a * (6371 * 6371)
}

func (Ω MultiPolygon) Contains(lat, lng float64) bool {
	p := s2.PointFromLatLng(s2.LatLngFromDegrees(lat, lng))
	for _, polygon := range Ω {
		if polygon.ContainsPoint(p) {
			return true
		}
	}
	return false
}

func (Ω MultiPolygon) ReferencePoints() [][2]float64 {
	res := [][2]float64{}
	for _, polygon := range Ω {
		rp := polygon.ReferencePoint()
		//ll := s2.LatLngFromPoint(rp.Point).Normalized()
		ll := s2.LatLngFromPoint(rp.Point)
		res = append(res, [2]float64{ll.Lat.Degrees(), ll.Lng.Degrees()})
	}
	return res
}

func (Ω MultiPolygon) LargestPolygon() *s2.Polygon {
	maxArea := 0.0
	var maxPolygon *s2.Polygon
	for _, polygon := range Ω {
		for _, l := range polygon.Loops() {
			if l.IsHole() {
				continue
			}
			a := l.Area()
			if a <= maxArea {
				continue
			}
			maxArea = a
			maxPolygon = polygon
			break
		}
	}
	return maxPolygon
}

func (Ω MultiPolygon) Decode(a [][]byte) (MultiPolygon, error) {
	mp := MultiPolygon{}
	for _, b := range a {
		r := bytes.NewReader(b)
		p := s2.Polygon{}
		if err := p.Decode(r); err != nil {
			return nil, err
		}
		var err error
		mp, err = mp.PushPolygon(&p)
		if err != nil {
			return nil, err
		}
	}
	return mp, nil
}

func (Ω MultiPolygon) Encode() ([][]byte, error) {

	res := [][]byte{}

	for _, p := range Ω {
		var b bytes.Buffer
		if err := p.Encode(&b); err != nil {
			return nil, err
		}
		res = append(res, b.Bytes())
	}

	return res, nil
}

func (Ω MultiPolygon) ToArray() [][][][]float64 {

	multipolygon_array := [][][][]float64{}
	for _, polygon := range Ω {
		multipolygon_array = append(multipolygon_array, PolygonToArray(polygon))
	}
	return multipolygon_array
}

func PolygonToArray(polygon *s2.Polygon) [][][]float64 {
	if polygon == nil {
		return nil
	}
	polygon_array := [][][]float64{}
	holes := [][][]float64{}
	for _, loop := range polygon.Loops() {
		loop_array := [][]float64{}
		for _, v := range loop.Vertices() {
			coords := s2.LatLngFromPoint(v)
			loop_array = append(loop_array, []float64{coords.Lng.Degrees(), coords.Lat.Degrees()})
		}
		if loop.IsHole() {
			holes = append(holes, loop_array)
		} else {
			polygon_array = append(polygon_array, loop_array)
		}
	}
	return append(polygon_array, holes...)
}

func (Ω MultiPolygon) PolylabelOfLargestPolygon() (*Point, error) {
	largestPolygon := Ω.LargestPolygon()
	if largestPolygon == nil {
		return nil, nil
	}
	return PolyLabel(PolygonToArray(largestPolygon), 0)
}

func (Ω MultiPolygon) CentroidOfLargestPolygon() (*Point, error) {
	largest := 0.0
	coord := &Point{}
	for _, polygon := range Ω {
		for _, l := range polygon.Loops() {
			if l.IsHole() {
				continue
			}
			a := l.Area()
			if a <= largest {
				continue
			}
			largest = a
			ll := s2.LatLngFromPoint(l.Centroid())
			var err error
			coord, err = NewPoint(ll.Lat.Degrees(), ll.Lng.Degrees())
			if err != nil {
				return nil, err
			}
		}
	}
	return coord, nil
}

func (Ω MultiPolygon) PushPolygon(_p *s2.Polygon) (MultiPolygon, error) {
	p := *_p
	return append(Ω, &p), nil
}

func (Ω MultiPolygon) PushMultiPolygon(mp MultiPolygon) (MultiPolygon, error) {
	var err error
	for _, _p := range mp {
		p := *_p
		Ω, err = Ω.PushPolygon(&p)
		if err != nil {
			return nil, err
		}
	}
	return Ω, nil
}

// Note that this function orders the points counter clockwise.
// Assumes geojson format in that outer ring is the only outer bound and all remaining rings are holes.
func NewPolygon(pArray [][][]float64) (*s2.Polygon, error) {
	var loops []*s2.Loop
	for _, ring := range pArray {
		if l := newLoop(ring); l != nil {
			loops = append([]*s2.Loop{l}, loops...)
		}
	}
	return s2.PolygonFromLoops(loops), nil
}

func newLoop(ar [][]float64) *s2.Loop {
	if len(ar) < 3 {
		return nil
	}
	points := []s2.Point{}
	for i := range ar {
		points = append(points, s2.PointFromLatLng(s2.LatLngFromDegrees(ar[i][1], ar[i][0])))
	}
	loop := s2.LoopFromPoints(points)
	if !isPositivelyOriented(loop) {
		loop.Invert()
	}
	return loop
}

// Outer ring must be positively oriented (counter-clockwise) while holes should be negatively oriented.
// https://stackoverflow.com/questions/1165647/how-to-determine-if-a-list-of-polygon-points-are-in-clockwise-order
func isPositivelyOriented(loop *s2.Loop) bool {

	vertices := loop.Vertices()

	k := len(vertices) - 1
	total := 0.0
	for i := range vertices {
		v1 := vertices[k]
		v1LatLng := s2.LatLngFromPoint(v1)
		v2 := vertices[i]
		v2LatLng := s2.LatLngFromPoint(v2)
		// Cross() -> x - y - z?
		total += (v1LatLng.Lng.Degrees()*v2LatLng.Lat.Degrees() - v2LatLng.Lng.Degrees()*v1LatLng.Lat.Degrees())
		//total += (v1.X * v2.Y - v2.X * v1.Y)
		//total += points[i].Dot(points[k].Vector)
		k = i
	}
	return total > 0
}
