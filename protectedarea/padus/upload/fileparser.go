package main

import (
	"github.com/heindl/floracast_process/protectedarea"
	"github.com/heindl/floracast_process/protectedarea/padus"
	"github.com/heindl/floracast_process/terra/ecoregions"
	"github.com/heindl/floracast_process/terra/geo"
	"github.com/heindl/floracast_process/utils"
	"github.com/dropbox/godropbox/errors"
	"gopkg.in/tomb.v2"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// ParseProtectedAreaDirectory reads all .geojson files in a directory and returns a list of ProtectedAreas
// Each file is expected to contain one FeatureCollection that becomes one ProtectedArea
func ParseProtectedAreaDirectory(directoryPath string) (protectedarea.ProtectedAreas, error) {
	w := walker{
		areas: protectedarea.ProtectedAreas{},
	}
	if err := filepath.Walk(directoryPath, w.recursiveSearchParse); err != nil {
		return nil, err
	}
	return w.areas, nil
}

type walker struct {
	sync.Mutex
	areas protectedarea.ProtectedAreas
}

func (Ω *walker) recursiveSearchParse(path string, f os.FileInfo, err error) error {

	if err != nil {
		return errors.Wrap(err, "passed error from file path walk")
	}

	if !isValidAreaFile(f) {
		return nil
	}

	fc, err := geo.ReadFeatureCollectionFromGeoJSONFile(path, nil)
	if err != nil {
		return errors.Wrapf(err, "Could not read FeatureCollection [%s]", path)
	}

	area, err := parseFeatureCollection(fc)
	if err != nil {
		return err
	}
	if area == nil {
		return nil
	}

	Ω.Lock()
	defer Ω.Unlock()
	Ω.areas = append(Ω.areas, area)

	return nil

}

func isValidAreaFile(f os.FileInfo) bool {
	return !f.IsDir() &&
		!strings.Contains(f.Name(), "state.geojson") &&
		strings.HasSuffix(f.Name(), ".geojson")
}

func parseFeatureCollection(fc *geo.FeatureCollection) (protectedarea.ProtectedArea, error) {

	polyLabel, err := fc.PolyLabel()
	if err != nil {
		return nil, err
	}

	areaMeters := fc.Area()
	if areaMeters == 0 {
		return nil, errors.New("Invalid Area [0] for ProtectedArea FeatureCollection")
	}

	area, err := protectedarea.NewProtectedArea(polyLabel.Latitude(), polyLabel.Longitude(), areaMeters)
	if utils.ContainsError(err, ecoregions.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	lock := sync.Mutex{}
	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, 𝝨 := range fc.Features() {
			feature := 𝝨
			tmb.Go(func() error {
				return parseFeature(area, feature, &lock)
			})
		}
		return nil
	})
	if err := tmb.Wait(); err != nil {
		return nil, err
	}

	return area, nil
}

func parseFeature(area protectedarea.ProtectedArea, feature *geo.Feature, lock sync.Locker) error {

	pa := padus.ProtectedArea{}
	if err := feature.GetProperties(&pa); err != nil {
		return err
	}

	protectionLevelInt, err := pa.GAPStatusCode.ProtectionLevel()
	if err != nil {
		return err
	}

	accessLevelInt, err := pa.Access.AccessLevel()
	if err != nil {
		return err
	}

	lock.Lock()
	defer lock.Unlock()

	if err := area.UpdateProtectionLevel(protectionLevelInt); err != nil {
		return err
	}

	if err := area.UpdateAccessLevel(accessLevelInt); err != nil {
		return err
	}

	if err := area.UpdateName(pa.UnitNm); err != nil {
		return err
	}

	if err := area.UpdateOwner(pa.DOwnName); err != nil {
		return err
	}

	if err := area.UpdateDesignation(pa.DDesTp); err != nil {
		return err
	}

	return nil

}
