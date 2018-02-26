package nameusage

import (
	"bitbucket.org/heindl/process/store"
	"context"
	"fmt"
	"bitbucket.org/heindl/process/utils"
	"github.com/dropbox/godropbox/errors"
)

func (Ω *usage) Upload(ctx context.Context, florastore store.FloraStore) (deletedUsageIDs NameUsageIDs, err error) {

	id, err := Ω.ID()
	if err != nil {
		return nil, err
	}

	col, err := florastore.FirestoreCollection(store.CollectionNameUsages)
	if err != nil {
		return nil, err
	}

	docRef := col.Doc(id.String())

	deletedUsageIDs, err = Ω.matchInStore(ctx, florastore)
	if err != nil {
		return nil, err
	}

	if err := clearStoreUsages(ctx, florastore, deletedUsageIDs); err != nil {
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

	return deletedUsageIDs, nil
}

func clearStoreUsages(ctx context.Context, florastore store.FloraStore, allUsageIDs NameUsageIDs) error {
	
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

func (Ω *usage) matchInStore(ctx context.Context, florastore store.FloraStore) (NameUsageIDs, error) {

	names, err := Ω.AllScientificNames()
	if err != nil {
		return nil, err
	}

	col, err := florastore.FirestoreCollection(store.CollectionNameUsages)
	if err != nil {
		return nil, err
	}

	wait := store.NewFirestoreLimiter()
	list, err := utils.ForEachStringToStrings(names, func(name string) ([]string, error){
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

	res, err := NameUsageIDsFromStrings(list)
	if err != nil {
		return nil, err
	}
	return res, err
}

//
//type storeKey string
//
//func (Ω storeKey) String() string {
//	return string(Ω)
//}

//const (
//	storeKeyCanonicalName   = storeKey("CanonicalName")
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
//		storeKeyCanonicalName:   Ω.CanonicalName().ScientificName(),
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
//		storeKeyCanonicalName:   Ω.CanonicalName().ScientificName(),
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
