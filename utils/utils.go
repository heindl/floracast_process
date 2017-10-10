package utils

import (
	"bytes"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"os"
	"strings"
	"time"
"github.com/kpawlik/geojson"
	"fmt"
)

var NorthAmericaBBOX = [2][2]float64{{-169.433594, 13.267549}, {-49.902344, 57.906568}}

func JsonOrSpew(o interface{}) string {
	j, err := json.Marshal(o)
	if err != nil {
		return spew.Sprintf("%+v", o)
	}
	var response bytes.Buffer
	if err := json.Indent(&response, j, "", "  "); err != nil {
		return spew.Sprintf("%+v", o)
	}
	return response.String()
}

func TimePtr(v time.Time)  *time.Time {
	return &v
}

func BoolPtr(v bool) *bool {
	return &v
}

func Contains(haystack []string, needle string) bool {
	for _, str := range haystack {
		if needle == str {
			return true
		}
	}
	return false
}

func MustParse(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}

func MustParseAsPointer(layout, value string) *time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return &t
}

func GetFileContents(path string) []byte {

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	return body

}

func AddIntToSet(haystack []int, needle int) []int {
	if !ContainsInt(haystack, needle) {
		return append(haystack, needle)
	}
	return haystack
}

func ContainsInt(haystack []int, needle int) bool {
	for _, h := range haystack {
		if h == needle {
			return true
		}
	}
	return false
}

func AddStringToSet(haystack []string, needle string) []string {
	if !Contains(haystack, needle) {
		return append(haystack, needle)
	}
	return haystack
}

func ContainsError(e error, d error) bool {
	fmt.Println("d", d.Error())
	fmt.Println("e", e.Error())
	return strings.Contains(d.Error(), e.Error())
}

func GeoJsonPoint(longitude float64, latitude float64) geojson.Point {
	p := geojson.NewPoint(geojson.Coordinate{
		geojson.Coord(longitude),
		geojson.Coord(latitude),
	})
	return *p
}