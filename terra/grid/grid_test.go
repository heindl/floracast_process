package grid

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGridGenerator(t *testing.T) {

	t.Parallel()

	Convey("should generate grid", t, func() {
		g, err := NewGridGenerator()
		So(err, ShouldBeNil)
		list, err := g.SubDivide(NorthAmerica, 3)
		So(err, ShouldBeNil)
		So(len(list), ShouldEqual, 16)
		_, err = list.ToGeoJSON()
		So(err, ShouldBeNil)
	})

}
