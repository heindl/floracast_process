package geofeatures

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"bitbucket.org/heindl/processors/utils"
)

func TestGeoFeatureGeneration(t *testing.T) {

	t.Parallel()

	Convey("should fetch geo features", t, func() {

		for _, a := range [][2]float64{
			{41.1491573, -115.4622611},
			{41.1491573, -115.4622611},
			{47.2600975, -120.2742729},
			{46.4411401, -117.8572807},
		} {
			s, err := NewGeoFeatureSet(a[0], a[1], false)
			So(err, ShouldBeNil)
			Println(utils.JsonOrSpew(s))
		}
	})
}
