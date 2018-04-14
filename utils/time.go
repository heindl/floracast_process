package utils

import (
	"strconv"
	"time"
)

type FormattedDate string

func (Ω FormattedDate) Valid() bool {
	if len(Ω) != 8 {
		return false
	}
	if _, err := strconv.Atoi(string(Ω)); err != nil {
		return false
	}
	return true
}

func MustParseTime(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}

func MustParseAsTimePointer(layout, value string) *time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return &t
}

func TimePtr(v time.Time) *time.Time {
	return &v
}
