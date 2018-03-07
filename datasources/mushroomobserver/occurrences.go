package mushroomobserver

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/datasources/providers"
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

func fetchObservationPage(targetID datasources.TargetID, since *time.Time, page int) ([]providers.Occurrence, int, error) {

	releaseLmtr := fetchLmtr.Go()
	defer releaseLmtr()

	url := occurrenceURL(targetID, since, page)

	res := observationsResult{}
	if err := utils.RequestJSON(url, &res); err != nil {
		return nil, 0, errors.Wrapf(err, "Could not fetch MushroomObserver Observations [%s]", url)
	}

	occurrenceList, err := filterObservations(targetID, res.Results)
	if err != nil {
		return nil, 0, err
	}
	return occurrenceList, res.NumberOfPages, nil

}

// FetchOccurrences implements the OccurrenceProvider interface.
func FetchOccurrences(_ context.Context, targetID datasources.TargetID, since *time.Time) ([]providers.Occurrence, error) {

	if !targetID.Valid(datasources.TypeMushroomObserver) {
		return nil, errors.New("Invalid TargetID")
	}

	res, numberOfPages, err := fetchObservationPage(targetID, since, 1)
	if err != nil {
		return nil, err
	}

	if numberOfPages > 1 {
		lock := sync.Mutex{}
		tmb := tomb.Tomb{}
		tmb.Go(func() error {
			for 𝝨 := 2; 𝝨 <= numberOfPages; 𝝨++ {
				page := 𝝨
				tmb.Go(func() error {
					list, _, err := fetchObservationPage(targetID, since, page)
					if err != nil {
						return err
					}
					lock.Lock()
					defer lock.Unlock()
					res = append(res, list...)
					return nil
				})
			}
			return nil
		})
		if err := tmb.Wait(); err != nil {
			return nil, err
		}
	}

	//if len(res) != initRes.NumberOfRecords {
	//	return nil, errors.Newf("MushroomObserver fetched Observations [%d] are different than those expected [%d]", len(res), initRes.NumberOfRecords)
	//}

	return res, nil
}

func filterObservations(targetID datasources.TargetID, given []*observation) ([]providers.Occurrence, error) {
	taxonID, err := targetID.ToInt()
	if err != nil {
		return nil, err
	}

	res := []providers.Occurrence{}
	for _, observation := range given {
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

type observationsResult struct {
	Version         float64        `json:"version"`
	RunDate         time.Time      `json:"run_date"`
	Query           string         `json:"query"`
	NumberOfRecords int            `json:"number_of_records"`
	NumberOfPages   int            `json:"number_of_pages"`
	PageNumber      int            `json:"page_number"`
	Results         []*observation `json:"results"`
	RunTime         float64        `json:"run_time"`
}

type observation struct {
	UpdatedAt         time.Time `json:"updated_at"`
	Type              string    `json:"type"`
	SpecimenAvailable bool      `json:"specimen_available"`
	//Sequences            []interface{} `json:"sequences"`
	PrimaryImage  image `json:"primary_image"`
	Owner         owner `json:"owner"`
	NumberOfViews int   `json:"number_of_views"`
	//NotesFields          struct{}      `json:"notes_fields,omitempty"`
	Notes                string      `json:"notes,omitempty"`
	Namings              namings     `json:"namings"`
	Longitude            customFloat `json:"longitude"`
	Location             location    `json:"location"`
	Latitude             customFloat `json:"latitude"`
	LastViewed           time.Time   `json:"last_viewed"`
	IsCollectionLocation bool        `json:"is_collection_location"`
	Images               []*image    `json:"images"`
	ID                   int         `json:"id"`
	//HerbariumRecords     []interface{} `json:"herbarium_records"`
	Date       string     `json:"date"`
	CreatedAt  time.Time  `json:"created_at"`
	Consensus  consensus  `json:"consensus"`
	Confidence float64    `json:"confidence"`
	Comments   []*comment `json:"comments"`
	//CollectionNumbers    []interface{} `json:"collection_numbers"`
	Altitude customFloat `json:"altitude"`
}

type comment struct {
	ID        int       `json:"id"`
	Type      string    `json:"type"`
	Summary   string    `json:"summary"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	//UpdatedAt  interface{} `json:"updated_at"`
	OwnerID    int    `json:"owner_id"`
	ObjectType string `json:"object_type"`
	ObjectID   int    `json:"object_id"`
}

type owner struct {
	Verified       time.Time `json:"verified"`
	Type           string    `json:"type"`
	Notes          string    `json:"notes"`
	MailingAddress string    `json:"mailing_address"`
	LoginName      string    `json:"login_name"`
	LocationID     int       `json:"location_id"`
	LegalName      string    `json:"legal_name"`
	LastLogin      time.Time `json:"last_login"`
	LastActivity   time.Time `json:"last_activity"`
	Joined         time.Time `json:"joined"`
	ImageID        int       `json:"image_id"`
	ID             int       `json:"id"`
	Contribution   int       `json:"contribution"`
}

type consensus struct {
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
}

type image struct {
	CopyrightHolder string    `json:"copyright_holder"`
	CreatedAt       time.Time `json:"created_at"`
	Date            string    `json:"date"`
	LastViewed      time.Time `json:"last_viewed"`
	LicenseID       int       `json:"license_id"`
	Notes           string    `json:"notes"`
	NumberOfViews   int       `json:"number_of_views"`
	OkForExport     bool      `json:"ok_for_export"`
	OwnerID         int       `json:"owner_id"`
	//Quality         interface{} `json:"quality"`
	Type      string    `json:"type"`
	UpdatedAt time.Time `json:"updated_at"`
	ID        int       `json:"id"`
}

type location struct {
	//AltitudeMaximum interface{} `json:"altitude_maximum"`
	//AltitudeMinimum interface{} `json:"altitude_minimum"`
	CreatedAt     time.Time `json:"created_at"`
	ID            int       `json:"id"`
	LastViewed    time.Time `json:"last_viewed"`
	LatitudeNorth float64   `json:"latitude_north"`
	LatitudeSouth float64   `json:"latitude_south"`
	LongitudeEast float64   `json:"longitude_east"`
	LongitudeWest float64   `json:"longitude_west"`
	Name          string    `json:"name"`
	Notes         string    `json:"notes"`
	NumberOfViews int       `json:"number_of_views"`
	OkForExport   bool      `json:"ok_for_export"`
	Type          string    `json:"type"`
	UpdatedAt     time.Time `json:"updated_at"`
}

var errDistanceToLarge = errors.New("Distance too large")

func (Ω *location) Coordinates() (lat, lng float64, err error) {
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
		return 0, 0, errDistanceToLarge
	}
	centroid, err := geo.Points{p1, p2}.Centroid()
	if err != nil {
		return 0, 0, err
	}
	return centroid.Latitude(), centroid.Longitude(), nil
}

type naming struct {
	ID            int       `json:"id"`
	Type          string    `json:"type"`
	Confidence    float64   `json:"confidence"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Name          name      `json:"name"`
	OwnerID       int       `json:"owner_id"`
	ObservationID int       `json:"observation_id"`
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

type name struct {
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
}

func (Ω *observation) Lat() (float64, error) {
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
func (Ω *observation) Lng() (float64, error) {
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
func (Ω *observation) DateString() string {
	return strings.Replace(Ω.Date, "-", "", -1)
}
func (Ω *observation) CoordinatesEstimated() bool {
	isEstimated := false
	if Ω.Latitude == 0 || Ω.Longitude == 0 {
		isEstimated = true
	}
	return isEstimated
}
func (Ω *observation) SourceOccurrenceID() string {
	return strconv.Itoa(Ω.ID)
}

type namings []*naming

func (Ω namings) VotesForTaxonID(taxonID int) int {
	for _, naming := range Ω {
		if naming.Name.ID == taxonID {
			return len(naming.Votes)
		}
	}
	return 0
}

type customFloat float64

func (t *customFloat) MarshalJSON() ([]byte, error) {
	return json.Marshal(*t)
}

func (t *customFloat) UnmarshalJSON(b []byte) error {

	s := strings.Trim(string(b), `"`)

	if s == "null" {
		s = "0"
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*t = customFloat(f)

	return nil
}
