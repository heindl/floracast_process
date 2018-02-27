package main

import (
	//"bitbucket.org/heindl/process/store"
	"flag"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"

	"bitbucket.org/heindl/process/geofeatures"
	"bitbucket.org/heindl/process/protected_areas"
	"bitbucket.org/heindl/process/protected_areas/pad_us"
	"bitbucket.org/heindl/process/terra"
	"github.com/dropbox/godropbox/errors"
	"golang.org/x/net/context"
	"gopkg.in/tomb.v2"
)

func main() {

	geojsonPath := flag.String("geojson", "/tmp/gap_analysis", "Path to geojson files to search recursively.")

	if *geojsonPath == "" {
		panic("A geojson directory must be specified.")
	}

	//florastore, err := store.NewFloraStore(context.Background())
	//if err != nil {
	//	panic(err)
	//}

	cxt := context.Background()

	parser := &Parser{
		Context: cxt,
		Areas:   protected_areas.ProtectedAreas{},
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

	return
}

type Parser struct {
	Context context.Context
	Tmb     tomb.Tomb
	sync.Mutex
	Areas protected_areas.ProtectedAreas
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

		gProtectionLevel := protected_areas.ProtectionLevelUnknown
		gAccessLevel := protected_areas.AccessLevelUnknown
		gName := ""
		gState := ""
		gDesignation := ""
		gOwner := ""

		for _, feature := range fc.Features() {
			pa := pad_us.ProtectedArea{}
			if err := feature.GetProperties(&pa); err != nil {
				return err
			}
			gsc, err := strconv.Atoi(string(pa.GAPStatusCode))
			if err != nil {
				return err
			}
			if protected_areas.ProtectionLevel(gsc-1) < gProtectionLevel {
				gProtectionLevel = protected_areas.ProtectionLevel(gsc - 1)
			}
			access := protected_areas.AccessLevelUnknown
			switch string(pa.Access) {
			case "OA":
				access = protected_areas.AccessLevelOpen
			case "RA":
				access = protected_areas.AccessLevelRestricted
			case "UK":
				access = protected_areas.AccessLevelUnknown
			case "XA":
				access = protected_areas.AccessLevelClosed
			}

			if gAccessLevel == protected_areas.AccessLevelUnknown && access != gAccessLevel || access < gAccessLevel {
				gAccessLevel = access
			}

			name, err := Ω.EscapeAreaName(parse_string_value(gName, pa.UnitNm))
			if err != nil {
				return err
			}
			gName = Ω.FormatAreaName(name)

			gState = pa.StateNm

			gDesignation = parse_string_value(gDesignation, pa.DDesTp)
			gOwner = parse_string_value(gOwner, pa.DOwnName)
		}

		area := protected_areas.ProtectedArea{
			Name:            gName,
			State:           gState,
			ProtectionLevel: &gProtectionLevel,
			AccessLevel:     &gAccessLevel,
			Designation:     gDesignation,
			Owner:           gOwner,
		}

		polyLabel, err := fc.PolyLabel()
		if err != nil {
			return err
		}

		area.GeoFeatureSet, err = geofeatures.NewGeoFeatureSet(polyLabel.Latitude(), polyLabel.Longitude(), false)
		if err != nil {
			return err
		}

		Ω.Lock()
		defer Ω.Unlock()
		Ω.Areas = append(Ω.Areas, &area)

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

	if (existing_has_negative_flag != given_has_negative_flag && existing_has_negative_flag) || (len(existing) < len(given)) {
		return given
	}
	return existing
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
