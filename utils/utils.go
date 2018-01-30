package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/kpawlik/geojson"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var NorthAmericaBBOX = [2][2]float64{{-169.433594, 13.267549}, {-49.902344, 57.906568}}

var EPSILON float64 = 0.00001

func CoordinatesEqual(a, b float64) bool {
	if (a-b) < EPSILON && (b-a) < EPSILON {
		return true
	}
	return false
}

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

func TimePtr(v time.Time) *time.Time {
	return &v
}

func BoolPtr(v bool) *bool {
	return &v
}

func FloatPtr(v float64) *float64 {
	return &v
}

func RemoveStringDuplicates(haystack []string) []string {

	res := []string{}

	for _, s := range haystack {
		if ContainsString(res, s) {
			continue
		}
		res = append(res, s)
	}

	return res
}

func IndexOfString(haystack []string, needle string) int {
	for i := range haystack {
		if needle == haystack[i] {
			return i
		}
	}
	return -1
}

func ContainsString(haystack []string, needle string) bool {
	return IndexOfString(haystack, needle) != -1
}

func IntersectsStrings(a []string, b []string) bool {
	for _, s := range a {
		if ContainsString(b, s) {
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

func StringsToLower(a ...string) []string {
	res := []string{}
	for _, s := range a {
		res = append(res, strings.ToLower(s))
	}
	return res
}

func AddStringToSet(haystack []string, needles ...string) []string {
	for _, needle := range needles {
		if needle == "" {
			continue
		}
		if !ContainsString(haystack, needle) {
			haystack = append(haystack, needle)
		}
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

type Limiter chan struct{}

func NewLimiter(amount int) Limiter {
	limiter := make(chan struct{}, amount)
	for i := 0; i < amount; i++ {
		limiter <- struct{}{}
	}
	return limiter
}

func (Ω Limiter) Go() func() {
	<-Ω
	return func() {
		Ω <- struct{}{}
	}
}
