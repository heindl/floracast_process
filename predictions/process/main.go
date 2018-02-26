package process

import (
	"bitbucket.org/heindl/process/predictions/parser"
	"context"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	"bitbucket.org/heindl/process/predictions/cache"
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/store"
)

const predictionUploadLimit = 2000

func main() {

	var err error
	//writeToCache := flag.Bool("cache", false, "write to buntdb cache and initiate server?")
	dates := flag.String("dates", "", "Dates for which to fetch latest predictions in format YYYYMMDD,YYYYMMDD. If blank will fetch all dates.")
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

	florastore, err := store.NewFloraStore(cxt)
	if err != nil {
		panic(err)
	}

	var src parser.PredictionSource
	if *bucket == "" {
		src, err = parser.NewLocalPredictionSource(cxt, "/tmp/floracast-datamining/")
		if err != nil {
			panic(err)
		}
	} else {
		src, err = parser.NewGCSPredictionSource(cxt, florastore)
		if err != nil {
			panic(err)
		}
	}

	predictionParser, err := parser.NewPredictionParser(src)
	if err != nil {
		panic(err)
	}

	date_list := strings.Split(*dates, ",")
	if len(date_list) == 0 {
		date_list = append(date_list, "") // AddUsage an empty value to make iteration simpler.
	}


	predictions, err := predictionParser.FetchPredictions(cxt, parsedUsageIDs, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Have Predictions")

	var predictionCache cache.PredictionCache

	switch *mode {
	case "write":
		predictionCache, _, err = cache.NewLocalFileCache()
		if err != nil {
			panic(err)
		}
	case "serve":
		predictionCache, _, err = cache.NewLocalGeoCache()
		if err != nil {
			panic(err)
		}
	default:
		panic("Expected mode to be write or serve")
	}
	defer func() {
		if err := predictionCache.Close(); err != nil {
			panic(err)
		}
	}()

	if err := predictions.Upload(cxt, florastore, predictionCache); err != nil {
		panic(err)
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

		l, err := predictionCache.ReadPredictions(lat, lng, rad, date, &usageID)
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
