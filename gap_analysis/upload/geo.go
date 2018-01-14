package main

import (
	"github.com/tidwall/tile38/geojson"
	"github.com/dropbox/godropbox/errors"
)

func SecondOpinionCentroid(gb []byte, pastLat, pastLng float64) (newLat, newLng float64, newContains, pastContains bool, err error) {
	geoObj, err := geojson.ObjectJSON(string(gb))
	if err != nil {
		err = errors.Wrap(err, "could not parse geojson object")
		return
	}

	p := geoObj.CalculatedPoint()
	newLat = p.Y
	newLng = p.X

	newContains = geojson.New2DPoint(newLng, newLat).Within(geoObj)
	pastContains = geojson.New2DPoint(pastLng, pastLat).Within(geoObj)

	return

}
