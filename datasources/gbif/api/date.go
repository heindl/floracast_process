package api

import (
	"github.com/dropbox/godropbox/errors"
	"strings"
	"time"
)

type Time struct {
	time.Time
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return []byte(t.Format("2006-01-02T15:04:05.999-0700")), nil
}

func (t *Time) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	tiempo, err := time.Parse("2006-01-02T15:04:05.999-0700", strings.Trim(string(b), `"`))
	if err != nil {
		return errors.Wrapf(err, "could not parse time value: %s", string(b))
	}
	t.Time = tiempo
	return nil
}
