package nameusage

import (
	"bitbucket.org/heindl/processors/store"
	"context"
	"fmt"
	"bitbucket.org/heindl/processors/utils"
	"github.com/dropbox/godropbox/errors"
	"encoding/json"
	"time"
)

func (Ω *NameUsage) Upload(ctx context.Context, florastore store.FloraStore) (deletedUsageIDs NameUsageIDs, err error) {

	docRef := florastore.FirestoreCollection(store.CollectionNameUsages).Doc(Ω.ID().String())

	deletedUsageIDs, err = Ω.matchInStore(ctx, florastore)
	if err != nil {
		return nil, err
	}

	if err := clearStoreUsages(ctx, florastore, deletedUsageIDs); err != nil {
		return nil, err
	}

	m, err := Ω.toMap()
	if err != nil {
		return nil, err
	}
	if _, err := docRef.Set(ctx, m); err != nil {
		return nil, err
	}

	return deletedUsageIDs, nil

}

func clearStoreUsages(ctx context.Context, florastore store.FloraStore, allUsageIDs NameUsageIDs) error {
	for _, usageIDs := range allUsageIDs.Batch(500) {
		if len(usageIDs) == 0 {
			return nil
		}
		batch := florastore.FirestoreBatch()
		for _, id := range usageIDs {
			docRef := florastore.FirestoreCollection(store.CollectionNameUsages).Doc(id.String())
			batch = batch.Delete(docRef)
		}
		if _, err := batch.Commit(ctx); err != nil {
			return errors.Wrap(err, "Could not commit materialized taxa")
		}
	}
	return nil
}

func (Ω *NameUsage) matchInStore(ctx context.Context, florastore store.FloraStore) (NameUsageIDs, error) {

	names := Ω.AllScientificNames()

	col := florastore.FirestoreCollection(store.CollectionNameUsages)

	wait := store.NewFirestoreLimiter()
	list, err := utils.ForEachStringToStrings(names, func(name string) ([]string, error){
		<-wait
		synonymMatch := fmt.Sprintf("%s.%s", storeKeyScientificNames.String())
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

	res, err := nameUsageIDsFromStrings(list)
	if err != nil {
		return nil, err
	}
	return res, err
}


type storeKey string

func (Ω storeKey) String() string {
	return string(Ω)
}

const (
	storeKeyCanonicalName   = storeKey("CanonicalName")
	storeKeyScientificNames = storeKey("ScientificNames")
	storeKeyOccurrences     = storeKey("Occurrences")
	storeKeySources     = storeKey("Sources")
	storeKeyModifiedAt     = storeKey("ModifiedAt")
	storeKeyCreatedAt     = storeKey("CreatedAt")
)

func (Ω *NameUsage) toMap() (map[storeKey]interface{}, error) {

	synonymMap := map[string]bool{}
	for _, s := range Ω.AllScientificNames() {
		synonymMap[s] = true
	}

	m := map[storeKey]interface{}{
		storeKeyCanonicalName:   Ω.CanonicalName().ScientificName(),
		storeKeyScientificNames: synonymMap,
		storeKeyOccurrences:     Ω.TotalOccurrenceCount(),
		storeKeySources:               Ω.sources,
		storeKeyModifiedAt: time.Now(),
		storeKeyCreatedAt: Ω.createdAt,
	}

	return m, nil
}

func (Ω *NameUsage) MarshalJSON() ([]byte, error) {
	if Ω == nil {
		return nil, nil
	}
	m, err := Ω.toMap()
	if err != nil {
		return nil, err
	}
	return json.Marshal(m)
}
