package terra

import (
	"github.com/golang/geo/s2"
	"fmt"
	"math"
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

func (Ω MultiPolygon) Contains(lat, lng float64) bool {
	p := s2.PointFromLatLng(s2.LatLngFromDegrees(lat, lng))
	for i, polygon := range Ω {

		for  _, l := range polygon.Loops() {
			fmt.Println(i, fmt.Sprintf("%.10f", l.Area()), len(l.Vertices()), l.IsHole(), math.Ceil(l.Area()), s2.LatLngFromPoint(l.Centroid()).String())
		}
		if polygon.ContainsPoint(p) {
			return true
		}
	}
	return false
}

func (Ω MultiPolygon) ReferencePoints() [][2]float64 {
	res := [][2]float64{}
	for _, polygon := range Ω {
		for _, l := range polygon.Loops() {
			rp := l.ReferencePoint()
			//ll := s2.LatLngFromPoint(rp.Point).Normalized()
			ll := s2.LatLngFromPoint(rp.Point)
			res = append(res, [2]float64{ll.Lat.Degrees(), ll.Lng.Degrees()})
		}
	}
	return res
}

func (Ω MultiPolygon) PushPolygon(p *s2.Polygon) MultiPolygon {
	return append(Ω, p)
}

func (Ω MultiPolygon) PushMultiPolygon(mp MultiPolygon) MultiPolygon {
	for _, p := range mp {
		Ω = Ω.PushPolygon(p)
	}
	return Ω
}

// Note that this function orders the points counter clockwise.
// Assumes geojson format in that outer ring is the only outer bound and all remaining rings are holes.
func NewPolygon(pArray [][][]float64) (*s2.Polygon, error) {
	var loops []*s2.Loop
	for _, ring := range pArray {
		if l := newLoop(ring); l != nil {
			loops = append([]*s2.Loop{l}, loops...)
		}
		break
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
		v2 := vertices[i]
		// Cross() -> x - y - z?
		total += (v1.X * v2.Y - v2.X * v1.Y)
		//total += points[i].Dot(points[k].Vector)
		k = i
	}
	return total > 0
}