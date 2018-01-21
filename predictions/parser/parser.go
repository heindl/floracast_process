package parser

import (
	"bitbucket.org/heindl/taxa/store"
	"bitbucket.org/heindl/taxa/utils"
	"context"
	"github.com/montanaflynn/stats"
	"github.com/saleswise/errors/errors"
	"google.golang.org/genproto/googleapis/type/latlng"
	"gopkg.in/tomb.v2"
	"strings"
	"sync"
	"time"
)

type Writer interface {
	WritePredictionLine(p store.Prediction) error
	Close() error
}

type PredictionParser interface {
	FetchPredictions(cxt context.Context, taxa []string, date []string) ([]*store.Prediction, error)
}

func NewPredictionParser(cxt context.Context, gcsBucketName string, localPath string) (PredictionParser, error) {

	taxastore, err := store.NewTaxaStore()
	if err != nil {
		return nil, err
	}

	gcsFetcher, err := NewGCSFetcher(cxt, gcsBucketName, localPath)
	if err != nil {
		return nil, err
	}

	return &predictionParser{
		WildernessAreaFetcher: NewWildernessAreaFetcher(taxastore),
		GCSFetcher:            gcsFetcher,
	}, nil
}

type predictionParser struct {
	WildernessAreaFetcher WildernessAreaFetcher
	GCSFetcher            GCSFetcher
}

func taxonFromPredictionFilePath(p string) store.INaturalistTaxonID {
	a := strings.Split(p, "/")
	for i, v := range a {
		if v == "predictions" {
			return store.INaturalistTaxonID(a[i+1])
		}
	}
	return store.INaturalistTaxonID("")
}

type PredictionAggregator struct {
	PredictionObjects       []*store.Prediction
	PredictionList          map[store.INaturalistTaxonID][]float64
	TotalProtectedAreaCount map[store.INaturalistTaxonID]float64
	sync.Mutex
}

func (Ω PredictionAggregator) calcTaxonScarcity() (map[store.INaturalistTaxonID]float64, error) {
	taxaRatios := stats.Float64Data{}
	taxaRatiosMap := map[store.INaturalistTaxonID]float64{}
	for taxon, predictionValues := range Ω.PredictionList {
		totalTaxonPredictionCount := float64(len(predictionValues))
		totalTaxonProtectedAreaCount := Ω.TotalProtectedAreaCount[taxon]
		taxaRatios = append(taxaRatios, totalTaxonPredictionCount/totalTaxonProtectedAreaCount)
		taxaRatiosMap[taxon] = totalTaxonPredictionCount / totalTaxonProtectedAreaCount
	}

	taxaRatioMean, err := stats.Mean(taxaRatios)
	if err != nil {
		return nil, errors.Wrap(err, "could not calculate mean")
	}

	var taxonRatioInvertedMin, taxonRatioInvertedMax float64
	// In order to scale, must calculate the min and the max values once we invert the value by subtracting the mean.
	// This is so rarer taxa have a higher intensity value.
	for taxon, ratio := range taxaRatiosMap {
		invertedValue := taxaRatioMean - ratio
		if invertedValue < taxonRatioInvertedMin || taxonRatioInvertedMin == 0 {
			taxonRatioInvertedMin = invertedValue
		}
		if invertedValue > taxonRatioInvertedMax || taxonRatioInvertedMax == 0 {
			taxonRatioInvertedMax = invertedValue
		}
		taxaRatiosMap[taxon] = invertedValue
	}
	for taxon, ratio := range taxaRatiosMap {
		// Scale between 1 and 0.5
		//((b-a)(x - min) / max - min) + a
		if ratio != 0 {
			taxaRatiosMap[taxon] = ((1 - 0.5) * (ratio - taxonRatioInvertedMin) / (taxonRatioInvertedMax - taxonRatioInvertedMin)) + 0.5
		}
	}

	return taxaRatiosMap, nil
}

func (Ω *predictionParser) FetchPredictions(cxt context.Context, taxa []string, dates []string) ([]*store.Prediction, error) {

	aggr := PredictionAggregator{
		PredictionObjects:       []*store.Prediction{},
		PredictionList:          make(map[store.INaturalistTaxonID][]float64),
		TotalProtectedAreaCount: make(map[store.INaturalistTaxonID]float64),
	}

	gcsFilePaths := []string{}
	for _, _taxonID := range taxa {
		taxonID := store.INaturalistTaxonID(_taxonID)
		aggr.PredictionList[taxonID] = []float64{}
		aggr.TotalProtectedAreaCount[taxonID] = 0
		gcsPaths, err := Ω.GCSFetcher.FetchLatestPredictionFileNames(cxt, taxonID, "*")
		if err != nil {
			return nil, err
		}
		gcsFilePaths = append(gcsFilePaths, gcsPaths...)
	}

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, _fpath := range gcsFilePaths {
			fpath := _fpath
			tmb.Go(func() error {
				predictions, err := Ω.GCSFetcher.FetchPredictions(cxt, fpath)
				if err != nil {
					return err
				}
				for _, _p := range predictions {
					p := _p
					tmb.Go(func() error {
						aggr.Lock()
						defer aggr.Unlock()
						aggr.TotalProtectedAreaCount[p.Taxon] += 1
						if p.Target <= p.Random {
							return nil
						}
						aggr.PredictionList[p.Taxon] = append(aggr.PredictionList[p.Taxon], p.Target)
						d, err := time.ParseInLocation("20060102", p.Date, time.UTC)
						if err != nil {
							return errors.Wrap(err, "could not parse date")
						}
						aggr.PredictionObjects = append(aggr.PredictionObjects, &store.Prediction{
							CreatedAt:       utils.TimePtr(time.Now()),
							Location:        latlng.LatLng{p.Latitude, p.Longitude},
							PredictionValue: p.Target,
							TaxonID:         p.Taxon,
							Date:            utils.TimePtr(d),
							FormattedDate:   p.Date,
							Month:           d.Month(),
							//WildernessAreaID: wa.ID,
							//WildernessAreaName: wa.Name,
						})
						return nil
					})

				}
				return nil
			})
		}
		return nil
	})
	if err := tmb.Wait(); err != nil {
		return nil, err
	}

	taxaScarcityMap, err := aggr.calcTaxonScarcity()
	if err != nil {
		return nil, err
	}

	lmtr := utils.NewLimiter(20)
	tmb = tomb.Tomb{}
	tmb.Go(func() error {
		for _, _predictionObject := range aggr.PredictionObjects {
			predictionObject := _predictionObject
			tmb.Go(func() error {
				done := lmtr.Go()
				defer done()

				wa, err := Ω.WildernessAreaFetcher.GetWildernessArea(cxt, predictionObject.Location.Latitude, predictionObject.Location.Longitude)
				if err != nil {
					return err
				}
				if wa.ID != "" {
					predictionObject.WildernessAreaID = wa.ID
					predictionObject.WildernessAreaName = wa.Name
				}
				predictionObject.ScaledPredictionValue = (predictionObject.PredictionValue - 0.5) / 0.5
				predictionObject.ScarcityValue = taxaScarcityMap[predictionObject.TaxonID] * predictionObject.ScaledPredictionValue
				return nil
			})
		}
		return nil
	})
	if err := tmb.Wait(); err != nil {
		return nil, err
	}

	return aggr.PredictionObjects, nil
}
