package terra

import (
	"math"
	"sort"
)

func PolyLabel(polygon [][][]float64, precision float64) (*Point, error) {

	if precision == 0 {
		precision = 0.001
		//precision = 1.0
	}

	// find the bounding box of the outer ring
	var minX, minY, maxX, maxY float64
	for i := 0; i < len(polygon[0]); i++ {
		p := polygon[0][i]
		if i == 0 || p[0] < minX {
			minX = p[0]
		}
		if i == 0 || p[1] < minY {
			minY = p[1]
		}
		if i == 0 || p[0] > maxX {
			maxX = p[0]
		}
		if i == 0 || p[1] > maxY {
			maxY = p[1]
		}
	}

	width := maxX - minX
	height := maxY - minY
	cellSize := math.Min(width, height)
	h := cellSize / 2

	// a priority queue of cells in order of their "potential" (max distance to polygon)
	//var cellQueue = new Queue(null, compareMax);
	cellQueue := cells{}

	if cellSize == 0 {
		return NewPoint(minY, minX)
	}

	// cover polygon with initial cells
	for x := minX; x < maxX; x += cellSize {
		for y := minY; y < maxY; y += cellSize {
			cellQueue = append(cellQueue, newCell(x+h, y+h, h, polygon))
		}
	}

	// take centroid as the first best guess
	var bestCell = getCentroidCell(polygon)

	// special case for rectangular polygons
	var bboxCell = newCell(minX+width/2, minY+height/2, 0, polygon)
	if bboxCell.d > bestCell.d {
		bestCell = bboxCell
	}

	var numProbes = len(cellQueue)

	for len(cellQueue) > 0 {

		sort.Sort(cellQueue)
		// pick the most promising cell from the queue
		cell := cellQueue[0] // Pop
		cellQueue = cellQueue[1:]

		// update the best cell if we found a better one
		if cell.d > bestCell.d {
			bestCell = cell
		}

		// do not drill down further if there's no chance of a better solution
		if cell.max-bestCell.d <= precision {
			continue
		}

		// split the cell into four cells
		h = cell.h / 2

		cellQueue = append(cellQueue, newCell(cell.x-h, cell.y-h, h, polygon))
		cellQueue = append(cellQueue, newCell(cell.x+h, cell.y-h, h, polygon))
		cellQueue = append(cellQueue, newCell(cell.x-h, cell.y+h, h, polygon))
		cellQueue = append(cellQueue, newCell(cell.x+h, cell.y+h, h, polygon))
		numProbes += 4
	}

	return NewPoint(bestCell.y, bestCell.x)
}

type cell struct {
	x   float64 // cell center x
	y   float64 // cell center y
	h   float64 // half the cell size
	d   float64 // distance from cell center to polygon
	max float64 // max distance to polygon within a cell
}

type cells []cell

func (a cells) Len() int           { return len(a) }
func (a cells) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a cells) Less(i, j int) bool { return a[i].max > a[j].max }

func newCell(x, y, h float64, polygon [][][]float64) cell {
	c := cell{
		x: x,
		y: y,
		h: h,
		d: pointToPolygonDist(x, y, polygon),
	}
	c.max = c.d + c.h*math.Sqrt2
	return c
}

// signed distance from point to polygon outline (negative if point is outside)
func pointToPolygonDist(x, y float64, polygon [][][]float64) float64 {
	inside := false
	minDistSq := math.Inf(1) // Infinity;

	for k := 0; k < len(polygon); k++ {
		ring := polygon[k]
		ring_length := len(ring)
		j := ring_length - 1
		for i := 0; i < ring_length; i++ {
			var a = ring[i]
			var b = ring[j]

			yGreater := (a[1] > y) != (b[1] > y)
			xLess := x < (b[0]-a[0])*(y-a[1])/(b[1]-a[1])+a[0]

			if yGreater && xLess {
				inside = !inside
			}
			minDistSq = math.Min(minDistSq, getSegDistSq(x, y, a, b))

			j = i
		}
	}

	sqrt := math.Sqrt(minDistSq)
	if inside {
		return sqrt
	} else {
		return sqrt * -1
	}
}

// get polygon centroid
func getCentroidCell(polygon [][][]float64) cell {
	area := 0.0
	x := 0.0
	y := 0.0
	points := polygon[0]

	j := len(points) - 1
	for i := 0; i < len(points); i++ {
		a := points[i]
		b := points[j]
		f := a[0]*b[1] - b[0]*a[1]
		x += (a[0] + b[0]) * f
		y += (a[1] + b[1]) * f
		area += f * 3
		j = i
	}
	if area == 0 {
		return newCell(points[0][0], points[0][1], 0, polygon)
	} else {
		return newCell(x/area, y/area, 0, polygon)
	}
}

// get squared distance from a point to a segment
func getSegDistSq(px, py float64, a, b []float64) float64 {

	x := a[0]
	y := a[1]
	dx := b[0] - x
	dy := b[1] - y

	if dx != 0 || dy != 0 {

		t := ((px-x)*dx + (py-y)*dy) / (dx*dx + dy*dy)

		if t > 1 {
			x = b[0]
			y = b[1]

		} else if t > 0 {
			x += dx * t
			y += dy * t
		}
	}

	dx = px - x
	dy = py - y

	return dx*dx + dy*dy
}
