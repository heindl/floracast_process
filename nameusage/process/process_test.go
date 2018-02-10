package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"context"
	"fmt"
	"bitbucket.org/heindl/processors/utils"
	"bitbucket.org/heindl/processors/nameusage/nameusage"
)

func TestTaxonFetcher(t *testing.T) {

	t.Parallel()

	Convey("should generate a list of sources", t, func() {

		cxt := context.Background()
		// 47348, 56830, 48701
		nameUsageAggr, err := InitialAggregation(cxt, 47348)
		if err != nil {
			panic(err)
		}

		So(nameUsageAggr.Each(cxt, func(ctx context.Context, usage *nameusage.NameUsage) error{
			fmt.Println(utils.JsonOrSpew(usage))
			return nil
		}), ShouldBeNil)

	})
}
