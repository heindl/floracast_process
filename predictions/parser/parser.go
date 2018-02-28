package parser

import (
	"bufio"
	"context"
	"io"
	"strconv"
	"strings"
	"sync"

	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/predictions"
	"bitbucket.org/heindl/process/utils"
	"github.com/dropbox/godropbox/errors"
	"gopkg.in/tomb.v2"
)

type PredictionParser interface {
	FetchPredictions(cxt context.Context, nameUsageIDs nameusage.NameUsageIDs, date []string) (predictions.Predictions, error)
}

func NewPredictionParser(src PredictionSource) (PredictionParser, error) {
	return &predictionParser{
		predictionSource: src,
	}, nil
}

type predictionParser struct {
	predictionSource PredictionSource
}

func parsePredictionReader(id nameusage.NameUsageID, reader io.Reader) ([]*PredictionResult, error) {

	scanner := bufio.NewScanner(reader)

	scanner.Split(bufio.ScanLines)
	response_list := []*PredictionResult{}
	for scanner.Scan() {
		var err error
		line := strings.Split(scanner.Text(), ",")
		res := PredictionResult{
			Date:        line[2],
			NameUsageID: id,
		}
		res.Latitude, err = strconv.ParseFloat(line[0], 64)
		if err != nil {
			return nil, errors.Wrap(err, "could not parse latitude")
		}
		res.Longitude, err = strconv.ParseFloat(line[1], 64)
		if err != nil {
			return nil, errors.Wrap(err, "could not parse longitude")
		}
		res.Target, err = strconv.ParseFloat(line[3], 64)
		if err != nil {
			return nil, errors.Wrap(err, "could not parse target")
		}
		res.Random, err = strconv.ParseFloat(line[4], 64)
		if err != nil {
			return nil, errors.Wrap(err, "could not parse random")
		}
		response_list = append(response_list, &res)
	}
	return response_list, nil
}

func parseNameUsageIDFromFilePath(p string) (nameusage.NameUsageID, error) {
	a := strings.Split(p, "/")
	for i, v := range a {
		if v == "predictions" {
			id := nameusage.NameUsageID(a[i+1])
			if !id.Valid() {
				return nameusage.NameUsageID(""), errors.Newf("Invalid NameUsageID [%s]", a[i+1])
			}
			return id, nil
		}
	}
	return nameusage.NameUsageID(""), errors.Newf("Invalid NameUsageID [%s]", "")
}

type aggregator struct {
	PredictionObjects       predictions.Predictions
	PredictionList          map[nameusage.NameUsageID][]float64
	TotalProtectedAreaCount map[nameusage.NameUsageID]float64
	sync.Mutex
}

func (Ω *predictionParser) parseFile(cxt context.Context, aggr *aggregator, fpath string) error {
	prediction_list, err := Ω.predictionSource.FetchPredictions(cxt, fpath)
	if err != nil {
		return err
	}
	for _, predictionResult := range prediction_list {

		parsedPrediction, err := predictions.NewPrediction(
			predictionResult.NameUsageID,
			predictionResult.Date,
			predictionResult.Latitude,
			predictionResult.Longitude,
			predictionResult.Target,
		)
		if err != nil {
			return err
		}

		aggr.Lock()
		defer aggr.Unlock()

		if _, ok := aggr.TotalProtectedAreaCount[predictionResult.NameUsageID]; !ok {
			aggr.TotalProtectedAreaCount[predictionResult.NameUsageID] = 0
		}
		aggr.TotalProtectedAreaCount[predictionResult.NameUsageID] += 1

		if predictionResult.Target <= predictionResult.Random {
			return nil
		}

		if _, ok := aggr.PredictionList[predictionResult.NameUsageID]; !ok {
			aggr.PredictionList[predictionResult.NameUsageID] = []float64{}
		}
		aggr.PredictionList[predictionResult.NameUsageID] = append(aggr.PredictionList[predictionResult.NameUsageID], predictionResult.Target)

		aggr.PredictionObjects = append(aggr.PredictionObjects, parsedPrediction)

	}
	return nil
}

func (Ω *predictionParser) FetchPredictions(cxt context.Context, nameUsageIDs nameusage.NameUsageIDs, dates []string) (predictions.Predictions, error) {

	aggr := aggregator{
		PredictionObjects:       predictions.Predictions{},
		PredictionList:          make(map[nameusage.NameUsageID][]float64),
		TotalProtectedAreaCount: make(map[nameusage.NameUsageID]float64),
	}

	gcsFilePaths := []string{}
	for _, usageID := range nameUsageIDs {
		if !usageID.Valid() {
			return nil, errors.New("Invalid NameUsageID")
		}
		aggr.PredictionList[usageID] = []float64{}
		aggr.TotalProtectedAreaCount[usageID] = 0
		gcsPaths, err := Ω.predictionSource.FetchLatestPredictionFileNames(cxt, usageID, "*")
		if err != nil {
			return nil, err
		}
		gcsFilePaths = utils.AddStringToSet(gcsFilePaths, gcsPaths...)
	}

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, _fpath := range gcsFilePaths {
			fpath := _fpath
			tmb.Go(func() error {
				return Ω.parseFile(cxt, &aggr, fpath)
			})
		}
		return nil
	})
	if err := tmb.Wait(); err != nil {
		return nil, err
	}

	return aggr.PredictionObjects, nil
}

//
//func (Ω *aggregator) calcTaxonScarcity() (map[nameusage.NameUsageID]float64, error) {
//	taxaRatios := stats.Float64Data{}
//	taxaRatiosMap := map[nameusage.NameUsageID]float64{}
//	for taxon, predictionValues := range Ω.PredictionList {
//		totalTaxonPredictionCount := float64(len(predictionValues))
//		totalTaxonProtectedAreaCount := Ω.TotalProtectedAreaCount[taxon]
//		taxaRatios = append(taxaRatios, totalTaxonPredictionCount/totalTaxonProtectedAreaCount)
//		taxaRatiosMap[taxon] = totalTaxonPredictionCount / totalTaxonProtectedAreaCount
//	}
//
//	taxaRatioMean, err := stats.Mean(taxaRatios)
//	if err != nil {
//		return nil, errors.Wrap(err, "could not calculate mean")
//	}
//
//	var taxonRatioInvertedMin, taxonRatioInvertedMax float64
//	// In order to scale, must calculate the min and the max values once we invert the value by subtracting the mean.
//	// This is so rarer taxa have a higher intensity value.
//	for taxon, ratio := range taxaRatiosMap {
//		invertedValue := taxaRatioMean - ratio
//		if invertedValue < taxonRatioInvertedMin || taxonRatioInvertedMin == 0 {
//			taxonRatioInvertedMin = invertedValue
//		}
//		if invertedValue > taxonRatioInvertedMax || taxonRatioInvertedMax == 0 {
//			taxonRatioInvertedMax = invertedValue
//		}
//		taxaRatiosMap[taxon] = invertedValue
//	}
//	for taxon, ratio := range taxaRatiosMap {
//		// Scale between 1 and 0.5
//		//((b-a)(x - min) / max - min) + a
//		if ratio != 0 {
//			taxaRatiosMap[taxon] = ((1 - 0.5) * (ratio - taxonRatioInvertedMin) / (taxonRatioInvertedMax - taxonRatioInvertedMin)) + 0.5
//		}
//	}
//
//	return taxaRatiosMap, nil
//}
