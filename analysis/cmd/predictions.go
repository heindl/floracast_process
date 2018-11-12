package cmd

import (
	"github.com/heindl/floracast_process/nameusage/nameusage"
	"github.com/heindl/floracast_process/predictions"
	"github.com/heindl/floracast_process/store"
	"github.com/heindl/floracast_process/terra/geo"
	"github.com/heindl/floracast_process/terra/geoembed"
	"github.com/heindl/floracast_process/utils"
	"bytes"
	"cloud.google.com/go/firestore"
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func init() {
	predictionCmd.AddCommand(countPredictionCmd)
	predictionCmd.AddCommand(aggregatePredictionsCmd)
	predictionCmd.AddCommand(uploadPredictionsCmd)
	rootCmd.AddCommand(predictionCmd)
}

var predictionCmd = &cobra.Command{
	Use: "predictions",
}

var countPredictionCmd = &cobra.Command{
	Use:   "count",
	Short: "Count all documents in the given collection",
	RunE:  CountPredictions,
}

func CountPredictions(cmd *cobra.Command, args []string) error {

	ctx := context.Background()

	floraStore, err := store.NewFloraStore(ctx)
	if err != nil {
		return err
	}
	col, err := floraStore.FirestoreCollection(store.CollectionPredictionIndex)
	if err != nil {
		return err
	}

	iter := &firestore.DocumentIterator{}
	if len(args) > 0 {
		iter = col.Where("NameUsageID", "==", args[0]).Documents(ctx)
	} else {
		iter = col.Documents(ctx)
	}

	i := 0
	for {
		_, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		i += 1
	}
	fmt.Println(i)
	return nil

}

var aggregatePredictionsCmd = &cobra.Command{
	Use:   "aggregate",
	Short: "Aggregate all predictions for a NameUsage into a GeoJSON file",
	RunE:  AggregatePredictions,
}

type TimelineField struct {
	Total      int `json:"𝝨"`
	Prediction int `json:"Ω"`
}

func AggregatePredictions(cmd *cobra.Command, args []string) error {

	ctx := context.Background()

	if len(args) != 1 {
		return errors.New("NameUsageID should be only argument")
	}

	floraStore, err := store.NewFloraStore(ctx)
	if err != nil {
		return err
	}

	list, err := predictions.FetchFromFireStore(ctx, floraStore, nameusage.ID(args[0]))
	if err != nil {
		return err
	}

	fmt.Println("Aggregating Predictions:", len(list))

	//geojson := map[utils.FormattedDate]geo.Points{}
	geojson := []string{}
	timeline := map[string][2]float32{}

	for _, pred := range list {

		wkday, err := pred.Date().Weekday()
		if err != nil {
			return err
		}

		if wkday.String() != "Friday" {
			continue
		}

		//if _, ok := geojson[pred.Date()]; !ok {
		//	geojson[pred.Date()] = geo.Points{}
		//}

		point, err := geo.NewPoint(pred.Latitude(), pred.Longitude())
		if err != nil {
			return err
		}

		timelineKey := fmt.Sprintf(`%s-%s`, pred.Date(), point.S2TokenMap()["7"])
		if _, ok := timeline[timelineKey]; !ok {
			timeline[timelineKey] = [2]float32{0, 0}
		}

		//reducedPrecision := [3]float32{}
		//for i, s := range [3]string{
		//	fmt.Sprintf("%.6f", pred.Latitude()),
		//	fmt.Sprintf("%.6f", pred.Longitude()),
		//	fmt.Sprintf("%.3f", pred.Value()),
		//} {
		//	f, err := strconv.ParseFloat(s, 32)
		//	if err != nil {
		//		return errors.Wrap(err, "Could not parse float")
		//	}
		//	reducedPrecision[i] = float32(f)
		//}

		//if err := point.SetProperty("p", float32(reducedPrecision[2])); err != nil {
		//	return err
		//}
		//timeline[timelineKey] = [2]float32{timeline[timelineKey][0] + 1, timeline[timelineKey][1] + reducedPrecision[2]}

		k, err := geoembed.NewS2Key(pred.Latitude(), pred.Longitude())
		if err != nil {
			return err
		}

		if err := point.SetProperty("t", k); err != nil {
			return err
		}

		//if err := point.SetProperty("id", k); err != nil {
		//	return err
		//}

		dateInt, err := strconv.Atoi(string(pred.Date()))
		if err != nil {
			return errors.Wrap(err, "Could not convert date to int.")
		}

		//geojson[pred.Date()] = append(geojson[pred.Date()], point)
		geojson = append(geojson, fmt.Sprintf(
			"%d,%.6f,%.6f,%.3f",
			dateInt,
			pred.Latitude(),
			pred.Longitude(),
			pred.Value()))
	}

	f1 := fmt.Sprintf("/Users/m/Desktop/%s.csv.gz", args[0])

	var buffer bytes.Buffer
	w, err := gzip.NewWriterLevel(&buffer, gzip.BestCompression)
	if err != nil {
		return errors.Wrap(err, "Could not write file")
	}
	for _, line := range geojson {
		fmt.Fprintln(w, line)
	}
	if err := w.Flush(); err != nil {
		return errors.Wrap(err, "Could not flush buffer")
	}

	if err := w.Close(); err != nil {
		return errors.Wrap(err, "Could not close GZIP compressor")
	}

	if err := ioutil.WriteFile(
		f1,
		buffer.Bytes(),
		os.ModePerm,
	); err != nil {
		return errors.Wrapf(err, "Could not write file: %s", f1)
	}

	//b1, err := json.Marshal(geojson)
	//if err != nil {
	//	return errors.Wrap(err, "Could not marshal geojson")
	//}

	//f2 := fmt.Sprintf("/Users/m/Desktop/prediction-timeline-%s.json", args[0])
	//
	//fmt.Println("SummaryDocs", len(timeline), len(geojson))
	//
	//ungrouped := map[string][][3]interface{}{}
	//for k, v := range timeline {
	//	date := strings.Split(k, "-")[0]
	//	token := strings.Split(k, "-")[1]
	//	if _, ok := ungrouped[token]; !ok {
	//		ungrouped[token] = [][3]interface{}{}
	//	}
	//	ungrouped[token] = append(ungrouped[token], [3]interface{}{date, v[0], v[1]})
	//}
	//
	//b2, err := json.Marshal(ungrouped)
	//if err != nil {
	//	return errors.Wrap(err, "Could not marshal timeline")
	//}
	//
	//if err := ioutil.WriteFile(f2, b2, os.ModePerm); err != nil {
	//	return errors.Wrapf(err, "Could not write file: %s", f2)
	//}

	return nil

}

var uploadPredictionsCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload predictions from all GCS files",
	RunE:  UploadPredictionsFromGCS,
}

func UploadPredictionsFromGCS(cmd *cobra.Command, args []string) error {

	ctx := context.Background()

	floraStore, err := store.NewFloraStore(ctx)
	if err != nil {
		return err
	}

	predictionCollection, err := predictions.NewGroupedCollection(floraStore)
	if err != nil {
		return err
	}

	handles, err := floraStore.CloudStorageObjects(ctx, "predictions/", ".csv")
	if err != nil {
		return err
	}

	for _, handle := range handles {

		attrs, err := handle.Attrs(ctx)
		if err != nil {
			return err
		}

		nameUsageID := nameusage.ID(strings.TrimSuffix(strings.TrimPrefix(attrs.Name, "predictions/"), ".csv"))
		if !nameUsageID.Valid() {
			return errors.Newf("Invalid NameUsageID [%s] with path prefix [predictions/]", nameUsageID)
		}

		reader, err := handle.NewReader(ctx)
		if err != nil {
			return errors.Wrapf(err, "Could not read GCS object")
		}

		rows, err := csv.NewReader(reader).ReadAll()
		if err != nil {
			return errors.Wrap(err, "Could not get csv rows")
		}

		for _, r := range rows {
			date := utils.FormattedDate(r[0])

			latitude, err := strconv.ParseFloat(r[1], 64)
			if err != nil {
				return err
			}

			longitude, err := strconv.ParseFloat(r[2], 64)
			if err != nil {
				return err
			}

			prediction, err := strconv.ParseFloat(r[3], 64)
			if err != nil {
				return err
			}

			predictionCollection.Add(nameUsageID, latitude, longitude, date, prediction)

		}
	}

	if err := predictionCollection.Upload(ctx); err != nil {
		return err
	}

	return nil

}
