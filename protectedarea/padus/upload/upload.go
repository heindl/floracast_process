package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unicode/utf8"

	"bitbucket.org/heindl/process/protectedarea"
	"bitbucket.org/heindl/process/protectedarea/padus"
	"bitbucket.org/heindl/process/terra/geo"
	"context"
	"github.com/dropbox/godropbox/errors"
	"gopkg.in/tomb.v2"
)

func main() {

	geojsonPath := flag.String("geojson", "/tmp/gap_analysis", "Path to geojson files to search recursively.")

	if *geojsonPath == "" {
		panic("A geojson directory must be specified.")
	}

	cxt := context.Background()

	//florastore, err := store.NewFloraStore(cxt)
	//if err != nil {
	//	panic(err)
	//}

	parser := &Parser{
		Context: cxt,
		Areas:   protectedarea.ProtectedAreas{},
		Tmb:     tomb.Tomb{},
	}

	parser.Tmb.Go(func() error {
		return filepath.Walk(*geojsonPath, parser.RecursiveSearchParse)
	})

	if err := parser.Tmb.Wait(); err != nil {
		panic(err)
	}

	//if err := parser.Areas.Upload(cxt, florastore); err != nil {
	//	panic(err)
	//}
}

type Parser struct {
	Context context.Context
	Tmb     tomb.Tomb
	sync.Mutex
	Areas protectedarea.ProtectedAreas
}

func (Ω *Parser) RecursiveSearchParse(path string, f os.FileInfo, err error) error {

	Ω.Tmb.Go(func() error {
		if err != nil {
			return errors.Wrap(err, "passed error from file path walk")
		}

		if !isValidAreaFile(f) {
			return nil
		}

		fc, err := geo.ReadFeatureCollectionFromGeoJSONFile(path, nil)
		if err != nil {
			return err
		}

		area, err := parseFeatureCollection(fc)
		if err != nil {
			return err
		}

		Ω.Lock()
		defer Ω.Unlock()
		Ω.Areas = append(Ω.Areas, area)

		return nil
	})

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

	area, err := protectedarea.NewProtectedArea(polyLabel.Latitude(), polyLabel.Latitude(), areaMeters)
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

func parseFeature(area protectedarea.ProtectedArea, feature *geo.Feature, lock *sync.Mutex) error {

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

	name, err := formatAreaName(pa.UnitNm)
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

	if err := area.UpdateName(name); err != nil {
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

func formatAreaName(name string) (string, error) {
	// There are a few invalid UTF8 characters in a few of the names, which will not save to firestore.
	// This is is the easiest way to correct them.
	if !utf8.ValidString(name) {
		return "", errors.Newf("invalid utf8 character in area name [%s] which will not save to firestore", name)
		//switch id {
		//
		//case "373061":
		//	return "Peña Blanca", nil
		//case "11115456":
		//	return "Año Nuevo SR", nil
		//case "11115421":
		//	return "Año Nuevo SP", nil
		//case "11116355", "555607510", "11116283", "11116280":
		//	return "Montaña de Oro SP", nil
		//default:
		//	return "", errors.Newf("invalid utf8 character which will not save to firestore: %s, %s", id, name)
		//}
	}
	// Convert abbreviated suffixes to full names.
	name = strings.ToLower(name)
	name = strings.Replace(name, "_", " ", -1)
	// Standardize Spaces
	name = strings.Join(strings.Fields(name), " ")
	for suffix, replacement := range map[string]string{
		"wa":  "Wilderness Area",
		"sp":  "State Park",
		"wma": "Wildlife Management Area",
	} {
		if strings.HasSuffix(name, suffix) {
			name = strings.Replace(name, suffix, replacement, -1)
		}
	}

	return strings.Title(name), nil
}
