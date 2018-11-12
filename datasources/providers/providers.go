package providers

import (
	"github.com/heindl/floracast_process/datasources"
	"time"
)

// Description is a shared provider for a taxon description
type Description interface {
	Citation() (string, error)
	Text() (string, error)
	SourceType() datasources.SourceType
}

// Photo is a shared provider for taxon photos.
type Photo interface {
	Citation() string
	Thumbnail() string
	Large() string
	SourceType() datasources.SourceType
}

// Occurrences is the standard interface for sources that fetch occurrences.
type Occurrence interface {
	Lat() (float64, error)
	Lng() (float64, error)
	DateString() string
	CoordinatesEstimated() bool
	SourceOccurrenceID() string
	Classes() ([]string, error)
	Confidence() float64
}

type Bounds struct {
	North float64
	South float64
	East  float64
	West  float64
}

var DefaultBounds = Bounds{
	//var cutset_geometry = ee.Geometry.Rectangle({
	//coords: [-145.1767463, 24.5465169,-49.0, 59.5747563],
	//geodesic: false,
	//});
	North: 59.5747563,
	South: 24.5465169,
}

type OccurrenceFetchRequest struct {
	StartTime     *time.Time
	EndTime       *time.Time
	ParentTaxonID string
}
