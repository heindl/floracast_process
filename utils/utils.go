package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"os"
)

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

func JsonOrSpew(o interface{}) string {
	j, err := json.Marshal(o)
	if err != nil {
		fmt.Println("MARSHAL ERROR", err)
		return spew.Sprintf("%+v", o)
	}
	var response bytes.Buffer
	if err := json.Indent(&response, j, "", "  "); err != nil {
		return spew.Sprintf("%+v", o)
	}
	return response.String()
}

func BoolPtr(v bool) *bool {
	return &v
}

func FloatPtr(v float64) *float64 {
	return &v
}

func IntPtr(v int) *int {
	return &v
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
