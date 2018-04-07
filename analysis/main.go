package main

import (
	"bitbucket.org/heindl/process/store"
	"bitbucket.org/heindl/process/utils"
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

	areas := protectedAreas{
		floraStore: floraStore,
		cxt:        context.Background(),
	}
	areas.PrintGeohashPrecisionGroups()

	return

	fmt.Println("Searching")
	col, err := floraStore.FirestoreCollection(store.CollectionGeoIndex)
	if err != nil {
		panic(err)
	}
	//
	snaps, err := col.Doc("dpp").Collection(string(store.CollectionPredictions)).
		Where("20180102.1.1.𝝨", ">", 0).
		Documents(cxt).GetAll()
	if err != nil {
		panic(err)
	}
	//
	for _, snap := range snaps {
		b, err := snap.DataAt("dates.20180102.3.m")
		if err != nil {
			panic(err)
		}
		fmt.Println(utils.JsonOrSpew(b))
	}
	//
	//for _, s := range snaps {
	//	if _, err := s.Ref.Delete(cxt); err != nil {
	//		panic(err)
	//	}
	//}
	//
	fmt.Println("Snaps", len(snaps))

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
