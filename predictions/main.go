package main

import (
	"context"
	"bitbucket.org/heindl/taxa/store"
	"flag"
	"bitbucket.org/heindl/taxa/predictions/parser"
	"bitbucket.org/heindl/taxa/predictions/geocache"
	"strings"
	"gopkg.in/tomb.v2"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
)

const predictionUploadLimit = 2000

func main() {

	var err error
	//writeToCache := flag.Bool("cache", false, "write to buntdb cache and initiate server?")
	dates := flag.String("dates", "", "Dates for which to fetch latest predictions in format YYYYMMDD,YYYYMMDD. If blank will fetch all dates.")
	taxa := flag.String("taxa", "", "Comma seperated list of taxa to fetch predictions for.")
	bucket := flag.String("bucket", "", "gcs bucket to fetch predictions from")
	flag.Parse()

	if *taxa == "" {
		panic("taxa required")
	}

	geocacheWriter, err := geocache.NewCacheWriter(strings.Split(*taxa, ","))
	if err != nil {
		panic(err)
	}

	cxt := context.Background()
	parser, err := parser.NewPredictionParser(cxt, *bucket, geocacheWriter)

	date_list := strings.Split(*dates, ",")
	if len(date_list) == 0 {
		date_list = append(date_list, "") // Add an empty value to make iteration simpler.
	}

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, _taxon := range strings.Split(*taxa, ",") {
			taxon := _taxon
			for _, _date := range date_list {
				date := _date
				tmb.Go(func() error {
					return parser.FetchWritePredictions(cxt, store.TaxonID(taxon), date)
				})
			}
		}
		return nil
	})
	if err := tmb.Wait(); err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/{taxon}/{location}/", func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		txn := vars["taxon"]
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
			l, err := geocacheWriter.ReadTaxa(lat, lng, rad)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Fprint(w, strings.Join(l, "\n"))
			return
		} else {
			l, err := geocacheWriter.ReadTaxon(store.TaxonID(txn), lat, lng, rad)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Fprint(w, strings.Join(l, "\n"))
			return
		}

	})

	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}

	// If we wrote to a geocache, hold it open as a web server.


}

