package protected_areas

import (
	"context"
	//"github.com/cenkalti/backoff"
	"bitbucket.org/heindl/process/store"
	"bitbucket.org/heindl/process/terra/geoembed"
	"github.com/saleswise/errors/errors"
	"gopkg.in/tomb.v2"
	"math"
)

type ProtectedArea struct {
	Name            string                  `firestore:"" json:""`
	State           string                  `firestore:"" json:""`
	Area            float64                 `firestore:"" json:""` // Kilometers
	ProtectionLevel *ProtectionLevel        `firestore:"" json:""`
	Designation     string                  `firestore:"" json:""`
	Owner           string                  `firestore:"" json:""`
	AccessLevel     *AccessLevel            `firestore:"" json:""`
	GeoFeatureSet   *geoembed.GeoFeatureSet `json:",omitempty"`
}

func (a *ProtectedArea) ID() (geoembed.CoordinateKey, error) {
	return geoembed.NewCoordinateKey(a.GeoFeatureSet.Lat(), a.GeoFeatureSet.Lat())
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

	if !Ω.Valid() {
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

var counter = 0

type ProtectedAreas []*ProtectedArea

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

func (Ω ProtectedAreas) Upload(cxt context.Context, florastore store.FloraStore) error {

	batches := Ω.batches(500)

	col, err := florastore.FirestoreCollection(store.CollectionProtectedAreas)
	if err != nil {
		return err
	}

	// TODO: Make parallel
	tmb := tomb.Tomb{}
	invalid_counter := 0
	tmb.Go(func() error {
		for _, _batch := range batches {
			batch := _batch
			tmb.Go(func() error {
				firestore_batch := florastore.FirestoreBatch()
				for _, area := range batch {
					if !area.Valid() {
						invalid_counter += 1
						continue
					}
					id, err := area.ID()
					if err != nil {
						return err
					}
					firestore_batch = firestore_batch.Set(col.Doc(string(id)), area)
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
	//	if err != nil && strings.ContainsString(err.Error(), "Internal error encountered") {
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
