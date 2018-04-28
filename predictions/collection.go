package predictions

import (
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/protectedarea/cache"
	"bitbucket.org/heindl/process/store"
	"bitbucket.org/heindl/process/terra/geo"
	"bitbucket.org/heindl/process/terra/geoembed"
	"bitbucket.org/heindl/process/utils"
	"context"
	"fmt"
	"github.com/golang/glog"
	"github.com/saleswise/errors/errors"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/tomb.v2"
	"math"
	"sync"
	"time"
)

type Collection interface {
	Add(lat, lng float64, date utils.FormattedDate, prediction float64) error
	Count() int
	Print() error
	Upload(cxt context.Context) error
}

func NewCollection(id nameusage.ID, floraStore store.FloraStore) (Collection, error) {
	if !id.Valid() {
		return nil, errors.New("Invalid NameUsageID")
	}

	protectedAreaCache, err := cache.NewProtectedAreaCache(floraStore)
	if err != nil {
		return nil, err
	}

	return &collection{
		floraStore:  floraStore,
		areaCache:   protectedAreaCache,
		nameUsageID: id,
		records:     map[geoembed.S2Key]*record{},
	}, nil
}

type collection struct {
	areaCache   cache.ProtectedAreaCache
	floraStore  store.FloraStore
	nameUsageID nameusage.ID
	records     map[geoembed.S2Key]*record
	total       int
	sync.Mutex
}

func (Ω *collection) Print() error {
	b, err := bson.Marshal(Ω.records)
	if err != nil {
		return err
	}

	r := map[string]interface{}{}
	if err := bson.Unmarshal(b, &r); err != nil {
		return err
	}
	fmt.Println(utils.JsonOrSpew(r))
	return nil
}

func (Ω *collection) Count() int {
	return Ω.total
}

type minimal struct {
	Exists bool    `firestore:"𝝨,omitempty" bson:"𝝨,omitempty"`
	Value  float64 `firestore:"Ω,omitempty" bson:"Ω,omitempty"`
}

type record struct {
	AreaKilometers float64                          `firestore:"AreaKilometers,omitempty" bson:"AreaKilometers,omitempty"`
	NameUsageID    nameusage.ID                     `firestore:"NameUsageID,omitempty" bson:"NameUsageID,omitempty"`
	ModifiedAt     int64                            `firestore:"ModifiedAt,omitempty" bson:"ModifiedAt,omitempty"`
	Timeline       map[utils.FormattedDate]*minimal `firestore:"Timeline,omitempty" bson:"Timeline,omitempty"`
	S2Tokens       map[string]string                `firestore:"S2Tokens,omitempty" bson:"S2Tokens,omitempty"`
}

func (Ω *collection) Upload(ctx context.Context) error {

	if len(Ω.records) == 0 {
		return nil
	}

	glog.Infof("Uploading Prediction Documents [%d] for NameUsage [%s]", len(Ω.records), Ω.nameUsageID)

	col, err := Ω.floraStore.FirestoreCollection(store.CollectionPredictionIndex)
	if err != nil {
		return err
	}

	startedAt := time.Now().UnixNano()

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		i := float64(0)
		fireStoreBatch := Ω.floraStore.FirestoreBatch()
		for s2Key, predictionRecord := range Ω.records {
			docRef := col.Doc(fmt.Sprintf("%s-%s", Ω.nameUsageID, s2Key))

			if km, ok := cache.SquareKilometers[s2Key]; ok {
				predictionRecord.AreaKilometers = km
			} else {
				area, err := Ω.areaCache.GetProtectedAreaWithToken(context.Background(), s2Key)
				if err != nil {
					return err
				}
				if area != nil {
					fmt.Println(fmt.Sprintf(`"%s": %f,`, s2Key, area.Kilometers()))
					predictionRecord.AreaKilometers = area.Kilometers()
				}
			}

			predictionRecord.ModifiedAt = time.Now().UnixNano()

			fireStoreBatch = fireStoreBatch.Set(docRef, predictionRecord)
			if i != 0 && math.Mod(i, 400) == 0 {
				_fireStoreBatch := *fireStoreBatch
				fireStoreBatch = Ω.floraStore.FirestoreBatch()
				tmb.Go(func() error {
					if _, err := _fireStoreBatch.Commit(ctx); err != nil {
						return errors.Wrapf(err, "Could not commit set batch for Predictions [%d]", i)
					}
					return nil
				})
			}
			i++
		}
		// Commit what is left
		if _, err := fireStoreBatch.Commit(ctx); err != nil {
			return errors.Wrapf(err, "Could not commit set batch for Predictions [%d]", i)
		}
		return nil
	})

	if err := tmb.Wait(); err != nil {
		return err
	}

	docRefsToDeleted, err := col.
		Where("ModifiedAt", "<", startedAt).
		Where("NameUsageID", "==", Ω.nameUsageID).Documents(ctx).GetAll()

	if err != nil {
		return errors.Wrap(err, "Could not get features to delete")
	}

	if len(docRefsToDeleted) == 0 {
		return nil
	}

	glog.Infof("Deleting stale Prediction Documents [%d] for NameUsage [%s]", len(docRefsToDeleted), Ω.nameUsageID)

	deleteBatchCount := float64(0)
	fireStoreBatch := Ω.floraStore.FirestoreBatch()
	for _, doc := range docRefsToDeleted {
		fireStoreBatch = fireStoreBatch.Delete(doc.Ref)
		if math.Mod(deleteBatchCount, 400) == 0 {
			if _, err := fireStoreBatch.Commit(ctx); err != nil {
				return errors.Wrapf(err, "Could not delete [%d] Prediction documents", len(docRefsToDeleted))
			}
			fireStoreBatch = Ω.floraStore.FirestoreBatch()
		}
		deleteBatchCount++
	}
	if _, err := fireStoreBatch.Commit(ctx); err != nil {
		return errors.Wrapf(err, "Could not delete [%d] Prediction documents", len(docRefsToDeleted))
	}

	return nil
}

func (Ω *collection) Add(lat, lng float64, date utils.FormattedDate, prediction float64) error {

	if prediction == 0 {
		return errors.New("Prediction required for GeoHashIndex Prediction Point")
	}

	if !date.Valid() {
		return errors.New("Prediction requires valid date")
	}

	s2Key, err := geoembed.NewS2Key(lat, lng)
	if err != nil {
		return err
	}

	Ω.Lock()
	defer Ω.Unlock()

	if _, ok := Ω.records[s2Key]; !ok {

		point, err := geo.NewPoint(lat, lng)
		if err != nil {
			return err
		}

		Ω.records[s2Key] = &record{
			NameUsageID: Ω.nameUsageID,
			Timeline:    map[utils.FormattedDate]*minimal{},
			S2Tokens:    point.S2TokenMap(),
		}
	}

	Ω.records[s2Key].Timeline[date] = &minimal{
		Exists: true,
		Value:  prediction,
	}

	Ω.total += 1

	return nil

}
