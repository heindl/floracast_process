package mushroomobserver

import (
	"fmt"
	"strings"
	"bitbucket.org/heindl/taxa/occurrences"
	"context"
	"time"
	"github.com/dropbox/godropbox/errors"
	"bitbucket.org/heindl/taxa/utils"
	"bitbucket.org/heindl/taxa/terra"
	"strconv"
	"bitbucket.org/heindl/taxa/geofeatures"
	"github.com/mongodb/mongo-tools/common/json"
	"bitbucket.org/heindl/taxa/datasources"
)

func FetchOccurrences(cxt context.Context, targetID datasources.DataSourceTargetID, since *time.Time) (*occurrences.OccurrenceAggregation, error) {

	if !targetID.Valid(datasources.DataSourceTypeMushroomObserver) {
		return nil, errors.New("Invalid TargetID")
	}

	res := occurrences.OccurrenceAggregation{}

	page := 1
	for {

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

		u := "http://mushroomobserver.org/api/observations?" + strings.Join(parameters, "&")

		apiResult := ObservationsResult{}
		if err := utils.RequestJSON(u, &apiResult); err != nil {
			return nil, errors.Wrap(err, "could not fetch mushroom observer observations")
		}

		for _, observation := range apiResult.Results {

			o, err := parseOccurrenceFromObservation(targetID, observation)
			if err != nil {
				return nil, err
			}

			if err := res.AddOccurrence(o); err != nil && !utils.ContainsError(err, occurrences.ErrCollision) {
				return nil, err
			}
		}

		if page >= apiResult.NumberOfPages {
			break
		}
		page += 1
	}

	return &res, nil
}

func parseOccurrenceFromObservation(targetID datasources.DataSourceTargetID, observation *Observation) (*occurrences.Occurrence, error) {

	taxonID, err := targetID.ToInt()
	if err != nil {
		return nil, err
	}

	lat := float64(observation.Latitude)
	lng := float64(observation.Longitude)
	isEstimated := false
	if lat == 0 || lng == 0 {
		isEstimated = true
		lat, lng = observation.Location.Coordinates()
	}

	// Should be covered in search, but just in case.
	if observation.Consensus.ID != taxonID {
		return nil, errors.Newf("WARNING: MushroomObserver consensus id [%d] does not equal taxon id [%d] in query.", observation.Consensus.ID, taxonID)
	}

	// Confidence should be covered by request, but just to be safe ...
	if observation.Confidence < 2 || observation.Namings.VotesForTaxonID(taxonID) < 2 {
		return nil, nil
	}

	o, err := occurrences.NewOccurrence(datasources.DataSourceTypeMushroomObserver, targetID, strconv.Itoa(observation.ID))
	if err != nil {
		return nil, err
	}

	err = o.SetGeospatial(lat, lng, strings.Replace(observation.Date, "-", "", -1), isEstimated)
	if err != nil && utils.ContainsError(err, geofeatures.ErrInvalidCoordinate) {
		return nil, nil
	}
	if err != nil && utils.ContainsError(err, occurrences.ErrInvalidDate) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return o, nil
}

type ObservationsResult struct {
	Version         float64   `json:"version"`
	RunDate         time.Time `json:"run_date"`
	Query           string    `json:"query"`
	NumberOfRecords int       `json:"number_of_records"`
	NumberOfPages   int       `json:"number_of_pages"`
	PageNumber      int       `json:"page_number"`
	Results         []*Observation     `json:"results"`
	RunTime         float64   `json:"run_time"`
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
	Location Location `json:"location"`
	CollectionNumbers []interface{} `json:"collection_numbers"`
	HerbariumRecords  []interface{} `json:"herbarium_records"`
	Sequences         []interface{} `json:"sequences"`
	Namings           Namings `json:"namings"`
	PrimaryImage struct {
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

func (Ω *Location) Coordinates() (lat, lng float64) {
	p1 := terra.NewPoint(Ω.LatitudeNorth, Ω.LongitudeEast)
	if p1.IsZero() {
		return 0, 0
	}
	p2 := terra.NewPoint(Ω.LatitudeSouth, Ω.LongitudeWest)
	if p2.IsZero() {
		return 0, 0
	}
	distance := p1.DistanceKilometers(p2)
	if distance > 20 {
		return 0, 0
	}
	centroid := terra.Points{p1, p2}.Centroid()
	return centroid.Latitude(), centroid.Longitude()
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
	c := CustomFloat(f)
	t = &c

	return nil
}