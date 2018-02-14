package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"context"
	"fmt"
	"bitbucket.org/heindl/processors/utils"
	"bitbucket.org/heindl/processors/nameusage/aggregate"
	"encoding/json"
	"bitbucket.org/heindl/processors/datasources/sourcefetchers"
	"bitbucket.org/heindl/processors/datasources"
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
		So(json.Unmarshal(utils.GetMorchellaAggregateTestData(), &aggr), ShouldBeNil)

		So(aggr.Count(), ShouldEqual, 1)

		occurrenceAggr, err := OccurrenceFetch(context.Background(), &aggr)
		So(err, ShouldBeNil)

		So(occurrenceAggr.Count(), ShouldEqual, 238)

		fmt.Println(utils.JsonOrSpew(&aggr))

	})
}
