package elevation

import (
	"bitbucket.org/heindl/process/terra/geo"
	"bitbucket.org/heindl/process/utils"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"strings"
	"sync"
)

var elevationBatchSize = 20

var global = processor{
	queued:  []string{},
	fetched: map[string]*int{},
}

// Queue an elevation to look as a batch.
func Queue(lat, lng float64) error {
	if err := geo.ValidateCoordinates(lat, lng); err != nil {
		return err
	}
	return global.queue(lat, lng)
}

// Get an elevation given a Latitude/Longitude coordinate.
// Must have been previously queued, and will flush cache if not already fetched.
func Get(lat, lng float64) (*int, error) {
	if err := geo.ValidateCoordinates(lat, lng); err != nil {
		return nil, err
	}
	return global.get(lat, lng)
}

type processor struct {
	sync.Mutex
	queued  []string
	fetched map[string]*int
}

func (Ω *processor) get(lat, lng float64) (*int, error) {

	Ω.Lock()
	defer Ω.Unlock()

	if Ω.isQueued(lat, lng) {
		if err := Ω.fetch(); err != nil {
			return nil, err
		}
	}

	if !Ω.isFetched(lat, lng) {
		return nil, errors.Newf("Lat/Lng is neither queued or fetched [%f, %f]", lat, lng)
	}

	return Ω.fetched[key(lat, lng)], nil

}

func (Ω *processor) queue(lat, lng float64) error {

	Ω.Lock()
	defer Ω.Unlock()

	if Ω.isQueued(lat, lng) || Ω.isFetched(lat, lng) {
		return nil
	}
	Ω.queued = append(Ω.queued, key(lat, lng))
	shouldFetch := len(Ω.queued) >= elevationBatchSize

	if shouldFetch {
		return Ω.fetch()
	}

	return nil
}

func (Ω *processor) isQueued(lat, lng float64) bool {
	return utils.ContainsString(Ω.queued, key(lat, lng))
}

func (Ω *processor) isFetched(lat, lng float64) bool {
	_, ok := Ω.fetched[key(lat, lng)]
	return ok
}

func key(lat, lng float64) string {
	return fmt.Sprintf("%.4f,%.4f", lat, lng)
}

var fetchCount = 0

func (Ω *processor) fetch() error {

	fetchCount++

	if len(Ω.queued) == 0 {
		return nil
	}

	if len(Ω.queued) > elevationBatchSize {
		return errors.Newf("Batch size [%d] greater than expected [%d]", len(Ω.queued), elevationBatchSize)
	}

	var res struct {
		Results []struct {
			Lat       float64 `json:"latitude"`
			Lng       float64 `json:"longitude"`
			Elevation int     `json:"elevation"` // Meters
		} `json:"results"`
	}

	if err := utils.RequestJSON("https://api.open-elevation.com/api/v1/lookup?locations="+strings.Join(Ω.queued, "|"), &res); err != nil {
		return errors.Wrap(err, "Could not fetch elevation api")
	}

	for _, r := range res.Results {
		k := key(r.Lat, r.Lng)
		if r.Elevation == 0 {
			fmt.Println(fmt.Sprintf("Elevation value is 0 for key [%s], so may not have resolved.", k))
		}
		Ω.fetched[k] = utils.IntPtr(r.Elevation)
	}

	Ω.queued = []string{}
	return nil
}
