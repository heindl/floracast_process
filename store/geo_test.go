package store

import (
	"bitbucket.org/heindl/taxa/utils"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/context"
	"testing"
)

func TestOccurrenceFetcher(t *testing.T) {

	t.Parallel()

	Convey("should fetch geo features", t, func() {

		processor, err := NewGeoFeaturesProcessor()
		So(err, ShouldBeNil)

		_ = PredictableLocation(&TestGeoLocation{lat: 41.1491573, lng: -115.4622611})

		locations := []PredictableLocation{}

		locations = append(locations,
			&TestGeoLocation{lat: 41.1491573, lng: -115.4622611},
			&TestGeoLocation{lat: 47.2600975, lng: -120.2742729},
			&TestGeoLocation{lat: 46.4411401, lng: -117.8572807},
		)

		So(processor.ProcessLocations(context.Background(), locations...), ShouldBeNil)

		fmt.Println(utils.JsonOrSpew(locations))
	})
}

type TestGeoLocation struct {
	lat, lng float64
	Features GeoFeatures
}

func (Ω *TestGeoLocation) Lat() float64 {
	return Ω.lat
}

func (Ω *TestGeoLocation) Lng() float64 {
	return Ω.lng
}

func (Ω *TestGeoLocation) SetGeoFeatures(g *GeoFeatures) {
	Ω.Features = *g
}
