package providers

import "bitbucket.org/heindl/process/datasources"

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
}
