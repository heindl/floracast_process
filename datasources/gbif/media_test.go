package gbif

import (
	"bitbucket.org/heindl/process/datasources"
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMediaFetcher(t *testing.T) {

	t.Parallel()

	Convey("Should fetch PhotoProviders", t, func() {
		photos, err := FetchPhotos(context.Background(), datasources.TargetID("2594602"))
		So(err, ShouldBeNil)
		So(len(photos), ShouldEqual, 4)
		for i, url := range []string{
			"http://www.mycokey.com/MycoKeySolidState/pictures/asco/disc/oper/Morchella/escu16L.jpg",
			"http://www.mycokey.com/MycoKeySolidState/pictures/asco/disc/oper/Morchella/escu3L.jpg",
			"http://www.mycokey.com/MycoKeySolidState/pictures/asco/disc/oper/Morchella/escu4L.jpg",
			"http://www.mycokey.com/MycoKeySolidState/pictures/asco/disc/oper/Morchella/escu5L.jpg",
		} {
			So(photos[i].Citation(), ShouldEqual, "Jens H. Petersen, Checklist of Danish Fungi")
			So(photos[i].Thumbnail(), ShouldEqual, "")
			So(photos[i].Large(), ShouldEqual, url)
			So(photos[i].Source(), ShouldEqual, datasources.TypeGBIF)
		}
	})

}
