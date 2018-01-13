package main

import (
	"bitbucket.org/heindl/taxa/predictions/filecache"
	"bitbucket.org/heindl/taxa/predictions/geocache"
	"bitbucket.org/heindl/taxa/predictions/parser"
	"bitbucket.org/heindl/taxa/store"
	"context"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

const predictionUploadLimit = 2000

type CacheWriter interface {
	WritePredictionLine(p *store.Prediction) error
	ReadTaxa(lat, lng, radius float64, qDate string, taxon string) ([]string, error)
	Close() error
}

func main() {

	var err error
	//writeToCache := flag.Bool("cache", false, "write to buntdb cache and initiate server?")
	dates := flag.String("dates", "", "Dates for which to fetch latest predictions in format YYYYMMDD,YYYYMMDD. If blank will fetch all dates.")
	taxa := flag.String("taxa", "", "Comma seperated list of taxa to fetch predictions for.")
	bucket := flag.String("bucket", "", "gcs bucket to fetch predictions from")
	mode := flag.String("mode", "serve", "mode to handle predictions: write to temp file for javascript geofire uploader or serve for testing in local web router.")
	flag.Parse()

	if *taxa == "" {
		panic("taxa required")
	}

	cxt := context.Background()
	predictionParser, err := parser.NewPredictionParser(cxt, *bucket, "/tmp")
	if err != nil {
		panic(err)
	}

	date_list := strings.Split(*dates, ",")
	if len(date_list) == 0 {
		date_list = append(date_list, "") // Add an empty value to make iteration simpler.
	}

	var cache CacheWriter

	switch *mode {
	case "write":
		cache, err = filecache.NewFileCache()
		if err != nil {
			panic(err)
		}
	case "serve":
		cache, err = geocache.NewCacheWriter(strings.Split(*taxa, ","))
		if err != nil {
			panic(err)
		}
	}
	defer cache.Close()

	predictions, err := predictionParser.FetchPredictions(cxt, strings.Split(*taxa, ","), nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Have Predictions")

	for _, p := range predictions {
		if err := cache.WritePredictionLine(p); err != nil {
			panic(err)
		}
	}

	if *mode != "serve" {
		return
	}

	//tmb := tomb.Tomb{}
	//tmb.Go(func() error {
	//	for _, _taxon := range strings.Split(*taxa, ",") {
	//		taxon := _taxon
	//		for _, _date := range date_list {
	//			date := _date
	//			tmb.Go(func() error {
	//				return predictionParser.WritePredictions(cxt, store.TaxonID(taxon), date)
	//			})
	//		}
	//	}
	//	return nil
	//})
	//if err := tmb.Wait(); err != nil {
	//	panic(err)
	//}

	router := mux.NewRouter()

	router.HandleFunc("/{taxon}/{location}", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		vars := mux.Vars(r)

		txn := vars["taxon"]

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

		if txn == "taxa" {
			txn = ""
		}

		l, err := cache.ReadTaxa(lat, lng, rad, date, txn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, strings.Join(l, "\n"))
		return

	})

	fmt.Println("Server Ready at http://localhost:8081")

	if err := http.ListenAndServe(":8081", router); err != nil {
		panic(err)
	}

	// If we wrote to a geocache, hold it open as a web server.

}
