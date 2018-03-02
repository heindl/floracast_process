package mushroomobserver

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/terra/geo"
	"bitbucket.org/heindl/process/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"gopkg.in/tomb.v2"
	"strconv"
	"strings"
	"sync"
	"time"
)

var fetchLmtr = utils.NewLimiter(5) // Should keep at five concurrently, because the API can not handle many requests.

func FetchOccurrences(cxt context.Context, targetID datasources.TargetID, since *time.Time) ([]*Observation, error) {

	if !targetID.Valid(datasources.TypeMushroomObserver) {
		return nil, errors.New("Invalid TargetID")
	}

	initRes := ObservationsResult{}
	initURL := occurrenceURL(targetID, since, 1)
	releaseOuterLmtr := fetchLmtr.Go()
	if err := utils.RequestJSON(initURL, &initRes); err != nil {
		releaseOuterLmtr()
		return nil, errors.Wrapf(err, "Could not fetch MushroomObserver Observations [%s]", initURL)
	}
	releaseOuterLmtr()

	observations := initRes.Results

	if initRes.NumberOfPages > 1 {
		lock := sync.Mutex{}
		tmb := tomb.Tomb{}
		tmb.Go(func() error {
			for 𝝨 := 2; 𝝨 <= initRes.NumberOfPages; 𝝨++ {
				releaseInnerLmtr := fetchLmtr.Go()
				i := 𝝨
				tmb.Go(func() error {
					defer releaseInnerLmtr()
					localRes := ObservationsResult{}
					localURL := occurrenceURL(targetID, since, i)
					if err := utils.RequestJSON(localURL, &localRes); err != nil {
						return errors.Wrapf(err, "Could not fetch MushroomObserver Observations [%s]", localURL)
					}
					lock.Lock()
					defer lock.Unlock()
					observations = append(observations, localRes.Results...)
					return nil
				})
			}
			return nil
		})
		if err := tmb.Wait(); err != nil {
			return nil, err
		}
	}

	if len(observations) != initRes.NumberOfRecords {
		return nil, errors.Newf("MushroomObserver fetched Observations [%d] are different than those expected [%d]", len(observations), initRes.NumberOfRecords)
	}

	res := []*Observation{}
	for _, observation := range observations {
		taxonID, err := targetID.ToInt()
		if err != nil {
			return nil, err
		}
		// Should be covered in search, but just in case.
		if observation.Consensus.ID != taxonID {
			return nil, errors.Newf("WARNING: MushroomObserver consensus id [%d] does not equal TaxonID [%d] in query.", observation.Consensus.ID, taxonID)
		}
		// Confidence should be covered by request, but just to be safe ...
		if observation.Confidence < 2 && observation.Namings.VotesForTaxonID(taxonID) < 2 {
			fmt.Println(fmt.Sprintf("WARNING: Insufficient Confidence [%f, %d]", observation.Confidence, observation.Namings.VotesForTaxonID(taxonID)))
			continue
		}
		res = append(res, observation)
	}

	return res, nil
}

func occurrenceURL(targetID datasources.TargetID, since *time.Time, page int) string {
	parameters := []string{
		fmt.Sprintf("name=%s", string(targetID)),
		"format=json",
		"detail=high",
		"has_images=true",
		"has_location=true",
		"is_collection_location=true",
		"east=-49.0",
		"north=83.3",
		"west=-178.2",
		"south=6.6",
		"confidence=2",
		fmt.Sprintf("page=%d", page),
	}

	if since != nil && !since.IsZero() {
		parameters = append(parameters, fmt.Sprintf("updated_at=%s-%s", since.Format("2006-01-02"), time.Now().Format("2006-01-02")))
	}

	return "http://mushroomobserver.org/api/observations?" + strings.Join(parameters, "&")
}

type ObservationsResult struct {
	Version         float64        `json:"version"`
	RunDate         time.Time      `json:"run_date"`
	Query           string         `json:"query"`
	NumberOfRecords int            `json:"number_of_records"`
	NumberOfPages   int            `json:"number_of_pages"`
	PageNumber      int            `json:"page_number"`
	Results         []*Observation `json:"results"`
	RunTime         float64        `json:"run_time"`
}

type Observation struct {
	ID                   int         `json:"id"`
	Type                 string      `json:"type"`
	Date                 string      `json:"date"`
	Latitude             CustomFloat `json:"latitude"`
	Longitude            CustomFloat `json:"longitude"`
	Altitude             CustomFloat `json:"altitude"`
	SpecimenAvailable    bool        `json:"specimen_available"`
	IsCollectionLocation bool        `json:"is_collection_location"`
	Confidence           float64     `json:"confidence"`
	NotesFields          struct {
	} `json:"notes_fields,omitempty"`
	Notes         string    `json:"notes,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	NumberOfViews int       `json:"number_of_views"`
	LastViewed    time.Time `json:"last_viewed"`
	Owner         struct {
		ID             int       `json:"id"`
		Type           string    `json:"type"`
		LoginName      string    `json:"login_name"`
		LegalName      string    `json:"legal_name"`
		Joined         time.Time `json:"joined"`
		Verified       time.Time `json:"verified"`
		LastLogin      time.Time `json:"last_login"`
		LastActivity   time.Time `json:"last_activity"`
		Contribution   int       `json:"contribution"`
		Notes          string    `json:"notes"`
		MailingAddress string    `json:"mailing_address"`
		LocationID     int       `json:"location_id"`
		ImageID        int       `json:"image_id"`
	} `json:"owner"`
	Consensus struct {
		ID            int       `json:"id"`
		Type          string    `json:"type"`
		Name          string    `json:"name"`
		Author        string    `json:"author"`
		Rank          string    `json:"rank"`
		Deprecated    bool      `json:"deprecated"`
		Misspelled    bool      `json:"misspelled"`
		Citation      string    `json:"citation"`
		Notes         string    `json:"notes"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
		NumberOfViews int       `json:"number_of_views"`
		LastViewed    time.Time `json:"last_viewed"`
		OkForExport   bool      `json:"ok_for_export"`
		SynonymID     int       `json:"synonym_id"`
	} `json:"consensus"`
	Location          Location      `json:"location"`
	CollectionNumbers []interface{} `json:"collection_numbers"`
	HerbariumRecords  []interface{} `json:"herbarium_records"`
	Sequences         []interface{} `json:"sequences"`
	Namings           Namings       `json:"namings"`
	PrimaryImage      struct {
		ID              int         `json:"id"`
		Type            string      `json:"type"`
		Date            string      `json:"date"`
		CopyrightHolder string      `json:"copyright_holder"`
		Notes           string      `json:"notes"`
		Quality         interface{} `json:"quality"`
		CreatedAt       time.Time   `json:"created_at"`
		UpdatedAt       time.Time   `json:"updated_at"`
		NumberOfViews   int         `json:"number_of_views"`
		LastViewed      time.Time   `json:"last_viewed"`
		OkForExport     bool        `json:"ok_for_export"`
		LicenseID       int         `json:"license_id"`
		OwnerID         int         `json:"owner_id"`
	} `json:"primary_image"`
	Images []struct {
		ID              int         `json:"id"`
		Type            string      `json:"type"`
		Date            string      `json:"date"`
		CopyrightHolder string      `json:"copyright_holder"`
		Notes           string      `json:"notes"`
		Quality         interface{} `json:"quality"`
		CreatedAt       time.Time   `json:"created_at"`
		UpdatedAt       time.Time   `json:"updated_at"`
		NumberOfViews   int         `json:"number_of_views"`
		LastViewed      time.Time   `json:"last_viewed"`
		OkForExport     bool        `json:"ok_for_export"`
		LicenseID       int         `json:"license_id"`
		OwnerID         int         `json:"owner_id"`
	} `json:"images"`
	Comments []struct {
		ID         int         `json:"id"`
		Type       string      `json:"type"`
		Summary    string      `json:"summary"`
		Content    string      `json:"content"`
		CreatedAt  time.Time   `json:"created_at"`
		UpdatedAt  interface{} `json:"updated_at"`
		OwnerID    int         `json:"owner_id"`
		ObjectType string      `json:"object_type"`
		ObjectID   int         `json:"object_id"`
	} `json:"comments"`
}

type Location struct {
	ID              int         `json:"id"`
	Type            string      `json:"type"`
	Name            string      `json:"name"`
	LatitudeNorth   float64     `json:"latitude_north"`
	LatitudeSouth   float64     `json:"latitude_south"`
	LongitudeEast   float64     `json:"longitude_east"`
	LongitudeWest   float64     `json:"longitude_west"`
	AltitudeMaximum interface{} `json:"altitude_maximum"`
	AltitudeMinimum interface{} `json:"altitude_minimum"`
	Notes           string      `json:"notes"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	NumberOfViews   int         `json:"number_of_views"`
	LastViewed      time.Time   `json:"last_viewed"`
	OkForExport     bool        `json:"ok_for_export"`
}

var ErrDistanceToLarge = errors.New("Distance too large")

func (Ω *Location) Coordinates() (lat, lng float64, err error) {
	p1, err := geo.NewPoint(Ω.LatitudeNorth, Ω.LongitudeEast)
	if err != nil {
		return 0, 0, err
	}
	p2, err := geo.NewPoint(Ω.LatitudeSouth, Ω.LongitudeWest)
	if err != nil {
		return 0, 0, err
	}
	distance := p1.DistanceKilometers(p2)
	if distance > 20 {
		return 0, 0, ErrDistanceToLarge
	}
	centroid, err := geo.Points{p1, p2}.Centroid()
	if err != nil {
		return 0, 0, err
	}
	return centroid.Latitude(), centroid.Longitude(), nil
}

type Naming struct {
	ID         int       `json:"id"`
	Type       string    `json:"type"`
	Confidence float64   `json:"confidence"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Name       struct {
		ID            int       `json:"id"`
		Type          string    `json:"type"`
		Name          string    `json:"name"`
		Author        string    `json:"author"`
		Rank          string    `json:"rank"`
		Deprecated    bool      `json:"deprecated"`
		Misspelled    bool      `json:"misspelled"`
		Citation      string    `json:"citation"`
		Notes         string    `json:"notes"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
		NumberOfViews int       `json:"number_of_views"`
		LastViewed    time.Time `json:"last_viewed"`
		OkForExport   bool      `json:"ok_for_export"`
		SynonymID     int       `json:"synonym_id"`
	} `json:"name"`
	OwnerID       int `json:"owner_id"`
	ObservationID int `json:"observation_id"`
	Votes         []struct {
		ID            int       `json:"id"`
		Type          string    `json:"type"`
		Confidence    float64   `json:"confidence"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
		NamingID      int       `json:"naming_id"`
		ObservationID int       `json:"observation_id"`
	} `json:"votes"`
	Reasons []interface{} `json:"reasons"`
}

func (Ω *Observation) Lat() (float64, error) {
	lat := float64(Ω.Latitude)
	if lat == 0 {
		var err error
		lat, _, err = Ω.Location.Coordinates()
		if err != nil {
			return 0, err
		}
	}
	return lat, nil
}
func (Ω *Observation) Lng() (float64, error) {
	lng := float64(Ω.Longitude)
	if lng == 0 {
		var err error
		_, lng, err = Ω.Location.Coordinates()
		if err != nil {
			return 0, err
		}
	}
	return lng, nil
}
func (Ω *Observation) DateString() string {
	return strings.Replace(Ω.Date, "-", "", -1)
}
func (Ω *Observation) CoordinatesEstimated() bool {
	isEstimated := false
	if Ω.Latitude == 0 || Ω.Longitude == 0 {
		isEstimated = true
	}
	return isEstimated
}
func (Ω *Observation) SourceOccurrenceID() string {
	return strconv.Itoa(Ω.ID)
}

type Namings []*Naming

func (Ω Namings) VotesForTaxonID(taxonID int) int {
	for _, naming := range Ω {
		if naming.Name.ID == taxonID {
			return len(naming.Votes)
		}
	}
	return 0
}

type CustomFloat float64

func (t *CustomFloat) MarshalJSON() ([]byte, error) {
	return json.Marshal(*t)
}

func (t *CustomFloat) UnmarshalJSON(b []byte) error {

	s := strings.Trim(string(b), `"`)

	if s == "null" {
		s = "0"
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*t = CustomFloat(f)

	return nil
}
