package main

import (
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/predictions/generate"
	"bitbucket.org/heindl/process/store"
	"context"
	"flag"
	"fmt"
)

func main() {
	//writeToCache := flag.Bool("cache", false, "write to buntdb cache and initiate server?")
	//dates := flag.String("dates", "", "Dates for which to fetch latest predictions in format YYYYMMDD,YYYYMMDD. If blank will fetch all dates.")
	nameUsageIDPtr := flag.String("name_usage", "", "NameUsageID")
	//mode := flag.String("mode", "serve", "mode to handle predictions: write to temp file for javascript geofire uploader or serve for testing in local web router.")
	//dbPtr := flag.String("db", "", "existing database to use")

	flag.Parse()

	cxt := context.Background()

	nameUsageID := nameusage.ID(*nameUsageIDPtr)

	if !nameUsageID.Valid() {
		panic("NameUsageID required")
	}

	floraStore, err := store.NewFloraStore(cxt)
	if err != nil {
		panic(err)
	}

	//h, err := newHandler(*mode, *dbPtr)
	//if err != nil {
	//	panic(err)
	//}
	//defer func() {
	//	if err := h.close(); err != nil {
	//		panic(err)
	//	}
	//}()
	//
	//if *dbPtr == "" {
	collection, err := generate.GeneratePredictions(cxt, nameUsageID, floraStore, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Predictions Generated: ", collection.Count())
	if err := collection.Upload(cxt); err != nil {
		panic(err)
	}
	//}

	//if *mode != "serve" {
	//	return
	//}

	//router := mux.NewRouter()
	//
	//router.HandleFunc("/{bboxString}", h.handleRequest)
	//
	//fmt.Println("Server Ready at http://localhost:8081")
	//
	//if err := http.ListenAndServe(":8081", router); err != nil {
	//	panic(err)
	//}

}

//func newHandler(mode string, cachePath string) (res *handler, err error) {
//
//	res = &handler{}
//
//	switch mode {
//	case "write":
//		res.cache, _, err = cache.NewLocalFileCache()
//		if err != nil {
//			return nil, err
//		}
//	case "serve":
//		res.cache, _, err = cache.NewLocalGeoCache(cachePath)
//		if err != nil {
//			return nil, err
//		}
//	default:
//		panic("Expected mode to be write or serve")
//	}
//
//	return res, nil
//}
//
//type handler struct {
//	cache cache.PredictionCache
//}
//
//func (Ω *handler) close() error {
//	return Ω.cache.Close()
//}
//
//func (Ω *handler) handleRequest(w http.ResponseWriter, r *http.Request) {
//
//	w.Header().Set("Content-Type", "text/html; charset=utf-8")
//	w.Header().Set("Access-Control-Allow-Origin", "*")
//
//	vars := mux.Vars(r)
//
//	l, err := Ω.cache.ReadPredictions(vars["bboxString"])
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	b, err := json.Marshal(l)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	fmt.Fprint(w, string(b))
//
//}
