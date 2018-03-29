package main

import (
	"bitbucket.org/heindl/process/store"
	"cloud.google.com/go/firestore"
	"golang.org/x/net/context"
)

func migrateOccurrenceToNameUsageMonth(floraStore store.FloraStore) error {

	ctx := context.Background()

	ref, err := floraStore.FirestoreCollection(store.CollectionOccurrences)
	if err != nil {
		return err
	}
	snaps, err := ref.Documents(ctx).GetAll()
	if err != nil {
		return err
	}
	for _, snap := range snaps {
		if _, err := snap.Ref.Update(ctx, []firestore.Update{{Path: "NameUsageID", Value: "9sYKdRe6OUgzTwabsjjuFiwVU"}}); err != nil {
			return err
		}
	}
	return nil
}
