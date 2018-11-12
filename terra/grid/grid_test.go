package grid

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/tomb.v2"
	"testing"
)

func TestGridGenerator(t *testing.T) {

	t.Parallel()

	Convey("should generate grid", t, func() {

		total := 0
		tmb := tomb.Tomb{}
		tmb.Go(func() error {
			for i := 1; i < 10; i++ {
				_i := i
				tmb.Go(func() error {
					atLevel := CountNorthAmericaAtLevel(_i)
					fmt.Println("Level", i, atLevel)
					total += atLevel
					return nil
				})
			}
			return nil
		})
		So(tmb.Wait(), ShouldBeNil)

		fmt.Println("Total", total)

		//g, err := NewGridGenerator()
		//So(err, ShouldBeNil)
		//list, err := g.SubDivide(NorthAmerica, 3)
		//So(err, ShouldBeNil)
		//So(len(list), ShouldEqual, 16)
		//_, err = list.ToGeoJSON()
		//So(err, ShouldBeNil)
	})

}
