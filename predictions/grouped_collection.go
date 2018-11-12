package predictions

import (
	"github.com/heindl/floracast_process/nameusage/nameusage"
	"github.com/heindl/floracast_process/protectedarea/cache"
	"github.com/heindl/floracast_process/store"
	"github.com/heindl/floracast_process/terra/geo"
	"github.com/heindl/floracast_process/terra/geoembed"
	"github.com/heindl/floracast_process/utils"
	"context"
	"fmt"
	"github.com/saleswise/errors/errors"
	"gopkg.in/tomb.v2"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type GroupedCollection interface {
	Add(nameUsageId nameusage.ID, lat, lng float64, date utils.FormattedDate, prediction float64) error
	Count() (int, int)
	Print() error
	Upload(cxt context.Context) error
}

func NewGroupedCollection(floraStore store.FloraStore) (GroupedCollection, error) {

	protectedAreaCache, err := cache.NewProtectedAreaCache(floraStore)
	if err != nil {
		return nil, err
	}

	return &groupedCollection{
		floraStore:   floraStore,
		areaCache:    protectedAreaCache,
		records:      map[string]map[nameusage.ID][]*groupRecord{},
		missingAreas: []string{},
	}, nil
}

type groupRecord struct {
	AreaKilometers float64 `firestore:"km" json:"km"`
	Latitude       float64 `firestore:"lat" json:"lat"`
	Longitude      float64 `firestore:"lng" json:"lng"`
	Prediction     float64 `firestore:"Ω" json:"Ω"`
	SortValue      float64
}

type ByPrediction []*groupRecord

func (a ByPrediction) Len() int           { return len(a) }
func (a ByPrediction) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPrediction) Less(i, j int) bool { return a[i].SortValue > a[j].SortValue }

type groupedCollection struct {
	areaCache    cache.ProtectedAreaCache
	floraStore   store.FloraStore
	nameUsageIDs []nameusage.ID
	records      map[string]map[nameusage.ID][]*groupRecord // Token-Date
	total        int
	missingAreas []string
	sync.Mutex
}

func (Ω *groupedCollection) Print() error {
	//b, err := json.Marshal(Ω.Docs())
	//if err != nil {
	//	return err
	//}
	//
	//r := map[string]interface{}{}
	//if err := bson.Unmarshal(b, &r); err != nil {
	//	return err
	//}
	fmt.Println(utils.JsonOrSpew(Ω.Docs()))
	return nil
}

func (Ω *groupedCollection) Count() (int, int) {
	return Ω.total, len(Ω.records)
}

type fireStorePredictionDoc struct {
	Token string                          `firestore:"Token" json:"Token"`
	Date  string                          `firestore:"Date" json:"Date"`
	Taxa  map[nameusage.ID][]*groupRecord `firestore:"Taxa" json:"Taxa"`
}

func (Ω *groupedCollection) Docs() []fireStorePredictionDoc {

	res := []fireStorePredictionDoc{}

	for tokenDateKey, nameUsageMap := range Ω.records {
		doc := fireStorePredictionDoc{
			Token: strings.Split(tokenDateKey, "-")[0],
			Date:  strings.Split(tokenDateKey, "-")[1],
			Taxa:  map[nameusage.ID][]*groupRecord{},
		}

		for nameUsageId, records := range nameUsageMap {

			if len(records) == 0 {
				continue
			}

			sort.Sort(ByPrediction(records))

			if len(records) > 3 {
				records = records[:3]
			}

			doc.Taxa[nameUsageId] = records
		}

		res = append(res, doc)
	}

	return res
}

func (Ω *groupedCollection) Upload(ctx context.Context) error {

	if len(Ω.records) == 0 {
		return nil
	}

	//glog.Infof("Uploading Prediction Documents [%d] for NameUsage [%s]", len(Ω.records), Ω.nameUsageID)

	col, err := Ω.floraStore.FirestoreCollection(store.CollectionPredictionIndex)
	if err != nil {
		return err
	}

	//startedAt := time.Now().UnixNano()

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		i := float64(0)
		fireStoreBatch := Ω.floraStore.FirestoreBatch()
		for _, doc := range Ω.Docs() {

			docRef := col.Doc(doc.Token + "-" + doc.Date)

			fireStoreBatch = fireStoreBatch.Set(docRef, doc)
			if i != 0 && math.Mod(i, 400) == 0 {
				_fireStoreBatch := *fireStoreBatch
				fireStoreBatch = Ω.floraStore.FirestoreBatch()
				tmb.Go(func() error {
					if _, err := _fireStoreBatch.Commit(ctx); err != nil {
						return errors.Wrapf(err, "Could not commit set batch for Predictions [%f]", i)
					}
					return nil
				})
			}
			i++
		}
		// Commit what is left
		if _, err := fireStoreBatch.Commit(ctx); err != nil {
			return errors.Wrapf(err, "Could not commit set batch for Predictions [%f]", i)
		}
		return nil
	})

	if err := tmb.Wait(); err != nil {
		return err
	}

	//docRefsToDeleted, err := col.
	//	Where("ModifiedAt", "<", startedAt).
	//	Where("NameUsageID", "==", Ω.nameUsageID).Documents(ctx).GetAll()
	//
	//if err != nil {
	//	return errors.Wrap(err, "Could not get features to delete")
	//}
	//
	//if len(docRefsToDeleted) == 0 {
	//	return nil
	//}
	//
	//glog.Infof("Deleting stale Prediction Documents [%d] for NameUsage [%s]", len(docRefsToDeleted), Ω.nameUsageID)
	//
	//deleteBatchCount := float64(0)
	//fireStoreBatch := Ω.floraStore.FirestoreBatch()
	//for _, doc := range docRefsToDeleted {
	//	fireStoreBatch = fireStoreBatch.Delete(doc.Ref)
	//	if math.Mod(deleteBatchCount, 400) == 0 {
	//		if _, err := fireStoreBatch.Commit(ctx); err != nil {
	//			return errors.Wrapf(err, "Could not delete [%d] Prediction documents", len(docRefsToDeleted))
	//		}
	//		fireStoreBatch = Ω.floraStore.FirestoreBatch()
	//	}
	//	deleteBatchCount++
	//}
	//if _, err := fireStoreBatch.Commit(ctx); err != nil {
	//	return errors.Wrapf(err, "Could not delete [%d] Prediction documents", len(docRefsToDeleted))
	//}

	return nil
}

func (Ω *groupedCollection) Add(nameUsageId nameusage.ID, lat, lng float64, date utils.FormattedDate, prediction float64) error {

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

	point, err := geo.NewPoint(lat, lng)
	if err != nil {
		return err
	}

	r := &groupRecord{
		Latitude:   lat,
		Longitude:  lng,
		Prediction: prediction,
	}

	if km, ok := cache.SquareKilometers[s2Key]; ok {
		r.AreaKilometers = km
	}

	r.SortValue = (r.Prediction * 50) * (r.AreaKilometers * 0.25)

	tokenMap := point.S2TokenMap()

	Ω.Lock()
	defer Ω.Unlock()

	if r.AreaKilometers == 0 {
		Ω.missingAreas = append(Ω.missingAreas, string(s2Key))
	}

	for k, v := range tokenMap {
		i, err := strconv.Atoi(k)
		if err != nil {
			return errors.Wrap(err, "Could not parse token map key as int")
		}
		if i >= 10 {
			continue
		}

		key := fmt.Sprintf("%s-%s", v, date)

		if _, ok := Ω.records[key]; !ok {
			Ω.records[key] = map[nameusage.ID][]*groupRecord{}
		}

		if _, ok := Ω.records[key][nameUsageId]; !ok {
			Ω.records[key][nameUsageId] = []*groupRecord{}
		}

		Ω.records[key][nameUsageId] = append(Ω.records[key][nameUsageId], r)
	}

	Ω.total += 1

	return nil

}
