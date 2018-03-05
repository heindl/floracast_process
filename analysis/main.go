package main

import (
	"bitbucket.org/heindl/process/store"
	"fmt"
	"golang.org/x/net/context"
)

// https://github.com/galeone/tfgo
func main() {
	cxt := context.Background()

	floraStore, err := store.NewFloraStore(cxt)
	if err != nil {
		panic(err)
	}

	s := "gs://floracast-datamining/random/1510706694"
	//s := "gs://floracast-datamining/protected_areas/20180302/1520210026.tfrecord.gz"

	total, err := countTFRecordsInCloudStoragePath(cxt, floraStore, s)
	if err != nil {
		panic(err)
	}

	fmt.Println("TOTAL", total)
}
