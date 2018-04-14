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
	"github.com/saleswise/errors/errors"
	"gopkg.in/mgo.v2/bson"
	"math"
	"sync"
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
	Timeline       map[utils.FormattedDate]*minimal `firestore:"Timeline,omitempty" bson:"Timeline,omitempty"`
	Cells          map[string]string                `firestore:"Cells,omitempty" bson:"Cells,omitempty"`
}

func (Ω *collection) batches(maxBatchSize float64) []map[geoembed.S2Key]*record {

	if len(Ω.records) == 0 {
		return nil
	}

	batchCount := math.Ceil(float64(len(Ω.records)) / maxBatchSize)

	res := []map[geoembed.S2Key]*record{}

Outer:
	for {
		m := map[geoembed.S2Key]*record{}
		for k, v := range Ω.records {
			m[k] = v
			delete(Ω.records, k)
			if float64(len(m)) >= batchCount {
				res = append(res, m)
				continue Outer
			}
		}
		break
	}

	return res
}

func (Ω *collection) Upload(cxt context.Context) error {

	col, err := Ω.floraStore.FirestoreCollection(store.CollectionPredictionIndex)
	if err != nil {
		return err
	}

	for _, m := range Ω.batches(250) {
		fireStoreBatch := Ω.floraStore.FirestoreBatch()
		for k, v := range m {
			doc := col.Doc(fmt.Sprintf("%s-%s", Ω.nameUsageID, k))
			fireStoreBatch = fireStoreBatch.Set(doc, v)
		}
		if _, err := fireStoreBatch.Commit(cxt); err != nil {
			return err
		}
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

		area, err := Ω.areaCache.GetProtectedArea(context.Background(), lat, lng)
		if err != nil {
			return err
		}

		point, err := geo.NewPoint(lat, lng)
		if err != nil {
			return err
		}

		Ω.records[s2Key] = &record{
			NameUsageID:    Ω.nameUsageID,
			AreaKilometers: area.Kilometers(),
			Timeline:       map[utils.FormattedDate]*minimal{},
			Cells:          point.S2TokenMap(),
		}
	}

	Ω.records[s2Key].Timeline[date] = &minimal{
		Exists: true,
		Value:  prediction,
	}

	Ω.total += 1

	return nil

}
