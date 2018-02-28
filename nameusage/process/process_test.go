package main

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/datasources/sourcefetchers"
	"bitbucket.org/heindl/process/nameusage/aggregate"
	"bitbucket.org/heindl/process/store"
	"bitbucket.org/heindl/process/utils"
	"context"
	"encoding/json"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNameUsageProcessor(t *testing.T) {

	t.Parallel()

	SkipConvey("should generate occurrence usages", t, func() {
		usages, err := sourcefetchers.FetchNameUsages(
			context.Background(),
			datasources.TypeINaturalist,
			nil,
			datasources.TargetIDs{"58682"},
		)
		So(err, ShouldBeNil)
		So(len(usages), ShouldEqual, 2)
	})

	SkipConvey("should execute nameusage aggregation among three sources.", t, func() {

		// 47348, 56830, 48701
		aggr, err := InitialAggregation(context.Background(), 58682)
		if err != nil {
			panic(err)
		}
		So(aggr.Count(), ShouldEqual, 1)

		fmt.Println(utils.JsonOrSpew(&aggr))

	})

	Convey("should fetch occurrences from name usages", t, func() {

		aggr := aggregate.Aggregate{}
		So(json.Unmarshal(utils.GetUnFetchedMorchellaAggregateTestData(), &aggr), ShouldBeNil)

		So(aggr.Count(), ShouldEqual, 1)

		occurrenceAggr, err := OccurrenceFetch(context.Background(), &aggr)
		So(err, ShouldBeNil)

		So(occurrenceAggr.Count(), ShouldEqual, 238)

		fmt.Println(utils.JsonOrSpew(occurrenceAggr))

	})

	SkipConvey("should upload occurrence count", t, func() {

		aggr := aggregate.Aggregate{}
		So(json.Unmarshal(utils.GetFetchedMorchellaAggregateTestData(), &aggr), ShouldBeNil)

		So(aggr.Count(), ShouldEqual, 1)

		cxt := context.Background()

		florastore, err := store.NewFloraStore(cxt)
		So(err, ShouldBeNil)

		So(aggr.Upload(cxt, florastore), ShouldBeNil)

	})
}
