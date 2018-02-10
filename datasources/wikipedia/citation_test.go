package wikipedia

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTaxonFetcher(t *testing.T) {

	t.Parallel()

	Convey("should generate wikipedia citation", t, func() {
		c, err := Citation("https://en.wikipedia.org/wiki/Morchella_esculenta")
		So(err, ShouldBeNil)
		Println(c)
	})
}
