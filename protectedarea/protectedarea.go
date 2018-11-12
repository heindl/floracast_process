package protectedarea

import (
	"github.com/heindl/floracast_process/terra/geoembed"
	"github.com/heindl/floracast_process/utils"
	"github.com/saleswise/errors/errors"
	"strings"
)

// ProtectedArea is the point recommended to the user with predictions
type ProtectedArea interface {
	ID() (geoembed.S2Key, error)
	UpdateAccessLevel(int) error
	UpdateProtectionLevel(int) error
	UpdateOwner(string) error
	UpdateDesignation(string) error
	UpdateName(string) error
	Kilometers() float64
	Valid() bool
}

// NewProtectedArea creates one, and fails if the location is invalid.
func NewProtectedArea(lat, lng float64, squareKilometers float64) (ProtectedArea, error) {
	geoFeatures, err := geoembed.NewGeoFeatureSet(lat, lng, true)
	if err != nil {
		return nil, err
	}
	return &protectedArea{
		GeoFeatureSet:    geoFeatures,
		ProtectionLevel:  ProtectionLevelUnknown,
		AccessLevel:      AccessLevelUnknown,
		SquareKilometers: squareKilometers,
	}, nil
}

type protectedArea struct {
	Name             string                  `json:""`
	SquareKilometers float64                 `json:""` // Kilometers
	ProtectionLevel  ProtectionLevel         `json:""`
	Designation      string                  `json:""`
	Owner            string                  `json:""`
	AccessLevel      AccessLevel             `json:""`
	GeoFeatureSet    *geoembed.GeoFeatureSet `json:""`
}

func (Ω *protectedArea) ID() (geoembed.S2Key, error) {
	return geoembed.NewS2Key(Ω.GeoFeatureSet.Lat(), Ω.GeoFeatureSet.Lng())
}

func (Ω *protectedArea) UpdateName(name string) error {
	var err error
	Ω.Name, err = selectBetweenTitleSentences(Ω.Name, name)
	return err
}

func (Ω *protectedArea) UpdateDesignation(name string) error {
	var err error
	Ω.Designation, err = selectBetweenTitleSentences(Ω.Designation, name)
	return err
}

func (Ω *protectedArea) UpdateOwner(name string) error {
	var err error
	Ω.Owner, err = selectBetweenTitleSentences(Ω.Owner, name)
	return err
}

func (Ω *protectedArea) UpdateProtectionLevel(level int) error {

	if level < 1 || level > 4 {
		return errors.Newf("Invalid ProtectionLevel [%d]", level)
	}

	if ProtectionLevel(level) < Ω.ProtectionLevel {
		Ω.ProtectionLevel = ProtectionLevel(level)
	}

	return nil
}

func (Ω *protectedArea) UpdateAccessLevel(level int) error {

	if level < 1 || level > 4 {
		return errors.Newf("Invalid AccessLevel [%d]", level)
	}

	nlevel := AccessLevel(level)

	if Ω.AccessLevel == AccessLevelUnknown || (Ω.AccessLevel > nlevel && nlevel != AccessLevelUnknown) {
		Ω.AccessLevel = nlevel
	}

	return nil
}

func (Ω *protectedArea) Kilometers() float64 {
	return Ω.SquareKilometers
}

func (Ω *protectedArea) Valid() bool {

	if Ω.GeoFeatureSet == nil {
		return false
	}

	if Ω.Name == "" {
		return false
	}

	if Ω.AccessLevel == 0 {
		return false
	}

	if Ω.ProtectionLevel == 0 {
		return false
	}

	return true
}

// AccessLevel is the openness of a ProtectedArea to the public.
// Open (1), Restricted (2), Unknown (3), Closed (4)
// States have different reporting measures so these are all over the place.
type AccessLevel int

const (
	// AccessLevelOpen means the area is open to the public, though inconclusive
	AccessLevelOpen AccessLevel = iota + 1 // 1

	// AccessLevelRestricted means the area has restricted access, though inconclusive
	AccessLevelRestricted // 2

	// AccessLevelUnknown means we don't know, though inconclusive
	AccessLevelUnknown // 3

	// AccessLevelClosed means the area is closed, though inconclusive
	AccessLevelClosed // 4
)

// ProtectionLevel is the protection status (1-4) of an ecosystem.
// One is the highest level of protection, and four is unknown.
// States have different reporting measures so these are all over the place.
type ProtectionLevel int

const (
	// ProtectionLevelHighest [1]
	ProtectionLevelHighest ProtectionLevel = iota

	// ProtectionLevelHigh [2]
	ProtectionLevelHigh

	// ProtectionLevelMultipleUse  [3]
	ProtectionLevelMultipleUse

	// ProtectionLevelUnknown  [4]
	ProtectionLevelUnknown
)

//type ProtectedAreaState struct {
//	Name string `json:""`
//	Abbr string `json:""`
//}

//
//var validProtectedAreaStates = map[string]string{
//	"AL": "Alabama",
//	//"AK": "Alaska",
//	"AZ": "Arizona",
//	"AR": "Arkansas",
//	"CA": "California",
//	"CO": "Colorado",
//	"CT": "Connecticut",
//	"DE": "Delaware",
//	"FL": "Florida",
//	"GA": "Georgia",
//	//"HI": "Hawaii",
//	"ID": "Idaho",
//	"IL": "Illinois",
//	"IN": "Indiana",
//	"IA": "Iowa",
//	"KS": "Kansas",
//	"KY": "Kentucky",
//	"LA": "Louisiana",
//	"ME": "Maine",
//	"MD": "Maryland",
//	"MA": "Massachusetts",
//	"MI": "Michigan",
//	"MN": "Minnesota",
//	"MS": "Mississippi",
//	"MO": "Missouri",
//	"MT": "Montana",
//	"NE": "Nebraska",
//	"NV": "Nevada",
//	"NH": "New Hampshire",
//	"NJ": "New Jersey",
//	"NM": "New Mexico",
//	"NY": "New York",
//	"NC": "North Carolina",
//	"ND": "North Dakota",
//	"OH": "Ohio",
//	"OK": "Oklahoma",
//	"OR": "Oregon",
//	"PA": "Pennsylvania",
//	"RI": "Rhode Island",
//	"SC": "South Carolina",
//	"SD": "South Dakota",
//	"TN": "Tennessee",
//	"TX": "Texas",
//	"UT": "Utah",
//	"VT": "Vermont",
//	"VA": "Virginia",
//	"WA": "Washington",
//	"WV": "West Virginia",
//	"WI": "Wisconsin",
//	"WY": "Wyoming",
//	// Territories
//	//"AS": "American Samoa",
//	"DC": "District of Columbia",
//	//"FM": "Federated States of Micronesia",
//	//"GU": "Guam",
//	//"MH": "Marshall Islands",
//	//"MP": "Northern Mariana Islands",
//	//"PW": "Palau",
//	//"PR": "Puerto Rico",
//	//"VI": "Virgin Islands",
//	// Armed Forces (AE includes Europe, Africa, Canada, and the Middle East)
//	//"AA": "Armed Forces Americas",
//	//"AE": "Armed Forces Europe",
//	//"AP": "Armed Forces Pacific",
//}

func selectBetweenTitleSentences(a, b string) (string, error) {

	b, err := formatTitleWord(b)
	if err != nil {
		return "", err
	}

	wordFlags := []string{"unknown", "other", "easement", "private"}

	aHasFlag := utils.ContainsString(wordFlags, strings.ToLower(a))
	bHasFlag := utils.ContainsString(wordFlags, strings.ToLower(b))

	if a == "" || (aHasFlag && !bHasFlag && b != "") || len(a) < len(b) {
		return b, nil
	}

	return a, nil
}

func formatTitleWord(s string) (string, error) {

	s = strings.Replace(s, "_", " ", -1)
	s = strings.ToLower(strings.TrimSpace(s))
	for suffix, replacement := range map[string]string{
		"wa":  "Wilderness Area",
		"sp":  "State Park",
		"wma": "Wildlife Management Area",
	} {
		if strings.HasSuffix(s, suffix) {
			s = strings.Replace(s, suffix, replacement, -1)
		}
	}

	return utils.FormatTitle(s)
}
