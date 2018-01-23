package main

import (
	"bitbucket.org/heindl/taxa/pad_us"
	"bitbucket.org/heindl/taxa/store"
	"bitbucket.org/heindl/taxa/terra"
	"flag"
	"github.com/saleswise/errors/errors"
	"golang.org/x/net/context"
	"gopkg.in/tomb.v2"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"
	"fmt"
)

func main() {

	geojsonPath := flag.String("geojson", "/tmp/gap_analysis", "Path to geojson files to search recursively.")

	if *geojsonPath == "" {
		panic("A geojson directory must be specified.")
	}

	taxaStore, err := store.NewTaxaStore()
	if err != nil {
		panic(err)
	}

	parser := &Parser{
		Context: context.Background(),
		Holder:  []*store.ProtectedArea{},
		Tmb:     tomb.Tomb{},
	}

	parser.Tmb.Go(func() error {
		return filepath.Walk(*geojsonPath, parser.RecursiveSearchParse)
	})

	if err := parser.Tmb.Wait(); err != nil {
		panic(err)
	}

	fmt.Println("Setting", len(parser.Holder), "Protected Areas")

	if err := taxaStore.SetProtectedAreas(parser.Context, parser.Holder...); err != nil {
		panic(err)
	}

	return
}

type Parser struct {
	Context context.Context
	Tmb     tomb.Tomb
	sync.Mutex
	Holder []*store.ProtectedArea
}

func (Ω *Parser) RecursiveSearchParse(path string, f os.FileInfo, err error) error {

	Ω.Tmb.Go(func() error {
		if err != nil {
			return errors.Wrap(err, "passed error from file path walk")
		}

		if f.IsDir() {
			return nil
		}

		if strings.Contains(f.Name(), "state.geojson") {
			return nil
		}

		if !strings.HasSuffix(f.Name(), ".geojson") {
			return nil
		}

		fc, err := terra.ReadFeatureCollectionFromGeoJSONFile(path, nil)
		if err != nil {
			return err
		}

		// PushSynonym features into single store ProtectedArea.

		spa := store.ProtectedArea{
			Area:      fc.Area(),
			PolyLabel: fc.PolyLabel().AsArray(),
		}

		for _, feature := range fc.Features() {
			pa := pad_us.ProtectedArea{}
			if err := feature.GetProperties(&pa); err != nil {
				return err
			}
			gsc, err := strconv.Atoi(string(pa.GAPStatusCode))
			if err != nil {
				return err
			}
			if spa.ProtectionLevel == nil || store.ProtectionLevel(gsc-1) < *spa.ProtectionLevel {
				pl := store.ProtectionLevel(gsc - 1)
				spa.ProtectionLevel = &pl
			}

			access := store.AccessLevelUnknown
			switch string(pa.Access) {
			case "OA":
				access = store.AccessLevelOpen
			case "RA":
				access = store.AccessLevelRestricted
			case "UK":
				access = store.AccessLevelUnknown
			case "XA":
				access = store.AccessLevelClosed
			}

			if spa.AccessLevel == nil || access < *spa.AccessLevel {
				spa.AccessLevel = &access
			}

			name, err := Ω.EscapeAreaName(parse_string_value(spa.Name, pa.UnitNm))
			if err != nil {
				return err
			}
			spa.Name = Ω.FormatAreaName(name)

			spa.State = pa.StateNm

			spa.Designation = parse_string_value(spa.Designation, pa.DDesTp)
			spa.Owner = parse_string_value(spa.Owner, pa.DOwnName)
		}

		Ω.Lock()
		defer Ω.Unlock()
		Ω.Holder = append(Ω.Holder, &spa)

		return nil
	})

	return nil

}

func parse_string_value(existing, given string) string {

	if existing == given {
		return existing
	}

	if existing == "" && given != "" {
		return given
	}

	existing_has_negative_flag := false
	given_has_negative_flag := false
	for _, f := range []string{"unknown", "other", "easement", "private"} {
		if strings.Contains(strings.ToLower(existing), f) {
			existing_has_negative_flag = true
		}
		if strings.Contains(strings.ToLower(given), f) {
			given_has_negative_flag = true
		}
	}

	if existing_has_negative_flag != given_has_negative_flag && existing_has_negative_flag {
		return given
	} else {
		return existing
	}

	if len(existing) > len(given) {
		return existing
	} else {
		return given
	}

	return given
}

func (p *Parser) FormatAreaName(name string) string {
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
	return strings.Title(name)
}

func (p *Parser) EscapeAreaName(name string) (string, error) {
	// There are a few invalid UTF8 characters in a few of the names, which will not save to firestore.
	// This is is the easiest way to correct them.
	if !utf8.ValidString(name) {
		return "", errors.Newf("invalid utf8 character which will not save to firestore: %s", name)
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

	return name, nil
}
