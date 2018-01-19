package store

import (
	"context"
	"fmt"
	//"github.com/cenkalti/backoff"
	"bitbucket.org/heindl/taxa/utils"
	"github.com/saleswise/errors/errors"
	"strings"
	"gopkg.in/tomb.v2"
)

type ProtectedArea struct {
	Name            string           `firestore:"" json:""`
	State           string           `firestore:"" json:""`
	PolyLabel       []float64        `firestore:"" json:""` // Longitude, Latitude
	Area            float64          `firestore:"" json:""` // Kilometers
	ProtectionLevel *ProtectionLevel `firestore:"" json:""`
	Designation     string           `firestore:"" json:""`
	Owner           string           `firestore:"" json:""`
	AccessLevel     *AccessLevel     `firestore:"" json:""`
	GeoFeatures     *GeoFeatures     `firestore:"" json:""`
}

func NewProtectedAreaID(lat, lng float64) string {
	return strings.Replace(fmt.Sprintf("%.6f_%.6f", lat, lng), ".", "|", -1)
}

func (a *ProtectedArea) ID() string {
	return NewProtectedAreaID(a.PolyLabel[1], a.PolyLabel[0])
}
func (a *ProtectedArea) Lat() float64 {
	return a.PolyLabel[1]
}

func (a *ProtectedArea) Lng() float64 {
	return a.PolyLabel[0]
}

func (a *ProtectedArea) SetGeoFeatures(f *GeoFeatures) {
	a.GeoFeatures = f
}

type AccessLevel int

const (
	AccessLevelOpen       AccessLevel = iota // 0
	AccessLevelRestricted                    // 1
	AccessLevelUnknown                       // 2
	AccessLevelClosed                        // 3
)

type ProtectionLevel int

const (
	ProtectionLevelHighest     ProtectionLevel = iota // 0
	ProtectionLevelHigh                               // 1
	ProtectionLevelMultipleUse                        // 2
	ProtectionLevelUnknown                            // 3
)

func (Ω *ProtectedArea) Valid() bool {

	if _, ok := ValidProtectedAreaStates[Ω.State]; !ok {
		return false
	}

	if len(Ω.PolyLabel) == 0 || Ω.PolyLabel[0] == 0 || Ω.PolyLabel[1] == 0 {
		return false
	}

	if Ω.GeoFeatures == nil || !Ω.GeoFeatures.Valid() {
		return false
	}

	if Ω.Name == "" {
		return false
	}

	if Ω.AccessLevel == nil {
		return false
	}

	if Ω.ProtectionLevel == nil {
		return false
	}

	return true
}

type ProtectedAreaState struct {
	Name string `json:""`
	Abbr string `json:""`
}

//
//type ProtectedAreas []ProtectedArea
//
//func (a ProtectedAreas) Len() int           { return len(a) }
//func (a ProtectedAreas) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
//func (a ProtectedAreas) Less(i, j int) bool { return a[i].GISAcres > a[j].GISAcres }

func (Ω *store) ReadProtectedArea(cxt context.Context, lat, lng float64) (*ProtectedArea, error) {

	doc, err := Ω.FirestoreClient.Collection(CollectionTypeProtectedAreas).Doc(NewProtectedAreaID(lat, lng)).Get(cxt)
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

var counter = 0

func (Ω *store) SetProtectedAreas(cxt context.Context, areas ...*ProtectedArea) error {

	counter = counter + 1

	batches := [][]*ProtectedArea{}
	batch_count := (len(areas) + 500 - 1) / 500

	for i := 0; i < len(areas); i += batch_count {
		end := i + batch_count
		if end > len(areas) {
			end = len(areas)
		}
		batches = append(batches, areas[i:end])
	}

	// TODO: Make parallel
	tmb := tomb.Tomb{}
	invalid_counter := 0
	tmb.Go(func() error {
		for _, _batch := range batches {
			batch := _batch
			tmb.Go(func() error {
				locations := []PredictableLocation{}
				for _, a := range batch {
					locations = append(locations, a)
				}
				if err := Ω.GeoFeaturesProcessor.ProcessLocations(cxt, locations...); err != nil {
					return err
				}
				firestore_batch := Ω.FirestoreClient.Batch()
				for _, area := range batch {
					if !area.Valid() {
						invalid_counter += 1
						fmt.Println("invalid", invalid_counter, utils.JsonOrSpew(area))
						continue
					}
					firestore_batch = firestore_batch.Set(
						Ω.FirestoreClient.Collection(CollectionTypeProtectedAreas).Doc(area.ID()),
						area)
				}
				if _, err := firestore_batch.Commit(cxt); err != nil {
					return errors.Wrap(err, "could not commit firestore batch")
				}
				return nil
			})
		}
		return nil
	})

	return tmb.Wait()


	//// Validate
	//if wa.ID == "" {
	//	return errors.New("invalid wilderness area id")
	//}
	//
	//bkf := backoff.NewExponentialBackOff()
	//bkf.InitialInterval = time.Second * 1
	//ticker := backoff.NewTicker(bkf)
	//for _ = range ticker.C {
	//	_, err := Ω.FirestoreClient.Collection(CollectionTypeProtectedAreas).Doc(string(wa.ID)).Set(cxt, wa)
	//	if err != nil && strings.Contains(err.Error(), "Internal error encountered") {
	//		fmt.Println("Internal error encountered", err)
	//		continue
	//	}
	//	if err != nil {
	//		ticker.Stop()
	//		return errors.Wrap(err, "could not set protected area")
	//	}
	//	ticker.Stop()
	//	break
	//}
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
