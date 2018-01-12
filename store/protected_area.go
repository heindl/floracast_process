package store

import (
	"google.golang.org/genproto/googleapis/type/latlng"
	"github.com/saleswise/errors/errors"
	"context"
	"math"
	"bitbucket.org/heindl/taxa/utils"
	"time"
	"github.com/cenkalti/backoff"
	"strings"
	"fmt"
	"github.com/paulmach/go.geojson"
)

type ProtectedArea struct {
	ID                    string             `json:",omitempty"`
	State                 ProtectedAreaState `json:",omitempty"`
	Acres                 float64            `json:",omitempty"`
	Name                  string             `json:",omitempty"`
	Centre                latlng.LatLng      `json:",omitempty"`
	HeightMeters float64      `json:",omitempty"`
	WidthMeters float64 `json:",omitempty"`
	Bounds string `json:",omitempty"`
	ManagerType           string             `json:",omitempty"`
	ManagerName           string             `json:",omitempty"`
	ManagementDesignation string             `json:",omitempty"`
	OwnerType             string             `json:",omitempty"`
	OwnerName             string             `json:",omitempty"`
	Category              string             `json:",omitempty"`
	YearEstablished       int                `json:",omitempty"`
	PublicAccess          string             `json:",omitempty"`
	GapAnalysisProjectStatus string `json:",omitempty"`
}

func(Ω *ProtectedArea) Valid() (isValid bool, reason string) {

	if Ω.ID == "" {
		return false, "id"
	}

	if v, ok := ValidProtectedAreaStates[Ω.State.Abbr]; !ok || Ω.State.Name != v {
		return false, "state"
	}

	if Ω.Centre.Latitude == 0 || Ω.Centre.Longitude == 0 {
		return false, "centre"
	}

	if Ω.Acres < 50 {
		return false, "acres"
	}

	if !strings.Contains(Ω.GapAnalysisProjectStatus, "managed for biodiversity") {
		return false, "gap_analysis"
	}

	if strings.ToLower(Ω.PublicAccess) == "closed" {
		return false, "public_access"
	}

	if !utils.Contains([]string{
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
		Ω.ManagementDesignation,
	) {
		return false, "management_designation"
	}

	if !utils.Contains([]string{
		"Federal",
		"Joint",
		"Local Government",
		"Non-Governmental Organization",
		"State",
		},
		Ω.ManagerType,
	) {
		return false, "manager_type"
	}

	return true, ""
}

type ProtectedAreaState struct {
	Name string `json:""`
	Abbr string `json:""`
}

type ProtectedAreas []ProtectedArea

func (a ProtectedAreas) Len() int           { return len(a) }
func (a ProtectedAreas) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ProtectedAreas) Less(i, j int) bool { return a[i].Acres > a[j].Acres }

func (Ω *store) ReadProtectedAreaByID(cxt context.Context, id string) (*ProtectedArea, error) {
	doc, err := Ω.FirestoreClient.Collection(CollectionTypeProtectedAreas).Doc(id).Get(cxt)
	if err != nil {
		errors.Wrap(err, "could not get ProtectedArea")
	}
	w := ProtectedArea{}
	if err := doc.DataTo(&w); err != nil {
		return nil, errors.Wrap(err, "could not type cast ProtectedArea")
	}
	return &w, nil
}

func (Ω *store) ReadProtectedAreas(cxt context.Context) ([]ProtectedArea, error) {

	docs, err := Ω.FirestoreClient.Collection(CollectionTypeProtectedAreas).
		Documents(cxt).
		GetAll()

	if err != nil {
		return nil, errors.Wrap(err, "could not find wilderness area")
	}

	res := []ProtectedArea{}
	for _, d := range docs {
		w := ProtectedArea{}
		if err := d.DataTo(&w); err != nil {
			return nil, errors.Wrap(err, "could not type cast ProtectedArea")
		}
		res = append(res, w)
	}

	return res, nil
}

func (Ω *store) ReadProtectedAreaByLatLng(cxt context.Context, lat, lng float64) (*ProtectedArea, error) {

	// Validate
	if lat == 0 || lng == 0 {
		return nil, errors.New("invalid protected area id")
	}

	docs, err := Ω.FirestoreClient.Collection(CollectionTypeProtectedAreas).
		// TODO: Would be great to use a geo query here or at least an approximation.
		Where("Centre.Longitude", ">", math.Floor(lng)).
		Where("Centre.Longitude", "<=", math.Ceil(lng)).
		Documents(cxt).
		GetAll()

	if err != nil {
		return nil, errors.Wrap(err, "could not find wilderness area")
	}

	for _, d := range docs {
		w := ProtectedArea{}
		if err := d.DataTo(&w); err != nil {
			return nil, errors.Wrap(err, "could not type cast ProtectedArea")
		}
		if utils.CoordinatesEqual(lat, w.Centre.Latitude) && utils.CoordinatesEqual(lat, w.Centre.Latitude) {
			return &w, nil
		}
	}

	return nil, errors.Newf("no wilderness area found: [%f, %f]", lat, lng)
}

var counter = 0;

func (Ω *store) SetProtectedArea(cxt context.Context, wa ProtectedArea) error {

	counter = counter + 1;

	// Validate
	if wa.ID == "" {
		return errors.New("invalid wilderness area id")
	}

	bkf := backoff.NewExponentialBackOff()
	bkf.InitialInterval = time.Second * 1
	ticker := backoff.NewTicker(bkf)
	for _ = range ticker.C {
		_, err := Ω.FirestoreClient.Collection(CollectionTypeProtectedAreas).Doc(wa.ID).Set(cxt, wa)
		if err != nil && strings.Contains(err.Error(), "Internal error encountered") {
			fmt.Println("Internal error encountered", err)
			continue
		}
		if err != nil {
			ticker.Stop()
			return errors.Wrap(err, "could not set protected area")
		}
		ticker.Stop()
		break
	}

	return nil
}

func (Ω *store) SetProtectedAreaGeometry(cxt context.Context, areaID string, geoJSONGeometry geojson.Geometry) error {

	////if !geoJSONGeometry.IsMultiPolygon() {
	////	return errors.New("Unsupported geojson geometry type.")
	////}
	//
	//b, err := json.Marshal(geoJSONGeometry)
	//if err != nil {
	//	return errors.Wrap(err, "Could not marshal geojson multipolygon.")
	//}
	//
	////"Geometry": base64.StdEncoding.EncodeToString([]byte(geometry)),
	//if _, err := Ω.FirestoreClient.Collection(CollectionTypeProtectedAreas).Doc(areaID).Update(cxt, []firestore.Update{
	//	{Path:"MultiPolygon", Value: b},
	//}); err != nil {
	//	return errors.Wrapf(err, "could not update protected area [%s] geometry", areaID)
	//}
	return nil
}

var ValidProtectedAreaStates = map[string]string{
	"AL": "Alabama",
	//"AK": "Alaska",
	"AZ": "Arizona",
	"AR": "Arkansas",
	"CA": "California",
	"CO": "Colorado",
	"CT": "Connecticut",
	"DE": "Delaware",
	"FL": "Florida",
	"GA": "Georgia",
	//"HI": "Hawaii",
	"ID": "Idaho",
	"IL": "Illinois",
	"IN": "Indiana",
	"IA": "Iowa",
	"KS": "Kansas",
	"KY": "Kentucky",
	"LA": "Louisiana",
	"ME": "Maine",
	"MD": "Maryland",
	"MA": "Massachusetts",
	"MI": "Michigan",
	"MN": "Minnesota",
	"MS": "Mississippi",
	"MO": "Missouri",
	"MT": "Montana",
	"NE": "Nebraska",
	"NV": "Nevada",
	"NH": "New Hampshire",
	"NJ": "New Jersey",
	"NM": "New Mexico",
	"NY": "New York",
	"NC": "North Carolina",
	"ND": "North Dakota",
	"OH": "Ohio",
	"OK": "Oklahoma",
	"OR": "Oregon",
	"PA": "Pennsylvania",
	"RI": "Rhode Island",
	"SC": "South Carolina",
	"SD": "South Dakota",
	"TN": "Tennessee",
	"TX": "Texas",
	"UT": "Utah",
	"VT": "Vermont",
	"VA": "Virginia",
	"WA": "Washington",
	"WV": "West Virginia",
	"WI": "Wisconsin",
	"WY": "Wyoming",
	// Territories
	//"AS": "American Samoa",
	"DC": "District of Columbia",
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