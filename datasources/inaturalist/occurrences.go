package inaturalist

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/datasources/providers"
	"bitbucket.org/heindl/process/utils"
	"context"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"google.golang.org/api/iterator"
	"strconv"
	"strings"
	"time"
)

type occurrence struct {
	//Annotations                    []interface{}        `json:"annotations"`
	Application      application `json:"application"`
	CachedVotesTotal int         `json:"cached_votes_total"`
	Captive          bool        `json:"captive"`
	//Comments                       []interface{}        `json:"comments"`
	CommentsCount int `json:"comments_count"`
	//CommunityTaxonID               interface{}          `json:"community_taxon_id"`
	CreatedAt        string      `json:"created_at"`
	CreatedAtDetails timeDetails `json:"created_at_details"`
	CreatedTimeZone  string      `json:"created_time_zone"`
	Description      string      `json:"description"`
	//Faves                          []interface{}        `json:"faves"`
	FavesCount int `json:"faves_count"`
	//Flags                          []interface{}        `json:"flags"`
	Geojson *geoJSON `json:"geojson"`
	//Geoprivacy                     interface{}          `json:"geoprivacy"`
	ID                          int             `json:"id"`
	IDPlease                    bool            `json:"id_please"`
	Identifications             identifications `json:"identifications"`
	IdentificationsCount        int             `json:"identifications_count"`
	IdentificationsMostAgree    bool            `json:"identifications_most_agree"`
	IdentificationsMostDisagree bool            `json:"identifications_most_disagree"`
	IdentificationsSomeAgree    bool            `json:"identifications_some_agree"`
	LicenseCode                 string          `json:"license_code"`
	Location                    string          `json:"location"`
	MapScale                    int             `json:"map_scale"`
	Mappable                    bool            `json:"mappable"`
	//NonOwnerIds                    []interface{}        `json:"non_owner_ids"`
	NumIdentificationAgreements    int  `json:"num_identification_agreements"`
	NumIdentificationDisagreements int  `json:"num_identification_disagreements"`
	OauthApplicationID             int  `json:"oauth_application_id"`
	Obscured                       bool `json:"obscured"`
	//ObservationPhotos              []interface{}        `json:"observation_photos"`
	ObservedOn        string      `json:"observed_on"`
	ObservedOnDetails timeDetails `json:"observed_on_details"`
	ObservedOnString  string      `json:"observed_on_string"`
	ObservedTimeZone  string      `json:"observed_time_zone"`
	//Ofvs                           []interface{}        `json:"ofvs"`
	OutOfRange bool `json:"out_of_range"`
	//Outlinks                       []interface{}        `json:"outlinks"`
	OwnersIdentificationFromVision bool `json:"owners_identification_from_vision"`
	//Photos                         []interface{}        `json:"photos"`
	PlaceGuess                 string               `json:"place_guess"`
	PlaceIds                   []int                `json:"place_ids"`
	PositionalAccuracy         int                  `json:"positional_accuracy"`
	Preferences                preferences          `json:"preferences"`
	ProjectIds                 []int                `json:"project_ids"`
	ProjectIdsWithCuratorID    []int                `json:"project_ids_with_curator_id"`
	ProjectIdsWithoutCuratorID []int                `json:"project_ids_without_curator_id"`
	ProjectObservations        []projectObservation `json:"project_observations"`
	PublicPositionalAccuracy   int                  `json:"public_positional_accuracy"`
	QualityGrade               string               `json:"quality_grade"`
	//QualityMetrics                 []interface{}        `json:"quality_metrics"`
	ReviewedBy []int `json:"reviewed_by"`
	SiteID     int   `json:"site_id"`
	//Sounds                         []interface{}        `json:"sounds"`
	SpeciesGuess string `json:"species_guess"`
	//Tags                           []interface{}        `json:"tags"`
	Taxon          taxon     `json:"taxon"`
	TimeObservedAt time.Time `json:"time_observed_at"` // 2012-03-18T17:31:53-04:00
	TimeZoneOffset string    `json:"time_zone_offset"`
	URI            string    `json:"uri"`
	UUID           string    `json:"uuid"`
	UpdatedAt      string    `json:"updated_at"`
	User           user      `json:"user"`
	//Votes                          []interface{}        `json:"votes"`
}

type preferences struct {
	PrefersCommunityTaxon interface{} `json:"prefers_community_taxon"`
}

type application struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
	Icon string `json:"icon"`
}

type geoJSON struct {
	Coordinates []string `json:"coordinates"`
	Type        string   `json:"type"`
}

type timeDetails struct {
	Date  string `json:"date"`
	Week  int    `json:"week"`
	Month int    `json:"month"`
	Hour  int    `json:"hour"`
	Year  int    `json:"year"`
	Day   int    `json:"day"`
}

// Lat returns the latitude.
func (Ω *occurrence) Lat() (float64, error) {
	lat, err := strconv.ParseFloat(Ω.Geojson.Coordinates[1], 64)
	if err != nil {
		return 0, errors.Wrap(err, "Could not parse latitude from INaturalist occurrence")
	}
	return lat, nil
}

// Lng returns the longitude.
func (Ω *occurrence) Lng() (float64, error) {
	lng, err := strconv.ParseFloat(Ω.Geojson.Coordinates[0], 64)
	if err != nil {
		return 0, errors.Wrap(err, "Could not parse longitude from INaturalist occurrence")
	}
	return lng, nil
}

// DateString converts the date to YYYYMMDD.
func (Ω *occurrence) DateString() string {
	return strings.Replace(Ω.ObservedOnDetails.Date, "-", "", -1)
}
func (Ω *occurrence) CoordinatesEstimated() bool {
	// Rounded to 5 decimal place. Not what I expected.
	// isEstimated := s.Issues.hasIssue(ogbif.OCCURRENCE_ISSUE_COORDINATE_ROUNDED)
	return false
}

// SourceOccurrenceID upholds the occurrence interface.
func (Ω *occurrence) SourceOccurrenceID() string {
	return strconv.Itoa(Ω.ID)
}

type projectObservation struct {
	Preferences struct {
		AllowsCuratorCoordinateAccess bool `json:"allows_curator_coordinate_access"`
	} `json:"preferences"`
	UserID    interface{} `json:"user_id"`
	ID        int         `json:"id"`
	UUID      string      `json:"uuid"`
	Project   project     `json:"project"`
	ProjectID int         `json:"project_id"`
}

type project struct {
	ID                       int           `json:"id"`
	ProjectObservationFields []interface{} `json:"project_observation_fields"`
	Icon                     string        `json:"icon"`
	Description              string        `json:"description"`
	Location                 string        `json:"location"`
	Title                    string        `json:"title"`
	Slug                     string        `json:"slug"`
	Latitude                 string        `json:"latitude"`
	Longitude                string        `json:"longitude"`
}

type user struct {
	ID          int    `json:"id"`
	Login       string `json:"login"`
	Preferences struct {
	} `json:"preferences"`
	LoginAutocomplete    string        `json:"login_autocomplete"`
	Name                 string        `json:"name"`
	NameAutocomplete     string        `json:"name_autocomplete"`
	Icon                 string        `json:"icon"`
	ObservationsCount    int           `json:"observations_count"`
	IdentificationsCount int           `json:"identifications_count"`
	JournalPostsCount    int           `json:"journal_posts_count"`
	ActivityCount        int           `json:"activity_count"`
	Roles                []interface{} `json:"roles"`
	SiteID               interface{}   `json:"site_id"`
	IconURL              string        `json:"icon_url"`
}

type identifications []*identification

func (Ω identifications) LatestIdentification() time.Time {
	latest := time.Time{}
	for _, identification := range Ω {
		if identification.CreatedAt.After(latest) {
			latest = identification.CreatedAt
		}
	}
	return latest
}

type identification struct {
	Disagreement     interface{}   `json:"disagreement"`
	Flags            []interface{} `json:"flags"`
	CreatedAt        time.Time     `json:"created_at"`
	TaxonID          int           `json:"taxon_id"`
	Body             interface{}   `json:"body"`
	OwnObservation   bool          `json:"own_observation"`
	UUID             string        `json:"uuid"`
	TaxonChange      interface{}   `json:"taxon_change"`
	Vision           bool          `json:"vision"`
	Current          bool          `json:"current"`
	ID               int           `json:"id"`
	CreatedAtDetails struct {
		Date  string `json:"date"`
		Week  int    `json:"week"`
		Month int    `json:"month"`
		Hour  int    `json:"hour"`
		Year  int    `json:"year"`
		Day   int    `json:"day"`
	} `json:"created_at_details"`
	Category                   string      `json:"category"`
	User                       user        `json:"user"`
	PreviousObservationTaxonID interface{} `json:"previous_observation_taxon_id"`
	Taxon                      taxon       `json:"taxon"`
}

// FetchOccurrences returns returns a slice of OccurrenceProviders.
func FetchOccurrences(_ context.Context, targetID datasources.TargetID, since *time.Time) ([]providers.Occurrence, error) {

	if !taxonIDFromTargetID(targetID).Valid() {
		return nil, errors.New("Invalid taxonID")
	}

	res := []providers.Occurrence{}
	page := 1
	for {
		list, err := fetchOccurrences(page, taxonIDFromTargetID(targetID), since)
		if err != nil && err != iterator.Done {
			return nil, err
		}
		for _, o := range list {
			res = append(res, o)
		}
		if err != nil && err == iterator.Done {
			break
		}
		page++
	}

	return res, nil
}

var throttle = time.NewTicker(time.Second / 20)

func fetchOccurrences(page int, txnID taxonID, since *time.Time) ([]*occurrence, error) {
	var response struct {
		TotalResults int           `json:"total_results"`
		Page         int           `json:"page"`
		PerPage      int           `json:"per_page"`
		Results      []*occurrence `json:"results"`
	}

	u := "https://api.inaturalist.org/v1/observations?place_id=97394&quality_grade=research&captive=false&per_page=200&geoprivacy=open"
	u += fmt.Sprintf("&taxon_id=%d", txnID)
	u += fmt.Sprintf("&page=%d", page)
	if since != nil && !since.IsZero() {
		u += fmt.Sprintf("&updated_since=%s", since.Format("2006-01-02"))
	}

	<-throttle.C

	if err := utils.RequestJSON(u, &response); err != nil {
		return nil, err
	}

	res := []*occurrence{}

	for _, inatOccurrence := range response.Results {

		// https://www.inaturalist.org/pages/help#quality
		// All of these are covered in query, but just to be safe ...
		if inatOccurrence.QualityGrade != "research" ||
			inatOccurrence.Captive ||
			inatOccurrence.Taxon.ID != txnID { // Ignore descendents
			continue
		}

		if inatOccurrence.Geojson == nil {
			return nil, errors.Newf("Invalid Geojson [%d]. Should have been solved by changing privacy setting", inatOccurrence.ID)
		}

		res = append(res, inatOccurrence)
	}

	if (response.Page * response.PerPage) < response.TotalResults {
		return res, nil
	}

	return res, iterator.Done
}
