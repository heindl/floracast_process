package main

import (
	"bitbucket.org/heindl/taxa/ecoregions"
	"bitbucket.org/heindl/taxa/store"
	"flag"
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/saleswise/errors/errors"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
	"bitbucket.org/heindl/taxa/terra"
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

	ecoRegionCache, err := ecoregions.NewEcoRegionsCache()
	if err != nil {
		panic(err)
	}

	parser := &Parser{
		Store: taxaStore,
		Stats: NewStatsContainer(),
		EcoRegionCache: ecoRegionCache,
	}

	if err := filepath.Walk(*geojsonPath, parser.RecursiveSearchParse); err != nil {
		panic(err)
	}

	fmt.Println("Total", total)
	fmt.Println("Missing Ecoregion", missing_ecoregion)

	return
}

type Parser struct {
	Store          store.TaxaStore
	Stats          *StatsContainer
	EcoRegionCache ecoregions.EcoRegionsCache
}

var total = 0;
var missing_ecoregion = 0;

func (Ω *Parser) RecursiveSearchParse(path string, f os.FileInfo, err error) error {

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

	fc, err := terra.ReadFeatureCollectionFromGeoJSONFile(path)
	if err != nil {
		return err
	}

	total +=1

	ecoID := Ω.EcoRegionCache.EcoID(fc.PolyLabel().Latitude(), fc.PolyLabel().Longitude())
	if !ecoID.Valid() {
		missing_ecoregion += 1
	}

	return nil

}

func (p *Parser) FormatAreaName(name string) string {
	// Convert abbreviated suffixes to full names.
	name = strings.TrimSpace(strings.ToLower(name))
	for suffix, replacement := range map[string]string{
		"wa": "Wilderness Area",
		"sp": "State Park",
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

func (Ω *Parser) Parse(gb []byte) (*store.ProtectedArea, error) {


	return nil, nil

	//fc, err := geojson.UnmarshalFeatureCollection(gb)
	//if err != nil {
	//	return nil, errors.Wrap(err, "could not unmarshal field collection")
	//}
	//
	//f := fc.Features[0]



	//Ensure gap status code is being properly converted because that can signal how careful to be with the area
	//
	//pa := store.ProtectedArea{
	//	ID:                  strconv.Itoa(f.PropertyMustInt("WDPA_Cd")),
	//	StateAbbr:           f.PropertyMustString("State_Nm"),
	//	Category:            store.AreaCategory(f.PropertyMustString("Category")),     // Category, d_Category
	//	Designation:         store.AreaDesignation(f.PropertyMustString("Des_Tp")),    // d_Des_Tp, Des_Tp
	//	ManagerStandardName: store.AreaManagerName(f.PropertyMustString("Mang_Nam")),  //  d_Mang_Nam, Mang_Name
	//	ManagerLocalName:    f.PropertyMustString("Loc_Mang"),                         // Loc_Mang
	//	ManagerType:         store.AreaManagerType(f.PropertyMustString("Mang_Type")), // d_Mang_Typ, Mang_Type
	//	OwnerName:               store.AreaOwnerName(f.PropertyMustString("Own_Name")),    // Own_Name, d_Own_Name,
	//	OwnerNameLocal: 	f.PropertyMustString("Loc_Own"), //  Loc_Own
	//	OwnerType:           store.AreaOwnerType(f.PropertyMustString("Own_Type")),    // Own_Type, d_Own_Type
	//	PublicAccess:        store.AreaPublicAccess(f.PropertyMustString("Access")),
	//	IUCNCategory:        store.AreaIUCNCategory(f.PropertyMustString("IUCN_Cat")),
	//	AreaGAPStatus:       store.AreaGAPStatus(f.PropertyMustString("GAP_Sts")),
	//}

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

	//pa.NameStandard, err = Ω.FormatAreaName(pa.ID, f.PropertyMustString("Unit_Nm")) // Unit_Nm
	//if err != nil {
	//	return nil, err
	//}
	//
	//pa.NameLocal, err = Ω.FormatAreaName(pa.ID, f.PropertyMustString("Loc_Nm")) // Unit_Nm
	//if err != nil {
	//	return nil, err
	//}
	//
	//// Total GISAcres
	//for _, _f := range fc.Features {
	//	pa.GISAcres += _f.PropertyMustFloat64("GIS_Acres")
	//}


	//if pa.NameStandard != "Eastshore State Park" {
	//	return nil, nil
	//}



	// Another option https://github.com/mapbox/polylabel
	//centre, bounds, minDistanceFromNearestPoint, err := p.ParsePolygon(fc)
	//poly, centroid, err := Ω.ParsePolygon(fc)
	//if err != nil {
	//	return nil, err
	//}
	//pa.PolyLabel = [2]float64{poly.Lat, poly.Lng}
	//pa.Centroid = [2]float64{centroid.Lat, centroid.Lng}
	//
	//
	//
	//return nil, nil

	//return nil, nil

	//newLat, newLng, newContains, pastContains, err := SecondOpinionCentroid(gb, centroid.Lat(), centroid.Lng())
	//
	//distanceBetweenCentroids := centroid.GeoDistanceFrom(&geo.Point{newLng, newLat})
	//
	//
	//if newContains != pastContains {
	//	fmt.Println(utils.JsonOrSpew(map[string]interface{}{
	//		"intial_centroid":           []float64{centroid.Lat(), centroid.Lng()},
	//		"intial_contained":          pastContains,
	//		"second_centroid":           []float64{newLat, newLng},
	//		"second_contained":          newContains,
	//		"initial_min_from_boundary": minDistanceFromNearestPoint,
	//		"distance_between":          distanceBetweenCentroids,
	//	}))
	//}
	//
	//return nil, nil


	//pa.Centroid = [2]float64{centroid.Lat(), centroid.Lng()}
	//pa.Height = bounds.GeoHeight()
	//pa.Width = bounds.GeoWidth()
	//ne := bounds.NorthEast()
	//sw := bounds.SouthWest()
	//pa.Bounds = [2][2]float64{{sw.Lat(), sw.Lng()}, {ne.Lat(), ne.Lng()}}

	//Ω.EcoRegionCache.PointWithin(centroid.Lat(), centroid.Lng())

	//valid, reason, invalidValue := pa.Valid()
	//
	//if !valid {
	//	fmt.Println("invalid", reason, invalidValue)
	//	return nil, nil
	//}
	//
	//if err := Ω.PrintCSV(&pa); err != nil {
	//	return nil, err
	//}
	//
	//return &pa, nil
}



func (p *Parser) PrintCSV(area *store.ProtectedArea) error {
	s, err := gocsv.MarshalString(area)
	if err != nil {
		return errors.Wrap(err, "could not marshal csv")
	}
	fmt.Println(s)
	return nil
}
