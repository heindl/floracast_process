package nameusage

import (
	"bitbucket.org/heindl/process/store"
	"bitbucket.org/heindl/process/utils"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"github.com/golang/glog"
)

// Upload validates and saves new NameUsage to FireStore
func (Ω *usage) Upload(ctx context.Context, floraStore store.FloraStore) (deletedUsageIDs IDs, err error) {

	glog.Infof("Uploading NameUsage [%s]", Ω.CanonicalName())

	id, err := Ω.ID()
	if err != nil {
		return nil, err
	}

	col, err := floraStore.FirestoreCollection(store.CollectionNameUsages)
	if err != nil {
		return nil, err
	}

	docRef := col.Doc(id.String())

	deletedUsageIDs, err = Ω.matchInStore(ctx, floraStore)
	if err != nil {
		return nil, err
	}

	if err = clearStoreUsages(ctx, floraStore, deletedUsageIDs); err != nil {
		return nil, err
	}

	Ω.Occrrncs, err = Ω.Occurrences()
	if err != nil {
		return nil, err
	}

	Ω.SciNames = map[string]bool{}
	sciNames, err := Ω.AllScientificNames()
	if err != nil {
		return nil, err
	}

	for _, s := range sciNames {
		Ω.SciNames[s] = true
	}

	if _, err := docRef.Set(ctx, Ω); err != nil {
		return nil, err
	}

	glog.Infof("Completed Uploading NameUsage [%s]", Ω.CanonicalName())

	return deletedUsageIDs, nil
}

func clearStoreUsages(ctx context.Context, florastore store.FloraStore, allUsageIDs IDs) error {

	if len(allUsageIDs) == 0 {
		return nil
	}

	for _, usageIDs := range allUsageIDs.Batch(500) {
		if len(usageIDs) == 0 {
			return nil
		}
		batch := florastore.FirestoreBatch()
		for _, id := range usageIDs {
			col, err := florastore.FirestoreCollection(store.CollectionNameUsages)
			if err != nil {
				return err
			}
			batch = batch.Delete(col.Doc(id.String()))
		}
		if _, err := batch.Commit(ctx); err != nil {
			return errors.Wrap(err, "Could not commit materialized taxa")
		}
	}
	return nil
}

func (Ω *usage) matchInStore(ctx context.Context, florastore store.FloraStore) (IDs, error) {

	names, err := Ω.AllScientificNames()
	if err != nil {
		return nil, err
	}

	col, err := florastore.FirestoreCollection(store.CollectionNameUsages)
	if err != nil {
		return nil, err
	}

	wait := store.NewFirestoreLimiter()
	list, err := utils.ForEachStringToStrings(names, func(name string) ([]string, error) {
		<-wait
		synonymMatch := fmt.Sprintf("%s.%s", storeKeyScientificName, name)
		snaps, err := col.Where(synonymMatch, "==", true).Documents(ctx).GetAll()
		if err != nil {
			return nil, err
		}
		res := []string{}
		for _, snap := range snaps {
			res = append(res, snap.Ref.ID)
		}
		return res, nil
	})
	if err != nil {
		return nil, err
	}

	res, err := IDsFromStrings(list)
	if err != nil {
		return nil, err
	}
	return res, err
}

// Iterator fits the Go standard interface for fetching NameUsages.
type Iterator struct {
	error    error
	iterator *firestore.DocumentIterator
}

// Next returns the next NameUsage.
func (i *Iterator) Next() (NameUsage, error) {
	if i.error != nil {
		return nil, i.error
	}
	snap, err := i.iterator.Next()
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(snap.Data())
	if err != nil {
		return nil, errors.Wrap(err, "Could not Marshal NameUsage")
	}
	usage, err := FromJSON(ID(snap.Ref.ID), b)
	if err != nil {
		return nil, err
	}
	return usage, nil
}

// FetchAll returns an iterator for all NameUsages in FireStore.
func FetchAll(ctx context.Context, floraStore store.FloraStore) *Iterator {
	ref, err := floraStore.FirestoreCollection(store.CollectionNameUsages)
	return &Iterator{
		error:    err,
		iterator: ref.Documents(ctx),
	}
}

//
//type storeKey string
//
//func (Ω storeKey) String() string {
//	return string(Ω)
//}

//const (
//	storeKeyCanonicalName   = storeKey("Name")
//	storeKeyScientificNames = storeKey("ScientificNames")
//	storeKeyOccurrences     = storeKey("Occurrences")
//	storeKeySources     = storeKey("Sources")
//	storeKeyModifiedAt     = storeKey("ModifiedAt")
//	storeKeyCreatedAt     = storeKey("CreatedAt")
//)
//
//func (Ω *usage) toMap() (map[storeKey]interface{}, error) {
//
//	synonymMap := map[string]bool{}
//	for _, s := range Ω.AllScientificNames() {
//		synonymMap[s] = true
//	}
//
//	m := map[storeKey]interface{}{
//		storeKeyCanonicalName:   Ω.Name().ScientificName(),
//		storeKeyScientificNames: synonymMap,
//		storeKeyOccurrences:     Ω.TotalOccurrenceCount(),
//		storeKeySources:               Ω.sources,
//		storeKeyModifiedAt: time.Now(),
//		storeKeyCreatedAt: Ω.createdAt,
//	}
//
//	return m, nil
//}
//
//func fromMap(æ map[string]interface{}) (*NameUsage, error) {
//
//	u := &NameUsage{
//		id:
//	}
//
//	m := map[storeKey]interface{}{
//		storeKeyCanonicalName:   Ω.Name().ScientificName(),
//		storeKeyScientificNames: synonymMap,
//		storeKeyOccurrences:     Ω.TotalOccurrenceCount(),
//		storeKeySources:               Ω.sources,
//		storeKeyModifiedAt: time.Now(),
//		storeKeyCreatedAt: Ω.createdAt,
//	}
//
//	var stuff map[string]string
//	err := json.Unmarshal(b, &stuff)
//	if err != nil {
//		return err
//	}
//	for key, value := range stuff {
//		numericKey, err := strconv.Atoi(key)
//		if err != nil {
//			return err
//		}
//		this[numericKey] = value
//	}
//	return nil
//}
//
//func (Ω *NameUsage) MarshalJSON() ([]byte, error) {
//	if Ω == nil {
//		return nil, nil
//	}
//	m, err := Ω.toMap()
//	if err != nil {
//		return nil, err
//	}
//	return json.Marshal(m)
//}
