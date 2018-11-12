package cache

import (
	"github.com/heindl/floracast_process/protectedarea"
	"github.com/heindl/floracast_process/store"
	"github.com/heindl/floracast_process/terra/geoembed"
	"context"
	"github.com/dropbox/godropbox/errors"
	"github.com/golang/glog"
	"strings"
	"sync"
)

// ProtectedAreaCache fetches and stores ProtectedAreas locally in case the same record
// is requested multiple times by a process.
type ProtectedAreaCache interface {
	GetProtectedArea(cxt context.Context, latitude, longitude float64) (protectedarea.ProtectedArea, error)
	GetProtectedAreaWithToken(cxt context.Context, key geoembed.S2Key) (protectedarea.ProtectedArea, error)
}

// NewProtectedAreaCache creates a new ProtectedAreaCache
func NewProtectedAreaCache(florastore store.FloraStore) (ProtectedAreaCache, error) {
	return &cache{
		florastore: florastore,
		areas:      map[geoembed.S2Key]protectedarea.ProtectedArea{},
	}, nil
}

type cache struct {
	florastore store.FloraStore
	sync.Mutex
	areas map[geoembed.S2Key]protectedarea.ProtectedArea
}

func (Ω *cache) GetProtectedAreaWithToken(cxt context.Context, s2Key geoembed.S2Key) (protectedarea.ProtectedArea, error) {

	Ω.Lock()
	if v, ok := Ω.areas[s2Key]; ok {
		Ω.Unlock()
		return v, nil
	}
	Ω.Unlock()

	area, err := protectedarea.FetchOne(cxt, Ω.florastore, s2Key)
	if err != nil && !strings.Contains(err.Error(), "NotFound") {
		return nil, err
	}

	if err != nil && strings.Contains(err.Error(), "NotFound") {
		glog.Warningf("ProtectedArea not found for S2Token [%s]", s2Key)
		return nil, nil
	}

	Ω.Lock()
	defer Ω.Unlock()

	Ω.areas[s2Key] = area

	return area, nil
}

func (Ω *cache) GetProtectedArea(cxt context.Context, latitude, longitude float64) (protectedarea.ProtectedArea, error) {

	s2Key, err := geoembed.NewS2Key(latitude, longitude)
	if err != nil {
		return nil, err
	}

	Ω.Lock()
	if v, ok := Ω.areas[s2Key]; ok {
		Ω.Unlock()
		return v, nil
	}
	Ω.Unlock()

	area, err := protectedarea.FetchOne(cxt, Ω.florastore, s2Key)
	if err != nil && !strings.Contains(err.Error(), "NotFound") {
		return nil, errors.Wrapf(err, "Position [%f, %f]", latitude, longitude)
	}

	if err != nil && strings.Contains(err.Error(), "NotFound") {
		glog.Warningf("ProtectedArea not found for S2Token [%s] at Location [%f, %f]", s2Key, latitude, longitude)
		return nil, nil
	}

	Ω.Lock()
	defer Ω.Unlock()

	Ω.areas[s2Key] = area

	return area, nil
}
