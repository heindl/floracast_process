package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"bitbucket.org/heindl/taxa/store"
	"golang.org/x/net/context"
	"fmt"
	"sync"
	"bitbucket.org/heindl/taxa/predictions/geocache"
)

type TestWriter struct {
	Count int
	sync.Mutex
}

func (Ω *TestWriter) WritePredictionLine(p store.Prediction) error {
	Ω.Lock()
	defer Ω.Unlock()
	Ω.Count = Ω.Count + 1
	fmt.Println(p.Month, p.WildernessAreaName)
	return nil
}

func (Ω *TestWriter) Close() error {
	fmt.Println("Prediction Count", Ω.Count)
	return nil
}

func TestPredictionParser(t *testing.T) {

	t.Parallel()

	SkipConvey("should fetch occurrences and add to queue", t, func() {

		cxt := context.Background()
		fetcher, err := NewGCSFetcher(cxt, "floracast-datamining")
		So(err, ShouldBeNil)

		files, err := fetcher.FetchLatestPredictionFileNames(cxt, store.TaxonID("58682"), "")
		So(err, ShouldBeNil)
		So(len(files), ShouldEqual, 1)

		predictions, err := fetcher.FetchPredictions(cxt, files[0])
		So(err, ShouldBeNil)
		So(len(predictions), ShouldEqual, 812)

	})

	SkipConvey("should create fetcher and parse predictions", t, func() {

		cxt := context.Background()
		writer := TestWriter{}
		parser, err := NewPredictionParser(cxt, "floracast-datamining", &writer)
		So(err, ShouldBeNil)
		So(parser.FetchWritePredictions(cxt, store.TaxonID("58682"), ""), ShouldBeNil)
		writer.Close()

	})

	Convey("should add items to geocache", t, func() {

		cxt := context.Background()
		writer, err := geocache.NewCacheWriter([]string{"58682"})
		So(err, ShouldBeNil)
		parser, err := NewPredictionParser(cxt, "floracast-datamining", writer)
		So(err, ShouldBeNil)
		So(parser.FetchWritePredictions(cxt, store.TaxonID("58682"), ""), ShouldBeNil)

		taxa, err := writer.ReadTaxon(store.TaxonID("58682"), 38.6530169,-90.3835463, 2000)
		So(err, ShouldBeNil)
		fmt.Println(taxa)

	})
}
