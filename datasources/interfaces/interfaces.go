package interfaces

import "bitbucket.org/heindl/process/datasources"

// Description is a shared provider for a taxon description
type Description interface {
	Citation() (string, error)
	Text() (string, error)
	Source() datasources.SourceType
}
