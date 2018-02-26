package protected_areas

import (
	"context"
	"bitbucket.org/heindl/process/store"
	"bitbucket.org/heindl/process/geofeatures"
	"github.com/dropbox/godropbox/errors"
)

var ErrNotFound = errors.New("Not Found")

func FetchOne(cxt context.Context, florastore store.FloraStore, coordinateKey geofeatures.CoordinateKey) (*ProtectedArea, error) {

	if !coordinateKey.Valid() {
		return nil, errors.Newf("Invalid CoordinateKey [%s]", coordinateKey)
	}

	col, err := florastore.FirestoreCollection(store.CollectionProtectedAreas)
	if err != nil {
		return nil, err
	}

	snap, err := col.Doc(string(coordinateKey)).Get(cxt)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not get ProtectedArea [%s]", coordinateKey)
	}
	w := ProtectedArea{}
	if err := snap.DataTo(&w); err != nil {
		return nil, errors.Wrap(err, "could not type cast ProtectedArea")
	}
	return &w, nil
}
