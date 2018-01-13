package main

import (
	"bitbucket.org/heindl/taxa/ecoregions"
	"bitbucket.org/heindl/taxa/store"
	"bitbucket.org/heindl/taxa/utils"
	"flag"
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/paulmach/go.geo"
	"github.com/paulmach/go.geojson"
	"github.com/saleswise/errors/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
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

	//ecoRegionCache, err := ecoregions.NewEcoRegionCache(*ecoRegionFile)
	//if err != nil {
	//	panic(err)
	//}

	parser := &Parser{
		Store: taxaStore,
		Stats: NewStatsContainer(),
		//EcoRegionCache: ecoRegionCache,
	}

	return

	if err := filepath.Walk(*geojsonPath, parser.RecursiveSearchParse); err != nil {
		panic(err)
	}

	return
}

type Parser struct {
	Store          store.TaxaStore
	Stats          *StatsContainer
	EcoRegionCache ecoregions.EcoRegionCache
}

func (p *Parser) RecursiveSearchParse(path string, f os.FileInfo, err error) error {

	if err != nil {
		return errors.Wrap(err, "passed error from file path walk")
	}

	if f.IsDir() {
		return nil
	}

	if !strings.HasSuffix(f.Name(), ".feature_collection.geojson") {
		return nil
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrapf(err, "could not read file [%s]", path)
	}

	pa, err := p.Parse(b)
	if err != nil {
		return err
	}

	if pa == nil {
		return nil
	}

	return nil

}

func (p *Parser) FormatAreaName(id, name string) (string, error) {
	// There are a few invalid UTF8 characters in a few of the names, which will not save to firestore.
	// This is is the easiest way to correct them.
	if !utf8.ValidString(name) {
		switch id {
		case "373061":
			return "Peña Blanca", nil
		case "11115456":
			return "Año Nuevo SR", nil
		case "11115421":
			return "Año Nuevo SP", nil
		case "11116355", "555607510", "11116283", "11116280":
			return "Montaña de Oro SP", nil
		default:
			return "", errors.Newf("invalid utf8 character which will not save to firestore: %s, %s", id, name)
		}
	}
	return name, nil
}

const (
	CategoryFee         = "Fee"
	CategoryDesignation = "Designation"
	CategoryEasement    = "Easement"
	CategoryOther       = "Other"
	CategoryUnknown     = "Unknown"
)

func (Ω *Parser) Parse(gb []byte) (*store.ProtectedArea, error) {

	fc, err := geojson.UnmarshalFeatureCollection(gb)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal field collection")
	}

	f := fc.Features[0]

	pa := store.ProtectedArea{
		ID:                  strconv.Itoa(f.PropertyMustInt("WDPA_Cd")),
		StateAbbr:           f.PropertyMustString("State_Nm"),
		Category:            store.AreaCategory(f.PropertyMustString("Category")),     // Category, d_Category
		Designation:         store.AreaDesignation(f.PropertyMustString("Des_Tp")),    // d_Des_Tp, Des_Tp
		ManagerStandardName: store.AreaManagerName(f.PropertyMustString("Mang_Nam")),  //  d_Mang_Nam, Mang_Name
		ManagerLocalName:    f.PropertyMustString("Loc_Mang"),                         // Loc_Mang
		ManagerType:         store.AreaManagerType(f.PropertyMustString("Mang_Type")), // d_Mang_Typ, Mang_Type
		Owner:               store.AreaOwnerName(f.PropertyMustString("Own_Name")),    // Own_Name, d_Own_Name, Loc_Own
		OwnerType:           store.AreaOwnerType(f.PropertyMustString("Own_Type")),    // Own_Type, d_Own_Type
		PublicAccess:        store.AreaPublicAccess(f.PropertyMustString("Access")),
		IUCNCategory:        store.AreaIUCNCategory(f.PropertyMustString("IUCN_Cat")),
		AreaGAPStatus:       store.AreaGAPStatus(f.PropertyMustString("GAP_Sts")),
	}

	// d_GAP_Sts

	//if !utils.Contains([]string{CategoryFee, CategoryDesignation, CategoryEasement, CategoryOther, CategoryUnknown}, pa.Category) {
	//	fmt.Println("New Category", pa.Category)
	//}

	//if _, ok := designs[f.PropertyMustString("Mang_Nam")]; !ok {
	//	designs[f.PropertyMustString("d_Mang_Nam")] = f.PropertyMustString("d_Mang_Nam")
	//} else {
	//	designs[f.PropertyMustString("d_Mang_Nam")] += 1
	//}

	//fmt.Println(
	//	"designation",
	//	f.PropertyMustString("d_Des_Tp"), ",",
	//	//f.PropertyMustString("Loc_Ds"), ",",
	//	f.PropertyMustString("Des_Tp"), ",",
	//)

	pa.NameStandard, err = Ω.FormatAreaName(pa.ID, f.PropertyMustString("Unit_Nm")) // Unit_Nm
	if err != nil {
		return nil, err
	}

	pa.NameLocal, err = Ω.FormatAreaName(pa.ID, f.PropertyMustString("Loc_Nm")) // Unit_Nm
	if err != nil {
		return nil, err
	}

	// Total GISAcres
	for _, _f := range fc.Features {
		pa.GISAcres += _f.PropertyMustFloat64("GIS_Acres")
	}

	// Another option https://github.com/mapbox/polylabel
	//centre, bounds, minDistanceFromNearestPoint, err := p.ParsePolygon(fc)
	centroid, bounds, _, err := Ω.ParsePolygon(fc)
	if err != nil {
		return nil, err
	}

	pa.Centroid = [2]float64{centroid.Lat(), centroid.Lng()}
	pa.Height = bounds.GeoHeight()
	pa.Width = bounds.GeoWidth()
	ne := bounds.NorthEast()
	sw := bounds.SouthWest()
	pa.Bounds = [2][2]float64{{sw.Lat(), sw.Lng()}, {ne.Lat(), ne.Lng()}}

	//Ω.EcoRegionCache.PointWithin(centroid.Lat(), centroid.Lng())

	valid, reason, invalidValue := pa.Valid()

	if !valid {
		fmt.Println("invalid", reason, invalidValue)
		return nil, nil
	}

	if err := Ω.PrintCSV(&pa); err != nil {
		return nil, err
	}

	return &pa, nil
}

var AccumulatedDefinitions = map[string]map[string]string{}
var AccumulatedCounters = map[string]map[string]int{}

func (p *Parser) AccumulateMaps(f *geojson.Feature) {
	for _, a := range [][]string{
		{"AreaCategory", "Category", "d_Category"},
		{"AreaDesignation", "Des_Tp", "d_Des_Tp"},
		{"AreaManagerName", "Mang_Name", "d_Mang_Nam"},
		{"AreaManagerType", "Mang_Type", "d_Mang_Typ"},
		{"AreaOwnerName", "Own_Name", "d_Own_Name"},
		{"AreaOwnerType", "Own_Type", "d_Own_Type"},
		{"AreaPublicAccess", "Access", "d_Access"},
		{"AreaIUCNCategory", "IUCN_Cat", "d_IUCN_Cat"},
		{"AreaGAPStatus", "GAP_Sts", "d_GAP_Sts"},
	} {
		if _, ok := AccumulatedDefinitions[a[0]]; !ok {
			AccumulatedDefinitions[a[0]] = map[string]string{}
		}
		if _, ok := AccumulatedDefinitions[a[0]][f.PropertyMustString(a[1])]; !ok {
			AccumulatedDefinitions[a[0]][f.PropertyMustString(a[1])] = f.PropertyMustString(a[2])
		}
	}
	if f.PropertyMustString("Unit_Nm") != f.PropertyMustString("Loc_Nm") {
		fmt.Println("{", f.PropertyMustString("Unit_Nm"), " : ", f.PropertyMustString("Loc_Nm"), "}")
	}

	for _, a := range [][]string{
		{"GapStatus", "GAP_Sts"},
		{"IUCNCategory", "IUCN_Cat"},
		{"PublicAccess", "Access"},
		{"Category", "Category"},
	} {
		if _, ok := AccumulatedCounters[a[0]]; !ok {
			AccumulatedCounters[a[0]] = map[string]int{}
		}
		v := f.PropertyMustString(a[1])
		if _, ok := AccumulatedCounters[a[0]][v]; !ok {
			AccumulatedCounters[a[0]][v] = 1
		} else {
			AccumulatedCounters[a[0]][v] += 1
		}
	}
	return
}

func (p *Parser) PrintAccumulatedMaps() {
	for a, b := range AccumulatedDefinitions {
		keys := []string{}
		for c, _ := range b {
			keys = append(keys, c)
		}
		sort.Strings(keys)
		fmt.Println(fmt.Sprintf("type %s string", a))
		fmt.Println(fmt.Sprintf("var %sDict = map[%s]string{", a, a))
		for _, k := range keys {
			fmt.Println(fmt.Sprintf(`%s("%s"): "%s",`, a, k, b[k]))
		}
		fmt.Println(fmt.Sprintf("}"))
	}

	fmt.Println(utils.JsonOrSpew(AccumulatedCounters))
}

func (p *Parser) ParsePolygon(fc *geojson.FeatureCollection) (centre *geo.Point, bounds *geo.Bound, minDistanceFromNearestPoint float64, err error) {

	pointSet := geo.NewPointSet()

	for _, f := range fc.Features {
		if !f.Geometry.IsMultiPolygon() && !f.Geometry.IsPolygon() {
			return nil, nil, 0, errors.Newf("unsupported geometry type: %s", f.Geometry.Type)
		}

		// MultiPolygon    [][][][]float64
		if f.Geometry.IsMultiPolygon() {
			for _, a := range f.Geometry.MultiPolygon {
				// [][][]float64
				for _, b := range a {
					// [][]float64
					for _, c := range b {
						pointSet = pointSet.Push(geo.NewPointFromLatLng(c[1], c[0]))
					}
				}
			}
		}
		// Polygon         [][][]float64
		if f.Geometry.IsPolygon() {
			for _, a := range f.Geometry.Polygon {
				// [][]float64
				for _, b := range a {
					// []float64
					pointSet = pointSet.Push(geo.NewPointFromLatLng(b[1], b[0]))
				}
			}
		}

	}
	centroid := pointSet.GeoCentroid()
	min, _ := pointSet.GeoDistanceFrom(centroid)
	bound := pointSet.Bound()

	return centroid, bound, min, nil

}

func (p *Parser) PrintCSV(area *store.ProtectedArea) error {
	s, err := gocsv.MarshalString(area)
	if err != nil {
		return errors.Wrap(err, "could not marshal csv")
	}
	fmt.Println(s)
	return nil
}
