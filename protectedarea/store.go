package protectedarea

import (
	"bitbucket.org/heindl/process/store"
	"bitbucket.org/heindl/process/terra/geoembed"
	"context"
	"encoding/json"
	"github.com/dropbox/godropbox/errors"
	"gopkg.in/tomb.v2"
	"math"
)

// FetchOne fetches a ProtectedArea from Cloud Firestore
func FetchOne(cxt context.Context, floraStore store.FloraStore, coordinateKey geoembed.CoordinateKey) (ProtectedArea, error) {

	if !coordinateKey.Valid() {
		return nil, errors.Newf("Invalid CoordinateKey [%s]", coordinateKey)
	}

	col, err := floraStore.FirestoreCollection(store.CollectionProtectedAreas)
	if err != nil {
		return nil, err
	}

	snap, err := col.Doc(string(coordinateKey)).Get(cxt)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not get ProtectedArea [%s]", coordinateKey)
	}

	b, err := json.Marshal(snap.Data())
	if err != nil {
		return nil, errors.Wrap(err, "Could not Marshal ProtectedArea FireStore data")
	}

	area := protectedArea{}
	if err := json.Unmarshal(b, &area); err != nil {
		return nil, errors.Wrap(err, "Could not Unmarshal ProtectedArea map")
	}

	return &area, nil
}

//
//type ProtectedAreas []ProtectedArea
//
//func (a ProtectedAreas) Len() int           { return len(a) }
//func (a ProtectedAreas) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
//func (a ProtectedAreas) Less(i, j int) bool { return a[i].GISAcres > a[j].GISAcres }

// ProtectedAreas is intended for bulk uploading.
type ProtectedAreas []ProtectedArea

func (Ω ProtectedAreas) batches(maxBatchSize float64) []ProtectedAreas {

	if len(Ω) == 0 {
		return nil
	}

	batchCount := math.Ceil(float64(len(Ω)) / maxBatchSize)

	res := []ProtectedAreas{}
	for i := 0.0; i <= batchCount-1; i++ {
		start := int(i * maxBatchSize)
		end := int(((i + 1) * maxBatchSize) - 1)
		if end > len(Ω) {
			end = len(Ω)
		}
		o := Ω[start:end]
		res = append(res, o)
	}

	return res
}

// Upload validates all ProtectedAreas and saves them to Firestore
func (Ω ProtectedAreas) Upload(cxt context.Context, floraStore store.FloraStore) (int, error) {

	batches := Ω.batches(500)

	col, err := floraStore.FirestoreCollection(store.CollectionProtectedAreas)
	if err != nil {
		return 0, err
	}

	tmb := tomb.Tomb{}
	invalidCount := 0
	tmb.Go(func() error {
		for _, 𝝨 := range batches {
			batch := 𝝨
			tmb.Go(func() error {
				fireStoreBatch := floraStore.FirestoreBatch()
				for _, area := range batch {
					if !area.Valid() {
						invalidCount++
						continue
					}
					id, err := area.ID()
					if err != nil {
						return err
					}

					b, err := json.Marshal(area)
					if err != nil {
						return err
					}

					m := map[string]interface{}{}
					if err := json.Unmarshal(b, &m); err != nil {
						return err
					}

					fireStoreBatch = fireStoreBatch.Set(col.Doc(string(id)), m)
				}
				if _, err := fireStoreBatch.Commit(cxt); err != nil {
					return errors.Wrap(err, "could not commit firestore batch")
				}
				return nil
			})
		}
		return nil
	})

	if err := tmb.Wait(); err != nil {
		return 0, err
	}

	return len(Ω) - invalidCount, nil
}
