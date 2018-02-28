package protectedarea

import (
	"bitbucket.org/heindl/process/store"
	"bitbucket.org/heindl/process/terra/geoembed"
	"context"
	"github.com/dropbox/godropbox/errors"
)

// FetchOne fetches a ProtectedArea from Cloud Firestore
func FetchOne(cxt context.Context, florastore store.FloraStore, coordinateKey geoembed.CoordinateKey) (ProtectedArea, error) {

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
	w := protectedArea{}
	if err := snap.DataTo(&w); err != nil {
		return nil, errors.Wrap(err, "could not type cast ProtectedArea")
	}
	return &w, nil
}
