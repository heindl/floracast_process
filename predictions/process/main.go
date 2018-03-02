package main

import (
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/predictions/cache"
	"bitbucket.org/heindl/process/predictions/parser"
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
	requestedUsageIDs := flag.String("usageIDs", "", "Comma seperated list of NameUsageIDs to fetch predictions for.")
	bucket := flag.String("bucket", "", "gcs bucket to fetch predictions from")
	mode := flag.String("mode", "serve", "mode to handle predictions: write to temp file for javascript geofire uploader or serve for testing in local web router.")
	flag.Parse()

	if *requestedUsageIDs == "" {
		panic("NameUsageIDs required")
	}

	cxt := context.Background()

	parsedUsageIDs, err := nameusage.NameUsageIDsFromStrings(strings.Split(*requestedUsageIDs, ","))
	if err != nil {
		panic(err)
	}

	floraStore, err := store.NewFloraStore(cxt)
	if err != nil {
		panic(err)
	}

	h, err := newHandler(cxt, floraStore, *bucket, *mode)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := h.close(); err != nil {
			panic(err)
		}
	}()

	predictions, err := h.prsr.FetchPredictions(cxt, parsedUsageIDs, nil)
	if err != nil {
		panic(err)
	}

	if err := predictions.Upload(cxt, floraStore, h.cache); err != nil {
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

func newHandler(ctx context.Context, floraStore store.FloraStore, bucket, mode string) (res *handler, err error) {

	res = &handler{}
	if bucket == "" {
		res.source, err = parser.NewLocalPredictionSource(ctx, "/tmp/floracast-datamining/")
		if err != nil {
			return nil, err
		}
	} else {
		res.source, err = parser.NewGCSPredictionSource(ctx, floraStore)
		if err != nil {
			return nil, err
		}
	}

	res.prsr, err = parser.NewPredictionParser(res.source)
	if err != nil {
		return nil, err
	}

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
	cache  cache.PredictionCache
	prsr   parser.PredictionParser
	source parser.PredictionSource
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

	var usageID nameusage.NameUsageID
	if _, ok := vars["nameUsageID"]; ok {
		usageID = nameusage.NameUsageID(vars["nameUsageID"])
		if !usageID.Valid() {
			http.Error(w, "Invalid NameUsageID", http.StatusBadRequest)
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
