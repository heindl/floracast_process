package elevation

import (
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/tomb.v2"
	"testing"
)

func TestElevationFetcher(t *testing.T) {

	Convey("should fetch elevations", t, func(c C) {

		coordinates := [][3]float64{
			{42.722702, -87.784225, 192},
			{32.346596, -106.787720, 1188},
			{42.326515, -122.875595, 424},
			{37.765205, -122.241638, 11},
			{37.910076, -122.065186, 52},
			{33.487007, -117.143784, 304},
			{41.653934, -81.450394, 189},
			{46.602070, -120.505898, 326},
			{28.018349, -82.764473, 17},
			{44.949642, -93.093124, 240},
			{47.380932, -122.234840, 13},
			{42.336983, -83.273262, 191},
			{34.704929, -81.210251, 151},
			{43.224194, -86.235809, 197},
			{34.426388, -117.300880, 969},
			{34.686787, -118.154160, 726},
			{25.468721, -80.477554, 7},
			{36.850769, -76.285873, 6},
			{36.974117, -122.030792, 7},
			{37.687923, -122.470207, 49},
			{42.902615, -78.744576, 201},
			{39.932117, -105.281639, 1831},
			{42.562786, -114.460503, 1136},
			{38.413651, -82.446732, 171},
			{41.394966, -73.454079, 123},
			{41.308273, -72.927887, 15},
			{43.565468, -116.560822, 763},
			{30.005417, -91.818665, 12},
			{41.245708, -75.881241, 173},
			{37.227928, -77.401924, 22},
			{42.963795, -85.670006, 199},
			{38.047989, -84.501640, 295},
			{41.385826, -72.904243, 27},
			{46.731705, -116.999939, 789},
			{34.224400, -92.019211, 70},
			{47.610378, -122.200676, 27},
			{37.412762, -79.146042, 209},
			{39.161079, -75.525681, 14},
			{36.136906, -111.240173, 1519},
			{42.326515, -122.875595, 424},
			{35.925064, -86.868889, 200},
			{45.523064, -122.676483, 14},
			{39.952583, -75.165222, 4},
			{25.761681, -80.191788, 22},
			{42.361145, -71.057083, 11},
			{33.448376, -112.074036, 351},
			{42.129223, -80.085060, 203},
			{34.536217, -117.292763, 835},
			{39.739071, -75.539787, 9},
			{34.092232, -117.435051, 380},
		}

		for _, c := range coordinates {
			coordinate := c
			So(Queue(coordinate[0], coordinate[1]), ShouldBeNil)
		}
		//
		tmb := tomb.Tomb{}
		tmb.Go(func() error {
			for _, coord := range coordinates {
				coordinate := coord
				tmb.Go(func() error {
					e, err := Get(coordinate[0], coordinate[1])
					c.So(err, ShouldBeNil)
					c.So(err, ShouldBeNil)
					c.So(e, ShouldNotBeNil)
					c.So(coordinate[2], ShouldEqual, float64(*e))
					return nil
				})
			}
			return nil
		})
		So(tmb.Wait(), ShouldBeNil)

	})
}
