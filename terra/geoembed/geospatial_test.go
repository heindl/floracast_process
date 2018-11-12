package geoembed

import (
	"github.com/heindl/floracast_process/utils"
	"encoding/json"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGeoFeatureGeneration(t *testing.T) {

	Convey("GeoFeatures", t, func() {

		_, err := NewS2Key(44.094082, -117.869965)
		So(err, ShouldBeNil)

	})

	SkipConvey("Should properly create GeoFeatureSet and marshal/unmarshal JSON on struct where it is embedded", t, func() {

		type geoEmbeddedStruct struct {
			Name          string         `json:""`
			GeoFeatureSet *GeoFeatureSet `json:""`
		}

		initialSet, err := NewGeoFeatureSet(33.5724059, -84.7170809, false)
		So(err, ShouldBeNil)

		fmt.Println(utils.JsonOrSpew(initialSet))

		embed := geoEmbeddedStruct{
			Name:          "test",
			GeoFeatureSet: initialSet,
		}

		b, err := json.Marshal(embed)
		So(err, ShouldBeNil)

		unembedded := geoEmbeddedStruct{}
		So(json.Unmarshal(b, &unembedded), ShouldBeNil)

		So(unembedded.Name, ShouldEqual, "test")
		So(unembedded.GeoFeatureSet.Lat(), ShouldEqual, 33.5724059)
		So(unembedded.GeoFeatureSet.Lng(), ShouldEqual, -84.7170809)

	})

	SkipConvey("Should consistently generate GeoFeatureSets", t, func() {

		for _, a := range [][6]float64{
			// Lat, Lng, Elevation, Biome, Realm, EcoNum
			{43.4732679, -110.7998022, 1874, 5, 5, 28}, // Jackson Hole, Wyoming
			{40.0292888, -105.310018, 1932, 5, 5, 11},  // Boulder, Colorado
			{33.5309219, -87.1303357, 140, 4, 5, 2},    // Birmingham, Alabama
			{40.4313473, -80.050541, 340, 4, 5, 2},     // Pittsburg, Pennsylvania
			{37.9318439, -122.295833, 207, 12, 5, 2},   // Berkeley, California
			{38.57654, -109.5816315, 1236, 13, 5, 4},   // Moab, Utah
		} {
			initialSet, err := NewGeoFeatureSet(a[0], a[1], false)
			So(err, ShouldBeNil)

			b, err := json.Marshal(initialSet)
			So(err, ShouldBeNil)

			newSet := GeoFeatureSet{}
			So(json.Unmarshal(b, &newSet), ShouldBeNil)

			So(newSet.geoPoint.GetLongitude(), ShouldEqual, a[1])
			So(newSet.geoPoint.GetLatitude(), ShouldEqual, a[0])
			So(newSet.elevation, ShouldNotBeNil)
			So(*newSet.elevation, ShouldEqual, a[2])
			So(newSet.biome, ShouldEqual, a[3])
			So(newSet.realm, ShouldEqual, a[4])
			So(newSet.ecoNum, ShouldEqual, a[5])

		}
	})
}
