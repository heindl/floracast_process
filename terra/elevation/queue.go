package elevation

import (
	"github.com/heindl/floracast_process/utils"
	"bytes"
	"encoding/json"
	"github.com/dropbox/godropbox/errors"
	"github.com/golang/glog"
	"google.golang.org/genproto/googleapis/type/latlng"
	"gopkg.in/tomb.v2"
	"sync"
	"time"
)

type queue struct {
	sync.Mutex
	m           map[token]struct{}
	ll          []*latlng.LatLng
	lastFetched time.Time
	length      int
}

func (Ω *queue) count() int {
	return Ω.length
}

func (Ω *queue) has(t token) bool {
	Ω.Lock()
	defer Ω.Unlock()
	_, ok := Ω.m[t]
	return ok
}

func (Ω *queue) add(lat, lng float64) {
	Ω.Lock()
	defer Ω.Unlock()
	if Ω.m == nil {
		Ω.m = map[token]struct{}{}
	}
	t := toToken(lat, lng)
	_, ok := Ω.m[t]
	if ok {
		return
	}
	Ω.m[t] = struct{}{}
	Ω.length = len(Ω.m)
	Ω.ll = append(Ω.ll, &latlng.LatLng{Latitude: lat, Longitude: lng})
}

type elevationPostRequest struct {
	Locations []*latlng.LatLng `json:"locations"`
}

type postResult struct {
	*latlng.LatLng
	Elevation float64 `json:"elevation"`
}

func (Ω *queue) fetch() (map[token]*int, error) {
	Ω.Lock()
	defer Ω.Unlock()

	glog.Infof("Fetching Elevation Batch [%d]", Ω.length)

	if len(Ω.ll) == 0 {
		return nil, nil
	}

	locker := sync.Mutex{}
	limiter := utils.NewLimiter(5)
	combined := []*postResult{}
	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _i := 0; len(Ω.ll) > _i*200; _i++ {
			release := limiter.Go()
			i := _i
			start := i * 200
			end := (i + 1) * 200
			if end > len(Ω.ll) {
				end = len(Ω.ll)
			}
			tmb.Go(func() error {
				defer release()
				b, err := json.Marshal(elevationPostRequest{
					Locations: Ω.ll[start:end],
				})
				if err != nil {
					return errors.Wrap(err, "could not marshal elevation post request")
				}
				postData, err := utils.PostJSON("https://api.open-elevation.com/api/v1/lookup", bytes.NewReader(b))
				if err != nil {
					return err
				}
				postResult := struct {
					Results []*postResult `json:"results"`
				}{}
				if err := json.Unmarshal(postData, &postResult); err != nil {
					return errors.Wrap(err, "could not marshal elevation post result")
				}
				locker.Lock()
				defer locker.Unlock()
				combined = append(combined, postResult.Results...)
				return nil
			})
		}
		return nil
	})
	if err := tmb.Wait(); err != nil {
		return nil, err
	}

	glog.Infof("Received Elevation Results [%d]", len(combined))

	Ω.lastFetched = time.Now()

	res := map[token]*int{}
	for _, r := range combined {
		t := toToken(r.Latitude, r.Longitude)
		delete(Ω.m, t)
		res[t] = utils.IntPtr(int(r.Elevation))
	}
	// Ensure even those that were not returned are accounted for in the results.
	for t := range Ω.m {
		res[t] = nil
	}

	Ω.m = nil
	Ω.length = 0
	Ω.ll = []*latlng.LatLng{}

	return res, nil
}
