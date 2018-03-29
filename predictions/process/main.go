package main

import (
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/predictions/cache"
	"bitbucket.org/heindl/process/predictions/generate"
	"bitbucket.org/heindl/process/store"
	"context"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	//writeToCache := flag.Bool("cache", false, "write to buntdb cache and initiate server?")
	//dates := flag.String("dates", "", "Dates for which to fetch latest predictions in format YYYYMMDD,YYYYMMDD. If blank will fetch all dates.")
	nameUsageIDPtr := flag.String("name_usage", "", "NameUsageID")
	mode := flag.String("mode", "serve", "mode to handle predictions: write to temp file for javascript geofire uploader or serve for testing in local web router.")
	flag.Parse()

	if *nameUsageIDPtr == "" {
		panic("NameUsageID required")
	}

	cxt := context.Background()

	nameUsageID := nameusage.ID(*nameUsageIDPtr)

	floraStore, err := store.NewFloraStore(cxt)
	if err != nil {
		panic(err)
	}

	h, err := newHandler(*mode)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := h.close(); err != nil {
			panic(err)
		}
	}()

	list, err := generate.GeneratePredictions(cxt, nameUsageID, floraStore, nil)
	if err != nil {
		panic(err)
	}

	if err := list.Upload(cxt, h.cache); err != nil {
		panic(err)
	}

	if *mode != "serve" {
		return
	}

	router := mux.NewRouter()

	router.HandleFunc("/{taxon}/{location}", h.handleRequest)

	fmt.Println("Server Ready at http://localhost:8081")

	if err := http.ListenAndServe(":8081", router); err != nil {
		panic(err)
	}

}

func newHandler(mode string) (res *handler, err error) {

	res = &handler{}

	switch mode {
	case "write":
		res.cache, _, err = cache.NewLocalFileCache()
		if err != nil {
			return nil, err
		}
	case "serve":
		res.cache, _, err = cache.NewLocalGeoCache()
		if err != nil {
			return nil, err
		}
	default:
		panic("Expected mode to be write or serve")
	}

	return res, nil
}

type handler struct {
	cache cache.PredictionCache
}

func (Ω *handler) close() error {
	return Ω.cache.Close()
}

func (Ω *handler) handleRequest(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)

	date := r.URL.Query().Get("date")

	fmt.Println("recieving request", vars["taxon"], vars["location"], date)

	loc := strings.Split(vars["location"], ",")

	lat, err := strconv.ParseFloat(loc[0], 64)
	if err != nil {
		http.Error(w, "Invalid latitude", http.StatusBadRequest)
		return
	}
	lng, err := strconv.ParseFloat(loc[1], 64)
	if err != nil {
		http.Error(w, "Invalid longitude", http.StatusBadRequest)
		return
	}
	rad, err := strconv.ParseFloat(loc[2], 64)
	if err != nil {
		http.Error(w, "Invalid radius", http.StatusBadRequest)
		return
	}

	var usageID nameusage.ID
	if _, ok := vars["nameUsageID"]; ok {
		usageID = nameusage.ID(vars["nameUsageID"])
		if !usageID.Valid() {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
	}

	l, err := Ω.cache.ReadPredictions(lat, lng, rad, date, &usageID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, strings.Join(l, "\n"))

}
