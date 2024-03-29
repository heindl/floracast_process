package utils

import (
	"github.com/dropbox/godropbox/errors"
	"strconv"
	"time"
)

type RangeVal string

func (start RangeVal) ExpandTo(end RangeVal) ([]RangeVal, error) {
	s, err := start.Unmarshal()
	if err != nil {
		return nil, err
	}
	e, err := end.Unmarshal()
	if err != nil {
		return nil, err
	}
	if s.After(e) {
		return nil, errors.New("expected start to be before end time while expanding ranges")
	}
	if s.Equal(e) {
		return []RangeVal{MarshalRangeVal(s)}, nil
	}
	res := []RangeVal{}
	for i := s; i.Before(e) || i.Equal(e); i = i.Add(time.Hour * 24) {
		res = append(res, MarshalRangeVal(i))
	}
	return res, nil
}

func (start RangeVal) AddDays(days int) (RangeVal, error) {
	t, err := start.Unmarshal()
	if err != nil {
		return RangeVal(""), err
	}
	return MarshalRangeVal(t.Add(time.Duration(days) * time.Hour * 24)), nil
}

func MarshalRangeVal(t time.Time) RangeVal {
	return RangeVal(t.In(time.UTC).Format("20060102"))
}

func (v RangeVal) Validate() error {
	if len(v) != 8 {
		return errors.New("range value must be in the format YYYYMMDD")
	}
	if _, err := strconv.Atoi(string(v)); err != nil {
		return errors.New("range value must be in the format YYYYMMDD")
	}
	return nil
}

func (v RangeVal) Int() (int, error) {
	if err := v.Validate(); err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(string(v))
	if err != nil {
		return 0, errors.New("range value must be in the format YYYYMMDD")
	}
	return i, nil
}

func (r RangeVal) Unmarshal() (time.Time, error) {
	if err := r.Validate(); err != nil {
		return time.Time{}, err
	}
	v := string(r)
	y, err := strconv.Atoi(v[0:4])
	if err != nil {
		return time.Time{}, err
	}
	m, err := strconv.Atoi(v[4:6])
	if err != nil {
		return time.Time{}, err
	}
	d, err := strconv.Atoi(v[6:8])
	if err != nil {
		return time.Time{}, err
	}
	// Set central date as noon.
	return time.Date(y, time.Month(m), d, 12, 0, 0, 0, time.UTC), nil
}
