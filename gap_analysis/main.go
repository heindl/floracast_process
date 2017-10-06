package main

import (
	"bitbucket.org/heindl/taxa/utils"
	"github.com/jonas-p/go-shp"
	"sort"
	"strconv"
	"github.com/saleswise/errors/errors"
	"github.com/kellydunn/golang-geo"
	"strings"
	"fmt"
	"bitbucket.org/heindl/taxa/store"
	"context"
	"google.golang.org/genproto/googleapis/type/latlng"
)

// https://gapanalysis.usgs.gov/padus/data/metadata/

// Download shape file:
// "the PAD-US database strives to be a complete “best available" inventory of areas dedicated to the
// preservation of biological diversity, and other natural, recreational or cultural uses, managed for
// these purposes through legal or other effective means."

// The data is updated annually.

// // https://gapanalysis.usgs.gov/padus/data/download/
// https://gapanalysis.usgs.gov/padus/data/statistics/

// May be best for big query because the data only changes once a year.

func main() {

	store, err := store.NewTaxaStore()
	if err != nil {
		panic(err)
	}

	parser := &Parser{store}

	areas, err := parser.Parse("./shapefiles/parsed.shp", 20)
	if err != nil {
		panic(err)
	}

	fmt.Println("len", len(areas))

	for _, a := range areas {
		if err := store.UpsertWildernessArea(context.Background(), a); err != nil {
			panic(err)
		}
	}



}

type Parser struct {
	Store store.TaxaStore
}

type State string

func (p *Parser) Parse(filename string, perstate int) (store.WildernessAreas, error) {
	reader, err := shp.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, "could not open shape file")
	}
	defer reader.Close()

	//for _, f := range reader.Fields() {
	//	fmt.Println(f)
	//}
	//
	//return nil, nil

	mapStateAreas := map[store.WildernessAreaState]store.WildernessAreas{}
	for _, v := range stateMap {
		mapStateAreas[store.WildernessAreaState(v)] = store.WildernessAreas{}
	}

	// FieldIDs are a map from known field names to their shapefile keys
	attrs := map[string]int{}
	for k, f := range reader.Fields() {
		for _, s := range []string{
			"d_State_Nm",
			"d_GAP_Sts",
			"GIS_Acres",
			"Loc_Nm",
			"WDPA_Cd", // World Database of Protected Areas Code, https://www.protectedplanet.net/
			"d_Mang_Nam", // Manager Name, "Natural Resources Conservation Service"
			"d_Mang_Typ", // Management Type, "Federal"
			"d_Own_Type", // Owner type, "Private"
			"d_Own_Name", // Owner name
			"d_Des_Tp", // The unit’s land management description or designation, standardized for nation (e.g. Area of Critical Environmental Concern, Wilderness Area, State Park, Local Rec Area, Conservation Easement). See the PAD-US Data Standard for a crosswalk of "Designation Type" from source data files or the geodatabase look up table for "Designation Type" for domain descriptions. "Designation Type" supports PAD-US queries and categorical conservation measures or public access assignments in the absence of other information.
			"Date_Est", // The Year (yyyy) the protected area was designated, decreed or otherwise established. Date is assigned to each unit by name, without event status(e.g. Yellowstone National Park: 1872, Frank Church-River of No Return Wilderness Area: 1980)
			"d_Access", // Level of public access permitted. Open requires no special requirements for public access to the property (may include regular hours available); Restricted requires a special permit from the owner for access, a registration permit on public land or has highly variable times when open to use; Closed occurs where no public access allowed (land bank property, special ecological study areas, military bases, etc. Unknown is assigned where information is not currently available. Access is assigned categorically by Designation Type or provided by PAD-US State Data Stewards, federal or NGO partners. Contact the PAD-US Coordinator with available public access information.
			"d_Category", // General category for the protection mechanism associated with the protected area. ‘Fee’ is the most common way real estate is owned. A conservation ‘easement’ creates a legally enforceable land preservation agreement between a landowner and government agency or qualified land protection organization (i.e. land trust). ‘Other’ types of protection include leases, agreements or those over marine waters. ‘Designation’ is applied to designations in the federal theme not tied to title documents (e.g. National Monument, Wild and Scenic River). These may be removed to reduce overlaps for area based analyses.
		} {
			if f.String() == s {
				attrs[s] = k
			}
		}
	}

	filterCounter := map[string]int{
		"d_GAP_Sts": 0,
		"d_Access": 0,
		"GIS_Acres": 0,
		"d_Des_Tp": 0,
		"d_Mang_Typ": 0,
		"WDPA_Cd": 0,
	}

	// Not certain why there are duplicates, except that perhaps different sources are not merged correctly.
	// No time to dig into it now, keep a match to avoid.
	duplicateFilter := map[string]struct{}{}
	isDuplicate := func(id string) bool {
		if _, ok := duplicateFilter[id]; !ok {
			duplicateFilter[id] = struct{}{}
			return false
		}
		return true
	}

	for reader.Next() {
		i, shp := reader.Shape()

		if !utils.Contains([]string{
			"1 - managed for biodiversity - disturbance events proceed or are mimicked",
			"2 - managed for biodiversity - disturbance events suppressed",
		}, reader.ReadAttribute(i, attrs["d_GAP_Sts"])) {
			filterCounter["d_GAP_Sts"]++
			continue
		}

		area := store.WildernessArea{
			PublicAccess: reader.ReadAttribute(i, attrs["d_Access"]),
		}
		if area.PublicAccess == "Closed" {
			filterCounter["d_Access"]++
			continue
		}

		area.Acres, err = strconv.ParseFloat(reader.ReadAttribute(i, attrs["GIS_Acres"]), 64)
		if err != nil {
			return nil, errors.Wrap(err, "could not parse GIS_Acres")
		}
		if area.Acres < 200 {
			filterCounter["GIS_Acres"]++
			continue
		}

		area.ManagementDesignation = reader.ReadAttribute(i, attrs["d_Des_Tp"])

		if !utils.Contains(
			[]string{
				"Area of Critical Environmental Concern",
				"Conservation Area",
				"Conservation Easement",
				"Local Conservation Area",
				"National Forest",
				"National Grassland",
				"National Lakeshore or Seashore",
				"National Park",
				"National Public Lands",
				"National Recreation Area",
				"National Scenic, Botanical or Volcanic Area",
				"National Wildlife Refuge",
				"State Conservation Area",
				"Recreation Management Area",
				"State Park",
				"State Wilderness",
				"State Recreation Area",
				"State Resource Management Area",
				"Watershed Protection Area",
				"Wild and Scenic River",
				"Wilderness Area",
				"Wilderness Study Area"},
			area.ManagementDesignation,
		) {
			filterCounter["d_Des_Tp"]++
			continue
		}

		area.ManagerType = reader.ReadAttribute(i, attrs["d_Mang_Typ"])
		area.ManagerName = reader.ReadAttribute(i, attrs["d_Mang_Nam"])

		if !utils.Contains(
			[]string{"Federal", "Joint", "Local Government", "Non-Governmental Organization", "State"},
			area.ManagerType,
		) {
			filterCounter["d_Mang_Typ"]++
			continue
		}

		area.State = store.WildernessAreaState(reader.ReadAttribute(i, attrs["d_State_Nm"]))

		if _, ok := mapStateAreas[area.State]; !ok {
			continue
		}

		id := reader.ReadAttribute(i, attrs["WDPA_Cd"])
		if id == "" {
			filterCounter["WDPA_Cd"]++
			continue
		}

		if isDuplicate(id) {
			continue
		}

		area.ID = id

		area.Name = reader.ReadAttribute(i, attrs["Loc_Nm"])
		area.Category = reader.ReadAttribute(i, attrs["d_Category"])
		area.OwnerType = reader.ReadAttribute(i, attrs["d_Own_Type"])
		area.OwnerName = reader.ReadAttribute(i, attrs["d_Own_Name"])


		area.YearEstablished, err = strconv.Atoi(reader.ReadAttribute(i, attrs["Date_Est"]))
		if err != nil && !strings.Contains(err.Error(), strconv.ErrSyntax.Error()) {
			return nil, errors.Wrap(err, "could not parse year established")
		}

		bbox := shp.BBox()
		sw := geo.NewPoint(bbox.MinY, bbox.MinX)
		centre := sw.MidpointTo(geo.NewPoint(bbox.MaxY, bbox.MaxX))
		area.Centre = latlng.LatLng{centre.Lat(), centre.Lng()}
		area.RadiusKilometers = sw.GreatCircleDistance(centre) / 2

		mapStateAreas[area.State] = append(mapStateAreas[area.State], area)
	}

	res := store.WildernessAreas{}

	for k := range mapStateAreas {
		sort.Sort(mapStateAreas[k])
		lastIndex := len(mapStateAreas[k])
		if lastIndex > 20 {
			lastIndex = 20
		}
		res = append(res, mapStateAreas[k][:lastIndex]...)
	}

	return res, nil
}

var stateMap = map[string]string{
	//"AL": "Alabama",
	//"AK": "Alaska",
	//"AZ": "Arizona",
	//"AR": "Arkansas",
	//"CA": "California",
	//"CO": "Colorado",
	//"CT": "Connecticut",
	"DE": "Delaware",
	//"FL": "Florida",
	"GA": "Georgia",
	//"HI": "Hawaii",
	//"ID": "Idaho",
	//"IL": "Illinois",
	//"IN": "Indiana",
	//"IA": "Iowa",
	//"KS": "Kansas",
	"KY": "Kentucky",
	//"LA": "Louisiana",
	//"ME": "Maine",
	"MD": "Maryland",
	//"MA": "Massachusetts",
	//"MI": "Michigan",
	//"MN": "Minnesota",
	//"MS": "Mississippi",
	//"MO": "Missouri",
	//"MT": "Montana",
	//"NE": "Nebraska",
	//"NV": "Nevada",
	//"NH": "New Hampshire",
	"NJ": "New Jersey",
	//"NM": "New Mexico",
	"NY": "New York",
	"NC": "North Carolina",
	//"ND": "North Dakota",
	"OH": "Ohio",
	//"OK": "Oklahoma",
	//"OR": "Oregon",
	"PA": "Pennsylvania",
	//"RI": "Rhode Island",
	"SC": "South Carolina",
	//"SD": "South Dakota",
	"TN": "Tennessee",
	//"TX": "Texas",
	//"UT": "Utah",
	//"VT": "Vermont",
	"VA": "Virginia",
	//"WA": "Washington",
	"WV": "West Virginia",
	//"WI": "Wisconsin",
	//"WY": "Wyoming",
	// Territories
	//"AS": "American Samoa",
	//"DC": "District of Columbia",
	//"FM": "Federated States of Micronesia",
	//"GU": "Guam",
	//"MH": "Marshall Islands",
	//"MP": "Northern Mariana Islands",
	//"PW": "Palau",
	//"PR": "Puerto Rico",
	//"VI": "Virgin Islands",
	// Armed Forces (AE includes Europe, Africa, Canada, and the Middle East)
	//"AA": "Armed Forces Americas",
	//"AE": "Armed Forces Europe",
	//"AP": "Armed Forces Pacific",
}
