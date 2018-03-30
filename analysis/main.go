package main

import (
	"bitbucket.org/heindl/process/store"
	"context"
	"fmt"
)

// https://github.com/galeone/tfgo
func main() {
	cxt := context.Background()

	floraStore, err := store.NewFloraStore(cxt)
	if err != nil {
		panic(err)
	}

	colRef, err := floraStore.FirestoreCollection(store.CollectionRandom)
	if err != nil {
		panic(err)
	}

	snaps, err := colRef.Documents(cxt).GetAll()
	if err != nil {
		panic(err)
	}

	total := 0
	for _ = range snaps {
		total++
	}

	fmt.Println("Random Total", total)

	//if err := migrateOccurrenceToNameUsageMonth(floraStore); err != nil {
	//	panic(err)
	//}

	//s := "gs://floracast-datamining/occurrences/aho2iyxvo37rjezikho6xbwmq/1519504349.tfrecords"
	////s := "gs://floracast-datamining/random/1520448273.tfrecords"
	//
	//recorder, err := newTFRecorder(cxt, floraStore, s)
	//if err != nil {
	//	panic(err)
	//}
	//
	//count, err := recorder.CountRecords(cxt)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println("Count", count)
	//
	//if err := recorder.PrintFeature(cxt, "eco_num", featureTypeInt); err != nil {
	//	panic(err)
	//}
}
