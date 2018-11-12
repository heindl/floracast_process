package cmd

import (
	"github.com/heindl/floracast_process/nameusage/nameusage"
	"github.com/heindl/floracast_process/occurrence"
	"github.com/heindl/floracast_process/store"
	"github.com/heindl/floracast_process/utils"
	"bytes"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"github.com/spf13/cobra"
	"github.com/twpayne/go-kml"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"time"
)

func init() {
	occurrenceCmd.AddCommand(printOccurrenceKMLCmd)
	occurrenceCmd.AddCommand(aggregateOccurrencesCmd)
	occurrenceCmd.AddCommand(aggregateOccurrencesByYearCmd)
	occurrenceCmd.AddCommand(countOccurrenceCmd)
	rootCmd.AddCommand(occurrenceCmd)
}

var occurrenceCmd = &cobra.Command{
	Use: "occurrences",
}

var aggregateOccurrencesByYearCmd = &cobra.Command{
	Use:   "aggregate-by-year",
	Short: "Map all occurrences by year",
	RunE:  AggregateOccurrencesByYear,
}

var printOccurrenceKMLCmd = &cobra.Command{
	Use:   "kml",
	Short: "Print csv file into terminal",
	RunE:  PrintOccurrenceKML,
}

func PrintOccurrenceKML(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if len(args) != 1 {
		return errors.New("NameUsageID should be only argument")
	}

	floraStore, err := store.NewFloraStore(ctx)
	if err != nil {
		return err
	}

	nameUsageId := nameusage.ID(args[0])

	list, err := occurrence.FetchFromFireStore(ctx, floraStore, nameUsageId)
	if err != nil {
		return err
	}

	elements := []kml.Element{}

	for _, o := range list {

		date, err := o.Date()
		if err != nil {
			return err
		}

		lat, lng, err := o.Coordinates()
		if err != nil {
			return err
		}

		t, err := time.Parse("20060102", date)
		if err != nil {
			return err
		}

		elements = append(elements, kml.Placemark(
			kml.ExtendedData(
				kml.SchemaData(
					"occurrences",
					kml.SimpleData("occurrenceId",
						o.SourceOccurrenceID(),
					),
					kml.SimpleData(
						"sourceType",
						string(o.SourceType()),
					),
					kml.SimpleData(
						"targetId",
						string(o.TargetID()),
					),
					kml.SimpleData(
						"nameUsageId",
						string(nameUsageId),
					),
					kml.SimpleData(
						"time_start",
						t.Format("2006-01-02"),
					),
					kml.SimpleData(
						"time_end",
						t.Format("2006-01-02"),
					),
					kml.SimpleData(
						"system:time_start",
						t.Format("2006-01-02"),
					),
					kml.SimpleData(
						"system:time_end",
						t.Format("2006-01-02"),
					),
				),
			),
			kml.TimeStamp(
				kml.When(t),
			),
			kml.Point(
				kml.Coordinates(kml.Coordinate{Lon: lng, Lat: lat}),
			),
		))
	}

	if err := kml.KML(elements...).WriteIndent(os.Stdout, "", "  "); err != nil {
		log.Fatal(err)
	}

	return nil
}

func AggregateOccurrencesByYear(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	floraStore, err := store.NewFloraStore(ctx)
	if err != nil {
		return err
	}
	col, err := floraStore.FirestoreCollection(store.CollectionOccurrences)
	if err != nil {
		return err
	}
	iter := col.Documents(ctx)
	yearMap := map[string]int{}
	decadeMap := map[string]int{}
	years := []string{}
	decades := []string{}
	for {
		snap, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		date, err := snap.DataAt("FormattedDate")
		if err != nil {
			return errors.Wrap(err, "Could not get FormattedDate")
		}

		decade := date.(string)[:3]
		year := date.(string)[:4]

		decades = utils.AddStringToSet(decades, decade)
		years = utils.AddStringToSet(years, year)

		if _, ok := yearMap[year]; !ok {
			yearMap[year] = 0
		}
		yearMap[year] += 1

		if _, ok := decadeMap[decade]; !ok {
			decadeMap[decade] = 0
		}
		decadeMap[decade] += 1
	}

	sort.Strings(years)
	sort.Strings(decades)

	for _, y := range years {
		fmt.Println(fmt.Sprintf("%s: %d", y, yearMap[y]))
	}
	for _, d := range decades {
		fmt.Println(fmt.Sprintf("%s: %d", d, decadeMap[d]))
	}
	return nil
}

var countOccurrenceCmd = &cobra.Command{
	Use:   "count",
	Short: "Count all documents in the given collection",
	RunE:  CountOccurrences,
}

func CountOccurrences(cmd *cobra.Command, args []string) error {

	ctx := context.Background()

	floraStore, err := store.NewFloraStore(ctx)
	if err != nil {
		return err
	}
	col, err := floraStore.FirestoreCollection(store.CollectionOccurrences)
	if err != nil {
		return err
	}
	iter := col.Documents(ctx)
	m := map[string]int{}

	for {
		snap, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		date, err := snap.DataAt("FormattedDate")
		if err != nil {
			return errors.Wrap(err, "Could not get FormattedDate")
		}
		year := date.(string)[:4]
		if year < "1999" {
			continue
		}

		usageID, err := snap.DataAt("NameUsageID")
		if err != nil {
			return errors.Wrap(err, "Could not get NameUsageID")
		}

		if _, ok := m[usageID.(string)]; !ok {
			m[usageID.(string)] = 0
		}
		m[usageID.(string)] += 1
	}

	for usageID, count := range m {

		col, err := floraStore.FirestoreCollection(store.CollectionTaxa)
		if err != nil {
			return err
		}
		doc, err := col.Doc(usageID).Get(ctx)
		if err != nil {
			return errors.Wrapf(err, "Could not get Taxa [%s]", usageID)
		}
		n, err := doc.DataAt("CommonName")
		if err != nil {
			return errors.Wrapf(err, "Could not get CommonName [%s]", usageID)
		}
		fmt.Println(fmt.Sprintf("%s: %d", n.(string), count))
	}

	return nil

}

var aggregateOccurrencesCmd = &cobra.Command{
	Use:   "aggregate",
	Short: "Aggregate all occurrences for a NameUsage into a GeoJSON file",
	RunE:  AggregateOccurrences,
}

func AggregateOccurrences(cmd *cobra.Command, args []string) error {

	ctx := context.Background()

	if len(args) != 1 {
		return errors.New("NameUsageID should be only argument")
	}

	floraStore, err := store.NewFloraStore(ctx)
	if err != nil {
		return err
	}

	list, err := occurrence.FetchFromFireStore(ctx, floraStore, nameusage.ID(args[0]))
	if err != nil {
		return err
	}

	fmt.Println("Aggregating Occurrences:", len(list))

	//docs := map[utils.FormattedDate][][3]string{}
	docs := []string{}

	for _, o := range list {

		date, err := o.Date()
		if err != nil {
			return err
		}

		id, err := o.ID()
		if err != nil {
			return err
		}

		//key, err := o.LocationKey()
		//if err != nil {
		//	return err
		//}
		//
		//token := strings.Split(key, "|")[0]

		lat, lng, err := o.Coordinates()
		if err != nil {
			return err
		}

		docs = append(docs, fmt.Sprintf("%s,%.6f,%.6f,0,%s", date, lat, lng, id))
	}

	f1 := fmt.Sprintf("/Users/m/Desktop/%s.csv", args[0])

	b := bytes.NewBuffer([]byte{})
	for _, line := range docs {
		fmt.Fprintln(b, line)
	}

	if err := ioutil.WriteFile(
		f1,
		b.Bytes(),
		os.ModePerm,
	); err != nil {
		return errors.Wrapf(err, "Could not write file: %s", f1)
	}

	//var buffer bytes.Buffer
	//w, err := gzip.NewWriterLevel(&buffer, gzip.BestCompression)
	//if err != nil {
	//	return errors.Wrap(err, "Could not write file")
	//}
	//if _, err := w.Write(b); err != nil {
	//	return errors.Wrap(err, "Could not write GZIP")
	//}
	//if err := w.Close(); err != nil {
	//	return errors.Wrap(err, "Could not close GZIP compressor")
	//}

	//filename := fmt.Sprintf("/tmp/%s.json", args[0])
	//
	//if err := ioutil.WriteFile(filename, b, os.ModePerm); err != nil {
	//	return errors.Wrapf(err, "Could not write file: %s", filename)
	//}

	return nil

}
