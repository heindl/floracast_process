package sourcefetchers

import (
	"bitbucket.org/heindl/process/datasources"
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNameUsageFetcher(t *testing.T) {

	t.Parallel()

	Convey("should generate occurrence usages", t, func() {
		usages, err := FetchNameUsages(
			context.Background(),
			datasources.TypeINaturalist,
			nil,
			datasources.TargetIDs{"58682"},
		)
		So(err, ShouldBeNil)
		So(len(usages), ShouldEqual, 2)

		srcs, err := usages[0].Sources()
		So(err, ShouldBeNil)
		So(len(srcs), ShouldEqual, 2)
		names, err := usages[0].AllScientificNames()
		So(err, ShouldBeNil)
		So(names, ShouldContain, "morchella esculenta umbrina")

		srcs, err = usages[1].Sources()
		So(err, ShouldBeNil)
		So(len(srcs), ShouldEqual, 2)
		names, err = usages[1].AllScientificNames()
		So(err, ShouldBeNil)
		So(names, ShouldContain, "morchella esculenta")

	})

}
