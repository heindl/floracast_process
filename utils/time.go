package utils

import (
	"github.com/dropbox/godropbox/errors"
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

func (Ω FormattedDate) Weekday() (time.Weekday, error) {
	t, err := time.Parse("20060102", string(Ω))
	if err != nil {
		return 0, errors.Wrapf(err, "Could not parse date: %s", Ω)
	}
	return t.Weekday(), nil
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
