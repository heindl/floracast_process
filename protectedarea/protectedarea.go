package protectedarea

import (
	"bitbucket.org/heindl/process/store"
	"bitbucket.org/heindl/process/terra/geoembed"
	"bitbucket.org/heindl/process/utils"
	"context"
	"github.com/saleswise/errors/errors"
	"gopkg.in/tomb.v2"
	"math"
	"strings"
)

// ProtectedArea is the point recommended to the user with predictions
type ProtectedArea interface {
	ID() (geoembed.CoordinateKey, error)
	UpdateAccessLevel(int) error
	UpdateProtectionLevel(int) error
	UpdateOwner(string) error
	UpdateDesignation(string) error
	UpdateName(string) error
	Valid() bool
}

// NewProtectedArea creates one, and fails if the location is invalid.
func NewProtectedArea(lat, lng float64, squareKilometers float64) (ProtectedArea, error) {
	geofeatures, err := geoembed.NewGeoFeatureSet(lat, lng, true)
	if err != nil {
		return nil, err
	}
	return &protectedArea{
		GeoFeatureSet:    geofeatures,
		ProtectionLevel:  ProtectionLevelUnknown,
		AccessLevel:      AccessLevelUnknown,
		SquareKilometers: squareKilometers,
	}, nil
}

type protectedArea struct {
	Name             string                  `firestore:"" json:""`
	SquareKilometers float64                 `firestore:"" json:""` // Kilometers
	ProtectionLevel  ProtectionLevel         `firestore:"" json:""`
	Designation      string                  `firestore:"" json:""`
	Owner            string                  `firestore:"" json:""`
	AccessLevel      AccessLevel             `firestore:"" json:""`
	GeoFeatureSet    *geoembed.GeoFeatureSet `json:",omitempty"`
}

func (Ω *protectedArea) ID() (geoembed.CoordinateKey, error) {
	return geoembed.NewCoordinateKey(Ω.GeoFeatureSet.Lat(), Ω.GeoFeatureSet.Lat())
}

func selectBetweenWordsForUpdate(a, b string) string {

	wordFlags := []string{"unknown", "other", "easement", "private"}

	a = strings.TrimSpace(a)
	b = strings.TrimSpace(b)

	aHasFlag := utils.ContainsString(wordFlags, strings.ToLower(a))
	bHasFlag := utils.ContainsString(wordFlags, strings.ToLower(b))

	if a == b || b == "" || (!aHasFlag && bHasFlag && a != "") {
		return a
	}

	if a == "" || (aHasFlag && !bHasFlag) {
		return b
	}

	if len(a) < len(b) {
		return b
	}

	return a

}

func (Ω *protectedArea) UpdateName(name string) error {
	Ω.Name = selectBetweenWordsForUpdate(Ω.Name, name)
	return nil
}

func (Ω *protectedArea) UpdateDesignation(name string) error {
	Ω.Designation = selectBetweenWordsForUpdate(Ω.Designation, name)
	return nil
}

func (Ω *protectedArea) UpdateOwner(name string) error {
	Ω.Owner = selectBetweenWordsForUpdate(Ω.Owner, name)
	return nil
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

func (Ω *protectedArea) Valid() bool {

	if !Ω.Valid() {
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
	AccessLevelOpen       AccessLevel = iota + 1 // 1
	AccessLevelRestricted                        // 2
	AccessLevelUnknown                           // 3
	AccessLevelClosed                            // 4
)

// ProtectionLevel is the protection status (1-4) of an ecosystem.
// One is the highest level of protection, and four is unknown.
// States have different reporting measures so these are all over the place.
type ProtectionLevel int

const (
	ProtectionLevelHighest     ProtectionLevel = iota // 1
	ProtectionLevelHigh                               // 2
	ProtectionLevelMultipleUse                        // 3
	ProtectionLevelUnknown                            // 4
)

//type ProtectedAreaState struct {
//	Name string `json:""`
//	Abbr string `json:""`
//}

//
//type ProtectedAreas []ProtectedArea
//
//func (a ProtectedAreas) Len() int           { return len(a) }
//func (a ProtectedAreas) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
//func (a ProtectedAreas) Less(i, j int) bool { return a[i].GISAcres > a[j].GISAcres }

// ProtectedAreas is intended for bulk uploading.
type ProtectedAreas []ProtectedArea

func (Ω ProtectedAreas) batches(maxBatchSize float64) []ProtectedAreas {

	if len(Ω) == 0 {
		return nil
	}

	batchCount := math.Ceil(float64(len(Ω)) / maxBatchSize)

	res := []ProtectedAreas{}
	for i := 0.0; i <= batchCount-1; i++ {
		start := int(i * maxBatchSize)
		end := int(((i + 1) * maxBatchSize) - 1)
		if end > len(Ω) {
			end = len(Ω)
		}
		o := Ω[start:end]
		res = append(res, o)
	}

	return res
}

// Upload validates all ProtectedAreas and saves them to Firestore
func (Ω ProtectedAreas) Upload(cxt context.Context, florastore store.FloraStore) (int, error) {

	batches := Ω.batches(500)

	col, err := florastore.FirestoreCollection(store.CollectionProtectedAreas)
	if err != nil {
		return 0, err
	}

	// TODO: Make parallel
	tmb := tomb.Tomb{}
	invalidCount := 0
	tmb.Go(func() error {
		for _, _batch := range batches {
			batch := _batch
			tmb.Go(func() error {
				firestoreBatch := florastore.FirestoreBatch()
				for _, area := range batch {
					if !area.Valid() {
						invalidCount++
						continue
					}
					id, err := area.ID()
					if err != nil {
						return err
					}
					firestoreBatch = firestoreBatch.Set(col.Doc(string(id)), area)
				}
				if _, err := firestoreBatch.Commit(cxt); err != nil {
					return errors.Wrap(err, "could not commit firestore batch")
				}
				return nil
			})
		}
		return nil
	})

	if err := tmb.Wait(); err != nil {
		return 0, err
	}

	return len(Ω) - invalidCount, nil
}

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
