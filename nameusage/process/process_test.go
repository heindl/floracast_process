package main

import (
	"github.com/heindl/floracast_process/nameusage/aggregate"
	"github.com/heindl/floracast_process/utils"
	"context"
	"encoding/json"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNameUsageProcessor(t *testing.T) {

	Convey("Should execute NameUsage Aggregation", t, func() {

		// 47348, 56830, 48701
		aggr, err := aggregateInitialNameUsages(context.Background(), 58682)
		if err != nil {
			panic(err)
		}
		So(aggr.Count(), ShouldEqual, 1)

		fmt.Println(utils.JsonOrSpew(&aggr))

	})

	Convey("Should fetch occurrences for NameUsages", t, func() {

		aggr := aggregate.Aggregate{}
		So(json.Unmarshal(utils.GetUnFetchedMorchellaAggregateTestData(), &aggr), ShouldBeNil)

		So(aggr.Count(), ShouldEqual, 1)

		occurrenceAggr, err := fetchOccurrences(context.Background(), &aggr)
		So(err, ShouldBeNil)

		So(occurrenceAggr.Count(), ShouldEqual, 238)

		fmt.Println(utils.JsonOrSpew(occurrenceAggr))

	})
}
