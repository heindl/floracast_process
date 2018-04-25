package elevation

import (
	"bitbucket.org/heindl/process/terra/geo"
	"github.com/golang/geo/s2"
	"gopkg.in/tomb.v2"
	"sync"
	"time"
)

var globalProcessor *processor

// Queue an elevation to look as a batch.
func Queue(lat, lng float64) error {

	if globalProcessor == nil {
		globalProcessor = &processor{
			Queue: &queue{},
			Cache: map[token]*int{},
			Tmb:   &tomb.Tomb{},
		}
		globalProcessor.Tmb.Go(globalProcessor.monitor)
	}

	if err := globalProcessor.Tmb.Err(); err != tomb.ErrStillAlive {
		return err
	}

	if err := geo.ValidateCoordinates(lat, lng); err != nil {
		return err
	}

	globalProcessor.Tmb.Go(func() error {
		globalProcessor.Queue.add(lat, lng)
		return nil
	})
	return nil
}

// Get an elevation given a Latitude/Longitude coordinate.
// Must have been previously queued, and will flush cache if not already fetched.
func Get(lat, lng float64) (*int, error) {
	if err := geo.ValidateCoordinates(lat, lng); err != nil {
		return nil, err
	}
	t := toToken(lat, lng)
	for {
		v, ok, err := globalProcessor.get(t)
		if err != nil {
			return nil, err
		}
		if ok {
			return v, nil
		}
		time.Sleep(time.Second * 2)
	}
}

type processor struct {
	sync.Mutex
	Queue *queue
	Cache map[token]*int
	Tmb   *tomb.Tomb
}

func (Ω *processor) monitor() error {
	for {
		time.Sleep(time.Second)
		if Ω.Queue.count() == 0 {
			continue
		}
		if Ω.Queue.count() >= 200 || time.Now().Sub(Ω.Queue.lastFetched).Seconds() > 10 {
			if err := Ω.fetch(); err != nil {
				return err
			}
		}
	}
}

func (Ω *processor) get(t token) (*int, bool, error) {
	if err := Ω.Tmb.Err(); err != tomb.ErrStillAlive {
		return nil, false, err
	}
	Ω.Lock()
	defer Ω.Unlock()
	v, ok := Ω.Cache[t]
	return v, ok, nil
}

func (Ω *processor) fetch() error {

	Ω.Lock()
	defer Ω.Unlock()

	res, err := Ω.Queue.fetch()
	if err != nil {
		return err
	}

	for k, v := range res {
		Ω.Cache[k] = v
	}

	return nil

}

type token string

func toToken(lat, lng float64) token {
	return token(s2.CellIDFromLatLng(s2.LatLngFromDegrees(lat, lng)).Parent(10).ToToken())
}
