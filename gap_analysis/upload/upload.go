package main

import (
	"github.com/saleswise/errors/errors"
	"strings"
	"fmt"
	"bitbucket.org/heindl/taxa/store"
	"google.golang.org/genproto/googleapis/type/latlng"
	"unicode/utf8"
	"github.com/paulmach/go.geojson"
	"github.com/paulmach/go.geo"
	"io/ioutil"
	"flag"
	"path/filepath"
	"os"
)


func main() {

	geojsonPath := flag.String("geojson", "/tmp/gap_analysis", "Path to geojson files to search recursively.")

	flag.Parse()

	if *geojsonPath == "" {
		panic("A geojson directory must be specified.")
	}

	store, err := store.NewTaxaStore()
	if err != nil {
		panic(err)
	}

	parser := &Parser{Store: store}


	if err := filepath.Walk(*geojsonPath, parser.RecursiveSearchParse); err != nil {
		panic(err)
	}



	return
}

type Parser struct {
	Store store.TaxaStore
}

func (p *Parser) RecursiveSearchParse(path string, f os.FileInfo, err error) error {

	fmt.Println(path, f.Name(), err)

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

	fmt.Println(pa.Name)

	return nil

}

/* GEOJSON KEY */

//"d_State_Nm",
//"d_GAP_Sts",
//"GIS_Acres",
//"Loc_Nm",
//"WDPA_Cd", // World Database of Protected Areas Code, https://www.protectedplanet.net/
//"d_Mang_Nam", // Manager Name, "Natural Resources Conservation Service"
//"d_Mang_Typ", // Management Type, "Federal"
//"d_Own_Type", // Owner type, "Private"
//"d_Own_Name", // Owner name
//"d_Des_Tp", // The unit’s land management description or designation, standardized for nation (e.g. Area of Critical Environmental Concern, Wilderness Area, State Park, Local Rec Area, Conservation Easement). See the PAD-US Data Standard for a crosswalk of "Designation Type" from source data files or the geodatabase look up table for "Designation Type" for domain descriptions. "Designation Type" supports PAD-US queries and categorical conservation measures or public access assignments in the absence of other information.
//"Date_Est", // The Year (yyyy) the protected area was designated, decreed or otherwise established. Date is assigned to each unit by name, without event status(e.g. Yellowstone National Park: 1872, Frank Church-River of No Return Wilderness Area: 1980)
//"d_Access", // Level of public access permitted. Open requires no special requirements for public access to the property (may include regular hours available); Restricted requires a special permit from the owner for access, a registration permit on public land or has highly variable times when open to use; Closed occurs where no public access allowed (land bank property, special ecological study areas, military bases, etc. Unknown is assigned where information is not currently available. Access is assigned categorically by Designation Type or provided by PAD-US State Data Stewards, federal or NGO partners. Contact the PAD-US Coordinator with available public access information.
//"d_Category", // General category for the protection mechanism associated with the protected area. ‘Fee’ is the most common way real estate is owned. A conservation ‘easement’ creates a legally enforceable land preservation agreement between a landowner and government agency or qualified land protection organization (i.e. land trust). ‘Other’ types of protection include leases, agreements or those over marine waters. ‘Designation’ is applied to designations in the federal theme not tied to title documents (e.g. National Monument, Wild and Scenic River). These may be removed to reduce overlaps for area based analyses.

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

func (p *Parser) Parse(gb []byte) (*store.ProtectedArea, error) {

	fc, err := geojson.UnmarshalFeatureCollection(gb)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal field collection")
	}

	if len(fc.Features) > 1 {
		fmt.Println("have feature length greater than one")
	}

	f := fc.Features[0]

	pa := store.ProtectedArea{
		ID: f.PropertyMustString("WDPA_Cd"),
		PublicAccess: f.PropertyMustString("d_Access"),
		GapAnalysisProjectStatus: f.PropertyMustString("d_GAP_Sts"),
		Acres: f.PropertyMustFloat64("GIS_Acres"),
		ManagementDesignation: f.PropertyMustString("d_Des_Tp"),
		ManagerType: f.PropertyMustString("d_Mang_Typ"),
		ManagerName: f.PropertyMustString("d_Mang_Nam"),
		Category: f.PropertyMustString("d_Category"),
		OwnerType: f.PropertyMustString("d_Own_Type"),
		OwnerName: f.PropertyMustString("d_Own_Name"),
		YearEstablished: f.PropertyMustInt("Date_Est"),
	}

	pa.Name, err = p.FormatAreaName(pa.ID, f.PropertyMustString("Loc_Nm"))
	if err != nil {
		return nil, err
	}

	pa.State = store.ProtectedAreaState{
		Name: f.PropertyMustString("d_State_Nm"),
	}

	for k, v := range store.ValidProtectedAreaStates {
		if v == pa.State.Name {
			pa.State.Abbr = k
		}
	}

	// Another option https://github.com/mapbox/polylabel
	centre, bounds, err := p.ParsePolygon(f)
	if err != nil {
		return nil, err
	}

	pa.Centre = latlng.LatLng{centre.Lat(), centre.Lng()}
	pa.HeightMeters = bounds.GeoHeight()
	pa.WidthMeters = bounds.GeoWidth()
	pa.Bounds = bounds.String()

	if valid, _ := pa.Valid(); !valid {
		return nil, nil
	}

	return &pa, nil
}

func (p *Parser) ParsePolygon(f *geojson.Feature) (centre *geo.Point, bounds *geo.Bound, err error) {
		if !f.Geometry.IsMultiPolygon() && !f.Geometry.IsPolygon() {
			return nil, nil, errors.Newf("unsupported geometry type: %s", f.Geometry.Type)
		}

		pointSet := geo.NewPointSet()

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

		min, i := pointSet.GeoDistanceFrom(pointSet.GeoCentroid())
		fmt.Println("Distance From", f.PropertyMustString("WDPA_Cd"), min, pointSet.GetAt(i))

		return pointSet.GeoCentroid(), pointSet.Bound(), nil

}